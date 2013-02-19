package utils

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ernestokarim/cb/config"
)

func Ask(q string) bool {
	q = q + " [y/N]: "

	buf := bufio.NewReader(os.Stdin)
	for {
		var ans string
		if !*config.AlwaysY && !*config.AlwaysN {
			fmt.Printf("%s", q)

			line, _, err := buf.ReadLine()
			if err != nil {
				panic(err)
			}
			ans = string(line)
		}

		if ans == "Y" || ans == "y" || *config.AlwaysY {
			return true
		} else if ans == "n" || ans == "N" || ans == "" || *config.AlwaysN {
			// Redudant check for ans == "" && *config.AlwaysN, leave for cleaner code
			return false
		}

		fmt.Println("answer y or n")
	}
	panic("should not reach here")
}
