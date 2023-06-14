package comply

import (
	"bytes"
	"fmt"
	"github.com/tdeken/easygo/build/mkservicetmp/config"
	"github.com/tdeken/easygo/build/mkservicetmp/gin"
	"github.com/tdeken/easygo/build/mkservicetmp/utils"
	"html/template"
	"os"
	"os/exec"
	"strings"
)

type Proto struct {
	moduleName string          //模块名称
	doc        Doc             //模块说明文档
	parses     []Parse         //解析的接口结构体
	messages   []template.HTML //结构体数据
}

// NewProto 实例化proto
func NewProto(moduleName string, doc Doc, parses []Parse) *Proto {
	return &Proto{
		moduleName: moduleName,
		doc:        doc,
		parses:     parses,
		messages:   nil,
	}
}

// Build 构建api业务处理
func (b *Proto) Build() (err error) {
	if config.Conf.Proto.Analysis == "" {
		return
	}

	dir, err := b.dir()
	if err != nil {
		return
	}

	for _, v := range b.parses {
		err = b.append(dir, v)
		if err != nil {
			return
		}
	}

	err = exec.Command("gofmt", "-w", dir).Run()

	return
}

// 创建目录
func (b *Proto) dir() (dir string, err error) {
	//检查目录是否存在
	dirName := strings.ToLower(b.moduleName)
	dir = fmt.Sprintf("%s/%s", config.Conf.Proto.Analysis, dirName)
	exist, err := utils.IsFileOrDirExist(dir)
	if err != nil {
		return
	}

	if exist {
		return
	}

	//不存在就创建目录
	err = utils.MkDir(dir)
	if err != nil {
		return
	}

	return
}

// 追加结构体
func (b *Proto) append(dir string, parse Parse) (err error) {
	defer b.reset()
	for _, v := range parse.Interfaces {
		var uriName string
		if len(v.Uri) > 0 {
			uriName = parse.Group.Name + v.Name + "Uri"
			b.parseUri(uriName, v.Uri)
		}

		v.Req.Name = parse.Group.Name + v.Name
		b.parseReq(v.Req, "Req", uriName)
		v.Res.Name = parse.Group.Name + v.Name
		b.parseRes(v.Res, "Res")
	}

	//模板内容
	var content = gin.Pb{
		ModuleName: b.moduleName,
		Messages:   b.messages,
	}

	t, err := template.New("pb.tpl").Parse(content.ProtoPbTmp())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	filename := dir + "/" + utils.MidString(parse.Group.Name, '_') + ".go"

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

// 解析请求结构体
func (b *Proto) parseUri(Name string, fields []Field) {
	sn := "type " + Name + " struct { \n"
	for _, field := range fields {
		var bind string
		var jsonTag string

		sn += fmt.Sprintf("\t%s %s `json:\"%s%s\" %s%s` //%s\n",
			utils.CamelString(field.Name),
			field.Class,
			field.Name,
			jsonTag,
			b.parseFromTag("uri", field.Name),
			bind,
			field.Desc,
		)
	}
	sn += "}"

	b.messages = append(b.messages, template.HTML(sn))

}

// 解析请求结构体
func (b *Proto) parseReq(req Message, suffix, uriName string) {
	sn := "type " + req.Name + suffix + " struct { \n"
	if uriName != "" {
		sn += uriName + "\n"
	}

	for _, field := range req.Fields {
		var bind string
		if field.Binding != "" {
			bind = " binding:\"" + field.Binding + "\""
		}

		var jsonTag string
		if len(field.JsonTag) > 0 {
			jsonTag = "," + strings.Join(field.JsonTag, ",")
		}

		sn += fmt.Sprintf("\t%s %s `json:\"%s%s\" %s%s` //%s\n",
			utils.CamelString(field.Name),
			utils.GetClass(req.Name, field.Class),
			field.Name,
			jsonTag,
			b.parseFromTag(field.From, field.Name),
			bind,
			field.Desc,
		)
	}
	sn += "}"

	b.messages = append(b.messages, template.HTML(sn))

	for _, v := range req.Message {
		v.Name = req.Name + v.Name
		b.parseReq(v, "", "")
	}
}

func (b *Proto) parseFromTag(from string, fieldName string) string {
	switch from {
	case "none":
		return ""
	case "uri":
		return fmt.Sprintf("uri:\"%s\"", fieldName)
	default:
		return fmt.Sprintf("form:\"%s\"", fieldName)
	}
}

// 解析返回结构体
func (b *Proto) parseRes(res Message, suffix string) {
	sn := "type " + res.Name + suffix + " struct { \n"
	for _, field := range res.Fields {
		var jsonTag string
		if len(field.JsonTag) > 0 {
			jsonTag = "," + strings.Join(field.JsonTag, ",")
		}

		sn += fmt.Sprintf("\t%s %s `json:\"%s%s\"` //%s\n", utils.CamelString(field.Name), utils.GetClass(res.Name, field.Class), field.Name, jsonTag, field.Desc)
	}
	sn += "}"

	b.messages = append(b.messages, template.HTML(sn))

	for _, v := range res.Message {
		v.Name = res.Name + v.Name
		b.parseRes(v, "")
	}
}

// 重置
func (b *Proto) reset() {
	b.messages = nil
}
