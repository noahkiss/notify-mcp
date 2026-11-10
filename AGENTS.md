# notify-mcp Agent Instructions

## Release Workflow (ALWAYS follow after semver changes)

After ANY change that warrants a version bump (bug fixes, features, breaking changes):

1. **Bump version** following semver (major.minor.patch)

2. **Tag and push the release**
   ```bash
   git tag v<version>
   git push origin v<version>
   ```

3. **Update the Homebrew formula** `Formula/notify-mcp.rb` in the [noahkiss/homebrew-tap](https://github.com/noahkiss/homebrew-tap) repo:
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

The server sends its instructions over the MCP protocol with `server.WithInstructions()`, so a client receives them when it connects, with no manual config.

The instructions tell Claude to send notifications when:
1. A task takes longer than 60 seconds to complete
2. User input is needed to continue
3. A task fails with an error the user should know about

The `init` command also writes instructions to `~/.claude.json` as a fallback.

## Formula Location

`Formula/notify-mcp.rb` in the [noahkiss/homebrew-tap](https://github.com/noahkiss/homebrew-tap) repo.

## Installation

Users install via: `brew install noahkiss/tap/notify-mcp`
