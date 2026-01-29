# Quick Start Guide

## Start the Server

```bash
go run cmd/agent-house/main.go --serve
```

The server will start on http://localhost:8080

## Open the Dashboard

Open your browser to: http://localhost:8080

## Try It Out

### 1. Submit a Task

In the dashboard:
- Enter a project name (e.g., "todo-app")
- Enter a task (e.g., "Build a React todo application")
- Click "Start Task"

### 2. Watch Parallel Execution

You'll see in the logs:
```
🚀 Executing 5 agents in parallel...
✅ PM completed: Research findings...
✅ UX completed: User research...
⏱️  Phase completed in 2m15s (speedup: 4.8x)
```

### 3. Click on Agent Cards

In the right sidebar:
- Click on any agent card (PM, UX, UI, etc.)
- See real-time status, tasks, and outputs
- Watch progress bars update live

### 4. Use the API

```bash
# Check status
curl http://localhost:8080/api/status

# Get agent status
curl http://localhost:8080/api/agents/pm/status

# Get all agent tasks
curl http://localhost:8080/api/agent-tasks

# Submit task via API
curl -X POST http://localhost:8080/api/task \
  -H "Content-Type: application/json" \
  -d '{"project_id": "test", "task": "Build a blog"}'
```

## Key Features

✅ **Parallel Execution** - Agents run simultaneously (3-5x faster)
✅ **Real-Time Updates** - WebSocket events show live progress
✅ **Agent Details** - Click any agent to see status and tasks
✅ **Multiple Projects** - Run several projects concurrently
✅ **Task Tracking** - See progress, outputs, and history

## Documentation

- `IMPLEMENTATION_SUMMARY.md` - Technical details
- `PARALLEL_EXECUTION_GUIDE.md` - User guide
- `FINAL_IMPLEMENTATION_SUMMARY.md` - Complete overview

## Testing

```bash
# Run all tests
go test ./...

# Run specific tests
go test ./internal/orchestrator/ -v
go test ./internal/agentask/ -v

# Run benchmarks
go test ./internal/orchestrator/ -bench=.
```

## Troubleshooting

### Port already in use?
```bash
go run cmd/agent-house/main.go --serve --port 8081
```

### See detailed logs?
The `--verbose` flag is enabled by default. Logs show:
- Parallel execution indicators
- Agent completions
- Phase durations with speedup metrics

### Multiple projects at once?
Just submit tasks with different `project_id` values - they'll run independently!

## Performance

**Before:** ~26 minutes for typical multi-agent task
**After:** ~7 minutes (3.7x faster!)

Enjoy the speed boost! 🚀
