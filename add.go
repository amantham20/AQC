package main

import (
	"flag"
	"fmt"
	"os"
)

// AddSubcommand handles the "add" subcommand to append a new command to the file.
func AddSubcommand() {
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

	if err := AppendCommand(newCommand); err != nil {
		fmt.Printf("%sError adding command: %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}
	fmt.Println(ColorGreen + "Command added successfully!" + ColorReset)
}
