package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const commandsFile = ".commands.aqc"

// Command holds the shell command, its display name, and a short description.
type Command struct {
	Cmd         string
	Name        string
	Description string
}

// LoadCommands reads the commands file, parses its content, and returns a slice of Command.
func LoadCommands() []Command {
	if _, err := os.Stat(commandsFile); os.IsNotExist(err) {
		fmt.Printf("%sError: %s not found in the current directory.%s\n", ColorRed, commandsFile, ColorReset)
		os.Exit(1)
	}
	data, err := ioutil.ReadFile(commandsFile)
	if err != nil {
		fmt.Printf("%sError reading file: %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}
	blocks := parseBlocks(string(data))
	return parseCommands(blocks)
}

// parseBlocks splits the file content into separate command blocks.
// Blocks are separated by a line containing exactly "---".
func parseBlocks(data string) []string {
	var blocks []string
	lines := strings.Split(data, "\n")
	var currentBlock []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if len(currentBlock) > 0 {
				blocks = append(blocks, strings.Join(currentBlock, "\n"))
				currentBlock = []string{}
			}
		} else if trimmed != "" {
			currentBlock = append(currentBlock, line)
		}
	}
	if len(currentBlock) > 0 {
		blocks = append(blocks, strings.Join(currentBlock, "\n"))
	}
	return blocks
}

// parseCommands converts each block into a Command struct.
// Each block must have at least two lines: the first is the command,
// the second starts with a hyphen and contains the name and description.
func parseCommands(blocks []string) []Command {
	var commands []Command
	for _, block := range blocks {
		lines := strings.Split(block, "\n")
		if len(lines) < 2 {
			continue
		}
		cmdText := strings.TrimSpace(lines[0])
		secondLine := strings.TrimSpace(lines[1])
		if !strings.HasPrefix(secondLine, "-") {
			continue
		}
		// Remove the hyphen and any leading spaces.
		info := strings.TrimSpace(secondLine[1:])
		// Split the info into a name and description by the first colon.
		parts := strings.SplitN(info, ":", 2)
		name := strings.TrimSpace(parts[0])
		description := ""
		if len(parts) > 1 {
			description = strings.TrimSpace(parts[1])
		}
		commands = append(commands, Command{
			Cmd:         cmdText,
			Name:        name,
			Description: description,
		})
	}
	return commands
}

// RunCommand executes the provided shell command using sh -c.
func RunCommand(command string) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("%sError executing command: %v%s\n", ColorRed, err, ColorReset)
	}
	fmt.Printf("Here is me pringint something")
}

// AppendCommand appends a new command block to the commands file.
func AppendCommand(c Command) error {
	block := fmt.Sprintf("%s\n- %s: %s\n---\n", c.Cmd, c.Name, c.Description)
	f, err := os.OpenFile(commandsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(block); err != nil {
		return err
	}
	return nil
}
