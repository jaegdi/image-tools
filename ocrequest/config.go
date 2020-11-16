package ocrequest

var FamilyNamespaces = T_famNs{
	"pkp": map[string][]string{
		"cid": {"ms-jenkins", "openshift", "images-pkp"},
		"int": {"images-pkp"},
		"ppr": {"images-pkp"},
		"vpt": {"images-pkp"},
		"pro": {"images-pkp"},
	},
	"ssp": map[string][]string{
		"cid": {"ssp-jenkins", "images-ssp"},
		"int": {"images-ssp"},
		"ppr": {"images-ssp"},
		"vpt": {"images-ssp"},
		"pro": {"images-ssp"},
	},
	"aps": map[string][]string{
		"cid": {"aps-jenkins", "images-aps"},
		"int": {"images-aps"},
		"ppr": {"images-aps"},
		"vpt": {"images-aps"},
		"pro": {"images-aps"},
	},
	"fpc": map[string][]string{
		"cid": {"fpc-basis-1-1-20", "fpc-basis-1-1-21"},
		"int": {"fpc-fa1", "fpc-fa2", "fpc-int1", "fpc-int2", "fpc-loadtest"},
		"ppr": {"fpc-ppr"},
		"vpt": {"vptest-fpc"},
		"pro": {"fpc"},
	},
}

var Clusters = T_ClusterConfig{
	Config: map[string]T_Cluster{
		"cid": T_Cluster{
			Name:  "cid-apc0",
			Url:   "https://console.cid-apc0.sf-rz.de:8443",
			Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi16cTZ3ZyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJlMTJiNjI0Zi1kOTAwLTExZTgtODI5NS0wMDUwNTY5MDM3ZmIiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.soU7K1s7wMSW63VLMrxMUMo2jZIOp8pbf11biz7RFmylpHdvPx72yETs2cBevncxfhcmJXOnWtiSBtaxpHsQwbsMLwy5vNR1Aoajy7uR-_TktWFhfQu6ak-fuoHOPZml5dL8WZJCR_wvkov40k1kNeBSRjH0aXd_YUECk8jOmn9kxHWmHjcTuhoF8_mH9UCU7fGWPMa0ahUrllZlttqf1ZQcmk4oLi4X2JGIHN6pG9hQV0nOqutmkxdbshiH4od-aKljh5sX1pcql5NK9FliwFYKYPfRB8QRh5SBftu_VxGbqfRdrjdxbKgZzoZk002o5S-PYegagxZZAqzu9wti2Q"},
		"int": T_Cluster{
			Name:  "int-apc0",
			Url:   "https://console.int-apc0.sf-rz.de:8443",
			Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1nNjdxcCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJhYmY2ODM4Ni0wNDQ1LTExZTktOGFmYi0wMDUwNTY5MGQ4MjEiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.Tb5e_b3QiBUHmKzBOS0h6r_NeIw_GMERoGcEuyEjSOk-LZPU6q8atsAOEYg9QZdlo0rghxZhPGl8plxk-EVWgkWhQZpOVpp1HOZj8w097E9vRfRRkpt1e0I6ozq_PBx0K8758lxWa-LDQGV1vf8yTGpEUupps08ElSgwApgoVm6z1ZkmiBl_B_jb-MAH5PYnbiAQWu57AE8WXr_rOU8EunYK0fLJcBljMU7bhyxKW1WChIe7COYATCggQuRsu22ySF7dlTeGq3NLAFe2BpmeJavLle8FUQ1pjOem0uKrPEGat2H03KRKyuSP-7t34Jue5kdd3Qwf4rOraiVJ6-3T3w"},
		"ppr": T_Cluster{
			Name:  "ppr-apc0",
			Url:   "https://console.ppr-apc0.sf-rz.de:8443",
			Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1qMnJ4eCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJjMDdkM2ZkNi0wNDRlLTExZTktYmJiYi0wMDUwNTZhMTkwNjQiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.OHHuZp-90WOBCaXjwTCCbvxW87lSsaBaH7q59JSDLM5oSWmca_aUbA3J_2nB6ruOOrvEbbQayb-UMcGGXhiVUv27Cg9TVb9Gl3AZym1gk6DUhma2wB3IwI98RKhX3XMr2J7pfte00p1Za8jBcVx_vdYaxebDQImpyUOFaEC3MT4ueirRXLTNSOQLnMARTVpr5r3s0Z0W3g8MGre7MDMNOO4BcGyEC1X8O8XzObuekzSErjJ6zugQMjIYw7_YkPULYNTJLrGsd136TFfFNGEtXow2RuLynIU8_1Zatvs64zNRirgvplGwxg4f15cHp7yUFJ_q3krkUegRji5tZdOBJw"},
		"vpt": T_Cluster{
			Name:  "pro-apc0",
			Url:   "https://console.pro-apc0.sf-rz.de:8443",
			Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1qMnJ4eCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJjMDdkM2ZkNi0wNDRlLTExZTktYmJiYi0wMDUwNTZhMTkwNjQiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.OHHuZp-90WOBCaXjwTCCbvxW87lSsaBaH7q59JSDLM5oSWmca_aUbA3J_2nB6ruOOrvEbbQayb-UMcGGXhiVUv27Cg9TVb9Gl3AZym1gk6DUhma2wB3IwI98RKhX3XMr2J7pfte00p1Za8jBcVx_vdYaxebDQImpyUOFaEC3MT4ueirRXLTNSOQLnMARTVpr5r3s0Z0W3g8MGre7MDMNOO4BcGyEC1X8O8XzObuekzSErjJ6zugQMjIYw7_YkPULYNTJLrGsd136TFfFNGEtXow2RuLynIU8_1Zatvs64zNRirgvplGwxg4f15cHp7yUFJ_q3krkUegRji5tZdOBJw"},
		"pro": T_Cluster{
			Name:  "pro-apc0",
			Url:   "https://console.pro-apc0.sf-rz.de:8443",
			Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi16N3NrdyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI3ZWY0YjE4OC0zNDM2LTExZTktOTIzZC0wMDUwNTY4MjA2Y2UiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.PX5Ghkmc4CgsvQgXuuIBBS95W0tE0YsRuCbymp4eVUBo4l1woEkKcds7yghx9PwcLHNMWWVki1DUUw-TbsWJY5MrnNwWCGwxTixeh3HSv3RBNAOkGUX4JPt6Bs0jEL_sELccLq9df6pZTQyDjidoSFnT4rid7R-7NyoRjlm0PhZRXEKxlGDt2eZNT0OJGl7btWqh1I_R7MTeuvya-KGSaMtNOlgpvLeyfH2UAvSY6INCu9ca88CBNWScpvZ_X6uxF75bLVGb1UJurPMbrG0isWge6kVnCaCpX_bSmlr8nyo7N6L1feGEgZZpwquTugbwfte8fLvbBgWCB9d8pt-ZhA"},
	},
	Stages:     []string{"cid", "int", "ppr", "pro"},
	Buildstage: "cid",
	Teststages: []string{"int", "ppr", "vpt"},
	Prodstage:  "pro",
}

var OcClient bool
