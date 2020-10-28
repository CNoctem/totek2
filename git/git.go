package git

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing"
	"log"
	"os"
	"totek2/util"
)
import "github.com/go-git/go-git/v5"

func CloneProject(workspace, url, branch string, user *util.User) (*string, error) {
	log.Println(url + " -> " + workspace)
	log.Println("Cloning as user " + user.Username)
	r, err := tryClone(url, workspace, branch)
	var commitid string
	if err != nil {
		if err.Error() == "authentication required" {
			if len(user.Password) == 0 {
				*user = util.Authenticate()
			}
			newUrl, err := util.URLWithUser(url, *user)
			if err != nil {
				log.Fatal(err)
			}
			//url = fmt.Sprintf("https://%s:%s@%s", user.Username, user.Password, url)
			r, err = tryClone(*newUrl, workspace, branch)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	ref, err := r.Head()
	if err != nil {
		log.Fatal(err)
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatal(err)
	}

	commit, err := cIter.Next()
	if err != nil {
		log.Fatal(err)
	}
	commitid = commit.Hash.String()

	fmt.Printf("commitid %s for %s %d\n", commitid, url, len(commitid))
	return &commitid, nil
}

func tryClone(url string, path string, branch string) (*git.Repository, error) {
	if len(branch) == 0 {
		r, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		return r, err
	}
	r, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})

	return r, err
}
