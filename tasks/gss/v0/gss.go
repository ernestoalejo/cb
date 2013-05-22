package v0

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("gss", 0, gss)
}

func gss(c *config.Config, q *registry.Queue) error {
	if !*config.ClosureMode {
		return fmt.Errorf("closure mode only task")
	}

	compilerPath, err := c.Get("closure.stylesheets")
	if err != nil {
		return fmt.Errorf("cannot get compiler path: %s", err)
	}
	compilerPath = filepath.Join(compilerPath, "build", "closure-stylesheets.jar")

	size := c.CountRequired("gss")
	for i := 0; i < size; i++ {
		src, err := c.GetStringf("gss[%d].source", i)
		if err != nil {
			return fmt.Errorf("get gss source failed: %s", err)
		}
		dest, err := c.GetStringf("gss[%d].dest", i)
		if err != nil {
			return fmt.Errorf("get gss dest failed: %s", err)
		}

		dir := filepath.Dir(dest)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create dest dir failed (%s): %s", dir, err)
		}
		args := []string{
			"-jar", compilerPath,
			"--output-renaming-map-format", "CLOSURE_COMPILED",
			"--rename", "CLOSURE",
			"--output-renaming-map", filepath.Join("temp", "gssmap.js"),
			"--output-file", dest,
			filepath.Join("temp", "styles", src),
		}
		output, err := utils.Exec("java", args)
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("compiler error: %s", err)
		}
		if *config.Verbose {
			log.Printf("created file %s\n", dest)
		}
	}

	return nil
}
