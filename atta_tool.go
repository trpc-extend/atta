package atta

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/codec"
	"git.code.oa.com/trpc-go/trpc-go/errs"
	attaCodec "git.woa.com/trpc-extend/trpc-go/atta/codec"
)

// Status 指007系统使用的Status
type Status int

const (
	StatusSuccess   = 0 // StatusSuccess 成功
	StatusException = 1 // StatusException 异常
	StatusTimeout   = 2 // StatusTimeout 超时
)

const (
	LetterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	LetterIdxBits = 6                    // 6 bits to represent a letter index
	LetterIdxMask = 1<<LetterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	LetterIdxMax  = 63 / LetterIdxBits   // # of letter indices fitting in 63 bits
)

// RandomString inspired by
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandomString(n int) []byte {
	//src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), LetterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), LetterIdxMax
		}
		if idx := int(cache & LetterIdxMask); idx < len(LetterBytes) {
			b[i] = LetterBytes[idx]
			i--
		}
		cache >>= LetterIdxBits
		remain--
	}
	return b
}

// Json 任意类型转换成json类型
func Json(i interface{}) string {
	b, _ := json.Marshal(i)
	return string(b)
}

// GetDecodeReportData 根据rpcName获取对应解码后的上报数据
func GetDecodeReportData(ctx context.Context, rpcName string, req, rsp interface{}) (
	interface{}, interface{}, []string) {
	var dReq, dRsp interface{}
	var extend []string
	reportCodec := attaCodec.GetReportCodec(rpcName)
	if reportCodec != nil {
		dReq, dRsp, extend = reportCodec.ReportDecode(ctx, req, rsp)
	} else {
		dReq, dRsp = req, rsp
	}
	return dReq, dRsp, extend
}

// GetStatusAndRetCode 默认从err中获取错误
func GetStatusAndRetCode(err error) (Status, string) {
	var (
		status  Status
		retCode string
	)
	if err != nil {
		e, ok := err.(*errs.Error)
		if ok {
			if e.Code == errs.RetClientTimeout && e.Type == errs.ErrorTypeFramework {
				status = StatusTimeout //超时
			} else {
				status = StatusException //异常
			}
			if e.Desc != "" {
				retCode = fmt.Sprintf("%s_%d", e.Desc, e.Code)
			} else {
				retCode = strconv.Itoa(int(e.Code))
			}
		} else {
			// 兼容 业务没有使用框架error,上报固定错误值
			status = StatusException //异常
			retCode = "007_999"
		}
	} else {
		retCode = "0"
		status = StatusSuccess
	}
	return status, retCode
}

// ToString 将任意数据转换成字符串，对非基本类型通过json格式转换
func ToString(v interface{}) (dst string) {
	switch s := v.(type) {
	case bool:
		dst = strconv.FormatBool(s)
	case float64:
		dst = strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		dst = strconv.FormatFloat(float64(s), 'f', -1, 32)
	case int:
		dst = strconv.Itoa(s)
	case int64:
		dst = strconv.FormatInt(s, 10)
	case int32:
		dst = strconv.Itoa(int(s))
	case int16:
		dst = strconv.FormatInt(int64(s), 10)
	case int8:
		dst = strconv.FormatInt(int64(s), 10)
	case uint:
		dst = strconv.FormatInt(int64(s), 10)
	case uint64:
		dst = strconv.FormatInt(int64(s), 10)
	case uint32:
		dst = strconv.FormatInt(int64(s), 10)
	case uint16:
		dst =  strconv.FormatInt(int64(s), 10)
	case uint8:
		dst =  strconv.FormatInt(int64(s), 10)
	case error:
		dst =  s.Error()
	case string:
		dst = s
	case []byte:
		dst = string(s)
	default:
		dst = Json(v)
	}
	return dst
}

// PInterfaceName 获取被调接口名称
func PInterfaceName(m codec.Msg) string {
	if m.CalleeMethod() == "" {
		return m.ServerRPCName()
	}
	return m.CalleeMethod()
}

// LocalAddr 获取被调主机信息
func LocalAddr(m codec.Msg) string {
	if m.LocalAddr() == nil {
		return ""
	}
	return m.LocalAddr().String()
}

// LocalAddr 获取主调主机信息
func RemoteAddr(m codec.Msg) string {
	if m.RemoteAddr() == nil {
		return ""
	}
	return m.RemoteAddr().String()
}

// ContainerName 获取被调容器名称
func ContainerName(m codec.Msg) string {
	if m.CalleeContainerName() == "" {
		return trpc.GlobalConfig().Global.ContainerName
	}
	return m.CalleeContainerName()
}
