package ocrequest

import (
	//#"github.com/imdario/mergo"
	"dario.cat/mergo"
)

// MergoNestedMaps merges two nested maps using the mergo.Merge function.
// The destination map is modified to include the contents of the source map.
// If an error occurs during the merge, it is logged using the ErrorLogger.
//
// Parameters:
// - dest: The destination map that will be modified to include the contents of the source map.
// - m1: The source map whose contents will be merged into the destination map.
func MergoNestedMaps(dest, m1 interface{}) {
	if err := mergo.Merge(dest, m1); err != nil {
		ErrorMsg("merge m1: failed:", err)
	}
}
