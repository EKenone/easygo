package main

import (
	"easygo/build"
	"fmt"
	"os"
)

func main() {
	var err error
	switch os.Args[1] {
	case "mkdir":
		err = build.Mkdir()
	case "service":
		err = build.MkService()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s Error: %v", os.Args[1], err)
		os.Exit(1)
	}
}
