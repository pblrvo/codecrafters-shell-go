package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var _ = fmt.Fprint

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	reader := bufio.NewReader(os.Stdin)

	for {
		message, err := reader.ReadString('\n')
		message = strings.Trim(message, "\n")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		commands := strings.Split(message, " ")

		switch commands[0] {

		case "exit":
			code, err := strconv.Atoi(commands[1])
			if err != nil {
				os.Exit(1)
			}
			os.Exit(code)
		case "echo":
			fmt.Fprintf(os.Stdout, "%s\n", strings.Join(commands[1:], " "))
			fmt.Fprint(os.Stdout, "$ ")

		default:
			fmt.Fprintf(os.Stdout, "%s: command not found\n", commands[0])
			fmt.Fprint(os.Stdout, "$ ")
		}
	}
}
