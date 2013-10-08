package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

const updateURL = `https://api.github.com/repos/ernestokarim/cb/commits?per_page=1`

func init() {
	registry.NewUserTask("update", 0, update)
	registry.NewTask("update:check", 0, updateCheck)
}

func update(c *config.Config, q *registry.Queue) error {
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
		if *config.Verbose {
			log.Printf("local or remote version was not retrieved correctly\n")
		}
		return nil
	}

	if err := writeCheckUpdate(); err != nil {
		return err
	}

	// No update, return directly
	if latestSha == currentSha {
		if *config.Verbose {
			log.Printf("same version detected\n")
		}
		return nil
	}

	// Perform the update
	args := []string{"get", "-u", "github.com/ernestokarim/cb"}
	output, err := utils.Exec("go", args)
	if err != nil {
		return err
	}
	if len(output) > 0 {
		fmt.Println(output)
	}

	log.Printf("%sUpdated correctly to commit: %s%s", colors.Green, latestSha[:10], colors.Reset)

	return nil
}

type commitInfo struct {
	Sha string
}

func updateCheck(c *config.Config, q *registry.Queue) error {
	// Check update-check file before updating again
	shouldCheck, err := checkShouldCheckUpdate()
	if err != nil {
		return err
	}
	if !shouldCheck {
		return nil
	}

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
		if *config.Verbose {
			log.Printf("local or remote version was not retrieved correctly\n")
		}
		return nil
	}

	if err := writeCheckUpdate(); err != nil {
		return err
	}

	// No update, return directly
	if latestSha == currentSha {
		if *config.Verbose {
			log.Printf("same version detected\n")
		}
		return nil
	}

	// Ask for update
	if utils.Ask("There's a new CB version. Do you want to auto-update it?") {
		return q.RunTasks(c, []string{"update@0"})
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

func checkShouldCheckUpdate() (bool, error) {
	p := config.GetUserConfigsPath()
	info, err := os.Stat(filepath.Join(p, "update-check"))
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("cannot stat update check file: %s", err)
	}

	if err == nil && time.Now().Sub(info.ModTime()) < 24*time.Hour {
		if *config.Verbose {
			log.Printf("ignoring update because it has been checked in the last 24 hours\n")
		}
		return false, nil
	}

	return true, nil
}

func writeCheckUpdate() error {
	p := config.GetUserConfigsPath()
	f, err := os.Create(filepath.Join(p, "update-check"))
	if err != nil {
		return fmt.Errorf("cannot create update check file: %s", err)
	}
	defer f.Close()

	if *config.Verbose {
		log.Printf("writing update-check file\n")
	}

	return nil
}
