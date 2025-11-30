# AQC - Aman's Quick Command Tool

[![Test](https://github.com/amantham20/AQC/actions/workflows/test.yml/badge.svg)](https://github.com/amantham20/AQC/actions/workflows/test.yml)
[![Release](https://github.com/amantham20/AQC/actions/workflows/release.yml/badge.svg)](https://github.com/amantham20/AQC/actions/workflows/release.yml)

A fast, interactive command-line tool for managing and executing frequently used shell commands. Save your commonly used commands once, then quickly access and run them with a beautiful TUI (Terminal User Interface).

```
            __
           / _)
    .-^^^-/ /
 __/       /
<__.|_|-|_|
```

## ✨ Features

- **Interactive Menu**: Navigate through your saved commands with arrow keys
- **Quick Selection**: Press number keys (1-9) to instantly select and execute commands
- **Scrollable Interface**: Handle large command lists with automatic scrolling
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Colorful TUI**: Beautiful terminal interface with syntax highlighting
- **Simple File Format**: Commands stored in a human-readable `.commands.aqc` file

## 🚀 Installation

### Quick Install (Recommended)

#### Linux (AMD64)
```bash
curl -L https://github.com/amantham20/AQC/releases/latest/download/aqc-linux-amd64 -o aqc
chmod +x aqc
sudo mv aqc /usr/local/bin/
```

#### Linux (ARM64 - Raspberry Pi, AWS Graviton, etc.)
```bash
curl -L https://github.com/amantham20/AQC/releases/latest/download/aqc-linux-arm64 -o aqc
chmod +x aqc
sudo mv aqc /usr/local/bin/
```

#### macOS (Apple Silicon - M1/M2/M3)
```bash
curl -L https://github.com/amantham20/AQC/releases/latest/download/aqc-darwin-arm64 -o aqc
chmod +x aqc
sudo mv aqc /usr/local/bin/
```

#### macOS (Intel)
```bash
curl -L https://github.com/amantham20/AQC/releases/latest/download/aqc-darwin-amd64 -o aqc
chmod +x aqc
sudo mv aqc /usr/local/bin/
```

#### Windows (PowerShell)
```powershell
Invoke-WebRequest -Uri "https://github.com/amantham20/AQC/releases/latest/download/aqc-windows-amd64.exe" -OutFile "aqc.exe"
Move-Item aqc.exe C:\Windows\aqc.exe
```

### Build from Source

```bash
git clone https://github.com/amantham20/AQC.git
cd AQC
go build -o aqc .
sudo mv aqc /usr/local/bin/
```

## 📖 Usage

### Interactive Mode

Simply run `aqc` without arguments to launch the interactive menu:

```bash
aqc
```

This opens a beautiful TUI where you can:
- Use **↑/↓ arrow keys** to navigate
- Press **Enter** to execute the selected command
- Press **1-9** to quickly select and execute a command by number
- Press **q** or **Esc** to quit

### Add a New Command

```bash
aqc add --cmd="docker build -t myapp ." --name="Build Docker" --desc="Build the Docker image"
```

Parameters:
- `--cmd` (required): The shell command to save
- `--name` (required): A short name for the command
- `--desc` (optional): A description of what the command does

### List Commands

```bash
aqc list
```

### Show Help

```bash
aqc help
# or
aqc --help
aqc -h
```

### Show Version

```bash
aqc version
# or
aqc --version
aqc -v
```

## 📁 Command File Format

Commands are stored in `.commands.aqc` in your current directory. The format is simple and human-readable:

```
docker build -t name
- Docker Build: Build a Docker image tagged as "name"
---
pytest
- Run Tests: Execute the test suite using pytest
---
ls -la
- List All: List all files with details
---
```

Each command block contains:
1. The shell command (first line)
2. A hyphen followed by the name and description: `- Name: Description`
3. A separator: `---`

You can manually edit this file if needed!

## 🎯 Examples

### Setting Up a Project

```bash
# Create a .commands.aqc file with your common commands
aqc add --cmd="npm install" --name="Install" --desc="Install dependencies"
aqc add --cmd="npm run dev" --name="Dev Server" --desc="Start development server"
aqc add --cmd="npm run build" --name="Build" --desc="Build for production"
aqc add --cmd="npm test" --name="Test" --desc="Run test suite"
aqc add --cmd="git status" --name="Git Status" --desc="Check git status"
```

### Quick Workflow

```bash
# Launch interactive mode
aqc

# Press 1 to install dependencies
# Press 2 to start dev server
# Or use arrow keys to navigate
```

## ⌨️ Keyboard Shortcuts

| Key | Action |
|-----|--------|
| ↑ / ↓ | Navigate up/down |
| Enter | Execute selected command |
| 1-9 | Quick select and execute command |
| q | Quit |
| Esc | Quit |
| Ctrl+C | Quit |

## 🔧 Development

### Prerequisites

- Go 1.21 or later

### Building

```bash
# Build for current platform
go build -o aqc .

# Build for all platforms
./build.sh
```

### Running Tests

```bash
go test -v ./...
```

### Project Structure

```
AQC/
├── main.go           # Entry point and CLI handling
├── commands.go       # Command file parsing and management
├── interactive.go    # Interactive TUI menu
├── utils.go          # Utility functions and colors
├── add.go            # Add command subcommand
├── build.sh          # Cross-platform build script
├── *_test.go         # Test files
└── .commands.aqc     # Your saved commands (created on first use)
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 👤 Author

**Aman Dhruva Thamminana**

- GitHub: [@amantham20](https://github.com/amantham20)
- Email: thammina@msu.edu

## 🙏 Feedback

If you have any feedback, please reach out at thammina@msu.edu or open an issue on GitHub!

---

Made with ❤️ by Aman Dhruva Thamminana
