package ocrequest

type T_completeResults struct {
	AllIstags  T_result
	UsedIstags T_usedIstagsResult
}

// getExistingIstags.go
type T_shaStreams map[string]map[string]T_istag
type T_shaNames map[string]map[string]interface{}

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

// type T_sha map[string]interface{}
type T_sha struct {
	Istags      interface{}
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
	IsTag     string
	Namespace string
	Sha       string
	Is        string
	Cluster   string
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
