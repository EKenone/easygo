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
	"regexp"
	"strings"
)

type Api struct {
	moduleName  string  //模块名称
	moduleRoute string  //模块路由
	parses      []Parse //解析的接口结构体
	pkgName     string  //包名字
	errCode     string  //错误码包
	interPath   string  //资源包
}

// NewApi 实例化api
func NewApi(moduleName, moduleRoute string, parses []Parse) *Api {
	idx := strings.LastIndex(config.Conf.Internal.Entrance, "/")

	return &Api{
		moduleName:  moduleName,
		moduleRoute: moduleRoute,
		parses:      parses,
		pkgName:     config.Conf.Internal.Entrance[idx+1:],
		errCode:     config.Conf.Internal.Entrance[:idx] + "/errcode",
		interPath:   config.Conf.Internal.Entrance[:idx],
	}
}

// Build 构建api业务处理
func (a *Api) Build() (err error) {
	if config.Conf.Internal.Entrance == "" {
		return
	}

	dir, err := a.dir()
	if err != nil {
		return
	}

	err = a.init()
	if err != nil {
		return
	}

	err = a.base(dir)
	if err != nil {
		return
	}

	for _, v := range a.parses {
		err = a.append(dir, v)
		if err != nil {
			return
		}
	}

	return
}

func (a *Api) init() (err error) {
	exist, err := utils.IsFileOrDirExist(config.Conf.Internal.Entrance + "/api.go")
	if err != nil {
		return
	}

	//模板内容
	var content = gin.EntranceBase{
		ModName: config.Conf.ModName,
		Package: a.pkgName,
		ErrCode: a.errCode,
	}

	if !exist {
		var t *template.Template
		t, err = template.New("base.tpl").Parse(content.Tmp())
		if err != nil {
			return
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, content)
		if err != nil {
			return
		}

		//把模板写入文件
		var file *os.File
		file, err = os.OpenFile(config.Conf.Internal.Entrance+"/api.go", os.O_CREATE|os.O_RDWR, 0777)
		defer file.Close()
		if err != nil {
			return
		}

		_, err = file.WriteString(buf.String())
		if err != nil {
			return
		}
	}

	exist, err = utils.IsFileOrDirExist(a.interPath + "/errcode")
	if err != nil {
		return
	}

	if exist {
		return
	}

	err = os.MkdirAll(a.interPath+"/errcode", 0777)
	if err != nil {
		return
	}

	//把模板写入文件
	file2, err := os.OpenFile(a.interPath+"/errcode/err_code.go", os.O_CREATE|os.O_RDWR, 0777)
	defer file2.Close()
	if err != nil {
		return
	}

	_, err = file2.WriteString(content.ErrCodeTmp())
	if err != nil {
		return
	}

	//把模板写入文件
	file3, err := os.OpenFile(a.interPath+"/errcode/define.go", os.O_CREATE|os.O_RDWR, 0777)
	defer file3.Close()
	if err != nil {
		return
	}

	_, err = file3.WriteString(content.ErrCodeDefine())
	if err != nil {
		return
	}

	//把模板写入文件
	file4, err := os.OpenFile(a.interPath+"/errcode/method.go", os.O_CREATE|os.O_RDWR, 0777)
	defer file4.Close()
	if err != nil {
		return
	}

	_, err = file4.WriteString(content.ErrCodeMethod())
	if err != nil {
		return
	}

	return
}

// 创建目录
func (a *Api) dir() (dir string, err error) {
	//检查目录是否存在
	dirName := strings.ToLower(a.moduleName)
	dir = fmt.Sprintf("%s/%s", config.Conf.Internal.Entrance, dirName)

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
func (a *Api) base(dir string) (err error) {
	var controllers = make([]string, 0, len(a.parses))
	for _, v := range a.parses {
		controllers = append(controllers, utils.CamelString(v.Group.Name))
	}

	//检查基础控制器文件是否存在
	filename := fmt.Sprintf("%s/controller.go", dir)
	exist, err := utils.IsFileOrDirExist(filename)
	if err != nil {
		return
	}

	//存在的时候就
	if exist {
		err = os.Remove(filename)
		if err != nil {
			return
		}
	}

	//模板内容
	var content = gin.Controller{
		ModName:       config.Conf.ModName,
		ModuleName:    strings.ToLower(a.moduleName),
		ModuleRoute:   a.moduleRoute,
		HasController: len(controllers) > 0,
		Controllers:   controllers,
		ApiPath:       config.Conf.Internal.Entrance,
		ApiPkg:        a.pkgName,
		ServerPath:    a.interPath,
	}

	exist, err = utils.IsFileOrDirExist(a.interPath + "/ginserver")
	if err != nil {
		return
	}

	if !exist {
		err = os.MkdirAll(a.interPath+"/ginserver", 0777)
		if err != nil {
			return
		}

		//把模板写入文件
		file5, err := os.OpenFile(a.interPath+"/ginserver/gin_server.go", os.O_CREATE|os.O_RDWR, 0777)
		defer file5.Close()
		if err != nil {
			return err
		}

		_, err = file5.WriteString(content.GinServerTmp())
		if err != nil {
			return err
		}
	}

	t, err := template.New("controller.tpl").Parse(content.ControllerTmp())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	//把模板写入文件
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
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

// 追加接口
func (a *Api) append(dir string, parse Parse) (err error) {
	//检查基础控制器文件是否存在
	filename := fmt.Sprintf("%s/%s.go", dir, utils.MidString(parse.Group.Name, '_'))

	err = a.createGroup(filename, parse.Group)
	if err != nil {
		return
	}

	file, err := os.OpenFile(filename, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		return
	}

	body, err := io.ReadAll(file)
	if err != nil {
		return
	}

	err = file.Truncate(0)
	if err != nil {
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return
	}

	content := string(body)

	regStr := `return []ginaction.Action{` + "\n"

	for _, v := range parse.Interfaces {

		var method = "GET"
		var lastPath, midType string
		if v.LastPath != "" {
			lastPath = fmt.Sprintf(", ginaction.UseLastPath(\"%s\")", v.LastPath)
		}

		if v.MidType != nil {
			midType = fmt.Sprintf(", ginaction.UseMidType(%v)", v.MidType)
		}

		if v.Method != "" {
			method = strings.ToUpper(v.Method)
		}

		regStr += fmt.Sprintf("\t\tginaction.NewAction(\"%s\", s.%s%s%s),\n", method, v.Name, lastPath, midType)

		if strings.Contains(content, fmt.Sprintf("%s(ctx *gin.Context)", v.Name)) {
			continue
		}

		var str string
		str, err = a.createApi(parse.Group.Name, v)
		if err != nil {
			return
		}

		content += str
	}

	regStr += "\t}"

	reg, err := regexp.Compile(`return \[]ginaction\.Action\{([\s\S]*?)}`)
	if err != nil {
		return
	}

	_, err = file.WriteString(reg.ReplaceAllLiteralString(content, regStr))

	return
}

// 创建接口字符串
func (a Api) createApi(group string, info Interface) (str string, err error) {

	//模板内容
	var content = gin.EntranceApi{
		Group:  group,
		Name:   info.Name,
		Desc:   info.Desc,
		Method: info.Method,
		Req:    fmt.Sprintf("%s%sReq", group, info.Name),
	}

	t, err := template.New("api.tpl").Parse(content.Tmp())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	str = buf.String()
	return
}

// 创建控制器组
func (a Api) createGroup(filename string, group Group) (err error) {
	//检查基础控制器文件是否存在
	exist, err := utils.IsFileOrDirExist(filename)
	if err != nil {
		return
	}

	if exist {
		return
	}

	//模板内容
	var content = gin.Group{
		ModName:     config.Conf.ModName,
		ModuleName:  strings.ToLower(a.moduleName),
		Name:        group.Name,
		Desc:        group.Desc,
		Route:       utils.MidString(group.Name, '-'),
		Path:        config.Conf.Proto.Analysis,
		ServicePath: config.Conf.Internal.Service,
	}

	t, err := template.New("group.tpl").Parse(content.GroupTmp())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	//把模板写入文件
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
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
