package utils

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var (
	alwaysY = flag.Bool("y", false, "answer yes to all overwrites")
	alwaysN = flag.Bool("n", false, "answer no to all overwrites")
)

func Ask(q string) bool {
	q = q + " [y/N]: "

	buf := bufio.NewReader(os.Stdin)
	for {
		var ans string
		if !*alwaysY && !*alwaysN {
			fmt.Printf("%s", q)

			line, _, err := buf.ReadLine()
			if err != nil {
				panic(err)
			}
			ans = string(line)
		}

		if ans == "Y" || ans == "y" || *alwaysY {
			return true
		} else if ans == "n" || ans == "N" || ans == "" || *alwaysN {
			// Redudant check for ans == "" && *alwaysN, leave for cleaner code
			return false
		}

		fmt.Println("answer y or n")
	}
	panic("should not reach here")
}
