package atta

import (
	"context"
	"fmt"
	"sync"

	"git.code.oa.com/atta/attaapi-go/v2"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/codec"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	attaCodec "git.woa.com/trpc-extend/trpc-go/atta/codec"
)

// AttaReport atta上报实例对象
type AttaReport struct {
	attaID        string          `yaml:"atta_id"`     // atta ID
	attaToken     string          `yaml:"atta_token"`  // atta token值
	retryTime     int             `yaml:"retry_time"`  // 上报atta失败重试次数
	autoEscape    bool            `yaml:"auto_escape"` // 是否打开自动转义

	attaObj       attaapi.AttaApi // atta上报实例
	reportMsgPool *sync.Pool      // 申请上报消息的消息池
}

// createAttaReport 创建AttaReport
func createAttaReport(name string, opts ...Option) *AttaReport {
	cfg := &Config{RetryTime: 1, AutoEscape: true}
	for _, opt := range opts {
		opt(cfg)
	}
	log.Infof("createAttaReport name:%v, cfg:%+v", name, cfg)
	// 初始化atta上报实例
	attaReport := &AttaReport{
		attaID:        cfg.AttaID,
		attaToken:     cfg.AttaToken,
		retryTime:     cfg.RetryTime,
		autoEscape:    cfg.AutoEscape,
		reportMsgPool: &sync.Pool{New: func() interface{} { return &ReportMsg{} }},
	}
	for i := 0; i < attaReport.retryTime; i++ {
		if initResult := attaReport.attaObj.InitUDP(); initResult != attaapi.AttaReportCodeSuccess {
			metrics.Counter("atta_init_udp_fail").Incr()
			log.Errorf("attaObj.InitUDP failed,initResult:%d,at retry time:%d", initResult, i)
		} else {
			log.Infof("attaObj.InitUDP succ")
			break
		}
	}
	return attaReport
}

// ReportMsgToAtta 主动上报 ReportMsg 消息到 atta
// ctx: 可通过ctx提取储通用协议头信息上报,
// data: 待上报的数据
func (t *AttaReport) ReportMsgToAtta(ctx context.Context, data *ReportMsg) {
	// 填充业务扩展字段信息
	reportCodec := attaCodec.GetReportCodec(trpc.Message(ctx).ServerRPCName())
	if len(data.Extend) <= 0 && reportCodec != nil {
		_, _, data.Extend = reportCodec.ReportDecode(ctx, nil, nil)
	}
	// 填充自定义字段信息
	data.Env = trpc.GlobalConfig().Global.EnvName
	msg := trpc.Message(ctx)
	metadata := msg.ServerMetaData()
	if metadata == nil {
		metadata = codec.MetaData{}
	}
	if data.TraceID == "" {
		data.TraceID = string(metadata[TraceIDKey])
	}
	if data.UID == "" {
		data.UID = string(metadata[UIDKey])
	}
	if data.ExtraField == "" {
		data.ExtraField = string(metadata[ExtraFieldKey])
	}
	// 设置被调信息
	if data.PInterface == "" {
		data.PInterface = PInterfaceName(msg)
	}
	data.PModule = fmt.Sprintf("%s.%s.%s", msg.CalleeApp(), msg.CalleeServer(), msg.CalleeService())
	data.PHost = fmt.Sprintf("%s.%s", RemoteAddr(msg), ContainerName(msg))
	// 设置主调信息
	if data.AInterface == "" {
		data.AInterface = msg.CallerMethod()
	}
	data.AModule = fmt.Sprintf("%s.%s.%s", msg.CallerApp(), msg.CallerServer(), msg.CallerService())
	data.AHost = fmt.Sprintf("%s.%s", LocalAddr(msg), trpc.GlobalConfig().Global.ContainerName)
	t.ReportFields(ReportMsgToList(data))
}

// ReportFields 封装对fieldValues数组消息的上报逻辑
func (t *AttaReport) ReportFields(fieldValues []string) {
	for i := 0; i < t.retryTime; i++ {
		ret := t.attaObj.SendFields(t.attaID, t.attaToken, fieldValues, t.autoEscape)
		if ret != attaapi.AttaReportCodeSuccess {
			metrics.Counter("atta_send_fields_fail").Incr()
			log.Errorf("SendFields failed! ret:%d, attaID:%s, fieldValues:%#v", ret, t.attaID, fieldValues)
		} else {
			log.Debugf("SendFields succ! attaID:%s, fieldValues:%#v", t.attaID, fieldValues)
			break
		}
	}
}

// ReportString 封装对字符串分割消息的上报，需要自己按照规定分隔符组装数据
func (t *AttaReport) ReportString(value string) {
	for i := 0; i < t.retryTime; i++ {
		ret := t.attaObj.SendString(t.attaID, t.attaToken, value)
		if ret != attaapi.AttaReportCodeSuccess {
			metrics.Counter("atta_send_string_fail").Incr()
			log.Errorf("SendString failed! ret:%d, attaID:%s, value:%s", ret, t.attaID, value)
		} else {
			log.Debugf("SendString succ! attaID:%s, value:%s", t.attaID, value)
			break
		}
	}
}
