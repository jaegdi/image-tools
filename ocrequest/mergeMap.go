package ocrequest

import (
	// "log"
	"encoding/json"
	. "fmt"
	"github.com/imdario/mergo"
	// "reflect"
)

type T_M map[string]interface{}

type T map[string]interface{}

var map1 = T{
	"A": T_M{
		"a": []int{3, 5, 7, 11},
		"b": []string{"xx", "yy"},
	},
	"B": T_M{
		"c": 3,
		"d": "xx",
	},
	"C": T_M{
		"c": true,
		"m": nil,
	},
}
var map2 = T{
	"A": T_M{
		"e": 5,
		"a": []int{1, 2, 3, 4, 5},
		"b": []string{"ww", "zz"},
	},
	"C": T_M{
		"c": 4,
		"d": "nil",
	},
}

func MergoNestedMaps(dest, m1, m2 interface{}) {
	if err := mergo.Merge(dest, m1, mergo.WithSliceDeepCopy, mergo.WithAppendSlice, mergo.WithTypeCheck); err != nil {
		ErrorLogger.Println("merge m1 m2" + ": failed: " + err.Error())
	}
	if err := mergo.Merge(dest, m2, mergo.WithSliceDeepCopy, mergo.WithAppendSlice, mergo.WithTypeCheck); err != nil {
		ErrorLogger.Println("merge m1 m2" + ": failed: " + err.Error())
	}
}

func Test_MergeNestedMaps() {
	d := T{}
	MergoNestedMaps(&d, map1, map2)
	Println("Map: ", d)
	JsonStr, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		Println("Json: Marshal error", err)
	}
	Println("Json: ", string(JsonStr))
}
