package main

import (
	"os"
	"github.com/ajduncan/rivulet/lib"
)

func main() {
	pwd, _ := os.Getwd()
	rivulet.NewRivulet(pwd)
}
