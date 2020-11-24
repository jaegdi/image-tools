package ocrequest

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

var (
	WarningLogger       *log.Logger
	InfoLogger          *log.Logger
	ErrorLogger         *log.Logger
	Multiproc           bool
	regexValidNamespace *regexp.Regexp
)

func Init() {

	logfile, err := os.OpenFile("log_ocimagetools.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("HTTP_PROXY", "")
	InfoLogger = log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	WarningLogger = log.New(logfile, "WARNING: ", log.Ldate|log.Ltime|log.Llongfile)
	ErrorLogger = log.New(logfile, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)

	LogMsg("------------------------------------------------------------")
	LogMsg("Starting execution of clean-istags")

	EvalFlags()

	Multiproc = true
	LogMsg("disable proxy: " + fmt.Sprint(CmdParams.Options.NoProxy))
	LogMsg("Multithreading: " + fmt.Sprint(Multiproc))

	regexValidNamespace = regexp.MustCompile(`^` + CmdParams.Family + `-..|..-` + CmdParams.Family + `-..|..-` + CmdParams.Family + `$`)

	if len(Clusters.Config["cid"].Token) < 10 {
		if err := readTokens("clusterconfig.json"); err != nil {
			LogMsg("Read Clusterconfig is failed, try to get the tokens from clusters with oc login")
			for _, cluster := range Clusters.Stages {
				ocGetToken(cluster)
			}
			saveTokens(Clusters, "clusterconfig.json")
		} else {
			LogMsg("Clusterconfig and Tokens loaded from clusterconfig.json")
		}
	}
	InitIsNamesForFamily(CmdParams.Family)
}
