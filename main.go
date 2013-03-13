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

	c, found, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("config loading failed: %s", err)
	}
	if !found && (len(args) != 1 || args[0] != "init") {
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
