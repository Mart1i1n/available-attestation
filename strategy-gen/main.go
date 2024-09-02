package main

import (
	"github.com/tsinghua-cel/strategy-gen/command/root"
)

func main() {
	root.NewRootCommand().Execute()
}
