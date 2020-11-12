package ocrequest

import (
	// "log"
	// "encoding/json"
	// . "fmt"
	"github.com/imdario/mergo"
	// "reflect"
)

// type T_M map[string]interface{}

// type T map[string]interface{}

// var map1 = T{
// 	"A": T_M{
// 		"a": []int{3, 5, 7, 11},
// 		"b": []string{"xx", "yy"},
// 	},
// 	"B": T_M{
// 		"c": 3,
// 		"d": "xx",
// 	},
// 	"C": T_M{
// 		"c": true,
// 		"m": nil,
// 	},
// }
// var map2 = T{
// 	"A": T_M{
// 		"e": 5,
// 		"a": []int{1, 2, 3, 4, 5},
// 		"b": []string{"ww", "zz"},
// 	},
// 	"C": T_M{
// 		"c": 4,
// 		"d": "nil",
// 	},
// }

// MergoNestedMaps merges two maps.
// As first parameter the __address__ of an empty map of the same type as the two maps, they should be merged, must given.
// the second parameter is the destination map, the third parameter is the map that should be merged into destination map.
// The merged result is in the first parameter. After execution of this function, the var of the first parameter must
// assigned to the destination map.
// Example:
// t := T_usedIstagsResult{}; MergoNestedMaps(&t, result, r); result = t
func MergoNestedMaps(dest, m1, m2 interface{}) {
	if err := mergo.Merge(dest, m1, mergo.WithSliceDeepCopy, mergo.WithAppendSlice, mergo.WithTypeCheck); err != nil {
		ErrorLogger.Println("merge m1 m2" + ": failed: " + err.Error())
	}
	if err := mergo.Merge(dest, m2, mergo.WithSliceDeepCopy, mergo.WithAppendSlice, mergo.WithTypeCheck); err != nil {
		ErrorLogger.Println("merge m1 m2" + ": failed: " + err.Error())
	}
}

// func Test_MergeNestedMaps() {
// 	d := T{}
// 	MergoNestedMaps(&d, map1, map2)
// 	Println("Map: ", d)
// 	JsonStr, err := json.MarshalIndent(d, "", "  ")
// 	if err != nil {
// 		Println("Json: Marshal error", err)
// 	}
// 	Println("Json: ", string(JsonStr))
// }
