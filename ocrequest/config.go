package ocrequest

//  The token is from sa image-pruner of namespace cluster-tasks
var Clusters = T_ClusterConfig{
	Config: map[T_clName]T_Cluster{
		"cid-scp0": {
			Name:  "cid-scp0",
			Url:   "https://api.cid-scp0.sf-rz.de:6443",
			Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6InFGN0U5SGRZRkNMOHN3ME5qMjVsOTI1VVpJOUlHc0tDTXR0dl9JV2JSVEUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi16Z3RsZCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIzMjY3ODkxYi1mNjAxLTRiYzctYTZkMS00MGY5Y2U3ODhjZTEiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.cSAa0k4EaY-yngdtcNcUMgKGx6B5S3cVDqx18D37CMN2UTXRl5M2KmSPqdAq8SLyLCgjTCfybSQi8pX8ZhlD6Mflc7rqd38QI1JIC6JHGSWBI5PqMN7rqSIMRKT3JqKSDzH2ErN0toLDGv8OBndq9X1u4woK8M1cKIKqOtAnezE0jJF5AAhJz7PLTQDC96UiHa4bw-bNII6iG9QOQoQ8g8-QWAB6bJH4obCvmaODzbmSEsW4fE_kNvVusgtem5pR2ms7JWH_0cT4ZOJ_U3wrp7yUlrDJS_QhivlusTedDORepHi3X6wxJG-dZQlSxq70Yx6TavuG2dTlF5VzCfeqTmVPP3xnMiFtkiCkIygXnvi0AUG6nRJPqtmC9JQSNAcLuXXMGk1GSakJPhtGOPoe9M2IaG3O77dBvb1Qy0TgiTabt7zAGZaKQCjYpFbstnM2_wngOzOc0M0orZjcHWmJO_ZamZdYyMjQJOsx0JDeE7OFuD_ZETzKWDpSYqr8ptFJWobMea38xbY67B8CVzNXc00_BILLi0wyjHhYI0f_PJpM1VY4Bz8l3xjPpgWFMf6o5jxLgw2KgNhAA73xuOdM5jdCma1Wv-r5PPMlcv-n_wiV7ZZ3ZURamzIUPfxUzIbAxKyRNJWAXUvcZBKIZ8R0BJ5ieIuFLm5a9BGdcr33G6M"},
		"ppr-scp0": {
			Name:  "ppr-scp0",
			Url:   "https://api.ppr-scp0.sf-rz.de:6443",
			Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6IlNlNTQxeUdFVTcwNHFhVDd3UWEtOWdoSU02X1AzWElMZkhOWmU5Z2dRRjAifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1ydG5obiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI2N2FjZGI4NC1lMmZmLTQzMjgtOGFlMC1hNjdjZTgzOWYxYTkiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.SzlBaJwPqgvbfSPh80La8d2wcJatPu-Y0Q0wUjhYowS_fcGMtmwBVgvLM9zjZ89CWiS-Z--h9srWANyOb_TK5YwBokOM0Jr9hLgnPWgdH6AaD-mM0iQH05khV-w7zEwb_t226zTGln_YwL8-Dkr8jsO1uWNmx6kb2aW7ZhGuGXJz6ONuGFmMV8Dk78O3jTOVXbpE_rFMcvA29WWVeWk2Z0croNxI02h0FsRVZ8Um13WJs5R1kzf5iwlX2aSyRGqOOaSitkHIWowbG2HIlyWt-Nf_Q6Ws4X6Ua0r94xqqiQ6dyGmQVj6hb6-E5jxl4w-jyoxPGLCYRz8k3UyGyGqZ4dPyiWfRifEjPbo1gsDpfmMF8nYeRAFFHgL42oIsG0H3JQ3Gsdz2pRaMRwHtT_6pYlMyfVauDXnpD0Wr7_nPv1icZy2sttOBzygXwPwz6iwidglAM6CQEcuLNnmRev0sCM8KWGSFKGqff_qNdoLjSUfW2CqFfA_Wr4_EUSqGXyOuPQc0DCzjgMgKW8Eeiqu32-C_a4fmrnaYd_8ZaPwyp7DDshVNB3sBasAmr7RmcIb2Fg0hj0yjOi2ucoXZuUsEUm93DxQf_4bsMkwK4EGKWvKamqwye-pGu-dVGfWBeqcHZNFg7tpdjyGlpUIY0_lJJtNwiR3hFCNQI5oSURjzsd8"},
		"vpt-scp0": {
			Name:  "vpt-scp0",
			Url:   "https://api.vpt-scp0.sf-rz.de:6443",
			Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6Ikl5eG9kREQtQlVIQTBMcjRDYm5pWnZ1bXlXTWExdjkxaEVyVEx1N20wMkUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1sOWNxYiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIwMTA2MTg0Ny0yNTA0LTRkNzMtYmE1ZC02NGMzZDczMjFhYjciLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.eZNfnreIVd4f_6sF6yPLh1EIoGQzCLOjcLgYEjZHOvzmX1pi6DizzgFKMF_ca_gHxZ3pit8-v8nYnJwcwC-pcU5sDOo17e0v3MCshP7_A2LtOLL8vt_OkODL6USMx1zaST_6N7gDC0M3Xkc5FlyN8K9Vsrk-76xycbw0XqK7LTNqvcVHsijJCgQJaAwq8asrbnBMJRN-c9CI4Eot7bAYX1oKEB2H5LdzgkDU3pfzpHcjNs3R7YYARd5w1O-DgqkCillRxplVFx2_K92h1xX2bQdGwDyFDelaJisLvZYiaQeLXKoPalE1sJ4fnUl_eV8tXyQGgax9QDvm2OOVd46w7AbgYRiu5Um9U_ZQgTb3bSmXy1Gpusci0oll_EMQ0lR0qyU4DtAE78MxvSa3MdjrtnR-ZJDOs93NX-Og8DJhthml8edpKgP-Aiq2j1s-7Imttx_zeoMW6vTgdl1lL-LN7ypj5aTAVWUtuOF9fNE6VSf8LA0uhfi0EVk3dWoJ5dCSbkPXybn6t9nYbLu_Bpgsm-C21_sljh_gOH6XnyJcqgbkuswW-qDOHA3kXplufghQutnrW7jQO3zVyw-DQoJmApO6YsnfwewQB8UU4P7H_z87ubQaa_WxLO0pDbPuwudDolPpR2XZglKMvp_lK6K3tRyZl6ytUfc_LpLiQBReptw"},
		"pro-scp0": {
			Name:  "pro-scp0",
			Url:   "https://api.pro-scp0.sf-rz.de:6443",
			Token: ""},
		"cid-apc0": {
			Name:  "cid-apc0",
			Url:   "https://console.cid-apc0.sf-rz.de:8443",
			Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi16cTZ3ZyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJlMTJiNjI0Zi1kOTAwLTExZTgtODI5NS0wMDUwNTY5MDM3ZmIiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.soU7K1s7wMSW63VLMrxMUMo2jZIOp8pbf11biz7RFmylpHdvPx72yETs2cBevncxfhcmJXOnWtiSBtaxpHsQwbsMLwy5vNR1Aoajy7uR-_TktWFhfQu6ak-fuoHOPZml5dL8WZJCR_wvkov40k1kNeBSRjH0aXd_YUECk8jOmn9kxHWmHjcTuhoF8_mH9UCU7fGWPMa0ahUrllZlttqf1ZQcmk4oLi4X2JGIHN6pG9hQV0nOqutmkxdbshiH4od-aKljh5sX1pcql5NK9FliwFYKYPfRB8QRh5SBftu_VxGbqfRdrjdxbKgZzoZk002o5S-PYegagxZZAqzu9wti2Q"},
		"int-apc0": {
			Name:  "int-apc0",
			Url:   "https://console.int-apc0.sf-rz.de:8443",
			Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1nNjdxcCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJhYmY2ODM4Ni0wNDQ1LTExZTktOGFmYi0wMDUwNTY5MGQ4MjEiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.Tb5e_b3QiBUHmKzBOS0h6r_NeIw_GMERoGcEuyEjSOk-LZPU6q8atsAOEYg9QZdlo0rghxZhPGl8plxk-EVWgkWhQZpOVpp1HOZj8w097E9vRfRRkpt1e0I6ozq_PBx0K8758lxWa-LDQGV1vf8yTGpEUupps08ElSgwApgoVm6z1ZkmiBl_B_jb-MAH5PYnbiAQWu57AE8WXr_rOU8EunYK0fLJcBljMU7bhyxKW1WChIe7COYATCggQuRsu22ySF7dlTeGq3NLAFe2BpmeJavLle8FUQ1pjOem0uKrPEGat2H03KRKyuSP-7t34Jue5kdd3Qwf4rOraiVJ6-3T3w"},
		"ppr-apc0": {
			Name:  "ppr-apc0",
			Url:   "https://console.ppr-apc0.sf-rz.de:8443",
			Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi1qMnJ4eCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJjMDdkM2ZkNi0wNDRlLTExZTktYmJiYi0wMDUwNTZhMTkwNjQiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.OHHuZp-90WOBCaXjwTCCbvxW87lSsaBaH7q59JSDLM5oSWmca_aUbA3J_2nB6ruOOrvEbbQayb-UMcGGXhiVUv27Cg9TVb9Gl3AZym1gk6DUhma2wB3IwI98RKhX3XMr2J7pfte00p1Za8jBcVx_vdYaxebDQImpyUOFaEC3MT4ueirRXLTNSOQLnMARTVpr5r3s0Z0W3g8MGre7MDMNOO4BcGyEC1X8O8XzObuekzSErjJ6zugQMjIYw7_YkPULYNTJLrGsd136TFfFNGEtXow2RuLynIU8_1Zatvs64zNRirgvplGwxg4f15cHp7yUFJ_q3krkUegRji5tZdOBJw"},
		"vpt-apc0": {
			Name:  "pro-apc0",
			Url:   "https://console.pro-apc0.sf-rz.de:8443",
			Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi16N3NrdyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI3ZWY0YjE4OC0zNDM2LTExZTktOTIzZC0wMDUwNTY4MjA2Y2UiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.PX5Ghkmc4CgsvQgXuuIBBS95W0tE0YsRuCbymp4eVUBo4l1woEkKcds7yghx9PwcLHNMWWVki1DUUw-TbsWJY5MrnNwWCGwxTixeh3HSv3RBNAOkGUX4JPt6Bs0jEL_sELccLq9df6pZTQyDjidoSFnT4rid7R-7NyoRjlm0PhZRXEKxlGDt2eZNT0OJGl7btWqh1I_R7MTeuvya-KGSaMtNOlgpvLeyfH2UAvSY6INCu9ca88CBNWScpvZ_X6uxF75bLVGb1UJurPMbrG0isWge6kVnCaCpX_bSmlr8nyo7N6L1feGEgZZpwquTugbwfte8fLvbBgWCB9d8pt-ZhA"},
		"pro-apc0": {
			Name:  "pro-apc0",
			Url:   "https://console.pro-apc0.sf-rz.de:8443",
			Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbHVzdGVyLXRhc2tzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImltYWdlLXBydW5lci10b2tlbi16N3NrdyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJpbWFnZS1wcnVuZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI3ZWY0YjE4OC0zNDM2LTExZTktOTIzZC0wMDUwNTY4MjA2Y2UiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2x1c3Rlci10YXNrczppbWFnZS1wcnVuZXIifQ.PX5Ghkmc4CgsvQgXuuIBBS95W0tE0YsRuCbymp4eVUBo4l1woEkKcds7yghx9PwcLHNMWWVki1DUUw-TbsWJY5MrnNwWCGwxTixeh3HSv3RBNAOkGUX4JPt6Bs0jEL_sELccLq9df6pZTQyDjidoSFnT4rid7R-7NyoRjlm0PhZRXEKxlGDt2eZNT0OJGl7btWqh1I_R7MTeuvya-KGSaMtNOlgpvLeyfH2UAvSY6INCu9ca88CBNWScpvZ_X6uxF75bLVGb1UJurPMbrG0isWge6kVnCaCpX_bSmlr8nyo7N6L1feGEgZZpwquTugbwfte8fLvbBgWCB9d8pt-ZhA"},
	},
}

var FamilyNamespaces = T_famNs{
	"pkp": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-apc0": {"ms-jenkins", "openshift", "images-pkp"},
			// "cid-scp0": {},
			"int-apc0": {"images-pkp"},
			"ppr-apc0": {"images-pkp"},
			"vpt-apc0": {"images-pkp"},
			"pro-apc0": {"images-pkp"},
		},
		Stages:      []T_clName{"cid-apc0", "int-apc0", "ppr-apc0", "pro-apc0"},
		Buildstages: []T_clName{"cid-apc0"},
		Teststages:  []T_clName{"int-apc0", "ppr-apc0", "vpt-apc0"},
		Prodstage:   "pro-apc0",
	},
	"ssp": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-apc0": {"ssp-jenkins", "images-ssp"},
			"int-apc0": {"images-ssp"},
			"ppr-apc0": {"images-ssp"},
			"vpt-apc0": {"images-ssp"},
			"pro-apc0": {"images-ssp"},
		},
		Stages:      []T_clName{"cid-apc0", "int-apc0", "ppr-apc0", "pro-apc0"},
		Buildstages: []T_clName{"cid-apc0"},
		Teststages:  []T_clName{"int-apc0", "ppr-apc0", "vpt-apc0"},
		Prodstage:   "pro-apc0",
	},
	"aps": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-apc0": {"aps-jenkins", "images-aps"},
			"int-apc0": {"images-aps"},
			"ppr-apc0": {"images-aps"},
			"vpt-apc0": {"images-aps"},
			"pro-apc0": {"images-aps"},
		},
		Stages:      []T_clName{"cid-apc0", "int-apc0", "ppr-apc0", "pro-apc0"},
		Buildstages: []T_clName{"cid-apc0"},
		Teststages:  []T_clName{"int-apc0", "ppr-apc0", "vpt-apc0"},
		Prodstage:   "pro-apc0",
	},
	"vps": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-apc0": {"zpe-jenkins", "images-vps", "images-dfn"},
			"int-apc0": {"images-vps"},
			"ppr-apc0": {"images-vps"},
			"vpt-apc0": {"images-vps"},
			"pro-apc0": {"images-vps"},
		},
		Stages:      []T_clName{"cid-apc0", "int-apc0", "ppr-apc0", "pro-apc0"},
		Buildstages: []T_clName{"cid-apc0"},
		Teststages:  []T_clName{"int-apc0", "ppr-apc0", "vpt-apc0"},
		Prodstage:   "pro-apc0",
	},
	"dca": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-apc0": {"images-dca"},
			"int-apc0": {"images-dca"},
			"ppr-apc0": {"images-dca"},
			"vpt-apc0": {"images-dca"},
			"pro-apc0": {"images-dca"},
		},
		Stages:      []T_clName{"cid-apc0", "int-apc0", "ppr-apc0", "pro-apc0"},
		Buildstages: []T_clName{"cid-apc0"},
		Teststages:  []T_clName{"int-apc0", "ppr-apc0", "vpt-apc0"},
		Prodstage:   "pro-apc0",
	},
	"fpc": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-apc0": {"fpc-basis-1-1-20", "openshift", "fpc-basis-1-1-21"},
			"int-apc0": {"fpc-fa1", "fpc-fa2", "fpc-int1", "fpc-int2", "fpc-loadtest"},
			"ppr-apc0": {"fpc-ppr"},
			"vpt-apc0": {"vptest-fpc"},
			"pro-apc0": {"fpc", "vptest-fpc"},
		},
		Stages:      []T_clName{"cid-apc0", "int-apc0", "ppr-apc0", "pro-apc0"},
		Buildstages: []T_clName{"cid-apc0"},
		Teststages:  []T_clName{"int-apc0", "ppr-apc0", "vpt-apc0"},
		Prodstage:   "pro-apc0",
	},
	"fpc-scp0": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-scp0": {},
			"ppr-scp0": {"fpc-ppr"},
			"vpt-scp0": {},
			// "pro-scp0": {},
		},
		Stages:      []T_clName{"cid-scp0", "ppr-scp0", "vpt-scp0"},
		Buildstages: []T_clName{"cid-scp0"},
		Teststages:  []T_clName{"cid-scp0", "ppr-scp0", "vpt-scp0"},
		Prodstage:   "ppr-scp0",
	},
	"hub": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-apc0": {"images-hub", "hub-dev1", "hub-dev2"},
			"int-apc0": {},
			"ppr-apc0": {"images-hub", "ppr-hub-mvp"},
			"vpt-apc0": {"images-hub", "vpt-hub-mvp"},
			"pro-apc0": {"images-hub", "pro-hub"},
		},
		Stages:      []T_clName{"cid-apc0", "ppr-apc0", "vpt-apc0", "pro-apc0"},
		Buildstages: []T_clName{"cid-apc0"},
		Teststages:  []T_clName{"vpt-apc0"},
		Prodstage:   "pro-apc0",
	},
	"hub-scp0": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-scp0": {"images-hub", "hub-dev1", "hub-dev2"},
			"ppr-scp0": {},
			"vpt-scp0": {},
			"pro-scp0": {},
		},
		Stages: []T_clName{"cid-scp0", "ppr-scp0", "vpt-scp0", "pro-scp0"},
		// Stages:     []T_clName{"cid-apc0", "cid-scp0"},
		Buildstages: []T_clName{"cid-scp0"},
		Teststages:  []T_clName{"cid-scp0", "ppr-scp0", "vpt-scp0"},
		Prodstage:   "pro-scp0",
	},
	"base": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-apc0": {"scpas-jenkins", "scptools-jenkins", "scp-operations", "uniserv"},
			"int-apc0": {"scp-operations", "uniserv"},
			"ppr-apc0": {"scp-operations", "uniserv"},
			"vpt-apc0": {"scp-operations", "uniserv"},
			"pro-apc0": {"scp-operations", "uniserv"},
			"cid-scp0": {"scp-operations", "scp-baseimages"},
			"ppr-scp0": {"scp-operations"},
			"vpt-scp0": {"scp-operations"},
			"pro-scp0": {"scp-operations"},
			// "cid-scp0": {"scp-operations", "scp-images", "scp-baseimages"},
			// "ppr-scp0": {"scp-operations", "scp-images"},
			// "vpt-scp0": {"scp-operations", "scp-images"},
			// "pro-scp0": {"scp-operations", "scp-images"},
		},
		Stages:      []T_clName{"cid-apc0", "int-apc0", "ppr-apc0", "pro-apc0", "cid-scp0", "ppr-scp0", "vpt-scp0"},
		Buildstages: []T_clName{"cid-apc0", "cid-scp0"},
		Teststages:  []T_clName{"int-apc0", "ppr-apc0", "vpt-apc0"},
		Prodstage:   "pro-apc0",
	},
	"base-scp0": {
		ClusterNamespaces: map[T_clName][]T_nsName{
			"cid-scp0": {"scp-operations", "scpas-jenkins", "scp-baseimages"},
			"ppr-scp0": {"scp-operations"},
			"vpt-scp0": {"scp-operations"},
			"pro-scp0": {"scp-operations"},
		},
		Stages:      []T_clName{"cid-scp0", "ppr-scp0", "vpt-scp0"},
		Buildstages: []T_clName{"cid-scp0"},
		Teststages:  []T_clName{"cid-scp0", "ppr-scp0", "vpt-scp0"},
		Prodstage:   "vpt-scp0",
	},
}

var OcClient bool

var LogFileName = "image-tools.log"
