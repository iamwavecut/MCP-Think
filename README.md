# MCP-Think

MCP-Think is a Model Context Protocol (MCP) server that implements a "[Think Tool](https://www.anthropic.com/engineering/claude-think-tool)" for LLMs. This tool allows LLMs to record and retrieve their thinking processes during reasoning.

# YOLO
```bash
curl -fsSL https://raw.githubusercontent.com/iamwavecut/MCP-Think/main/install.sh | bash
```

## Features

- **Think Tool**: Record thoughts and reasoning steps
- **Get Thoughts**: Retrieve all previously recorded thoughts
- **Clear Thoughts**: Clear all recorded thoughts
- **Get Thought Stats**: Get statistics about recorded thoughts

## Installation

There are several ways to install and run MCP-Think:

**1. Pre-built Binaries (Recommended for Standalone Use)**

Ready-to-use binaries for Linux, Windows, and macOS (amd64 & arm64) are automatically built and attached to each [GitHub Release](https://github.com/iamwavecut/MCP-Think/releases). This is the easiest way to get started if you don't need to modify the code.

*   **macOS/Linux Auto-Install Script: (see YOLO)**
    *The script automatically detects your OS and architecture, downloads the appropriate binary, and guides you through installation.*

*   **Manual Installation (incl. Windows):**
    1.  Go to the [Releases page](https://github.com/iamwavecut/MCP-Think/releases).
    2.  Download the appropriate binary for your system (e.g., `think-tool-linux-amd64`, `think-tool-windows-amd64.exe`, `think-tool-darwin-arm64`).
    3.  (Optional) Rename it: `mv think-tool-linux-amd64 think-tool`
    4.  Make it executable (Linux/macOS): `chmod +x think-tool`
    5.  Run it: `./think-tool` (See Usage section)



**2. Using `go install` (Requires Go)**

This command compiles and installs the binary into your Go bin directory (`$GOPATH/bin` or `$HOME/go/bin`).

```bash
go install github.com/iamwavecut/MCP-Think@latest
```

*   **Note:** Ensure your Go bin directory is in your system's `PATH`. You might need to add `export PATH=$PATH:$(go env GOPATH)/bin` or `export PATH=$PATH:$HOME/go/bin` to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.).
*   Run the installed binary: `MCP-Think`

**3. Using `go run` (Quick Testing, Requires Go)**

This command compiles and runs the `main` package directly from the source code without installing a binary. It's useful for quick tests.

```bash
go run github.com/iamwavecut/MCP-Think@latest
```
*   This downloads the module and its dependencies temporarily if needed.

## Requirements

-   Go 1.24 or higher (if building from source or using `go install`/`go run`)

## Usage

## Running the Standalone Server**

If you installed via **Pre-built Binary** or **`go install`**:

```bash
# If using pre-built binary in current directory:
./think-tool

# If installed via 'go install' or the install script to /usr/local/bin:
think-tool
```

If you are using **`go run`**:

```bash
go run github.com/iamwavecut/MCP-Think@latest
```

The server will print `Starting Think Tool MCP Server with stdio transport...` and wait for MCP requests on stdin.