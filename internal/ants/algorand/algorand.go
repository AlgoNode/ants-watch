package algorand

import (
	"context"
	"time"

	ds "github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/probe-lab/ants-watch/internal/ants/common"
	"github.com/probe-lab/ants-watch/internal/utils"
)

var logger = logging.Logger("ants-queen")

type AlgorandAnt struct {
	common.CommonAnt
	utils.Stoppable
}

func SpawnAlgorandAnt(ctx context.Context, ps peerstore.Peerstore, ds ds.Batching, cfg *common.AntConfig) (common.Ant, error) {
	commonAnt, err := common.SpawnAnt(ctx, ps, ds, cfg)
	if err != nil {
		return nil, err
	}
	algorandAnt := &AlgorandAnt{
		CommonAnt: *commonAnt,
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
			logger.Infow("Advertising gossip service", "service", advertiseList[i], "ant", a.CommonAnt.Host.ID())

			// Advertise with a TTL (time-to-live) - the advertisement will be valid for the specified duration
			ttl, err := routingDiscovery.Advertise(a.Context(), string(advertiseList[i]), discovery.TTL(time.Hour))
			if err != nil {
				logger.Errorf("advertise gossip service %s: %s", advertiseList[i], err)
				continue
			}

			logger.Infow("Successfully advertised service", "service", advertiseList[i], "ant", a.CommonAnt.Host.ID(), "ttl", ttl)

		}
		a.Sleep(time.Minute * 14)
	}

}
