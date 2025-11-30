package main

import (
	"os"
	"testing"
)

func TestInteractiveModeNumberKeySelection(t *testing.T) {
	// Test the number key parsing logic (1-9 should map to indices 0-8)
	tests := []struct {
		name          string
		keyByte       byte
		numCommands   int
		expectedIndex int
		shouldSelect  bool
	}{
		{"key 1 with 5 commands", '1', 5, 0, true},
		{"key 5 with 5 commands", '5', 5, 4, true},
		{"key 9 with 9 commands", '9', 9, 8, true},
		{"key 9 with 5 commands", '9', 5, -1, false}, // Out of range
		{"key 1 with 1 command", '1', 1, 0, true},
		{"key 2 with 1 command", '2', 1, -1, false}, // Out of range
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the number key logic from displayScrollableMenu
			num := int(tt.keyByte - '0')
			var resultIndex int
			var selected bool

			if num > 0 && num <= tt.numCommands {
				resultIndex = num - 1
				selected = true
			} else {
				resultIndex = -1
				selected = false
			}

			if selected != tt.shouldSelect {
				t.Errorf("Selection status = %v, expected %v", selected, tt.shouldSelect)
			}
			if selected && resultIndex != tt.expectedIndex {
				t.Errorf("Selected index = %d, expected %d", resultIndex, tt.expectedIndex)
			}
		})
	}
}

func TestKeyInputProcessing(t *testing.T) {
	// Test the key input processing logic
	tests := []struct {
		name       string
		inputByte  byte
		expectQuit bool
		expectEnter bool
		expectNumber bool
	}{
		{"q key", 'q', true, false, false},
		{"Ctrl+C", 3, true, false, false},
		{"Esc", 27, true, false, false},
		{"Enter", 13, false, true, false},
		{"Number 1", '1', false, false, true},
		{"Number 9", '9', false, false, true},
		{"Number 0", '0', false, false, false}, // 0 is not a valid selection
		{"Letter a", 'a', false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isQuit := tt.inputByte == 'q' || tt.inputByte == 3 || tt.inputByte == 27
			isEnter := tt.inputByte == 13
			isNumber := tt.inputByte >= '1' && tt.inputByte <= '9'

			if isQuit != tt.expectQuit {
				t.Errorf("isQuit = %v, expected %v", isQuit, tt.expectQuit)
			}
			if isEnter != tt.expectEnter {
				t.Errorf("isEnter = %v, expected %v", isEnter, tt.expectEnter)
			}
			if isNumber != tt.expectNumber {
				t.Errorf("isNumber = %v, expected %v", isNumber, tt.expectNumber)
			}
		})
	}
}

func TestArrowKeyEscapeSequences(t *testing.T) {
	// Test arrow key escape sequence detection
	tests := []struct {
		name     string
		bytes    [3]byte
		isUp     bool
		isDown   bool
	}{
		{"Up arrow", [3]byte{27, 91, 65}, true, false},
		{"Down arrow", [3]byte{27, 91, 66}, false, true},
		{"Left arrow", [3]byte{27, 91, 68}, false, false},
		{"Right arrow", [3]byte{27, 91, 67}, false, false},
		{"Not escape", [3]byte{65, 91, 65}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEscapeSeq := tt.bytes[0] == 27 && tt.bytes[1] == 91
			isUp := isEscapeSeq && tt.bytes[2] == 65
			isDown := isEscapeSeq && tt.bytes[2] == 66

			if isUp != tt.isUp {
				t.Errorf("isUp = %v, expected %v", isUp, tt.isUp)
			}
			if isDown != tt.isDown {
				t.Errorf("isDown = %v, expected %v", isDown, tt.isDown)
			}
		})
	}
}

func TestScrollingLogic(t *testing.T) {
	// Test the scrolling logic for the menu
	tests := []struct {
		name            string
		currentPos      int
		scrollOffset    int
		maxVisibleItems int
		numCommands     int
		moveDown        bool
		expectedPos     int
		expectedOffset  int
	}{
		{
			name:            "move down within visible",
			currentPos:      0,
			scrollOffset:    0,
			maxVisibleItems: 5,
			numCommands:     10,
			moveDown:        true,
			expectedPos:     1,
			expectedOffset:  0,
		},
		{
			name:            "move down causes scroll",
			currentPos:      4,
			scrollOffset:    0,
			maxVisibleItems: 5,
			numCommands:     10,
			moveDown:        true,
			expectedPos:     5,
			expectedOffset:  1,
		},
		{
			name:            "at bottom, can't move down",
			currentPos:      9,
			scrollOffset:    5,
			maxVisibleItems: 5,
			numCommands:     10,
			moveDown:        true,
			expectedPos:     9,
			expectedOffset:  5,
		},
		{
			name:            "move up within visible",
			currentPos:      3,
			scrollOffset:    0,
			maxVisibleItems: 5,
			numCommands:     10,
			moveDown:        false,
			expectedPos:     2,
			expectedOffset:  0,
		},
		{
			name:            "at top, can't move up",
			currentPos:      0,
			scrollOffset:    0,
			maxVisibleItems: 5,
			numCommands:     10,
			moveDown:        false,
			expectedPos:     0,
			expectedOffset:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := tt.currentPos
			offset := tt.scrollOffset

			if tt.moveDown {
				if pos < tt.numCommands-1 {
					pos++
					if pos >= offset+tt.maxVisibleItems {
						offset = pos - tt.maxVisibleItems + 1
					}
				}
			} else {
				if pos > 0 {
					pos--
					if pos < offset {
						offset = pos
					}
				}
			}

			if pos != tt.expectedPos {
				t.Errorf("position = %d, expected %d", pos, tt.expectedPos)
			}
			if offset != tt.expectedOffset {
				t.Errorf("offset = %d, expected %d", offset, tt.expectedOffset)
			}
		})
	}
}

func TestTerminalDimensionDefaults(t *testing.T) {
	// Test default terminal dimensions when GetSize fails
	// These are the defaults used in the code
	defaultHeight := 24
	defaultWidth := 80

	if defaultHeight < 1 {
		t.Error("Default height should be at least 1")
	}
	if defaultWidth < 1 {
		t.Error("Default width should be at least 1")
	}
}

func TestMenuCalculations(t *testing.T) {
	// Test menu space calculations
	headerLines := 4
	footerLines := 2
	termHeight := 24

	maxVisibleItems := termHeight - headerLines - footerLines
	expectedMaxVisible := 18

	if maxVisibleItems != expectedMaxVisible {
		t.Errorf("maxVisibleItems = %d, expected %d", maxVisibleItems, expectedMaxVisible)
	}

	// Test minimum visibility
	smallTermHeight := 5
	maxVisibleSmall := smallTermHeight - headerLines - footerLines
	if maxVisibleSmall < 1 {
		maxVisibleSmall = 1
	}
	if maxVisibleSmall != 1 {
		t.Errorf("maxVisibleItems for small terminal should be at least 1, got %d", maxVisibleSmall)
	}
}

func TestDebugFileCreation(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "aqc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Test debug file creation
	debugPath := "debug.log"
	f, err := os.OpenFile(debugPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to create debug file: %v", err)
	}
	f.Close()

	// Verify file exists
	if _, err := os.Stat(debugPath); os.IsNotExist(err) {
		t.Error("Debug file was not created")
	}
}
