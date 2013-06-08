package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s%s%s\n", colors.RED, err, colors.RESET)
	}
}

func run() error {
	flag.Parse()
	log.SetFlags(log.Ltime)

	if *config.NoColors {
		colors.SetNoColors()
	}

	if *config.Help {
		usage()
		return nil
	}
	args := flag.Args()
	if len(args) == 0 {
		usage()
		return nil
	}

	c, err := config.Load()
	if err != nil {
		return fmt.Errorf("config loading failed: %s", err)
	}
	if c == nil && !isNoConfigTask(args[0]) {
		return fmt.Errorf("config file not found")
	}

	q := &registry.Queue{}
	for _, task := range args {
		q.AddTask(task)
	}
	if err := q.RunWithTimer(c); err != nil {
		return err
	}

	return nil
}

func usage() {
	fmt.Println("\n Usage: cb [target] [options...]")
	flag.PrintDefaults()
	registry.PrintTasks()
}

func isNoConfigTask(t string) bool {
	switch t {
	case "validator":
		return true

	case "init":
		return true
		
	case "init:client":
		return true
	}
	return false
}
