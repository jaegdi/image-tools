package ocrequest

// getExistingIstags.go
type T_shaStreams map[string]map[string]T_istag
type T_shaNames map[string]map[string]interface{}

// type T_istag map[string]interface{}
type T_istag struct {
	Imagestream string
	Tagname     string
	Namespace   string
	Link        string
	Date        string
	AgeInDays   string
	Sha         string
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

type T_isShaTagnames []string
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
}

// usedIstagsResult[Is][Tag]T_usedIstag
type T_usedIstagsResult map[string]map[string]T_usedIstag

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
	Cluster string
	Token   string
	Family  string
	Output  T_flagOut
	Filter  T_flagFilt
}

//------------------------------------------
