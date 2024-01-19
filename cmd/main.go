package main

import (
	"context"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/nilorg/nat/internal/module"
	_ "github.com/nilorg/nat/internal/module/conf"
	"github.com/nilorg/nat/pkg/config"
	"github.com/nilorg/nat/pkg/natx"
	"github.com/nilorg/nat/pkg/runtimex"
	"github.com/nilorg/nat/pkg/watch"
	"github.com/nilorg/pkg/zlog"
	"github.com/spf13/viper"
)

func init() {
	// 初始化线程数量
	runtime.GOMAXPROCS(runtime.NumCPU())
	zlog.InitForViper(viper.GetViper())
}

func main() {
	// 监控系统信号和创建 Context 现在一步搞定
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// 在收到信号的时候，会自动触发 ctx 的 Done ，这个 stop 是不再捕获注册的信号的意思，算是一种释放资源。
	defer stop()

	module.Init(ctx)

	ch := watch.GetChannel()
	go natx.AutoSet(ctx, ch)

	portForwards := viper.Get("port-forward").([]interface{})
	if portForwards == nil {
		zlog.WithSugared(ctx).Fatalln("port-forward is nil")
		return
	}

	for _, portForwardConf := range portForwards {
		values := portForwardConf.(map[string]interface{})
		typ := values["type"].(string)
		port := values["port"].(int)
		remoteDomain := values["remote_domain"].(string)
		remotePort := values["remote_port"].(int)
		timing := values["timing"].(int)
		portForward := &config.PortForward{
			Type:         typ,
			Port:         port,
			RemoteDomain: remoteDomain,
			RemotePort:   remotePort,
			Timing:       timing,
		}
		go watch.Watch(ctx, portForward)
	}

	<-ctx.Done()

	zlog.Sugared.Infof("Stopped monitoring.")
	zlog.Sync()
}

// 检查系统环境是否支持
func checkSystem(ctx context.Context) {
	if runtimex.IsWindows() {

	}
}
