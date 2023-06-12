package build

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func MkService() (err error) {
	var mk = &mkService{
		init:   "",
		conf:   "",
		module: "",
	}

	fs := flag.NewFlagSet("service", flag.ExitOnError)
	fs.StringVar(&mk.init, "i", "", "init service conf in this path")
	fs.StringVar(&mk.conf, "c", "", "set your config path when use service")
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

	return
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

	err := os.MkdirAll(path, 0666)
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

	file, err := os.OpenFile(filename, os.O_CREATE, 0666)
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
