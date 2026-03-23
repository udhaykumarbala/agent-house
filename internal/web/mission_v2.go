package web

var missionV2HTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Agent House</title>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
:root {
  --void:#0B0E17;--base:#111827;--elevated:#1A2332;--hover:#1F2B3D;
  --glass:rgba(17,24,39,0.75);--glass-border:rgba(255,255,255,0.04);
  --border-subtle:rgba(255,255,255,0.04);--border:rgba(255,255,255,0.08);--border-active:rgba(129,140,248,0.3);
  --text:#F3F4F6;--text2:#9CA3AF;--text3:#4B5563;
  --accent:#818CF8;--glow:rgba(129,140,248,0.12);
  --green:#34D399;--yellow:#FBBF24;--red:#F87171;--idle:#374151;
  --font:'Inter',system-ui,sans-serif;--mono:'JetBrains Mono',monospace;
  --sidebar-w:220px;--header-h:48px;--status-h:32px;--tab-h:36px;
}
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
body{font-family:var(--font);background:var(--void);color:var(--text);height:100vh;overflow:hidden;-webkit-font-smoothing:antialiased}
::selection{background:var(--accent);color:var(--void)}
::-webkit-scrollbar{width:5px}::-webkit-scrollbar-track{background:transparent}::-webkit-scrollbar-thumb{background:var(--border);border-radius:3px}

/* ═══ LAYOUT ═══ */
.app{display:grid;height:100vh;grid-template-rows:var(--header-h) var(--tab-h) 1fr var(--status-h);grid-template-columns:var(--sidebar-w) 1fr;grid-template-areas:"hdr hdr" "side tabs" "side main" "status status"}

/* ═══ HEADER ═══ */
.hdr{grid-area:hdr;display:flex;align-items:center;gap:14px;padding:0 16px;background:var(--glass);backdrop-filter:blur(12px);border-bottom:1px solid var(--glass-border);z-index:50}
.logo{font-size:13px;font-weight:700;letter-spacing:-0.02em}.logo span{color:var(--accent)}
.hdr select{background:var(--elevated);border:1px solid var(--border);color:var(--text);padding:4px 8px;border-radius:6px;font-size:11px;font-family:var(--font)}
.hdr-btn{background:none;border:1px solid var(--border);color:var(--text2);width:30px;height:30px;border-radius:6px;cursor:pointer;font-size:12px;display:flex;align-items:center;justify-content:center}
.hdr-btn:hover{background:var(--hover);color:var(--text)}
.hdr-right{margin-left:auto;display:flex;align-items:center;gap:8px}
.hdr-cost{font-family:var(--mono);font-size:11px;color:var(--green)}
.ws-dot{width:7px;height:7px;border-radius:50%;background:var(--idle)}.ws-dot.live{background:var(--green);box-shadow:0 0 6px var(--green)}
.nav-link{color:var(--text3);font-size:11px;text-decoration:none;padding:3px 8px;border-radius:4px}.nav-link:hover{color:var(--text2)}

/* ═══ SIDEBAR ═══ */
.side{grid-area:side;background:var(--glass);backdrop-filter:blur(16px) saturate(1.2);border-right:1px solid var(--glass-border);overflow-y:auto;padding:12px 0;display:flex;flex-direction:column;gap:4px}
.side-section{padding:0 12px;margin-bottom:4px}
.side-title{font-size:10px;font-weight:600;text-transform:uppercase;letter-spacing:0.08em;color:var(--text3);padding:8px 0 4px}
.agent-item{display:flex;align-items:center;gap:8px;padding:6px 12px;border-radius:6px;cursor:pointer;border-left:2px solid transparent;transition:all 0.1s}
.agent-item:hover{background:var(--hover)}
.agent-item.active{border-left-color:var(--accent);background:var(--hover)}
.agent-dot{width:7px;height:7px;border-radius:50%;flex-shrink:0}
.agent-name{font-size:12px;font-weight:500;flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.agent-mode{font-size:8px;padding:1px 4px;border-radius:3px;flex-shrink:0}
.agent-mode.oneshot{background:rgba(251,191,36,0.15);color:var(--yellow)}
.agent-mode.session{background:rgba(129,140,248,0.15);color:var(--accent)}
.agent-activity{font-size:10px;color:var(--text3);padding-left:15px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;margin-top:-2px}
.mail-item{display:flex;align-items:center;gap:8px;padding:6px 12px;border-radius:6px;cursor:pointer;font-size:12px;color:var(--text2)}
.mail-item:hover{background:var(--hover)}
.mail-badge{background:var(--red);color:#fff;font-size:9px;font-weight:600;padding:1px 5px;border-radius:8px;min-width:16px;text-align:center}
.mail-badge.warn{background:var(--yellow);color:var(--void)}
.phase-bar{display:flex;gap:3px;padding:4px 0}
.phase-block{height:4px;flex:1;border-radius:2px;background:var(--idle)}
.phase-block.done{background:var(--green)}.phase-block.current{background:var(--accent);animation:pulse 1.5s infinite}
.phase-label{font-size:10px;color:var(--text3);font-family:var(--mono)}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:0.5}}

/* ═══ TABS ═══ */
.tabs{grid-area:tabs;display:flex;align-items:center;gap:0;padding:0 20px;background:var(--base);border-bottom:1px solid var(--border-subtle)}
.tab{padding:8px 14px;font-size:12px;color:var(--text3);cursor:pointer;border-bottom:2px solid transparent;transition:all 0.15s;user-select:none;position:relative}
.tab:hover{color:var(--text2)}
.tab.active{color:var(--text);border-bottom-color:var(--accent)}
.tab .badge{position:absolute;top:4px;right:4px;background:var(--red);color:#fff;font-size:8px;padding:0 4px;border-radius:6px;min-width:14px;text-align:center;line-height:14px}

/* ═══ MAIN ═══ */
.main{grid-area:main;overflow:hidden;display:flex;flex-direction:column}
.view{display:none;flex:1;overflow-y:auto;flex-direction:column}
.view.active{display:flex}

/* ═══ CHAT VIEW ═══ */
.chat-messages{flex:1;overflow-y:auto;padding:20px;display:flex;flex-direction:column;gap:12px}
.chat-input-area{padding:12px 20px;border-top:1px solid var(--border-subtle);display:flex;gap:10px;background:var(--base)}
.chat-input{flex:1;background:var(--elevated);border:1px solid var(--border);color:var(--text);padding:10px 14px;border-radius:10px;font-size:13px;font-family:var(--font);resize:none;min-height:42px;max-height:120px}
.chat-input:focus{outline:none;border-color:var(--accent);box-shadow:0 0 0 3px var(--glow)}
.chat-input::placeholder{color:var(--text3)}
.chat-send{background:var(--accent);color:var(--void);border:none;padding:10px 18px;border-radius:10px;font-size:13px;font-weight:600;cursor:pointer;align-self:flex-end}
.chat-send:hover{filter:brightness(1.1)}.chat-send:disabled{opacity:0.4}

/* Chat bubbles */
.msg{display:flex;gap:10px;max-width:85%}
.msg.user{align-self:flex-end;flex-direction:row-reverse}
.msg-avatar{width:28px;height:28px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:12px;flex-shrink:0;font-weight:600}
.msg.user .msg-avatar{background:var(--accent);color:var(--void)}
.msg.brain .msg-avatar{background:var(--elevated);border:1px solid var(--border)}
.msg-body{background:var(--elevated);border:1px solid var(--border-subtle);border-radius:4px 12px 12px 12px;padding:10px 14px;font-size:13px;line-height:1.6}
.msg.user .msg-body{background:rgba(129,140,248,0.1);border:1px solid rgba(129,140,248,0.15);border-radius:12px 4px 12px 12px}
.msg-meta{font-size:10px;color:var(--text3);margin-bottom:4px;display:flex;align-items:center;gap:6px}
.msg-content{color:var(--text2)}
.msg-content strong{color:var(--text)}
.msg-content code{font-family:var(--mono);font-size:11px;background:rgba(255,255,255,0.05);padding:1px 4px;border-radius:3px}
.msg-content table{width:100%;border-collapse:collapse;margin:8px 0;font-size:11px}
.msg-content th,.msg-content td{padding:4px 8px;border:1px solid var(--border);text-align:left}
.msg-content th{background:rgba(255,255,255,0.03);font-weight:500}

/* Checkpoint card in chat */
.checkpoint-card{background:rgba(251,191,36,0.04);border:1px solid rgba(251,191,36,0.15);border-left:3px solid var(--yellow);border-radius:8px;padding:14px;margin:4px 0}
.checkpoint-card .cp-title{font-size:12px;font-weight:600;color:var(--yellow);margin-bottom:6px}
.checkpoint-card .cp-body{font-size:12px;color:var(--text2);margin-bottom:10px}
.checkpoint-card .cp-actions{display:flex;gap:8px}
.cp-btn{padding:6px 14px;border-radius:6px;font-size:11px;font-weight:600;border:none;cursor:pointer}
.cp-btn.approve{background:var(--green);color:var(--void)}.cp-btn.approve:hover{filter:brightness(1.1)}
.cp-btn.reject{background:transparent;border:1px solid var(--red);color:var(--red)}.cp-btn.reject:hover{background:rgba(248,113,113,0.1)}

/* Agent inline activity */
.agent-inline{display:flex;align-items:center;gap:6px;padding:6px 10px;background:rgba(129,140,248,0.05);border-radius:6px;font-size:11px;font-family:var(--mono);color:var(--text2);margin:4px 0}
.agent-inline .dot{width:6px;height:6px;border-radius:50%;animation:pulse 1.5s infinite}

/* ═══ TEAM VIEW ═══ */
.team-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:12px;padding:20px}
.team-card{background:var(--elevated);border:1px solid var(--border-subtle);border-radius:10px;padding:14px;cursor:pointer;transition:all 0.15s;position:relative;overflow:hidden}
.team-card:hover{border-color:var(--border);background:var(--hover)}
.team-card .indicator{position:absolute;top:0;left:0;right:0;height:2px;background:var(--idle)}
.team-card.working .indicator{background:var(--accent);animation:pulse 1.5s infinite}
.team-card.complete .indicator{background:var(--green)}
.tc-header{display:flex;align-items:center;gap:8px;margin-bottom:8px}
.tc-avatar{width:32px;height:32px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:14px}
.tc-name{font-size:13px;font-weight:600}
.tc-mode{font-size:9px;padding:1px 5px;border-radius:3px}
.tc-status{font-size:11px;color:var(--text3);margin-bottom:8px}
.tc-metrics{display:flex;gap:12px;font-family:var(--mono);font-size:10px;color:var(--text3)}
.tc-tools{font-size:10px;color:var(--text3);margin-top:6px;font-family:var(--mono);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}

/* ═══ MAIL VIEW ═══ */
.mail-list{padding:20px;display:flex;flex-direction:column;gap:6px}
.mail-row{display:flex;align-items:center;gap:12px;padding:10px 14px;background:var(--elevated);border:1px solid var(--border-subtle);border-radius:8px;cursor:pointer;transition:all 0.1s}
.mail-row:hover{border-color:var(--border);background:var(--hover)}
.mail-row.unread{border-left:3px solid var(--accent)}
.mail-row.alert{border-left:3px solid var(--red)}
.mail-from{width:160px;font-size:12px;font-weight:500;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.mail-subject{flex:1;font-size:12px;color:var(--text2);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.mail-trust{font-size:9px;padding:2px 6px;border-radius:4px;font-weight:600;flex-shrink:0}
.mail-trust.trusted{background:rgba(52,211,153,0.12);color:var(--green)}
.mail-trust.new_contact{background:rgba(251,191,36,0.12);color:var(--yellow)}
.mail-trust.impersonation{background:rgba(248,113,113,0.12);color:var(--red);animation:pulse 2s infinite}
.mail-time{font-size:10px;color:var(--text3);font-family:var(--mono);flex-shrink:0}

/* ═══ PROJECT VIEW ═══ */
.proj-grid{display:grid;grid-template-columns:1fr 1fr;gap:12px;padding:20px}
.proj-card{background:var(--elevated);border:1px solid var(--border-subtle);border-radius:10px;padding:14px}
.proj-card-title{font-size:10px;font-weight:600;text-transform:uppercase;letter-spacing:0.06em;color:var(--text3);margin-bottom:10px}
.proj-file{display:flex;gap:6px;padding:3px 0;cursor:pointer;color:var(--text2);font-family:var(--mono);font-size:11px}
.proj-file:hover{color:var(--accent)}

/* ═══ STATUS BAR ═══ */
.status{grid-area:status;display:flex;align-items:center;gap:16px;padding:0 16px;background:var(--base);border-top:1px solid var(--border-subtle);font-family:var(--mono);font-size:11px;color:var(--text3)}
.status .sep{color:var(--text3);opacity:0.3}
.status .active-agent{color:var(--accent)}

/* ═══ TOAST ═══ */
.toast-container{position:fixed;top:56px;right:16px;z-index:100;display:flex;flex-direction:column;gap:8px}
.toast{background:var(--elevated);border:1px solid var(--border);border-radius:8px;padding:10px 14px;font-size:12px;animation:slideIn 0.3s ease;max-width:320px;backdrop-filter:blur(12px)}
@keyframes slideIn{from{transform:translateX(100%);opacity:0}to{transform:translateX(0);opacity:1}}

/* ═══ EMPTY STATE ═══ */
.empty{display:flex;flex-direction:column;align-items:center;justify-content:center;height:100%;color:var(--text3);gap:8px;font-size:13px}
.empty .icon{font-size:40px;opacity:0.3}
</style>
</head>
<body>
<div class="app">
<!-- HEADER -->
<div class="hdr">
  <div class="logo">Agent<span>House</span></div>
  <select id="projectSel" onchange="switchProject(this.value)"></select>
  <button class="hdr-btn" onclick="newProject()" title="New project">+</button>
  <div class="hdr-right">
    <span class="hdr-cost" id="totalCost">$0.00</span>
    <div class="ws-dot" id="wsDot"></div>
    <a class="nav-link" href="/live">Live</a>
    <a class="nav-link" href="/email-sim">Email Sim</a>
    <button class="hdr-btn" id="cmdKBtn" title="Cmd+K">⌘K</button>
  </div>
</div>

<!-- SIDEBAR -->
<div class="side" id="sidebar">
  <div class="side-section">
    <div class="side-title">Team</div>
    <div id="agentList"></div>
  </div>
  <div class="side-section">
    <div class="side-title">Mail</div>
    <div class="mail-item" onclick="switchTab('mail')">
      <span>📧</span><span>Inbox</span>
      <span class="mail-badge" id="mailCount" style="display:none">0</span>
      <span class="mail-badge warn" id="alertCount" style="display:none">0</span>
    </div>
  </div>
  <div class="side-section">
    <div class="side-title">Phase</div>
    <div class="phase-bar" id="phaseBar"></div>
    <div class="phase-label" id="phaseLabel">Idle</div>
  </div>
</div>

<!-- TABS -->
<div class="tabs">
  <div class="tab active" data-tab="chat" onclick="switchTab('chat')">💬 Chat</div>
  <div class="tab" data-tab="team" onclick="switchTab('team')">👥 Team</div>
  <div class="tab" data-tab="board" onclick="switchTab('board')">📋 Board</div>
  <div class="tab" data-tab="mail" onclick="switchTab('mail')">📧 Mail<span class="badge" id="mailTabBadge" style="display:none"></span></div>
  <div class="tab" data-tab="project" onclick="switchTab('project')">📊 Project</div>
</div>

<!-- MAIN -->
<div class="main">
  <!-- Chat View -->
  <div class="view active" id="view-chat">
    <div class="chat-messages" id="chatMessages">
      <div class="msg brain">
        <div class="msg-avatar">🧠</div>
        <div class="msg-body"><div class="msg-meta">Brain</div><div class="msg-content">Hello! I'm the Brain of Agent House. Type anything — ask questions, create projects, check status, or just chat.</div></div>
      </div>
    </div>
    <div class="chat-input-area">
      <textarea class="chat-input" id="chatInput" placeholder="Chat with Brain — ask anything..." rows="1" onkeydown="if(event.key==='Enter'&&!event.shiftKey){event.preventDefault();sendChat()}"></textarea>
      <button class="chat-send" id="chatSend" onclick="sendChat()">Send</button>
    </div>
  </div>

  <!-- Team View -->
  <div class="view" id="view-team">
    <div class="team-grid" id="teamGrid"></div>
  </div>

  <!-- Board View -->
  <div class="view" id="view-board">
    <div class="empty"><div class="icon">📋</div>Kanban board — coming soon</div>
  </div>

  <!-- Mail View -->
  <div class="view" id="view-mail">
    <div class="mail-list" id="mailList"><div class="empty"><div class="icon">📧</div>No emails yet. Use the <a href="/email-sim" style="color:var(--accent)">Email Simulator</a> to send test emails.</div></div>
  </div>

  <!-- Project View -->
  <div class="view" id="view-project">
    <div class="proj-grid" id="projGrid"><div class="empty"><div class="icon">📊</div>Select a project to view details</div></div>
  </div>
</div>

<!-- STATUS BAR -->
<div class="status">
  <span id="statusAgent">Ready</span>
  <span class="sep">│</span>
  <span id="statusTools">0 tools</span>
  <span class="sep">│</span>
  <span id="statusFiles">0 files</span>
  <span class="sep">│</span>
  <span id="statusPhase">Idle</span>
</div>
</div>

<div class="toast-container" id="toasts"></div>

<script>
// ═══ STATE ═══
const S = {
  ws: null, projectId: 'default',
  agents: {}, agentModes: {},
  mailCount: 0, alertCount: 0,
  phase: '', totalCost: 0,
  chatHistory: []
};

// ═══ INIT ═══
async function init() {
  connectWS();
  await Promise.all([fetchProjects(), fetchAgentModes(), fetchInbox(), fetchActiveSessions()]);
  renderAgentList();
  renderTeamGrid();
  // Refresh sessions every 10s to catch new agents
  setInterval(fetchActiveSessions, 10000);
}

async function fetchActiveSessions() {
  try {
    const d = await(await fetch('/api/sessions')).json();
    (d.sessions||[]).forEach(s => {
      const r = s.agent_role;
      if(!r) return;
      if(!S.agentModes[r]) S.agentModes[r] = 'session';
      if(!S.agents[r]) S.agents[r] = {state:'idle',tools:0,activity:''};
      // Update from session stats
      const a = S.agents[r];
      a.tools = s.stats?.total_tool_calls || a.tools;
      if(s.state === 'active') a.state = 'working';
    });
    renderAgentList();
  } catch(e) {}
}

// ═══ WEBSOCKET ═══
function connectWS() {
  const proto = location.protocol==='https:'?'wss:':'ws:';
  S.ws = new WebSocket(proto+'//'+location.host+'/ws');
  S.ws.onopen = () => { document.getElementById('wsDot').classList.add('live'); };
  S.ws.onclose = () => { document.getElementById('wsDot').classList.remove('live'); setTimeout(connectWS,3000); };
  S.ws.onerror = () => S.ws.close();
  S.ws.onmessage = e => { try{handleWS(JSON.parse(e.data))}catch(x){console.error(x)} };
}

function handleWS(data) {
  if(data.type==='agent_event'&&data.event) handleAgentEvent(data.event);
  if(data.type==='message'&&data.message) {
    if(data.message.type==='lifecycle') handleLifecycle(data.message);
  }
  if(data.type==='checkpoint'&&data.event) handleCheckpoint(data.event);
  if(data.type==='brain_event'&&data.event) handleBrainEvent(data.event);
}

function handleAgentEvent(ev) {
  const r = ev.agent_role; if(!r) return;
  if(!S.agents[r]) S.agents[r]={state:'idle',tools:0,activity:''};
  // Ensure agent appears in modes (for sidebar)
  if(!S.agentModes[r]) { S.agentModes[r]='session'; }
  const a = S.agents[r];
  switch(ev.type) {
    case 'text_delta': a.state='working'; a.activity=ev.content?.substring(0,50)||''; break;
    case 'thinking_delta': a.state='working'; a.activity='Thinking...'; break;
    case 'tool_use': a.state='working'; a.tools++; a.activity=(ev.tool_name||'')+': '+(formatToolShort(ev.tool_name,ev.input)||''); updateStatus(r,ev.tool_name); break;
    case 'tool_result': break;
    case 'turn_complete':
      a.state='complete'; a.activity='Done';
      S.totalCost += ev.cost_usd||0; document.getElementById('totalCost').textContent='$'+S.totalCost.toFixed(2);
      setTimeout(()=>{if(S.agents[r]?.state==='complete'){S.agents[r].state='idle';S.agents[r].activity='';renderAgentList();renderTeamGrid();}},5000);
      break;
    case 'session_meta': a.state='working'; a.activity='Session started'; break;
    case 'error': a.state='error'; a.activity='Error'; break;
  }
  renderAgentList(); renderTeamGrid();
}

function handleLifecycle(msg) {
  const x = msg.metadata?.extra||{};
  if(x.event_type==='task_completed') {
    Object.keys(S.agents).forEach(k=>{S.agents[k].state='idle';S.agents[k].activity='';});
    renderAgentList(); renderTeamGrid();
    addChatMsg('brain','✅ Task completed: '+(x.turns||0)+' turns, '+(x.files_count||0)+' files created.');
    toast('Task Complete','All agents finished','success');
  }
  if(x.event_type==='phase_started') {
    S.phase = x.phase_name||'';
    document.getElementById('statusPhase').textContent = S.phase;
    document.getElementById('phaseLabel').textContent = S.phase;
    addChatMsg('system','Phase: '+S.phase);
  }
  if(x.event_type==='phase_completed') {
    addChatMsg('system','Phase completed: '+(x.phase_name||''));
  }
}

function handleCheckpoint(evt) {
  if(evt.event_type==='checkpoint_reached') {
    const cp = evt;
    const html = '<div class="checkpoint-card"><div class="cp-title">✋ CHECKPOINT: '+(cp.checkpoint_type||'Approval')+'</div>'+
      '<div class="cp-body">'+(cp.data?.artifact_summary||'Review required')+'</div>'+
      '<div class="cp-actions">'+
      '<button class="cp-btn approve" onclick="resolveCP(\''+cp.checkpoint_id+'\',\'approved\')">✅ Approve</button>'+
      '<button class="cp-btn reject" onclick="resolveCP(\''+cp.checkpoint_id+'\',\'rejected\')">❌ Reject</button>'+
      '</div></div>';
    addChatHTML(html);
  }
}

function handleBrainEvent(ev) {
  // Brain activity events from WebSocket
}

// ═══ CHAT ═══
async function sendChat() {
  const inp = document.getElementById('chatInput');
  const msg = inp.value.trim(); if(!msg) return;
  inp.value=''; inp.style.height='42px';
  document.getElementById('chatSend').disabled=true;

  addChatMsg('user', msg);

  try {
    const res = await fetch('/api/chat',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({message:msg,user_id:'default'})});
    const d = await res.json();
    const icon = d.action==='create_project'?'🚀':d.action==='delegate'?'📌':d.action==='list_projects'?'📋':d.action==='project_status'?'📊':d.action==='delete_email'?'🗑️':d.action==='send_reply'?'📤':'';
    addChatMsg('brain', (icon?icon+' ':'')+(d.response||'Done'));

    // Render suggestion buttons (handle null/undefined)
    var suggestions = d.suggestions || [];
    if(suggestions.length) {
      const el = document.getElementById('chatMessages');
      const div = document.createElement('div');
      div.style.cssText = 'display:flex;flex-wrap:wrap;gap:6px;padding:4px 38px;margin-top:-4px';
      div.innerHTML = suggestions.map(s =>
        '<button onclick="useSuggestion(this)" style="background:var(--elevated);border:1px solid var(--border);color:var(--text2);padding:5px 12px;border-radius:16px;font-size:11px;cursor:pointer;font-family:var(--font);transition:all 0.1s" onmouseover="this.style.borderColor=\'var(--accent)\';this.style.color=\'var(--text)\'" onmouseout="this.style.borderColor=\'var(--border)\';this.style.color=\'var(--text2)\'">' + esc(s) + '</button>'
      ).join('');
      el.appendChild(div);
      el.scrollTop = el.scrollHeight;
    }

    // Refresh mail badges after email-related actions
    var r = (d.response||'').toLowerCase();
    if(d.action==='delete_email'||d.action==='send_reply'||d.action==='archive_email'||
       r.indexOf('deleted')>-1||r.indexOf('sent')>-1||r.indexOf('replied')>-1||
       r.indexOf('email')>-1||r.indexOf('inbox')>-1||r.indexOf('mail')>-1) {
      fetchInbox();
    }

    if(d.action==='create_project'&&d.success&&d.project_id) {
      S.projectId=d.project_id;
      await fetchProjects();
      document.getElementById('projectSel').value=d.project_id;
    }
  } catch(e) {
    addChatMsg('brain','❌ Error: '+e.message);
  }
  document.getElementById('chatSend').disabled=false;
}

function useSuggestion(btn) {
  const text = btn.textContent;
  // Remove all suggestion buttons
  const parent = btn.parentElement;
  if(parent) parent.remove();
  // Set as input and send
  document.getElementById('chatInput').value = text;
  sendChat();
}

function addChatMsg(role, text) {
  const el = document.getElementById('chatMessages');
  const div = document.createElement('div');
  if(role==='user') {
    div.className='msg user';
    div.innerHTML='<div class="msg-avatar">U</div><div class="msg-body"><div class="msg-meta">You</div><div class="msg-content">'+esc(text)+'</div></div>';
  } else if(role==='brain') {
    div.className='msg brain';
    div.innerHTML='<div class="msg-avatar">🧠</div><div class="msg-body"><div class="msg-meta">Brain</div><div class="msg-content">'+renderMd(text)+'</div></div>';
  } else {
    div.className='msg brain';
    div.innerHTML='<div class="msg-avatar">⚙</div><div class="msg-body"><div class="msg-meta" style="color:var(--text3)">System</div><div class="msg-content" style="color:var(--text3);font-size:11px">'+esc(text)+'</div></div>';
  }
  el.appendChild(div);
  el.scrollTop=el.scrollHeight;
}

function addChatHTML(html) {
  const el = document.getElementById('chatMessages');
  const div = document.createElement('div');
  div.innerHTML = html;
  el.appendChild(div);
  el.scrollTop = el.scrollHeight;
}

async function resolveCP(cpId, action) {
  try {
    await fetch('/api/checkpoints/'+cpId+'/decide?project='+S.projectId,{
      method:'POST',headers:{'Content-Type':'application/json'},
      body:JSON.stringify({action,decided_by:'human',feedback:''})
    });
    addChatMsg('system','Checkpoint '+action);
  } catch(e) { addChatMsg('brain','Error: '+e.message); }
}

// ═══ AGENT LIST (sidebar) ═══
function renderAgentList() {
  const el = document.getElementById('agentList');
  const roles = Object.keys(S.agentModes).sort();
  if(!roles.length) { el.innerHTML='<div style="font-size:11px;color:var(--text3);padding:4px 0">No agents loaded</div>'; return; }
  el.innerHTML = roles.map(r => {
    const a = S.agents[r]||{state:'idle',tools:0,activity:''};
    const mode = S.agentModes[r]||'session';
    const dotColor = a.state==='working'?'var(--accent)':a.state==='complete'?'var(--green)':a.state==='error'?'var(--red)':'var(--idle)';
    const activeClass = a.state==='working'?' active':'';
    return '<div class="agent-item'+activeClass+'" onclick="selectAgent(\''+r+'\')">'+
      '<div class="agent-dot" style="background:'+dotColor+(a.state==='working'?';animation:pulse 1.5s infinite':'')+'"></div>'+
      '<div class="agent-name">'+esc(r.replace(/_/g,' '))+'</div>'+
      '<div class="agent-mode '+mode+'">'+(mode==='oneshot'?'⚡':'🔗')+'</div>'+
      '</div>'+(a.activity?'<div class="agent-activity">'+esc(a.activity.substring(0,40))+'</div>':'');
  }).join('');
}

// ═══ TEAM GRID ═══
function renderTeamGrid() {
  const el = document.getElementById('teamGrid');
  const roles = Object.keys(S.agentModes).sort();
  if(!roles.length) { el.innerHTML='<div class="empty"><div class="icon">👥</div>No agents</div>'; return; }
  el.innerHTML = roles.map(r => {
    const a = S.agents[r]||{state:'idle',tools:0,activity:''};
    const mode = S.agentModes[r]||'session';
    const stateClass = a.state==='working'?' working':a.state==='complete'?' complete':'';
    return '<div class="team-card'+stateClass+'"><div class="indicator"></div>'+
      '<div class="tc-header"><div class="tc-name">'+esc(r.replace(/_/g,' '))+'</div>'+
      '<div class="tc-mode '+(mode==='oneshot'?'agent-mode oneshot':'agent-mode session')+'">'+(mode==='oneshot'?'⚡ API':'🔗 Session')+'</div></div>'+
      '<div class="tc-status">'+(a.state||'idle')+(a.activity?' — '+esc(a.activity.substring(0,40)):'')+'</div>'+
      '<div class="tc-metrics"><span>'+a.tools+' tools</span></div>'+
      '</div>';
  }).join('');
}

// ═══ MAIL VIEW ═══
async function fetchInbox() {
  try {
    const d = await(await fetch('/api/email/inbox')).json();
    const emails = d.emails||[];
    S.mailCount = emails.filter(e=>!e.read).length;
    S.alertCount = emails.filter(e=>e.trust_status==='impersonation').length;
    renderMailBadges();
    renderMailList(emails);
  } catch(e) {}
}

function renderMailBadges() {
  const mc=document.getElementById('mailCount');
  const ac=document.getElementById('alertCount');
  const tb=document.getElementById('mailTabBadge');
  if(S.mailCount>0){mc.textContent=S.mailCount;mc.style.display='';}else mc.style.display='none';
  if(S.alertCount>0){ac.textContent=S.alertCount;ac.style.display='';}else ac.style.display='none';
  if(S.mailCount>0){tb.textContent=S.mailCount;tb.style.display='';}else tb.style.display='none';
}

function renderMailList(emails) {
  S.emails = emails;
  const el = document.getElementById('mailList');
  if(!emails.length) { el.innerHTML='<div class="empty"><div class="icon">📧</div>No emails yet. <a href="/email-sim" style="color:var(--accent)">Send test emails</a></div>'; return; }
  el.innerHTML = '<div id="mailListInner">' + emails.map((e,i) => {
    const unread = !e.read?' unread':'';
    const alert = e.trust_status==='impersonation'?' alert':'';
    const time = e.date?new Date(e.date).toLocaleTimeString('en-US',{hour:'2-digit',minute:'2-digit'}):'';
    var replied = e.replied ? '<span style="color:var(--green);font-size:9px;margin-left:4px">✓ replied</span>' : '';
    return '<div class="mail-row'+unread+alert+'" onclick="openMail('+i+')">'+
      '<div class="mail-from">'+esc(e.from_name||e.from)+replied+'</div>'+
      '<div class="mail-subject">'+esc(e.subject)+'</div>'+
      '<div class="mail-trust '+(e.trust_status||'')+'">'+esc((e.trust_status||'').replace(/_/g,' '))+'</div>'+
      '<div class="mail-time">'+time+'</div></div>';
  }).join('') + '</div><div id="mailDetail" style="display:none"></div>';
}

function openMail(idx) {
  const e = S.emails[idx]; if(!e) return;
  document.getElementById('mailListInner').style.display='none';
  const det = document.getElementById('mailDetail');
  det.style.display='block';

  const isThreat = e.trust_status==='impersonation';
  const isNew = e.trust_status==='new_contact';

  let html = '<div style="padding:20px">';
  // Back button
  html += '<div style="margin-bottom:16px"><button onclick="closeMail()" style="background:none;border:1px solid var(--border);color:var(--text2);padding:5px 12px;border-radius:6px;cursor:pointer;font-size:11px">← Back to Inbox</button></div>';

  // Trust alert banner
  if(isThreat) {
    html += '<div style="background:rgba(248,113,113,0.08);border:1px solid rgba(248,113,113,0.25);border-left:3px solid var(--red);border-radius:8px;padding:14px;margin-bottom:16px">';
    html += '<div style="font-weight:600;color:var(--red);margin-bottom:6px">🚨 IMPERSONATION RISK</div>';
    html += '<div style="font-size:12px;color:var(--text2)">'+esc(e.trust_reason||'Sender not verified')+'</div>';
    html += '<div style="margin-top:10px;display:flex;gap:8px">';
    html += '<button onclick="rejectMail(\''+e.id+'\')" style="background:var(--red);color:#fff;border:none;padding:6px 14px;border-radius:6px;font-size:11px;font-weight:600;cursor:pointer">🚫 Reject & Flag</button>';
    html += '<button onclick="verifyMail(\''+e.id+'\')" style="background:none;border:1px solid var(--border);color:var(--text2);padding:6px 14px;border-radius:6px;font-size:11px;cursor:pointer">🔍 Verify Manually</button>';
    html += '</div></div>';
  }
  if(isNew) {
    html += '<div style="background:rgba(251,191,36,0.08);border:1px solid rgba(251,191,36,0.25);border-left:3px solid var(--yellow);border-radius:8px;padding:14px;margin-bottom:16px">';
    html += '<div style="font-weight:600;color:var(--yellow);margin-bottom:6px">⚠️ New Contact</div>';
    html += '<div style="font-size:12px;color:var(--text2)">'+esc(e.trust_reason||'')+'</div>';
    html += '<div style="margin-top:10px">';
    html += '<button onclick="addToTrusted(\''+esc(e.vendor_id||'')+'\',\''+esc(e.from)+'\')" style="background:var(--green);color:var(--void);border:none;padding:6px 14px;border-radius:6px;font-size:11px;font-weight:600;cursor:pointer">✅ Add to Trusted List</button>';
    html += '</div></div>';
  }

  // Email header
  html += '<div style="background:var(--elevated);border:1px solid var(--border-subtle);border-radius:10px;padding:16px;margin-bottom:16px">';
  html += '<div style="font-size:15px;font-weight:600;margin-bottom:10px">'+esc(e.subject)+'</div>';
  html += '<div style="display:flex;gap:16px;font-size:12px;color:var(--text2);margin-bottom:12px">';
  html += '<span><strong>From:</strong> '+esc(e.from_name||'')+' &lt;'+esc(e.from)+'&gt;</span>';
  html += '<span><strong>To:</strong> '+esc(e.to)+'</span>';
  html += '<span style="color:var(--text3)">'+new Date(e.date).toLocaleString()+'</span>';
  html += '</div>';
  html += '<div style="font-size:13px;line-height:1.7;color:var(--text2);white-space:pre-wrap;border-top:1px solid var(--border-subtle);padding-top:12px">'+esc(e.body)+'</div>';
  html += '</div>';

  // Actions
  html += '<div style="display:flex;gap:8px;margin-bottom:16px">';
  html += '<button onclick="suggestReply('+idx+')" style="background:var(--accent);color:var(--void);border:none;padding:8px 16px;border-radius:8px;font-size:12px;font-weight:600;cursor:pointer">🧠 Suggest Reply</button>';
  html += '<button onclick="chatAboutMail('+idx+')" style="background:none;border:1px solid var(--border);color:var(--text2);padding:8px 16px;border-radius:8px;font-size:12px;cursor:pointer">💬 Chat About This</button>';
  if(e.vendor_id && !isThreat) {
    html += '<button onclick="addToTrusted(\''+esc(e.vendor_id)+'\',\''+esc(e.from)+'\')" style="background:none;border:1px solid var(--border);color:var(--text2);padding:8px 16px;border-radius:8px;font-size:12px;cursor:pointer">📋 Add to Trusted</button>';
  }
  html += '</div>';

  // Reply area
  html += '<div id="replyArea"></div>';
  html += '</div>';
  det.innerHTML = html;
}

function closeMail() {
  document.getElementById('mailListInner').style.display='';
  document.getElementById('mailDetail').style.display='none';
}

async function suggestReply(idx) {
  const e = S.emails[idx]; if(!e) return;
  const area = document.getElementById('replyArea');
  area.innerHTML = '<div style="padding:14px;background:var(--elevated);border:1px solid var(--border-subtle);border-radius:10px"><div style="color:var(--accent);font-size:12px;margin-bottom:8px">🧠 Generating reply suggestion...</div></div>';

  try {
    const res = await fetch('/api/chat',{method:'POST',headers:{'Content-Type':'application/json'},
      body:JSON.stringify({message:'Draft a professional reply to this email:\n\nFrom: '+e.from_name+' <'+e.from+'>\nSubject: '+e.subject+'\n\n'+e.body,user_id:'default'})});
    const d = await res.json();
    area.innerHTML = '<div style="padding:14px;background:var(--elevated);border:1px solid var(--border-subtle);border-radius:10px">'+
      '<div style="font-size:11px;font-weight:600;color:var(--accent);margin-bottom:8px">Suggested Reply</div>'+
      '<div style="font-size:13px;line-height:1.6;color:var(--text2);white-space:pre-wrap">'+renderMd(d.response||'')+'</div>'+
      '<div style="margin-top:12px;display:flex;gap:8px">'+
      '<button onclick="sendReply()" style="background:var(--green);color:var(--void);border:none;padding:6px 14px;border-radius:6px;font-size:11px;font-weight:600;cursor:pointer">📤 Send Reply</button>'+
      '<button onclick="editReply()" style="background:none;border:1px solid var(--border);color:var(--text2);padding:6px 14px;border-radius:6px;font-size:11px;cursor:pointer">✏️ Edit</button>'+
      '</div></div>';
  } catch(e) {
    area.innerHTML = '<div style="color:var(--red);font-size:12px">Failed: '+e.message+'</div>';
  }
}

function chatAboutMail(idx) {
  const e = S.emails[idx]; if(!e) return;
  switchTab('chat');
  const inp = document.getElementById('chatInput');
  inp.value = 'Regarding email from '+e.from_name+' about "'+e.subject+'" — ';
  inp.focus();
}

async function addToTrusted(vendorId, emailAddr) {
  if(!vendorId){toast('Error','No vendor ID','error');return;}
  try {
    await fetch('/api/email/trust',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({vendor_id:vendorId,email:emailAddr})});
    toast('Added',''+emailAddr+' added to trusted list','success');
  } catch(e) { toast('Error',e.message,'error'); }
}

function rejectMail(id) { toast('Rejected','Email flagged as fraud attempt','success'); }
function verifyMail(id) { toast('Verify','Please contact vendor via known phone number','info'); }
function sendReply() { toast('Sent','Reply sent (RESEND_API_KEY required for real delivery)','success'); }
function editReply() { toast('Edit','Edit feature coming soon','info'); }

// ═══ PROJECT VIEW ═══
async function loadProjectView() {
  const el = document.getElementById('projGrid');
  try {
    const d = await(await fetch('/api/projects/'+encodeURIComponent(S.projectId))).json();
    const files = d.files||[];
    const sessions = d.sessions||{};
    const agents = sessions.agent_metrics||{};
    let html = '<div class="proj-card"><div class="proj-card-title">Files ('+files.length+')</div>';
    if(files.length) {
      html += files.slice(0,20).map(f => '<div class="proj-file">📄 '+esc(f.path)+'</div>').join('');
    } else html += '<div style="color:var(--text3);font-size:12px">No files</div>';
    html += '</div>';
    html += '<div class="proj-card"><div class="proj-card-title">Agent Metrics</div>';
    Object.entries(agents).sort().forEach(([r,m]) => {
      html += '<div style="display:flex;gap:8px;font-size:11px;padding:3px 0;color:var(--text2)"><span style="width:80px">'+esc(m.name||r)+'</span>'+
        '<span style="font-family:var(--mono)">'+(m.turns||0)+'t '+(m.tool_calls||0)+'tc $'+(m.cost_usd||0).toFixed(2)+'</span></div>';
    });
    if(!Object.keys(agents).length) html += '<div style="color:var(--text3);font-size:12px">No session data</div>';
    html += '</div>';
    el.innerHTML = html;
  } catch(e) { el.innerHTML = '<div class="empty">Failed to load project</div>'; }
}

// ═══ TABS ═══
function switchTab(tab) {
  document.querySelectorAll('.tab').forEach(t=>t.classList.toggle('active',t.dataset.tab===tab));
  document.querySelectorAll('.view').forEach(v=>v.classList.toggle('active',v.id==='view-'+tab));
  if(tab==='mail') fetchInbox();
  if(tab==='project') loadProjectView();
}

// ═══ FETCH ═══
async function fetchProjects() {
  try {
    const d = await(await fetch('/api/projects')).json();
    const sel = document.getElementById('projectSel'); sel.innerHTML='';
    (d.projects||[]).forEach(p=>{const o=document.createElement('option');o.value=p.id;o.textContent=p.id;if(p.id===S.projectId)o.selected=true;sel.appendChild(o);});
  }catch(e){}
}

async function fetchAgentModes() {
  try {
    const d = await(await fetch('/api/agents/modes')).json();
    S.agentModes = d.modes||{};
    renderAgentList();
  }catch(e){}
}

function switchProject(id) { S.projectId=id; S.agents={}; renderAgentList(); renderTeamGrid(); }
function newProject() { const n=prompt('Project name:'); if(n){S.projectId=n.toLowerCase().replace(/\s+/g,'-');fetchProjects();} }
function selectAgent(role) { switchTab('team'); }

// ═══ STATUS BAR ═══
function updateStatus(agentRole, toolName) {
  document.getElementById('statusAgent').innerHTML='<span class="active-agent">◉ '+esc(agentRole.replace(/_/g,' '))+': '+esc(toolName||'working')+'</span>';
  const totalTools = Object.values(S.agents).reduce((s,a)=>s+(a.tools||0),0);
  document.getElementById('statusTools').textContent=totalTools+' tools';
}

// ═══ HELPERS ═══
function esc(s){return(s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;')}
function formatToolShort(name,input){try{const p=JSON.parse(input||'{}');switch(name){case'Read':return p.file_path||'';case'Write':return(p.file_path||'')+'';case'Edit':return p.file_path||'';case'Bash':return'$ '+(p.command||'').substring(0,40);case'Grep':return'"'+(p.pattern||'')+'"';default:return''}}catch(e){return''}}
function renderMd(s){
  if(!s)return'';
  // Parse tables BEFORE escaping (pipes and dashes need to be raw)
  var parts=[];
  var lines=s.split('\n');
  var i=0;
  while(i<lines.length){
    // Detect table: line contains | and has table-like structure
    var trimLine = lines[i].trim();
    if(trimLine.indexOf('|')>=0 && (trimLine.indexOf('|')===0 || trimLine.split('|').length>=3)){
      var tableLines=[];
      while(i<lines.length && lines[i].trim().indexOf('|')>=0 && lines[i].trim().split('|').length>=3){
        tableLines.push(lines[i].trim()); i++;
      }
      // Build HTML table
      var dataRows=tableLines.filter(function(r){return !/^[\s|:-]+$/.test(r.replace(/[^|\-:\s]/g,''));});
      if(dataRows.length>0){
        var t='<table style="width:100%;border-collapse:collapse;margin:8px 0;font-size:11px">';
        dataRows.forEach(function(row,ri){
          var cells=row.split('|').filter(function(c,ci,a){return ci>0&&ci<a.length-1||c.trim()!=='';});
          // Clean edge empty cells
          if(cells.length>0&&cells[0].trim()==='')cells.shift();
          if(cells.length>0&&cells[cells.length-1].trim()==='')cells.pop();
          var tag=ri===0?'th':'td';
          var bg=ri===0?'background:rgba(255,255,255,0.04);font-weight:500;color:var(--text)':'';
          t+='<tr>'+cells.map(function(c){return'<'+tag+' style="padding:5px 10px;border:1px solid var(--border);'+bg+'">'+esc(c.trim())+'</'+tag+'>';}).join('')+'</tr>';
        });
        t+='</table>';
        parts.push(t);
      } else {
        tableLines.forEach(function(l){parts.push(esc(l));});
      }
    } else {
      parts.push(null); // placeholder — process later
      i++;
    }
  }
  // Now process non-table lines
  var lineIdx=0;
  var result=[];
  for(var p=0;p<parts.length;p++){
    if(parts[p]!==null){result.push(parts[p]);continue;}
    var line=lines[lineIdx]||'';lineIdx++;
    // Skip if we already consumed this line in table parsing
    while(parts[lineIdx]!==undefined&&parts[lineIdx]!==null)lineIdx++;
    var l=esc(line);
    // Headers
    if(/^### /.test(line))l='<div style="font-size:13px;font-weight:600;margin:10px 0 4px;color:var(--text)">'+esc(line.substring(4))+'</div>';
    else if(/^## /.test(line))l='<div style="font-size:14px;font-weight:600;margin:10px 0 4px;color:var(--text)">'+esc(line.substring(3))+'</div>';
    // HR
    else if(/^---+$/.test(line.trim()))l='<hr style="border:none;border-top:1px solid var(--border);margin:8px 0">';
    // Blockquote
    else if(/^> /.test(line))l='<div style="border-left:2px solid var(--accent);padding-left:10px;color:var(--text2);font-style:italic;margin:4px 0">'+esc(line.substring(2))+'</div>';
    // Bullet
    else if(/^[-*] /.test(line.trim())){var t=line.trim().substring(2);l='<div style="padding-left:12px">• '+esc(t)+'</div>';}
    // Bold (inline)
    l=l.replace(/\*\*(.+?)\*\*/g,'<strong>$1</strong>');
    // Inline code
    var bt=String.fromCharCode(96);
    var codeRe=new RegExp(bt+'([^'+bt+']+)'+bt,'g');
    l=l.replace(codeRe,'<code style="font-family:var(--mono);font-size:11px;background:rgba(255,255,255,0.05);padding:1px 4px;border-radius:3px">$1</code>');
    // Emoji status
    l=l.replace(/✅/g,'<span style="color:var(--green)">✅</span>');
    l=l.replace(/⚠️/g,'<span style="color:var(--yellow)">⚠️</span>');
    l=l.replace(/🚨/g,'<span style="color:var(--red)">🚨</span>');
    result.push(l);
  }
  return result.join('<br>').replace(/<br><(div|table|hr)/g,'<$1').replace(/<\/(div|table)><br>/g,'</$1>').replace(/<br><br><br>/g,'<br><br>');
}
function toast(title,text,type){const c=document.getElementById('toasts');const t=document.createElement('div');t.className='toast';t.innerHTML='<strong>'+esc(title)+'</strong><br><span style="color:var(--text2)">'+esc(text)+'</span>';c.appendChild(t);setTimeout(()=>t.remove(),5000)}

init();
</script>
</body>
</html>
`
