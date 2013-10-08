package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

const updateURL = `https://api.github.com/repos/ernestokarim/cb/commits?per_page=1`

func init() {
	registry.NewUserTask("self-update", 0, update)
	registry.NewUserTask("update", 0, update)
	registry.NewUserTask("update:check", 0, updateCheck)
}

func update(c *config.Config, q *registry.Queue) error {
	fmt.Println("Hello World!")
	return nil
}

type commitInfo struct {
	Sha string
}

func updateCheck(c *config.Config, q *registry.Queue) error {
	// Fetch last commits, both localy & remotely
	latestSha, err := fetchLatestCommit()
	if err != nil {
		return err
	}
	currentSha, err := fetchCurrentCommit()
	if err != nil {
		return err
	}

	// Couldn't retrieve current/latest commit, ignore update
	if latestSha == "" || currentSha == "" {
		return nil
	}

	// No update, return directly
	if latestSha == currentSha {
		return nil
	}

	// Ask for update
	if utils.Ask("There's a new CB version. Do you want to auto-update it?") {
		q.AddTask("update@0")
	}

	return nil
}

func fetchLatestCommit() (string, error) {
	resp, err := http.Get(updateURL)
	if err != nil {
		// If there's no Internet connection, don't return an error
		if e, ok := err.(*url.Error); ok {
			if e.Err.Error() == "dial tcp: lookup api.github.com: no such host" {
				log.Printf("%scannot check for updates, there's no connection%s\n",
					colors.Yellow, colors.Reset)
				return "", nil
			}
		}

		// Fatal error otherwise
		return "", fmt.Errorf("cannot check update url: %T", err)
	}
	defer resp.Body.Close()

	// Extract the commit info
	var data []*commitInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("cannot decode github api data: %s", err)
	}

	return data[0].Sha, nil
}

func fetchCurrentCommit() (string, error) {
	path := utils.PackagePath("github.com/ernestokarim/cb")

	args := []string{
		"--git-dir", filepath.Join(path, ".git"),
		"--work-tree", path,
		"rev-parse",
		"HEAD",
	}
	output, err := utils.Exec("git", args)
	if err != nil {
		return "", fmt.Errorf("cannot parse git head revision: %s", err)
	}

	return strings.TrimSpace(output), nil
}
