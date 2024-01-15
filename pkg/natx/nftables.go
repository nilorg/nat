package natx

import (
	"context"

	"github.com/nilorg/pkg/zlog"
)

var _ Nater = &nftablesNAT{}

type nftablesNAT struct {
}

func (n *nftablesNAT) Set(ctx context.Context, typ string, port int, remoteIp string, remotePort int) error {
	zlog.WithSugared(ctx).Infoln("nftablesNAT Set.")
	return nil
}

func (n *nftablesNAT) Remove(ctx context.Context, typ string, port int) error {
	zlog.WithSugared(ctx).Infoln("nftablesNAT Remove.")
	return nil
}

func (n *nftablesNAT) Exist(ctx context.Context, typ string, port int) (bool, error) {
	zlog.WithSugared(ctx).Infoln("nftablesNAT Exist.")
	return false, nil
}
