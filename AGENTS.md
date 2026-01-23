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

As of v0.2.3, instructions are sent via the MCP protocol using `server.WithInstructions()`. This means Claude receives the instructions automatically when connecting - no manual config needed.

The instructions tell Claude to send notifications when:
1. A task takes longer than 60 seconds to complete
2. User input is needed to continue
3. A task fails with an error the user should know about

The `init` command also writes instructions to `~/.claude.json` as a fallback.

## Current Formula Location

`~/develop/homebrew-tap/Formula/notify-mcp.rb`

## Installation

Users install via: `brew install noahkiss/tap/notify-mcp`
