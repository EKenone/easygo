package build

import (
	"os"
)

func fileExist(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return
}
