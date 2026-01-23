# notify-mcp Agent Instructions

## Releasing a New Version

When making changes that should be released:

1. **Tag the release**
   ```bash
   git tag v<version>
   git push origin v<version>
   ```

2. **Update the Homebrew formula** in `~/develop/homebrew-tap/Formula/notify-mcp.rb`:
   - Update the `url` line with the new tag version
   - Get the new SHA256: `curl -sL https://github.com/noahkiss/notify-mcp/archive/refs/tags/v<version>.tar.gz | sha256sum`
   - Update the `sha256` line with the new hash
   - Commit and push the homebrew-tap repo

## Current Formula Location

```
~/develop/homebrew-tap/Formula/notify-mcp.rb
```

## Installation

Users install via: `brew install noahkiss/tap/notify-mcp`
