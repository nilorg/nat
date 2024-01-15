package natx

import (
	"context"

	"github.com/nilorg/nat/pkg/runtimex"
	"github.com/nilorg/nat/pkg/watch"
	"github.com/nilorg/pkg/zlog"
)

type Nater interface {
	Set(ctx context.Context, typ string, port int, remoteIp string, remotePort int) error
	Remove(ctx context.Context, typ string, port int) error
	Exist(ctx context.Context, typ string, port int) (bool, error)
}

// NewNAT 根据不同的系统配置端口转发
func NewNAT() Nater {
	if runtimex.IsWindows() {
		return &netshNAT{}
	} else {
		return &nftablesNAT{}
	}
}

// AutoSet 自动设置端口转发
func AutoSet(ctx context.Context, ch <-chan *watch.WatchChannel) {
	nat := NewNAT()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c := <-ch
			var exist bool
			var err error
			exist, err = nat.Exist(ctx, c.PortForward.Type, c.PortForward.Port)
			if err != nil {
				zlog.WithSugared(ctx).Errorf("检测端口转发失败: %s", err)
				continue
			}
			if exist {
				err = nat.Remove(ctx, c.PortForward.Type, c.PortForward.Port)
				if err != nil {
					zlog.WithSugared(ctx).Errorf("删除端口转发失败: %s", err)
					continue
				}
			}
			err = nat.Set(ctx, c.PortForward.Type, c.PortForward.Port, c.Domain.Ip, c.PortForward.RemotePort)
			if err != nil {
				zlog.WithSugared(ctx).Errorf("设置端口转发失败: %s", err)
			}
		}
	}
}
