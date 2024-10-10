package main

import (
	. "image-tool/ocrequest"
	// _ "net/http/pprof"
)

func init() {
	Init()
}

func main() {
	// // --%<---- Begin Profiling
	// var wg sync.WaitGroup
	// if CmdParams.Options.Profiler {
	// 	f, perr := os.Create("cpu.pprof")
	// 	if perr != nil {
	// 		log.Fatal(perr)
	// 	}
	// 	pprof.StartCPUProfile(f)
	// 	defer pprof.StopCPUProfile()
	// 	go func() {
	// 		runtime.SetBlockProfileRate(1)
	// 		InfoMsg(http.ListenAndServe("localhost:6060", nil))
	// 	}()
	// }
	// // --%<---- End Profiling

	if CmdParams.Options.ServerMode {
		StartServer()
	} else {
		_ = CmdlineMode()
	}

	// // --%<---- Begin Profiling
	// if CmdParams.Options.Profiler {
	// 	wg.Add(1)
	// 	wg.Wait()
	// }
	// // --%<---- End Profiling
}
