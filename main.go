package main

import (
	"github.com/ajduncan/rivulet/lib"
	"os"
)

func main() {
	// where are we?
	pwd, _ := os.Getwd()

	// bootstrap a rivulet from here;
	rivulet.NewRivulet(pwd)
}
