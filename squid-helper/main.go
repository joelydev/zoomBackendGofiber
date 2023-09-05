package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner((os.Stdin))
	scanner.Scan()
	params := scanner.Text()

	if params != "" {
		fmt.Print("OK")
	}

	fmt.Print("ERR")
}
