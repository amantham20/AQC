package main

import (
	// "flag"
	"fmt"
	"os"
	"strconv"
)

// Version is set at build time via ldflags
var Version = "dev"

func main() {
	// If no subcommand is provided, use interactive mode.
	if len(os.Args) < 2 {
		InteractiveModeWithDefault()
		return
	}

	// Switch based on the provided subcommand.
	// Check if first argument is a number
	if num, err := strconv.Atoi(os.Args[1]); err == nil {
		InteractiveMode(num)
		return
	}

	switch os.Args[1] {
	case "add":
		AddSubcommand()
	case "list":
		ListSubcommand()
	case "help", "--help", "-h":
		PrintHelp()
	case "version", "--version", "-v":
		version()
	default:
		// Fallback to interactive mode for unknown arguments.
		fmt.Println("Unknown subcommand.")
		PrintHelp()

	}
}

func version() {
	fmt.Printf("AQC - Aman's Quick Command Tool %s\n", Version)

	out := `
            __
           / _)
    .-^^^-/ /
 __/       /
<__.|_|-|_|
`
	fmt.Println(out)
	fmt.Println("Developed by Aman Dhruva Thamminana")
	fmt.Println("Help me with feedback at thammina@msu.edu or contribute at https://github.com/amantham20/AQC")

}

func PrintHelp() {
	fmt.Println(ColorCyan + "AQC - Quick Command Tool" + ColorReset)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  aqc                     Launch interactive mode to select and run a command")
	fmt.Println("  aqc add --cmd=\"<command>\" --name=\"<name>\" --desc=\"<description>\"")
	fmt.Println("                          Add a new command to the command file")
	fmt.Println("  aqc list                List available commands")
	fmt.Println("  aqc help                Show this help message")
	fmt.Println("  aqc version             Show the version information")

}
