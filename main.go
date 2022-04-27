package main

import (
	"fmt"
	. "image-tools/ocrequest"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
)

func init() {
	Init()
}

func main() {
	// --%<---- Profiling
	var wg sync.WaitGroup
	if CmdParams.Options.Profiler {
		f, perr := os.Create("cpu.pprof")
		if perr != nil {
			log.Fatal(perr)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		go func() {
			runtime.SetBlockProfileRate(1)
			InfoLogger.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	// --%<---- Profiling
	if CmdParams.Options.ServerMode {
		fmt.Println("Will be startet in server mode")

	} else {
		CmdlineMode()
	}
	// --%<---- Profiling
	if CmdParams.Options.Profiler {
		wg.Add(1)
		wg.Wait()
	}
	// --%<---- Profiling
}
