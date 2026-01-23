# notify-mcp

An MCP (Model Context Protocol) server that sends push notifications via [Apprise](https://github.com/caronc/apprise). Useful for notifying when long-running tasks complete or when Claude needs user input.

## Installation

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/noahkiss/notify-mcp/releases).

### Build from Source

```bash
go install github.com/noahkiss/notify-mcp@latest
```

Or clone and build:

```bash
git clone https://github.com/noahkiss/notify-mcp.git
cd notify-mcp
go build -o notify-mcp
```

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APPRISE_URL` | Yes | - | Full Apprise API endpoint (e.g., `http://localhost:8000/notify/apprise`) |
| `APPRISE_TAG` | No | `general` | Apprise routing tag |

### Claude Code Setup

Add to `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "notify": {
      "command": "/path/to/notify-mcp",
      "env": {
        "APPRISE_URL": "http://your-apprise-server:8000/notify/apprise"
      }
    }
  }
}
```

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
