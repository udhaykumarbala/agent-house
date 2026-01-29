# Parallel Execution Quick Start Guide

## What Changed?

The agent-house-2 system now executes agents **in parallel** within each phase, dramatically improving performance. Previously, agents ran one-at-a-time (sequential). Now they run simultaneously (parallel).

## Performance Comparison

### Before (Sequential)
```
Phase: Research
  [1] PM executes...           (2 min)
  [2] UX executes...           (2 min)
  [3] UI executes...           (2 min)
  [4] Security executes...     (2 min)
  [5] Architect executes...    (2 min)
Total: 10 minutes
```

### After (Parallel)
```
Phase: Research
  [1] PM executes...
  [2] UX executes...
  [3] UI executes...           All 5 run simultaneously
  [4] Security executes...
  [5] Architect executes...
Total: 2 minutes (5x faster!)
```

## How It Works

### Phase-Based Parallel Execution

Each phase runs agents in parallel, but phases run sequentially:

```
Phase 0: Template Selection (sequential - needs decisions)
  └─ Architect proposes → CEO/PM review

Phase 1: Research (PARALLEL)
  ├─ PM researches
  ├─ UX researches
  ├─ UI researches
  ├─ Security researches
  └─ Architect researches

Phase 2: Planning (PARALLEL)
  ├─ PM creates spec
  ├─ UX creates spec
  ├─ UI creates spec
  ├─ Security creates spec
  └─ Architect creates spec

Phase 3: Discussion (sequential - needs consensus)
  └─ CEO reviews and approves

Phase 4: Development (PARALLEL)
  ├─ Senior Dev implements
  └─ Junior Dev implements
```

## Observing Parallel Execution

### Server Logs

Watch for these indicators in the server output:

```
🚀 Executing 5 agents in parallel...

✅ PM completed:
   └─ Research findings...

✅ UX completed:
   └─ User research...

✅ UI completed:
   └─ Design analysis...

✅ Security completed:
   └─ Security assessment...

✅ Architect completed:
   └─ Technical architecture...

⏱️  Phase completed in 2m15s (max agent: 2m10s, speedup: 4.8x)
```

### Key Metrics

- **Phase Duration**: Total time for the phase
- **Max Agent Duration**: Longest individual agent
- **Speedup**: How much faster than sequential (e.g., 4.8x)

## Running Multiple Projects Concurrently

You can now run multiple projects **at the same time**:

### Terminal 1
```bash
curl -X POST http://localhost:8080/api/task \
  -H "Content-Type: application/json" \
  -d '{"project_id": "project-a", "task": "Build a blog"}'
```

### Terminal 2
```bash
curl -X POST http://localhost:8080/api/task \
  -H "Content-Type: application/json" \
  -d '{"project_id": "project-b", "task": "Build an e-commerce site"}'
```

Both projects run **independently** without blocking each other!

## API Endpoints for Agent Tasks

### Get Real-Time Agent Status

```bash
# Check what PM is doing
curl http://localhost:8080/api/agents/pm/status

# Response:
{
  "role": "pm",
  "state": "working",              # idle, working, waiting, completed, error
  "current_task_id": "at_123",
  "current_task_type": "research",
  "tasks_completed": 5,
  "tasks_failed": 0,
  "last_active_at": "2025-01-29T12:34:56Z"
}
```

### Get Agent Tasks

```bash
# All tasks for PM
curl http://localhost:8080/api/agents/pm/tasks

# All pending tasks
curl http://localhost:8080/api/agent-tasks?status=pending

# All tasks for a specific project
curl http://localhost:8080/api/agent-tasks?project=my-project
```

### Get Agent Outputs

```bash
# All outputs from UX designer
curl http://localhost:8080/api/agents/ux/outputs

# Response shows research findings, specs, etc.
{
  "role": "ux",
  "outputs": [
    {
      "type": "research_findings",
      "title": "User Research Summary",
      "file_path": ".tasks/agent-tasks/ux/outputs/research.md"
    }
  ]
}
```

## WebSocket Real-Time Updates

Connect to the WebSocket to get live updates:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);

  // Regular message updates
  if (data.type === 'message') {
    console.log('Message:', data.message);
  }

  // Agent task events
  if (data.type === 'agent_task_event') {
    const event = data.event;

    switch (event.event_type) {
      case 'task_started':
        console.log(`${event.agent_role} started: ${event.data.title}`);
        break;

      case 'task_progress':
        console.log(`${event.agent_role} progress: ${event.data.progress}%`);
        break;

      case 'task_completed':
        console.log(`${event.agent_role} completed in ${event.data.duration}`);
        break;

      case 'task_failed':
        console.error(`${event.agent_role} failed: ${event.data.error}`);
        break;
    }
  }
};
```

## File Organization

Agent task data is stored in `.tasks/agent-tasks/`:

```
projects/my-project/
├── .tasks/
│   ├── history.json                    # Main task history
│   └── agent-tasks/                    # NEW: Agent-specific tasks
│       ├── index.json                  # Task index
│       ├── pm/
│       │   ├── at_123.json            # PM's task
│       │   └── outputs/
│       │       └── research.md        # PM's research output
│       ├── ux/
│       │   ├── at_124.json            # UX's task
│       │   └── outputs/
│       │       └── research.md        # UX's research output
│       └── architect/
│           ├── at_125.json            # Architect's task
│           └── outputs/
│               └── architecture.md    # Architecture doc
```

## Configuration

Parallel execution is **enabled by default**. You can configure it:

```go
// In orchestrator.go or server initialization
executorConfig := ExecutorConfig{
    EnableParallel: true,           // Enable parallel execution
    MaxWorkers:     0,              // 0 = unlimited (use CPU cores)
    AgentTimeout:   5 * time.Minute, // Per-agent timeout
    PhaseTimeout:   15 * time.Minute, // Per-phase timeout
}
```

## Benefits You'll See

### 1. Faster Task Completion
- Research phase: ~5x faster
- Planning phase: ~5x faster
- Development phase: ~2x faster
- Overall: 3-5x improvement

### 2. Better Resource Utilization
- All CPU cores used during parallel phases
- Agents no longer wait for each other
- Better throughput

### 3. Concurrent Projects
- Run multiple projects at once
- Each project has independent execution
- No blocking between projects

### 4. Better Visibility
- See what each agent is doing in real-time
- Track progress per agent
- Structured outputs by agent

## Troubleshooting

### All agents seem sequential?

Check logs for the parallel execution message:
```
🚀 Executing N agents in parallel...
```

If you don't see it, parallel execution might be disabled.

### File conflicts?

The ConcurrentFileManager prevents duplicates with "first-wins" policy:
```
⚠️  Skipping file.md - already created by pm
```

This is expected behavior when multiple agents try to create the same file.

### Agent timeouts?

If agents take too long:
```
⚠️  pm timed out or cancelled: context deadline exceeded
```

Increase timeouts in configuration.

### Memory usage high?

With 5-10 agents running in parallel, memory usage will be higher. This is expected and normal. Monitor with:
```bash
ps aux | grep agent-house
```

## Example: Full Task Flow

```bash
# 1. Start server
go run cmd/agent-house/main.go --serve

# 2. Submit task
curl -X POST http://localhost:8080/api/task \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "todo-app",
    "task": "Build a React todo application with authentication"
  }'

# 3. Watch logs for parallel execution
# You'll see:
#   Phase: Research
#   🚀 Executing 5 agents in parallel...
#   [Multiple agents working simultaneously]
#   ⏱️  Phase completed in 2m (speedup: 4.5x)

# 4. Check agent status
curl http://localhost:8080/api/agents/pm/status

# 5. Get all outputs
curl http://localhost:8080/api/agent-tasks?project=todo-app

# 6. View results
ls projects/todo-app/
ls projects/todo-app/.tasks/agent-tasks/
```

## Performance Tips

### 1. Use Appropriate Hardware
- More CPU cores = better parallelism
- SSD helps with file I/O during parallel writes

### 2. Monitor Resource Usage
```bash
# CPU usage
top -p $(pgrep -f agent-house)

# Memory usage
ps aux | grep agent-house
```

### 3. Adjust Worker Limits
If system resources are constrained:
```go
MaxWorkers: 3,  // Limit to 3 concurrent agents instead of all
```

### 4. Project Isolation
Each project runs independently, so you can:
- Run different-sized projects together
- Stop one without affecting others
- Isolate resource-intensive tasks

## Summary

✅ **Parallel execution** speeds up phases by 3-5x
✅ **Concurrent projects** run simultaneously
✅ **Real-time tracking** via API and WebSocket
✅ **Better organization** with agent-specific outputs
✅ **Thread-safe** file operations prevent conflicts
✅ **Configurable** timeouts and worker limits

The system automatically manages parallelism - just submit tasks and watch them complete faster!
