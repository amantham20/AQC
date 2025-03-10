package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

// At the top of interactive.go, add a global variable for the debug log file.
var debugFile *os.File

// In InteractiveMode, open the debug log file.
func InteractiveMode() {
	// Open (or create) the debug log file.
	var err error
	debugFile, err = os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open debug log file: %v\n", err)
	}
	defer func() {
		if debugFile != nil {
			debugFile.Close()
		}
	}()

	commands := LoadCommands()
	if len(commands) == 0 {
		fmt.Println(ColorRed + "No commands found in the file." + ColorReset)
		os.Exit(1)
	}

	// Enter alternate screen mode so the application takes over the terminal.
	EnterAlternateScreen()
	defer ExitAlternateScreen()

	// Initialize terminal for raw mode to capture individual keystrokes
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(debugFile, "Error setting up terminal in raw mode: %v\n", err)
		fmt.Printf("%sError setting up terminal: %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Display the menu with scrolling
	selectedIndex := displayScrollableMenu(commands)

	// Exit if no valid selection
	if selectedIndex < 0 || selectedIndex >= len(commands) {
		return
	}

	// Get selected command
	selected := commands[selectedIndex]

	// Clear the screen to display the command details
	// ClearScreen()
	// PrintHeader()

	printLine(ColorCyan + "Shell Command:" + ColorReset)
	printLine(ColorWhite + selected.Cmd + ColorReset)
	printLine("")
	printLine(ColorGreen + "Command Name:" + ColorReset + " " + selected.Name)
	printLine(ColorYellow + "Description:" + ColorReset + " " + selected.Description)

	// Restore terminal to normal mode for command execution
	term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Println("\nPress ENTER to execute the command...")
	bufio.NewReader(os.Stdin).ReadString('\n')

	fmt.Println(ColorCyan + "\nExecuting...\n" + ColorReset)
	RunCommand(selected.Cmd)

}

// Update getTerminalHeight to log to the debug file.
func getTerminalHeight() int {
	_, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		if debugFile != nil {
			fmt.Fprintf(debugFile, "DEBUG: Failed to get terminal height, using default 24: %v\n", err)
		}
		return 24
	}
	if debugFile != nil {
		fmt.Fprintf(debugFile, "DEBUG: Terminal height = %d\n", height)
	}
	return height
}

// Update getTerminalWidth to log to the debug file.
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		if debugFile != nil {
			fmt.Fprintf(debugFile, "DEBUG: Failed to get terminal width, using default 80: %v\n", err)
		}
		return 80
	}
	if debugFile != nil {
		fmt.Fprintf(debugFile, "DEBUG: Terminal width = %d\n", width)
	}
	return width
}

// In displayScrollableMenu, log the dimensions.
func displayScrollableMenu(commands []Command) int {
	termHeight := getTerminalHeight()
	termWidth := getTerminalWidth()
	if debugFile != nil {
		fmt.Fprintf(debugFile, "DEBUG: Using terminal dimensions: height=%d, width=%d\n", termHeight, termWidth)
	}

	// Calculate available space for menu items (accounting for header and footer)
	headerLines := 4 // Header + blank line + title + blank line
	footerLines := 2 // Help text + input prompt
	maxVisibleItems := termHeight - headerLines - footerLines

	if maxVisibleItems < 1 {
		maxVisibleItems = 1
	}

	currentPos := 0   // Current cursor position
	scrollOffset := 0 // Current scroll offset

	// Main display loop
	for {
		ClearScreen()
		PrintHeader()
		printLine(ColorYellow + "Quick Command Menu:" + ColorReset)

		// Display visible commands
		displayEnd := scrollOffset + maxVisibleItems
		if displayEnd > len(commands) {
			displayEnd = len(commands)
		}

		// Show scroll indicator if needed
		if scrollOffset > 0 {
			printLine(ColorBlue + "  ▲ (more commands above)" + ColorReset)
		}

		// Display commands in the visible window
		for i := scrollOffset; i < displayEnd; i++ {
			prefix := "  "
			if i == currentPos {
				prefix = ColorCyan + "→ " + ColorReset // Highlight current selection
			}

			cmdName := commands[i].Name
			if len(cmdName) > 30 {
				cmdName = cmdName[:27] + "..."
			}

			desc := commands[i].Description
			// Calculate max description length and enforce a minimum length
			maxDescLen := termWidth - 40
			if maxDescLen < 10 {
				maxDescLen = 10
			}
			if len(desc) > maxDescLen && maxDescLen > 3 {
				desc = desc[:maxDescLen-3] + "..."
			}

			line := fmt.Sprintf("%s[%d] %s: %s", prefix, i+1, ColorGreen+cmdName+ColorReset, desc)
			printLine(line)
		}

		// Show scroll indicator if needed
		if displayEnd < len(commands) {
			fmt.Println(ColorBlue + "  ▼ (more commands below)" + ColorReset)
		}

		// Show help text
		line := ColorYellow + "Navigate: ↑/↓ arrows | Select: Enter | Quit: q/Esc" + ColorReset
		printLine(line)

		// Read a single key
		b := make([]byte, 3)
		n, err := os.Stdin.Read(b)
		if err != nil {
			break
		}

		// Process key input
		if n == 1 {
			switch b[0] {
			case 'q', 3, 27: // q, Ctrl+C, Esc
				return -1
			case 13: // Enter
				return currentPos
			}
		} else if n >= 3 {
			// Arrow keys send escape sequences
			if b[0] == 27 && b[1] == 91 {
				switch b[2] {
				case 65: // Up arrow
					if currentPos > 0 {
						currentPos--
						if currentPos < scrollOffset {
							scrollOffset = currentPos
						}
					}
				case 66: // Down arrow
					if currentPos < len(commands)-1 {
						currentPos++
						if currentPos >= scrollOffset+maxVisibleItems {
							scrollOffset = currentPos - maxVisibleItems + 1
						}
					}
				case 53: // Page Up
					currentPos -= maxVisibleItems
					if currentPos < 0 {
						currentPos = 0
					}
					scrollOffset = currentPos
				case 54: // Page Down
					currentPos += maxVisibleItems
					if currentPos >= len(commands) {
						currentPos = len(commands) - 1
					}
					if currentPos >= scrollOffset+maxVisibleItems {
						scrollOffset = currentPos - maxVisibleItems + 1
					}
				}
			}
		}
	}
	os.Stdout.Sync()
	return -1
}

func ListSubcommand() {}
