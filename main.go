package main

import (
	// log "github.com/sirupsen/logrus"
	_ "net/http/pprof"

	cmd "github.com/import-yuefeng/BGPParser/cmd"
)

func main() {
	cmd.Execute()
}
