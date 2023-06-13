package gin

const routeTmp = "package {{ .RoutePkg }}\n\nimport ({{ range $value := .Pkg }}\n    \"{{ $.ModName }}/{{ $.ApiPath }}/{{ $value }}\"{{ end }}\n)\n\n// InitRoute 启动路由\nfunc InitRoute() { {{ range $value := .Pkg }}\n\t{{ $value }}.Controller{}.Route(){{ end }}\n}\n"

type Route struct {
	ModName  string
	ApiPath  string
	Pkg      []string
	RoutePkg string
}

func (Route) RouteTmp() string {
	return routeTmp
}
