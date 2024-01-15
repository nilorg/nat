package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func main() {
	ctx := context.Background()
	domain := "example.com"
	dnsServer := "8.8.8.8:53"
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, address string) (net.Conn, error) {
			dialer := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return dialer.DialContext(ctx, "udp", dnsServer)
		},
	}
	ips, err := resolver.LookupIP(ctx, "ip4", domain)
	if err != nil {
		fmt.Println("域名解析错误:", err)
		return
	}
	fmt.Println("初始解析结果:")
	printIPs(ips)

	time.Sleep(30 * time.Second)

	newIPs, err := resolver.LookupIP(ctx, "ip4", domain)
	if err != nil {
		fmt.Println("域名解析错误:", err)
		return
	}

	fmt.Println("30秒后的解析结果:")
	printIPs(newIPs)

	if !ipsEqual(ips, newIPs) {
		fmt.Println("域名解析发生了变化")
	} else {
		fmt.Println("域名解析未发生变化")
	}
}

func printIPs(ips []net.IP) {
	for _, ip := range ips {
		fmt.Println(ip.String())
	}
}

func ipsEqual(ips1, ips2 []net.IP) bool {
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
