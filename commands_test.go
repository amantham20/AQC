package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseBlocks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single block",
			input:    "ls -la\n- List Files: List all files in directory",
			expected: []string{"ls -la\n- List Files: List all files in directory"},
		},
		{
			name:     "two blocks",
			input:    "ls -la\n- List Files: List all files\n---\npwd\n- Current Dir: Show current directory",
			expected: []string{"ls -la\n- List Files: List all files", "pwd\n- Current Dir: Show current directory"},
		},
		{
			name:     "blocks with empty lines",
			input:    "ls -la\n- List Files: List all files\n\n---\n\npwd\n- Current Dir: Show current directory",
			expected: []string{"ls -la\n- List Files: List all files", "pwd\n- Current Dir: Show current directory"},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
		{
			name:     "only separator",
			input:    "---",
			expected: nil,
		},
		{
			name:     "multiple separators",
			input:    "cmd1\n- Name1: Desc1\n---\n---\ncmd2\n- Name2: Desc2",
			expected: []string{"cmd1\n- Name1: Desc1", "cmd2\n- Name2: Desc2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseBlocks(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("parseBlocks() returned %d blocks, expected %d", len(result), len(tt.expected))
				return
			}
			for i, block := range result {
				if block != tt.expected[i] {
					t.Errorf("parseBlocks()[%d] = %q, expected %q", i, block, tt.expected[i])
				}
			}
		})
	}
}

func TestParseCommands(t *testing.T) {
	tests := []struct {
		name     string
		blocks   []string
		expected []Command
	}{
		{
			name:   "single command",
			blocks: []string{"ls -la\n- List Files: List all files in directory"},
			expected: []Command{
				{Cmd: "ls -la", Name: "List Files", Description: "List all files in directory"},
			},
		},
		{
			name:   "multiple commands",
			blocks: []string{"ls -la\n- List Files: List all files", "pwd\n- Current Dir: Show current directory"},
			expected: []Command{
				{Cmd: "ls -la", Name: "List Files", Description: "List all files"},
				{Cmd: "pwd", Name: "Current Dir", Description: "Show current directory"},
			},
		},
		{
			name:   "command without description",
			blocks: []string{"ls\n- List Files"},
			expected: []Command{
				{Cmd: "ls", Name: "List Files", Description: ""},
			},
		},
		{
			name:     "invalid block - missing hyphen",
			blocks:   []string{"ls -la\nList Files: Description"},
			expected: nil,
		},
		{
			name:     "invalid block - single line",
			blocks:   []string{"ls -la"},
			expected: nil,
		},
		{
			name:     "empty blocks",
			blocks:   []string{},
			expected: nil,
		},
		{
			name:   "command with colons in description",
			blocks: []string{"echo 'hello:world'\n- Echo Test: Prints hello:world to stdout"},
			expected: []Command{
				{Cmd: "echo 'hello:world'", Name: "Echo Test", Description: "Prints hello:world to stdout"},
			},
		},
		{
			name:   "command with extra whitespace",
			blocks: []string{"  ls -la  \n  - List Files  :  List all files  "},
			expected: []Command{
				{Cmd: "ls -la", Name: "List Files", Description: "List all files"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommands(tt.blocks)
			if len(result) != len(tt.expected) {
				t.Errorf("parseCommands() returned %d commands, expected %d", len(result), len(tt.expected))
				return
			}
			for i, cmd := range result {
				if cmd.Cmd != tt.expected[i].Cmd {
					t.Errorf("parseCommands()[%d].Cmd = %q, expected %q", i, cmd.Cmd, tt.expected[i].Cmd)
				}
				if cmd.Name != tt.expected[i].Name {
					t.Errorf("parseCommands()[%d].Name = %q, expected %q", i, cmd.Name, tt.expected[i].Name)
				}
				if cmd.Description != tt.expected[i].Description {
					t.Errorf("parseCommands()[%d].Description = %q, expected %q", i, cmd.Description, tt.expected[i].Description)
				}
			}
		})
	}
}

func TestAppendCommand(t *testing.T) {
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

	tests := []struct {
		name        string
		command     Command
		expectError bool
	}{
		{
			name: "append simple command",
			command: Command{
				Cmd:         "ls -la",
				Name:        "List Files",
				Description: "List all files",
			},
			expectError: false,
		},
		{
			name: "append command without description",
			command: Command{
				Cmd:         "pwd",
				Name:        "Current Dir",
				Description: "",
			},
			expectError: false,
		},
		{
			name: "append command with special characters",
			command: Command{
				Cmd:         "echo 'hello world' | grep hello",
				Name:        "Grep Test",
				Description: "Test with pipes and quotes",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AppendCommand(tt.command)
			if (err != nil) != tt.expectError {
				t.Errorf("AppendCommand() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// Verify the file was created and contains the command
			data, err := os.ReadFile(commandsFile)
			if err != nil {
				t.Errorf("Failed to read commands file: %v", err)
				return
			}

			content := string(data)
			if !contains(content, tt.command.Cmd) {
				t.Errorf("Commands file does not contain expected command: %s", tt.command.Cmd)
			}
			if !contains(content, tt.command.Name) {
				t.Errorf("Commands file does not contain expected name: %s", tt.command.Name)
			}
		})
	}
}

func TestParseBlocksAndCommands_Integration(t *testing.T) {
	// Test the full flow from raw file content to parsed commands
	fileContent := `ls -la
- List Files: List all files in the current directory
---
git status
- Git Status: Show the working tree status
---
docker ps
- Docker List: List running containers
`

	blocks := parseBlocks(fileContent)
	if len(blocks) != 3 {
		t.Fatalf("Expected 3 blocks, got %d", len(blocks))
	}

	commands := parseCommands(blocks)
	if len(commands) != 3 {
		t.Fatalf("Expected 3 commands, got %d", len(commands))
	}

	expectedCommands := []Command{
		{Cmd: "ls -la", Name: "List Files", Description: "List all files in the current directory"},
		{Cmd: "git status", Name: "Git Status", Description: "Show the working tree status"},
		{Cmd: "docker ps", Name: "Docker List", Description: "List running containers"},
	}

	for i, cmd := range commands {
		if cmd.Cmd != expectedCommands[i].Cmd {
			t.Errorf("Command[%d].Cmd = %q, expected %q", i, cmd.Cmd, expectedCommands[i].Cmd)
		}
		if cmd.Name != expectedCommands[i].Name {
			t.Errorf("Command[%d].Name = %q, expected %q", i, cmd.Name, expectedCommands[i].Name)
		}
		if cmd.Description != expectedCommands[i].Description {
			t.Errorf("Command[%d].Description = %q, expected %q", i, cmd.Description, expectedCommands[i].Description)
		}
	}
}

func TestLoadCommandsWithFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "aqc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test commands file
	testContent := `echo hello
- Hello World: Print hello to stdout
---
date
- Show Date: Display current date and time
`
	testFilePath := filepath.Join(tempDir, commandsFile)
	if err := os.WriteFile(testFilePath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Note: LoadCommands calls os.Exit on error, so we test the underlying functions instead
	data, err := os.ReadFile(commandsFile)
	if err != nil {
		t.Fatalf("Failed to read commands file: %v", err)
	}

	blocks := parseBlocks(string(data))
	commands := parseCommands(blocks)

	if len(commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(commands))
	}

	if commands[0].Cmd != "echo hello" {
		t.Errorf("First command = %q, expected %q", commands[0].Cmd, "echo hello")
	}
	if commands[1].Name != "Show Date" {
		t.Errorf("Second command name = %q, expected %q", commands[1].Name, "Show Date")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
