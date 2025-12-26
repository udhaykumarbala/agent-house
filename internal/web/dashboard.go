package web

// dashboardHTML contains the embedded dashboard HTML with WebSocket support
var dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Agent House - Multi-Agent Dashboard</title>
    <style>
        :root {
            --bg-primary: #0f0f23;
            --bg-secondary: #1a1a3e;
            --bg-tertiary: #252550;
            --text-primary: #eee;
            --text-secondary: #aaa;
            --accent: #4ECDC4;
            --accent-hover: #45b7aa;
            --border: #333366;
        }

        * { box-sizing: border-box; margin: 0; padding: 0; }

        body {
            font-family: 'Segoe UI', system-ui, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            min-height: 100vh;
        }

        .header {
            background: var(--bg-secondary);
            padding: 15px 30px;
            border-bottom: 1px solid var(--border);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .header h1 { font-size: 1.5rem; color: var(--accent); }

        .status-bar { display: flex; gap: 20px; font-size: 0.85rem; }
        .status-item { display: flex; align-items: center; gap: 5px; }
        .status-dot { width: 8px; height: 8px; border-radius: 50%; background: #666; }
        .status-dot.connected { background: #2ECC71; }
        .status-dot.running { background: #F39C12; animation: pulse 1s infinite; }

        @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }

        .container {
            display: grid;
            grid-template-columns: 250px 1fr 300px;
            height: calc(100vh - 60px);
        }

        .sidebar {
            background: var(--bg-secondary);
            border-right: 1px solid var(--border);
            padding: 20px;
            overflow-y: auto;
        }

        .sidebar h2 {
            font-size: 0.9rem;
            color: var(--text-secondary);
            margin-bottom: 15px;
            text-transform: uppercase;
            letter-spacing: 1px;
        }

        .agent-card {
            background: var(--bg-tertiary);
            border-radius: 8px;
            padding: 12px;
            margin-bottom: 10px;
            border-left: 3px solid var(--border);
            transition: all 0.2s;
        }

        .agent-card.active { border-left-color: var(--accent); }
        .agent-card h3 { font-size: 0.95rem; margin-bottom: 4px; }
        .agent-card .role { font-size: 0.75rem; color: var(--text-secondary); }

        .main-content { display: flex; flex-direction: column; height: 100%; }

        .task-input-area {
            padding: 20px;
            background: var(--bg-secondary);
            border-bottom: 1px solid var(--border);
        }

        .task-form { display: flex; gap: 10px; }

        .task-form input {
            flex: 1;
            padding: 12px 16px;
            border: 1px solid var(--border);
            border-radius: 8px;
            background: var(--bg-tertiary);
            color: var(--text-primary);
            font-size: 1rem;
        }

        .task-form input:focus { outline: none; border-color: var(--accent); }

        .task-form button {
            padding: 12px 24px;
            border: none;
            border-radius: 8px;
            background: var(--accent);
            color: var(--bg-primary);
            font-weight: 600;
            cursor: pointer;
            transition: background 0.2s;
        }

        .task-form button:hover { background: var(--accent-hover); }
        .task-form button:disabled { background: #666; cursor: not-allowed; }

        .messages-area { flex: 1; overflow-y: auto; padding: 20px; }

        .message {
            background: var(--bg-secondary);
            border-radius: 8px;
            padding: 15px;
            margin-bottom: 15px;
            border-left: 3px solid var(--border);
            animation: slideIn 0.3s ease;
        }

        @keyframes slideIn {
            from { opacity: 0; transform: translateY(-10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .message-header {
            display: flex;
            justify-content: space-between;
            margin-bottom: 10px;
            font-size: 0.85rem;
        }

        .message-from { font-weight: 600; }
        .message-time { color: var(--text-secondary); }

        .message-content {
            white-space: pre-wrap;
            line-height: 1.5;
            font-size: 0.9rem;
        }

        .message-content.collapsed {
            max-height: 200px;
            overflow: hidden;
            position: relative;
        }

        .message-content.collapsed::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 0;
            right: 0;
            height: 50px;
            background: linear-gradient(transparent, var(--bg-secondary));
        }

        .expand-btn {
            background: none;
            border: none;
            color: var(--accent);
            cursor: pointer;
            margin-top: 5px;
            font-size: 0.85rem;
        }

        .files-panel {
            background: var(--bg-secondary);
            border-left: 1px solid var(--border);
            padding: 20px;
            overflow-y: auto;
        }

        .files-panel h2 {
            font-size: 0.9rem;
            color: var(--text-secondary);
            margin-bottom: 15px;
            text-transform: uppercase;
            letter-spacing: 1px;
        }

        .file-item {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 8px 10px;
            border-radius: 4px;
            font-size: 0.85rem;
            cursor: pointer;
        }

        .file-item:hover { background: var(--bg-tertiary); }
        .file-icon { opacity: 0.7; }

        .empty-state { text-align: center; padding: 40px; color: var(--text-secondary); }
        .empty-state-icon { font-size: 3rem; margin-bottom: 15px; }

        .agent-ceo { --agent-color: #4A90A4; }
        .agent-pm { --agent-color: #7B68EE; }
        .agent-ux { --agent-color: #FF6B6B; }
        .agent-ui { --agent-color: #4ECDC4; }
        .agent-security { --agent-color: #F39C12; }
        .agent-architect { --agent-color: #9B59B6; }
        .agent-senior_dev { --agent-color: #2ECC71; }
        .agent-junior_dev { --agent-color: #3498DB; }

        .message[data-from] { border-left-color: var(--agent-color, var(--border)); }
    </style>
</head>
<body>
    <header class="header">
        <h1>Agent House</h1>
        <div class="status-bar">
            <div class="status-item">
                <div class="status-dot" id="wsStatus"></div>
                <span id="wsStatusText">Connecting...</span>
            </div>
            <div class="status-item">
                <div class="status-dot" id="taskStatus"></div>
                <span id="taskStatusText">Idle</span>
            </div>
            <div class="status-item">
                <span id="messageCount">0 messages</span>
            </div>
        </div>
    </header>

    <div class="container">
        <aside class="sidebar">
            <h2>Agents</h2>
            <div id="agentsList"></div>
        </aside>

        <main class="main-content">
            <div class="task-input-area">
                <form class="task-form" id="taskForm">
                    <input type="text" id="projectInput" placeholder="Project name" autocomplete="off" style="max-width: 180px; flex: none;">
                    <input type="text" id="taskInput" placeholder="Enter a task for the team..." autocomplete="off">
                    <button type="submit" id="submitBtn">Start Task</button>
                </form>
            </div>

            <div class="messages-area" id="messagesArea">
                <div class="empty-state" id="emptyState">
                    <div class="empty-state-icon">&#128172;</div>
                    <p>No messages yet. Submit a task to get started!</p>
                </div>
            </div>
        </main>

        <aside class="files-panel">
            <h2>Project Files</h2>
            <div id="filesList">
                <div class="empty-state"><p>No files created yet</p></div>
            </div>
        </aside>
    </div>

    <script>
        let ws = null;
        let messages = [];
        let isTaskRunning = false;

        const agentEmojis = {
            'ceo': '&#128084;', 'pm': '&#128203;', 'ux': '&#127912;', 'ui': '&#128444;',
            'security': '&#128274;', 'architect': '&#127959;', 'senior_dev': '&#128104;&#8205;&#128187;', 'junior_dev': '&#128105;&#8205;&#128187;'
        };

        const wsStatus = document.getElementById('wsStatus');
        const wsStatusText = document.getElementById('wsStatusText');
        const taskStatus = document.getElementById('taskStatus');
        const taskStatusText = document.getElementById('taskStatusText');
        const messageCount = document.getElementById('messageCount');
        const agentsList = document.getElementById('agentsList');
        const messagesArea = document.getElementById('messagesArea');
        const emptyState = document.getElementById('emptyState');
        const filesList = document.getElementById('filesList');
        const taskForm = document.getElementById('taskForm');
        const projectInput = document.getElementById('projectInput');
        const taskInput = document.getElementById('taskInput');
        const submitBtn = document.getElementById('submitBtn');

        function connectWebSocket() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            ws = new WebSocket(protocol + '//' + window.location.host + '/ws');

            ws.onopen = () => {
                wsStatus.classList.add('connected');
                wsStatusText.textContent = 'Connected';
            };

            ws.onclose = () => {
                wsStatus.classList.remove('connected');
                wsStatusText.textContent = 'Disconnected';
                setTimeout(connectWebSocket, 3000);
            };

            ws.onerror = () => {
                wsStatus.classList.remove('connected');
                wsStatusText.textContent = 'Error';
            };

            ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                if (data.type === 'message') {
                    addMessage(data.message);
                }
            };
        }

        function addMessage(msg) {
            messages.push(msg);
            renderMessages();
            updateMessageCount();
            messagesArea.scrollTop = messagesArea.scrollHeight;
        }

        function renderMessages() {
            if (messages.length === 0) {
                emptyState.style.display = 'block';
                return;
            }
            emptyState.style.display = 'none';

            const html = messages.map(msg => {
                const time = new Date(msg.timestamp).toLocaleTimeString();
                const emoji = agentEmojis[msg.from] || '&#129302;';
                const content = escapeHtml(msg.content);
                const isLong = content.length > 500;

                return '<div class="message agent-' + msg.from + '" data-from="' + msg.from + '">' +
                    '<div class="message-header">' +
                    '<span class="message-from">' + emoji + ' ' + msg.from + ' &#8594; ' + msg.to + '</span>' +
                    '<span class="message-time">' + time + '</span>' +
                    '</div>' +
                    '<div class="message-content' + (isLong ? ' collapsed' : '') + '">' + content + '</div>' +
                    (isLong ? '<button class="expand-btn" onclick="toggleExpand(this)">Show more</button>' : '') +
                    '</div>';
            }).join('');

            messagesArea.innerHTML = html;
        }

        function toggleExpand(btn) {
            const content = btn.previousElementSibling;
            content.classList.toggle('collapsed');
            btn.textContent = content.classList.contains('collapsed') ? 'Show more' : 'Show less';
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        function updateMessageCount() {
            messageCount.textContent = messages.length + ' message' + (messages.length !== 1 ? 's' : '');
        }

        async function fetchAgents() {
            const res = await fetch('/api/agents');
            const data = await res.json();
            agentsList.innerHTML = data.agents.map(a =>
                '<div class="agent-card' + (a.active ? ' active' : '') + '" style="border-left-color: ' + a.color + '">' +
                '<h3>' + (agentEmojis[a.role] || '&#129302;') + ' ' + a.name + '</h3>' +
                '<div class="role">' + a.role + '</div>' +
                '</div>'
            ).join('');
        }

        async function fetchFiles() {
            let projectId = projectInput.value.trim().toLowerCase().replace(/[^a-z0-9-_]/g, '-') || 'default';
            const res = await fetch('/api/files?project=' + encodeURIComponent(projectId));
            const data = await res.json();
            if (data.files.length === 0) {
                filesList.innerHTML = '<div class="empty-state"><p>No files in ' + projectId + '</p></div>';
                return;
            }
            filesList.innerHTML = '<div style="margin-bottom:10px;font-size:0.8rem;color:var(--text-secondary);">Project: ' + projectId + '</div>' +
                data.files.map(f =>
                    '<div class="file-item">' +
                    '<span class="file-icon">' + (f.is_dir ? '&#128193;' : '&#128196;') + '</span>' +
                    '<span>' + f.path + '</span>' +
                    '</div>'
                ).join('');
        }

        async function fetchStatus() {
            try {
                const res = await fetch('/api/status');
                const data = await res.json();
                isTaskRunning = data.task_running;
                if (isTaskRunning) {
                    taskStatus.classList.add('running');
                    taskStatusText.textContent = 'Running';
                    submitBtn.disabled = true;
                    submitBtn.textContent = 'Running...';
                } else {
                    taskStatus.classList.remove('running');
                    taskStatusText.textContent = 'Idle';
                    submitBtn.disabled = false;
                    submitBtn.textContent = 'Start Task';
                }
            } catch (e) {
                console.error('Failed to fetch status:', e);
            }
        }

        async function fetchMessages() {
            try {
                const res = await fetch('/api/messages');
                const data = await res.json();
                messages = data.messages || [];
                renderMessages();
                updateMessageCount();
            } catch (e) {
                console.error('Failed to fetch messages:', e);
            }
        }

        async function submitTask(e) {
            e.preventDefault();
            const task = taskInput.value.trim();
            let projectId = projectInput.value.trim().toLowerCase().replace(/[^a-z0-9-_]/g, '-') || 'default';
            if (!task || isTaskRunning) return;
            try {
                await fetch('/api/task', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ task: task, project_id: projectId })
                });
                taskInput.value = '';
                fetchStatus();
                fetchFiles();
            } catch (e) {
                console.error('Failed to submit task:', e);
            }
        }

        connectWebSocket();
        fetchAgents();
        fetchFiles();
        fetchMessages();
        fetchStatus();

        setInterval(fetchStatus, 2000);
        setInterval(fetchFiles, 5000);

        taskForm.addEventListener('submit', submitTask);
        projectInput.addEventListener('change', fetchFiles);
        projectInput.addEventListener('blur', fetchFiles);
    </script>
</body>
</html>`
