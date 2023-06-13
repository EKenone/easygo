package gin

const serviceBusiness = "\n// {{ .Name }} {{ .Desc }}\nfunc (s *{{ .Group }}) {{ .Name }}(req *pb.{{ .Group }}{{ .Name }}Req) (res *pb.{{ .Group }}{{ .Name }}Res, errCode errcode.Error) {\n    //TODO 实现业务\n\n    res = &pb.{{ .Group }}{{ .Name }}Res{}\n\treturn\n}\n"

type ServiceBusiness struct {
	Group string `json:"group"` //接口组名称
	Name  string `json:"name"`  //接口名称
	Desc  string `json:"desc"`  //接口备注
}

func (ServiceBusiness) ServiceBusiness() string {
	return serviceBusiness
}
