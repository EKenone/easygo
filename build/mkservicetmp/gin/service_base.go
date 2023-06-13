package gin

const serviceBaseTmp = "package {{ .PkgName }}\n\nimport (\n\t\"context\"\n\t\"github.com/gin-gonic/gin\"\n)\n\ntype Service struct {\n    GinCtx *gin.Context\n\tCtx context.Context\n}\n\nfunc (s *Service) Init(ctx *gin.Context) {\n\ts.GinCtx = ctx\n    s.Ctx = ctx.Request.Context()\n}\n"

type ServiceBase struct {
	ModName string //go.mod的module名字
	PkgName string //包名称
}

func (ServiceBase) ServiceBaseTmp() string {
	return serviceBaseTmp
}
