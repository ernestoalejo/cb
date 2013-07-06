package utils

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
)

// Ask the user a yes/no question by console.
func Ask(q string) bool {
	q = q + " [y/N]: "

	buf := bufio.NewReader(os.Stdin)
	for {
		var ans string
		if !*config.AlwaysY && !*config.AlwaysN {
			fmt.Fprintf(os.Stderr, "%s%s%s", colors.YELLOW, q, colors.RESET)

			var err error
			ans, err = buf.ReadString('\n')
			if err != nil {
				panic(err)
			}
			ans = ans[:len(ans)-1]
		}

		if ans == "Y" || ans == "y" || *config.AlwaysY {
			return true
		} else if ans == "n" || ans == "N" || ans == "" || *config.AlwaysN {
			// Redudant check for ans == "" && *config.AlwaysN, leave for cleaner code
			return false
		}

		fmt.Fprintf(os.Stderr, "%sanswer y or n%s\n", colors.RED, colors.RESET)
	}
	panic("should not reach here")
}
