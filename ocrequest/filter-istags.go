package ocrequest

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// matchIsIstagToFilterParams returns true when the filters are empty or a defined filter matches its corresponding item
//
// Parameters:
// - is: The image stream name (T_isName).
// - tag: The tag name (T_tagName).
// - istag: The image stream tag name (T_istagName).
// - namespace: The namespace name (T_nsName).
// - age: The age of the tag in days.
//
// Returns:
// - A boolean indicating whether the item matches the filter parameters.
func matchIsIstagToFilterParams(is T_isName, tag T_tagName, istag T_istagName, namespace T_nsName, age int) bool {
	VerifyMsg("filtering:", is, tag, istag, namespace)
	return ((CmdParams.Filter.Isname == "" ||
		(CmdParams.Filter.Isname != "" && is == CmdParams.Filter.Isname) ||
		(CmdParams.Filter.Isname != "" && CmdParams.FilterReg.Isname.MatchString(string(is)))) &&
		(CmdParams.Filter.Tagname == "" ||
			(CmdParams.Filter.Tagname != "" && tag == CmdParams.Filter.Tagname) ||
			(CmdParams.Filter.Tagname != "" && CmdParams.FilterReg.Tagname.MatchString(string(tag)))) &&
		(CmdParams.Filter.Istagname == "" ||
			(CmdParams.Filter.Istagname != "" && istag == CmdParams.Filter.Istagname) ||
			(CmdParams.Filter.Istagname != "" && CmdParams.FilterReg.Istagname.MatchString(string(istag)))) &&
		(CmdParams.Filter.Namespace == "" ||
			(CmdParams.Filter.Namespace != "" && namespace == CmdParams.Filter.Namespace) ||
			(CmdParams.Filter.Namespace != "" && CmdParams.FilterReg.Namespace.MatchString(string(namespace)))) &&
		(CmdParams.Filter.Minage == -1 ||
			(CmdParams.Filter.Minage > -1 && age >= CmdParams.Filter.Minage)) &&
		(CmdParams.Filter.Maxage == -1 ||
			(CmdParams.Filter.Maxage > -1 && age <= CmdParams.Filter.Maxage)))
}

// logUsedIstags logs the details of usedIstags to the logfile
//
// Parameters:
// - usedIstags: A slice of T_usedIstag containing the used image stream tags.
// - is: The image stream name (T_isName).
// - tag: The tag name (T_tagName).
// - istag: The image stream tag name (T_istagName).
func logUsedIstags(usedIstags []T_usedIstag, is T_isName, tag T_tagName, istag T_istagName) {
	VerifyMsg("logUsedIstags::", " ## is: ", is, " ### tag: ", tag, "  #### Istag: ", istag, " is used.")
	// Iterate through each used image stream tag and log its details
	for _, istagdetails := range usedIstags {
		VerifyLogger.Println("logUsedIstags::", "   # -->",
			"  Cluster:", istagdetails.Cluster,
			"  UsedInNamespace:", istagdetails.UsedInNamespace,
			"  FromNamespace:", istagdetails.FromNamespace,
			"  Image:", istagdetails.Image,
			"  AgeInDays:", istagdetails.AgeInDays)
	}
}

// printShellCmds prints a map of shell commands sorted by the map key
//
// Parameters:
// - commands: A map of shell commands with the command description as the key and the command as the value.
func printShellCmds(commands map[string]string) {
	keys := make([]string, 0, len(commands))
	DebugMsg("printShellCmds::", "printShellCmds", commands)
	// Collect all keys from the commands map
	for key := range commands {
		keys = append(keys, key)
	}
	// Sort the keys
	sort.Strings(keys)
	// Print each command in the sorted order
	for _, key := range keys {
		fmt.Print(commands[key])
	}
}

// isImageRefencedByLatestTag checks if an image is referenced by a "latest" tag
//
// Parameters:
// - image: A map of image stream tag names (T_istagName) to their corresponding SHA (T_sha).
//
// Returns:
// - A boolean indicating whether the image is referenced by a "latest" tag.
func isImageRefencedByLatestTag(image map[T_istagName]T_sha) bool {
	tagLatestRegexp := regexp.MustCompile("^.*?:latest$")
	// Iterate through each image stream tag name
	for istag := range image {
		// Check if the tag name matches the "latest" pattern
		if tagLatestRegexp.MatchString(istag.str()) {
			return true
		}
	}
	return false
}

// FilterIstagsToDelete generates shell commands to delete image stream tags that fit the specified conditions
// The conditions are:
// - The tag must fit the tagPattern
// - The tag must be older or equal to the given minAge
// - The image stream tag must not be used in any of the clusters by any of: dc, pod, job, cronjob, or daemonset
// - The image stream tag must not reference an image that is referenced by a "latest" tag
//
// Parameters:
// - data: The complete results for all families (T_completeResultsFamilies).
// - family: The family name (T_familyName).
// - clusters: A slice of cluster names (T_clNames).
// - tagPattern: A string representing the tag pattern to match.
// - minAge: An integer representing the minimum age of the tag in days.
// - cause: A string representing the cause for deletion.
func FilterIstagsToDelete(data T_completeResultsFamilies, family T_familyName, clusters T_clNames, tagPattern string, minAge int, cause string) {
	result := map[string]string{}
	tagPatternRegexp := regexp.MustCompile(tagPattern)
	// Iterate through each cluster in the list
	for _, cluster := range clusters {
		// Iterate through each image stream tag in the current cluster
		for istag, nsTags := range data[family].AllIstags[cluster].Istag {
			is, tag := istag.split()
			// Check if the tag matches the pattern and is not "latest"
			if (tagPatternRegexp.MatchString(istag.str()) || tagPattern == "") && tag != "latest" {
				// Iterate through each namespace in the current image stream tag
				for ns, tagMap := range nsTags {
					DebugMsg("FilterIstagsToDelete::", "ns:", ns, "tagMap:", GetJsonFromMap(tagMap))
					// Check if the tag meets the age and filter parameters
					if tagMap.AgeInDays >= minAge && matchIsIstagToFilterParams(is, tag, istag, tagMap.Namespace, tagMap.AgeInDays) {
						// Check if the tag is not used and does not reference an image used by a "latest" tag
						if data[family].UsedIstags[is][tag] == nil && !isImageRefencedByLatestTag(data[family].AllIstags[cluster].Image[tagMap.Image]) {
							s := (string(ns) + "/" + string(istag))
							value := fmt.Sprintln(
								"oc -n", tagMap.Namespace, "delete istag", tagMap.Imagestream.str()+":"+tagMap.Tagname.str(),
								"   #", cause, "-->", tagMap.Image,
								",  Commit.Ref:", tagMap.Build.CommitRef,
								",  Age:", tagMap.AgeInDays)
							DebugMsg("FilterIstagsToDelete::", "key:", s, "value:", value)
							// Add the delete command to the result map
							result[s] = value
						} else {
							// Log the used image stream tags
							logUsedIstags(data[family].UsedIstags[is][tag], is, tag, istag)
						}
					}
				}
			}
		}
	}
	// Print the generated shell commands
	printShellCmds(result)
}

// FilterNonbuildIstagsToDelete filters out all image stream tags when there is no build-tag on the same image
//
// Parameters:
// - data: The complete results for all families (T_completeResultsFamilies).
// - family: The family name (T_familyName).
// - clusters: A slice of cluster names (T_clNames).
// - minAge: An integer representing the minimum age of the tag in days.
func FilterNonbuildIstagsToDelete(data T_completeResultsFamilies, family T_familyName, clusters T_clNames, minAge int) {
	result := map[string]string{}
	buildPatternRegexp := regexp.MustCompile("^.*?:.*?[A-Za-z]")
	// Iterate through each cluster in the list
	for _, cluster := range clusters {
		// Iterate through each image in the current cluster
		for image, tags := range data[family].AllIstags[cluster].Image {
			hasBuildTag := false
			var istagname T_istagName
			// Iterate through each tag in the current image
			for istagName := range tags {
				if !hasBuildTag {
					// Check if the tag matches the build pattern
					hasBuildTag = buildPatternRegexp.MatchString(istagName.str())
				}
				istagname = istagName
			}
			// If no build tag is found and the tag is older than the minimum age
			if !hasBuildTag {
				if tags[istagname].AgeInDays > minAge {
					// Iterate through each image stream tag in the current tag
					for tn := range tags[istagname].Istags {
						tn := tn
						imageParts := strings.Split(tn.str(), "/")
						fromNamespace := T_nsName(imageParts[len(imageParts)-2])
						istag := T_istagName(imageParts[len(imageParts)-1])
						is, tag := istag.split()
						// Check if the tag meets the filter parameters
						if matchIsIstagToFilterParams(is, tag, istag, fromNamespace, tags[istagname].AgeInDays) {
							// Check if the tag is not used
							if data[family].UsedIstags[is][tag] == nil {
								// Add the delete command to the result map
								result[tn.str()] = fmt.Sprintln(
									"oc -n", fromNamespace, "delete istag", istag,
									"   # nonbuild -->", image,
									",  Commit.Ref:", data[family].AllIstags[cluster].Istag[istag][fromNamespace].Build.CommitRef,
									",  Age:", tags[istagname].AgeInDays)
							} else {
								// Log the used image stream tags
								logUsedIstags(data[family].UsedIstags[is][tag], is, tag, istag)
							}
						}
					}
				}
			}
		}
	}
	// Print the generated shell commands
	printShellCmds(result)
}

// FilterAllIstags removes all parts from the complete result that are not specified for output in the CmdParams.Output flags
//
// Parameters:
// - result: A pointer to the complete results (T_completeResults).
func FilterAllIstags(result *T_completeResults) {
	outputflags := CmdParams.Output
	if !outputflags.All {
		// Iterate through each cluster in the family namespaces
		for _, cluster := range FamilyNamespaces[CmdParams.Family].Stages {
			x := result.AllIstags[cluster]
			if !CmdParams.Delete {
				// Remove parts of the result based on the output flags
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
		// Remove used image stream tags if not specified for output
		if !outputflags.Used && !CmdParams.Delete {
			result.UsedIstags = T_usedIstagsResult{}
		}
	}
}

// FilterUnusedIstags finds all image stream tags (istags) that exist in the defined clusters but are not used in those clusters.
// It updates the UnUsedIstags field in the provided T_completeResults struct.
//
// Parameters:
// - result: A pointer to a T_completeResults struct that contains all istags and used istags information.
func FilterUnusedIstags(result *T_completeResults) {
	// Initialize the UnUsedIstags field as an empty T_unUsedIstagsResult map
	result.UnUsedIstags = T_unUsedIstagsResult{}
	// Iterate through each cluster defined in the command parameters
	for _, cluster := range CmdParams.Cluster {
		// Get all istags for the current cluster
		istags := result.AllIstags[cluster].Istag
		// Get the used istags
		used := result.UsedIstags
		// Iterate through each istag in the current cluster
		for x := range istags {
			// Extract the image stream (is) and tag from the istag
			_, _, is, tag, _ := GetIsAndTag(x.str())
			// Check if the istag is not used
			if used[is][tag] == nil {
				// If the istag is unused, create a new T_unUsedIstag struct
				u := T_unUsedIstag{}
				// Set the cluster field in the T_unUsedIstag struct
				u.Cluster = cluster
				// Add the unused istag to the UnUsedIstags map in the result
				result.UnUsedIstags[x] = u
			}
		}
	}
}
