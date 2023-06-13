package gin

import "html/template"

const protoPbTmp = "package {{ .ModuleName }}\n{{ range $value := .Messages }}\n{{ $value }}\n{{ end }}"

type Pb struct {
	ModuleName string          `json:"module_name"`
	Messages   []template.HTML `json:"messages"`
}

func (Pb) ProtoPbTmp() string {
	return protoPbTmp
}
