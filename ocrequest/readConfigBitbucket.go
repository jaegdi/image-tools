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

// getBitbucketData ruft die Konfigurationsdaten von Bitbucket ab und gibt sie als JSON-Byte-Array zurück.
// Die Funktion führt folgende Schritte aus:
// 1. Erzeugt die URL für die Bitbucket-Anfrage.
// 2. Holt die YAML-Daten von der URL.
// 3. Unmarshalt die YAML-Daten in eine Map.
// 4. Konvertiert die Map in ein JSON-Byte-Array und gibt es zurück.
func getBitbucketData(filename string) []byte {
	// Erzeuge die URL für die Bitbucket-Anfrage basierend auf dem Dateinamen
	url := getBitbucketUrl(filename)
	DebugMsg("url: ", url)

	// Hole die YAML-Daten von der URL unter Verwendung des Bitbucket-Tokens
	yamlstr := getHttpAnswer(url, bitbucket_token)
	DebugMsg(filename, "yaml: ", string(yamlstr))

	// Unmarshal die YAML-Daten in eine Map
	yamlmap := []interface{}{}
	if err := UnmarshalMultidocYaml(yamlstr, &yamlmap); err != nil {
		// Logge eine Fehlermeldung, falls das Unmarshalling fehlschlägt
		ErrorMsg("Unmarshal multidoc yaml:", yamlstr)
		ErrorMsg("Unmarshal multidoc yaml:", err.Error())
	}
	DebugMsg("yamlmap: ", yamlmap)

	// Konvertiere die Map in ein JSON-Byte-Array
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonstr, err := json.Marshal(&yamlmap)
	if err != nil {
		// Logge eine Fehlermeldung, falls die Konvertierung fehlschlägt
		ErrorMsg("yamlmap:    ", yamlmap)
		ErrorMsg("err:    ", err)
	}
	VerifyMsg(filename, "jsonstr: ", GetJsonFromMap(yamlmap))

	// Logge die URL und die JSON-Daten
	DebugMsg("Config from scp-infra-config url:", url)
	DebugMsg("Config from scp-infra-config json:", string(jsonstr))

	// Gib das JSON-Byte-Array zurück
	return jsonstr
}

// GetClusters ruft die Cluster-Konfigurationsdaten von Bitbucket ab und unmarshalt sie in eine T_cft_clusters-Struktur.
func GetClusters() T_cft_clusters {
	// Hole die JSON-Daten von Bitbucket
	jsonbytes := getBitbucketData("clusters.yaml")
	data := T_cft_clusters{}
	// Unmarshal die JSON-Daten in die T_cft_clusters-Struktur
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		// Logge eine Fehlermeldung, falls das Unmarshalling fehlschlägt
		ErrorMsg("Unmarshal jsonstr:", string(jsonbytes))
		ErrorMsg("Unmarshal jsonstr err:", err.Error())
	}
	return data
}

// GetFamilies ruft die Familien-Konfigurationsdaten von Bitbucket ab und unmarshalt sie in eine T_cft_families-Struktur.
func GetFamilies() T_cft_families {
	// Hole die JSON-Daten von Bitbucket
	VerifyMsg("GetFamilies")
	jsonbytes := getBitbucketData("families.yaml")
	data := T_cft_families{}
	// Unmarshal die JSON-Daten in die T_cft_families-Struktur
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		// Logge eine Fehlermeldung, falls das Unmarshalling fehlschlägt
		ErrorMsg("Unmarshal jsonstr:", string(jsonbytes))
		ErrorMsg("Unmarshal jsonstr err:", err.Error())
	}
	return data
}

// GetEnvironments ruft die Umgebungs-Konfigurationsdaten von Bitbucket ab und unmarshalt sie in eine T_cft_environments-Struktur.
func GetEnvironments() T_cft_environments {
	// Hole die JSON-Daten von Bitbucket
	VerifyMsg("GetEnvironments")
	jsonbytes := getBitbucketData("environments.yaml")
	data := T_cft_environments{}
	// Unmarshal die JSON-Daten in die T_cft_environments-Struktur
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		// Logge eine Fehlermeldung, falls das Unmarshalling fehlschlägt
		ErrorMsg("Unmarshal jsonstr:", string(jsonbytes))
		ErrorMsg("Unmarshal jsonstr err:", err.Error())
	}
	return data
}

// GetNamespaces ruft die Namespace-Konfigurationsdaten von Bitbucket ab und unmarshalt sie in eine T_cft_namespaces-Struktur.
func GetNamespaces() T_cft_namespaces {
	// Hole die JSON-Daten von Bitbucket
	VerifyMsg("GetNamespaces")
	jsonbytes := getBitbucketData("namespaces.yaml")
	data := T_cft_namespaces{}
	// Unmarshal die JSON-Daten in die T_cft_namespaces-Struktur
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		// Logge eine Fehlermeldung, falls das Unmarshalling fehlschlägt
		ErrorMsg("Unmarshal jsonstr:", string(jsonbytes))
		ErrorMsg("Unmarshal jsonstr err:", err.Error())
	}
	// Logge die Daten für Debugging-Zwecke
	DebugMsg("data: ", data)
	return data
}

// GetPipelines ruft die Pipeline-Konfigurationsdaten von Bitbucket ab und unmarshalt sie in eine T_cft_pipelines-Struktur.
func GetPipelines() T_cft_pipelines {
	// Hole die JSON-Daten von Bitbucket
	jsonbytes := getBitbucketData("pipelines.yaml")
	data := T_cft_pipelines{}
	// Unmarshal die JSON-Daten in die T_cft_pipelines-Struktur
	if err := json.Unmarshal(jsonbytes, &data); err != nil {
		// Logge eine Fehlermeldung, falls das Unmarshalling fehlschlägt
		ErrorMsg("Unmarshal jsonstr:", string(jsonbytes))
		ErrorMsg("Unmarshal jsonstr err:", err.Error())
	}
	return data
}

// genClusterConfig generiert die Cluster-Konfiguration basierend auf den übergebenen Cluster-Daten.
func genClusterConfig(clusters T_cft_clusters) T_ClusterConfig {
	cfg := T_ClusterConfig{}
	cfg.Config = map[T_clName]T_Cluster{}
	// Iteriere über alle Cluster und setze die Konfigurationswerte
	for _, cluster := range clusters {
		cl := T_Cluster{}
		cl.Name = string(cluster.Name)
		// Setze die URL und ConfigToolUrl basierend auf dem Cluster-Namen
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
	// Konvertiere die Konfiguration in ein JSON-Byte-Array für Debugging-Zwecke
	jsonstr, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		// Logge eine Fehlermeldung, falls die Konvertierung fehlschlägt
		DebugMsg("JsonMarshal failes:", err)
	}
	// Logge die generierte Cluster-Konfiguration
	DebugMsg("ClusterConfig:", string(jsonstr))
	return cfg
}

// genFamilyNamespacesConfig generiert die Konfiguration der Familien-Namespaces basierend auf den übergebenen Cluster-, Familien-,
// Umgebungs-, Namespace- und Pipeline-Daten. Die Funktion gibt eine T_famNsList-Struktur zurück.
func genFamilyNamespacesConfig(clusters T_cft_clusters,
	families T_cft_families,
	environments T_cft_environments,
	namespaces T_cft_namespaces,
	pipelines T_cft_pipelines) T_famNsList {

	// Logge die Eingabedaten für Debugging-Zwecke
	DebugMsg("genFamilyNamespacesConfig", clusters, families, namespaces, pipelines)

	// Initialisiere die Rückgabestruktur
	fnc := T_famNsList{}
	all := T_familyKeys{
		ImageNamespaces: T_appNamespaceList{},
		Stages:          []T_clName{},
		Config:          T_ClusterConfig{},
		Buildstages:     []T_clName{},
		Teststages:      []T_clName{},
		Prodstages:      []T_clName{},
		Apps:            map[T_appName]T_appKeys{},
	}

	// Iteriere über alle Familien
	for _, fam := range families {
		// Initialisiere Maps für die verschiedenen Stages
		stages := map[T_appName][]T_clName{}
		buildstages := map[T_appName][]T_clName{}
		teststages := map[T_appName][]T_clName{}
		prodstages := map[T_appName][]T_clName{}

		// Konvertiere den Familiennamen in den Typ T_familyName
		family := T_familyName(fam.Name)

		// Initialisiere die Struktur für die Familien-Keys
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

		// Iteriere über alle Umgebungen
		for _, environment := range environments {
			if environment.Family == family || CmdParams.Family == "all" {
				// Initialisiere die ImageNamespaces für die Cluster und Pre-Cluster
				if famMap.ImageNamespaces[environment.Cluster] == nil {
					famMap.ImageNamespaces[environment.Cluster] = []T_nsName{}
					if all.ImageNamespaces[environment.Cluster] == nil {
						all.ImageNamespaces[environment.Cluster] = []T_nsName{}
					}
				}
				if famMap.ImageNamespaces[environment.Pre_Cluster] == nil {
					famMap.ImageNamespaces[environment.Pre_Cluster] = []T_nsName{}
					if all.ImageNamespaces[environment.Pre_Cluster] == nil {
						all.ImageNamespaces[environment.Pre_Cluster] = []T_nsName{}
					}
				}

				// Iteriere über alle Pipelines
				for _, pipeline := range pipelines {
					if pipeline.Name == environment.Pipeline {
						// Füge die ImageNamespaces hinzu, falls sie nicht bereits enthalten sind
						if !slice.Contains(famMap.ImageNamespaces[environment.Cluster], pipeline.Image_Ns) {
							famMap.ImageNamespaces[environment.Cluster] = append(famMap.ImageNamespaces[environment.Cluster], pipeline.Image_Ns)
							if !slice.Contains(all.ImageNamespaces[environment.Cluster], pipeline.Image_Ns) {
								all.ImageNamespaces[environment.Cluster] = append(all.ImageNamespaces[environment.Cluster], pipeline.Image_Ns)
							}
						}
						// Füge die Deployer_Ns hinzu, falls es sich um eine Dev-Pipeline handelt
						if strings.Contains(environment.Pipeline, "-dev-pipeline") {
							if !slice.Contains(famMap.ImageNamespaces[environment.Cluster], pipeline.Deployer_Ns) {
								famMap.ImageNamespaces[environment.Cluster] = append(famMap.ImageNamespaces[environment.Cluster], pipeline.Deployer_Ns)
								if !slice.Contains(all.ImageNamespaces[environment.Cluster], pipeline.Deployer_Ns) {
									all.ImageNamespaces[environment.Cluster] = append(all.ImageNamespaces[environment.Cluster], pipeline.Deployer_Ns)
								}
							}
							// Füge den Cluster zu den Buildstages hinzu
							if !slice.Contains(famMap.Buildstages, environment.Cluster) {
								famMap.Buildstages = append(famMap.Buildstages, environment.Cluster)
								if !slice.Contains(all.Buildstages, environment.Cluster) {
									all.Buildstages = append(all.Buildstages, environment.Cluster)
								}
							}
						}
					}
					if pipeline.Name == environment.Pre_Pipeline {
						// Füge die ImageNamespaces für den Pre-Cluster hinzu
						if !slice.Contains(famMap.ImageNamespaces[environment.Pre_Cluster], pipeline.Image_Ns) {
							famMap.ImageNamespaces[environment.Pre_Cluster] = append(famMap.ImageNamespaces[environment.Pre_Cluster], pipeline.Image_Ns)
							all.ImageNamespaces[environment.Pre_Cluster] = append(all.ImageNamespaces[environment.Pre_Cluster], pipeline.Image_Ns)
						}
					}
				}

				// Füge den Cluster zu den Stages hinzu
				if !slice.Contains(famMap.Stages, environment.Cluster) {
					famMap.Stages = append(famMap.Stages, environment.Cluster)
					if !slice.Contains(all.Stages, environment.Cluster) {
						all.Stages = append(all.Stages, environment.Cluster)
					}

				}

				// Füge den Cluster zu den Teststages hinzu, falls es sich nicht um einen CID- oder Prod-Cluster handelt
				if !(strings.Contains(environment.Cluster.str(), "cid-") || strings.Contains(environment.Cluster.str(), "pro-")) {
					if !slice.Contains(famMap.Teststages, environment.Cluster) {
						famMap.Teststages = append(famMap.Teststages, environment.Cluster)
						if !slice.Contains(all.Teststages, environment.Cluster) {
							all.Teststages = append(all.Teststages, environment.Cluster)
						}
					}
				}

				// Füge den Cluster zu den Prodstages hinzu, falls es sich um einen Prod-Cluster handelt
				if strings.Contains(environment.Cluster.str(), "pro-") {
					if !slice.Contains(famMap.Prodstages, environment.Cluster) {
						famMap.Prodstages = append(famMap.Prodstages, environment.Cluster)
						if !slice.Contains(all.Prodstages, environment.Cluster) {
							all.Prodstages = append(all.Prodstages, environment.Cluster)
						}
					}
				}

				// Sammle Informationen über die Anwendungen
				for _, app := range fam.Applications {
					appname := T_appName(app)
					app_appns := T_appNamespaceList{}
					app_buildns := T_appNamespaceList{}

					// Initialisiere die Stages-Maps für die Anwendungen
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

					// Iteriere über alle Namespaces
					for _, namespace := range namespaces {
						// Überprüfe, ob die Anwendung im aktuellen Namespace enthalten ist und
						// ob die Umgebung des Namespaces mit der aktuellen Umgebung übereinstimmt
						if slice.Contains(namespace.Applications, app) && namespace.Environment == environment.Name {
							// Füge den Cluster zu den Stages der Anwendung hinzu
							if !slice.Contains(stages[appname], environment.Cluster) {
								stages[appname] = append(stages[appname], environment.Cluster)
							}

							// Verarbeite CID-Namespaces
							if strings.Contains(string(environment.Cluster), "cid-") {
								if app_buildns[environment.Cluster] == nil {
									app_buildns[environment.Cluster] = []T_nsName{}
								}
								app_buildns[environment.Cluster] = append(app_buildns[environment.Cluster], namespace.Name)
								if !slice.Contains(buildstages[appname], environment.Cluster) {
									buildstages[appname] = append(buildstages[appname], environment.Cluster)
								}
							} else {
								// Verarbeite Dev-, Int-, PPR- und VPT- Namespaces
								if strings.Contains(string(environment.Cluster), "dev-") ||
									strings.Contains(string(environment.Cluster), "ppr-") ||
									strings.Contains(string(environment.Cluster), "vpt-") {
									if !slice.Contains(teststages[appname], environment.Cluster) {
										teststages[appname] = append(teststages[appname], environment.Cluster)
									}
								} else {
									// Verarbeite Prod-Namespaces
									if !slice.Contains(prodstages[appname], environment.Cluster) {
										prodstages[appname] = append(prodstages[appname], environment.Cluster)
									}
								}
								app_appns[environment.Cluster] = append(app_appns[environment.Cluster], namespace.Name)
							}
						} else {
							if namespace.Environment == environment.Name && strings.Contains(namespace.Name.str(), environment.Name) {
								if !slice.Contains(app_appns[environment.Cluster], namespace.Name) {
									app_appns[environment.Cluster] = append(app_appns[environment.Cluster], namespace.Name)
								}
							}
						}
					}
					if len(app_appns) == 0 {
						app_buildns[environment.Cluster] = append(app_buildns[environment.Cluster], fam.Metadata.Image_ns)
						app_buildns[environment.Cluster] = append(app_buildns[environment.Cluster], fam.Metadata.Pipeline_namespaces_cid...)
					}
					// Füge die Anwendungsinformationen zur Familien-Map hinzu
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
				MergoNestedMaps(&all.Apps, famMap.Apps)
			}
		}
		// Füge die Familien-Map zur Rückgabestruktur hinzu
		fnc[T_familyName(fam.Name)] = famMap
	}
	delete(all.ImageNamespaces, "")
	VerifyMsg("all", GetJsonFromMap(all))
	fnc["all"] = all
	return fnc
}
