package ioutils

import (
	"io/ioutil"
	"os"

	berror "github.com/maplain/bulletin/pkg/error"
)

func ReadFile(name string) string {
	dat, err := ioutil.ReadFile(name)
	berror.CheckError(err)
	return string(dat)
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		berror.CheckError(err)
	}
}

func CreateFileIfNotExist(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err = os.Create(filename)
		berror.CheckError(err)
	}
}
