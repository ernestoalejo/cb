package config

import (
	"flag"
)

var (
	Verbose     = flag.Bool("v", false, "verbose mode")
	AlwaysY     = flag.Bool("y", false, "answer yes to all overwrites")
	AlwaysN     = flag.Bool("n", false, "answer no to all overwrites")
	Compiled    = flag.Bool("compiled", false, "test the compiled version of the app")
	Help        = flag.Bool("help", false, "show this help message")
	AngularMode = flag.Bool("angular", false, "use angular tasks")
	ClosureMode = flag.Bool("closure", false, "use closure tasks")
	NoColors    = flag.Bool("no-color", false, "don't use colors in the output")
)
