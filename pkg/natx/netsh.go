package natx

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/nilorg/pkg/zlog"
)

var _ Nater = &netshNAT{}

type netshNAT struct {
	exec.Cmd
}

func (n *netshNAT) Set(ctx context.Context, typ Type, localIPV IPV, localIP string, localPort int, remoteIPV IPV, remoteIP string, remotePort int) (err error) {
	zlog.WithSugared(ctx).Infoln("netshNAT Set.")
	if typ != TypeTCP && typ != TypeUDP {
		err = ErrNatSet
		return
	}
	if localIPV != IPv4 && localIPV != IPv6 {
		err = ErrNatSet
		return
	}
	// netsh interface portproxy set v4tov4 listenport=8080 listenaddress=127.0.0.1 connectport=8080 connectaddress=192.168.0.123
	_, err = exec.CommandContext(ctx, "netsh", "interface", "portproxy", "set", n.ipv(localIPV, remoteIPV), fmt.Sprintf("listenport=%d", localPort), fmt.Sprintf("listenaddress=%s", localIP), fmt.Sprintf("connectport=%d", remotePort), fmt.Sprintf("connectaddress=%s", remoteIP)).Output()
	if err != nil {
		zlog.WithSugared(ctx).Errorf("netshNAT Set error: %s", err)
		err = ErrNatSet
		return
	}
	return
}

func (n *netshNAT) Remove(ctx context.Context, typ Type, localIPV IPV, localIP string, localPort int, remoteIPV IPV, remoteIP string, remotePort int) (err error) {
	zlog.WithSugared(ctx).Infoln("netshNAT Remove.")
	if typ != TypeTCP && typ != TypeUDP {
		err = ErrNatRemove
		return
	}
	if localIPV != IPv4 && localIPV != IPv6 {
		err = ErrNatRemove
		return
	}
	// netsh interface portproxy delete v4tov4 listenport=8080 listenaddress=127.0.0.1
	_, err = exec.CommandContext(ctx, "netsh", "interface", "portproxy", "delete", n.ipv(localIPV, remoteIPV), fmt.Sprintf("listenport=%d", localPort), fmt.Sprintf("listenaddress=%s", localIP)).Output()
	if err != nil {
		zlog.WithSugared(ctx).Errorf("netshNAT Remove error: %s", err)
		err = ErrNatRemove
		return
	}
	return
}

func (n *netshNAT) Exist(ctx context.Context, typ Type, localIPV IPV, localIP string, localPort int, remoteIPV IPV, remoteIP string, remotePort int) (flag bool, err error) {
	zlog.WithSugared(ctx).Infoln("netshNAT Exist.")
	if typ != TypeTCP && typ != TypeUDP {
		err = ErrNatExist
		return
	}
	if localIPV != IPv4 && localIPV != IPv6 {
		err = ErrNatExist
		return
	}
	// netsh interface portproxy show v4tov4 listenport=8080 listenaddress=127.0.0.1
	_, err = exec.CommandContext(ctx, "netsh", "interface", "portproxy", "show", n.ipv(localIPV, remoteIPV), fmt.Sprintf("listenport=%d", localPort), fmt.Sprintf("listenaddress=%s", localIP), fmt.Sprintf("connectport=%d", remotePort), fmt.Sprintf("connectaddress=%s", remoteIP)).Output()
	if err != nil {
		zlog.WithSugared(ctx).Errorf("netshNAT Exist error: %s", err)
		err = ErrNatExist
		return
	}
	return
}

func (n *netshNAT) ipv(localIPV, remoteIPV IPV) string {
	// v4tov4/v4tov6/v6tov4/v6tov6
	if localIPV == IPv4 && remoteIPV == IPv4 {
		return "v4tov4"
	} else if localIPV == IPv4 && remoteIPV == IPv6 {
		return "v4tov6"
	} else if localIPV == IPv6 && remoteIPV == IPv4 {
		return "v6tov4"
	} else if localIPV == IPv6 && remoteIPV == IPv6 {
		return "v6tov6"
	}
	return "v4tov4"
}
