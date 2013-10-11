package v0

import (
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

const selfPkg = "github.com/ernestokarim/cb/tasks/init/v0/templates"

func init() {
	registry.NewUserTask("init:*", 0, initTask)
}

func initTask(c *config.Config, q *registry.Queue) error {
	// Retrieve the current working directory
	cur, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd failed: %s", err)
	}

	// Go back one folder if we're inside the client one
	if filepath.Base(cur) == "client" {
		cur = filepath.Dir(cur)
		if pathErr := os.Chdir(cur); pathErr != nil {
			return fmt.Errorf("chdir to root folder failed: %s", err)
		}
	}

	// Prepare some paths and start copying from the root templates folder
	parts := strings.Split(q.CurTask, ":")
	base := utils.PackagePath(filepath.Join(selfPkg, parts[1]))
	appname := filepath.Base(cur)

	if err := copyFiles(c, appname, base, cur, cur); err != nil {
		return fmt.Errorf("copy files failed: %s", err)
	}

	// Post-init steps
	if err := postInit(); err != nil {
		return fmt.Errorf("post init failed: %s", err)
	}

	return nil
}

// Copy recursively all the files in src to the dest folder.
//   - appname: Name of the app extracted from the root folder
//   - src: Source folder path
//   - dest: Dest folder path
//   - root: Root folder path
func copyFiles(c *config.Config, appname, src, dest, root string) error {
	// Read the list of files of the source folder
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read folder failed (%s): %s", src, err)
	}

	for _, entry := range files {
		// Full paths to source & dest files
		fullsrc := filepath.Join(src, entry.Name())
		fulldest := filepath.Join(dest, entry.Name())

		// Check if source is a directory or a file
		if entry.IsDir() {
			info, err := os.Stat(fulldest)
			if err != nil {
				// Unknown error
				if !os.IsNotExist(err) {
					return fmt.Errorf("stat dest failed: %s", err)
				}

				// Create dest directory
				if *config.Verbose {
					log.Printf("create folder `%s`\n", dest)
				}
				if err := os.Mkdir(fulldest, entry.Mode()); err != nil {
					return fmt.Errorf("create folder failed: %s", err)
				}
			} else if !info.IsDir() {
				// Dest already present and not a folder
				return fmt.Errorf("dest is present and is not a folder: %s", fulldest)
			}

			// Copy recursively the folder files
			if err := copyFiles(c, appname, fullsrc, fulldest, root); err != nil {
				return fmt.Errorf("recursive copy failed: %s", err)
			}
		} else {
			// Copy only one file
			fulldest, err = copyFile(c, appname, fullsrc, fulldest, root)
			if err != nil {
				return fmt.Errorf("copy file failed: %s", err)
			}
		}
		if err := os.Chmod(fulldest, entry.Mode()); err != nil {
			return fmt.Errorf("change mode failed: %s", err)
		}
	}
	return nil
}

// Copy a file, using templates if needed, from srcPath to destPath.
func copyFile(c *config.Config, appname, srcPath, destPath, root string) (string, error) {
	// Use a template for the file if needed
	if strings.HasPrefix(filepath.Base(srcPath), "cbtmpl.") {
		srcName, err := copyFileTemplate(appname, srcPath)
		if err != nil {
			return destPath, fmt.Errorf("copy file template failed: %s", err)
		}

		srcPath = srcName
		destPath = filepath.Join(filepath.Dir(destPath), filepath.Base(destPath)[7:])
	}

	// Open source file
	src, err := os.Open(srcPath)
	if err != nil {
		return destPath, fmt.Errorf("open source file failed: %s", err)
	}
	defer src.Close()

	// Path of the file relative to the root
	relDest, err := filepath.Rel(root, destPath)
	if err != nil {
		return destPath, fmt.Errorf("cannot rel dest path: %s", err)
	}

	if _, err := os.Stat(destPath); err != nil {
		// Stat failed
		if !os.IsNotExist(err) {
			return destPath, fmt.Errorf("stat failed: %s", err)
		}

		// If it doesn't exists, but the config file is present, we're updating
		// the contents, ask for creation perms
		if c != nil {
			q := fmt.Sprintf("Do you want to create `%s`?", relDest)
			if !utils.Ask(q) {
				return destPath, nil
			}
		}
	} else {
		// If it exists, but they're equal, ignore the copy of this file
		if equal, err := compareFiles(srcPath, destPath); err != nil {
			return destPath, fmt.Errorf("compare files failed: %s", err)
		} else if equal {
			return destPath, nil
		}

		// Otherwise ask the user to overwrite the file
		q := fmt.Sprintf("Do you want to overwrite `%s`?", relDest)
		if !utils.Ask(q) {
			return destPath, nil
		}
	}

	// Open dest file
	dest, err := os.Create(destPath)
	if err != nil {
		return destPath, err
	}
	defer dest.Close()

	// Copy the file contents
	if *config.Verbose {
		log.Printf("copy file `%s`\n", relDest)
	}
	if _, err := io.Copy(dest, src); err != nil {
		return destPath, fmt.Errorf("copy file failed: %s", err)
	}

	return destPath, nil
}

func copyFileTemplate(appname, srcPath string) (string, error) {
	t, err := template.New(filepath.Base(srcPath)).Delims(`{{%`, `%}}`).ParseFiles(srcPath)
	if err != nil {
		return "", fmt.Errorf("parse template failed: %s", err)
	}

	f, err := ioutil.TempFile("", "cb-init:")
	if err != nil {
		return "", fmt.Errorf("cannot create temp file: %s", err)
	}
	defer f.Close()

	data := map[string]interface{}{
		"AppName": appname,
	}
	if err := t.Execute(f, data); err != nil {
		return "", fmt.Errorf("execute template failed: %s", err)
	}

	return f.Name(), nil
}

func compareFiles(srcPath, destPath string) (bool, error) {
	src, err := os.Open(srcPath)
	if err != nil {
		return false, fmt.Errorf("open source failed: %s", err)
	}
	defer src.Close()

	dest, err := os.Open(destPath)
	if err != nil {
		return false, fmt.Errorf("open dest failed: %s", err)
	}
	defer dest.Close()

	srcContents, err := ioutil.ReadAll(src)
	if err != nil {
		return false, fmt.Errorf("read source failed: %s", err)
	}
	destContents, err := ioutil.ReadAll(dest)
	if err != nil {
		return false, fmt.Errorf("read dest failed: %s", err)
	}

	contentsHash := fmt.Sprintf("%x", crc32.ChecksumIEEE(srcContents))
	srcHash := fmt.Sprintf("%x", crc32.ChecksumIEEE(destContents))

	return srcHash == contentsHash, nil
}

func postInit() error {
	// Test if the post-init file exists
	if _, err := os.Stat("post-init.sh"); err != nil {
		if os.IsNotExist(err) {
			if *config.Verbose {
				log.Println("post-init.sh file doesn't exist")
			}
			return nil
		}
		return fmt.Errorf("err failed")
	}

	if *config.Verbose {
		log.Println("running post-init.sh file")
	}

	// Run it
	if err := utils.ExecCopyOutput("bash", []string{"./post-init.sh"}); err != nil {
		return fmt.Errorf("post init exec failed: %s", err)
	}

	return nil
}
