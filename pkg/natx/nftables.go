package natx

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/nilorg/pkg/zlog"
)

const nftablesPrefix = `#!/usr/sbin/nft -f

delete table ip nat
add table ip nat
add chain nat PREROUTING { type nat hook prerouting priority -100 \; }
add chain nat POSTROUTING { type nat hook postrouting priority 100 \; }

delete table ip6 nat
add table ip6 nat
add chain nat PREROUTING { type nat hook prerouting priority -100 \; }
add chain nat POSTROUTING { type nat hook postrouting priority 100 \; }
`

// tcp
// add rule ip nat PREROUTING tcp dport {port1} counter dnat to {remoteIP}:{port2}
// add rule ip nat POSTROUTING ip daddr {remoteIP} tcp dport {port2} counter snat to {localIP}
// udp
// add rule ip nat PREROUTING udp dport {port1} counter dnat to {remoteIP}:{port2}
// add rule ip nat POSTROUTING ip daddr {remoteIP} udp dport {port2} counter snat to {localIP}

func InitNfTables() {

}

var _ Nater = &nftablesNAT{}

type nftablesNAT struct {
}

func (n *nftablesNAT) Set(ctx context.Context, typ Type, localIPV IPV, localIP string, localPort int, remoteIPV IPV, remoteIP string, remotePort int) (err error) {
	zlog.WithSugared(ctx).Infoln("nftablesNAT Set.")
	if typ != TypeTCP && typ != TypeUDP {
		err = ErrNatSet
		return
	}
	if localIPV != IPv4 && localIPV != IPv6 {
		err = ErrNatSet
		return
	}
	// 向 prerouting 链添加一条规则
	if localIPV == IPv6 {
		_, err = exec.CommandContext(ctx, "nft", "--", "add", "rule", "ip6", "nat", "prerouting", string(typ), "dport", fmt.Sprintf("%d", localPort), "counter", "dnat", "to", fmt.Sprintf("%s:%d", remoteIP, remotePort)).Output()
	} else {
		_, err = exec.CommandContext(ctx, "nft", "--", "add", "rule", "ip", "nat", "prerouting", string(typ), "dport", fmt.Sprintf("%d", localPort), "counter", "dnat", "to", fmt.Sprintf("%s:%d", remoteIP, remotePort)).Output()
	}
	if err != nil {
		zlog.WithSugared(ctx).Errorf("nftablesNAT Set error: %s", err)
		err = ErrNatSet
		return
	}
	// 向 prerouting 链添加一条规则
	if localIPV == IPv6 {
		_, err = exec.CommandContext(ctx, "nft", "--", "add", "rule", "ip6", "nat", "postrouting", "ip", "daddr", remoteIP, string(typ), "dport", fmt.Sprintf("%d", remotePort), "counter", "snat", "to", localIP).Output()
	} else {
		_, err = exec.CommandContext(ctx, "nft", "--", "add", "rule", "ip", "nat", "postrouting", "ip", "daddr", remoteIP, string(typ), "dport", fmt.Sprintf("%d", remotePort), "counter", "snat", "to", localIP).Output()
	}
	if err != nil {
		zlog.WithSugared(ctx).Errorf("nftablesNAT Set error: %s", err)
		err = ErrNatSet
		return
	}
	return
}

func (n *nftablesNAT) Remove(ctx context.Context, typ Type, localIPV IPV, localIP string, localPort int, remoteIPV IPV, remoteIP string, remotePort int) (err error) {
	zlog.WithSugared(ctx).Infoln("nftablesNAT Remove.")
	if typ != TypeTCP && typ != TypeUDP {
		err = ErrNatRemove
		return
	}
	if localIPV != IPv4 && localIPV != IPv6 {
		err = ErrNatRemove
		return
	}
	// 向 prerouting 链删除一条规则
	if localIPV == IPv6 {
		_, err = exec.CommandContext(ctx, "nft", "--", "delete", "rule", "ip6", "nat", "prerouting", string(typ), "dport", fmt.Sprintf("%d", localPort), "counter", "dnat", "to", fmt.Sprintf("%s:%d", remoteIP, remotePort)).Output()
	} else {
		_, err = exec.CommandContext(ctx, "nft", "--", "delete", "rule", "ip", "nat", "prerouting", string(typ), "dport", fmt.Sprintf("%d", localPort), "counter", "dnat", "to", fmt.Sprintf("%s:%d", remoteIP, remotePort)).Output()
	}
	if err != nil {
		zlog.WithSugared(ctx).Errorf("nftablesNAT Remove error: %s", err)
		err = ErrNatRemove
		return
	}
	// 向 prerouting 链删除一条规则
	if localIPV == IPv6 {
		_, err = exec.CommandContext(ctx, "nft", "--", "delete", "rule", "ip6", "nat", "postrouting", "ip", "daddr", remoteIP, string(typ), "dport", fmt.Sprintf("%d", remotePort), "counter", "snat", "to", localIP).Output()
	} else {
		_, err = exec.CommandContext(ctx, "nft", "--", "delete", "rule", "ip", "nat", "postrouting", "ip", "daddr", remoteIP, string(typ), "dport", fmt.Sprintf("%d", remotePort), "counter", "snat", "to", localIP).Output()
	}
	if err != nil {
		zlog.WithSugared(ctx).Errorf("nftablesNAT Remove error: %s", err)
		err = ErrNatRemove
		return
	}
	return
}

func (n *nftablesNAT) Exist(ctx context.Context, typ Type, localIPV IPV, localIP string, localPort int, remoteIPV IPV, remoteIP string, remotePort int) (flag bool, err error) {
	zlog.WithSugared(ctx).Infoln("nftablesNAT Exist.")

	return false, nil
}
