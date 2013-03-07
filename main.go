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

	if *config.Help {
		usage()
		return nil
	}
	args := flag.Args()
	if len(args) == 0 {
		usage()
		return nil
	}

	c, err := config.LoadConfig()
	if err != nil {
		return err
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
