package algorand

import (
	"context"
	"fmt"
	"time"

	ds "github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	kad "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptcp "github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/probe-lab/ants-watch/internal/ants/common"
	"github.com/probe-lab/ants-watch/internal/keys"
	"github.com/probe-lab/ants-watch/internal/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var logger = logging.Logger("ants-queen")

type AlgorandAnt struct {
	common.CommonAnt
	utils.Stoppable
}

func SpawnAlgorandAnt(ctx context.Context, ps peerstore.Peerstore, ds ds.Batching, cfg *common.AntConfig) (common.Ant, error) {
	if cfg == nil {
		return nil, fmt.Errorf("no config given")
	} else if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Configure the resource manager to not limit anything
	noSubnetLimit := []rcmgr.ConnLimitPerSubnet{}
	noNetPrefixLimit := []rcmgr.NetworkPrefixLimit{}
	limiter := rcmgr.NewFixedLimiter(rcmgr.InfiniteLimits)
	rm, err := rcmgr.NewResourceManager(limiter,
		rcmgr.WithLimitPerSubnet(noSubnetLimit, noSubnetLimit),
		rcmgr.WithNetworkPrefixLimit(noNetPrefixLimit, noNetPrefixLimit),
	)
	if err != nil {
		return nil, fmt.Errorf("new resource manager: %w", err)
	}

	listenAddrs := []string{
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", cfg.Port),
		fmt.Sprintf("/ip6/::/tcp/%d", cfg.Port),
	}
	ymx := *yamux.DefaultTransport

	opts := []libp2p.Option{
		libp2p.UserAgent(cfg.UserAgent),
		libp2p.Identity(cfg.PrivateKey),
		libp2p.Peerstore(ps),
		libp2p.DisableRelay(),
		libp2p.Muxer("/yamux/1.0.0", &ymx),
		libp2p.ListenAddrStrings(listenAddrs...),
		libp2p.DisableMetrics(),
		libp2p.ShareTCPListener(),
		libp2p.ResourceManager(rm),
		libp2p.ConnectionManager(connmgr.NullConnMgr{}),
		libp2p.Transport(libp2ptcp.NewTCPTransport),
		libp2p.Security(noise.ID, noise.New),
		libp2p.AddrsFactory(addressFilter),
	}

	if cfg.Port == 0 {
		opts = append(opts, libp2p.NATPortMap()) // enable NAT port mapping if no port is specified
	}

	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("new libp2p host: %w", err)
	}
	h.SetStreamHandler(protocol.ID("/nodely-honeypot/1.0.0"), handleStreamBH)
	h.SetStreamHandler(protocol.ID(AlgorandWsProtocolV1), handleStreamBH)
	h.SetStreamHandler(protocol.ID(AlgorandWsProtocolV22), handleStreamBH)

	h.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, conn network.Conn) {
			cfg.Telemetry.ConnectCounter.Add(ctx, 1, metric.WithAttributes(
				attribute.String("direction", conn.Stat().Direction.String()),
			))
		},
		DisconnectedF: func(n network.Network, conn network.Conn) {
			cfg.Telemetry.DisconnectCounter.Add(ctx, 1, metric.WithAttributes(
				attribute.String("direction", conn.Stat().Direction.String()),
			))
		},
	})

	dhtOpts := []kad.Option{
		kad.Mode(kad.ModeServer),
		kad.BootstrapPeers(cfg.BootstrapPeers...),
		kad.V1ProtocolOverride(protocol.ID(cfg.ProtocolID)),
		kad.Datastore(ds),
		kad.OnRequestHook(common.OnRequestHook(h, cfg)),
	}
	dht, err := kad.New(ctx, h, dhtOpts...)
	if err != nil {
		return nil, fmt.Errorf("new libp2p dht: %w", err)
	}
	logger.Debugf("spawned ant. kadid: %s, peerid: %s", keys.PeerIDToKadID(h.ID()).HexString(), h.ID())

	if err = dht.Bootstrap(ctx); err != nil {
		logger.Warn("bootstrap failed: %s", err)
	}

	sub, err := h.EventBus().Subscribe([]interface{}{
		new(event.EvtLocalAddressesUpdated),
		new(event.EvtLocalReachabilityChanged),
	})
	if err != nil {
		return nil, fmt.Errorf("subscribe to event bus: %w", err)
	}

	go func() {
		for out := range sub.Out() {
			switch evt := out.(type) {
			case event.EvtLocalAddressesUpdated:
				if !evt.Diffs {
					continue
				}

				logger.Infow("Ant now listening on:", "ant", h.ID())
				for i, maddr := range evt.Current {
					actionStr := ""
					switch maddr.Action {
					case event.Added:
						actionStr = "ADD"
					case event.Removed:
						actionStr = "REMOVE"
					case event.Maintained:
						actionStr = "MAINTAINED"
					default:
						continue
					}
					logger.Infof("[%d] %s %s/p2p/%s", i, actionStr, maddr.Address, h.ID())
				}
			case event.EvtLocalReachabilityChanged:
				logger.Infow("Reachability changed", "ant", h.ID(), "reachability", evt.Reachability)
			}
		}
	}()

	algorandAnt := &AlgorandAnt{
		CommonAnt: common.CommonAnt{
			Cfg:   cfg,
			Host:  h,
			Dht:   dht,
			Sub:   sub,
			KadID: keys.PeerIDToKadID(h.ID()),
		},
		Stoppable: *utils.MakeStoppable(ctx),
	}
	go algorandAnt.AdvertiseLoop()

	return algorandAnt, nil
}

func (a *AlgorandAnt) AdvertiseLoop() {
	// Create a routing discovery using the DHT
	routingDiscovery := routing.NewRoutingDiscovery(a.CommonAnt.Dht)
	logger.Info("Advertising loop is up.")

	for !a.Stopped() {
		if !a.Sleep(time.Minute) {
			return
		}
		for i := range advertiseList {
			if !a.Sleep(time.Second) {
				return
			}
			// Advertise the "gossip" service
			logger.Infow("Advertising service", "service", advertiseList[i], "ant", a.CommonAnt.Host.ID())

			// Advertise with a TTL (time-to-live) - the advertisement will be valid for the specified duration
			ttl, err := routingDiscovery.Advertise(a.Context(), string(advertiseList[i]), discovery.TTL(time.Hour))
			if err != nil {
				logger.Errorf("advertise service %s: %s", advertiseList[i], err)
				continue
			}

			logger.Infow("Successfully advertised service", "service", advertiseList[i], "ant", a.CommonAnt.Host.ID(), "ttl", ttl)

		}
		a.Sleep(time.Minute * 14)
	}

}

func handleStreamBH(s network.Stream) {
	logger.Infof("Received new stream from peer: %s using protocol: %s\n", s.Conn().RemotePeer(), s.Protocol())

	// Read one message and close the stream
	buf := make([]byte, 1024)
	if n, _ := s.Read(buf); n >= 2 {
		logger.Warnf("Received payload", "Peer", s.Conn().RemotePeer(), "Legth", n, "MSG", string(buf[0:2]))
	}

	s.Close()
}
