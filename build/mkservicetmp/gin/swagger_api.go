package gin

import "html/template"

const swaggerApiTmp = "\n{{ range $value := .Messages }}\n{{ $value }}\n{{ end }}\n\n// {{ .Name }}\n// @Tags {{ .GroupDesc }}\n// @Summary {{ .Desc }}\n// {{if .SecurityTitle}}@Security {{ .SecurityTitle }} {{ end }}\n// @accept {{ .ReqContentType }}\n// @Produce {{ .ResContentType }}\n// @Param data {{ .Body }} {{ .Req }} true \"数据\"\n// @Success 200 {object} {{ .Res }}\n// @Router {{ .Route }} [{{ .Method }}]\nfunc ({{ .Group }}) {{ .Name }}() {\n\n}"

type SwaggerApi struct {
	Name           string          //方法名称
	Desc           string          //方法备注
	Group          string          //组名称
	GroupDesc      string          //组备注
	Body           string          //数据提交方式
	Req            string          //请求结构体
	Res            string          //返回结构体
	Route          string          //路由
	Method         string          //请求方法
	ReqContentType string          //请求数据格式
	ResContentType string          //返回数据格式
	SecurityTitle  string          //身份校验
	Messages       []template.HTML `json:"messages"`
}

func (SwaggerApi) SwaggerApiTmp() string {
	return swaggerApiTmp
}
