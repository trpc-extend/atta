# 提供基于trpc框架的通用日志上报组件

支持被动插件方式上报和主动调用上报

## 上报字段

基于atta通用上报字段，支持在该字段后面进行扩展字段，根据如下字段申请对应的atta数据表，记录申请的attaID和attaToken
- TraceID,全链路ID,varchar,64,
- Env,环境名,varchar,32,
- UID,用户ID,varchar,64,
- Status,状态信息,varchar,64,
- Time,耗时毫秒,varchar,32,
- RetCode,业务返回码,varchar,32,
- ErrMsg,错误信息,varchar,1024,
- PModule,被调模块,varchar,128,
- PInterface,被调接口,varchar,128,
- PHost,被调主机信息,varchar,128,
- AModule,主调模块,varchar,128,
- AInterface,主调接口,varchar,128,
- AHost,主调主机信息,varchar,128,
- ReqBody,请求信息体,varchar,8192,
- RspBody,响应信息体,varchar,8192,
- ExtraField,附加字段信息,varchar,4096,
  
## 被动上报

被动上报是通过注册trpc插件实现的，因此分两个步骤：  
- 首先在main函数地方引入插件包
```
  import _ "github.com/trpc-extend/trpc-go/atta/plugin"
```

- 其次增加trpc-go框架配置：
```
server:     # 按照需要支持服务端配置
  filter:
    - rpc_atta # 插件名称
    
client:     # 按照需要支持客户端配置，默认下不配置，对性能消耗比较大
  filter:
    - rpc_atta # 插件名称

plugins:
 log: # plugin type
     rpc_atta: # plugin name
         atta_id: 03900066415 # 申请表 atta ID
         atta_token: 2818548356 # 申请表 atta token
         retry_time: 1 #上报atta失败重试次数 默认1次及不重试
         auto_escape: true #是否打开自动转义
```

## 主动上报

主动上报是使用者在程序代码中进行主动调用上报，主要依赖组件提供的API调用。
- 首先在调用地方引入插件包
```
  import "github.com/trpc-extend/trpc-go/atta"
```
- 在程序起始地方，注册对应的上报实例，支持多实例注册，以下选择其一即可：
```
  // 注册默认实例上报
  atta.RegisterDefault(atta.WithAttaInfo("03900066415", "2818548356"))
  // 注册以"log_report"为名称的自定义实例上报
  atta.Register("log_report", atta.WithAttaInfo("03900066415", "2818548356"))
``` 
  使用atta.RegisterDefault()注册的是默认上报实例，注册这个好处是在使用的时候代码简单，比如如下调用：
```
  reportMsg := &atta.ReportMsg{}
  atta.ReportMsgToAtta(ctx, reportMsg)
``` 
  否则调用方式如下：
```
  reportMsg := &atta.ReportMsg{}
  atta.GetReport("log_report").ReportMsgToAtta(ctx, reportMsg)
```  

## 常用接口

- Register 注册AttaReport对象
```
// name: 实例名称, opts: 可选注册参数
func Register(name string, opts ...Option)
```  
- RegisterDefault 注册默认AttaReport对象
```
// name: 实例名称, opts: 可选注册参数
func Register(name string, opts ...Option)
```
- GetReport 获取通过 Register 注册的AttaReport对象, 如果没有返回nil
```
// name: 获取实例名称
func GetReport(name string) *AttaReport 
```
- GetDefaultReport 获取默认的AttaReport
```
func GetDefaultReport() *AttaReport
```
- (t *AttaReport) ReportMsgToAtta 主动上报 ReportMsg消息到atta
```
// ctx: 可通过ctx提取储通用协议头信息上报, data: 上报的数据
func (t *AttaReport) ReportMsgToAtta(ctx context.Context, data *ReportMsg) 
```
- ReportMsgToAtta 对默认实例对象ReportMsgToAtta方法封装，方便使用
```
// ctx: 可通过ctx提取储通用协议头信息上报, data: 上报的数据
func ReportMsgToAtta(ctx context.Context, data *ReportMsg)
```
- SetTraceID 设置链路追踪ID
```
// ctx: 承载上报信息内容上下文, traceID: 链路追踪ID
func SetTraceID(ctx context.Context, traceID string)
``` 
- GetTraceID 获取设置链路追踪ID
```
// ctx: 承载上报信息内容上下文
func GetTraceID(ctx context.Context) string
```
- SetUID 设置业务自定义用户ID
```
// ctx: 承载上报信息内容上下文, id: 用户自定义ID
SetUID(ctx context.Context, id string)
``` 
- GetUID 获取设置业务自定义用户ID
```
// ctx: 承载上报信息内容上下文
func GetUID(ctx context.Context) string
```
- SetForbidReport 设置禁止当前接口调用atta上报
```
// ctx: 承载上报信息内容上下文
func SetForbidReport(ctx context.Context)
``` 
- ResetForbidReport 重置禁止当前接口调用atta上报状态
```
// ctx: 承载上报信息内容上下文
func ResetForbidReport(ctx context.Context)
``` 
- IsForbidReport 获取是否禁止上报
```
// ctx: 承载上报信息内容上下文
func IsForbidReport(ctx context.Context) bool
``` 
- SetExtraField 设置附加额外字段信息
```
// ctx: 承载上报信息内容上下文, value: 设置附加字段信息
func SetExtraField(ctx context.Context, value string)
``` 
- GetExtraField 获取设置附加额外字段信息
```
// ctx: 承载上报信息内容上下文
func GetExtraField(ctx context.Context) string
``` 
- AppendExtraValue 以英文分号为分隔符追加扩展字段的值
```
// ctx: 承载上报信息内容上下文 value: 追加附加字段信息
func AppendExtraValue(ctx context.Context, value string)
```
- AppendExtraValueSep 以sep为分隔符追加扩展字段的值
```
// ctx: 承载上报信息内容上下文 value: 追加附加字段信息 sep: 分隔符
func AppendExtraValueSep(ctx context.Context, value string, sep string)
```
- RegisterDefault(atta/codec包) 注册默认ReportCodec实例 实现自定义解码和扩展字段解析
```
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

// RegisterDefault 注册默认ReportCodec
func RegisterDefault(rpCodec ReportCodec)
```
- Register(atta/codec包) 以trpc框架rpcName为key注册ReportCodec实例 实现自定义解码和扩展字段解析
```
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

// Register 注册ReportCodec
// rpcName: 为trpc框架的rpcName 例如: /trpc.video_detail.sport_national_group.NationalGroupService/GetIntroduction
func Register(rpcName string, reportCodec ReportCodec) 
```
- GetReportCodec(atta/codec包) 获取对应的编解码器
```
// GetReportCodec 获取对应的编解码器
func GetReportCodec(rpcName string) ReportCodec
``` 

## 扩展字段上报

扩展上报是为了通用字段不能满足业务需求而设计，组件支持在通用字段的后面追加定义任意个数扩展字段上报，以下以视频实现的扩展字段为例介绍说明如何使用。

- 首先在通用字段表的后面，申请追加扩展字段信息，如下从ReqHead字段开始都是业务自定义的扩展字段
  - TraceID,全链路ID,varchar,64,
  - Env,环境名,varchar,32,
  - UID,用户ID,varchar,64,
  - Status,状态信息,varchar,64,
  - Time,耗时毫秒,varchar,32,
  - RetCode,业务返回码,varchar,32,
  - ErrMsg,错误信息,varchar,1024,
  - PModule,被调模块,varchar,128,
  - PInterface,被调接口,varchar,128,
  - PHost,被调主机信息,varchar,128,
  - AModule,主调模块,varchar,128,
  - AInterface,主调接口,varchar,128,
  - AHost,主调主机信息,varchar,128,
  - ReqBody,请求信息体,varchar,8192,
  - RspBody,响应信息体,varchar,8192,
  - ExtraField,附加字段信息,varchar,4096,
  - `ReqHead,请求消息头,varchar,4096,`
  - `QQ,用户QQ,varchar,64,`
  - `Vuid,用户Vuid,varchar,64,`
  - `OmgID,用户OmgID,varchar,64,`
  - `Guid,用户Guid,varchar,64,`
  - `Vid,视频ID,varchar,64,`
  - `Cid,专辑ID,varchar,64,`
  - `Lid,栏目ID,varchar,64`
- 实现业务自定义的解码函数，自定义解码函数可以是整个框架级别的解码函数，也可以是针对某个接口实现的解码函数，通过注册名称区分。如下是视频默认框架级别扩展解码函数TVAttaDecodeFunc，通过返回值[]string返回对应的扩展字段信息，各个业务可以自行扩展实现自己的解码函数。
```
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
```
- 在程序初始化的地方，将自定义的解码函数注册到组件中即可。
```
   import attaCodec "git.woa.com/trpc-extend/trpc-go/atta/codec"
   
   // 注册整个框架默认解码函数
   attaCodec.RegisterDefault(&TestReportCodec{})
   // 注册trpc接口为ReadRoomInfo对应的解码函数
   attaCodec.Register("/trpc.video_detail.sport_national_group.NationalGroupService/GetIntroduction", &TestReportCodec{})
```