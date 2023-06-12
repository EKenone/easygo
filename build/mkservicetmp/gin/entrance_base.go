package gin

const entranceBase = "package api\n\nimport (\n\t\"hongru/internal/errcode\"\n\t\"hongru/internal/ginserver/action\"\n\t\"hongru/internal/plugins/logger\"\n\t\"hongru/internal/plugins/validator\"\n\t\"encoding/json\"\n\t\"github.com/gin-gonic/gin\"\n)\n\ntype Api struct {\n}\n\n// ChooseMid 可以选择的服务中间件\nfunc (a Api) ChooseMid(router *gin.RouterGroup, t action.MidType) gin.IRoutes {\n\tswitch t {\n\tdefault:\n\t\treturn router\n\t}\n}\n\n// Group 接口组标识\nfunc (a Api) Group() string {\n\treturn \"\"\n}\n\n// Register 接口注册\nfunc (a Api) Register() []action.Action {\n\treturn nil\n}\n\n// ValidateRequest 统一校验请求数据\nfunc (a Api) ValidateRequest(ctx *gin.Context, rt validator.RequestInterface) errcode.Error {\n\terr := validator.CheckParams(ctx, rt)\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tb, _ := json.Marshal(rt)\n\tlogger.InfoFCtx(ctx.Request.Context(), \"请求路由: %s, 请求方法: %s, 请求数据: %s\", ctx.Request.URL.Path, ctx.Request.Method, string(b))\n\n\treturn nil\n}\n\n// Response 服务返回值\nfunc (a Api) Response(ctx *gin.Context, res interface{}, errCode errcode.Error) {\n\tdefer func() {\n\t\tb, _ := json.Marshal(res)\n\t\tlogger.InfoFCtx(ctx.Request.Context(), \"返回数据: %s, 返回错误: %v\", string(b), errCode)\n\t}()\n\n\t//设置code\n\tvar rsp = gin.H{\"code\": errcode.OK, \"msg\": \"ok\", \"data\": res}\n\tif errCode != nil {\n\t\trsp[\"code\"] = errCode.Code()\n\t\trsp[\"msg\"] = errCode.Msg()\n\t}\n\tctx.JSON(200, rsp)\n\n}\n"

type EntranceBase struct {
	ModName string //go.mod的module名字
	Package string //包名称
}

func (EntranceBase) Tmp() string {
	return entranceBase
}
