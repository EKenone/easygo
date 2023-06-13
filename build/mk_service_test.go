package build

import "testing"

func TestCreateTmp(t *testing.T) {
	var mk = mkService{
		conf:   "",
		module: "",
		help:   false,
		init:   "../example/config",
	}

	t.Log(mk.createConf())
}

func TestMkService(t *testing.T) {
	var mk = mkService{
		conf:   "./service.yaml",
		module: "admin",
		help:   false,
		init:   "",
	}

	t.Log(mk.create())
}
