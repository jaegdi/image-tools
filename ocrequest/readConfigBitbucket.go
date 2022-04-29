package ocrequest

import (
	"encoding/json"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/stew/slice"
)

func getBitbucketUrl(urlpath string) string {
	return "https://bitbucket.sf-bk.de/projects/SCPDEPLOY/repos/scp-infra-config/raw/" + urlpath + "?at=refs%2Fheads%2Fmaster"
}

func getBitbucketData(filename string) []byte {
	url := getBitbucketUrl(filename)
	// DebugLogger.Println("url: ", url)
	yamlstr := getHttpAnswer(url)
	// DebugLogger.Println("yaml: ", string(yamlstr))
	yamlmap := []interface{}{}
	if err := UnmarshalMultidocYaml(yamlstr, &yamlmap); err != nil {
		ErrorLogger.Println("Unmarshal multidoc yaml err:", yamlstr)
		ErrorLogger.Println("Unmarshal multidoc yaml err:", err.Error())
	}
	// DebugLogger.Println("yamlmap: ", yamlmap)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonstr, err := json.Marshal(&yamlmap)
	// jsonstr, err := json.Marshal(yamlmap)
	if err != nil {
		ErrorLogger.Println("yamlmap:    ", yamlmap)
		ErrorLogger.Println("err:    ", err)
	}
	DebugLogger.Println("Config from scp-infra-config url:", url)
	DebugLogger.Println("Config from scp-infra-config json:", string(jsonstr))
	return jsonstr
}

func GetClusters() T_cft_clusters {
	jsonbytes := getBitbucketData("clusters.yaml")
	data := T_cft_clusters{}
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		ErrorLogger.Println("Unmarshal jsonstr\n", string(jsonbytes), "err\n", err.Error())
	}
	return data
}

func GetFamilies() T_cft_families {
	jsonbytes := getBitbucketData("families.yaml")
	data := T_cft_families{}
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		ErrorLogger.Println("Unmarshal jsonstr\n", string(jsonbytes), "err\n", err.Error())
	}
	return data
}

func GetEnvironments() T_cft_environments {
	jsonbytes := getBitbucketData("environments.yaml")
	data := T_cft_environments{}
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		ErrorLogger.Println("Unmarshal jsonstr\n", string(jsonbytes), "err\n", err.Error())
	}
	return data
}

func GetNamespaces() T_cft_namespaces {
	jsonbytes := getBitbucketData("namespaces.yaml")
	data := T_cft_namespaces{}
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		ErrorLogger.Println("Unmarshal jsonstr\n", string(jsonbytes), "err\n", err.Error())
	}
	// DebugLogger.Println("data \n", data)
	return data
}

func GetPipelines() T_cft_pipelines {
	jsonbytes := getBitbucketData("pipelines.yaml")
	data := T_cft_pipelines{}
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		ErrorLogger.Println("Unmarshal jsonstr\n", string(jsonbytes), "err\n", err.Error())
	}
	return data
}

func genClusterConfig(clusters T_cft_clusters) T_ClusterConfig {
	// func body
	cfg := T_ClusterConfig{}
	cfg.Config = map[T_clName]T_Cluster{}
	for _, cluster := range clusters {
		cl := T_Cluster{}
		cl.Name = string(cluster.Name)
		if strings.Contains(cl.Name, "scp0") {
			cl.Url = "https://api." + cl.Name + ".sf-rz.de:6443"
			cl.ConfigToolUrl = "https://scpconfig-service-master.apps." + cl.Name + ".sf-rz.de"
		} else {
			cl.Url = "https://console." + cl.Name + ".sf-rz.de:8443"
			cl.ConfigToolUrl = "https://scpconfig-service-master.default." + cl.Name + ".sf-rz.de"
		}
		cl.Token = "x"
		cfg.Config[cluster.Name] = cl
	}
	jsonstr, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		ErrorLogger.Println("\n  JsonMarshal failes\n  ERROR:  ", err)
	}
	DebugLogger.Println("\nClusterConfig:\n  ", string(jsonstr))
	return cfg
}

func genFamilyNamespacesConfig(clusters T_cft_clusters,
	families T_cft_families,
	environments T_cft_environments,
	namespaces T_cft_namespaces,
	pipelines T_cft_pipelines) T_famNs {
	// func body
	fnc := T_famNs{}
	for _, fam := range families {
		family := T_family(fam.Name)
		fnc[family] = T_familyKeys{}
		famMap := T_familyKeys{}
		famMap.ImageNamespaces = map[T_clName][]T_nsName{}
		famMap.Buildstages = []T_clName{}
		famMap.Teststages = []T_clName{}
		for _, environment := range environments {
			if environment.Family == family {
				if famMap.ImageNamespaces[environment.Cluster] == nil {
					famMap.ImageNamespaces[environment.Cluster] = []T_nsName{}
				}
				if famMap.ImageNamespaces[environment.Pre_Cluster] == nil {
					famMap.ImageNamespaces[environment.Pre_Cluster] = []T_nsName{}
				}
				for _, pipeline := range pipelines {
					if pipeline.Name == environment.Pipeline {
						if !slice.Contains(famMap.ImageNamespaces[environment.Cluster], pipeline.Image_Ns) {
							famMap.ImageNamespaces[environment.Cluster] = append(famMap.ImageNamespaces[environment.Cluster], pipeline.Image_Ns)
						}
						//  dev pipelines
						// if strings.Contains(environment.Cluster.str(), "cid-") && pipeline.Name == string(family)+"-dev-pipeline" {
						// 	if !slice.Contains(famMap.ImageNamespaces[environment.Cluster], pipeline.Deployer_Ns) {
						// 		famMap.ImageNamespaces[environment.Cluster] = append(famMap.ImageNamespaces[environment.Cluster], pipeline.Deployer_Ns)
						// 	}
						// 	if !slice.Contains(famMap.Buildstages, environment.Cluster) {
						// 		famMap.Buildstages = append(famMap.Buildstages, environment.Cluster)
						// 	}
						// }
						// for _, app := range fam.Applications {
						if strings.Contains(environment.Pipeline, "-dev-pipeline") {
							if !slice.Contains(famMap.ImageNamespaces[environment.Cluster], pipeline.Deployer_Ns) {
								famMap.ImageNamespaces[environment.Cluster] = append(famMap.ImageNamespaces[environment.Cluster], pipeline.Deployer_Ns)
							}
							if !slice.Contains(famMap.Buildstages, environment.Cluster) {
								famMap.Buildstages = append(famMap.Buildstages, environment.Cluster)
							}
						}
						// if pipeline.Name == string(family)+"-"+app+"-pipeline" {
						// 	if !slice.Contains(famMap.ImageNamespaces[environment.Cluster], pipeline.Image_Ns) {
						// 		famMap.ImageNamespaces[environment.Cluster] = append(famMap.ImageNamespaces[environment.Cluster], pipeline.Image_Ns)
						// 	}
						// 	if !slice.Contains(famMap.Buildstages, environment.Cluster) {
						// 		famMap.Buildstages = append(famMap.Buildstages, environment.Cluster)
						// 	}
						// }
						// }
					}
					if pipeline.Name == environment.Pre_Pipeline {
						if !slice.Contains(famMap.ImageNamespaces[environment.Pre_Cluster], pipeline.Image_Ns) {
							famMap.ImageNamespaces[environment.Pre_Cluster] = append(famMap.ImageNamespaces[environment.Pre_Cluster], pipeline.Image_Ns)
						}
					}
				}
				if !(strings.Contains(environment.Cluster.str(), "cid-") || strings.Contains(environment.Cluster.str(), "pro-")) {
					if !slice.Contains(famMap.Teststages, environment.Cluster) {
						famMap.Teststages = append(famMap.Teststages, environment.Cluster)
					}
				}
				if strings.Contains(environment.Cluster.str(), "pro-") {
					if !slice.Contains(famMap.Prodstages, environment.Cluster) {
						famMap.Prodstages = append(famMap.Prodstages, environment.Cluster)
					}
				}
			}
		}
		famMap.Stages = []T_clName{}
		for cl, nslist := range famMap.ImageNamespaces {
			if len(nslist) > 0 {
				if !slice.Contains(famMap.Stages, cl) {
					famMap.Stages = append(famMap.Stages, cl)
				}
			}
		}
		fnc[family] = famMap
	}
	return fnc
}
