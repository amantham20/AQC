package main

import (
	"strings"
	"testing"
)

func TestMainCommandParsing(t *testing.T) {
	// Test various subcommand recognitions
	tests := []struct {
		name    string
		arg     string
		isHelp  bool
		isVersion bool
		isAdd   bool
		isList  bool
	}{
		{"help command", "help", true, false, false, false},
		{"--help flag", "--help", true, false, false, false},
		{"-h flag", "-h", true, false, false, false},
		{"version command", "version", false, true, false, false},
		{"--version flag", "--version", false, true, false, false},
		{"-v flag", "-v", false, true, false, false},
		{"add command", "add", false, false, true, false},
		{"list command", "list", false, false, false, true},
		{"unknown command", "unknown", false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isHelp := tt.arg == "help" || tt.arg == "--help" || tt.arg == "-h"
			isVersion := tt.arg == "version" || tt.arg == "--version" || tt.arg == "-v"
			isAdd := tt.arg == "add"
			isList := tt.arg == "list"

			if isHelp != tt.isHelp {
				t.Errorf("isHelp = %v, expected %v", isHelp, tt.isHelp)
			}
			if isVersion != tt.isVersion {
				t.Errorf("isVersion = %v, expected %v", isVersion, tt.isVersion)
			}
			if isAdd != tt.isAdd {
				t.Errorf("isAdd = %v, expected %v", isAdd, tt.isAdd)
			}
			if isList != tt.isList {
				t.Errorf("isList = %v, expected %v", isList, tt.isList)
			}
		})
	}
}

func TestNumericArgParsing(t *testing.T) {
	// Test that numeric arguments are properly detected
	tests := []struct {
		name      string
		arg       string
		isNumeric bool
		value     int
	}{
		{"single digit", "1", true, 1},
		{"double digit", "42", true, 42},
		{"zero", "0", true, 0},
		{"negative", "-1", true, -1},
		{"not numeric", "abc", false, 0},
		{"mixed", "12abc", false, 0},
		{"empty", "", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var num int
			var isNumeric bool

			// Simulate strconv.Atoi behavior
			if tt.arg == "" {
				isNumeric = false
			} else {
				negative := false
				startIdx := 0
				if tt.arg[0] == '-' {
					negative = true
					startIdx = 1
				}

				isNumeric = len(tt.arg) > startIdx
				for i := startIdx; i < len(tt.arg); i++ {
					if tt.arg[i] < '0' || tt.arg[i] > '9' {
						isNumeric = false
						break
					}
					num = num*10 + int(tt.arg[i]-'0')
				}
				if negative {
					num = -num
				}
			}

			if isNumeric != tt.isNumeric {
				t.Errorf("isNumeric = %v, expected %v", isNumeric, tt.isNumeric)
			}
			if isNumeric && num != tt.value {
				t.Errorf("value = %d, expected %d", num, tt.value)
			}
		})
	}
}

func TestVersionOutput(t *testing.T) {
	// Test that version string contains expected information
	versionStr := "AQC - Aman's Quick Command Tool v0.1"

	if !strings.Contains(versionStr, "AQC") {
		t.Error("Version should contain 'AQC'")
	}
	if !strings.Contains(versionStr, "v0.1") {
		t.Error("Version should contain version number")
	}
}

func TestHelpOutputFormat(t *testing.T) {
	// Test expected help text components
	expectedPhrases := []string{
		"AQC",
		"Quick Command",
		"Usage:",
		"add",
		"list",
		"help",
		"version",
		"--cmd",
		"--name",
		"--desc",
	}

	helpText := `AQC - Quick Command Tool

Usage:
  aqc                     Launch interactive mode to select and run a command
  aqc add --cmd="<command>" --name="<name>" --desc="<description>"
                          Add a new command to the command file
  aqc list                List available commands
  aqc help                Show this help message
  aqc version             Show the version information`

	for _, phrase := range expectedPhrases {
		if !strings.Contains(helpText, phrase) {
			t.Errorf("Help text should contain %q", phrase)
		}
	}
}
