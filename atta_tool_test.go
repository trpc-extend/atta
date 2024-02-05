package atta

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	mocker "github.com/tencent/goom"
	"github.com/trpc-extend/atta/codec"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/errs"
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

// ReportDecode 业务自定义扩展字段解析
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

// TestToString 单测 ToString
func TestToString(t *testing.T) {
	convey.Convey("TestToString", t, func() {
		tests := []struct {
			name string
			args interface{}
			want string
		}{
			{name: "TestToString bool: true", args: true, want: "true"},
			{name: "TestStr2Int64 float64: 3.1415", args: float64(3.1415), want: "3.1415"},
			{name: "TestStr2Int64 float32: 3.14", args: float32(3.14), want: "3.14"},
			{name: "TestStr2Int64 int: -1000", args: int(-1000), want: "-1000"},
			{name: "TestStr2Int64 int64: 123456789012345", args: int64(123456789012345), want: "123456789012345"},
			{name: "TestStr2Int64 int32: -123456789", args: int32(-123456789), want: "-123456789"},
			{name: "TestStr2Int64 int16: 1234", args: int16(1234), want: "1234"},
			{name: "TestStr2Int64 int8: 127", args: int8(127), want: "127"},
			{name: "TestStr2Int64 uint: 1000", args: uint(1000), want: "1000"},
			{name: "TestStr2Int64 uint64: 123456789012345", args: uint64(123456789012345), want: "123456789012345"},
			{name: "TestStr2Int64 uint32: 123456789", args: uint32(123456789), want: "123456789"},
			{name: "TestStr2Int64 uint16: 1234", args: uint16(1234), want: "1234"},
			{name: "TestStr2Int64 uint8: 127", args: uint8(127), want: "127"},
			{name: "TestStr2Int64 error: err", args: fmt.Errorf("mock exception"), want: "mock exception"},
			{name: "TestStr2Int64 string: abc123", args: "abc123", want: "abc123"},
			{name: "TestStr2Int64 []byte: abc456", args: []byte("abc456"), want: "abc456"},
			{name: "TestStr2Int64 nil: nil", args: nil, want: ""},
		}
		for _, tt := range tests {
			dst := ToString(tt.args)
			fmt.Printf("TestToString args:%+v, dst:%s, want:%s \n", tt.args, dst, tt.want)
			convey.So(dst, convey.ShouldEqual, tt.want)
		}
	})
}
