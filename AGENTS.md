# notify-mcp Agent Instructions

## Release Workflow (ALWAYS follow after semver changes)

After ANY change that warrants a version bump (bug fixes, features, breaking changes):

1. **Bump version** following semver (major.minor.patch)

2. **Tag and push the release**
   ```bash
   git tag v<version>
   git push origin v<version>
   ```

3. **Update Homebrew formula** in `~/develop/homebrew-tap/Formula/notify-mcp.rb`:
   ```bash
   # Get new SHA256
   curl -sL https://github.com/noahkiss/notify-mcp/archive/refs/tags/v<version>.tar.gz | sha256sum
   ```
   - Update `url` with new tag
   - Update `sha256` with new hash
   - Commit and push homebrew-tap

4. **Verify** the release is installable:
   ```bash
   brew update && brew upgrade notify-mcp
   ```

## MCP Server Instructions

When users configure this server, they should add an `instructions` field to ensure Claude uses it proactively. Add to `~/.claude.json` under the notify server config:

```json
"instructions": "Send a notification when: (1) a task takes longer than 60 seconds to complete, (2) you need user input to continue, (3) a task fails with an error the user should know about. Keep notifications concise."
```

The `init` command writes this automatically, but `claude mcp add` does not support instructions, so manual addition may be needed.

## Current Formula Location

`~/develop/homebrew-tap/Formula/notify-mcp.rb`

## Installation

Users install via: `brew install noahkiss/tap/notify-mcp`
