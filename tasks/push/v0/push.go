package v0

import (
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

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
	remoteHashes, err := retrieveRemoteHashes(scriptsPath, user, password, host)
	if err != nil {
		return fmt.Errorf("retrieve remote hashes failed: %s", err)
	}
	log.Printf("Hashing remote files... %s[SUCCESS]%s\n", colors.Green, colors.Reset)

	if err := saveLocalHashes(localHashes); err != nil {
		return fmt.Errorf("save local hashes failed: %s", err)
	}

	// Prepare FTP commands
	log.Printf("Preparing FTP commands... ")
	if err := prepareFTPCommands(localHashes, remoteHashes); err != nil {
		return fmt.Errorf("prepare FTP commands failed: %s", err)
	}
	log.Printf("Preparing FTP commands... %s[SUCCESS]%s\n", colors.Green, colors.Reset)

	// Upload files
	log.Printf("Uploading files... ")
	if err := uploadFiles(scriptsPath, user, password, host); err != nil {
		return fmt.Errorf("uploading files failed: %s", err)
	}
	log.Printf("Uploading files... %s[SUCCESS]%s\n", colors.Green, colors.Reset)

	return nil
}

func hashLocalFiles() (map[string]string, error) {
	hashes := map[string]string{}
	rootPath := "../deploy"

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		h := sha1.New()
		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("open file failed: %s", err)
			}
			defer f.Close()
			content, err := ioutil.ReadAll(f)
			if err != nil {
				return fmt.Errorf("read file failed: %s", err)
			}
			if _, err := h.Write(content); err != nil {
				return fmt.Errorf("hash failed: %s", err)
			}
		}
		if _, err := h.Write([]byte(fmt.Sprintf("%s", info.Mode()))); err != nil {
			return fmt.Errorf("hash perms failed: %s", err)
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

func retrieveRemoteHashes(scriptsPath, user, password, host string) (map[string]string, error) {
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

func prepareFTPCommands(localHashes, remoteHashes map[string]string) error {
	rootPath := "../deploy"
	changed := 0

	// Prepare commands file
	f, err := os.Create("temp/upload-commands")
	if err != nil {
		return fmt.Errorf("cannot create commands file: %s", err)
	}
	defer f.Close()

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == rootPath {
			return nil
		}

		// Prepare paths
		folder, err := filepath.Rel(rootPath, filepath.Dir(path))
		if err != nil {
			return fmt.Errorf("cannot rel path: %s", err)
		}
		basename := filepath.Base(path)
		basepath := filepath.Join("..", "deploy", path)
		dest := filepath.Join(folder, basename)

		// Ignore similar files
		if remoteHashes[dest] != "" && remoteHashes[dest] == localHashes[dest] {
			return nil
		}
		changed++

		if info.IsDir() {
			// Make dirs
			fmt.Fprintf(f, "mkdir \"%s\" || cd .;\n", dest)
		} else {
			// Upload files
			fmt.Fprintf(f, "echo \".......... upload %s\"\n", filepath.Join(folder, basename))
			fmt.Fprintf(f, "put \"%s\" -o \"%s/\";\n", basepath, folder)
		}

		// Change perms of the files
		fmt.Fprintf(f, "chmod %o \"%s\";\n", info.Mode().Perm(), dest)

		return nil
	}
	if err := filepath.Walk(rootPath, walkFn); err != nil {
		return fmt.Errorf("hash walk failed: %s", err)
	}

	mk := make([]string, len(remoteHashes))
	i := 0
	for k := range remoteHashes {
		mk[i] = k
		i++
	}
	sort.Sort(sort.Reverse(sort.StringSlice(mk)))

	for _, path := range mk {
		if localHashes[path] == "" {
			fmt.Fprintf(f, "echo \".......... remove \"%s\"\"\n", path)
			fmt.Fprintf(f, "rm \"%s\" || cd .\n", path)
			changed++
		}
	}

	if changed > 0 {
		fmt.Fprintf(f, "echo \".......... push hashes\"\n")
		fmt.Fprintf(f, "put temp/hashes -o push-hashes;\n")
	} else {
		fmt.Fprintf(f, "echo \">>>>>>>>>> no changes\"\n")
	}

	return nil
}

func saveLocalHashes(hashes map[string]string) error {
	f, err := os.Create("temp/hashes")
	if err != nil {
		return fmt.Errorf("create file failed: %s", err)
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(hashes); err != nil {
		return fmt.Errorf("go failed encoder failed: %s", err)
	}
	return nil
}
