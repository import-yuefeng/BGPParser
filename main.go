package main

import (
	_ "net/http/pprof"

	"github.com/import-yuefeng/BGPParser/cmd"
)

func main() {
	cmd.Execute()
}
