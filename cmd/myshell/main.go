package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var _ = fmt.Fprint

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	reader := bufio.NewReader(os.Stdin)

	for {
		command, err := reader.ReadString('\n')
		command = strings.Trim(command, "\n")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		switch {

		case command == "exit 0":
			os.Exit(0)

		case strings.HasPrefix(command, "echo"):
			fmt.Fprint(os.Stdout, command[5:]+"\n")
			fmt.Fprint(os.Stdout, "$ ")

		default:
			fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
			fmt.Fprint(os.Stdout, "$ ")
		}
	}
}
