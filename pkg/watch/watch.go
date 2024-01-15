package watch

import (
	"context"
	"time"

	"github.com/nilorg/nat/pkg/config"
	"github.com/nilorg/nat/pkg/dnsx"
	"github.com/nilorg/pkg/zlog"
)

var (
	channel = make(chan *WatchChannel)
)

type WatchChannel struct {
	PortForward *config.PortForward
	Domain      *config.Domain
}

// Watch 检测DNS中解析到的IP是否有变化
func Watch(ctx context.Context, p *config.PortForward) {
	var lastIP string
	d, err := config.GetDomain(p.RemoteDomain)
	if err != nil {
		zlog.WithSugared(ctx).Errorf("获取域名信息失败: %s", err)
		return
	}
	if d != nil {
		lastIP = d.Ip
	}
	lastIPFlag := false
	if lastIP == "" {
		ips, err := dnsx.LookupIPv4(ctx, p.RemoteDomain)
		if err != nil {
			zlog.WithSugared(ctx).Errorf("获取域名IP失败: %s", err)
			return
		}
		if len(ips) == 0 {
			zlog.WithSugared(ctx).Errorln("获取域名IP失败: 无IP")
			return
		}
		lastIP = ips[0].String()
		domain := &config.Domain{
			Domain: p.RemoteDomain,
			Ip:     lastIP,
		}
		err = config.SetDomain(domain)
		if err != nil {
			zlog.WithSugared(ctx).Errorf("更新域名IP失败: %s", err)
			return
		}
		lastIPFlag = true
		channel <- &WatchChannel{
			PortForward: p,
			Domain:      domain,
		}
	}
	if lastIPFlag {
		timeSleep(p.Timing)
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ips, err := dnsx.LookupIPv4(ctx, p.RemoteDomain)
			if err != nil {
				zlog.WithSugared(ctx).Errorf("获取域名IP失败: %s", err)
				continue
			}
			if len(ips) == 0 {
				zlog.WithSugared(ctx).Errorln("获取域名IP失败: 无IP")
				continue
			}
			ip := ips[0].String()
			if ip != lastIP {
				zlog.WithSugared(ctx).Infof("域名IP变更: %s -> %s", lastIP, ip)
				lastIP = ip
				domain := &config.Domain{
					Domain: p.RemoteDomain,
					Ip:     lastIP,
				}
				err := config.SetDomain(domain)
				if err != nil {
					zlog.WithSugared(ctx).Errorf("更新域名IP失败: %s", err)
					continue
				}
				channel <- &WatchChannel{
					PortForward: p,
					Domain:      domain,
				}
			}
			timeSleep(p.Timing)
		}
	}
}

func timeSleep(timing int) {
	if timing <= 0 {
		timing = 30000
	}
	time.Sleep(time.Duration(timing) * time.Millisecond)
}

func GetChannel() <-chan *WatchChannel {
	return channel
}
