package plugin

import (
	"math/rand"
	"time"

	"github.com/trpc-extend/trpc-go/atta"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginName = "rpc_atta"
	pluginType = "log"
)

func init() {
	plugin.Register(pluginName, &AttaReportPlugin{})
}

// AttaReportPlugin atta上报插件
type AttaReportPlugin struct {
}

// Type 插件类型名称
func (t *AttaReportPlugin) Type() string {
	return pluginType
}

// Setup 装载AttaReport拦截器实现
func (t *AttaReportPlugin) Setup(name string, decoder plugin.Decoder) error {
	rand.Seed(time.Now().UnixNano())
	cfg := atta.Config{}
	if err := decoder.Decode(&cfg); err != nil {
		log.Errorf("AttaReportPlugin decoder Decode Failed! name:%v, err:%v", name, err)
		return err
	}
	log.Infof("AttaReportPlugin Setup name:%v, config:%+v", name, cfg)
	atta.Register(name,
		atta.WithAttaInfo(cfg.AttaID, cfg.AttaToken),
		atta.WithRetryTime(cfg.RetryTime),
		atta.WithAutoEscape(cfg.AutoEscape),
	)
	if attaReport := atta.GetReport(name); attaReport != nil {
		filter.Register(name, atta.AttaServerFilter(attaReport), atta.AttaClientFilter(attaReport))
	} else {
		log.Errorf("AttaReportPlugin atta.GetReport Failed! name:%v", name)
	}
	return nil
}
