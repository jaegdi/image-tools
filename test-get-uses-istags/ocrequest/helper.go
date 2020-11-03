package ocrequest

import (
	"github.com/imdario/mergo"
	"log"
	"os"
	"time"
)

var today time.Time = time.Now()

func ageInDays(date string) int {
	t, _ := time.Parse(time.RFC3339, date)
	return int(today.Sub(t).Hours()) / 24
}

func mergeMaps(dest interface{}, source interface{}, msg string) interface{} {
	if err := mergo.Merge(&dest, source); err != nil {
		log.Println("ERROR: " + msg + ": failed: " + err.Error())
		return source
	}
	return dest
}

func ExitWithError(errormsg string) {
	log.Println(errormsg)
	os.Exit(1)
}
