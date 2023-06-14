package gin

const entranceGroupTmp = "package {{ .ModuleName }}\n\nimport (\n\tpb \"{{ .ModName }}/{{ .Path }}/{{ .ModuleName }}\"\n\t\"{{ .ModName }}/{{ .ServicePath }}/{{ .ModuleName }}\"\n\t\"github.com/tdeken/ginaction\"\n\t\"github.com/gin-gonic/gin\"\n)\n\n// {{ .Name }} {{ .Desc }}\ntype {{ .Name }} struct {\n\tController\n\tdep {{ .ModuleName }}.{{ .Name }}\n}\n\n// 依赖业务服务\nfunc (s {{ .Name }}) depService(ctx *gin.Context) {{ .ModuleName }}.{{ .Name }} {\n\ts.dep.Init(ctx)\n\treturn s.dep\n}\n\n// Group 基础请求组\nfunc (s {{ .Name }}) Group() string {\n\treturn \"{{ .Route }}\"\n}\n\n// Register 注册路由\nfunc (s {{ .Name }}) Register() []ginaction.Action {\n\treturn []ginaction.Action{\n\t}\n}\n"

type Group struct {
	ModName     string //go.mod的module名字
	ModuleName  string //模块名称
	Name        string //接口组名称
	Desc        string //接口组备注
	Route       string //接口组路由
	Path        string //协议路径
	ServicePath string //业务路径
}

func (Group) GroupTmp() string {
	return entranceGroupTmp
}
