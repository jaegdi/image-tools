package main

import (
	. "image-tool/ocrequest"
	_ "net/http/pprof"
)

func init() {
	Init()
}

func main() {

	if CmdParams.Options.ServerMode {
		StartServer()
	} else {
		_ = CmdlineMode()
	}

}
