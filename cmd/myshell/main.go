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
var builtinCommands = map[string]int{"echo": 0, "type": 1, "exit": 2, "pwd": 3, "cd": 4}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		// Wait for user input
		s, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		s = strings.Trim(s, "\r\n")
		var args []string
		command, argstr, _ := strings.Cut(s, " ")

		// Parse arguments with backslash escape handling
		args = parseArguments(argstr)

		// Handle the command
		switch command {
		case "cd":
			cdCommand(args)
		case "pwd":
			pwdCommand(args)
		case "exit":
			exitCommand(args)
		case "echo":
			echoCommand(args)
		case "type":
			typeCommand(args)
		default:
			var executed bool
			env := os.Getenv("PATH")
			paths := strings.Split(env, ":")
			for _, path := range paths {
				executable := path + "/" + command
				if _, err := os.Stat(executable); err == nil {
					// Create and run the command
					cmd := exec.Command(executable, args...)
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					if err := cmd.Run(); err != nil {
						fmt.Printf("Error executing file: %v\n", err)
					} else {
						executed = true
					}
				}
			}

			if !executed {
				fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
			}
		}
	}
}

// parseArguments handles splitting arguments, respecting backslashes as escape characters
func parseArguments(input string) []string {
	var result []string
	var current string
	inQuote := false
	escape := false
	quoteChar := byte(0) // Track whether we're inside single or double quotes

	for i := 0; i < len(input); i++ {
		char := input[i]

		// Handle escape character
		if escape {
			current += string(char)
			escape = false
			continue
		}

		// Handle backslash as escape character
		if char == '\\' {
			escape = true
			continue
		}

		// Handle quotes (single or double)
		if char == '"' || char == '\'' {
			if inQuote && char == quoteChar {
				// Closing quote for the current quote type
				inQuote = false
				result = append(result, current)
				current = ""
			} else if !inQuote {
				// Opening quote for a string
				inQuote = true
				quoteChar = char
			} else {
				// Add the quote as part of the current argument
				current += string(char)
			}
			continue
		}

		// Handle spaces outside of quotes
		if char == ' ' && !inQuote {
			if current != "" {
				result = append(result, current)
				current = ""
			}
			continue
		}

		// Regular character: just append to the current argument
		current += string(char)
	}

	// Add the last argument if present
	if current != "" {
		result = append(result, current)
	}

	return result
}

func cdCommand(commands []string) {
	if strings.TrimSpace(commands[0]) == "~" {
		commands[0] = os.Getenv("HOME")
	}

	err := os.Chdir(commands[0])

	if err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", commands[0])
	}
}

func pwdCommand(commands []string) {
	if len(commands) > 1 {
		fmt.Fprint(os.Stdout, "pwd: too many arguments\n")
		fmt.Fprint(os.Stdout, "$ ")
		return
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Fprint(os.Stdout, path+"\n")
}

func exitCommand(commands []string) {
	code, err := strconv.Atoi(commands[0])
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}

func echoCommand(commands []string) {
	// Print the arguments joined by spaces, preserving escape sequences and quotes
	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(commands, " "))
}

func typeCommand(commands []string) {
	// Logic for the type command (not needed for your current task)
	if _, exists := builtinCommands[commands[0]]; exists {
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", commands[0])
	} else {
		var found bool
		env := os.Getenv("PATH")
		paths := strings.Split(env, ":")
		for _, path := range paths {
			exec := path + "/" + commands[0]
			if _, err := os.Stat(exec); err == nil {
				fmt.Fprintf(os.Stdout, "%v is %v\n", commands[0], exec)
				found = true
			}
		}
		if !found {
			fmt.Fprintf(os.Stdout, "%s: not found\n", commands[0])
		}
	}
}
