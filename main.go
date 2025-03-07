package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

type Command struct {
	Cmd         string
	Name        string
	Description string
}

const commandsFile = ".commands.aqc"

func main() {
	if len(os.Args) < 2 {
		interactiveMode()
		return
	}

	switch os.Args[1] {
	case "add":
		addSubcommand()
	case "list":
		listSubcommand()
	case "help", "--help", "-h":
		printHelp()
	default:
		interactiveMode()
	}
}

func printHelp() {
	fmt.Println(ColorCyan + "AQC - Quick Command Tool" + ColorReset)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  aqc                     Launch interactive mode to select and run a command")
	fmt.Println("  aqc add --cmd=\"<command>\" --name=\"<name>\" --desc=\"<description>\"")
	fmt.Println("                          Add a new command to the command file")
	fmt.Println("  aqc list                List available commands")
	fmt.Println("  aqc help                Show this help message")
}

func addSubcommand() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	cmdPtr := addCmd.String("cmd", "", "The command to run")
	namePtr := addCmd.String("name", "", "The name of the command")
	descPtr := addCmd.String("desc", "", "A short description of the command")
	addCmd.Parse(os.Args[2:])

	if *cmdPtr == "" || *namePtr == "" {
		fmt.Println(ColorRed + "Error: --cmd and --name are required fields." + ColorReset)
		addCmd.Usage()
		os.Exit(1)
	}

	newCommand := Command{
		Cmd:         *cmdPtr,
		Name:        *namePtr,
		Description: *descPtr,
	}

	err := appendCommand(newCommand)
	if err != nil {
		fmt.Printf(ColorRed+"Error adding command: %v\n"+ColorReset, err)
		os.Exit(1)
	}
	fmt.Println(ColorGreen + "Command added successfully!" + ColorReset)
}

func listSubcommand() {
	commands := loadCommands()
	printHeader()
	fmt.Println(ColorYellow + "\nAvailable Commands:" + ColorReset)
	for i, cmd := range commands {
		fmt.Printf("[%d] %s: %s\n", i+1, ColorGreen+cmd.Name+ColorReset, cmd.Description)
	}
}

func interactiveMode() {
	commands := loadCommands()
	if len(commands) == 0 {
		fmt.Println(ColorRed + "No commands found in the file." + ColorReset)
		os.Exit(1)
	}
	clearScreen()
	printHeader()
	fmt.Println(ColorYellow + "\nQuick Command Menu:" + ColorReset)
	for i, cmd := range commands {
		fmt.Printf("[%d] %s: %s\n", i+1, ColorGreen+cmd.Name+ColorReset, cmd.Description)
	}

	fmt.Print("\nSelect a command by entering its number: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(commands) {
		fmt.Println(ColorRed + "Invalid selection. Exiting." + ColorReset)
		os.Exit(1)
	}

	selected := commands[index-1]
	// Show selected command details on top
	clearScreen()
	printHeader()
	fmt.Println(ColorYellow + "\nSelected Command:" + ColorReset)
	fmt.Printf("%s\n\n%s\n\n", ColorPurple+selected.Name+ColorReset, selected.Description)
	fmt.Println(ColorCyan + "Command: " + ColorReset + selected.Cmd)
	fmt.Println("\nExecuting...\n")
	runCommand(selected.Cmd)
}

func loadCommands() []Command {
	if _, err := os.Stat(commandsFile); os.IsNotExist(err) {
		fmt.Printf(ColorRed+"Error: %s not found in the current directory.\n"+ColorReset, commandsFile)
		os.Exit(1)
	}
	data, err := ioutil.ReadFile(commandsFile)
	if err != nil {
		fmt.Printf(ColorRed+"Error reading file: %v\n"+ColorReset, err)
		os.Exit(1)
	}
	blocks := parseBlocks(string(data))
	return parseCommands(blocks)
}

// clearScreen uses ANSI escape codes to clear the terminal screen.
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// printHeader prints a colorful header.
func printHeader() {
	fmt.Println(ColorCyan + "============================================" + ColorReset)
	fmt.Println(ColorGreen + "           AQC - Quick Command              " + ColorReset)
	fmt.Println(ColorCyan + "============================================" + ColorReset)
}

// parseBlocks splits the file content into separate command blocks.
// Blocks are separated by a line that contains exactly "---".
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
// Each block is expected to have at least two lines:
// the first is the actual command,
// the second starts with a hyphen (-) followed by the command name and description.
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
		// Split the info into a name and a description by the first colon.
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

// runCommand executes the provided shell command using sh -c.
func runCommand(command string) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf(ColorRed+"Error executing command: %v\n"+ColorReset, err)
	}
}

// appendCommand appends a new command block to the commands file.
func appendCommand(c Command) error {
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
