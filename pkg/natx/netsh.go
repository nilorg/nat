package natx

import (
	"context"

	"github.com/nilorg/pkg/zlog"
)

var _ Nater = &netshNAT{}

type netshNAT struct{}

func (n *netshNAT) Set(ctx context.Context, typ string, port int, remoteIp string, remotePort int) error {
	zlog.WithSugared(ctx).Infoln("netshNAT Set.")
	return nil
}

func (n *netshNAT) Remove(ctx context.Context, typ string, port int) error {
	zlog.WithSugared(ctx).Infoln("netshNAT Remove.")
	return nil
}

func (n *netshNAT) Exist(ctx context.Context, typ string, port int) (bool, error) {
	zlog.WithSugared(ctx).Infoln("netshNAT Exist.")
	return false, nil
}
