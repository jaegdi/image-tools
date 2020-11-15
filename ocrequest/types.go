package ocrequest

import (
	// "log"
	// "encoding/json"
	// "github.com/jedib0t/go-pretty/v6/table"
	// "log"
	// "os"
	"encoding/csv"
	"os"
	"reflect"
)

type T_completeResults struct {
	AllIstags  T_ResultExistingIstagsOverAllClusters
	UsedIstags T_usedIstagsResult
}

//                                 family
type T_completeResultsFamilies map[string]T_completeResults

// getExistingIstags.go
type T_shaStreams map[string]map[string]T_istag

func (a T_shaStreams) Add(is string, image string, istag T_istag) {
	if a == nil {
		a = T_shaStreams{}
	}
	if a[is] == nil {
		a[is] = T_resIstag{}
	}
	if a[is][image] == (T_istag{}) {
		a[is][image] = T_istag{}
	}
	// for k, v := range istag {
	a[is][image] = istag
}

type T_Istags_List map[string]bool
type T_shaNames map[string]T_Istags_List

func (a T_shaNames) Add(key string, b string) {
	if a == nil {
		a = T_shaNames{}
	}
	if a[key] == nil {
		a[key] = T_Istags_List{}
	}
	a[key][b] = true
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
	// 	LogError("Unmarshal unescaped String", err)
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
	Imagestream string
	Tagname     string
	Namespace   string
	Link        string
	Date        string
	AgeInDays   string
	Image       string
	Build       T_istagBuildLabels
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
	Imagestream string
	Namespace   string
	Link        string
	Date        string
	AgeInDays   string
}

type T_isShaTagnames map[string]interface{}
type T_is map[string]T_isShaTagnames

type T_resIstag map[string]T_istag
type T_resIs map[string]T_is
type T_resSha map[string]map[string]T_sha

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
	UsedInNamespace string
	Image           string
	Cluster         string
}

func (c T_usedIstag) Values() interface{} {
	l := []interface{}{}
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumField(); i++ {
		l = append(l, v.Field(i).String())
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

// usedIstagsResult[Is][Tag]T_usedIstag
type T_usedIstagsResult map[string]map[string][]T_usedIstag

// IsNamesForFamily[family][is]true
type T_IsNamesForFamily map[string]map[string]bool

//------------------------------------------

//               family  namespaces
type T_famNs map[string][]string

type T_flagOut struct {
	Is    bool
	Istag bool
	Image bool
	Used  bool
	All   bool
}
type T_flagFilt struct {
	Isname    string
	Istagname string
	Tagname   string
	Imagename string
	Namespace string
}

type T_flagOpts struct {
	OcClient bool
	NoProxy  bool
	Profiler bool
}
type T_flags struct {
	Cluster  string
	Token    string
	Family   string
	Json     bool
	Yaml     bool
	Csv      bool
	CsvFile  string
	Html     bool
	Table    bool
	TabGroup bool
	Output   T_flagOut
	Filter   T_flagFilt
	Options  T_flagOpts
}

//------------------------------------------

type T_Cluster struct {
	Url   string `json:"Url,omitempty"`
	Name  string `json:"Name,omitempty"`
	Token string `json:"Token,omitempty"`
}

type T_ClusterConfig struct {
	Stages     []string
	Config     map[string]T_Cluster `json:"Config.[],omitempty"`
	Buildstage string
	Teststages []string
	Prodstage  string
}

type T_csvLine []string

type T_csvDoc []T_csvLine

func (c T_csvDoc) csvDoc(typ string) {
	out := [][]string{}
	for _, l := range c {
		out = append(out, l)
	}
	if CmdParams.CsvFile == "" {
		w := csv.NewWriter(os.Stdout)
		if err := w.WriteAll(out); err != nil {
			LogError("writing csv failed" + err.Error())
		}
	} else {
		file := CmdParams.CsvFile + "-" + typ + ".csv"
		csvfile, err := os.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			LogError("failed to open file", file, err)
		}
		LogMsg("write CSV file for", typ, "to", file)
		w := csv.NewWriter(csvfile)
		if err := w.WriteAll(out); err != nil {
			LogError("writing csv failed" + err.Error())
		}
	}
}
