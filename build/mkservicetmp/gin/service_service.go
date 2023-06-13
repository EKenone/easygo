package gin

const serviceTmp = "package {{ .ModuleName }}\n\nimport (\n\tpb \"{{ .ModName }}/{{ .ProtocolPath }}/{{ .ModuleName }}\"\n\t\"{{ .ModName }}/{{ .ServicePath }}\"\n\t\"{{ .ModName }}/{{ .ErrCode }}\"\n)\n\ntype {{ .Group }} struct {\n\t{{ .PkgName }}.Service\n}\n"

type Service struct {
	ModName      string //go.mod的module名字
	ModuleName   string //模块名称
	Group        string //组名称
	ProtocolPath string //协议路径
	ServicePath  string //业务路径
	PkgName      string //包名称
	ErrCode      string //错误吗存放的地方
}

func (Service) ServiceTmp() string {
	return serviceTmp
}
