package main

import (
	"os"
	"github.com/ajduncan/rivulet/lib"
)

func main() {
	// where are we?
	pwd, _ := os.Getwd()

	// bootstrap a rivulet from here;
	rivulet.NewRivulet(pwd)
}
