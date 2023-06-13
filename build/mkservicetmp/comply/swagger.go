package comply

import (
	"bytes"
	"easygo/build/mkservicetmp/config"
	"easygo/build/mkservicetmp/gin"
	"easygo/build/mkservicetmp/utils"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"
)

type Swagger struct {
	moduleName  string          //模块名称
	moduleRoute string          //模块路由
	doc         Doc             //模块说明文档
	parses      []Parse         //解析的接口结构体
	messages    []template.HTML //结构体数据
}

// NewSwagger 实例化Swagger
func NewSwagger(moduleName, moduleRoute string, doc Doc, parses []Parse) *Swagger {
	return &Swagger{
		moduleName:  moduleName,
		moduleRoute: moduleRoute,
		doc:         doc,
		parses:      parses,
		messages:    nil,
	}
}

// Build 构建api业务处理
func (b *Swagger) Build() (err error) {
	if config.Conf.Doc.Swagger == "" {
		return
	}

	dir, err := b.dir()
	if err != nil {
		return
	}

	for _, v := range b.parses {
		err = b.create(dir, v)
		if err != nil {
			return
		}

		for _, v1 := range v.Interfaces {
			err = b.append(dir, v.Group, v1)
			if err != nil {
				return
			}
		}

	}

	err = exec.Command("gofmt", "-w", dir).Run()

	return
}

func (b *Swagger) dir() (dir string, err error) {
	dir = fmt.Sprintf("%s/%s", config.Conf.Doc.Swagger, b.moduleName)
	exist, err := utils.IsFileOrDirExist(dir)
	if err != nil {
		return
	}

	if !exist {
		//不存在就创建目录
		err = utils.MkDir(dir)
		if err != nil {
			return
		}

		var generateExist bool
		generateExist, err = utils.IsFileOrDirExist("generate.go")
		if err != nil || !generateExist {
			return
		}

		var generateFile *os.File
		generateFile, err = os.OpenFile("generate.go", os.O_CREATE|os.O_APPEND, 0777)
		defer generateFile.Close()
		if err != nil {
			return
		}

		_, err = generateFile.WriteString(fmt.Sprintf("//go:generate swag init  -o docs/%s -g doc.go -d %s \n", b.moduleName, dir))
		if err != nil {
			return
		}

	}

	var text string

	if b.doc.Title != "" {
		text += "\n" + "// @title " + b.doc.Title
	}

	if b.doc.Host != "" {
		text += "\n" + "// @host " + b.doc.Host
	}

	if b.doc.Schemes != nil {
		text += "\n" + "// @schemes " + strings.Join(b.doc.Schemes, " ")
	}

	if b.doc.Ver != "" {
		text += "\n" + "// @version " + b.doc.Ver
	}

	if b.doc.Desc != "" {
		text += "\n" + "// @description " + b.doc.Desc
	}

	if b.doc.Auth.Security != "" {
		text += "\n" + "// @securityDefinitions." + b.doc.Auth.Security + " " + b.doc.Auth.Title
	}

	if b.doc.Auth.In != "" {
		text += "\n" + "// @in " + b.doc.Auth.In
	}

	if b.doc.Auth.Name != "" {
		text += "\n" + "// @name " + b.doc.Auth.Name
	}

	if b.doc.Contact.Name != "" {
		text += "\n" + "// @contact.name " + b.doc.Contact.Name
	}

	if b.doc.Contact.Url != "" {
		text += "\n" + "// @contact.url " + b.doc.Contact.Url
	}

	if b.doc.Contact.Email != "" {
		text += "\n" + "// @contact.email " + b.doc.Contact.Email
	}

	//模板内容
	var content = gin.SwaggerDoc{
		Content:    template.HTML(text),
		ModuleName: b.moduleName,
	}

	t, err := template.New("doc.tpl").Parse(content.SwaggerDoc())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	filename := dir + "/doc.go"

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	defer file.Close()
	if err != nil {
		return
	}

	_, err = file.WriteString(buf.String())
	if err != nil {
		return
	}

	return
}

func (b *Swagger) create(dir string, v Parse) (err error) {
	//模板内容
	var content = gin.Swagger{
		ModuleName: b.moduleName,
		Group:      v.Group.Name,
	}

	t, err := template.New("swagger.tpl").Parse(content.SwaggerTmp())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	filename := fmt.Sprintf("%s/%s.go", dir, utils.MidString(v.Group.Name, '_'))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		return
	}

	_, err = file.WriteString(buf.String())
	if err != nil {
		return
	}

	return
}

func (b *Swagger) append(dir string, group Group, v Interface) (err error) {
	defer b.reset()
	v.Req.Name = group.Name + v.Name
	b.parseReq(v.Req, "Req")
	v.Res.Name = group.Name + v.Name
	b.parseRes(v.Res, "Res")

	var lastPath = v.LastPath
	if lastPath == "" {
		lastPath = utils.MidString(v.Name, '-')
	}

	var method = v.Method
	if method == "" {
		method = "GET"
	}

	var securityTitle = b.doc.Auth.Title
	if v.NotAuth {
		securityTitle = ""
	}

	var reqContentType = "application/json"
	if v.ReqContentType != "" {
		reqContentType = v.ReqContentType
	}

	var resContentType = "application/json"
	if v.ResContentType != "" {
		resContentType = v.ResContentType
	}

	var body = "body"
	if strings.ToUpper(method) == "GET" {
		body = "query"
	}

	//模板内容
	var content = gin.SwaggerApi{
		Name:           v.Name,
		Desc:           v.Desc,
		Group:          group.Name,
		GroupDesc:      group.Desc,
		Body:           body,
		Req:            v.Req.Name + "Req",
		Res:            v.Res.Name + "Res",
		Route:          fmt.Sprintf("/%s/%s/%s", b.moduleRoute, utils.MidString(group.Name, '-'), lastPath),
		Method:         strings.ToUpper(method),
		ReqContentType: reqContentType,
		ResContentType: resContentType,
		SecurityTitle:  securityTitle,
		Messages:       b.messages,
	}

	t, err := template.New("api.tpl").Parse(content.SwaggerApiTmp())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	filename := dir + "/" + utils.MidString(group.Name, '_') + ".go"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		return
	}

	_, err = file.WriteString(buf.String())
	if err != nil {
		return
	}

	return
}

// 解析请求结构体
func (b *Swagger) parseReq(req Message, suffix string) {
	sn := "type " + req.Name + suffix + " struct { \n"
	for _, field := range req.Fields {
		var bind string
		if field.Binding != "" {
			bind = " binding:\"" + field.Binding + "\""
		}

		sn += fmt.Sprintf("\t%s %s `json:\"%s\" %s` //%s\n", utils.CamelString(field.Name), utils.GetClass(req.Name, field.Class), field.Name, bind, field.Desc)
	}
	sn += "}"

	b.messages = append(b.messages, template.HTML(sn))

	for _, v := range req.Message {
		v.Name = req.Name + v.Name
		b.parseReq(v, "")
	}
}

// 解析返回结构体
func (b *Swagger) parseRes(res Message, suffix string) {
	sn := "type " + res.Name + suffix + " struct { \n"
	for _, field := range res.Fields {
		sn += fmt.Sprintf("\t%s %s `json:\"%s\"` //%s\n", utils.CamelString(field.Name), utils.GetClass(res.Name, field.Class), field.Name, field.Desc)
	}
	sn += "}"

	b.messages = append(b.messages, template.HTML(sn))

	for _, v := range res.Message {
		v.Name = res.Name + v.Name
		b.parseRes(v, "")
	}
}

// 重置
func (b *Swagger) reset() {
	b.messages = nil
}
