package web

// agentActivityHTML contains the live agent activity dashboard
var agentActivityHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Agent House — Live Activity</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
  <style>
:root {
  --bg-void:       #080B14;
  --bg-base:       #0D1117;
  --bg-raised:     #161B27;
  --bg-overlay:    #1E2535;
  --bg-hover:      #252D3D;
  --border-subtle: #1B2130;
  --border-default:#2A3345;
  --border-active: #3D4F6E;
  --text-primary:  #E6EDF3;
  --text-secondary:#8B949E;
  --text-muted:    #484F58;
  --accent:        #3B82F6;
  --accent-glow:   rgba(59,130,246,0.15);
  --c-ceo:      #60B4D4;
  --c-pm:       #8B7EF8;
  --c-ux:       #F47F7F;
  --c-ui:       #4FD1C5;
  --c-security: #F6A623;
  --c-architect:#A87CF5;
  --c-sr-dev:   #34D478;
  --c-jr-dev:   #4BA3E3;
  --st-idle:     #3A4A5A;
  --st-working:  #3B82F6;
  --st-thinking: #8B5CF6;
  --st-complete: #22C55E;
  --st-error:    #EF4444;
  --font-ui:   'Inter', system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', monospace;
}
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: var(--font-ui); background: var(--bg-void); color: var(--text-primary); height: 100vh; overflow: hidden; -webkit-font-smoothing: antialiased; }

/* Layout */
.app { display: flex; flex-direction: column; height: 100vh; }
.header { display: flex; align-items: center; gap: 16px; padding: 0 20px; height: 52px; background: var(--bg-base); border-bottom: 1px solid var(--border-subtle); flex-shrink: 0; }
.logo { font-size: 14px; font-weight: 700; letter-spacing: -0.02em; }
.logo span { color: var(--accent); }
.header-right { margin-left: auto; display: flex; align-items: center; gap: 12px; }
.ws-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--st-idle); }
.ws-dot.live { background: var(--st-complete); box-shadow: 0 0 6px var(--st-complete); }
.nav-links { display: flex; gap: 8px; }
.nav-links a { color: var(--text-secondary); text-decoration: none; font-size: 12px; padding: 4px 10px; border-radius: 6px; border: 1px solid var(--border-subtle); }
.nav-links a:hover { color: var(--text-primary); border-color: var(--border-default); }
.cost-badge { font-family: var(--font-mono); font-size: 12px; color: var(--st-complete); background: rgba(34,197,94,0.1); padding: 4px 10px; border-radius: 6px; }

/* Main content */
.main { display: flex; flex: 1; overflow: hidden; }

/* Agent grid */
.agent-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; padding: 16px; width: 100%; max-width: 900px; align-content: start; }
.agent-grid.has-detail { max-width: 400px; grid-template-columns: repeat(2, 1fr); }

/* Agent card */
.agent-card { background: var(--bg-raised); border: 1px solid var(--border-subtle); border-radius: 10px; padding: 14px; cursor: pointer; transition: all 0.15s; position: relative; overflow: hidden; min-height: 130px; }
.agent-card:hover { border-color: var(--border-default); background: var(--bg-overlay); }
.agent-card.active { border-color: var(--border-active); }
.agent-card.selected { border-color: var(--accent); box-shadow: 0 0 0 1px var(--accent), 0 0 20px var(--accent-glow); }
.agent-card .indicator { position: absolute; top: 0; left: 0; right: 0; height: 3px; background: var(--st-idle); transition: background 0.3s; }
.agent-card.active .indicator { background: var(--st-working); }
.agent-card.thinking .indicator { background: var(--st-thinking); animation: pulse 1.5s infinite; }
.agent-card.complete .indicator { background: var(--st-complete); }
.agent-card.error .indicator { background: var(--st-error); }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.5; } }

.card-header { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; }
.card-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.card-name { font-size: 12px; font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.card-role { font-size: 10px; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 6px; }
.card-status { font-size: 11px; color: var(--text-secondary); margin-bottom: 8px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.card-status.working { color: var(--st-working); }
.card-status.thinking { color: var(--st-thinking); }
.card-metrics { display: flex; gap: 10px; font-family: var(--font-mono); font-size: 10px; color: var(--text-muted); }
.card-metrics .metric { display: flex; flex-direction: column; gap: 2px; }
.card-metrics .metric-val { color: var(--text-secondary); }
.card-activity { font-size: 10px; color: var(--text-muted); margin-top: 6px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; font-family: var(--font-mono); }

/* Detail panel */
.detail-panel { flex: 1; background: var(--bg-base); border-left: 1px solid var(--border-subtle); display: flex; flex-direction: column; overflow: hidden; }
.detail-header { display: flex; align-items: center; gap: 12px; padding: 14px 18px; border-bottom: 1px solid var(--border-subtle); flex-shrink: 0; }
.detail-header .agent-dot { width: 12px; height: 12px; border-radius: 50%; }
.detail-header h2 { font-size: 15px; font-weight: 600; flex: 1; }
.detail-header .close-btn { background: none; border: 1px solid var(--border-subtle); color: var(--text-secondary); width: 28px; height: 28px; border-radius: 6px; cursor: pointer; font-size: 14px; display: flex; align-items: center; justify-content: center; }
.detail-header .close-btn:hover { background: var(--bg-raised); color: var(--text-primary); }
.abort-btn { background: rgba(239,68,68,0.1); border: 1px solid rgba(239,68,68,0.3); color: #EF4444; padding: 4px 12px; border-radius: 6px; font-size: 11px; cursor: pointer; font-weight: 500; }
.abort-btn:hover { background: rgba(239,68,68,0.2); }

.detail-stats { display: flex; gap: 16px; padding: 10px 18px; border-bottom: 1px solid var(--border-subtle); font-size: 11px; flex-shrink: 0; }
.detail-stats .stat { display: flex; flex-direction: column; gap: 2px; }
.detail-stats .stat-label { color: var(--text-muted); font-size: 10px; text-transform: uppercase; }
.detail-stats .stat-val { color: var(--text-secondary); font-family: var(--font-mono); }

/* Event timeline */
.timeline { flex: 1; overflow-y: auto; padding: 12px 18px; }
.timeline::-webkit-scrollbar { width: 6px; }
.timeline::-webkit-scrollbar-track { background: transparent; }
.timeline::-webkit-scrollbar-thumb { background: var(--border-default); border-radius: 3px; }

.event { display: flex; gap: 10px; padding: 6px 0; border-bottom: 1px solid var(--border-subtle); font-size: 12px; }
.event:last-child { border-bottom: none; }
.event-icon { width: 20px; height: 20px; display: flex; align-items: center; justify-content: center; font-size: 12px; flex-shrink: 0; margin-top: 1px; }
.event-body { flex: 1; min-width: 0; }
.event-header { display: flex; align-items: center; gap: 8px; }
.event-type { font-weight: 500; font-size: 11px; }
.event-time { margin-left: auto; color: var(--text-muted); font-size: 10px; font-family: var(--font-mono); }
.event-content { color: var(--text-secondary); font-size: 11px; margin-top: 3px; line-height: 1.5; }
.event-content.mono { font-family: var(--font-mono); font-size: 10px; }
.event-content.collapsed { max-height: 3em; overflow: hidden; cursor: pointer; }
.event-content.collapsed:hover { color: var(--text-primary); }
.event-content.expanded { max-height: none; }

.event.text .event-type { color: var(--text-primary); }
.event.thinking .event-type { color: var(--st-thinking); }
.event.tool-use .event-type { color: var(--accent); }
.event.tool-result .event-type { color: var(--st-complete); }
.event.tool-result.error .event-type { color: var(--st-error); }
.event.turn-complete .event-type { color: var(--st-complete); }

/* Tool approval */
.approval-card { background: rgba(246,166,35,0.08); border: 1px solid rgba(246,166,35,0.25); border-radius: 8px; padding: 12px; margin: 8px 0; }
.approval-card .tool-name { font-weight: 600; font-size: 12px; color: #F6A623; margin-bottom: 6px; }
.approval-card .tool-input { font-family: var(--font-mono); font-size: 10px; color: var(--text-secondary); background: var(--bg-raised); padding: 8px; border-radius: 4px; margin-bottom: 8px; max-height: 120px; overflow: auto; white-space: pre-wrap; word-break: break-all; }
.approval-btns { display: flex; gap: 8px; }
.approval-btns button { padding: 5px 16px; border-radius: 6px; border: none; font-size: 11px; font-weight: 600; cursor: pointer; }
.btn-approve { background: var(--st-complete); color: #000; }
.btn-approve:hover { filter: brightness(1.1); }
.btn-deny { background: transparent; border: 1px solid var(--st-error) !important; color: var(--st-error); }
.btn-deny:hover { background: rgba(239,68,68,0.1); }

/* Empty state */
.empty { display: flex; align-items: center; justify-content: center; height: 100%; color: var(--text-muted); font-size: 13px; flex-direction: column; gap: 8px; }
.empty .icon { font-size: 32px; opacity: 0.4; }

/* Toast */
.toast { position: fixed; top: 60px; right: 16px; background: var(--bg-overlay); border: 1px solid var(--border-default); border-radius: 8px; padding: 10px 16px; font-size: 12px; z-index: 1000; animation: slideIn 0.3s ease; max-width: 300px; }
@keyframes slideIn { from { transform: translateX(100%); opacity: 0; } to { transform: translateX(0); opacity: 1; } }
  </style>
</head>
<body>
<div class="app">
  <!-- Header -->
  <div class="header">
    <div class="logo">Agent<span>House</span> <span style="color:var(--text-muted);font-weight:400">Live Activity</span></div>
    <div class="nav-links">
      <a href="/">Dashboard</a>
      <a href="/mission">Mission</a>
      <a href="/office">Office</a>
    </div>
    <div class="header-right">
      <div class="cost-badge" id="totalCost">$0.0000</div>
      <div class="ws-dot" id="wsDot"></div>
    </div>
  </div>

  <!-- Main -->
  <div class="main">
    <!-- Agent Grid -->
    <div class="agent-grid" id="agentGrid">
    </div>

    <!-- Detail Panel (hidden by default) -->
    <div class="detail-panel" id="detailPanel" style="display:none">
      <div class="detail-header">
        <div class="agent-dot" id="detailDot"></div>
        <h2 id="detailName"></h2>
        <button class="abort-btn" id="abortBtn" onclick="abortAgent()">Abort</button>
        <button class="close-btn" onclick="closeDetail()">&#x2715;</button>
      </div>
      <div class="detail-stats" id="detailStats"></div>
      <div class="timeline" id="timeline"></div>
    </div>
  </div>
</div>

<script>
// ══════════════════════════════════════════════
// STATE
// ══════════════════════════════════════════════
const AGENTS = {
  ceo:        { name: 'CEO',              color: '#60B4D4', icon: '&#x1F454;' },
  pm:         { name: 'Product Manager',  color: '#8B7EF8', icon: '&#x1F4CB;' },
  ux:         { name: 'UX Designer',      color: '#F47F7F', icon: '&#x1F3A8;' },
  ui:         { name: 'UI Designer',      color: '#4FD1C5', icon: '&#x1F58C;' },
  security:   { name: 'Security Expert',  color: '#F6A623', icon: '&#x1F6E1;' },
  architect:  { name: 'Architect',        color: '#A87CF5', icon: '&#x1F3D7;' },
  senior_dev: { name: 'Senior Dev',       color: '#34D478', icon: '&#x1F4BB;' },
  junior_dev: { name: 'Junior Dev',       color: '#4BA3E3', icon: '&#x2328;'  }
};

const TOOL_ICONS = {
  Read: '&#x1F4D6;', Write: '&#x1F4DD;', Edit: '&#x270F;', Bash: '&#x1F4BB;',
  Grep: '&#x1F50D;', Glob: '&#x1F4C2;', Agent: '&#x1F916;', default: '&#x1F527;'
};

const S = {
  ws: null,
  selected: null,          // currently selected agent role
  agents: {},              // role -> { state, events, stats, lastActivity }
  totalCost: 0,
  projectId: ''
};

// Init agent state
Object.keys(AGENTS).forEach(k => {
  S.agents[k] = { state: 'idle', events: [], stats: { tokens: 0, cost: 0, tools: 0, turns: 0 }, lastActivity: '', pendingTool: null };
});

// ══════════════════════════════════════════════
// RENDER
// ══════════════════════════════════════════════
function renderGrid() {
  const grid = document.getElementById('agentGrid');
  grid.className = S.selected ? 'agent-grid has-detail' : 'agent-grid';
  grid.innerHTML = Object.entries(AGENTS).map(([role, info]) => {
    const a = S.agents[role];
    const sel = S.selected === role ? ' selected' : '';
    const stClass = a.state === 'active' ? ' active' : a.state === 'thinking' ? ' thinking' : a.state === 'complete' ? ' complete' : a.state === 'error' ? ' error' : '';
    const statusText = a.state === 'idle' ? 'Idle' : a.state === 'active' ? 'Working...' : a.state === 'thinking' ? 'Thinking...' : a.state === 'complete' ? 'Complete' : a.state === 'error' ? 'Error' : a.state;
    const statusClass = (a.state === 'active' || a.state === 'thinking') ? a.state : '';
    return '<div class="agent-card'+stClass+sel+'" onclick="selectAgent(\''+role+'\')">' +
      '<div class="indicator"></div>' +
      '<div class="card-header"><div class="card-dot" style="background:'+info.color+'"></div><div class="card-name">'+info.name+'</div></div>' +
      '<div class="card-role">'+role.replace('_',' ')+'</div>' +
      '<div class="card-status '+statusClass+'">'+statusText+'</div>' +
      '<div class="card-metrics">' +
        '<div class="metric"><span>Tokens</span><span class="metric-val">'+fmtNum(a.stats.tokens)+'</span></div>' +
        '<div class="metric"><span>Cost</span><span class="metric-val">$'+a.stats.cost.toFixed(4)+'</span></div>' +
        '<div class="metric"><span>Tools</span><span class="metric-val">'+a.stats.tools+'</span></div>' +
      '</div>' +
      (a.lastActivity ? '<div class="card-activity">'+escHtml(a.lastActivity)+'</div>' : '') +
    '</div>';
  }).join('');
}

function renderDetail() {
  const panel = document.getElementById('detailPanel');
  if (!S.selected) { panel.style.display = 'none'; return; }
  panel.style.display = 'flex';

  const role = S.selected;
  const info = AGENTS[role];
  const a = S.agents[role];

  document.getElementById('detailDot').style.background = info.color;
  document.getElementById('detailName').textContent = info.name;
  document.getElementById('abortBtn').style.display = (a.state === 'active' || a.state === 'thinking') ? '' : 'none';

  document.getElementById('detailStats').innerHTML =
    '<div class="stat"><span class="stat-label">Status</span><span class="stat-val">'+a.state+'</span></div>' +
    '<div class="stat"><span class="stat-label">Tokens</span><span class="stat-val">'+fmtNum(a.stats.tokens)+'</span></div>' +
    '<div class="stat"><span class="stat-label">Cost</span><span class="stat-val">$'+a.stats.cost.toFixed(4)+'</span></div>' +
    '<div class="stat"><span class="stat-label">Tools</span><span class="stat-val">'+a.stats.tools+'</span></div>' +
    '<div class="stat"><span class="stat-label">Turns</span><span class="stat-val">'+a.stats.turns+'</span></div>';

  renderTimeline();
}

function renderTimeline() {
  const tl = document.getElementById('timeline');
  const role = S.selected;
  if (!role) return;

  const events = S.agents[role].events;
  if (events.length === 0) {
    tl.innerHTML = '<div class="empty"><div class="icon">&#x1F4E1;</div>No activity yet</div>';
    return;
  }

  // Render last 200 events
  const slice = events.slice(-200);
  tl.innerHTML = slice.map(ev => renderEvent(ev)).join('');

  // Auto-scroll to bottom
  tl.scrollTop = tl.scrollHeight;
}

function renderEvent(ev) {
  const time = new Date(ev.timestamp).toLocaleTimeString('en-US', {hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit'});

  switch (ev.type) {
    case 'text_delta':
      return '<div class="event text"><div class="event-icon">&#x1F4AC;</div><div class="event-body">' +
        '<div class="event-header"><span class="event-type">Text</span><span class="event-time">'+time+'</span></div>' +
        '<div class="event-content collapsed" onclick="this.classList.toggle(\'collapsed\');this.classList.toggle(\'expanded\')">'+escHtml(ev.content)+'</div></div></div>';

    case 'thinking_delta':
      return '<div class="event thinking"><div class="event-icon">&#x1F4AD;</div><div class="event-body">' +
        '<div class="event-header"><span class="event-type">Thinking</span><span class="event-time">'+time+'</span></div>' +
        '<div class="event-content collapsed" onclick="this.classList.toggle(\'collapsed\');this.classList.toggle(\'expanded\')">'+escHtml(ev.content)+'</div></div></div>';

    case 'tool_use': {
      const icon = TOOL_ICONS[ev.tool_name] || TOOL_ICONS.default;
      let inputPreview = ev.input || '';
      try { const p = JSON.parse(inputPreview); inputPreview = formatToolInput(ev.tool_name, p); } catch(e) {}

      // Check if this is a pending approval
      const a = S.agents[ev.agent_role];
      if (a && a.pendingTool && a.pendingTool.tool_use_id === ev.tool_use_id) {
        return '<div class="event tool-use"><div class="event-icon">'+icon+'</div><div class="event-body">' +
          '<div class="event-header"><span class="event-type">'+ev.tool_name+'</span><span class="event-time">'+time+'</span></div>' +
          '<div class="approval-card"><div class="tool-name">Waiting for approval: '+ev.tool_name+'</div>' +
          '<div class="tool-input">'+escHtml(inputPreview)+'</div>' +
          '<div class="approval-btns">' +
          '<button class="btn-approve" onclick="approveTool(\''+ev.agent_role+'\',\''+ev.tool_use_id+'\')">Approve</button>' +
          '<button class="btn-deny" onclick="denyTool(\''+ev.agent_role+'\',\''+ev.tool_use_id+'\')">Deny</button>' +
          '</div></div></div></div>';
      }

      return '<div class="event tool-use"><div class="event-icon">'+icon+'</div><div class="event-body">' +
        '<div class="event-header"><span class="event-type">'+ev.tool_name+'</span><span class="event-time">'+time+'</span></div>' +
        '<div class="event-content mono collapsed" onclick="this.classList.toggle(\'collapsed\');this.classList.toggle(\'expanded\')">'+escHtml(inputPreview)+'</div></div></div>';
    }

    case 'tool_result': {
      const errClass = ev.is_error ? ' error' : '';
      const output = ev.output || '';
      const preview = output.length > 500 ? output.substring(0, 500) + '...' : output;
      return '<div class="event tool-result'+errClass+'"><div class="event-icon">'+(ev.is_error?'&#x274C;':'&#x2705;')+'</div><div class="event-body">' +
        '<div class="event-header"><span class="event-type">Result'+(ev.is_error?' (Error)':'')+'</span><span class="event-time">'+time+'</span></div>' +
        '<div class="event-content mono collapsed" onclick="this.classList.toggle(\'collapsed\');this.classList.toggle(\'expanded\')">'+escHtml(preview)+'</div></div></div>';
    }

    case 'turn_complete':
      return '<div class="event turn-complete"><div class="event-icon">&#x2705;</div><div class="event-body">' +
        '<div class="event-header"><span class="event-type">Turn Complete</span><span class="event-time">'+time+'</span></div>' +
        '<div class="event-content">Tokens: '+fmtNum(ev.input_tokens||0)+' in / '+fmtNum(ev.output_tokens||0)+' out &middot; Cost: $'+(ev.cost_usd||0).toFixed(4)+'</div></div></div>';

    case 'session_meta':
      return '<div class="event"><div class="event-icon">&#x2699;</div><div class="event-body">' +
        '<div class="event-header"><span class="event-type">Session Started</span><span class="event-time">'+time+'</span></div>' +
        '<div class="event-content">Model: '+(ev.model||'unknown')+'</div></div></div>';

    case 'error':
      return '<div class="event tool-result error"><div class="event-icon">&#x26A0;</div><div class="event-body">' +
        '<div class="event-header"><span class="event-type">Error</span><span class="event-time">'+time+'</span></div>' +
        '<div class="event-content">'+escHtml(ev.content||'Unknown error')+'</div></div></div>';

    default:
      return '';
  }
}

function formatToolInput(toolName, input) {
  switch (toolName) {
    case 'Read': return input.file_path || JSON.stringify(input);
    case 'Write': return (input.file_path || '?') + ' (' + ((input.content||'').length) + ' chars)';
    case 'Edit': return (input.file_path || '?') + ' (edit)';
    case 'Bash': return '$ ' + (input.command || JSON.stringify(input));
    case 'Grep': return 'grep "' + (input.pattern||'') + '" ' + (input.path||'.');
    case 'Glob': return 'glob "' + (input.pattern||'') + '"';
    default: return JSON.stringify(input).substring(0, 200);
  }
}

// ══════════════════════════════════════════════
// ACTIONS
// ══════════════════════════════════════════════
function selectAgent(role) {
  S.selected = S.selected === role ? null : role;
  renderGrid();
  renderDetail();
}

function closeDetail() {
  S.selected = null;
  renderGrid();
  renderDetail();
}

function abortAgent() {
  if (!S.selected) return;
  if (S.ws && S.ws.readyState === WebSocket.OPEN) {
    S.ws.send(JSON.stringify({ type: 'abort_agent', project_id: S.projectId, agent_role: S.selected }));
  }
  // Also via REST
  fetch('/api/sessions/' + encodeURIComponent(S.projectId) + '/' + S.selected + '/abort', { method: 'POST' }).catch(console.error);
}

function approveTool(agentRole, toolUseId) {
  if (S.ws && S.ws.readyState === WebSocket.OPEN) {
    S.ws.send(JSON.stringify({ type: 'approve_agent_tool', project_id: S.projectId, agent_role: agentRole, tool_use_id: toolUseId }));
  }
  const a = S.agents[agentRole];
  if (a) a.pendingTool = null;
  if (S.selected === agentRole) renderTimeline();
}

function denyTool(agentRole, toolUseId) {
  if (S.ws && S.ws.readyState === WebSocket.OPEN) {
    S.ws.send(JSON.stringify({ type: 'deny_agent_tool', project_id: S.projectId, agent_role: agentRole, tool_use_id: toolUseId }));
  }
  const a = S.agents[agentRole];
  if (a) a.pendingTool = null;
  if (S.selected === agentRole) renderTimeline();
}

// ══════════════════════════════════════════════
// WEBSOCKET
// ══════════════════════════════════════════════
function connectWS() {
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
  S.ws = new WebSocket(proto + '//' + location.host + '/ws');
  S.ws.onopen = () => { document.getElementById('wsDot').className = 'ws-dot live'; };
  S.ws.onclose = () => { document.getElementById('wsDot').className = 'ws-dot'; setTimeout(connectWS, 3000); };
  S.ws.onerror = () => { S.ws.close(); };
  S.ws.onmessage = (e) => {
    try { handleWS(JSON.parse(e.data)); } catch (x) { console.error('[WS]', x); }
  };
}

function handleWS(data) {
  // Handle agent session events
  if (data.type === 'agent_event' && data.event) {
    handleAgentEvent(data.event);
    return;
  }

  // Handle lifecycle events (phase/task start/complete)
  if (data.type === 'message' && data.message && data.message.type === 'lifecycle') {
    handleLifecycleEvent(data.message);
    return;
  }

  // Handle legacy task events (update agent state)
  if (data.type === 'agent_task_event' && data.event) {
    const evt = data.event;
    const role = normRole(evt.agent_role);
    if (!role || !S.agents[role]) return;
    const a = S.agents[role];
    switch (evt.event_type) {
      case 'task_started': a.state = 'active'; break;
      case 'task_completed': a.state = 'complete'; setTimeout(() => { a.state = 'idle'; renderGrid(); }, 5000); break;
      case 'task_failed': a.state = 'error'; break;
    }
    renderGrid();
  }
}

function handleLifecycleEvent(msg) {
  const extra = msg.metadata?.extra || {};
  const eventType = extra.event_type;

  if (eventType === 'task_completed') {
    const fc = extra.files_count || 0;
    const turns = extra.turns || 0;
    // Reset ALL agents to idle
    Object.keys(S.agents).forEach(k => {
      S.agents[k].state = 'idle';
      S.agents[k].lastActivity = '';
      S.agents[k].pendingTool = null;
    });
    renderGrid();
    renderDetail();
    // Show completion toast with project link
    if (S.projectId) {
      showCompletionBanner(turns, fc, S.projectId);
    }
  }

  if (eventType === 'phase_started') {
    console.log('[LIFECYCLE] Phase started:', extra.phase_name);
  }

  if (eventType === 'phase_completed') {
    console.log('[LIFECYCLE] Phase completed:', extra.phase_name, extra.duration_ms + 'ms');
  }
}

function showCompletionBanner(turns, files, projectId) {
  // Remove existing banner
  const old = document.getElementById('completionBanner');
  if (old) old.remove();

  const banner = document.createElement('div');
  banner.id = 'completionBanner';
  banner.style.cssText = 'position:fixed;bottom:20px;left:50%;transform:translateX(-50%);background:var(--bg-overlay);border:1px solid var(--st-complete);border-radius:10px;padding:14px 24px;z-index:1000;display:flex;align-items:center;gap:16px;font-size:13px;box-shadow:0 4px 20px rgba(0,0,0,0.4)';
  banner.innerHTML = '<span style="font-size:18px">&#x2705;</span>' +
    '<span>Task complete: <strong>' + turns + ' turns</strong>, <strong>' + files + ' files</strong></span>' +
    '<a href="/project/' + encodeURIComponent(projectId) + '" style="color:var(--accent);font-weight:600;white-space:nowrap">View Project &#x2192;</a>' +
    '<button onclick="this.parentElement.remove()" style="background:none;border:none;color:var(--text-muted);cursor:pointer;font-size:16px;margin-left:8px">&#x2715;</button>';
  document.body.appendChild(banner);
}

function handleAgentEvent(ev) {
  const role = normRole(ev.agent_role);
  if (!role || !S.agents[role]) return;

  const a = S.agents[role];

  // Capture project ID
  if (ev.project_id && !S.projectId) S.projectId = ev.project_id;

  // Add event to agent's timeline
  a.events.push(ev);
  if (a.events.length > 1000) a.events = a.events.slice(-500); // Keep last 500

  // Update agent state based on event type
  switch (ev.type) {
    case 'text_delta':
      a.state = 'active';
      a.lastActivity = ev.content ? ev.content.substring(0, 80) : '';
      break;
    case 'thinking_delta':
      a.state = 'thinking';
      a.lastActivity = 'Thinking...';
      break;
    case 'tool_use':
      a.state = 'active';
      a.stats.tools++;
      a.lastActivity = ev.tool_name + ': ' + (formatToolInput(ev.tool_name, tryParse(ev.input)) || '').substring(0, 60);
      // Check if this might need approval (permission_request)
      break;
    case 'tool_result':
      a.state = 'active';
      break;
    case 'turn_complete':
      a.state = 'complete';
      a.stats.tokens += (ev.input_tokens || 0) + (ev.output_tokens || 0);
      a.stats.cost += ev.cost_usd || 0;
      a.stats.turns++;
      S.totalCost += ev.cost_usd || 0;
      document.getElementById('totalCost').textContent = '$' + S.totalCost.toFixed(4);
      a.lastActivity = 'Turn complete';
      // Reset to idle after a moment
      setTimeout(() => { if (a.state === 'complete') { a.state = 'idle'; a.lastActivity = ''; renderGrid(); } }, 5000);
      break;
    case 'session_meta':
      a.state = 'active';
      a.lastActivity = 'Session started (' + (ev.model || 'unknown') + ')';
      break;
    case 'error':
      a.state = 'error';
      a.lastActivity = 'Error: ' + (ev.content || '').substring(0, 60);
      setTimeout(() => { if (a.state === 'error') { a.state = 'idle'; renderGrid(); } }, 10000);
      break;
  }

  renderGrid();
  if (S.selected === role) renderDetail();
}

// ══════════════════════════════════════════════
// HELPERS
// ══════════════════════════════════════════════
function normRole(r) {
  if (!r) return null;
  r = r.toLowerCase().replace(/[\s-]+/g, '_');
  if (AGENTS[r]) return r;
  // Try common aliases
  const map = { 'ceo': 'ceo', 'pm': 'pm', 'product_manager': 'pm', 'ux_designer': 'ux', 'ui_designer': 'ui',
    'security_expert': 'security', 'security': 'security', 'architect': 'architect',
    'senior_developer': 'senior_dev', 'senior_dev': 'senior_dev', 'sr_dev': 'senior_dev',
    'junior_developer': 'junior_dev', 'junior_dev': 'junior_dev', 'jr_dev': 'junior_dev' };
  return map[r] || null;
}

function fmtNum(n) {
  if (n >= 1000000) return (n/1000000).toFixed(1) + 'M';
  if (n >= 1000) return (n/1000).toFixed(1) + 'K';
  return String(n);
}

function escHtml(s) {
  if (!s) return '';
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

function tryParse(s) {
  try { return JSON.parse(s); } catch(e) { return s; }
}

// ══════════════════════════════════════════════
// INIT
// ══════════════════════════════════════════════
renderGrid();
connectWS();

// Fetch existing sessions to get initial state
fetch('/api/sessions').then(r => r.json()).then(data => {
  if (!data.sessions) return;
  data.sessions.forEach(sess => {
    const role = normRole(sess.agent_role);
    if (!role || !S.agents[role]) return;
    const a = S.agents[role];
    if (sess.state === 'active') a.state = 'active';
    if (sess.stats) {
      a.stats.tokens = (sess.stats.total_input_tokens||0) + (sess.stats.total_output_tokens||0);
      a.stats.cost = sess.stats.total_cost_usd || 0;
      a.stats.tools = sess.stats.total_tool_calls || 0;
      a.stats.turns = sess.stats.total_turns || 0;
      S.totalCost += a.stats.cost;
    }
    if (sess.project_id && !S.projectId) S.projectId = sess.project_id;
  });
  document.getElementById('totalCost').textContent = '$' + S.totalCost.toFixed(4);
  renderGrid();
}).catch(() => {});
</script>
</body>
</html>
`
