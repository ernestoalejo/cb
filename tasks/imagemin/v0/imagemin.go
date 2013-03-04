package v0

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("imagemin", 0, imagemin)
}

// Compress & optimize images. It does not run if the folder images does not
// exists inside the temp directory.
func imagemin(c config.Config, q *registry.Queue) error {
	root := filepath.Join("client", "temp", "images")
	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.New(err)
	}

	if err := filepath.Walk(root, walkFn); err != nil {
		return errors.New(err)
	}
	return nil
}

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return errors.New(err)
	}
	if info.IsDir() {
		return nil
	}

	base := filepath.Join("client", "temp", "images")
	dest, err := filepath.Rel(base, path)
	if err != nil {
		return errors.New(err)
	}
	dest = filepath.Join("client", "temp", "images", dest)

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return errors.New(err)
	}

	switch filepath.Ext(path) {
	case ".jpg":
		fallthrough
	case ".jpeg":
		if err := jpegtran(path, dest); err != nil {
			return err
		}

	case ".png":
		if err := optipng(path, dest); err != nil {
			return err
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
	if err == utils.ErrExec {
		fmt.Println(output)
		return errors.Format("tool error")
	} else if err != nil {
		return err
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
	if err == utils.ErrExec {
		fmt.Println(output)
		return errors.Format("tool error")
	} else if err != nil {
		return err
	}

	if err := os.Remove(dest + ".bak"); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.New(err)
	}

	return nil
}
