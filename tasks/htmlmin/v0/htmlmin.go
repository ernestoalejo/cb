package v0

import (
	"fmt"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

const selfPkg = "github.com/ernestokarim/cb/vendor"

func init() {
	registry.NewTask("htmlmin", 0, htmlmin)
}

func htmlmin(c *config.Config, q *registry.Queue) error {
	size := c.CountDefault("htmlmin")
	for i := 0; i < size; i++ {
		source := c.GetRequired("htmlmin[%d].source", i)
		dest := c.GetRequired("htmlmin[%d].dest", i)
		if err := htmlcompressor(source, dest); err != nil {
			return fmt.Errorf("html compress failed: %s", err)
		}
	}

	return nil
}

func htmlcompressor(src, dest string) error {
	base := utils.PackagePath(selfPkg)
	jarFile := filepath.Join(base, "htmlcompressor-1.5.3.jar")

	args := []string{
		"-jar", jarFile,
		"--type", "html",
		"-o", dest,
		"-r", src,
	}
	output, err := utils.Exec("java", args)
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("compressor error: %s", err)
	}

	return nil
}
