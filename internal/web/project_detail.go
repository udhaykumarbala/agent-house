package web

var projectDetailHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Agent House — Project Detail</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
  <style>
:root {
  --bg-void:#080B14;--bg-base:#0D1117;--bg-raised:#161B27;--bg-overlay:#1E2535;--bg-hover:#252D3D;
  --border-subtle:#1B2130;--border-default:#2A3345;--border-active:#3D4F6E;
  --text-primary:#E6EDF3;--text-secondary:#8B949E;--text-muted:#484F58;
  --accent:#3B82F6;--st-complete:#22C55E;--st-error:#EF4444;
  --font-ui:'Inter',system-ui,sans-serif;--font-mono:'JetBrains Mono',monospace;
}
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
body{font-family:var(--font-ui);background:var(--bg-void);color:var(--text-primary);min-height:100vh;-webkit-font-smoothing:antialiased}
a{color:var(--accent);text-decoration:none}a:hover{text-decoration:underline}

.header{display:flex;align-items:center;gap:16px;padding:0 24px;height:52px;background:var(--bg-base);border-bottom:1px solid var(--border-subtle)}
.logo{font-size:14px;font-weight:700}.logo span{color:var(--accent)}
.nav{display:flex;gap:8px;margin-left:auto}
.nav a{color:var(--text-secondary);font-size:12px;padding:4px 10px;border-radius:6px;border:1px solid var(--border-subtle);text-decoration:none}
.nav a:hover{color:var(--text-primary);border-color:var(--border-default)}

.container{max-width:1100px;margin:0 auto;padding:24px}

.project-header{display:flex;align-items:center;gap:16px;margin-bottom:24px}
.project-header h1{font-size:22px;font-weight:700}
.project-header .badge{font-size:11px;padding:3px 10px;border-radius:12px;background:rgba(59,130,246,0.15);color:var(--accent)}

/* Grid layout */
.grid{display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-bottom:24px}
.grid-full{grid-column:1/-1}

/* Card */
.card{background:var(--bg-raised);border:1px solid var(--border-subtle);border-radius:10px;overflow:hidden}
.card-header{display:flex;align-items:center;justify-content:space-between;padding:12px 16px;border-bottom:1px solid var(--border-subtle);font-size:13px;font-weight:600}
.card-header .count{font-weight:400;color:var(--text-muted);font-size:12px}
.card-body{padding:16px;max-height:400px;overflow-y:auto}
.card-body::-webkit-scrollbar{width:5px}
.card-body::-webkit-scrollbar-thumb{background:var(--border-default);border-radius:3px}

/* File list */
.file-item{display:flex;align-items:center;gap:10px;padding:6px 0;border-bottom:1px solid var(--border-subtle);font-size:12px;cursor:pointer}
.file-item:last-child{border-bottom:none}
.file-item:hover{color:var(--accent)}
.file-icon{font-size:14px;flex-shrink:0;width:20px;text-align:center}
.file-path{flex:1;font-family:var(--font-mono);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.file-size{color:var(--text-muted);font-family:var(--font-mono);font-size:11px;flex-shrink:0}

/* Task history */
.task-item{padding:10px 0;border-bottom:1px solid var(--border-subtle)}
.task-item:last-child{border-bottom:none}
.task-title{font-size:12px;font-weight:500;margin-bottom:4px}
.task-meta{display:flex;gap:12px;font-size:11px;color:var(--text-muted)}
.task-meta .tag{padding:1px 6px;border-radius:4px;font-size:10px}
.task-meta .tag.success{background:rgba(34,197,94,0.15);color:var(--st-complete)}
.task-meta .tag.error{background:rgba(239,68,68,0.15);color:var(--st-error)}

/* Agent metrics */
.agent-row{display:flex;align-items:center;gap:12px;padding:8px 0;border-bottom:1px solid var(--border-subtle);font-size:12px}
.agent-row:last-child{border-bottom:none}
.agent-dot{width:10px;height:10px;border-radius:50%;flex-shrink:0}
.agent-name{width:100px;font-weight:500}
.agent-stat{font-family:var(--font-mono);font-size:11px;color:var(--text-secondary);min-width:60px}
.agent-tools{font-family:var(--font-mono);font-size:10px;color:var(--text-muted);flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}

/* File viewer */
.viewer{position:fixed;inset:0;z-index:100;display:none;background:rgba(0,0,0,0.6);backdrop-filter:blur(4px)}
.viewer.open{display:flex;align-items:center;justify-content:center}
.viewer-box{background:var(--bg-base);border:1px solid var(--border-default);border-radius:12px;width:90%;max-width:800px;max-height:85vh;display:flex;flex-direction:column}
.viewer-header{display:flex;align-items:center;padding:12px 16px;border-bottom:1px solid var(--border-subtle);font-size:13px;font-weight:600;gap:10px}
.viewer-header .close{margin-left:auto;background:none;border:1px solid var(--border-subtle);color:var(--text-secondary);width:28px;height:28px;border-radius:6px;cursor:pointer;font-size:14px;display:flex;align-items:center;justify-content:center}
.viewer-header .close:hover{background:var(--bg-raised)}
.viewer-content{flex:1;overflow:auto;padding:16px}
.viewer-content pre{font-family:var(--font-mono);font-size:12px;line-height:1.6;white-space:pre-wrap;word-break:break-all;color:var(--text-secondary)}

/* Continue task */
.continue-section{margin-top:24px}
.continue-input{display:flex;gap:10px}
.continue-input textarea{flex:1;background:var(--bg-raised);border:1px solid var(--border-default);color:var(--text-primary);padding:10px 14px;border-radius:8px;font-family:var(--font-ui);font-size:13px;resize:vertical;min-height:44px}
.continue-input textarea:focus{outline:none;border-color:var(--accent)}
.continue-input button{background:var(--accent);color:#fff;border:none;padding:10px 20px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;white-space:nowrap}
.continue-input button:hover{filter:brightness(1.1)}
.continue-input button:disabled{opacity:0.5;cursor:not-allowed}

/* Subtasks */
.subtask{display:flex;align-items:center;gap:10px;padding:6px 0;border-bottom:1px solid var(--border-subtle);font-size:12px}
.subtask:last-child{border-bottom:none}
.subtask-status{width:8px;height:8px;border-radius:50%;flex-shrink:0}
.subtask-title{flex:1}
.subtask-agents{font-size:10px;color:var(--text-muted)}

.empty-msg{color:var(--text-muted);font-size:12px;text-align:center;padding:20px}
  </style>
</head>
<body>
<div class="header">
  <div class="logo">Agent<span>House</span></div>
  <div class="nav">
    <a href="/">Dashboard</a>
    <a href="/mission">Mission</a>
    <a href="/live">Live</a>
  </div>
</div>

<div class="container">
  <div class="project-header">
    <h1 id="projectName">Loading...</h1>
    <span class="badge" id="fileCount"></span>
  </div>

  <div class="grid">
    <!-- Files -->
    <div class="card">
      <div class="card-header">Project Files <span class="count" id="filesCount"></span></div>
      <div class="card-body" id="fileList"></div>
    </div>

    <!-- Agent Metrics -->
    <div class="card">
      <div class="card-header">Agent Activity</div>
      <div class="card-body" id="agentMetrics"></div>
    </div>

    <!-- Task History -->
    <div class="card">
      <div class="card-header">Task History <span class="count" id="taskCount"></span></div>
      <div class="card-body" id="taskHistory"></div>
    </div>

    <!-- Development Plan / Subtasks -->
    <div class="card">
      <div class="card-header">Development Plan</div>
      <div class="card-body" id="devPlan"></div>
    </div>
  </div>

  <!-- Continue building -->
  <div class="continue-section">
    <div class="card">
      <div class="card-header">Continue Building</div>
      <div class="card-body">
        <p style="font-size:12px;color:var(--text-secondary);margin-bottom:12px">Submit a follow-up task to continue developing this project.</p>
        <div class="continue-input">
          <textarea id="continueTask" placeholder="Describe what to change or add..."></textarea>
          <button id="continueBtn" onclick="submitContinuation()">Continue</button>
        </div>
        <div id="continueStatus" style="font-size:11px;color:var(--text-muted);margin-top:8px"></div>
      </div>
    </div>
  </div>
</div>

<!-- File Viewer Modal -->
<div class="viewer" id="viewer" onclick="if(event.target===this)closeViewer()">
  <div class="viewer-box">
    <div class="viewer-header">
      <span id="viewerTitle">File</span>
      <button class="close" onclick="closeViewer()">&#x2715;</button>
    </div>
    <div class="viewer-content"><pre id="viewerContent"></pre></div>
  </div>
</div>

<script>
const AGENT_COLORS = {
  ceo:'#60B4D4',pm:'#8B7EF8',ux:'#F47F7F',ui:'#4FD1C5',
  security:'#F6A623',architect:'#A87CF5',senior_dev:'#34D478',junior_dev:'#4BA3E3'
};
const FILE_ICONS = {
  html:'&#x1F310;',css:'&#x1F3A8;',js:'&#x26A1;',ts:'&#x1F535;',json:'&#x1F4CB;',
  md:'&#x1F4DD;',go:'&#x1F439;',py:'&#x1F40D;',yaml:'&#x2699;',yml:'&#x2699;',
  default:'&#x1F4C4;'
};

// Extract project ID from URL: /project/{id}
const PROJECT_ID = window.location.pathname.split('/project/')[1] || '';

async function loadProject() {
  if (!PROJECT_ID) {
    document.getElementById('projectName').textContent = 'No project selected';
    return;
  }

  document.getElementById('projectName').textContent = PROJECT_ID;

  try {
    const res = await fetch('/api/projects/' + encodeURIComponent(PROJECT_ID));
    const data = await res.json();
    renderFiles(data.files || []);
    renderTasks(data.tasks || []);
    renderSessions(data.sessions);
    renderDevPlan(data.dev_plan);
    document.getElementById('fileCount').textContent = (data.file_count || 0) + ' files';
  } catch (e) {
    console.error('Failed to load project:', e);
  }
}

function renderFiles(files) {
  const el = document.getElementById('fileList');
  document.getElementById('filesCount').textContent = files.length + ' files';
  if (!files.length) { el.innerHTML = '<div class="empty-msg">No files yet</div>'; return; }
  el.innerHTML = files.map(f => {
    const ext = (f.path || '').split('.').pop().toLowerCase();
    const icon = FILE_ICONS[ext] || FILE_ICONS.default;
    const size = f.size < 1024 ? f.size + ' B' : (f.size / 1024).toFixed(1) + ' KB';
    return '<div class="file-item" onclick="viewFile(\'' + esc(f.path) + '\')">'+
      '<span class="file-icon">' + icon + '</span>' +
      '<span class="file-path">' + esc(f.path) + '</span>' +
      '<span class="file-size">' + size + '</span></div>';
  }).join('');
}

function renderTasks(tasks) {
  const el = document.getElementById('taskHistory');
  document.getElementById('taskCount').textContent = tasks.length + ' tasks';
  if (!tasks.length) { el.innerHTML = '<div class="empty-msg">No tasks yet</div>'; return; }
  // tasks might be array of metadata objects
  el.innerHTML = tasks.map(t => {
    const status = t.status || t.Status || 'unknown';
    const isOk = status === 'completed' || status === 'success';
    const task = t.task || t.Task || 'Task';
    const dur = t.duration || t.Duration || '';
    const files = t.files_created || t.FilesCreated || [];
    return '<div class="task-item">' +
      '<div class="task-title">' + esc(task.substring(0, 120)) + '</div>' +
      '<div class="task-meta">' +
      '<span class="tag ' + (isOk ? 'success' : 'error') + '">' + status + '</span>' +
      (dur ? '<span>' + dur + '</span>' : '') +
      '<span>' + files.length + ' files</span>' +
      '</div></div>';
  }).join('');
}

function renderSessions(metrics) {
  const el = document.getElementById('agentMetrics');
  if (!metrics || !metrics.agent_metrics) {
    el.innerHTML = '<div class="empty-msg">No session data — run a task first</div>';
    return;
  }
  const agents = metrics.agent_metrics;
  const entries = Object.entries(agents).sort((a, b) => a[0].localeCompare(b[0]));
  el.innerHTML = entries.map(([role, m]) => {
    const color = AGENT_COLORS[role] || '#888';
    const toolStr = Object.entries(m.tools || {}).map(([k, v]) => k + ':' + v).join(', ') || '-';
    return '<div class="agent-row">' +
      '<div class="agent-dot" style="background:' + color + '"></div>' +
      '<div class="agent-name">' + (m.name || role) + '</div>' +
      '<div class="agent-stat">' + (m.turns || 0) + ' turns</div>' +
      '<div class="agent-stat">' + (m.tool_calls || 0) + ' tools</div>' +
      '<div class="agent-stat">$' + (m.cost_usd || 0).toFixed(3) + '</div>' +
      '<div class="agent-tools">' + esc(toolStr) + '</div></div>';
  }).join('');
  // Total row
  const totalCost = metrics.total_cost || '0';
  el.innerHTML += '<div class="agent-row" style="font-weight:600;border-top:2px solid var(--border-default);margin-top:4px;padding-top:10px">' +
    '<div class="agent-dot" style="background:transparent"></div>' +
    '<div class="agent-name">Total</div>' +
    '<div class="agent-stat">' + entries.length + ' agents</div>' +
    '<div class="agent-stat">' + (metrics.total_tokens || 0) + ' tok</div>' +
    '<div class="agent-stat" style="color:var(--st-complete)">$' + totalCost + '</div>' +
    '<div class="agent-tools"></div></div>';
}

function renderDevPlan(plan) {
  const el = document.getElementById('devPlan');
  if (!plan || !plan.phases) {
    el.innerHTML = '<div class="empty-msg">No development plan</div>';
    return;
  }
  el.innerHTML = plan.phases.map(phase => {
    const statusColors = {pending:'var(--text-muted)',in_progress:'var(--accent)',completed:'var(--st-complete)',needs_revision:'#F6A623'};
    let html = '<div style="margin-bottom:12px">' +
      '<div style="font-size:12px;font-weight:600;margin-bottom:6px">' +
      '<span style="display:inline-block;width:8px;height:8px;border-radius:50%;background:'+(statusColors[phase.status]||'var(--text-muted)')+';margin-right:6px"></span>' +
      'Phase ' + phase.index + ': ' + esc(phase.name) + ' <span style="color:var(--text-muted);font-weight:400">(' + (phase.status||'pending') + ')</span></div>';
    if (phase.subtasks && phase.subtasks.length) {
      html += phase.subtasks.map(st => {
        const stColor = st.status === 'completed' ? 'var(--st-complete)' : st.status === 'in_progress' ? 'var(--accent)' : 'var(--text-muted)';
        return '<div class="subtask">' +
          '<div class="subtask-status" style="background:'+stColor+'"></div>' +
          '<div class="subtask-title">' + esc(st.title || 'Subtask') + '</div>' +
          '<div class="subtask-agents">' + (st.assigned_agents || []).join(', ') + '</div></div>';
      }).join('');
    }
    html += '</div>';
    return html;
  }).join('');
}

async function viewFile(relPath) {
  document.getElementById('viewerTitle').textContent = relPath;
  document.getElementById('viewerContent').textContent = 'Loading...';
  document.getElementById('viewer').classList.add('open');
  try {
    const res = await fetch('/api/file-content?path=' + encodeURIComponent(PROJECT_ID + '/' + relPath));
    if (res.ok) {
      const data = await res.json();
      document.getElementById('viewerContent').textContent = data.content || '(empty)';
    } else {
      document.getElementById('viewerContent').textContent = 'Failed to load: ' + res.status;
    }
  } catch (e) {
    document.getElementById('viewerContent').textContent = 'Error: ' + e.message;
  }
}

function closeViewer() {
  document.getElementById('viewer').classList.remove('open');
}

async function submitContinuation() {
  const textarea = document.getElementById('continueTask');
  const task = textarea.value.trim();
  if (!task) return;
  const btn = document.getElementById('continueBtn');
  btn.disabled = true;
  document.getElementById('continueStatus').textContent = 'Submitting...';
  try {
    const res = await fetch('/api/task', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({task: task, project_id: PROJECT_ID})
    });
    const data = await res.json();
    if (data.success) {
      textarea.value = '';
      document.getElementById('continueStatus').innerHTML =
        'Task started! <a href="/live">Watch live</a> | <a href="/mission">Mission Control</a>';
    } else {
      document.getElementById('continueStatus').textContent = 'Error: ' + (data.error || 'Failed');
      btn.disabled = false;
    }
  } catch (e) {
    document.getElementById('continueStatus').textContent = 'Error: ' + e.message;
    btn.disabled = false;
  }
}

function esc(s) { return (s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;').replace(/'/g,'&#39;'); }

// Handle keyboard
document.addEventListener('keydown', e => {
  if (e.key === 'Escape') closeViewer();
});

// Load on start
loadProject();
</script>
</body>
</html>
`
