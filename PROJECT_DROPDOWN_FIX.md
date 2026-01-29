# Project Dropdown - Fixed & Working! ✅

## What Was Fixed

The project dropdown in the left sidebar was not showing previously created projects. This has been fixed by implementing:

### 1. **New API Endpoint** (`/api/projects`)
Returns a list of all projects with metadata:

```json
{
  "projects": [
    {
      "id": "todo-app",
      "name": "todo-app",
      "task_count": 5,
      "last_updated": "2025-01-29T12:34:56Z"
    },
    {
      "id": "blog-system",
      "name": "blog-system",
      "task_count": 3,
      "last_updated": "2025-01-29T11:20:00Z"
    }
  ],
  "count": 2
}
```

### 2. **Project Discovery**
The API scans the `projects/` directory and:
- ✅ Finds all project folders
- ✅ Reads task history from `.tasks/history.json`
- ✅ Counts total tasks per project
- ✅ Gets last updated timestamp
- ✅ Sorts by most recently updated

### 3. **Dropdown Population**
The dashboard now:
- ✅ Fetches projects on page load
- ✅ Populates the dropdown with project names
- ✅ Shows task count next to each project
- ✅ Auto-refreshes every 30 seconds
- ✅ Maintains selection when switching

## How It Works

### On Page Load:
```
1. Dashboard loads
2. fetchProjects() calls /api/projects
3. Dropdown populated with all projects
4. Current project is pre-selected
```

### When You Create a New Project:
```
1. Submit task with new project name
2. After 30 seconds (or page refresh)
3. New project appears in dropdown
4. Can switch to it immediately
```

### When You Switch Projects:
```
1. Select project from dropdown
2. Dashboard reloads:
   - Files for that project
   - Task history for that project
   - Messages for that project
3. Toast notification confirms switch
```

## Where Is It?

### Left Sidebar - Top Section:
```
┌─────────────────────────────────┐
│ Projects                        │
│ ┌─────────────────────────────┐ │
│ │ Select Project            ▼ │ │ ← Click here!
│ └─────────────────────────────┘ │
│                                 │
│ [Search tasks...]               │
│                                 │
│ + New Task                      │
│                                 │
│ Task History...                 │
└─────────────────────────────────┘
```

## Example Usage

### Scenario 1: View Existing Projects
```
1. Start server: go run cmd/agent-house/main.go --serve
2. Open: http://localhost:8080
3. Click the dropdown in left sidebar
4. See all your projects:
   - Default Project
   - todo-app (5 tasks)
   - blog-system (3 tasks)
   - e-commerce (1 task)
```

### Scenario 2: Switch Between Projects
```
1. Currently viewing: "todo-app"
2. Click dropdown
3. Select: "blog-system"
4. Dashboard updates to show blog-system files/tasks
5. Toast: "Project Switched - Now viewing: blog-system"
```

### Scenario 3: Create New Project
```
1. Type "new-project" in the project input field
2. Submit a task
3. Wait 30 seconds (or refresh page)
4. Open dropdown - "new-project" now appears!
```

## What Shows in the Dropdown

Each project entry shows:
```
project-name (X tasks)
```

For example:
```
┌─────────────────────────────────┐
│ Default Project                 │
│ todo-app (5 tasks)              │ ← Most recent
│ blog-system (3 tasks)           │
│ e-commerce (1 task)             │
│ test-project (0 tasks)          │ ← Oldest
└─────────────────────────────────┘
```

Projects are sorted by **most recently updated** at the top.

## API Testing

You can also query the projects API directly:

```bash
# Get all projects
curl http://localhost:8080/api/projects

# Response:
{
  "projects": [
    {
      "id": "todo-app",
      "name": "todo-app",
      "task_count": 5,
      "last_updated": "2025-01-29T12:34:56Z"
    },
    {
      "id": "blog-system",
      "name": "blog-system",
      "task_count": 3,
      "last_updated": "2025-01-29T11:20:00Z"
    }
  ],
  "count": 2
}
```

## How Data Is Stored & Retrieved

### Storage:
```
projects/
├── todo-app/
│   ├── .tasks/
│   │   └── history.json          ← Read for task count & timestamp
│   └── [project files...]
├── blog-system/
│   ├── .tasks/
│   │   └── history.json          ← Read for task count & timestamp
│   └── [project files...]
└── default/
    └── [project files...]
```

### Retrieval:
```javascript
1. API scans projects/ directory
2. For each folder:
   - Read .tasks/history.json
   - Count tasks
   - Get last updated time
3. Sort by most recent
4. Return as JSON
```

## Auto-Refresh

The dropdown automatically refreshes every **30 seconds**:
```javascript
setInterval(fetchProjects, 30000);
```

So new projects appear without needing a manual page refresh.

## Synchronization

Both project selectors stay in sync:
1. **Dropdown** (left sidebar) - Select from existing projects
2. **Text Input** (top of main area) - Type new project name

When you:
- Change dropdown → Text input updates
- Change text input → Dropdown selection updates (if project exists)

## Troubleshooting

### Dropdown is Empty?
```bash
# Check if projects directory exists
ls -la projects/

# Should see folders like:
# drwxr-xr-x  default
# drwxr-xr-x  todo-app
# drwxr-xr-x  blog-system
```

### Project Not Showing?
```bash
# Check if project has .tasks directory
ls -la projects/your-project/.tasks/

# Should see:
# -rw-r--r--  history.json
```

### API Not Working?
```bash
# Test the API endpoint
curl http://localhost:8080/api/projects

# Should return JSON with projects array
```

### Dropdown Not Updating?
- Refresh the page (or wait 30 seconds)
- Check browser console for errors (F12)
- Verify API endpoint returns data

## Code Changes Summary

### Files Modified:

1. **`internal/web/server.go`**
   - Added `/api/projects` endpoint
   - Added `handleProjects()` method
   - Scans projects directory and returns metadata

2. **`internal/web/dashboard.go`**
   - Added `fetchProjects()` function
   - Added `setupProjectSelector()` function
   - Wired up dropdown change event
   - Added to init() and polling

### New Functionality:

- ✅ `/api/projects` API endpoint
- ✅ Project discovery and metadata
- ✅ Dropdown population
- ✅ Auto-refresh every 30 seconds
- ✅ Project switching with notifications
- ✅ Task count display
- ✅ Sorted by most recent

## Benefits

### Before:
```
❌ Dropdown only showed "Default Project"
❌ Had to manually type project names
❌ Couldn't see what projects existed
❌ No way to discover previous projects
```

### After:
```
✅ Dropdown shows all projects
✅ Click to switch instantly
✅ See task count for each project
✅ Sorted by most recent activity
✅ Auto-updates every 30 seconds
✅ Toast notifications on switch
```

## Next Steps

Now you can:
1. **Browse Projects**: See all your existing projects
2. **Quick Switch**: One click to change projects
3. **Track Activity**: See which projects are most recent
4. **Monitor Tasks**: View task counts per project

Everything is now properly stored and retrieved! 🎉
