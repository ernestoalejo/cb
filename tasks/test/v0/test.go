package v0

import (
	"log"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("test", 0, test)
	registry.NewTask("e2e", 0, e2e)
}

func test(c config.Config, q *registry.Queue) error {
	var configFile string
	if *config.Compiled {
		configFile = "client/config/testacular-compiled.conf.js"
	} else {
		configFile = "client/config/testacular.conf.js"
	}

	if *config.Verbose {
		log.Printf("using config `%s`\n", configFile)
	}

	args := []string{"start", configFile}
	if err := utils.ExecCopyOutput("testacular", args); err != nil {
		return err
	}
	return nil
}

func e2e(c config.Config, q *registry.Queue) error {
	if *config.Compiled {
		q.AddTask("proxy")
	} else {
		q.AddTask("server")
	}
	return nil
}
