package gin

const swaggerTmp = "package {{ .ModuleName }}\n\ntype {{ .Group }} struct {\n\n}"

type Swagger struct {
	ModuleName string
	Group      string
}

func (Swagger) SwaggerTmp() string {
	return swaggerTmp
}
