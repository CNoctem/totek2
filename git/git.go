package git

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing"
	"log"
	"os"
	"totek2/util"
)
import "github.com/go-git/go-git/v5"

func CloneProject(workspace, url, branch string, user *util.User) error {
	log.Println(url + " -> " + workspace)
	err := tryClone("https://"+url, workspace, branch)
	if err != nil && err.Error() == "authentication required" {
		if len(user.Password) == 0 {
			*user = util.Authenticate()
		}
		url = fmt.Sprintf("https://%s:%s@%s", user.Username, user.Password, url)
		err = tryClone(url, workspace, branch)
		if err != nil {
			log.Fatal(err)
		}
	}
	return err
}

func tryClone(url string, path string, branch string) error {
	if len(branch) == 0 {
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		return err
	}
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})
	return err
}
