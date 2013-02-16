package utils

import (
	"bufio"
	"fmt"
	"os"
)

func Ask(q string) bool {
	q = q + " [y/N]: "

	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s", q)
		line, _, err := buf.ReadLine()
		if err != nil {
			panic(err)
		}

		ans := string(line)
		if ans == "Y" || ans == "y" {
			return true
		} else if ans == "n" || ans == "N" || ans == "" {
			return false
		}

		fmt.Println("answer y or n")
	}
	panic("should not reach here")
}
