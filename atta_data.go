package atta

const (
    DefReportName = "default" // 默认atta上报实现
)

const (
    TraceIDKey    = "atta_tid"         // 全链路追踪ID
    UIDKey        = "atta_uid"         // 业务用户ID
    ForbidKey     = "atta_forbid"      // 是否禁止atta上报标志
    ExtraFieldKey = "atta_extra_field" // 扩展额外字段
)

// ReportMsg 上报字段内容
type ReportMsg struct {
    TraceID    string   // 全链路跟踪ID
    Env        string   // 环境名
    UID        string   // 用户ID
    Status     string   // 状态
    Time       string   // 耗时，毫秒
    RetCode    string   // 业务返回码
    ErrMsg     string   // 错误信息，如有
    PModule    string   // 被调模块
    PInterface string   // 被调接口
    PHost      string   // 被调主机信息
    AModule    string   // 主调模块
    AInterface string   // 主调接口
    AHost      string   // 主调主机信息
    ReqBody    string   // 请求信息体
    RspBody    string   // 响应信息体
    ExtraField string   // 附加额外字段信息
    Extend     []string // 扩展字段信息
}

// ReportMsgToList 将ReportMsg消息转成列表格式
func ReportMsgToList(data *ReportMsg) []string {
    fieldValues := []string{
        data.TraceID,
        data.Env,
        data.UID,
        data.Status,
        data.Time,
        data.RetCode,
        data.ErrMsg,
        data.PModule,
        data.PInterface,
        data.PHost,
        data.AModule,
        data.AInterface,
        data.AHost,
        data.ReqBody,
        data.RspBody,
        data.ExtraField,
    }
    fieldValues = append(fieldValues, data.Extend...)
    return fieldValues
}

// ResetReportMsg 重置消息为默认值
func ResetReportMsg(data *ReportMsg) {
    data.TraceID = ""
    data.Env = ""
    data.UID = ""
    data.Status = ""
    data.Time = ""
    data.RetCode = ""
    data.ErrMsg = ""
    data.PModule = ""
    data.PInterface = ""
    data.PHost = ""
    data.AModule = ""
    data.AInterface = ""
    data.AHost = ""
    data.ReqBody = ""
    data.RspBody = ""
    data.ExtraField = ""
    data.Extend = make([]string, 0)
}