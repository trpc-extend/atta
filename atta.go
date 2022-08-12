package atta

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/atta/attaapi-go/v2"
)

func init() {
	// 上报过程中用到随机种子因此这里seed下
	rand.Seed(time.Now().UnixNano())
	// 初始化创建默认的上报对象
	defaultReport = createAttaReport(DefReportName, WithAttaInfo("03900066415", "2818548356"))
}

var (
	defaultReport  *AttaReport
	attaReportPool = make(map[string]*AttaReport)
	attaReportLock sync.RWMutex
)

// RegisterDefault 注册默认AttaReport对象
func RegisterDefault(opts ...Option) {
	defaultReport = createAttaReport(DefReportName, opts...)
}

// GetDefaultReport 获取默认的AttaReport
func GetDefaultReport() *AttaReport {
	return defaultReport
}

// Register 以name为key注册AttaReport对象
func Register(name string, opts ...Option) {
	attaReport := createAttaReport(name, opts...)
	attaReportLock.Lock()
	attaReportPool[name] = attaReport
	attaReportLock.Unlock()
}

// GetReport 获取通过 Register 注册的AttaReport对象, 如果没有返回nil
func GetReport(name string) *AttaReport {
	attaReportLock.RLock()
	defer attaReportLock.RUnlock()
	if v, ok := attaReportPool[name]; ok {
		return v
	}
	return nil
}

// GetDefaultAttaApiObj 获取默认的AttaApi对象, 如果没有返回nil
func GetDefaultAttaApiObj() *attaapi.AttaApi {
	if defaultReport != nil {
		return &defaultReport.attaObj
	}
	return nil
}

// GetAttaApiObj 获取通过 Register 注册的AttaApi对象, 如果没有返回nil
func GetAttaApiObj(name string) *attaapi.AttaApi {
	attaReportLock.RLock()
	defer attaReportLock.RUnlock()
	if v, ok := attaReportPool[name]; ok {
		return &v.attaObj
	}
	return nil
}

// ReportMsgToAtta 上报ReportMsg数据到atta
func ReportMsgToAtta(ctx context.Context, data *ReportMsg) {
	GetDefaultReport().ReportMsgToAtta(ctx, data)
}

// SetTraceID 设置链路追踪ID
func SetTraceID(ctx context.Context, traceID string) {
	trpc.SetMetaData(ctx, TraceIDKey, []byte(traceID))
}

// GetTraceID 获取设置链路追踪ID
func GetTraceID(ctx context.Context) string {
	return string(trpc.GetMetaData(ctx, TraceIDKey))
}

// SetUID 设置业务自定义用户ID
func SetUID(ctx context.Context, id string) {
	trpc.SetMetaData(ctx, UIDKey, []byte(id))
}

// GetUID 获取设置业务自定义用户ID
func GetUID(ctx context.Context) string {
	return string(trpc.GetMetaData(ctx, UIDKey))
}

// SetForbidReport 设置禁止当前接口调用atta上报
func SetForbidReport(ctx context.Context) {
	trpc.SetMetaData(ctx, ForbidKey, []byte("true"))
}

// ResetForbidReport 重置禁止当前接口调用atta上报状态
func ResetForbidReport(ctx context.Context) {
	trpc.SetMetaData(ctx, ForbidKey, []byte(""))
}

// IsForbidReport 获取是否禁止上报
func IsForbidReport(ctx context.Context) bool {
	v := trpc.GetMetaData(ctx, ForbidKey)
	if string(v) == "true" {
		return true
	}
	return false
}

// SetExtraField 设置附加额外字段信息
func SetExtraField(ctx context.Context, value string) {
	trpc.SetMetaData(ctx, ExtraFieldKey, []byte(value))
}

// GetExtraField 获取设置附加额外字段信息
func GetExtraField(ctx context.Context) string {
	return string(trpc.GetMetaData(ctx, ExtraFieldKey))
}

// AppendExtraValue 以英文分号为分隔符追加扩展字段的值
func AppendExtraValue(ctx context.Context, value string) {
	AppendExtraValueSep(ctx, value, ";")
}

// AppendExtraValueSep 以sep为分隔符追加扩展字段的值
func AppendExtraValueSep(ctx context.Context, value string, sep string) {
	ov := string(trpc.GetMetaData(ctx, ExtraFieldKey))
	if len(ov) == 0 {
		trpc.SetMetaData(ctx, ExtraFieldKey, []byte(value))
	} else {
		trpc.SetMetaData(ctx, ExtraFieldKey, []byte(ov+sep+value))
	}
}


