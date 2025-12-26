# PTY + Claude Code Integration Findings

## Summary

Successfully tested using Go's PTY package to create terminal sessions with Claude Code and capture output.

## Mode Comparison

| Aspect | `--print` Mode | Interactive Mode |
|--------|---------------|------------------|
| **Complexity** | Simple | Complex |
| **Setup** | Minimal | Requires timing/state detection |
| **Output Size** | ~23 bytes | ~129KB (with all UI updates) |
| **ANSI Codes** | Minimal (`[?25h` cursor) | Heavy (colors, cursor movement, redraws) |
| **Use Case** | Single prompt/response | Multi-turn conversations |
| **Recommended** | Yes, for automation | For interactive sessions only |

## Key Findings

### 1. `--print` Mode (Recommended for Automation)
- **Pros:**
  - Simple: single command, single response
  - Minimal output processing needed
  - Clean output with few escape codes
  - Fast and predictable
- **Cons:**
  - Single-turn only
  - No access to Claude's full UI features

### 2. Interactive Mode
- **Pros:**
  - Full Claude experience
  - Multi-turn conversations possible
  - Access to all features (/exit, /clear, etc.)
- **Cons:**
  - Complex: requires detecting ready state
  - Large output (lots of UI refresh data)
  - Character-by-character typing needed
  - Timing-sensitive

## Technical Notes

### PTY Package
Using `github.com/creack/pty` v1.1.24:
```go
cmd := exec.Command("claude", "--print", prompt)
ptmx, err := pty.Start(cmd)
// ptmx is read/write - output and input
```

### ANSI Code Handling
Common escape codes from Claude:
- `[?25h` / `[?25l` - Show/hide cursor
- `[38;2;R;G;Bm` - RGB foreground color
- `[2K` - Clear line
- `[1A` - Cursor up 1 line
- `[?2004h` - Bracketed paste mode

### Input to Interactive Mode
- Must type character-by-character for reliable input
- Use `\r` (carriage return) instead of `\n` for Enter key
- Wait for UI to stabilize before sending input

## Files

```
├── main.go           # Entry point with mode selection
├── pty_session.go    # PTY wrapper for interactive sessions
├── ansi_strip.go     # ANSI escape code removal
├── go.mod            # Go module definition
└── logs/             # Output logs
    ├── session_raw_*.log    # With ANSI codes
    └── session_clean_*.log  # Plain text
```

## Usage

```bash
# Simple prompt (recommended)
go run . -mode=print

# Interactive session
go run . -mode=interactive
```

## Recommendation

For automation purposes, **use `--print` mode**. It's simpler, faster, and produces cleaner output. Interactive mode is useful for testing but adds unnecessary complexity for automation.

For programmatic Claude interaction, also consider the Claude API directly instead of PTY.
