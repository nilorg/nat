package dnsx

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/spf13/viper"
)

var resolver *net.Resolver

func init() {
	var dnsServer = "8.8.8.8:53"

	dns := viper.GetString("dns")
	if dns != "" {
		dnsServer = dns
	}
	resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, address string) (net.Conn, error) {
			dialer := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return dialer.DialContext(ctx, "udp", dnsServer)
		},
	}
}

func LookupIPv4(ctx context.Context, domain string) ([]net.IP, error) {
	return resolver.LookupIP(ctx, "ip4", domain)
}

func LookupIPv6(ctx context.Context, domain string) ([]net.IP, error) {
	return resolver.LookupIP(ctx, "ip6", domain)
}

func PrintIPs(ips []net.IP) {
	for _, ip := range ips {
		fmt.Println(ip.String())
	}
}

func IpsEqual(ips1, ips2 []net.IP) bool {
	if len(ips1) != len(ips2) {
		return false
	}

	for i := range ips1 {
		if !ips1[i].Equal(ips2[i]) {
			return false
		}
	}

	return true
}
