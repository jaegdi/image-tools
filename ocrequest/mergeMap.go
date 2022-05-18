package ocrequest

import (
	"github.com/imdario/mergo"
)

// MergoNestedMaps merges two maps.
// As first parameter the __address__ of an empty map of the same type as the two maps, they should be merged, must given.
// the second parameter is the destination map, the third parameter is the map that should be merged into destination map.
// The merged result is in the first parameter. After execution of this function, the var of the first parameter must
// assigned to the destination map.
// Example:
// t := T_usedIstagsResult{}; MergoNestedMaps(&t, result, r); result = t
func MergoNestedMaps(dest, m1 interface{}) {
	if err := mergo.Merge(dest, m1); err != nil {
		ErrorLogger.Println("merge m1: failed:", err)
	}
}
