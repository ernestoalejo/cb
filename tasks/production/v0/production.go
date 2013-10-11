package v0

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

const selfPkg = "github.com/ernestokarim/cb/tasks/production/v0/scripts"

func init() {
	registry.NewUserTask("production", 0, production)
}

func production(c *config.Config, q *registry.Queue) error {
	scriptsPath := utils.PackagePath(selfPkg)

	log.Printf("Hashing local files... ")
	localHashes, err := hashLocalFiles()
	if err != nil {
		return fmt.Errorf("hash local files failed: %s", err)
	}
	log.Printf("Hashing local files... %s[SUCCESS]%s\n", colors.Green, colors.Reset)

	_ = scriptsPath
	_ = localHashes
	/*
	   args := []string{
	     filepath.Base(c.GetRequired("paths.base")),
	   }
	   if err := utils.ExecCopyOutput(base, args); err != nil {
	     return fmt.Errorf("deploy failed: %s", err)
	   }*/

	return nil
}

func hashLocalFiles() (map[string]string, error) {
	hashes := map[string]string{}
	rootPath := "../deploy"

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open file failed: %s", err)
		}
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("read file failed: %s", err)
		}
		h := sha1.New()
		if _, err := h.Write(content); err != nil {
			return fmt.Errorf("hash failed: %s", err)
		}

		rel, err := filepath.Rel(rootPath, path)
		if err != nil {
			return fmt.Errorf("rel path failed: %s", err)
		}

		hashes[rel] = fmt.Sprintf("%x", h.Sum(nil))

		return nil
	}
	if err := filepath.Walk(rootPath, walkFn); err != nil {
		return nil, fmt.Errorf("hash walk failed: %s", err)
	}

	return hashes, nil
}
