package ocrequest

var Clusters T_ClusterConfig
var FamilyNamespaces T_famNsList
var AppNamespaces T_appNsList
var FamilyNamespacesStat = T_famNsList{
	"aps": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-crp0": {"aps-images", "aps-build"},
			"ppr-crp0": {"aps-images"},
			"vpt-crp0": {"aps-images"},
			"pro-crp0": {"aps-images"},
		},
		Stages:      []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0", "pro-crp0"},
		Buildstages: []T_clName{"cid-crp0"},
		Teststages:  []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0"},
		Prodstages:  []T_clName{"pro-crp0"},
	},
	"b2b": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-crp0": {"b2b-images", "b2b-build"},
			"ppr-crp0": {"b2b-images"},
			"vpt-crp0": {"b2b-images"},
			"pro-crp0": {"b2b-images"},
		},
		Stages:      []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0", "pro-crp0"},
		Buildstages: []T_clName{"cid-crp0"},
		Teststages:  []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0"},
		Prodstages:  []T_clName{"pro-crp0"},
	},
	"b2c": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-crp0": {"b2c-images", "b2b-build"},
			"ppr-crp0": {"b2c-images"},
			"vpt-crp0": {"b2c-images"},
			"pro-crp0": {"b2c-images"},
		},
		Stages:      []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0", "pro-crp0"},
		Buildstages: []T_clName{"cid-crp0"},
		Teststages:  []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0"},
		Prodstages:  []T_clName{"pro-crp0"},
	},
	"cbs": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-crp0": {"cbs-images", "cbs-build"},
			"int-crp0": {"cbs-images"},
			"ppr-crp0": {"cbs-images"},
			"vpt-crp0": {"cbs-images"},
			"pro-crp0": {"cbs-images"},
		},
		Stages:      []T_clName{"cid-crp0", "ppr-crp0", "pro-crp0"},
		Buildstages: []T_clName{"cid-crp0"},
		Teststages:  []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0"},
		Prodstages:  []T_clName{"pro-crp0"},
	},
	"dca": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-crp0": {"dca-images", "dca-build"},
			"int-crp0": {"dca-images"},
			"ppr-crp0": {"dca-images"},
			"vpt-crp0": {"dca-images"},
			"pro-crp0": {"dca-images"},
		},
		Stages:      []T_clName{"cid-crp0", "ppr-crp0", "pro-crp0"},
		Buildstages: []T_clName{"cid-crp0"},
		Teststages:  []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0"},
		Prodstages:  []T_clName{"pro-crp0"},
	},
	"ebs": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-crp0": {"ebs-images", "ebs-build"},
			"ppr-crp0": {"ebs-images"},
			"vpt-crp0": {"ebs-images"},
			"pro-crp0": {"ebs-images"},
		},
		Stages:      []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0", "pro-crp0"},
		Buildstages: []T_clName{"cid-crp0"},
		Teststages:  []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0"},
		Prodstages:  []T_clName{"pro-crp0"},
	},
	"ibs": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-crp0": {"ibs-images", "ibs-build"},
			"ppr-crp0": {"ibs-images"},
			"vpt-crp0": {"ibs-images"},
			"pro-crp0": {"ibs-images"},
		},
		Stages:      []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0", "pro-crp0"},
		Buildstages: []T_clName{"cid-crp0"},
		Teststages:  []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0"},
		Prodstages:  []T_clName{"pro-crp0"},
	},
	"pkp": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-crp0": {"pkp-images", "pkp-build"},
			"ppr-crp0": {"pkp-images"},
			"vpt-crp0": {"pkp-images"},
			"pro-crp0": {"pkp-images"},
		},
		Stages:      []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0", "pro-crp0"},
		Buildstages: []T_clName{"cid-crp0"},
		Teststages:  []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0"},
		Prodstages:  []T_clName{"pro-crp0"},
	},
	"crp": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-crp0": {"crp-images", "crp-build", "crp-baseimages"},
			"ppr-crp0": {"crp-images", "crp-baseimages"},
			"vpt-crp0": {"crp-images", "crp-baseimages"},
			"pro-crp0": {"crp-images", "crp-baseimages"},
		},
		Stages:      []T_clName{"dev-crp0", "cid-crp0", "ppr-crp0", "vpt-crp0", "pro-crp0"},
		Buildstages: []T_clName{"cid-crp0"},
		Teststages:  []T_clName{"dev-crp0", "cid-crp0", "ppr-crp0", "vpt-crp0"},
		Prodstages:  []T_clName{"pro-crp0"},
	},
	"vps": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-crp0": {"vps-images", "vps-build"},
			"ppr-crp0": {"vps-images"},
			"vpt-crp0": {"vps-images"},
			"pro-crp0": {"vps-images"},
		},
		Stages:      []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0", "pro-crp0"},
		Buildstages: []T_clName{"cid-crp0"},
		Teststages:  []T_clName{"cid-crp0", "ppr-crp0", "vpt-crp0"},
		Prodstages:  []T_clName{"pro-crp0"},
	},
}

var OcClient bool

var LogFileName = "image-tools.log"
