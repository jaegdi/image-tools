package ocrequest

import (
	// "log"
	// "encoding/json"
	// "github.com/jedib0t/go-pretty/v6/table"
	// "log"
	// "os"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type T_completeResults struct {
	AllIstags    T_ResultExistingIstagsOverAllClusters
	UsedIstags   T_usedIstagsResult
	UnUsedIstags T_unUsedIstagsResult
}

type T_allImages map[string]interface{}

type T_family string

func (c T_family) str() string {
	return string(c)
}

//                                 family
type T_completeResultsFamilies map[T_family]T_completeResults

type T_isName string

func (c T_isName) str() string {
	return string(c)
}

type T_tagName string

func (c T_tagName) str() string {
	return string(c)
}

type T_istagName string

func (c T_istagName) str() string {
	return string(c)
}

func (istag T_istagName) split() (T_isName, T_tagName) {
	istagParts := strings.Split(istag.str(), ":")
	var is T_isName
	var tag T_tagName
	if len(istagParts) < 2 {
		is = T_isName(istagParts[0])
		tag = T_tagName("")
	} else {
		is = T_isName(istagParts[0])
		tag = T_tagName(istagParts[1])
	}
	return is, tag
}

type T_shaName string

func (c T_shaName) str() string {
	return string(c)
}

// getExistingIstags.go
//                     is        image
type T_shaStreams map[T_isName]map[T_shaName]T_resIstag

func (a T_shaStreams) Add(is T_isName, image T_shaName, istag T_istag) {
	if a == nil {
		a = T_shaStreams{}
	}
	if a[is] == nil {
		a[is] = map[T_shaName]T_resIstag{}
	}
	if a[is][image] == nil {
		a[is][image] = T_resIstag{}
	}
	istagname := T_istagName(istag.Imagestream.str() + ":" + istag.Tagname.str())
	if a[is][image][istagname] == nil {
		a[is][image][istagname] = map[T_nsName]T_istag{}
	}
	a[is][image][istagname][istag.Namespace] = istag
}

type T_Istags_List map[T_istagName]bool
type T_shaNames map[T_shaName]T_Istags_List

func (a T_shaNames) Add(sha T_shaName, istag T_istagName) {
	if a == nil {
		a = T_shaNames{}
	}
	if a[sha] == nil {
		a[sha] = T_Istags_List{}
	}
	a[sha][istag] = true
}

type T_istagBuildLabels struct {
	CommitAuthor   string `json:"io.openshift.build.commit.author,omitempty"`
	CommitDate     string `json:"io.openshift.build.commit.date,omitempty"`
	CommitId       string `json:"io.openshift.build.commit.id,omitempty"`
	CommitRef      string `json:"io.openshift.build.commit.ref,omitempty"`
	CommitVersion  string `json:"io.openshift.build.commit.version,omitempty"`
	IsProdImage    string `json:"io.openshift.build.isProdImage,omitempty"`
	BuildName      string `json:"io.openshift.build.name,omitempty"`
	BuildNamespace string `json:"io.openshift.build.namespace,omitempty"`
}

func (b T_istagBuildLabels) Values() interface{} {
	l := []interface{}{}
	v := reflect.ValueOf(b)
	for i := 0; i < v.NumField(); i++ {
		l = append(l, v.Field(i).String())
	}
	return l
}
func (b T_istagBuildLabels) Names() interface{} {
	l := []interface{}{}
	v := reflect.ValueOf(b)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		l = append(l, typeOfS.Field(i).Name)
	}
	return l
}
func (b T_istagBuildLabels) GetVal(s string) string {
	r := reflect.ValueOf(b)
	f := reflect.Indirect(r).FieldByName(s)
	return string(f.String())
}
func (buildLabels *T_istagBuildLabels) Set(buildLabelsMap map[string]interface{}) {
	// buildLabelsJSON := []byte(GetJsonFromMap(buildLabelsMap))
	// if err := json.Unmarshal(buildLabelsJSON, &buildLabels); err != nil {
	// 	ErrorLogger.Println("Unmarshal unescaped String", err)
	// }
	if buildLabelsMap["io.openshift.build.commit.author"] != nil {
		buildLabels.CommitAuthor = buildLabelsMap["io.openshift.build.commit.author"].(string)
	}
	if buildLabelsMap["io.openshift.build.commit.date"] != nil {
		buildLabels.CommitDate = buildLabelsMap["io.openshift.build.commit.date"].(string)
	}
	if buildLabelsMap["io.openshift.build.commit.id"] != nil {
		buildLabels.CommitId = buildLabelsMap["io.openshift.build.commit.id"].(string)
	}
	if buildLabelsMap["io.openshift.build.commit.ref"] != nil {
		buildLabels.CommitRef = buildLabelsMap["io.openshift.build.commit.ref"].(string)
	}
	if buildLabelsMap["io.openshift.build.commit.version"] != nil {
		buildLabels.CommitVersion = buildLabelsMap["io.openshift.build.commit.version"].(string)
	}
	if buildLabelsMap["io.openshift.build.isProdImage"] != nil {
		buildLabels.IsProdImage = buildLabelsMap["io.openshift.build.isProdImage"].(string)
	}
	if buildLabelsMap["io.openshift.build.name"] != nil {
		buildLabels.BuildName = buildLabelsMap["io.openshift.build.name"].(string)
	}
	if buildLabelsMap["io.openshift.build.namespace"] != nil {
		buildLabels.BuildNamespace = buildLabelsMap["io.openshift.build.namespace"].(string)
	}

}

// type T_istag map[string]interface{}
type T_istag struct {
	Imagestream T_isName
	Tagname     T_tagName
	Namespace   T_nsName
	// Link        string
	Date      string
	AgeInDays int
	Image     T_shaName
	Build     T_istagBuildLabels
}

func ToGenericArray(arr ...interface{}) []interface{} {
	return arr
}

func (c T_istag) Values() interface{} {
	l := []interface{}{}
	v := reflect.ValueOf(c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		switch typeOfS.Field(i).Name {
		case "Build":
			v := reflect.ValueOf(c.Build)
			for i := 0; i < v.NumField(); i++ {
				l = append(l, v.Field(i).String())
			}
		case "AgeInDays":
			s := fmt.Sprintf("%5d", int(v.Field(i).Int()))
			l = append(l, s)
		default:
			l = append(l, v.Field(i).String())
		}
	}
	return l
}

func (c T_istag) Names() interface{} {
	l := []interface{}{}
	v := reflect.ValueOf(c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		switch typeOfS.Field(i).Name {
		case "Build":
			v := reflect.ValueOf(c.Build)
			typeOfS := v.Type()
			for i := 0; i < v.NumField(); i++ {
				l = append(l, typeOfS.Field(i).Name)
			}
		default:
			l = append(l, typeOfS.Field(i).Name)
		}
	}
	return l
}

type T_sha struct {
	Istags      T_Istags_List
	Imagestream T_isName
	Namespace   T_nsName
	// Link        string
	Date      string
	AgeInDays int
}

type T_isShaTagnames map[T_istagName]interface{}
type T_is map[T_shaName]T_isShaTagnames

//                   istag     namespace
type T_resIstag map[T_istagName]map[T_nsName]T_istag
type T_resIs map[T_isName]T_is
type T_resSha map[T_shaName]map[T_istagName]T_sha

type T_resReport struct {
	Anz_ImageStreamTags int
	Anz_Images          int
	Anz_ImageStreams    int
}

type T_result struct {
	Is     T_resIs
	Istag  T_resIstag
	Image  T_resSha
	Report T_resReport
}

//------------------------------------------

// get UsedIstags.go
type T_DcResults map[string]interface{}
type T_JobResults map[string]interface{}
type T_CronjobResults map[string]interface{}
type T_Results map[string]interface{}

type T_runningObjects struct {
	Dc      T_DcResults
	Job     T_JobResults
	Cronjob T_CronjobResults
	Pod     T_Results
}

type T_usedIstag struct {
	Cluster         T_clName
	UsedInNamespace T_nsName
	FromNamespace   T_nsName
	AgeInDays       int
	Image           T_shaName
}

func (c T_usedIstag) Values() interface{} {
	l := []interface{}{}
	v := reflect.ValueOf(c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		switch typeOfS.Field(i).Name {
		case "AgeInDays":
			s := fmt.Sprintf("%5d", int(v.Field(i).Int()))
			l = append(l, s)
		default:
			l = append(l, v.Field(i).String())
		}
	}
	return l
}

func (c T_usedIstag) Names() interface{} {
	l := []interface{}{}
	v := reflect.ValueOf(c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		l = append(l, typeOfS.Field(i).Name)
	}
	return l
}

type T_unUsedIstag struct {
	Cluster T_clName
}

func (c T_unUsedIstag) Values() interface{} {
	l := []interface{}{}
	v := reflect.ValueOf(c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		switch typeOfS.Field(i).Name {
		case "AgeInDays":
			s := fmt.Sprintf("%5d", int(v.Field(i).Int()))
			l = append(l, s)
		default:
			l = append(l, v.Field(i).String())
		}
	}
	return l
}

func (c T_unUsedIstag) Names() interface{} {
	l := []interface{}{}
	v := reflect.ValueOf(c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		l = append(l, typeOfS.Field(i).Name)
	}
	return l
}

// usedIstagsResult[Is][Tag]T_usedIstag
type T_usedIstagsResult map[T_isName]map[T_tagName][]T_usedIstag

type T_unUsedIstagsResult map[T_istagName]T_unUsedIstag

// IsNamesForFamily[family][is]true
type T_IsNamesForFamily map[T_family]map[T_isName]bool

//------------------------------------------

type T_clName string

// str convert T_clName to string
func (c T_clName) str() string {
	return string(c)
}

type T_clNames []T_clName

func (c T_clName) list() T_clNames {
	clusters := strings.Split(string(c), ",")
	clusterlist := T_clNames{}
	for _, cl := range clusters {
		x := T_clName(cl)
		clusterlist = append(clusterlist, x)
	}
	print(clusterlist)
	return clusterlist
}

func (clusters T_clNames) contains(c T_clName) bool {
	for _, v := range clusters {
		if v == c {
			return true
		}
	}
	return false
}

type T_nsName string

// str convert T_nsName to string
func (c T_nsName) str() string {
	return string(c)
}

type T_familyKeys struct {
	ImageNamespaces map[T_clName][]T_nsName
	Stages          []T_clName
	Config          map[T_clName]T_Cluster `json:"config.[],omitempty"`
	Buildstages     []T_clName
	Teststages      []T_clName
	Prodstages      []T_clName
}

func (c T_familyKeys) clusterList() []string {
	clusters := []string{}
	for _, fam := range c.Stages {
		clusters = append(clusters, string(fam))
	}
	return clusters
}

func (c T_familyKeys) clusterListStr() string {
	return strings.Join(c.clusterList(), `, `)
}

//               family   cl-ns          cluster   namespaces
type T_famNs map[T_family]T_familyKeys

func (c T_famNs) familyList() []string {
	families := []string{}
	for fam := range c {
		families = append(families, string(fam))
	}
	return families
}

func (c T_famNs) familyListStr() string {
	return strings.Join(c.familyList(), `, `)
}

//------------------------------------------

type T_flagOut struct {
	Is     bool
	Istag  bool
	Image  bool
	Used   bool
	UnUsed bool
	All    bool
}
type T_flagFilt struct {
	Isname    T_isName
	Istagname T_istagName
	Tagname   T_tagName
	Imagename T_shaName
	Namespace T_nsName
}

type T_flagFiltRegexp struct {
	Isname    *regexp.Regexp
	Istagname *regexp.Regexp
	Tagname   *regexp.Regexp
	Namespace *regexp.Regexp
}

type T_flagDeleteOpts struct {
	Pattern     string
	MinAge      int
	NonBuild    bool
	Snapshots   bool
	PruneImages bool
	Confirm     bool
}

type T_flagOpts struct {
	InsecureSSL bool
	OcClient    bool
	NoProxy     bool
	Socks5Proxy string
	Profiler    bool
	NoLog       bool
	Debug       bool
	Verify      bool
}
type T_flags struct {
	Cluster    T_clNames
	Token      string
	Family     T_family
	Json       bool
	Yaml       bool
	Csv        bool
	CsvFile    string
	Delete     bool
	Html       bool
	Table      bool
	TabGroup   bool
	Output     T_flagOut
	Filter     T_flagFilt
	FilterReg  T_flagFiltRegexp
	DeleteOpts T_flagDeleteOpts
	Options    T_flagOpts
}

//------------------------------------------

type T_Cluster struct {
	Url           string `json:"url,omitempty"`
	Name          string `json:"name,omitempty"`
	Token         string `json:"token,omitempty"`
	ConfigToolUrl string `json:"configtoolurl,omitempty"`
}

type T_ClusterConfig struct {
	// Config map[T_clName]T_Cluster `json:"Config.[],omitempty"`
	Config map[T_clName]T_Cluster `json:"config,omitempty"`
}

//------------------------------------------

type T_csvLine []string

type T_csvDoc []T_csvLine

// csvDoc converts a T_csvDoc to csv and print it out or if CmdParams.CsvFile is defined it writes the csv to this file
func (c T_csvDoc) csvDoc(typ string) {
	out := [][]string{}
	for _, l := range c {
		out = append(out, l)
	}
	if CmdParams.CsvFile == "" {
		w := csv.NewWriter(os.Stdout)
		if err := w.WriteAll(out); err != nil {
			ErrorLogger.Println("writing csv failed" + err.Error())
		}
	} else {
		file := CmdParams.CsvFile + "-" + typ + ".csv"
		csvfile, err := os.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			ErrorLogger.Println("failed to open file", file, err)
		}
		DebugLogger.Println("write CSV file for", typ, "to", file)
		w := csv.NewWriter(csvfile)
		if err := w.WriteAll(out); err != nil {
			ErrorLogger.Println("writing csv failed" + err.Error())
		}
	}
}

//------------------------------------------

type T_cft_cluster struct {
	Name         T_clName `json:"name"`
	Stage        string   `json:"stage"`
	Console_Url  string   `json:"console_url,omitempty"`
	Consoleurl   string   `json:"console-url,omitempty"`
	Registry     string   `json:"registry"`
	Registry_Url string   `json:"registry_url"`
	Sccs         []string `json:"sccs,omitempty"`
}
type T_cft_clusters []T_cft_cluster

type T_cft_family struct {
	Name              string   `json:"name",omitempty`
	Description       string   `json:"description",omitempty`
	Applications      []string `json:"applications",omitempty`
	Repository_Layout string   `json:"repository_layout",omitempty`
}
type T_cft_families []T_cft_family

type T_cft_environment struct {
	Name         string   `json:"name"`
	Cluster      T_clName `json:"cluster"`
	Pre_Cluster  T_clName `json:"pre_cluster,omitempty"`
	Description  string   `json:"description,omitempty"`
	Family       T_family `json:"family"`
	Pipeline     string   `json:"pipeline,omitempty"`
	Pre_Pipeline string   `json:"pre_pipeline,omitempty"`
}
type T_cft_environments []T_cft_environment

type T_cft_roles struct {
	Admin []string `json:"admin,omitempty"`
	View  []string `json:"view,omitempty"`
}

type T_cft_labels struct {
	It_Services_Id           *string `json:"it_services_id,omitempty"`
	Router                   *string `json:"router,omitempty"`
	Allow_Friend_Environment *string `json:"allow_friend_environment,omitempty"`
	Access_Uniserv           *bool   `json:"access_uniserv,omitempty"`
	Argoproj                 *string `json:"argoproj,omitempty"`
}

type T_cft_namespace struct {
	Name           T_nsName      `json:"name"`
	Environment    string        `json:"environment,omitempty"`
	Applications   []string      `json:"applications,omitempty"`
	Roles          T_cft_roles   `json:"roles,omitempty"`
	Sa_scc_nonroot string        `json:"sa_scc_nonroot,omitempty"`
	Feature_kafka  *bool         `json:"feature_kafka,omitempty"`
	Labels         *T_cft_labels `json:"labels,omitempty"`
}
type T_cft_namespaces []T_cft_namespace

type T_cft_pipeline struct {
	Name        string   `json:"name"`
	Deployer_Ns T_nsName `json:"deployer_ns,omitempty"`
	Deployer_Sa string   `json:"deployer_sa,omitempty"`
	Image_Ns    T_nsName `json:"image_ns,omitempty"`
}
type T_cft_pipelines []T_cft_pipeline
