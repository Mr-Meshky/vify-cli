package cleanip

import (
	"encoding/binary"
	"math/rand/v2"
	"net"
)

// CloudflareCIDRs contains standard IPv4 ranges assigned to Cloudflare edge network
var CloudflareCIDRs = []string{
	"104.16.0.0/13",
	"104.24.0.0/14",
	"172.64.0.0/13",
	"162.158.0.0/15",
	"198.41.128.0/17",
	"197.234.240.0/22",
	"188.114.96.0/20",
	"190.93.240.0/20",
	"141.101.64.0/18",
	"108.162.192.0/18",
	"103.21.244.0/22",
	"103.22.200.0/22",
	"103.31.4.0/22",
}

// PopularCandidateIPs are high-probability clean edge IPs frequently operational across MCI, Irancell, TCI
var PopularCandidateIPs = []string{
	"104.16.132.229",
	"104.16.133.229",
	"104.17.150.100",
	"104.18.2.161",
	"104.19.241.93",
	"104.20.74.39",
	"104.21.49.191",
	"104.22.6.21",
	"104.25.158.11",
	"104.26.12.18",
	"162.159.130.105",
	"162.159.138.85",
	"172.64.155.209",
	"172.67.181.189",
	"172.67.74.152",
	"188.114.97.7",
	"188.114.98.224",
	"198.41.200.12",
	"198.41.214.162",
}

// RandomIPsFromCIDRs samples up to `count` pseudo-random, unique host IPs
// from the given CIDR ranges. This is what makes the scanner actually
// "discover" fresh clean IPs rather than only re-ranking the same static
// PopularCandidateIPs list every run — important because an ISP can block
// those well-known IPs specifically while leaving the rest of the range open.
func RandomIPsFromCIDRs(cidrs []string, count int) []string {
	if count <= 0 || len(cidrs) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	out := make([]string, 0, count)

	for attempts := 0; len(out) < count && attempts < count*20; attempts++ {
		_, ipnet, err := net.ParseCIDR(cidrs[rand.IntN(len(cidrs))])
		if err != nil {
			continue
		}

		ip := randomHostIP(ipnet)
		if ip == "" || seen[ip] {
			continue
		}
		seen[ip] = true
		out = append(out, ip)
	}

	return out
}

func randomHostIP(ipnet *net.IPNet) string {
	ip4 := ipnet.IP.To4()
	if ip4 == nil {
		return "" // IPv4-only: Cloudflare edge sampling doesn't need IPv6 here
	}

	ones, bits := ipnet.Mask.Size()
	hostBits := bits - ones
	if hostBits <= 1 || hostBits > 24 {
		return "" // no usable host range, or absurdly large (keeps sampling meaningful)
	}

	base := binary.BigEndian.Uint32(ip4)
	maxHost := uint32(1)<<uint(hostBits) - 1

	// Avoid the network (.0) and broadcast address
	offset := uint32(1) + uint32(rand.IntN(int(maxHost-1)))
	addr := base + offset

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, addr)
	return net.IP(b).String()
}
