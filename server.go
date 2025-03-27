package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ThoughtEntry represents a single thought with timestamp
type ThoughtEntry struct {
	Timestamp string `json:"timestamp"`
	Thought   string `json:"thought"`
}

// ThinkToolServer implements the Think Tool functionality
type ThinkToolServer struct {
	server     *server.MCPServer
	thoughts   []ThoughtEntry
	serverName string
}

// NewThinkToolServer creates a new instance of ThinkToolServer
func NewThinkToolServer(serverName string) *ThinkToolServer {
	if serverName == "" {
		serverName = "think-tool"
	}

	s := &ThinkToolServer{
		thoughts:   []ThoughtEntry{},
		serverName: serverName,
	}

	// Create MCP server
	s.server = server.NewMCPServer(
		serverName,
		"1.0.0",
	)

	// Register tools
	s.registerTools()

	return s
}

// registerTools sets up all the available tools
func (s *ThinkToolServer) registerTools() {
	// Register the think tool
	thinkTool := mcp.NewTool("think",
		mcp.WithDescription("Use this tool to think about something. It will not obtain new information or change anything, but just append the thought to the log. Use it when complex reasoning or cache memory is needed."),
		mcp.WithString("thought",
			mcp.Required(),
			mcp.Description("A thought to think about. This can be structured reasoning, step-by-step analysis, policy verification, or any other mental process that helps with problem-solving."),
		),
	)
	s.server.AddTool(thinkTool, s.handleThink)

	// Register the get_thoughts tool
	getThoughtsTool := mcp.NewTool("get_thoughts",
		mcp.WithDescription("Retrieve all thoughts recorded in the current session. This tool helps review the thinking process that has occurred so far."),
	)
	s.server.AddTool(getThoughtsTool, s.handleGetThoughts)

	// Register the clear_thoughts tool
	clearThoughtsTool := mcp.NewTool("clear_thoughts",
		mcp.WithDescription("Clear all recorded thoughts from the current session. Use this to start fresh if the thinking process needs to be reset."),
	)
	s.server.AddTool(clearThoughtsTool, s.handleClearThoughts)

	// Register the get_thought_stats tool
	getStatsThoughtsTool := mcp.NewTool("get_thought_stats",
		mcp.WithDescription("Get statistics about the thoughts recorded in the current session."),
	)
	s.server.AddTool(getStatsThoughtsTool, s.handleGetThoughtStats)
}

// handleThink implements the think tool functionality
func (s *ThinkToolServer) handleThink(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	thought, ok := request.Params.Arguments["thought"].(string)
	if !ok {
		return nil, fmt.Errorf("thought must be a string")
	}

	// Record thought with timestamp
	timestamp := time.Now().Format(time.RFC3339)
	s.thoughts = append(s.thoughts, ThoughtEntry{
		Timestamp: timestamp,
		Thought:   thought,
	})

	// Generate confirmation message
	var message string
	if len(thought) > 50 {
		message = fmt.Sprintf("Thought recorded: %s...", thought[:50])
	} else {
		message = fmt.Sprintf("Thought recorded: %s", thought)
	}

	return mcp.NewToolResultText(message), nil
}

// handleGetThoughts implements the get_thoughts tool functionality
func (s *ThinkToolServer) handleGetThoughts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if len(s.thoughts) == 0 {
		return mcp.NewToolResultText("No thoughts have been recorded yet."), nil
	}

	var formattedThoughts []string
	for i, entry := range s.thoughts {
		formattedThoughts = append(formattedThoughts,
			fmt.Sprintf("Thought #%d (%s):\n%s\n", i+1, entry.Timestamp, entry.Thought))
	}

	result := ""
	for _, thought := range formattedThoughts {
		result += thought + "\n"
	}

	return mcp.NewToolResultText(result), nil
}

// handleClearThoughts implements the clear_thoughts tool functionality
func (s *ThinkToolServer) handleClearThoughts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	count := len(s.thoughts)
	s.thoughts = []ThoughtEntry{}

	return mcp.NewToolResultText(fmt.Sprintf("Cleared %d recorded thoughts.", count)), nil
}

// handleGetThoughtStats implements the get_thought_stats tool functionality
func (s *ThinkToolServer) handleGetThoughtStats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if len(s.thoughts) == 0 {
		return mcp.NewToolResultText("No thoughts have been recorded yet."), nil
	}

	totalThoughts := len(s.thoughts)

	// Calculate average length
	var totalLength int
	for _, entry := range s.thoughts {
		totalLength += len(entry.Thought)
	}
	avgLength := float64(totalLength) / float64(totalThoughts)

	// Find longest thought
	longestLength := 0
	longestIndex := -1
	for i, entry := range s.thoughts {
		if len(entry.Thought) > longestLength {
			longestLength = len(entry.Thought)
			longestIndex = i
		}
	}

	// Create stats object
	stats := map[string]interface{}{
		"total_thoughts":         totalThoughts,
		"average_length":         float64(int(avgLength*100)) / 100, // Round to 2 decimal places
		"longest_thought_index":  longestIndex + 1,
		"longest_thought_length": longestLength,
	}

	// Convert stats to JSON
	statsJSON, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error creating stats JSON: %v", err)
	}

	return mcp.NewToolResultText(string(statsJSON)), nil
}

// Run starts the MCP server with the specified transport
func (s *ThinkToolServer) Run(transport string) error {
	fmt.Printf("Starting Think Tool MCP Server with %s transport...\n", transport)

	switch transport {
	case "stdio":
		return server.ServeStdio(s.server)
	default:
		return fmt.Errorf("unsupported transport: %s", transport)
	}
}

func main() {
	server := NewThinkToolServer("")
	if err := server.Run("stdio"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
