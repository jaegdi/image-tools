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
	yamlstr := getHttpAnswer(url, bitbucket_token)
	// DebugLogger.Println("yaml: ", string(yamlstr))
	yamlmap := []interface{}{}
	if err := UnmarshalMultidocYaml(yamlstr, &yamlmap); err != nil {
		ErrorLogger.Println("Unmarshal multidoc yaml:", yamlstr)
		ErrorLogger.Println("Unmarshal multidoc yaml:", err.Error())
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
		ErrorLogger.Println("Unmarshal jsonstr:", string(jsonbytes))
		ErrorLogger.Println("Unmarshal jsonstr err:", err.Error())
	}
	return data
}

func GetFamilies() T_cft_families {
	jsonbytes := getBitbucketData("families.yaml")
	data := T_cft_families{}
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		ErrorLogger.Println("Unmarshal jsonstr:", string(jsonbytes))
		ErrorLogger.Println("Unmarshal jsonstr err:", err.Error())
	}
	return data
}

func GetEnvironments() T_cft_environments {
	jsonbytes := getBitbucketData("environments.yaml")
	data := T_cft_environments{}
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		ErrorLogger.Println("Unmarshal jsonstr:", string(jsonbytes))
		ErrorLogger.Println("Unmarshal jsonstr err:", err.Error())
	}
	return data
}

func GetNamespaces() T_cft_namespaces {
	jsonbytes := getBitbucketData("namespaces.yaml")
	data := T_cft_namespaces{}
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		ErrorLogger.Println("Unmarshal jsonstr:", string(jsonbytes))
		ErrorLogger.Println("Unmarshal jsonstr err:", err.Error())
	}
	// DebugLogger.Println("data: ", data)
	return data
}

func GetPipelines() T_cft_pipelines {
	jsonbytes := getBitbucketData("pipelines.yaml")
	data := T_cft_pipelines{}
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		ErrorLogger.Println("Unmarshal jsonstr:", string(jsonbytes))
		ErrorLogger.Println("Unmarshal jsonstr err:", err.Error())
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
		ErrorLogger.Println("JsonMarshal failes:", err)
	}
	DebugLogger.Println("ClusterConfig:", string(jsonstr))
	return cfg
}

func genFamilyNamespacesConfig(clusters T_cft_clusters,
	families T_cft_families,
	environments T_cft_environments,
	namespaces T_cft_namespaces,
	pipelines T_cft_pipelines) T_famNsList {
	// func body
	fnc := T_famNsList{}
	for _, fam := range families {
		stages := map[T_appName][]T_clName{}
		buildstages := map[T_appName][]T_clName{}
		teststages := map[T_appName][]T_clName{}
		prodstages := map[T_appName][]T_clName{}
		family := T_familyName(fam.Name)
		fnc[family] = T_familyKeys{}
		famMap := T_familyKeys{
			ImageNamespaces: T_appNamespaceList{},
			Stages:          []T_clName{},
			Config:          T_ClusterConfig{},
			Buildstages:     []T_clName{},
			Teststages:      []T_clName{},
			Prodstages:      []T_clName{},
			Apps:            map[T_appName]T_appKeys{},
		}
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
						if strings.Contains(environment.Pipeline, "-dev-pipeline") {
							if !slice.Contains(famMap.ImageNamespaces[environment.Cluster], pipeline.Deployer_Ns) {
								famMap.ImageNamespaces[environment.Cluster] = append(famMap.ImageNamespaces[environment.Cluster], pipeline.Deployer_Ns)
							}
							if !slice.Contains(famMap.Buildstages, environment.Cluster) {
								famMap.Buildstages = append(famMap.Buildstages, environment.Cluster)
							}
						}
					}
					if pipeline.Name == environment.Pre_Pipeline {
						if !slice.Contains(famMap.ImageNamespaces[environment.Pre_Cluster], pipeline.Image_Ns) {
							famMap.ImageNamespaces[environment.Pre_Cluster] = append(famMap.ImageNamespaces[environment.Pre_Cluster], pipeline.Image_Ns)
						}
					}
				}
				if !slice.Contains(famMap.Stages, environment.Cluster) {
					famMap.Stages = append(famMap.Stages, environment.Cluster)
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
				// Collect applications info
				for _, app := range fam.Applications {
					appname := T_appName(app)
					// appnslist := T_appNsList{}
					// appkeys := T_appKeys{}
					app_appns := T_appNamespaceList{}
					app_buildns := T_appNamespaceList{}
					if stages[appname] == nil {
						stages[appname] = []T_clName{}
					}
					if buildstages[appname] == nil {
						buildstages[appname] = []T_clName{}
					}
					if teststages[appname] == nil {
						teststages[appname] = []T_clName{}
					}
					if prodstages[appname] == nil {
						prodstages[appname] = []T_clName{}
					}
					for _, namespace := range namespaces {
						if slice.Contains(namespace.Applications, app) && namespace.Environment == environment.Name {
							// every namespace added to stages
							if !slice.Contains(stages[appname], environment.Cluster) {
								stages[appname] = append(stages[appname], environment.Cluster)
							}
							if strings.Contains(string(environment.Cluster), "cid-") {
								// if cid namespace
								if app_buildns[environment.Cluster] == nil {
									app_buildns[environment.Cluster] = []T_nsName{}
								}
								// add to app_buildns
								app_buildns[environment.Cluster] = append(app_buildns[environment.Cluster], namespace.Name)
								// add to buildstages
								if !slice.Contains(buildstages[appname], environment.Cluster) {
									buildstages[appname] = append(buildstages[appname], environment.Cluster)
								}
							} else {
								if strings.Contains(string(environment.Cluster), "dev-") || strings.Contains(string(environment.Cluster), "int-") || strings.Contains(string(environment.Cluster), "ppr-") {
									// if int- or ppr- namespace, add to teststages
									if !slice.Contains(teststages[appname], environment.Cluster) {
										teststages[appname] = append(teststages[appname], environment.Cluster)
									}
								} else {
									// must be prod namespace, add to prodstages
									if !slice.Contains(prodstages[appname], environment.Cluster) {
										prodstages[appname] = append(prodstages[appname], environment.Cluster)
									}
								}
								// add to app_appns
								app_appns[environment.Cluster] = append(app_appns[environment.Cluster], namespace.Name)
							}
						}
					}
					famMap.Apps[appname] = T_appKeys{
						Namespaces: T_appNamespaces{
							Buildnamespaces: app_buildns,
							Appnamespaces:   app_appns,
						},
						Config:      T_ClusterConfig{},
						Stages:      stages[appname],
						Buildstages: buildstages[appname],
						Teststages:  teststages[appname],
						Prodstages:  prodstages[appname],
					}
				}
			}
		}
		fnc[T_familyName(fam.Name)] = famMap
	}
	return fnc
}
