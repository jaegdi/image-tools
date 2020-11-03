package ocrequest

import (
	"encoding/json"
	"flag"
	"github.com/imdario/mergo"
	"log"
)

type T_famNs map[string][]string
type T_flags map[string]interface{}
type T_flagOut map[string]bool
type T_flagFilt map[string]string

func EvalFlags(familyNamespaces T_famNs) T_flags {
	// Global Flags
	clusterPtr := flag.String("cluster", "", "shortname of cluster, eg. cid,int, ppr or pro")
	familyPtr := flag.String("family", "", "family name, eg. pkp, aps, ssp or fpc ")
	namespacePtr := flag.String("namespace", "", "namespace to look for istags")

	// Output flags
	isPtr := flag.Bool("is", false, "output of imageStreams as json")
	istagPtr := flag.Bool("istag", false, "output of imageStreamTags as json")
	shaPtr := flag.Bool("sha", false, "output of Sha's as json")
	allPtr := flag.Bool("all", false, "output all imageStreams imageStreamTags and Sha's as json")

	// Filter flags
	isnamePtr := flag.String("isname", "", "filter output of one imageStream as json, eg. -is=wvv-service")
	istagnamePtr := flag.String("istagname", "", "filter output of one imageStreamTag as json")
	shanamePtr := flag.String("shaname", "", "filter output of one Sha as json")
	flag.Parse()

	// define map with all flags
	flags := T_flags{
		"cluster": string(*clusterPtr),
		"family":  string(*familyPtr),
		// "namespace": *namespacePtr,
		"output": T_flagOut{
			"is":    *isPtr,
			"istag": *istagPtr,
			"sha":   *shaPtr,
			"all":   *allPtr,
		},
		"filter": T_flagFilt{
			"is":        *isnamePtr,
			"istag":     *istagnamePtr,
			"sha":       *shanamePtr,
			"namespace": *namespacePtr,
		},
	}

	log.Println("Cluster: ", flags["cluster"])
	log.Println("Family: ", flags["family"])
	log.Println("Output Is: ", flags["output"].(T_flagOut)["is"])
	log.Println("Output Istag: ", flags["output"].(T_flagOut)["istag"])
	log.Println("Output Sha: ", flags["output"].(T_flagOut)["sha"])
	log.Println("Output All: ", flags["output"].(T_flagOut)["all"])
	log.Println("Filter Namespace: ", flags["filter"].(T_flagFilt)["namespace"])
	log.Println("Filter IsName: ", flags["filter"].(T_flagFilt)["isname"])
	log.Println("Filter IstagName: ", flags["filter"].(T_flagFilt)["istagname"])
	log.Println("Filter ShaName: ", flags["filter"].(T_flagFilt)["shaname"])

	if flags["cluster"] == "" {
		exitWithError("a shortname for cluster must given like: '-cluster=cid'. Is now: " + flags["cluster"].(string))
	}
	if flags["family"] == "" {
		exitWithError("a name for family must given like: '-family=pkp'")
	}
	if familyNamespaces[flags["family"].(string)] == nil {
		exitWithError("Family " + flags["family"].(string) + " is not defined")
	}

	foundNamespace := false
	for _, v := range familyNamespaces[flags["family"].(string)] {
		if flags["filter"].(T_flagFilt)["namespace"] == v {
			foundNamespace = true
		}
	}
	if !foundNamespace && !(flags["filter"].(T_flagFilt)["namespace"] == "") {
		exitWithError("Namespace " + flags["filter"].(T_flagFilt)["namespace"] +
			" is no image namespace for family " + flags["family"].(string))
	}

	if !(*isPtr || *istagPtr || *shaPtr || *allPtr) {
		exitWithError("As least one of the output flags mus set")
	}
	return flags
}

func FilterAllIstags(list T_result, flags T_flags) T_result {
	outputflags := flags["output"].(T_flagOut)
	if !outputflags["all"] {
		if !outputflags["is"] {
			delete(list, "is")
		}
		if !outputflags["istag"] {
			delete(list, "istag")
		}
		if !outputflags["sha"] {
			delete(list, "sha")
		}
	}
	return list
}

// Generate json output depending on the commadline flags
func GetJsonFromMap(list interface{}, flags T_flags) string {
	jsonBytes, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Println(err.Error())
	}
	return string(jsonBytes)
}

// Generate map of all istags and return json string with the results
func GetAllIstagsForFamilyInCluster(flags T_flags, familyNamespaces T_famNs) T_result {
	family := flags["family"].(string)
	cluster := flags["cluster"].(string) + "-apc0"
	namespace := flags["filter"].(T_flagFilt)["namespace"]

	token, err := OcLogin(cluster)
	if err != nil {
		exitWithError("Login failed")
	}

	var result = map[string]interface{}{}
	if namespace == "" {
		for _, ns := range familyNamespaces[family] {
			r := OcGetAllIstagsOfNamespace(result, cluster, token, ns)
			if err := mergo.Merge(&result, r); err != nil {
				log.Println("ERROR: " + "merge mySha to resultSha" + ": failed: " + err.Error())
			}

		}
	} else {
		result = OcGetAllIstagsOfNamespace(result, cluster, token, namespace)
	}
	return result
}

func GetUsedIstagsForFamilyInCluster(flags T_flags, familyNamespaces T_famNs) map[string]map[string]interface{} {
	family := flags["family"].(string)
	cluster := flags["cluster"].(string) + "-apc0"
	namespace := flags["filter"].(T_flagFilt)["namespace"]
	token, err := OcLogin(cluster)
	if err != nil {
		exitWithError("Login failed")
	}
	// Println(cluster, token, namespace)

	var result = map[string]map[string]interface{}{}
	if namespace == "" {
		for _, ns := range familyNamespaces[family] {
			r := ocGetAllUsedIstagsOfNamespace(cluster, token, ns)
			if err := mergo.Merge(&result, r); err != nil {
				log.Println("ERROR: " + "merge mySha to resultSha" + ": failed: " + err.Error())
			}
		}
	} else {
		result = ocGetAllUsedIstagsOfNamespace(cluster, token, namespace)
	}
	return result
}
