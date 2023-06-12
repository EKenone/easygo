package build

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var baseDir = []string{"api", "config", "docs", "internal", "log", "pkg", "test", "cmd"}
var internalDir = []string{"app", "config", "init", "plugins", "utils"}
var appDir = []string{"api", "service", "model"}

func Mkdir() error {
	var mk = &mkdir{
		custom: "",
		path:   "",
		keep:   false,
		readme: false,
		help:   false,
		demo:   false,
	}

	fs := flag.NewFlagSet("mkdir", flag.ExitOnError)
	fs.StringVar(&mk.custom, "d", "", "custom your dir, split by ','")
	fs.StringVar(&mk.path, "p", "", "custom dir path")
	fs.BoolVar(&mk.keep, "k", false, "create .gitkeep")
	fs.BoolVar(&mk.demo, "demo", false, "create author dir demo")
	fs.BoolVar(&mk.readme, "r", false, "create README.md")
	fs.BoolVar(&mk.help, "help", false, "show help mkdir")
	fs.BoolVar(&mk.help, "h", false, "short var by help")

	// 解析命令行参数
	if err := fs.Parse(os.Args[2:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if mk.help {
		fmt.Fprintln(os.Stderr, "mkdir usage options:")
		fs.PrintDefaults()
		return nil
	}

	if mk.demo {
		return mk.buildDemo()
	}

	var dir = baseDir

	if mk.custom != "" {
		dir = strings.Split(mk.custom, ",")
	}

	return mk.create(dir)
}

type mkdir struct {
	custom string
	path   string
	keep   bool
	readme bool
	help   bool
	demo   bool
}

func (m *mkdir) create(dir []string) (err error) {
	if len(dir) == 0 {
		return
	}

	var prefix = m.prefix()

	for _, v := range dir {
		err = m.done(prefix + v)
		if err != nil {
			return
		}
	}

	return
}

func (m *mkdir) prefix() (prefix string) {
	if m.path != "" {
		prefix = fmt.Sprintf("%s/", strings.TrimSuffix(strings.TrimSuffix(m.path, "/"), "\\"))
	}
	return
}

func (m *mkdir) done(dirname string) (err error) {

	exist, err := fileExist(dirname)
	if err != nil {
		return errors.New(fmt.Sprintf("dir [%s] stat err: %v", dirname, err))
	}

	if exist {
		fmt.Printf("dir [%s] exist \n", dirname)
		return
	}

	err = os.Mkdir(dirname, 0666)
	if err != nil {
		return errors.New(fmt.Sprintf("dir [%s] create err: %v", dirname, err))
	}

	if m.keep {
		_, err = os.Create(dirname + "/.gitkeep")
		if err != nil {
			return errors.New(fmt.Sprintf("[%s] create keep err: %v", dirname, err))
		}
	}

	return
}

func (m *mkdir) buildDemo() (err error) {
	m.keep = true

	file, err := os.OpenFile(m.prefix()+"README.md", os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return
	}

	_, err = file.WriteString(mkdirReadmeContent)
	if err != nil {
		return
	}

	err = m.create(baseDir)
	if err != nil {
		return
	}

	m.keep = false
	m.path = m.prefix() + "internal"
	err = m.create(internalDir)
	if err != nil {
		return
	}

	m.path = m.prefix() + "app"
	err = m.create(appDir)
	if err != nil {
		return
	}

	return
}
