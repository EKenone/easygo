package gin

const entranceBase = "package {{ .Package }}\n\nimport (\n\t\"{{ .ModName }}/{{ .ErrCode }}\"\n\t\"github.com/tdeken/ginaction\"\n\t\"github.com/gin-gonic/gin\"\n)\n\ntype Api struct {\n}\n\n// ChooseMid 可以选择的服务中间件\nfunc (a Api) ChooseMid(router *gin.RouterGroup, t ginaction.MidType) gin.IRoutes {\n\tswitch t {\n\tdefault:\n\t\treturn router\n\t}\n}\n\n// Group 接口组标识\nfunc (a Api) Group() string {\n\treturn \"\"\n}\n\n// Register 接口注册\nfunc (a Api) Register() []ginaction.Action {\n\treturn nil\n}\n\n// ValidateRequest 统一校验请求数据\nfunc (a Api) ValidateRequest(ctx *gin.Context, rt any) errcode.Error {\n    if err := ctx.BindUri(rt); err != nil {\n        return errcode.VerifyErrorWithTip(err.Error())\n    }\n\n    if err := ctx.ShouldBindQuery(rt); err != nil {\n        return errcode.VerifyErrorWithTip(err.Error())\n    }\n\n\treturn nil\n}\n\n// Response 服务返回值\nfunc (a Api) Response(ctx *gin.Context, res interface{}, errCode errcode.Error) {\n\t//设置code\n\tvar rsp = gin.H{\"code\": errcode.OK, \"msg\": \"ok\", \"data\": res}\n\tif errCode != nil {\n\t\trsp[\"code\"] = errCode.Code()\n\t\trsp[\"msg\"] = errCode.Msg()\n\t}\n\t\n\tctx.JSON(200, rsp)\n}\n"

const (
	errCodeTmp    = "package errcode\n\nimport (\n\t\"fmt\"\n)\n\ntype Error interface {\n\terror\n\tCode() int32\n\tMsg() string\n}\n\n// Is 判断两个接口是否相等\nfunc Is(error1 Error, error2 ...Error) bool {\n\n\tif error1 == nil || len(error2) == 0 {\n\t\treturn false\n\t}\n\n\tfor _, v := range error2 {\n\t\tif v.Code() == error1.Code() {\n\t\t\treturn true\n\t\t}\n\t}\n\n\treturn false\n}\n\ntype CodeError struct {\n\tcode int32  //错误码\n\tmsg  string //错误信息\n}\n\n// Code 状态码\nfunc (e *CodeError) Code() int32 {\n\treturn e.code\n}\n\n// Msg 状态码说明\nfunc (e *CodeError) Msg() string {\n\treturn e.msg\n}\n\n// 实例化一个错误\nfunc newError(code int32, msg string) *CodeError {\n\treturn &CodeError{code: code, msg: msg}\n}\n\n// Error 错误信息\nfunc (e *CodeError) Error() string {\n\treturn fmt.Sprintf(\"错误码：%d, 错误信息：%s\", e.Code(), e.Msg())\n}"
	errCodeDefine = "package errcode\n\n// OK 服务请求正常标识\nconst OK = 200\n\n// 系统错误码\nconst (\n\tWorkErrorCode         = 0\n\tAuthErrorCode         = 401\n\tVerifyErrorCode       = 400\n\tDataNotFoundErrorCode = 404\n\tLimitErrorCode        = 405\n\tServerErrorCode       = 500\n\tServerBusyErrorCode   = 502\n)\n\n// 系统\nvar (\n\tServerError       = newError(ServerErrorCode, \"服务内部错误\")       //代码错误，连接错误等\n\tDataNotFoundError = newError(DataNotFoundErrorCode, \"数据不存在\")  //代码错误，连接错误等\n\tAuthError         = newError(AuthErrorCode, \"登陆已失效\")          //权限错误（例如：数据的企业ID与登录者的企业ID不一致，登录者身份信息异常等）\n\tVerifyError       = newError(VerifyErrorCode, \"参数错误\")         //验证前端上传上来的表单，如果需要带一些详细信息，请用 VerifyErrorWithTip 方法\n\tServerBusyError   = newError(ServerBusyErrorCode, \"服务繁忙，请重试\") //业务处理超时，如：生成订单多次都没生成成功，获取锁超时\n\tLimitError        = newError(LimitErrorCode, \"请勿频繁请求\")        //限流错误码，如：一定时间内不允许重复请求限制错误\n)\n"
	errCodeMethod = "package errcode\n\nimport (\n\t\"fmt\"\n)\n\n// VerifyErrorWithTip 校验错误，并返回错误提示\nfunc VerifyErrorWithTip(tip string) Error {\n\treturn newError(VerifyErrorCode, fmt.Sprintf(\"校验错误:%s\", tip))\n}\n"
)

type EntranceBase struct {
	ModName string //go.mod的module名字
	Package string //包名称
	ErrCode string //错误吗存放的地方
}

func (EntranceBase) Tmp() string {
	return entranceBase
}

func (EntranceBase) ErrCodeTmp() string {
	return errCodeTmp
}

func (EntranceBase) ErrCodeDefine() string {
	return errCodeDefine
}

func (EntranceBase) ErrCodeMethod() string {
	return errCodeMethod
}
