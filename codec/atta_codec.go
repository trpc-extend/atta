package codec

import (
	"context"
	"sync"
)

// ReportCodec atta上报编解码器
type ReportCodec interface {

	// ReportDecode 根据req rsp解码需要上报的数据
	// req: 请求协议包
	// req: 应答协议包
	// dReq: 解码后的需要上报请求协议数据
	// dRsp: 解码后的需要上报应答协议数据
	// extend: 自定义扩展数据，根据数组顺序依次对应扩展atta扩展字段
	ReportDecode(ctx context.Context, req, rsp interface{}) (dReq interface{}, dRsp interface{}, extend []string)
}

var (
	// 默认atta编解码器 可以做全局默认编解码上报
	defaultCodec ReportCodec = nil
	// key为trpc框架的 rpcName 例如：/trpc.video_detail.sport_national_group.NationalGroupService/GetIntroduction
	globalCodec = make(map[string]ReportCodec)
	rcLock      sync.RWMutex
)

// RegisterDefault 注册默认ReportCodec
func RegisterDefault(rpCodec ReportCodec) {
	defaultCodec = rpCodec
}

// Register 注册ReportCodec
// rpcName: 为trpc框架的rpcName 例如: /trpc.video_detail.sport_national_group.NationalGroupService/GetIntroduction
func Register(rpcName string, reportCodec ReportCodec) {
	rcLock.Lock()
	globalCodec[rpcName] = reportCodec
	rcLock.Unlock()
}

// GetReportCodec 获取对应的编解码器
func GetReportCodec(rpcName string) ReportCodec {
	rcLock.RLock()
	defer rcLock.RUnlock()
	reportCodec, ok := globalCodec[rpcName]
	if ok {
		return reportCodec
	}
	return defaultCodec
}
