package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
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
		if strings.Contains(s, "\"") {
			re := regexp.MustCompile("\"(.*?)\"")
			args = re.FindAllString(s, -1)
			fmt.Println(args)
			for i := range args {
				args[i] = strings.ReplaceAll(args[i], "\\\"", "\"")
				args[i] = strings.Trim(args[i], "\"")
			}
		} else if strings.Contains(s, "'") {
			re := regexp.MustCompile("'(.*?)'")
			args = re.FindAllString(s, -1)
			for i := range args {
				args[i] = strings.Trim(args[i], "'")
			}
		} else {
			if strings.Contains(argstr, "\\") {
				re := regexp.MustCompile(`[^\\] +`)
				args = re.Split(argstr, -1)
				for i := range args {
					args[i] = strings.ReplaceAll(args[i], "\\", "")
				}
			} else {
				args = strings.Fields(argstr)
			}
		}

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

	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(commands, " "))
}

func typeCommand(commands []string) {
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
