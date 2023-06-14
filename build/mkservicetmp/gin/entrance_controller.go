package gin

const entranceControllerTmp = "package {{ .ModuleName }}\n\nimport (\n\t\"{{ .ModName }}/{{ .ApiPath }}\"\n\t\"{{ .ModName }}/{{ .ServerPath }}/ginserver\"\n\t\"github.com/tdeken/ginaction\"\n)\n\ntype Controller struct {\n\t{{ .ApiPkg }}.Api\n}\n\n// Route 模块路由\nfunc (c Controller) Route() {\n\tr := ginserver.Server.Group(\"{{ .ModuleRoute }}\")\n\n\tginaction.AutoRegister(r, {{ if .HasController }}{{ range $value := .Controllers }}\n\t    {{ $value }}{},{{end}}{{ else }}nil,{{ end }}\n    )\n}"

const ginServerTmp = "package ginserver\n\nimport (\n\t\"context\"\n\t\"github.com/gin-gonic/gin\"\n)\n\nvar Server *gin.Engine\n\nfunc Init() {\n\tServer = gin.Default()\n}\n\nfunc Run(ctx context.Context) {\n\tgo func() {\n\t\tif err := Server.Run(\":8080\"); err != nil {\n\t\t\tpanic(err)\n\t\t}\n\t}()\n}"

type Controller struct {
	ModName       string   //go.mod的module名字
	ModuleName    string   //模块名称
	ModuleRoute   string   //模块路由
	HasController bool     //是否有控制器
	Controllers   []string //控制器名称
	ApiPath       string   //api结构体路径
	ApiPkg        string   //api包名字
	ServerPath    string   //全局的Gin服务变量
}

func (Controller) ControllerTmp() string {
	return entranceControllerTmp
}

func (Controller) GinServerTmp() string {
	return ginServerTmp
}
