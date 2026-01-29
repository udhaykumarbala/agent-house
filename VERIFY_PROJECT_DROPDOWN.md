# How to Verify Project Dropdown is Working

## Step 1: Start the Server

```bash
go run cmd/agent-house/main.go --serve
```

Wait for:
```
Starting server on http://localhost:8080
WebSocket available at ws://localhost:8080/ws
```

## Step 2: Test the API Endpoint

In another terminal:

```bash
curl http://localhost:8080/api/projects | python3 -m json.tool
```

**Expected Output:**
```json
{
  "projects": [
    {
      "id": "reorder-alert",
      "name": "reorder-alert",
      "task_count": 1,
      "last_updated": "2025-01-29T12:34:56Z"
    },
    {
      "id": "reorder-system",
      "name": "reorder-system",
      "task_count": 1,
      "last_updated": "2025-01-29T11:20:00Z"
    },
    {
      "id": "subscription-manager",
      "name": "subscription-manager",
      "task_count": 1,
      "last_updated": "2025-01-29T10:15:30Z"
    },
    {
      "id": "calculator",
      "name": "calculator",
      "task_count": 0,
      "last_updated": "2025-01-28T..."
    }
    // ... more projects
  ],
  "count": 8
}
```

✅ **If you see this**, the API is working correctly!

## Step 3: Check the Dashboard

1. Open browser: http://localhost:8080

2. Look at the **left sidebar** at the top

3. You should see a dropdown:
   ```
   ┌────────────────────────────┐
   │ Projects                   │
   │ ┌────────────────────────┐ │
   │ │ Select Project       ▼ │ │ ← Click this!
   │ └────────────────────────┘ │
   └────────────────────────────┘
   ```

4. Click the dropdown - you should see:
   ```
   ┌────────────────────────────┐
   │ Default Project            │
   │ reorder-alert (1 tasks)    │
   │ reorder-system (1 tasks)   │
   │ subscription-manager (1 task)│
   │ calculator (0 tasks)       │
   │ todo-test (0 tasks)        │
   │ ...                        │
   └────────────────────────────┘
   ```

## Step 4: Test Switching Projects

1. Select "reorder-alert" from dropdown
2. You should see a toast: "Project Switched - Now viewing: reorder-alert"
3. The task history updates to show reorder-alert tasks
4. The files panel updates to show reorder-alert files

## Step 5: Create a New Project

1. In the main area, type a new project name: `"test-dropdown"`
2. Enter a task: `"Test the dropdown"`
3. Submit
4. Wait 30 seconds (or refresh the page)
5. Click the dropdown again
6. You should see `"test-dropdown"` appear!

## Troubleshooting

### Dropdown Only Shows "Default Project"?

**Check browser console (F12):**
```javascript
// Open Console and run:
fetch('/api/projects').then(r => r.json()).then(console.log)

// Should show array of projects
```

**If you see an error:**
- Check server logs for errors
- Verify API endpoint is registered
- Check if projects directory exists

### Dropdown Not Updating?

**Force refresh:**
- Press F5 to reload page
- Or wait 30 seconds for auto-refresh

**Check network tab (F12 → Network):**
- Look for requests to `/api/projects`
- Should happen every 30 seconds
- Check response has your projects

### Projects Missing?

**Verify projects directory:**
```bash
ls -la projects/

# Should show all your project folders:
# drwxr-xr-x  calculator
# drwxr-xr-x  reorder-alert
# drwxr-xr-x  subscription-manager
# etc.
```

**Check project structure:**
```bash
# Projects WITH tasks will show:
ls -la projects/reorder-alert/.tasks/
# Should see: history.json

# Projects WITHOUT tasks yet won't have .tasks/ but will still appear in dropdown
```

## What Should Work Now

✅ **Storage** - All projects are stored in `projects/` directory
✅ **Retrieval** - API scans directory and reads task history
✅ **Dropdown** - Populated with all projects (with task counts)
✅ **Switching** - Click to switch between projects instantly
✅ **Auto-update** - Refreshes every 30 seconds
✅ **New projects** - Appear after creation (within 30 seconds)

## Expected Behavior Summary

### On Page Load:
1. ✅ Fetches `/api/projects`
2. ✅ Populates dropdown with all projects
3. ✅ Shows task count: `"project-name (5 tasks)"`
4. ✅ Sorted by most recent activity

### When Creating New Project:
1. ✅ Submit task with new project name
2. ✅ Project folder created in `projects/`
3. ✅ After 30 seconds: dropdown auto-updates
4. ✅ New project appears in list

### When Switching:
1. ✅ Select project from dropdown
2. ✅ Dashboard reloads files/tasks for that project
3. ✅ Toast notification confirms switch
4. ✅ Project input field syncs with dropdown

## Quick Test

**One command to test everything:**

```bash
# Terminal 1 - Start server
go run cmd/agent-house/main.go --serve

# Terminal 2 - Test API
curl http://localhost:8080/api/projects

# Browser
# Open http://localhost:8080
# Click dropdown in left sidebar
# See all 8 projects!
```

If you see all your projects in the dropdown, it's working! ✅
