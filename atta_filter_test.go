package atta

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	mocker "github.com/tencent/goom"
	"trpc.group/trpc-go/trpc-go"
)

// CreateReq 创建群组请求
type CreateReq struct {
	Type    string `json:"Type"`
	GroupID string `json:"GroupId"`
	Name    string `json:"Name"`
}

// CreateRsp 创建群组回包
type CreateRsp struct {
	ActionStatus string `json:"ActionStatus"`
	ErrorInfo    string `json:"ErrorInfo"`
	ErrorCode    int    `json:"ErrorCode"`
	GroupID      string `json:"GroupId"`
}

// TrpcServerFilterFunc trpc框架server拦截器函数
func TrpcServerFilterFunc(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	return &CreateRsp{}, nil
}

// TrpcClientFilterFunc trpc框架client拦截器函数
func TrpcClientFilterFunc(ctx context.Context, req interface{}, rsp interface{}) (err error) {
	return nil
}

// TestReportServerFilter 单测ReportServerFilter
func TestReportServerFilter(t *testing.T) {
	convey.Convey("TestReportServerFilter", t, func() {
		convey.Convey("TestReportServerFilter run succ", func() {
			mocker := mocker.Create()
			defer mocker.Reset()
			RegisterDefault()
			var attaObj interface{}
			mocker.Struct(&attaObj).Method("SendFields").Return(0)
			req := &CreateReq{Type: "1", GroupID: "66881234"}
			rsp, err := GetDefaultReport().ReportServerFilter(trpc.BackgroundContext(), req, TrpcServerFilterFunc)
			convey.So(err, convey.ShouldBeNil)
			convey.So(rsp, convey.ShouldNotBeNil)
		})
	})
}

// TestReportClientFilter 单测ReportClientFilter
func TestReportClientFilter(t *testing.T) {
	convey.Convey("TestReportClientFilter", t, func() {
		convey.Convey("TestReportClientFilter run succ", func() {
			mocker := mocker.Create()
			defer mocker.Reset()
			RegisterDefault()
			var attaObj interface{}
			mocker.Struct(&attaObj).Method("SendFields").Return(0)
			req := &CreateReq{Type: "1", GroupID: "66881234"}
			rsp := &CreateRsp{ActionStatus: "1"}
			err := GetDefaultReport().ReportClientFilter(trpc.BackgroundContext(), req, rsp, TrpcClientFilterFunc)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
