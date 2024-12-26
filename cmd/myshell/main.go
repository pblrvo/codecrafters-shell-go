package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
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

		case "pwd":
			path, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}
			fmt.Print(path + "\n")
			fmt.Fprint(os.Stdout, "$ ")
		case "exit":
			code, err := strconv.Atoi(commands[1])
			if err != nil {
				os.Exit(1)
			}
			os.Exit(code)
		case "echo":
			fmt.Fprintf(os.Stdout, "%s\n", strings.Join(commands[1:], " "))
			fmt.Fprint(os.Stdout, "$ ")

		case "type":
			switch commands[1] {
			case "echo", "type", "exit", "pwd":
				fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", commands[1])
			default:
				var found bool
				env := os.Getenv("PATH")
				paths := strings.Split(env, ":")
				for _, path := range paths {
					exec := path + "/" + commands[1]
					if _, err := os.Stat(exec); err == nil {
						fmt.Fprintf(os.Stdout, "%v is %v\n", commands[1], exec)
						found = true
					}
				}
				if !found {
					fmt.Fprintf(os.Stdout, "%s: not found\n", commands[1])
				}

			}
			fmt.Fprint(os.Stdout, "$ ")
		default:
			var executed bool
			env := os.Getenv("PATH")
			paths := strings.Split(env, ":")
			for _, path := range paths {
				executable := path + "/" + commands[0]
				args := commands[1:]
				if _, err := os.Stat(executable); err == nil {
					// Create the command
					cmd := exec.Command(executable, args...)

					// Set the command's standard input, output, and error streams
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					// Run the command
					if err := cmd.Run(); err != nil {
						fmt.Printf("Error executing file: %v\n", err)
					} else {
						fmt.Fprint(os.Stdout, "$ ")
						executed = true
					}
				}
			}

			if !executed {
				fmt.Fprintf(os.Stdout, "%s: command not found\n", commands[0])
				fmt.Fprint(os.Stdout, "$ ")
			}
		}
	}
}
