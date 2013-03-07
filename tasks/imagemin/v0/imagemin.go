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
	registry.NewTask("imagemin", 0, imagemin)
}

// Compress & optimize images. It does not run if the folder images does not
// exists inside the temp directory.
func imagemin(c config.Config, q *registry.Queue) error {
	var from string
	if *config.AngularMode {
		from = "client"
	}
	root := filepath.Join(from, "temp", "images")
	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat images folder failed: %s", err)
	}

	if err := filepath.Walk(root, walkFn); err != nil {
		return fmt.Errorf("walk images folder failed: %s", err)
	}
	return nil
}

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("walk error: %s", err)
	}
	if info.IsDir() {
		return nil
	}

	var from string
	if *config.AngularMode {
		from = "client"
	}
	base := filepath.Join(from, "temp", "images")
	dest, err := filepath.Rel(base, path)
	if err != nil {
		return fmt.Errorf("rel failed: %s", err)
	}
	dest = filepath.Join(from, "temp", "images", dest)

	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create folder failed (%s): %s", dir, err)
	}

	switch filepath.Ext(path) {
	case ".jpg":
		fallthrough
	case ".jpeg":
		if err := jpegtran(path, dest); err != nil {
			return fmt.Errorf("jpeg optimization failed: %s", err)
		}

	case ".png":
		if err := optipng(path, dest); err != nil {
			return fmt.Errorf("png optimization failed: %s", err)
		}
	}

	return nil
}

func jpegtran(src, dest string) error {
	log.Printf("optimizing jpeg `%s`\n", src)
	args := []string{
		"-copy", "none",
		"-optimize", "-progressive",
		"-outfile", dest, src,
	}
	output, err := utils.Exec("jpegtran", args)
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("jpeg optimizer error: %s", err)
	}

	return nil
}

func optipng(src, dest string) error {
	log.Printf("optimizing png `%s`\n", src)
	args := []string{
		"-strip", "all", "-clobber",
		"-out", dest, src,
	}
	output, err := utils.Exec("optipng", args)
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("png optimizer error: %s", err)
	}

	if err := os.Remove(dest + ".bak"); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("remove backup file failed: %s", err)
	}

	return nil
}
