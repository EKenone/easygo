package gin

import "html/template"

const swaggerDoc = "package {{ .ModuleName }}\n\n{{ .Content }}\nfunc docDesc() {\n\n}\n"

type SwaggerDoc struct {
	Content    template.HTML `json:"content"`
	ModuleName string        `json:"module_name"`
}

func (SwaggerDoc) SwaggerDoc() string {
	return swaggerDoc
}
