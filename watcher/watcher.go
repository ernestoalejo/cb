package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/howeyc/fsnotify"
)

var (
	watcher  *fsnotify.Watcher
	watched  = map[string]bool{}
	unblock  = make(chan bool)
	errorCh  = make(chan error)
	modified = map[string]bool{}
)

func FolderSync(path string) error {
	var err error
	if watcher, err = fsnotify.NewWatcher(); err != nil {
		return errors.New(err)
	}
	defer watcher.Close()

	go watchEvents()

	if err := filepath.Walk(path, walkFn); err != nil {
		return errors.New(err)
	}

	log.Println("Watching...")

	select {
	case <-unblock:
	case err := <-errorCh:
		return err
	}

	return nil
}

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return errors.New(err)
	}
	if !info.IsDir() {
		return nil
	}

	if err := watcher.Watch(path); err != nil {
		return errors.New(err)
	}
	watched[path] = true

	if *config.Verbose {
		log.Printf("watching `%s`\n", path)
	}

	return nil
}

func watchEvents() {
	for {
		select {
		case ev := <-watcher.Event:
			if err := processEvent(ev); err != nil {
				errorCh <- err
				return
			}
		case err := <-watcher.Error:
			errorCh <- errors.New(err)
			return
		}
	}
}

func processEvent(ev *fsnotify.FileEvent) error {
	log.Println(ev)
	info, err := os.Stat(ev.Name)
	if err != nil && !os.IsNotExist(err) {
		return errors.New(err)
	}

	// If it's a folder we were watching, remove it
	if err != nil && watched[ev.Name] {
		delete(watched, ev.Name)
		if err := watcher.RemoveWatch(ev.Name); err != nil {
			return errors.New(err)
		}
		if *config.Verbose {
			log.Printf("remove watcher `%s`\n", ev.Name)
		}
		return nil
	}

	// If it's a new folder, follow it too
	if err == nil && info.IsDir() && ev.IsCreate() {
		if err := watcher.Watch(ev.Name); err != nil {
			return errors.New(err)
		}
		watched[ev.Name] = true

		if *config.Verbose {
			log.Printf("watching `%s`\n", ev.Name)
		}
		return nil
	}

	// Otherwise add it to the list of modified files
	if *config.Verbose {
		log.Printf("modified `%s`\n", ev.Name)
	}
	modified[ev.Name] = true

	return nil
}
