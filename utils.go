package main

import "fmt"

// ANSI color codes.
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

// ClearScreen uses ANSI escape codes to clear the current screen within the alternate buffer.
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// EnterAlternateScreen switches the terminal to the alternate screen buffer.
func EnterAlternateScreen() {
	fmt.Print("\033[?1049h")
}

// ExitAlternateScreen reverts the terminal back to the main screen buffer.
func ExitAlternateScreen() {
	fmt.Print("\033[?1049l")
}

func printLine(line string) {
	fmt.Print(line + "\r\n")
}

// PrintHeader prints a colorful header for the tool.
func PrintHeader() {
	printLine(ColorCyan + "============================================" + ColorReset)
	printLine(ColorGreen + "           AQC - Quick Command              " + ColorReset)
	printLine(ColorCyan + "============================================" + ColorReset)
}
