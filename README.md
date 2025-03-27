# MCP-Think

MCP-Think is a Model Context Protocol (MCP) server that implements a "[Think Tool](https://www.anthropic.com/engineering/claude-think-tool)" for LLMs. This tool allows LLMs to record and retrieve their thinking processes during reasoning.
---

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

### 1. Pre-built Binaries (Recommended for Standalone Use)

Ready-to-use binaries for Linux, Windows, and macOS (amd64 & arm64) are automatically built and attached to each [GitHub Release](https://github.com/iamwavecut/MCP-Think/releases). This is the easiest way to get started if you don't need to modify the code.

*   #### macOS/Linux Auto-Install Script: (see YOLO)
    *The script automatically detects your OS and architecture, downloads the appropriate binary, and guides you through installation.*

*   #### Manual Installation (incl. Windows):
    1.  Go to the [Releases page](https://github.com/iamwavecut/MCP-Think/releases).
    2.  Download the appropriate binary for your system (e.g., `think-tool-linux-amd64`, `think-tool-windows-amd64.exe`, `think-tool-darwin-arm64`).
    3.  (Optional) Rename it: `mv think-tool-linux-amd64 think-tool`
    4.  Make it executable (Linux/macOS): `chmod +x think-tool`
    5.  Run it: `./think-tool` (See Usage section)



### 2. Using `go install` (Requires Go)

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

### Requirements

-   Go 1.24 or higher (if building from source or using `go install`/`go run`)

## Usage

### Running the Standalone Server

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

### Setting up in Cursor

To use MCP-Think with Cursor, follow these steps:

1. Install MCP-Think using one of the installation methods above
2. Create or update your Cursor MCP configuration file at `~/.cursor/mcp.json`:

```javascript
{
  "mcpServers": {
    "think-tool": {
      "command": "think-tool", // or absolute path, f.e.: /opt/bin/think-tool
      "transport": "stdio"
    }
  }
}
```
![Cursor MCP should be initialized](https://github.com/user-attachments/assets/addb8439-8259-4d3f-a055-773d9819468d)

3. Add the following rule to your Cursor rules:
<details>
<summary>Cursor settings > Rules > User rules</summary>
<pre>
## Using the think tool

Before taking any action or responding to the user after receiving tool results, use the think tool as a scratchpad to:
- List the specific rules that apply to the current request
- Check if all required information is collected
- Verify that the planned action complies with all policies
- Iterate over tool results for correctness 

Here are some examples of what to iterate over inside the think tool:
<think_tool_example_1>
User wants to cancel flight ABC123
- Need to verify: user ID, reservation ID, reason
- Check cancellation rules:
  * Is it within 24h of booking?
  * If not, check ticket class and insurance
- Verify no segments flown or are in the past
- Plan: collect missing info, verify rules, get confirmation
</think_tool_example_1>

<think_tool_example_2>
User wants to book 3 tickets to NYC with 2 checked bags each
- Need user ID to check:
  * Membership tier for baggage allowance
  * Which payments methods exist in profile
- Baggage calculation:
  * Economy class × 3 passengers
  * If regular member: 1 free bag each → 3 extra bags = $150
  * If silver member: 2 free bags each → 0 extra bags = $0
  * If gold member: 3 free bags each → 0 extra bags = $0
- Payment rules to verify:
  * Max 1 travel certificate, 1 credit card, 3 gift cards
  * All payment methods must be in profile
  * Travel certificate remainder goes to waste
- Plan:
1. Get user ID
2. Verify membership level for bag fees
3. Check which payment methods in profile and if their combination is allowed
4. Calculate total: ticket price + any bag fees
5. Get explicit confirmation for booking
</think_tool_example_2>
</pre>
</details>
3. Cursor can now use the Think Tool in your Cursor conversations with Claude 3.7 Sonnet
![image](https://github.com/user-attachments/assets/e90f61ab-0609-4bd7-961d-f64c49dd15c7)

