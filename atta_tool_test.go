package atta

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"git.code.oa.com/goom/mocker"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.woa.com/trpc-extend/trpc-go/atta/codec"
)

// TestRandomString 单测RandomString
func TestRandomString(t *testing.T) {
	convey.Convey("TestRandomString", t, func() {
		convey.Convey("TestRandomString run succ", func() {
			for i := 0; i < 3; i++ {
				resultByte := RandomString(32)
				fmt.Printf("%v\n", string(resultByte))
				convey.So(len(resultByte), convey.ShouldEqual, 32)
			}
		})
	})
}

// TestReportCodec 测试atta上报解码对象
type TestReportCodec struct {
}

// TVAttaDecodeFunc 长视频业务自定义扩展字段解析
func (t *TestReportCodec) ReportDecode(ctx context.Context, req, rsp interface{}) (interface{},
	interface{}, []string) {
	fieldList := make([]string, 3)
	fieldList[0] = "1"
	fieldList[1] = "2"
	fieldList[2] = "3"
	return req, rsp, fieldList
}

// TestGetDecodeReportData 单测GetDecodeReportData
func TestGetDecodeReportData(t *testing.T) {
	convey.Convey("TestGetDecodeReportData", t, func() {
		convey.Convey("TestGetDecodeReportData run succ", func() {
			mocker := mocker.Create()
			defer mocker.Reset()
			codec.Register("TestReportCodec", &TestReportCodec{})
			req, rsp, extend := GetDecodeReportData(trpc.BackgroundContext(),
				"TestReportCodec", "reqTest", 1000)
			convey.So(req, convey.ShouldEqual, "reqTest")
			convey.So(rsp, convey.ShouldEqual, 1000)
			convey.So(len(extend), convey.ShouldEqual, 3)
		})
	})
}

// TestGetStatusAndRetCode 单测GetStatusAndRetCode
func TestGetStatusAndRetCode(t *testing.T) {
	convey.Convey("TestGetStatusAndRetCode", t, func() {
		convey.Convey("TestGetStatusAndRetCode run succ", func() {
			mocker := mocker.Create()
			defer mocker.Reset()
			status, retCode := GetStatusAndRetCode(errs.New(1, "mock exception"))
			convey.So(status, convey.ShouldEqual, StatusException)
			convey.So(retCode, convey.ShouldEqual, "1")
			status, retCode = GetStatusAndRetCode(nil)
			convey.So(status, convey.ShouldEqual, StatusSuccess)
			convey.So(retCode, convey.ShouldEqual, "0")
		})
	})
}
