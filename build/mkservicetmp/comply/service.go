package comply

import (
	"bytes"
	"easygo/build/mkservicetmp/config"
	"easygo/build/mkservicetmp/gin"
	"easygo/build/mkservicetmp/utils"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
)

type Service struct {
	moduleName string  //模块名称
	parses     []Parse //解析的接口结构体
	pkgName    string  //包名字
	errCode    string  //错误码包
}

// NewService 实例化proto
func NewService(moduleName string, parses []Parse) *Service {

	idx := strings.LastIndex(config.Conf.Internal.Service, "/")

	return &Service{
		moduleName: moduleName,
		parses:     parses,
		pkgName:    config.Conf.Internal.Service[idx+1:],
		errCode:    config.Conf.Internal.Service[:idx] + "/errcode",
	}
}

// Build 构建api业务处理
func (b *Service) Build() (err error) {
	if config.Conf.Internal.Service == "" {
		return
	}

	dir, err := b.dir()
	if err != nil {
		return
	}

	err = b.init()
	if err != nil {
		return
	}

	for _, v := range b.parses {
		err = b.base(dir, v)
		if err != nil {
			return
		}

		err = b.append(dir, v)
		if err != nil {
			return
		}
	}

	return
}

func (b *Service) init() (err error) {
	exist, err := utils.IsFileOrDirExist(config.Conf.Internal.Service)
	if err != nil {
		return
	}

	if exist {
		return
	}

	//模板内容
	var content = gin.ServiceBase{
		ModName: config.Conf.ModName,
		PkgName: b.pkgName,
	}

	t, err := template.New("base.tpl").Parse(content.ServiceBaseTmp())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	//把模板写入文件
	file, err := os.OpenFile(config.Conf.Internal.Service+"/service.go", os.O_CREATE|os.O_RDWR, 0777)
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

// 创建目录
func (b *Service) dir() (dir string, err error) {
	//检查目录是否存在
	dirName := strings.ToLower(b.moduleName)
	dir = fmt.Sprintf("%s/%s", config.Conf.Internal.Service, dirName)
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

// 创建基础环境
func (b *Service) base(dir string, parse Parse) (err error) {

	//检查基础控制器文件是否存在
	filename := fmt.Sprintf("%s/%s.go", dir, utils.MidString(parse.Group.Name, '_'))

	exist, err := utils.IsFileOrDirExist(filename)
	if err != nil {
		return
	}

	//存在的时候就
	if exist {
		return
	}

	//模板内容
	var content = gin.Service{
		ModName:      config.Conf.ModName,
		ModuleName:   strings.ToLower(b.moduleName),
		Group:        parse.Group.Name,
		ProtocolPath: config.Conf.Proto.Analysis,
		ServicePath:  config.Conf.Internal.Service,
		PkgName:      b.pkgName,
		ErrCode:      b.errCode,
	}

	t, err := template.New("service.tpl").Parse(content.ServiceTmp())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	//把模板写入文件
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0777)
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

// 创建基础环境
func (b *Service) append(dir string, parse Parse) (err error) {

	t, err := template.New("business.tpl").Parse(gin.ServiceBusiness{}.ServiceBusiness())
	if err != nil {
		return
	}

	//检查基础控制器文件是否存在
	filename := fmt.Sprintf("%s/%s.go", dir, utils.MidString(parse.Group.Name, '_'))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		return
	}

	body, err := io.ReadAll(file)

	text := string(body)

	for _, v := range parse.Interfaces {

		if strings.Contains(text, fmt.Sprintf("func (s %s) %s(req", parse.Group.Name, v.Name)) {
			continue
		}

		//模板内容
		var content = gin.ServiceBusiness{
			Group: parse.Group.Name,
			Name:  v.Name,
			Desc:  v.Desc,
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, content)
		if err != nil {
			return
		}

		_, err = file.WriteString(buf.String())
		if err != nil {
			return
		}
	}

	return
}
