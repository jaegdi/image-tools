package ocrequest

// The token is from sa image-pruner of namespace cluster-tasks
var Clusters = T_ClusterConfig{
	Config: map[T_clName]T_Cluster{
		"dev-scp0": {
			Name:          "dev-scp0",
			Url:           "https://api.dev-scp0.sf-rz.de:6443",
			Token:         "eyJhbGciOiJSUzI1NiIsImtpZCI6IkJDcXdZNV95YjFFampiWFJGTW1LLVJJRnpZc29WLV9PVFY4eHFYMFlMUDAifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1oamY5NyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIyOTczNzY0OC1iNDVmLTQ0NDMtYTU3Zi1iYzk1NmU3YjdiM2EiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.y6C30QlWqOOX6UMwnDnBWc8MnH2DKEUa733nBfbGaN9YjLJn75eU_Lj2iq_OQY5Q-JTw_fzd4fR6C-_F8Yi8EKJJTzXMUwDFkJMY6chCMVjzw7ccgLzKebE7tWRsy5PocD0KdweRSBU38Ehr9xapvy1c4q_BD0zkgLqb8kxEtgBq0bwEpg3-RkLXfO8Jf8uiLZsBm5WZFBAx_PW1Rt2mvYuHpbb9x7jYS5__MaFrj6QVqKznONB3bQ9u5RYwE13lMOCZrd9zh-N9oa0-6VtHucngC-r88K5Ccy2yE0C_GNOAnGNQ-GLvJRcZt9GGNcJ8zu3OgtmBfgRzqxNbSnmL6OGOx3EGVv0ZK1LsNFF9auur0cqwyTK8N75YwB8r7mpiCSrYtUBrpnzI0k2zCFihV6YP3YiIsSjx2zkxVBBaA6Q3PcIx-miAmKghyou_aSBGFUEBHNKijVrpMXYBRrTd-Ou_sEu2uwUpYraSu49XnGy3gBgpHvuoOElzQan9P9xrfLnsrgvC1H9s90Dwp9dTRJ_2fICSufo_5F-6jhGtZNwUsjX-WclY6KWIQ3HkmCF5HEjGdxJ0fg1Vf96u4DXHDbX2BB8IK8dwyqJ-l9fBEKZ_Z4oXhO5oeydihheguDY_QiWQ1_ZN7t9cnekdLIAbT0BHlpF-Z9kdUaFlzltATpY",
			ConfigToolUrl: "https://scpconfig-service-master.apps.dev-scp0.sf-rz.de"},
		"cid-scp0": {
			Name:          "cid-scp0",
			Url:           "https://api.cid-scp0.sf-rz.de:6443",
			Token:         "eyJhbGciOiJSUzI1NiIsImtpZCI6InFGN0U5SGRZRkNMOHN3ME5qMjVsOTI1VVpJOUlHc0tDTXR0dl9JV2JSVEUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1ncTc1NiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIzMjY3ODkxYi1mNjAxLTRiYzctYTZkMS00MGY5Y2U3ODhjZTEiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.L_Q1kuyJq9dsbli-Qcw73pcoaC7RFduMnO-nzCRsBkIWrgyeTDbxr4yhmC2MLhIylg63BDKhmFDI2X6YK-PdPEURfYOZp5QU_JJDbm2wXXiOknNCmVEM2uTTNgatEvCDtjW7v88uLeh2mTtCMjcSe0tMoBGXvvNn3So2WX6oK9oRAycKIID2rPxrm87kASCQYRjYdY1NbtVDh7idqxgNITo3812zhI_ETR8EHouIBLx3RctYAxc9RTGfGAtq8k1dDo9xdfH0esC0IREJi4bsd2Tq_ZeQhFpzWH9VGSi588lwmuVkUFzoxAkkQwPXPsvpjpiZ0HSfPOdjdq_up1Nw1j7N8ThQv9kxvMaavFsBIywF0PpDWp_H_QtcJsGY43P6I_Hs8O7DzElFOgthbewn3SwgWbSOVqkQD3gE3gz1wfBNhamVasEKvbEc4_tfXXfFFAHVr7_ocXAE7RiN_GW3JswpchV3i5AwPz5JrSJds7rSHgOdr-ApfuerLafDax0REKVzrPheLCdjRj7cV8awLfOY6DpGcWmUbseiabfS0GW8MELsve97yJt3mSN2Y_s2rm4n3VRHPqSQYdI5eoAhQpKqeLr2QMCYEZrIviMpsJY5h_888cHTirucLihVlYhnhE-LN1FdCXQ9pLP-HMrvjFA7zLFcpJytq2FvU0GCabw",
			ConfigToolUrl: "https://scpconfig-service-master.apps.cid-scp0.sf-rz.de"},
		"ppr-scp0": {
			Name:          "ppr-scp0",
			Url:           "https://api.ppr-scp0.sf-rz.de:6443",
			Token:         "eyJhbGciOiJSUzI1NiIsImtpZCI6IjRnSWdKTnRzcVByVWdFNTJjLXVkSEZVdjI5M0JPWkZVVlZtdTRsc1lzdlkifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1yYm1rdiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI3NThjNmNlMi1iNTdmLTRjZjgtOTY3Mi1kOTg5MGU2MzdkYTciLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.aQR0D9t6XZT6B1dNRgPfOxhKc3wZFXvfjoIncw5EyFU23dzQYX8jXRA8c2rjAuIlf0HtiKs8VdA1l4Z6EB3GzO79EHTW4Ru87xUz_UKOdpwy73blYEJgxHnL87mH-2FKL4ecBLPQXaOpHUJ18r2kzSDavtg0sbyjBWpaHWy-D7cAOywJKESgZFYkwje3nUGYRzWlohuqw3_B-v-H0qkWpAjIhisO37SxZRxq4WKcBtpFFkzW9O3YQ5aJhUjjrZzdY7M7aDyH4LhKBmhHInybjbx7-oQ-VSOWGM_Fdh7WzCcxIKdJfufT6oH7KcaMpz30z7p8LIU79bcr3awmeRGMwIT61B8JlQ_NwoVg5DQoyr4sUEiVTg5rEzXXHfdbl-lmJGDFnhYONfA2RKxAYXMsk3BOH-0nRWWGi2v8fauvMdTH6-xFcGlj8SH5H7jmkF1QkDaf50e73hXlsfWaZSSf4ZfJH_rAAzRHmdQ8E4cHW_9zPVe-jX3pDrNESJat09M_BpP8BNR6vZT7ScJKkq4KEFbniQQIURKxjoIWkTCxEkCD3co-79KvnxRSofrufZFCDSZqhuo8uyDmGEDiCDB-Ynr6dqBwmpeZcFWNUEwdUFNZWS34SVvSosCGkmf1iuQy5dtHBq-pKvkMwZEhz7fLFj5qn_MMn2fQn8HSdWXORoQ",
			ConfigToolUrl: "https://scpconfig-service-master.apps.ppr-scp0.sf-rz.de"},
		"vpt-scp0": {
			Name:          "vpt-scp0",
			Url:           "https://api.vpt-scp0.sf-rz.de:6443",
			Token:         "eyJhbGciOiJSUzI1NiIsImtpZCI6Ikl5eG9kREQtQlVIQTBMcjRDYm5pWnZ1bXlXTWExdjkxaEVyVEx1N20wMkUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1sOWNxYiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIwMTA2MTg0Ny0yNTA0LTRkNzMtYmE1ZC02NGMzZDczMjFhYjciLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.eZNfnreIVd4f_6sF6yPLh1EIoGQzCLOjcLgYEjZHOvzmX1pi6DizzgFKMF_ca_gHxZ3pit8-v8nYnJwcwC-pcU5sDOo17e0v3MCshP7_A2LtOLL8vt_OkODL6USMx1zaST_6N7gDC0M3Xkc5FlyN8K9Vsrk-76xycbw0XqK7LTNqvcVHsijJCgQJaAwq8asrbnBMJRN-c9CI4Eot7bAYX1oKEB2H5LdzgkDU3pfzpHcjNs3R7YYARd5w1O-DgqkCillRxplVFx2_K92h1xX2bQdGwDyFDelaJisLvZYiaQeLXKoPalE1sJ4fnUl_eV8tXyQGgax9QDvm2OOVd46w7AbgYRiu5Um9U_ZQgTb3bSmXy1Gpusci0oll_EMQ0lR0qyU4DtAE78MxvSa3MdjrtnR-ZJDOs93NX-Og8DJhthml8edpKgP-Aiq2j1s-7Imttx_zeoMW6vTgdl1lL-LN7ypj5aTAVWUtuOF9fNE6VSf8LA0uhfi0EVk3dWoJ5dCSbkPXybn6t9nYbLu_Bpgsm-C21_sljh_gOH6XnyJcqgbkuswW-qDOHA3kXplufghQutnrW7jQO3zVyw-DQoJmApO6YsnfwewQB8UU4P7H_z87ubQaa_WxLO0pDbPuwudDolPpR2XZglKMvp_lK6K3tRyZl6ytUfc_LpLiQBReptw",
			ConfigToolUrl: "https://scpconfig-service-master.apps.vpt-scp0.sf-rz.de"},
		"pro-scp0": {
			Name:          "pro-scp0",
			Url:           "https://api.pro-scp0.sf-rz.de:6443",
			Token:         "eyJhbGciOiJSUzI1NiIsImtpZCI6IndGRFo3REpZNlh0bmd2SWtyOFVnNXQxclp5dElPS2RLbWNTWDlGSDhwS1kifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1nOXZtdiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI2MjZlMzU0YS0xYTRmLTQ0ZDYtOTdmOS1jNjAxMjA4ZTMxNzYiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.wuyxvfzeDVIfQkI012uh_Bi3MbWaiMy6-jFXRqbd2rhvudp84KbrPyGYIFNgzrtnvCe815vtTw9dWk7lr5M7Y6aUp1Vb5CsLRty7NbXKa2J8LPHHmpwwEdMWcbBVLEHCUVmZee4nZ0WsTKjxRsvFrJ99a0-ikLIVPkUdZG9yon5dgPpy2j0j5PD7zk0hQAF1wgs7LjsR11uOJkZso0XcoXyDdiKPNMDI6ySBg0q1pJ6TaLkihWCw3OPB-aDvOFJBL_6FWCJV5MBV10LIZYXV5TsNRvQKoID6qHftZIQvEhKYtQu7Q3d1kNIq71bYkF09vRhIctP1gBkmbtGq9c9ZBsGHTWWgwIwNI192ztjPHpKMidbfpDmCuKn3-iIF4dDxFS8ApUfjDB79ib0TsoARp3lFrGQV9Ea-V7IRqgndcoKbYZ9Ebu9m0gUUV8W95d8Q0VX1vNyKzhYVSrsp-3ot6OIbHJDlJuljaOv-jGGS3063MT_lYkypEySYUe3dQDrur2--BnTF2hA7SOP57y-AuQdIOMT9LvxEUxqL_JRpu-lx5CY7NH6VCXI2kkbbqeJxMgNw5f1AMlcVvxVW7WX-VYTv0RH6iCw02OKMydVPZF-aDloXQK6_yHGj3tCd73i3cwAezgEEIxHE4-XEQq_gqZNxo6FngqOdV4MkvXwH_y0",
			ConfigToolUrl: "https://scpconfig-service-master.apps.pro-scp0.sf-rz.de"},
	},
}

var FamilyNamespaces T_famNsList
var AppNamespaces T_appNsList
var FamilyNamespacesStat = T_famNsList{
	"pkp": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-scp0": {"pkp-images"},
			"ppr-scp0": {"pkp-images"},
			"vpt-scp0": {"pkp-images"},
			"pro-scp0": {"pkp-images"},
		},
		Stages:      []T_clName{"cid-scp0", "ppr-scp0", "pro-scp0"},
		Buildstages: []T_clName{"cid-scp0"},
		Teststages:  []T_clName{"ppr-scp0", "vpt-scp0"},
		Prodstages:  []T_clName{"pro-scp0"},
	},
	"ebs": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-scp0": {"ebs-images"},
			"ppr-scp0": {"ebs-images"},
			"vpt-scp0": {"ebs-images"},
			"pro-scp0": {"ebs-images"},
		},
		Stages:      []T_clName{"cid-scp0", "ppr-scp0", "pro-scp0"},
		Buildstages: []T_clName{"cid-scp0"},
		Teststages:  []T_clName{"ppr-scp0", "vpt-scp0"},
		Prodstages:  []T_clName{"pro-scp0"},
	},
	"aps": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-scp0": {"aps-images", "openshift"},
			"ppr-scp0": {"aps-images"},
			"vpt-scp0": {"aps-images"},
			"pro-scp0": {"aps-images"},
		},
		Stages:      []T_clName{"cid-scp0", "ppr-scp0", "pro-scp0"},
		Buildstages: []T_clName{"cid-scp0"},
		Teststages:  []T_clName{"ppr-scp0", "vpt-scp0"},
		Prodstages:  []T_clName{"pro-scp0"},
	},
	"vps": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-scp0": {"vps-images"},
			"ppr-scp0": {"vps-images"},
			"vpt-scp0": {"vps-images"},
			"pro-scp0": {"vps-images"},
		},
		Stages:      []T_clName{"cid-scp0", "ppr-scp0", "pro-scp0"},
		Buildstages: []T_clName{"cid-scp0"},
		Teststages:  []T_clName{"ppr-scp0", "vpt-scp0"},
		Prodstages:  []T_clName{"pro-scp0"},
	},
	"dca": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-scp0": {"images-dca"},
			"int-scp0": {"images-dca"},
			"ppr-scp0": {"images-dca"},
			"vpt-scp0": {"images-dca"},
			"pro-scp0": {"images-dca"},
		},
		Stages:      []T_clName{"cid-scp0", "ppr-scp0", "pro-scp0"},
		Buildstages: []T_clName{"cid-scp0"},
		Teststages:  []T_clName{"ppr-scp0", "vpt-scp0"},
		Prodstages:  []T_clName{"pro-scp0"},
	},
	"fpc-scp0": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-scp0": {"b2c-images"},
			"ppr-scp0": {"b2c-images"},
			// "vpt-scp0": {"b2c-images"},
			"pro-scp0": {"b2c-images"},
		},
		Stages:      []T_clName{"cid-scp0", "ppr-scp0", "vpt-scp0"},
		Buildstages: []T_clName{"cid-scp0"},
		Teststages:  []T_clName{"cid-scp0", "ppr-scp0", "vpt-scp0"},
		Prodstages:  []T_clName{"ppr-scp0"},
	},
	"scp": {
		ImageNamespaces: map[T_clName][]T_nsName{
			"cid-scp0": {"scp-images", "scp-baseimages"},
			"ppr-scp0": {"scp-images", "scp-baseimages"},
			"vpt-scp0": {"scp-images", "scp-baseimages"},
			"pro-scp0": {"scp-images", "scp-baseimages"},
		},
		Stages:      []T_clName{"dev-scp0", "cid-scp0", "ppr-scp0", "vpt-scp0", "pro-scp0"},
		Buildstages: []T_clName{"cid-scp0"},
		Teststages:  []T_clName{"dev-scp0", "cid-scp0", "ppr-scp0", "vpt-scp0"},
		Prodstages:  []T_clName{"vpt-scp0"},
	},
}

var OcClient bool

var LogFileName = "image-tools.log"

var bitbucket_token = "NDk4MzkwOTAxNzgzOifliHsGrHWZUpubeoNMZbhM/7qT"
