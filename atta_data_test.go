package atta

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"git.code.oa.com/goom/mocker"
)

// TestReportMsgToList 单测ReportMsgToList
func TestReportMsgToList(t *testing.T) {
	convey.Convey("TestResetReportMsg", t, func() {
		convey.Convey("TestResetReportMsg run succ", func() {
			mocker := mocker.Create()
			defer mocker.Reset()
			dataReport := &ReportMsg{
				ReqBody: "testReqBody",
				RspBody: "testRspBody",
				RetCode: "0",
				ErrMsg:  "succ",
			}
			result := ReportMsgToList(dataReport)
			convey.So(len(result), convey.ShouldEqual, 16)
			dataReport.Extend = []string{"1", "2", "3"}
			result = ReportMsgToList(dataReport)
			convey.So(len(result), convey.ShouldEqual, 19)
		})
	})
}

// TestResetReportMsg 单测ResetReportMsg
func TestResetReportMsg(t *testing.T) {
	convey.Convey("TestResetReportMsg", t, func() {
		convey.Convey("TestResetReportMsg run succ", func() {
			mocker := mocker.Create()
			defer mocker.Reset()
			dataReport := &ReportMsg{
				ReqBody: "testReqBody",
				RspBody: "testRspBody",
				RetCode: "0",
				ErrMsg:  "succ",
			}
			ResetReportMsg(dataReport)
		})
	})
}
