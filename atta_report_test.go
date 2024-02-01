package atta

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	mocker "github.com/tencent/goom"

	"trpc.group/trpc-go/trpc-go"
)

// TestReportMsgToAtta 单测ReportMsgToAtta
func TestReportMsgToAtta(t *testing.T) {
	convey.Convey("TestReportMsgToAtta", t, func() {
		convey.Convey("TestReportMsgToAtta run succ", func() {
			mocker := mocker.Create()
			defer mocker.Reset()
			RegisterDefault()
			var attaObj interface{}
			mocker.Struct(&attaObj).Method("SendFields").Return(0)
			dataReport := &ReportMsg{
				ReqBody: "testReqBody",
				RspBody: "testRspBody",
				RetCode: "0",
				ErrMsg:  "succ",
			}
			GetDefaultReport().ReportMsgToAtta(trpc.BackgroundContext(), dataReport)
		})
	})
}

// TestAttaInterface 单测atta批量接口工具函数
func TestAttaInterface(t *testing.T) {
	convey.Convey("TestToolFunc", t, func() {
		convey.Convey("TestToolFunc run succ", func() {
			ctx := trpc.BackgroundContext()
			// SetTraceID
			str := string(RandomString(32))
			SetTraceID(ctx, str)
			convey.So(str, convey.ShouldEqual, GetTraceID(ctx))
			// SetTraceID
			str = "117551784"
			SetUID(ctx, str)
			convey.So(str, convey.ShouldEqual, GetUID(ctx))
			// SetForbidReport
			SetForbidReport(ctx)
			convey.So(IsForbidReport(ctx), convey.ShouldBeTrue)
			// ResetForbidReport
			ResetForbidReport(ctx)
			convey.So(IsForbidReport(ctx), convey.ShouldBeFalse)
			// SetExtraField
			str = "290541533"
			SetExtraField(ctx, str)
			convey.So(str, convey.ShouldEqual, GetExtraField(ctx))
			// AppendExtraValue
			AppendExtraValue(ctx, "117551784")
			convey.So(GetExtraField(ctx), convey.ShouldEqual, "290541533;117551784")
		})
	})
}
