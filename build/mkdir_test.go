package build

import "testing"

func TestCreate(t *testing.T) {
	var d = &mkdir{
		custom: "",
		path:   "",
		keep:   false,
		readme: false,
		help:   false,
		demo:   false,
	}
	t.Log(d.create(baseDir))
}

func TestBuildDemo(t *testing.T) {
	var d = &mkdir{
		custom: "",
		path:   "",
		keep:   false,
		readme: false,
		help:   false,
		demo:   false,
	}
	t.Log(d.buildDemo())
}
