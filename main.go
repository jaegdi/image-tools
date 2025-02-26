package main

import (
	. "image-tool/ocrequest"
	_ "net/http/pprof"
)

func init() {
	Init()
}

func main() {
	// If CmdParams.Clusters is empty, set it to all clusters
	// if len(CmdParams.Cluster) == 0 {
	// 	CmdParams.Cluster = Clusters.list() //.clusterNames()
	// }
	// fmt.Println("CmdParams.Cluster: ", CmdParams.Cluster)

	if CmdParams.Options.ServerMode {
		StartServer()
	} else {
		_ = CmdlineMode()
	}

}
