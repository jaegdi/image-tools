package ocrequest

var FamilyNamespaces = T_famNs{
	"pkp": []string{"ms-jenkins", "openshift", "images-pkp"},
	"ssp": []string{"ssp-jenkins", "images-ssp"},
	"aps": []string{"aps-jenkins", "images-aps"},
	"fpc": []string{"fpc-jenkins", "images-fpc"},
}

var Clusters = map[string]interface{}{
	"cid": T_Cluster{
		Name: "cid-apc0",
		Url:  "https://console.cid-apc0.sf-rz.de:8443"},
	"int": T_Cluster{
		Name: "int-apc0",
		Url:  "https://console.int-apc0.sf-rz.de:8443"},
	"ppr": T_Cluster{
		Name: "ppr-apc0",
		Url:  "https://console.ppr-apc0.sf-rz.de:8443"},
	"vpt": T_Cluster{
		Name: "pro-apc0",
		Url:  "https://console.pro-apc0.sf-rz.de:8443"},
	"pro": T_Cluster{
		Name: "pro-apc0",
		Url:  "https://console.pro-apc0.sf-rz.de:8443"},
	"stages":     []string{"cid", "int", "ppr", "vpt", "pro"},
	"buildstage": "cid",
	"teststages": []string{"int", "ppr", "vpt"},
	"prodstage":  "pro",
}

var OcClient bool
