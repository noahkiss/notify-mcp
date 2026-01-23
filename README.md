# notify-mcp

An MCP (Model Context Protocol) server that sends push notifications via [Apprise](https://github.com/caronc/apprise). Useful for notifying when long-running tasks complete or when Claude needs user input.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install noahkiss/tap/notify-mcp
```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/noahkiss/notify-mcp/releases).

### Build from Source

```bash
go install github.com/noahkiss/notify-mcp@latest
```

## Setup

After installation, run the setup wizard:

```bash
notify-mcp init
```

This will:
1. Prompt for your Apprise URL
2. Prompt for your notification tag
3. Automatically configure Claude Code

Then restart Claude Code to load the MCP server.

## Manual Configuration

If you prefer manual setup, use the Claude CLI:

```bash
claude mcp add -s user notify /path/to/notify-mcp \
  -e APPRISE_URL="http://your-apprise-server:8000/notify/apprise" \
  -e APPRISE_TAG="apprise"
```

The `-s user` flag makes the server available across all projects. Alternatively, add a `.mcp.json` file to your project root for project-specific configuration:

```json
{
  "mcpServers": {
    "notify": {
      "command": "/path/to/notify-mcp",
      "env": {
        "APPRISE_URL": "http://your-apprise-server:8000/notify/apprise",
        "APPRISE_TAG": "apprise"
      }
    }
  }
}
```

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APPRISE_URL` | Yes | - | Full Apprise API endpoint |
| `APPRISE_TAG` | No | `apprise` | Apprise routing tag |

## Usage

The server exposes a single tool:

### send_notification

Send a push notification to configured devices.

**Parameters:**
- `title` (string, required): Short notification title
- `message` (string, required): Notification body with details

**Example:**
```json
{
  "title": "Build Complete",
  "message": "The project compiled successfully with no errors."
}
```

## Apprise Setup

This server requires a running [Apprise API](https://github.com/caronc/apprise-api) instance. You can run one with Docker:

```bash
docker run -d --name apprise \
  -p 8000:8000 \
  -v /path/to/config:/config \
  caronc/apprise:latest
```

Then configure your notification services (Pushover, Telegram, Discord, etc.) in the Apprise web interface.

## License

MIT
