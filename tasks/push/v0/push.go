package v0

import (
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"code.google.com/p/gopass"
	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

const selfPkg = "github.com/ernestokarim/cb/tasks/push/v0/scripts"

func init() {
	registry.NewUserTask("push", 0, push)
}

func push(c *config.Config, q *registry.Queue) error {
	scriptsPath := utils.PackagePath(selfPkg)
	host := c.GetRequired("push")

	// FTP User & password
	user := q.NextTask()
	if user == "" {
		return fmt.Errorf("ftp user required as the first argument")
	}
	q.RemoveNextTask()

	password, err := gopass.GetPass(fmt.Sprintf("Enter \"%s\" password: ", user))
	if err != nil {
		return fmt.Errorf("cannot read password: %s", err)
	}
	if password == "" {
		return fmt.Errorf("ftp password is required")
	}

	// Hash local files
	log.Printf("Hashing local files... ")
	localHashes, err := hashLocalFiles()
	if err != nil {
		return fmt.Errorf("hash local files failed: %s", err)
	}
	log.Printf("Hashing local files... %s[SUCCESS]%s\n", colors.Green, colors.Reset)

	// Hash remote files
	log.Printf("Hashing remote files... ")
	remoteHashes, err := retrieveRemoveHashes(scriptsPath, user, password, host)
	if err != nil {
		return fmt.Errorf("retrieve remote hashes failed: %s", err)
	}
	log.Printf("Hashing remote files... %s[SUCCESS]%s\n", colors.Green, colors.Reset)

	// Remove similar files
	if remoteHashes != nil {
		fmt.Println("here")
	}

	// Prepare FTP commands
	log.Printf("Preparing FTP commands... ")
	if err := prepareFTPCommands(); err != nil {
		return fmt.Errorf("prepare FTP commands failed: %s", err)
	}
	log.Printf("Preparing FTP commands... %s[SUCCESS]%s\n", colors.Green, colors.Reset)

	// Upload files
	log.Printf("Uploading files... ")
	if err := uploadFiles(scriptsPath, user, password, host); err != nil {
		return fmt.Errorf("uploading files failed: %s", err)
	}
	log.Printf("Uploading files... %s[SUCCESS]%s\n", colors.Green, colors.Reset)

	_ = localHashes

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

func retrieveRemoveHashes(scriptsPath, user, password, host string) (map[string]string, error) {
	args := []string{user, password, host}
	output, err := utils.Exec(filepath.Join(scriptsPath, "download-hashes.sh"), args)
	if err != nil {
		fmt.Println(output)
		return nil, fmt.Errorf("download hashes script failed: %s", err)
	}

	f, err := os.Open("temp/hashes")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("stat hashes failed: %s", err)
	}

	hashes := map[string]string{}
	if err := gob.NewDecoder(f).Decode(&hashes); err != nil {
		return nil, fmt.Errorf("cannot decode hashes: %s", err)
	}

	return hashes, nil
}

func uploadFiles(scriptsPath, user, password, host string) error {
	args := []string{user, password, host}
	if err := utils.ExecCopyOutput(filepath.Join(scriptsPath, "upload.sh"), args); err != nil {
		return fmt.Errorf("upload script failed: %s", err)
	}

	return nil
}

func prepareFTPCommands() error {
	rootPath := "../deploy"

	f, err := os.Create("temp/upload-commands")
	if err != nil {
		return fmt.Errorf("cannot create commands file: %s", err)
	}
	defer f.Close()

	curFolder := "."
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == rootPath {
			return nil
		}

		folder, err := filepath.Rel(rootPath, filepath.Dir(path))
		if err != nil {
			return fmt.Errorf("cannot rel path: %s", err)
		}
		basename := filepath.Base(path)
		basepath := filepath.Join("..", "deploy", path)

		for curFolder != folder {
			curFolder = filepath.Dir(curFolder)
			fmt.Fprintf(f, "cd ..;")
		}

		if info.IsDir() {
			curFolder = filepath.Join(curFolder, basename)
			fmt.Fprintf(f, "(mkdir \"%s\" && cd \"%s\") || cd \"%s\";\n",
				basename, basename, basename)
		} else {
			fmt.Fprintf(f, "echo \".......... uploading %s...\"\n", filepath.Join(folder, basename))
			fmt.Fprintf(f, "put \"%s\" -o \"%s\";\n", basepath, basename)
		}

		return nil
	}
	if err := filepath.Walk(rootPath, walkFn); err != nil {
		return fmt.Errorf("hash walk failed: %s", err)
	}

	return nil
}
