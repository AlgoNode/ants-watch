package algorand

import (
	"net"

	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

var private6 = parseCIDR([]string{
	"100::/64",
	"2001:2::/48",
})

// parseCIDR converts string CIDRs to net.IPNet.
// function panics on errors so that it is only called during initialization.
func parseCIDR(cidrs []string) []*net.IPNet {
	result := make([]*net.IPNet, 0, len(cidrs))
	var ipnet *net.IPNet
	var err error
	for _, cidr := range cidrs {
		if _, ipnet, err = net.ParseCIDR(cidr); err != nil {
			panic(err)
		}
		result = append(result, ipnet)
	}
	return result
}

// addressFilter filters out private and unroutable addresses
func addressFilter(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {

	res := make([]multiaddr.Multiaddr, 0, len(addrs))
	for _, addr := range addrs {
		if manet.IsPublicAddr(addr) {
			if _, err := addr.ValueForProtocol(multiaddr.P_IP4); err == nil {
				// no rules for IPv4 at the moment, accept
				res = append(res, addr)
				continue
			}

			isPrivate := false
			a, err := addr.ValueForProtocol(multiaddr.P_IP6)
			if err != nil {
				logger.Warnf("failed to get IPv6 addr from %s: %v", addr, err)
				continue
			}
			addrIP := net.ParseIP(a)
			for _, ipnet := range private6 {
				if ipnet.Contains(addrIP) {
					isPrivate = true
					break
				}
			}
			if !isPrivate {
				res = append(res, addr)
			}
		}
	}
	return res
}
