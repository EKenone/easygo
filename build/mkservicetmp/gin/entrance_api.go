package gin

const entranceApiTmp = "\n// {{ .Name }} {{ .Desc }}\nfunc (s {{ .Group }}) {{ .Name }}(ctx *gin.Context) {\n\tvar form = &pb.{{ .Req }}{}\n\tif err := s.ValidateRequest(ctx, form); err != nil {\n\t\ts.Response(ctx, nil, err)\n\t\treturn\n\t}\n\n\tres, err := s.depService(ctx).{{ .Name }}(form)\n\ts.Response(ctx, res, err)\n}\n"

type EntranceApi struct {
	Group  string `json:"group"`  //接口组名称
	Name   string `json:"name"`   //接口名称
	Desc   string `json:"desc"`   //接口备注
	Method string `json:"method"` //请求方法
	Req    string `json:"req"`    //请求表单
}

func (EntranceApi) Tmp() string {
	return entranceApiTmp
}
