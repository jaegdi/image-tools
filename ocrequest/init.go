package ocrequest

import (
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
	file, err := os.OpenFile("logs_ocimagetools.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Llongfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)

	InfoLogger.Println("------------------------------------------------------------")
	InfoLogger.Println("Starting execution of clean-istags")

	EvalFlags()

	Multiproc = true
	regexValidNamespace = regexp.MustCompile(`^` + CmdParams.Family + `-..|..-` + CmdParams.Family + `-..|..-` + CmdParams.Family + `$`)

	if len(Clusters.Config["cid"].Token) < 10 {
		if err := readTokens("clusterconfig.json"); err != nil {
			log.Println("Read Clusterconfig is failed, try to get the tokens from clusters with oc login")
			for _, cluster := range Clusters.Stages {
				ocGetToken(cluster)
			}
			saveTokens(Clusters, "clusterconfig.json")
		} else {
			log.Println("Clusterconfig and Tokens loaded from clusterconfig.json")
		}
	}
	InitIsNamesForFamily(CmdParams.Family)
}
