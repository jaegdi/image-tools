package ocrequest

import (
	// "fmt"
	"log"
	"os/exec"
)

func OcLogin(cluster string) (string, error) {
	app := "ocl"
	cmd := exec.Command(app, cluster)
	stdout, err := cmd.Output()

	if err != nil {
		log.Println(err.Error() + ":" + string(stdout))
		return "login failed", err
	} else {
		cmd := exec.Command("oc", "whoami", "-t")
		token, err := cmd.Output()
		if err != nil {
			return "token failed:" + string(token), err
		}
		// token has a linefeed, that must removed
		return string(token[:len(token)-1]), err
	}
}
