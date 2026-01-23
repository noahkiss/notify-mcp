package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const version = "0.2.1"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			runInit()
			return
		case "version", "--version", "-v":
			fmt.Printf("notify-mcp %s\n", version)
			return
		case "help", "--help", "-h":
			printHelp()
			return
		}
	}

	runServer()
}

func printHelp() {
	fmt.Println(`notify-mcp - MCP server for push notifications via Apprise

Usage:
  notify-mcp          Start the MCP server (stdio)
  notify-mcp init     Configure Claude Code integration
  notify-mcp version  Show version

Environment variables (when running as MCP server):
  APPRISE_URL   Apprise API endpoint (required)
  APPRISE_TAG   Notification routing tag (default: apprise)`)
}

func runServer() {
	s := server.NewMCPServer("notify-mcp", version,
		server.WithToolCapabilities(true),
	)

	tool := mcp.NewTool("send_notification",
		mcp.WithDescription("Send a push notification to the user's devices. Use when: long-running tasks complete, you need user input, or important errors occur."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Short notification title")),
		mcp.WithString("message", mcp.Required(), mcp.Description("Notification body with details")),
	)

	s.AddTool(tool, handleNotification)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func handleNotification(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	message, err := req.RequireString("message")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	appriseURL := os.Getenv("APPRISE_URL")
	if appriseURL == "" {
		return mcp.NewToolResultError("APPRISE_URL environment variable not set"), nil
	}

	tag := os.Getenv("APPRISE_TAG")
	if tag == "" {
		tag = "apprise"
	}

	payload := map[string]string{
		"tag":   tag,
		"title": title,
		"body":  message,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(appriseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to send: %v", err)), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 204 {
		return mcp.NewToolResultText("Notification sent successfully"), nil
	}
	return mcp.NewToolResultError(fmt.Sprintf("HTTP %d", resp.StatusCode)), nil
}

func runInit() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("notify-mcp setup")
	fmt.Println("================")
	fmt.Println()

	// Get binary path
	binaryPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding executable path: %v\n", err)
		os.Exit(1)
	}
	binaryPath, _ = filepath.EvalSymlinks(binaryPath)

	fmt.Printf("Binary location: %s\n\n", binaryPath)

	// Prompt for Apprise URL
	fmt.Print("Apprise URL (e.g., http://localhost:8000/notify/apprise): ")
	appriseURL, _ := reader.ReadString('\n')
	appriseURL = strings.TrimSpace(appriseURL)
	if appriseURL == "" {
		fmt.Println("Error: Apprise URL is required")
		os.Exit(1)
	}

	// Prompt for tag
	fmt.Print("Apprise tag [apprise]: ")
	tag, _ := reader.ReadString('\n')
	tag = strings.TrimSpace(tag)
	if tag == "" {
		tag = "apprise"
	}

	// Build the MCP server config
	serverConfig := map[string]any{
		"command": binaryPath,
		"env": map[string]string{
			"APPRISE_URL": appriseURL,
			"APPRISE_TAG": tag,
		},
		"instructions": "Send a notification when: (1) a task takes longer than 60 seconds to complete, (2) you need user input to continue, (3) a task fails with an error the user should know about. Keep notifications concise.",
	}

	// Load or create settings file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
		os.Exit(1)
	}

	settingsPath := filepath.Join(homeDir, ".claude", "settings.json")
	settingsDir := filepath.Dir(settingsPath)

	// Ensure .claude directory exists
	if err := os.MkdirAll(settingsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating settings directory: %v\n", err)
		os.Exit(1)
	}

	// Load existing settings or create new
	var settings map[string]any
	data, err := os.ReadFile(settingsPath)
	if err == nil {
		if err := json.Unmarshal(data, &settings); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing existing settings: %v\n", err)
			os.Exit(1)
		}
	} else {
		settings = make(map[string]any)
	}

	// Ensure mcpServers exists
	mcpServers, ok := settings["mcpServers"].(map[string]any)
	if !ok {
		mcpServers = make(map[string]any)
		settings["mcpServers"] = mcpServers
	}

	// Add/update notify server
	mcpServers["notify"] = serverConfig

	// Write back
	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding settings: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(settingsPath, output, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing settings: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("Configuration saved to %s\n", settingsPath)
	fmt.Println()
	fmt.Println("Restart Claude Code to load the new MCP server.")
}
