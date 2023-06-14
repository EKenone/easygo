package build

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/tdeken/easygo/build/mkservicetmp/comply"
	"github.com/tdeken/easygo/build/mkservicetmp/config"
	"github.com/tdeken/easygo/build/mkservicetmp/gin"
	"github.com/tdeken/easygo/build/mkservicetmp/utils"
	"gopkg.in/yaml.v3"
	"html/template"
	"os"
	"strings"
	"sync"
)

func MkService() (err error) {
	var mk = &mkService{
		init:   "",
		conf:   "",
		module: "",
	}

	fs := flag.NewFlagSet("service", flag.ExitOnError)
	fs.StringVar(&mk.init, "i", "", "init service conf in this path")
	fs.StringVar(&mk.conf, "c", "config/service.yaml", "set your config path when use service")
	fs.StringVar(&mk.module, "m", "all", "scan your module dir, default scan all")
	fs.BoolVar(&mk.help, "help", false, "show help service")
	fs.BoolVar(&mk.help, "h", false, "short var by help")

	// 解析命令行参数
	if err = fs.Parse(os.Args[2:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if mk.help {
		fmt.Fprintln(os.Stderr, "service usage options:")
		fs.PrintDefaults()
		return nil
	}

	if mk.init != "" {
		return mk.createConf()
	}

	return mk.create()
}

type mkService struct {
	conf   string
	module string
	help   bool
	init   string
}

func (i *mkService) createConf() error {
	var filename = i.init
	var path string
	if !strings.HasSuffix(filename, ".yaml") {
		if !strings.HasSuffix(filename, "/") {
			filename += "/"
		}
		path = filename
		filename += "service.yaml"
	} else {
		idx := strings.LastIndex(filename, "/")
		path = filename[:idx]
	}

	err := os.MkdirAll(path, 0777)
	if err != nil {
		return errors.New(fmt.Sprintf("create [%s] dir error : %v", i.init, err))
	}

	exist, err := fileExist(filename)
	if err != nil {
		return errors.New(fmt.Sprintf("create [%s] error : %v", i.init, err))
	}

	if exist {
		return errors.New(fmt.Sprintf("[%s] exist", i.init))
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("create [%s] error : %v", i.init, err))
	}

	_, err = file.WriteString(serviceConfTmp)
	if err != nil {
		return errors.New(fmt.Sprintf("write [%s] error : %v", i.init, err))
	}

	return nil
}

func (i *mkService) create() (err error) {
	data, err := os.ReadFile(i.conf)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &config.Conf)
	if err != nil {
		return
	}

	entries, err := os.ReadDir(config.Conf.Proto.Source)
	if err != nil {
		return errors.New(fmt.Sprintf("read [%s] source err: %v", config.Conf.Proto.Source, err))
	}
	var pkg []string

	for _, v := range entries {
		pkg = append(pkg, v.Name())
	}

	if i.module != "all" {
		md, err := NewModule(i.module)
		if err != nil {
			return errors.New(fmt.Sprintf("create [%s] module err :%v", i.module, err))
		}

		err = md.Build().Error()
		if err != nil {
			return errors.New(fmt.Sprintf("create [%s] module err :%v", i.module, err))
		}

	} else {
		for _, v := range entries {
			fmt.Println(v.Name())
			md, err := NewModule(v.Name())
			if err != nil {
				return errors.New(fmt.Sprintf("create [%s] module err :%v", v.Name(), err))
			}

			err = md.Build().Error()
			if err != nil {
				return errors.New(fmt.Sprintf("create [%s] module err :%v", i.module, err))
			}
		}
	}

	//模板内容
	var content = gin.Route{
		ModName:  config.Conf.ModName,
		ApiPath:  config.Conf.Internal.Entrance,
		Pkg:      pkg,
		RoutePkg: config.Conf.Internal.Route,
	}

	t, err := template.New("route.tpl").Parse(content.RouteTmp())
	if err != nil {
		fmt.Println(err)
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	exist, err := fileExist(config.Conf.Internal.Route)
	if err != nil {
		return
	}

	if !exist {
		err = os.MkdirAll(config.Conf.Internal.Route, 0777)
		if err != nil {
			return
		}
	}

	//把模板写入文件
	file, err := os.OpenFile(config.Conf.Internal.Route+"/route.go", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
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

func (i *mkService) build(md string) (err error) {
	parseDir := fmt.Sprintf("%s/%s", config.Conf.Proto.Source, md)
	entries, err := os.ReadDir(parseDir)
	if err != nil {
		return
	}

	build := &Build{
		moduleName: md,
	}

	data, err := os.ReadFile(strings.TrimSuffix(parseDir, "/") + "/doc.yaml")
	if err == nil {
		_ = yaml.Unmarshal(data, &build.Doc)
	}

	for _, v := range entries {
		str := strings.Split(v.Name(), ".")
		if len(str) != 2 || str[1] != "json" {
			continue
		}

		var parse comply.Parse
		err = parse.Parse(fmt.Sprintf("%s/%s", parseDir, v.Name()))
		if err != nil {
			return
		}

		build.parses = append(build.parses, parse)
	}

	return
}

type Build struct {
	moduleName string         //模块名称
	Doc        comply.Doc     //模块文档
	parses     []comply.Parse //解析的接口结构体
	err        error          //错误
}

func NewModule(moduleName string) (build *Build, err error) {
	parseDir := fmt.Sprintf("%s/%s", config.Conf.Proto.Source, moduleName)
	entries, err := os.ReadDir(parseDir)
	if err != nil {
		return
	}

	build = &Build{
		moduleName: moduleName,
	}

	data, err := os.ReadFile(strings.TrimSuffix(parseDir, "/") + "/doc.yaml")
	if err == nil {
		_ = yaml.Unmarshal(data, &build.Doc)
	}

	for _, v := range entries {
		str := strings.Split(v.Name(), ".")
		if len(str) != 2 || str[1] != "json" {
			continue
		}

		var parse comply.Parse
		err = parse.Parse(fmt.Sprintf("%s/%s", parseDir, v.Name()))
		if err != nil {
			return
		}

		build.parses = append(build.parses, parse)
	}

	return
}

func (b *Build) Error() error {
	return b.err
}

func (b *Build) Build() *Build {
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()

		var route = b.Doc.Route
		if route == "" {
			route = utils.MidString(b.moduleName, '-')
		}

		err := comply.NewApi(b.moduleName, route, b.parses).Build()
		if err != nil {
			b.err = errors.New(fmt.Sprintf("api: %v", err))
		}
	}()

	go func() {
		defer wg.Done()
		err := comply.NewProto(b.moduleName, b.Doc, b.parses).Build()
		if err != nil {
			b.err = errors.New(fmt.Sprintf("proto: %v", err))
		}
	}()

	go func() {
		defer wg.Done()
		err := comply.NewService(b.moduleName, b.parses).Build()
		if err != nil {
			b.err = errors.New(fmt.Sprintf("service: %v", err))
		}
	}()

	go func() {
		defer wg.Done()

		var route = b.Doc.Route
		if route == "" {
			route = utils.MidString(b.moduleName, '-')
		}

		err := comply.NewSwagger(b.moduleName, route, b.Doc, b.parses).Build()
		if err != nil {
			b.err = errors.New(fmt.Sprintf("swagger: %v", err))
		}
	}()

	wg.Wait()

	return b
}
