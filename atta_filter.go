package atta

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/filter"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

// AttaServerFilter 作为主端被调函数接口拦截器
func AttaServerFilter(t *AttaReport) filter.ServerFilter {
	return func(ctx context.Context, req interface{}, nextHandle filter.ServerHandleFunc) (interface{}, error) {
		return t.ReportServerFilter(ctx, req, nextHandle)
	}
}

// AttaClientFilter 作为客户端请求后端服务拦截器
func AttaClientFilter(t *AttaReport) filter.ClientFilter {
	return func(ctx context.Context, req interface{}, rsp interface{}, nextHandle filter.ClientHandleFunc) error {
		return t.ReportClientFilter(ctx, req, rsp, nextHandle)
	}
}

// ReportServerFilter 框架上报atta日志服务端拦截器实现函数
func (t *AttaReport) ReportServerFilter(ctx context.Context, req interface{}, nextHandle filter.ServerHandleFunc) (
	interface{}, error) {
	startTime := time.Now()
	// 前置处理
	t.PreHandle(ctx, req)
	rsp, err := nextHandle(ctx, req)
	// 检查是否主动禁止上报
	if IsForbidReport(ctx) {
		log.Debugf("filter: req:%v, rsp:%v", Json(req), Json(rsp))
		return rsp, err
	}
	// 以下接口方法后置拦截在这里填充业务扩展字段信息
	msg := trpc.Message(ctx)
	dReq, dRsp, extend := GetDecodeReportData(ctx, msg.ServerRPCName(), req, rsp)
	status, retCode := GetStatusAndRetCode(err)
	data, _ := t.reportMsgPool.Get().(*ReportMsg)
	defer func() {
		ResetReportMsg(data)
		t.reportMsgPool.Put(data)
	}()
	// 填充上报字段数据
	data.TraceID = string(trpc.GetMetaData(ctx, TraceIDKey))
	data.Env = trpc.GlobalConfig().Global.EnvName
	data.UID = string(trpc.GetMetaData(ctx, UIDKey))
	data.Status = strconv.Itoa(int(status))
	data.RetCode = retCode
	data.ErrMsg = errs.Msg(err)
	data.PModule = fmt.Sprintf("%s.%s.%s", msg.CalleeApp(), msg.CalleeServer(), msg.CalleeService())
	data.PInterface = PInterfaceName(msg)
	data.PHost = fmt.Sprintf("%s.%s", LocalAddr(msg), ContainerName(msg))
	data.AModule = fmt.Sprintf("%s.%s.%s", msg.CallerApp(), msg.CallerServer(), msg.CallerService())
	data.AInterface = msg.CallerMethod()
	data.AHost = fmt.Sprintf("%s", RemoteAddr(msg))
	data.ReqBody = Json(dReq)
	data.RspBody = Json(dRsp)
	data.ExtraField = string(trpc.GetMetaData(ctx, ExtraFieldKey))
	data.Extend = extend
	data.Time = strconv.FormatInt(time.Since(startTime).Milliseconds(), 10)
	t.ReportFields(ReportMsgToList(data))
	ResetReportMsg(data)
	return rsp, err
}

// ReportClientFilter 框架上报atta日志请求端拦截器实现函数
func (t *AttaReport) ReportClientFilter(ctx context.Context, req, rsp interface{},
	nextHandle filter.ClientHandleFunc) error {
	startTime := time.Now()
	// 前置处理
	t.PreHandle(ctx, req)
	err := nextHandle(ctx, req, rsp)
	// 检查是否主动禁止上报
	if IsForbidReport(ctx) {
		log.Debugf("filter: req:%v, rsp:%v", Json(req), Json(rsp))
		return err
	}
	// 以下接口方法后置拦截在这里填充业务扩展字段信息
	msg := trpc.Message(ctx)
	dReq, dRsp, extend := GetDecodeReportData(ctx, msg.ClientRPCName(), req, rsp)
	status, retCode := GetStatusAndRetCode(err)
	data, _ := t.reportMsgPool.Get().(*ReportMsg)
	defer func() {
		ResetReportMsg(data)
		t.reportMsgPool.Put(data)
	}()
	// 填充上报字段数据
	data.TraceID = string(trpc.GetMetaData(ctx, TraceIDKey))
	data.Env = trpc.GlobalConfig().Global.EnvName
	data.UID = string(trpc.GetMetaData(ctx, UIDKey))
	data.Status = strconv.Itoa(int(status))
	data.RetCode = retCode
	data.ErrMsg = errs.Msg(err)
	data.PModule = fmt.Sprintf("%s.%s.%s", msg.CalleeApp(), msg.CalleeServer(), msg.CalleeService())
	data.PInterface = PInterfaceName(msg)
	data.PHost = fmt.Sprintf("%s.%s", RemoteAddr(msg), ContainerName(msg))
	data.AModule = fmt.Sprintf("%s.%s.%s", msg.CallerApp(), msg.CallerServer(), msg.CallerService())
	data.AInterface = msg.CallerMethod()
	data.AHost = fmt.Sprintf("%s.%s", LocalAddr(msg), trpc.GlobalConfig().Global.ContainerName)
	data.ReqBody = Json(dReq)
	data.RspBody = Json(dRsp)
	data.ExtraField = string(trpc.GetMetaData(ctx, ExtraFieldKey))
	data.Extend = extend
	data.Time = strconv.FormatInt(time.Since(startTime).Milliseconds(), 10)
	t.ReportFields(ReportMsgToList(data))
	return err
}

// preHandle 进行一些前置预处理，可以在这里生成TraceID等等逻辑
func (t *AttaReport) PreHandle(ctx context.Context, req interface{}) {
	// 前置检查 如果上游没有traceid字段，则当前请求为链路root，生成traceid
	if len(trpc.GetMetaData(ctx, TraceIDKey)) <= 0 {
		trpc.SetMetaData(ctx, TraceIDKey, RandomString(32))
	}
}
