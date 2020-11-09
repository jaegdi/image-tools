package ocrequest

import (
	// "log"
)

type T_completeResults struct {
	AllIstags  T_result
	UsedIstags T_usedIstagsResult
}

// getExistingIstags.go
type T_shaStreams map[string]map[string]T_istag
type T_shaNames map[string]map[string]bool

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

func (b T_istagBuildLabels) List() []string {
	l := []string{}
	l = append(l, b.CommitAuthor)
	l = append(l, b.CommitDate)
	l = append(l, b.CommitId)
	l = append(l, b.CommitRef)
	l = append(l, b.CommitVersion)
	l = append(l, b.IsProdImage)
	l = append(l, b.Name)
	l = append(l, b.Namespace)
	return l
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

func (c T_istag) List() []string {
	line := []string{}
	line = append(line, c.Imagestream)
	line = append(line, c.Tagname)
	line = append(line, c.Namespace)
	line = append(line, c.Link)
	line = append(line, c.Date)
	line = append(line, c.AgeInDays)
	line = append(line, c.Sha)
	line = append(line, c.Build.List()...)
	return line
}

// type T_sha map[string]interface{}
type T_Istags_List map[string]bool
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

type T_Clusters struct {
	Cid T_Cluster
	Int T_Cluster
	Ppr T_Cluster
	Vpt T_Cluster
	Pro T_Cluster
}

type T_csvLine []string

type T_csvDoc []T_csvLine

func (c T_csvDoc) csvDoc() [][]string {
	out := [][]string{}
	for _, l := range c {
		out = append(out, l)
	}
	return out
}
