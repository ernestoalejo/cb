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

			var err error
			ans, err = buf.ReadString('\n')
			if err != nil {
				panic(err)
			}
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
