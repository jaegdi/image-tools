package ocrequest

import (
	// "log"
	"encoding/json"
	"reflect"
)

type T_completeResults struct {
	AllIstags  T_result
	UsedIstags T_usedIstagsResult
}

// getExistingIstags.go
type T_shaStreams map[string]map[string]T_istag

func (a T_shaStreams) Add(is string, sha string, istag T_istag) {
	if a == nil {
		a = T_shaStreams{}
	}
	if a[is] == nil {
		a[is] = T_resIstag{}
	}
	if a[is][sha] == (T_istag{}) {
		a[is][sha] = T_istag{}
	}
	// for k, v := range istag {
	a[is][sha] = istag
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
	CommitAuthor  string `json:"io.openshift.build.commit.author,omitempty"`
	CommitDate    string `json:"io.openshift.build.commit.date,omitempty"`
	CommitId      string `json:"io.openshift.build.commit.id,omitempty"`
	CommitRef     string `json:"io.openshift.build.commit.ref,omitempty"`
	CommitVersion string `json:"io.openshift.build.commit.version,omitempty"`
	IsProdImage   string `json:"io.openshift.build.isProdImage,omitempty"`
	Name          string `json:"io.openshift.build.name,omitempty"`
	Namespace     string `json:"io.openshift.build.namespace,omitempty"`
}

func (b T_istagBuildLabels) Values() []string {
	l := []string{}
	v := reflect.ValueOf(b)
	for i := 0; i < v.NumField(); i++ {
		l = append(l, v.Field(i).String())
	}
	return l
}
func (b T_istagBuildLabels) Names() []string {
	l := []string{}
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
	buildLabelsJSON := []byte(GetJsonFromMap(buildLabelsMap))
	if err := json.Unmarshal(buildLabelsJSON, &buildLabels); err != nil {
		ErrorLogger.Println("Unmarshal unescaped String", err)
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
	Sha         string
	Build       T_istagBuildLabels
}

func (c T_istag) Values() []string {
	l := []string{}
	v := reflect.ValueOf(c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		switch typeOfS.Field(i).Name {
		case "Build":
			l = append(l, c.Build.Values()...)
		default:
			l = append(l, v.Field(i).String())
		}
	}
	return l
}

func (c T_istag) Names() []string {
	l := []string{}
	v := reflect.ValueOf(c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		switch typeOfS.Field(i).Name {
		case "Build":
			l = append(l, c.Build.Names()...)
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
	AnzNames    int
	AnzShas     int
	AnzIstreams int
}

type T_result struct {
	Is     T_resIs
	Istag  T_resIstag
	Sha    T_resSha
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
	Sha             string
	Cluster         string
}

func (c T_usedIstag) Values() []string {
	l := []string{}
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumField(); i++ {
		l = append(l, v.Field(i).String())
	}
	return l
}

func (c T_usedIstag) Names() []string {
	l := []string{}
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

// userinterface.go
type T_famNs map[string][]string

type T_flagOut struct {
	Is    bool
	Istag bool
	Sha   bool
	Used  bool
	All   bool
}
type T_flagFilt struct {
	Isname    string
	Istagname string
	Shaname   string
	Namespace string
}
type T_flags struct {
	Cluster  string
	Token    string
	Family   string
	OcClient bool
	Json     bool
	Yaml     bool
	Csv      bool
	Output   T_flagOut
	Filter   T_flagFilt
}

//------------------------------------------

type T_Cluster struct {
	Name  string
	Url   string
	Token string
}

// type T_Clusters struct {
// 	Cid T_Cluster
// 	Int T_Cluster
// 	Ppr T_Cluster
// 	Vpt T_Cluster
// 	Pro T_Cluster
// }

type T_csvLine []string

type T_csvDoc []T_csvLine

func (c T_csvDoc) csvDoc() [][]string {
	out := [][]string{}
	for _, l := range c {
		out = append(out, l)
	}
	return out
}
