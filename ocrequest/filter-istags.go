package ocrequest

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func matchIsIstagToFilterParams(is T_isName, tag T_tagName, istag T_istagName, namespace T_nsName) bool {
	return ((CmdParams.Filter.Isname == "" || (CmdParams.Filter.Isname != "" && is == CmdParams.Filter.Isname)) &&
		(CmdParams.Filter.Tagname == "" || (CmdParams.Filter.Tagname != "" && tag == CmdParams.Filter.Tagname)) &&
		(CmdParams.Filter.Istagname == "" || (CmdParams.Filter.Istagname != "" && istag == CmdParams.Filter.Istagname)) &&
		(CmdParams.Filter.Namespace == "" || (CmdParams.Filter.Namespace != "" && namespace == CmdParams.Filter.Namespace)))
}

func logUsedIstags(usedIstags []T_usedIstag, is T_isName, tag T_tagName, istag T_istagName) {
	LogDebug("logUsedIstags::", "#### Istag:", istag, "is used.")
	for _, istagdetails := range usedIstags {
		LogDebug("logUsedIstags::", "   # -->",
			"Cluster:", istagdetails.Cluster,
			"UsedInNamespace:", istagdetails.UsedInNamespace,
			"FromNamespace:", istagdetails.FromNamespace,
			"Image:", istagdetails.Image,
			"AgeInDays:", istagdetails.AgeInDays)
	}
}

func printShellCmds(result map[string]string) {
	keys := make([]string, 0, len(result))
	LogDebug("printShellCmds::", "printShellCmds", result)
	for key := range result {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Print(result[key])
	}
}

func FilterIstagsToDelete(data T_completeResultsFamilies, family T_family, cluster T_clName, tagPattern string, minAge int, cause string) {
	result := map[string]string{}
	tagPatternRegexp := regexp.MustCompile(tagPattern)
	for istag, nsTags := range data[family].AllIstags[cluster].Istag {
		is, tag := istag.split()
		if tagPatternRegexp.MatchString(istag.str()) || tagPattern == "" {
			for ns, tagMap := range nsTags {
				if CmdParams.Options.Debug {
					LogDebug("FilterIstagsToDelete::", "ns:", ns, "tagMap:", GetJsonFromMap(tagMap))
				}
				if tagMap.AgeInDays >= minAge && matchIsIstagToFilterParams(is, tag, istag, tagMap.Namespace) {
					if data[family].UsedIstags[is][tag] == nil {
						s := (string(ns) + "/" + string(istag))
						value := fmt.Sprintln(
							"oc -n", tagMap.Namespace, "delete istag", tagMap.Imagestream.str()+":"+tagMap.Tagname.str(),
							"   #", cause, "-->", tagMap.Image,
							",  Commit.Ref:", tagMap.Build.CommitRef,
							",  Age:", tagMap.AgeInDays)
						LogDebug("FilterIstagsToDelete::", "key:", s, "value:", value)
						result[s] = value
					} else {
						logUsedIstags(data[family].UsedIstags[is][tag], is, tag, istag)
					}
				}
			}
		}
	}
	printShellCmds(result)
}

func FilterNonbuildIstagsToDelete(data T_completeResultsFamilies, family T_family, cluster T_clName, minAge int) {
	result := map[string]string{}
	buildPatternRegexp := regexp.MustCompile("^.*?:.*?[A-Za-z]")
	for image, tags := range data[family].AllIstags[cluster].Image {
		hasBuildTag := false
		var istagname T_istagName
		for istagName := range tags {
			if !hasBuildTag {
				hasBuildTag = buildPatternRegexp.MatchString(istagName.str())
			}
			istagname = istagName
		}
		if !hasBuildTag {
			if tags[istagname].AgeInDays > minAge {
				for tn := range tags[istagname].Istags {
					tn := tn
					imageParts := strings.Split(tn.str(), "/")
					fromNamespace := T_nsName(imageParts[len(imageParts)-2])
					istag := T_istagName(imageParts[len(imageParts)-1])
					is, tag := istag.split()
					if matchIsIstagToFilterParams(is, tag, istag, fromNamespace) {
						if data[family].UsedIstags[is][tag] == nil {
							result[tn.str()] = fmt.Sprintln(
								"oc -n", fromNamespace, "delete istag", istag,
								"   # nonbuild -->", image,
								",  Commit.Ref:", data[family].AllIstags[cluster].Istag[istag][fromNamespace].Build.CommitRef,
								",  Age:", tags[istagname].AgeInDays)
						} else {
							logUsedIstags(data[family].UsedIstags[is][tag], is, tag, istag)
						}
					}
				}
			}
		}
	}
	printShellCmds(result)
}

func FilterAllIstags(result *T_completeResults) {
	outputflags := CmdParams.Output
	if !outputflags.All {
		for _, cluster := range Clusters.Stages {
			x := result.AllIstags[cluster]
			if !CmdParams.Delete {
				if !outputflags.Is {
					x.Is = T_resIs{}
				}
				if !outputflags.Istag {
					x.Istag = T_resIstag{}
				}
				if !outputflags.Image {
					x.Image = T_resSha{}
				}
				result.AllIstags[cluster] = x
			}
		}
		if !outputflags.Used && !CmdParams.Delete {
			result.UsedIstags = T_usedIstagsResult{}
		}
	}
}

func FilterUnusedIstags(result *T_completeResults) {
	istags := result.AllIstags[CmdParams.Cluster].Istag
	for x := range istags {
		fmt.Println(x)
	}

}
