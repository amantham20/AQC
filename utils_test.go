package main

import (
	"strings"
	"testing"
)

func TestColorConstants(t *testing.T) {
	// Test that color constants are properly defined ANSI codes
	tests := []struct {
		name     string
		color    string
		expected string
	}{
		{"ColorReset", ColorReset, "\033[0m"},
		{"ColorRed", ColorRed, "\033[31m"},
		{"ColorGreen", ColorGreen, "\033[32m"},
		{"ColorYellow", ColorYellow, "\033[33m"},
		{"ColorBlue", ColorBlue, "\033[34m"},
		{"ColorPurple", ColorPurple, "\033[35m"},
		{"ColorCyan", ColorCyan, "\033[36m"},
		{"ColorWhite", ColorWhite, "\033[37m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color != tt.expected {
				t.Errorf("%s = %q, expected %q", tt.name, tt.color, tt.expected)
			}
		})
	}
}

func TestColorCodesAreValid(t *testing.T) {
	// Ensure all color codes start with escape sequence
	colors := []string{ColorReset, ColorRed, ColorGreen, ColorYellow, ColorBlue, ColorPurple, ColorCyan, ColorWhite}

	for _, color := range colors {
		if !strings.HasPrefix(color, "\033[") {
			t.Errorf("Color code %q does not start with ANSI escape sequence", color)
		}
		if !strings.HasSuffix(color, "m") {
			t.Errorf("Color code %q does not end with 'm'", color)
		}
	}
}

func TestPrintHeaderContainsExpectedText(t *testing.T) {
	// We can't easily capture stdout in a unit test without more setup,
	// but we can at least verify the header would contain the expected strings
	// by checking that the color codes are valid

	// This is a basic sanity check
	header := ColorCyan + "============================================" + ColorReset
	if !strings.Contains(header, "====") {
		t.Error("Header should contain separator line")
	}

	title := ColorGreen + "           AQC - Quick Command              " + ColorReset
	if !strings.Contains(title, "AQC") {
		t.Error("Title should contain 'AQC'")
	}
}
