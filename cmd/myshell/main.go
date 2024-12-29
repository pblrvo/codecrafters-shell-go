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

func parseCommand(s string) (string, []string) {
	parts := strings.SplitN(s, " ", 2)
	command := parts[0]

	if len(parts) == 1 {
		return command, nil
	}

	argstr := parts[1]
	var args []string

	var current strings.Builder
	inDoubleQuotes := false
	inSingleQuotes := false
	escaped := false

	for i := 0; i < len(argstr); i++ {
		c := argstr[i]

		if escaped {
			if !inSingleQuotes {
				if inDoubleQuotes {
					// In double quotes, backslash only special for \, $, ", and newline
					if c == '\\' || c == '$' || c == '"' || c == '\n' {
						current.WriteByte(c)
					} else {
						// For other characters, preserve both \ and the character
						current.WriteByte('\\')
						current.WriteByte(c)
					}
				} else {
					// Outside quotes, just preserve the escaped character
					current.WriteByte(c)
				}
			} else {
				// In single quotes, preserve everything literally
				current.WriteByte('\\')
				current.WriteByte(c)
			}
			escaped = false
			continue
		}

		switch c {
		case '\\':
			if !inSingleQuotes {
				escaped = true
			} else {
				current.WriteByte(c)
			}
		case '"':
			if !inSingleQuotes {
				inDoubleQuotes = !inDoubleQuotes
			} else {
				current.WriteByte(c)
			}
		case '\'':
			if !inDoubleQuotes {
				inSingleQuotes = !inSingleQuotes
			} else {
				current.WriteByte(c)
			}
		case ' ':
			if !inDoubleQuotes && !inSingleQuotes {
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			} else {
				current.WriteByte(c)
			}
		default:
			current.WriteByte(c)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return command, args
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		// Wait for user input
		s, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		s = strings.Trim(s, "\r\n")
		if s == "" {
			continue
		}

		command, args := parseCommand(s)

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
					cmd := exec.Command(executable, args...)
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					if err := cmd.Run(); err != nil {
						fmt.Printf("Error executing file: %v\n", err)
					} else {
						executed = true
					}
					break
				}
			}

			if !executed {
				fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
			}
		}
	}
}

func cdCommand(commands []string) {
	if len(commands) == 0 {
		home := os.Getenv("HOME")
		if err := os.Chdir(home); err != nil {
			fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", home)
		}
		return
	}

	if commands[0] == "~" {
		commands[0] = os.Getenv("HOME")
	}

	if err := os.Chdir(commands[0]); err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", commands[0])
	}
}

func pwdCommand(commands []string) {
	if len(commands) > 0 {
		fmt.Fprint(os.Stdout, "pwd: too many arguments\n")
		return
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintln(os.Stdout, path)
}

func exitCommand(commands []string) {
	if len(commands) == 0 {
		os.Exit(0)
	}

	code, err := strconv.Atoi(commands[0])
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}

func echoCommand(commands []string) {
	fmt.Fprintln(os.Stdout, strings.Join(commands, " "))
}

func typeCommand(commands []string) {
	if len(commands) == 0 {
		return
	}

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
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stdout, "%s: not found\n", commands[0])
		}
	}
}
