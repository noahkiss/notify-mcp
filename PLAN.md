# Plan: Create notify-mcp GitHub Repository

## Context

This plan creates an MCP (Model Context Protocol) server in Go that sends push notifications via Apprise. The goal is to allow Claude Code to notify the user when long-running tasks complete or input is needed.

**Why Go?** Compiles to a single binary with no runtime dependencies - users just download and run.

**Why MCP?** Gives Claude an explicit tool with clear semantics for when to notify, rather than relying on bash access.

## Project Location

`~/develop/notify-mcp`

## Repository Structure

```
notify-mcp/
├── main.go                        # MCP server implementation
├── go.mod
├── go.sum
├── README.md                      # Setup instructions
├── LICENSE                        # MIT
├── .goreleaser.yaml               # Release automation
└── .github/
    └── workflows/
        └── release.yml            # GitHub Actions for releases
```

## Technical Details

### MCP Server Implementation

Using [mcp-go](https://github.com/mark3labs/mcp-go) SDK (MIT licensed, well-documented):

```go
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    // Create MCP server
    s := server.NewMCPServer("notify-mcp", "0.1.0",
        server.WithToolCapabilities(true),
    )

    // Define the notification tool
    tool := mcp.NewTool("send_notification",
        mcp.WithDescription("Send a push notification to the user's devices. Use when: long-running tasks complete, you need user input, or important errors occur."),
        mcp.WithString("title", mcp.Required(), mcp.Description("Short notification title")),
        mcp.WithString("message", mcp.Required(), mcp.Description("Notification body with details")),
    )

    // Add tool handler
    s.AddTool(tool, handleNotification)

    // Run over stdio
    if err := server.ServeStdio(s); err != nil {
        fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
        os.Exit(1)
    }
}

func handleNotification(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    title, _ := req.RequireString("title")
    message, _ := req.RequireString("message")

    appriseURL := os.Getenv("APPRISE_URL")
    if appriseURL == "" {
        return mcp.NewToolResultError("APPRISE_URL environment variable not set"), nil
    }

    tag := os.Getenv("APPRISE_TAG")
    if tag == "" {
        tag = "general"
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
```

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APPRISE_URL` | Yes | - | Full Apprise API endpoint (e.g., `http://10.10.10.30:2777/notify/apprise`) |
| `APPRISE_TAG` | No | `general` | Apprise routing tag |

### Claude Code Configuration

Once built, users configure in `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "notify": {
      "command": "/path/to/notify-mcp",
      "env": {
        "APPRISE_URL": "http://your-apprise-server:8000/notify/apprise"
      },
      "instructions": "Use send_notification when long-running tasks complete or you need user input to continue."
    }
  }
}
```

---

## Execution Steps

1. **Initialize Go module**
   ```bash
   go mod init github.com/YOUR_USERNAME/notify-mcp
   ```

2. **Write main.go** - implement the MCP server (code above)

3. **Add dependency**
   ```bash
   go get github.com/mark3labs/mcp-go
   go mod tidy
   ```

4. **Build and test locally**
   ```bash
   go build -o notify-mcp
   # Test: echo '{"jsonrpc":"2.0","method":"initialize","id":1,"params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | ./notify-mcp
   ```

5. **Write README.md** - installation, configuration, usage

6. **Write LICENSE** - MIT

7. **Create GitHub repo**
   ```bash
   git init
   git add .
   git commit -m "Initial commit: notify-mcp server"
   gh repo create notify-mcp --public --source=. --push
   ```

8. **Add goreleaser config** (`.goreleaser.yaml`) for multi-platform builds

9. **Add GitHub Actions workflow** (`.github/workflows/release.yml`) to auto-build on tags

10. **Tag first release**
    ```bash
    git tag v0.1.0
    git push origin v0.1.0
    ```

11. **Configure locally** - add to `~/.claude/settings.json`

## Verification

1. Build succeeds: `go build -o notify-mcp`
2. MCP handshake works: send initialize JSON, get valid response
3. Tool works: configure in Claude Code, ask Claude to send a test notification
4. GitHub release: tag triggers workflow, binaries appear in releases

## Notes

- The original `/home/flight/bin/notify` bash script is unchanged - keep it for direct CLI use
- This MCP server is separate and purpose-built for Claude integration
