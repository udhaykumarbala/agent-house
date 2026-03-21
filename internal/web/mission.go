package web

// missionHTML contains the enterprise "Mission Control" dashboard
var missionHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Agent House — Mission Control</title>
  <meta name="theme-color" content="#080B14">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
  <style>
/* ═══════════════════════════════════════════
   DESIGN TOKENS
   ═══════════════════════════════════════════ */
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

  /* Agent role colors */
  --c-ceo:      #60B4D4;
  --c-pm:       #8B7EF8;
  --c-ux:       #F47F7F;
  --c-ui:       #4FD1C5;
  --c-security: #F6A623;
  --c-architect:#A87CF5;
  --c-sr-dev:   #34D478;
  --c-jr-dev:   #4BA3E3;

  /* Status */
  --st-idle:     #3A4A5A;
  --st-working:  #3B82F6;
  --st-thinking: #8B5CF6;
  --st-complete: #22C55E;
  --st-error:    #EF4444;

  /* Fonts */
  --font-ui:   'Inter', system-ui, -apple-system, sans-serif;
  --font-mono: 'JetBrains Mono', ui-monospace, monospace;

  /* Layout */
  --header-h:  52px;
  --roster-w:  200px;
  --panel-w:   340px;
  --log-h:     180px;
  --log-collapsed: 44px;
}

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

body {
  font-family: var(--font-ui);
  background: var(--bg-void);
  color: var(--text-primary);
  height: 100vh; overflow: hidden;
  -webkit-font-smoothing: antialiased;
}

/* ═══════════════════════════════════════════
   APP SHELL
   ═══════════════════════════════════════════ */
.app {
  display: grid; height: 100vh;
  grid-template-rows: var(--header-h) 1fr auto;
  grid-template-columns: var(--roster-w) 1fr 0;
  grid-template-areas:
    "header header header"
    "roster canvas panel"
    "log    log    log";
}
.app.panel-open {
  grid-template-columns: var(--roster-w) 1fr var(--panel-w);
}

/* ═══════════════════════════════════════════
   HEADER
   ═══════════════════════════════════════════ */
.header {
  grid-area: header;
  display: flex; align-items: center; gap: 16px;
  padding: 0 20px;
  background: var(--bg-base);
  border-bottom: 1px solid var(--border-subtle);
  z-index: 50;
}
.logo { font-size: 14px; font-weight: 700; color: var(--text-primary); letter-spacing: -0.02em; white-space: nowrap; }
.logo span { color: var(--accent); }

/* Phase Pipeline */
.pipeline { display: flex; align-items: center; gap: 0; flex: 1; justify-content: center; }
.pipe-step {
  display: flex; align-items: center; gap: 6px;
  padding: 4px 12px; font-size: 11px; font-weight: 500;
  letter-spacing: 0.04em; white-space: nowrap;
  color: var(--text-muted); position: relative;
}
.pipe-step .dot {
  width: 8px; height: 8px; border-radius: 50%;
  border: 1.5px solid var(--text-muted); background: transparent;
  flex-shrink: 0;
}
.pipe-step.active { color: var(--accent); }
.pipe-step.active .dot { border-color: var(--accent); background: var(--accent); box-shadow: 0 0 8px var(--accent-glow); }
.pipe-step.done { color: var(--st-complete); }
.pipe-step.done .dot { border-color: var(--st-complete); background: var(--st-complete); }
.pipe-connector { width: 32px; height: 1.5px; background: var(--border-subtle); flex-shrink: 0; }
.pipe-connector.done { background: var(--st-complete); }
.pipe-connector.active { background: linear-gradient(90deg, var(--st-complete), var(--accent)); }

/* KPI Strip */
.kpi-strip {
  display: flex; align-items: center; gap: 16px; margin-left: auto;
  font-family: var(--font-mono); font-size: 11px; color: var(--text-secondary);
}
.kpi { display: flex; align-items: center; gap: 4px; white-space: nowrap; }
.kpi-val { color: var(--text-primary); font-weight: 500; }
.kpi-dot { width: 6px; height: 6px; border-radius: 50%; }
.kpi-dot.live { background: var(--st-complete); animation: pulse-dot 2s ease-in-out infinite; }
.kpi-dot.off  { background: var(--st-error); }

/* Header controls */
.hdr-controls { display: flex; align-items: center; gap: 8px; margin-left: 12px; }
.hdr-btn {
  background: transparent; border: 1px solid var(--border-subtle);
  color: var(--text-secondary); font-family: var(--font-ui); font-size: 11px;
  padding: 4px 10px; border-radius: 4px; cursor: pointer; white-space: nowrap;
}
.hdr-btn:hover { border-color: var(--border-active); color: var(--text-primary); }
.project-sel {
  background: var(--bg-base); border: 1px solid var(--border-subtle);
  color: var(--text-primary); font-family: var(--font-mono); font-size: 11px;
  padding: 4px 8px; border-radius: 4px; cursor: pointer;
}

/* ═══════════════════════════════════════════
   AGENT ROSTER (Left)
   ═══════════════════════════════════════════ */
.roster {
  grid-area: roster;
  background: var(--bg-base);
  border-right: 1px solid var(--border-subtle);
  overflow-y: auto; padding: 8px 0;
}
.roster-title {
  font-size: 10px; font-weight: 600; letter-spacing: 0.08em;
  color: var(--text-muted); padding: 8px 14px 6px; text-transform: uppercase;
}
.roster-card {
  display: flex; align-items: center; gap: 10px;
  padding: 8px 14px; cursor: pointer;
  border-left: 2px solid transparent;
  transition: background 0.15s, border-color 0.15s;
}
.roster-card:hover { background: var(--bg-raised); }
.roster-card.selected { background: var(--bg-raised); border-left-color: var(--accent); }
.roster-card.working { border-left-color: var(--c-ceo); }
.roster-avatar {
  width: 32px; height: 32px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 12px; font-weight: 600; flex-shrink: 0;
  border: 1.5px solid; position: relative;
}
.roster-avatar .state-dot {
  position: absolute; bottom: -1px; right: -1px;
  width: 8px; height: 8px; border-radius: 50%;
  border: 1.5px solid var(--bg-base);
}
.roster-info { flex: 1; min-width: 0; }
.roster-name { font-size: 12px; font-weight: 600; line-height: 1.3; }
.roster-role { font-size: 10px; color: var(--text-muted); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.roster-task { font-size: 10px; color: var(--text-secondary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin-top: 1px; }

/* ═══════════════════════════════════════════
   CANVAS (Center)
   ═══════════════════════════════════════════ */
.canvas-area {
  grid-area: canvas;
  position: relative; overflow: hidden;
  background:
    radial-gradient(circle at 50% 50%, rgba(59,130,246,0.03) 0%, transparent 70%),
    var(--bg-void);
}
/* Dot grid background */
.canvas-area::before {
  content: '';
  position: absolute; inset: 0;
  background-image: radial-gradient(circle, rgba(255,255,255,0.04) 1px, transparent 1px);
  background-size: 28px 28px;
  pointer-events: none;
}

/* SVG arc layer */
.arc-layer {
  position: absolute; inset: 0; pointer-events: none; z-index: 2;
}

/* Agent nodes on canvas */
.canvas-node {
  position: absolute; z-index: 5;
  display: flex; flex-direction: column; align-items: center; gap: 6px;
  cursor: pointer; transform: translate(-50%, -50%);
  transition: transform 0.25s ease, filter 0.25s;
}
.canvas-node:hover { transform: translate(-50%, -50%) scale(1.08); z-index: 10; }
.canvas-node.dimmed { filter: opacity(0.35); }

.node-ring {
  width: 80px; height: 80px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  border: 2px solid var(--st-idle);
  background: rgba(13,17,23,0.8);
  position: relative;
  transition: border-color 0.4s, box-shadow 0.4s, background 0.4s;
}
.node-ring .icon-shape {
  width: 36px; height: 36px;
}
.node-ring .icon-shape svg { width: 100%; height: 100%; }

/* State: working */
.canvas-node[data-state="working"] .node-ring {
  border-color: var(--st-working);
  box-shadow: 0 0 20px rgba(59,130,246,0.25), inset 0 0 12px rgba(59,130,246,0.08);
  animation: ring-pulse 2.5s ease-in-out infinite;
}
/* State: thinking */
.canvas-node[data-state="thinking"] .node-ring {
  border-color: var(--st-thinking);
  box-shadow: 0 0 20px rgba(139,92,246,0.25);
  animation: ring-pulse 1.8s ease-in-out infinite;
}
/* State: completed */
.canvas-node[data-state="completed"] .node-ring {
  border-color: var(--st-complete);
  box-shadow: 0 0 12px rgba(34,197,94,0.2);
}
/* State: error */
.canvas-node[data-state="error"] .node-ring {
  border-color: var(--st-error);
  box-shadow: 0 0 16px rgba(239,68,68,0.3);
  animation: ring-error 0.8s ease-in-out 3;
}

@keyframes ring-pulse {
  0%, 100% { box-shadow: 0 0 16px rgba(59,130,246,0.2); transform: scale(1); }
  50% { box-shadow: 0 0 28px rgba(59,130,246,0.35); transform: scale(1.04); }
}
@keyframes ring-error {
  0%, 100% { box-shadow: 0 0 16px rgba(239,68,68,0.3); }
  50% { box-shadow: 0 0 32px rgba(239,68,68,0.6); }
}
@keyframes pulse-dot {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.node-label { font-size: 11px; font-weight: 600; color: var(--text-primary); text-align: center; line-height: 1.2; }
.node-status {
  font-size: 9px; font-weight: 500; letter-spacing: 0.06em;
  color: var(--text-muted); text-transform: uppercase;
  padding: 1px 6px; border-radius: 3px; background: var(--bg-raised);
  border: 1px solid var(--border-subtle);
}
.canvas-node[data-state="working"] .node-status { color: var(--st-working); border-color: rgba(59,130,246,0.3); }
.canvas-node[data-state="completed"] .node-status { color: var(--st-complete); border-color: rgba(34,197,94,0.3); }
.canvas-node[data-state="error"] .node-status { color: var(--st-error); border-color: rgba(239,68,68,0.3); }

/* Task label under node when working */
.node-task {
  font-size: 10px; color: var(--text-secondary); max-width: 120px;
  text-align: center; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  position: absolute; top: 100%; margin-top: 24px;
}

/* ═══════════════════════════════════════════
   DETAIL PANEL (Right)
   ═══════════════════════════════════════════ */
.detail-panel {
  grid-area: panel;
  background: var(--bg-base);
  border-left: 1px solid var(--border-subtle);
  display: flex; flex-direction: column;
  overflow: hidden;
  transition: width 0.25s cubic-bezier(0.4,0,0.2,1);
}
.dp-header {
  display: flex; align-items: center; gap: 10px;
  padding: 14px 16px; border-bottom: 1px solid var(--border-subtle);
}
.dp-avatar {
  width: 40px; height: 40px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 14px; font-weight: 700; border: 2px solid;
}
.dp-info { flex: 1; }
.dp-name { font-size: 14px; font-weight: 600; }
.dp-role { font-size: 11px; color: var(--text-secondary); }
.dp-close {
  background: none; border: none; color: var(--text-muted); font-size: 18px;
  cursor: pointer; padding: 4px;
}
.dp-close:hover { color: var(--text-primary); }

/* Panel tabs */
.dp-tabs {
  display: flex; border-bottom: 1px solid var(--border-subtle);
}
.dp-tab {
  flex: 1; padding: 8px; text-align: center; font-size: 11px; font-weight: 500;
  color: var(--text-muted); cursor: pointer; border-bottom: 2px solid transparent;
  transition: color 0.15s, border-color 0.15s;
}
.dp-tab:hover { color: var(--text-secondary); }
.dp-tab.active { color: var(--accent); border-bottom-color: var(--accent); }

/* Panel body */
.dp-body { flex: 1; overflow-y: auto; padding: 12px 16px; }
.dp-section { margin-bottom: 16px; }
.dp-section-title {
  font-size: 10px; font-weight: 600; letter-spacing: 0.06em;
  color: var(--text-muted); margin-bottom: 8px; text-transform: uppercase;
}

/* Status row */
.dp-status-row { display: flex; align-items: center; gap: 8px; font-size: 12px; }
.dp-status-dot { width: 8px; height: 8px; border-radius: 50%; }

/* Task card */
.dp-task {
  padding: 8px 10px; background: var(--bg-raised); border: 1px solid var(--border-subtle);
  border-radius: 6px; margin-bottom: 6px;
}
.dp-task-title { font-size: 12px; font-weight: 500; margin-bottom: 4px; }
.dp-task-bar { height: 4px; background: var(--bg-overlay); border-radius: 2px; overflow: hidden; }
.dp-task-fill { height: 100%; border-radius: 2px; transition: width 0.5s; }
.dp-task-meta { font-size: 10px; color: var(--text-muted); margin-top: 3px; display: flex; justify-content: space-between; }

/* Message item */
.dp-msg {
  padding: 6px 0; border-bottom: 1px solid var(--border-subtle); font-size: 12px;
}
.dp-msg-header { display: flex; justify-content: space-between; margin-bottom: 2px; }
.dp-msg-from { font-weight: 500; }
.dp-msg-time { font-family: var(--font-mono); font-size: 10px; color: var(--text-muted); }
.dp-msg-text { color: var(--text-secondary); line-height: 1.5; max-height: 80px; overflow: hidden; position: relative; }
.dp-msg-text code { font-family: var(--font-mono); font-size: 10px; background: var(--bg-overlay); padding: 1px 3px; border-radius: 2px; }
.dp-msg-text strong { color: var(--text-primary); }
.dp-msg-text em { font-style: italic; }
.dp-msg-text ul { margin: 2px 0; padding-left: 14px; }
.dp-msg-text h3,.dp-msg-text h4,.dp-msg-text h5 { font-size: 11px; color: var(--text-primary); margin: 2px 0; }
.dp-msg-more {
  font-size: 10px; color: var(--accent); cursor: pointer; background: none;
  border: none; padding: 2px 0; font-family: var(--font-ui); margin-top: 2px;
}
.dp-msg-more:hover { text-decoration: underline; }

/* Content Viewer Modal */
.viewer-overlay {
  position: fixed; inset: 0; background: rgba(8,11,20,0.85);
  backdrop-filter: blur(12px); z-index: 300;
  display: flex; justify-content: center; align-items: center; padding: 40px;
}
.viewer-box {
  width: 720px; max-width: 90vw; max-height: 80vh; background: var(--bg-raised);
  border: 1px solid var(--border-default); border-radius: 10px;
  display: flex; flex-direction: column; box-shadow: 0 16px 48px rgba(0,0,0,0.5);
}
.viewer-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 12px 18px; border-bottom: 1px solid var(--border-subtle);
  font-size: 13px; font-weight: 600; flex-shrink: 0;
}
.viewer-body {
  flex: 1; overflow-y: auto; padding: 18px;
  font-size: 13px; color: var(--text-secondary); line-height: 1.7;
}
.viewer-body h3 { font-size: 16px; color: var(--text-primary); margin: 12px 0 6px; }
.viewer-body h4 { font-size: 14px; color: var(--text-primary); margin: 10px 0 4px; }
.viewer-body h5 { font-size: 13px; color: var(--text-primary); margin: 8px 0 4px; }
.viewer-body pre { background: var(--bg-overlay); padding: 10px; border-radius: 6px; overflow-x: auto; margin: 8px 0; font-size: 11px; }
.viewer-body code { font-family: var(--font-mono); font-size: 11px; background: var(--bg-overlay); padding: 1px 4px; border-radius: 3px; }
.viewer-body pre code { background: none; padding: 0; }
.viewer-body table { border-collapse: collapse; margin: 8px 0; width: 100%; }
.viewer-body td, .viewer-body th { border: 1px solid var(--border-subtle); padding: 4px 8px; font-size: 11px; }
.viewer-body ul { padding-left: 18px; margin: 4px 0; }
.viewer-body li { margin: 2px 0; }
.viewer-body strong { color: var(--text-primary); }

/* ═══════════════════════════════════════════
   CHAT TAB
   ═══════════════════════════════════════════ */
.dp-chat { display: flex; flex-direction: column; height: 100%; }
.dp-chat-messages {
  flex: 1; overflow-y: auto; padding: 8px 0;
  display: flex; flex-direction: column; gap: 8px;
}
.chat-bubble {
  max-width: 85%; padding: 8px 12px; border-radius: 8px;
  font-size: 12px; line-height: 1.5; word-wrap: break-word;
}
.chat-bubble.agent {
  align-self: flex-start; background: var(--bg-raised);
  border: 1px solid var(--border-subtle); border-bottom-left-radius: 2px;
}
.chat-bubble.user {
  align-self: flex-end; background: rgba(59,130,246,0.15);
  border: 1px solid rgba(59,130,246,0.2); border-bottom-right-radius: 2px;
}
.chat-bubble .cb-from { font-size: 10px; font-weight: 600; margin-bottom: 2px; }
.chat-bubble .cb-time { font-size: 9px; color: var(--text-muted); margin-top: 3px; }
.chat-bubble .cb-text { color: var(--text-secondary); }
.chat-bubble .cb-text code { font-family: var(--font-mono); font-size: 10px; background: var(--bg-overlay); padding: 1px 3px; border-radius: 2px; }
.chat-bubble .cb-text strong { color: var(--text-primary); }
.chat-bubble.typing .cb-text::after { content: '...'; animation: typing-dots 1s step-end infinite; }
@keyframes typing-dots { 33% { content: '.'; } 66% { content: '..'; } }

.dp-chat-input {
  display: flex; gap: 6px; padding: 8px 0 0;
  border-top: 1px solid var(--border-subtle); margin-top: 4px;
}
.chat-input {
  flex: 1; background: var(--bg-raised); border: 1px solid var(--border-subtle);
  color: var(--text-primary); font-family: var(--font-ui); font-size: 12px;
  padding: 6px 10px; border-radius: 6px; outline: none; resize: none; min-height: 32px; max-height: 80px;
}
.chat-input:focus { border-color: var(--accent); }
.chat-input::placeholder { color: var(--text-muted); }
.chat-send {
  background: var(--accent); color: #fff; border: none; padding: 6px 12px;
  border-radius: 6px; font-family: var(--font-ui); font-size: 11px;
  font-weight: 600; cursor: pointer; align-self: flex-end;
}
.chat-send:hover { background: #4B92FF; }
.chat-send:disabled { opacity: 0.4; cursor: not-allowed; }

/* ═══════════════════════════════════════════
   TOAST NOTIFICATIONS
   ═══════════════════════════════════════════ */
.toast-container {
  position: fixed; top: 60px; right: 16px; z-index: 400;
  display: flex; flex-direction: column; gap: 8px; pointer-events: none;
}
.toast {
  background: var(--bg-raised); border: 1px solid var(--border-default);
  border-radius: 8px; padding: 10px 16px; min-width: 280px; max-width: 400px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.4); pointer-events: auto;
  display: flex; align-items: flex-start; gap: 10px;
  animation: toast-in 0.3s ease;
  border-left: 3px solid var(--accent);
}
.toast.success { border-left-color: var(--st-complete); }
.toast.warning { border-left-color: #F59E0B; }
.toast.error { border-left-color: var(--st-error); }
@keyframes toast-in { from { transform: translateX(120%); opacity: 0; } to { transform: translateX(0); opacity: 1; } }
@keyframes toast-out { to { transform: translateX(120%); opacity: 0; } }
.toast-icon { font-size: 16px; flex-shrink: 0; line-height: 1; }
.toast-body { flex: 1; }
.toast-title { font-size: 11px; font-weight: 600; color: var(--text-primary); margin-bottom: 2px; }
.toast-text { font-size: 11px; color: var(--text-secondary); line-height: 1.4; }

/* ═══════════════════════════════════════════
   LOG STREAM (Bottom)
   ═══════════════════════════════════════════ */
.log-area {
  grid-area: log;
  background: var(--bg-base);
  border-top: 1px solid var(--border-subtle);
  display: flex; flex-direction: column;
  max-height: var(--log-h);
  transition: max-height 0.25s ease;
}
.log-area.collapsed { max-height: var(--log-collapsed); }

.log-toolbar {
  display: flex; align-items: center; gap: 10px;
  padding: 6px 16px; border-bottom: 1px solid var(--border-subtle);
  flex-shrink: 0;
}
.log-toolbar-title {
  font-size: 10px; font-weight: 600; letter-spacing: 0.08em;
  color: var(--text-muted); text-transform: uppercase; cursor: pointer;
  user-select: none;
}
.log-toolbar-title:hover { color: var(--text-secondary); }
.log-tab { padding: 4px 0; border-bottom: 2px solid transparent; margin-right: 4px; }
.log-tab.active { color: var(--text-secondary); border-bottom-color: var(--accent); }

/* Project panel in log area */
.proj-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; padding: 12px; font-size: 11px; }
.proj-card { background: var(--bg-overlay); border-radius: 6px; padding: 10px; }
.proj-card-title { font-size: 10px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: var(--text-muted); margin-bottom: 8px; }
.proj-file { display: flex; gap: 6px; padding: 3px 0; cursor: pointer; color: var(--text-secondary); font-family: var(--font-mono); font-size: 10px; }
.proj-file:hover { color: var(--accent); }
.proj-agent { display: flex; gap: 8px; padding: 3px 0; font-size: 10px; }
.proj-agent-dot { width: 8px; height: 8px; border-radius: 50%; margin-top: 2px; flex-shrink: 0; }
.proj-agent-name { width: 80px; font-weight: 500; color: var(--text-secondary); }
.proj-agent-val { color: var(--text-muted); font-family: var(--font-mono); }

/* Task input inline */
.task-input-row { display: flex; align-items: center; gap: 8px; flex: 1; }
.task-input {
  flex: 1; background: var(--bg-raised); border: 1px solid var(--border-subtle);
  color: var(--text-primary); font-family: var(--font-ui); font-size: 12px;
  padding: 5px 10px; border-radius: 4px; outline: none;
}
.task-input:focus { border-color: var(--accent); }
.task-input::placeholder { color: var(--text-muted); }
.task-submit {
  background: var(--accent); color: #fff; border: none;
  font-family: var(--font-ui); font-size: 11px; font-weight: 600;
  padding: 6px 14px; border-radius: 4px; cursor: pointer; white-space: nowrap;
  transition: background 0.15s;
}
.task-submit:hover { background: #4B92FF; }
.task-submit:disabled { opacity: 0.4; cursor: not-allowed; }

.log-stream { flex: 1; overflow-y: auto; padding: 4px 16px; font-family: var(--font-mono); font-size: 11px; }
.log-entry {
  padding: 2px 0; display: flex; gap: 10px; line-height: 1.5;
  border-bottom: 1px solid rgba(27,33,48,0.5);
  animation: log-in 0.15s ease;
}
@keyframes log-in {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}
.log-time { color: var(--text-muted); min-width: 64px; flex-shrink: 0; }
.log-who { font-weight: 500; min-width: 60px; flex-shrink: 0; }
.log-text { color: var(--text-secondary); flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* ═══════════════════════════════════════════
   CHECKPOINT BANNER
   ═══════════════════════════════════════════ */
.checkpoint-banner {
  position: absolute; top: 0; left: 0; right: 0; z-index: 20;
  background: rgba(13,17,23,0.95); backdrop-filter: blur(12px);
  border-bottom: 2px solid var(--accent);
  padding: 20px 28px; display: flex; flex-direction: column; gap: 12px;
  animation: cp-slide-in 0.3s ease;
}
@keyframes cp-slide-in { from { transform: translateY(-100%); } to { transform: translateY(0); } }
.cp-badge {
  display: inline-flex; align-items: center; gap: 6px;
  font-size: 10px; font-weight: 600; letter-spacing: 0.08em;
  color: #F59E0B; text-transform: uppercase;
}
.cp-badge::before {
  content: ''; width: 8px; height: 8px; border-radius: 50%;
  background: #F59E0B; animation: pulse-dot 1.5s ease-in-out infinite;
}
.cp-title { font-size: 16px; font-weight: 600; color: var(--text-primary); }
.cp-summary { font-size: 12px; color: var(--text-secondary); line-height: 1.5; }
.cp-countdown {
  font-family: var(--font-mono); font-size: 11px; color: #F59E0B;
  padding: 4px 10px; background: rgba(245,158,11,0.1);
  border: 1px solid rgba(245,158,11,0.2); border-radius: 4px; width: fit-content;
}
.cp-content-wrap {
  max-height: 200px; overflow-y: auto; padding: 10px 12px;
  background: var(--bg-raised); border: 1px solid var(--border-subtle);
  border-radius: 6px; font-size: 11px;
  color: var(--text-secondary); line-height: 1.6;
}
.cp-content-wrap h3, .cp-content-wrap h4, .cp-content-wrap h5 { color: var(--text-primary); margin: 8px 0 4px; }
.cp-content-wrap h3 { font-size: 14px; } .cp-content-wrap h4 { font-size: 12px; } .cp-content-wrap h5 { font-size: 11px; }
.cp-content-wrap pre { background: var(--bg-overlay); padding: 8px; border-radius: 4px; overflow-x: auto; margin: 6px 0; }
.cp-content-wrap code { font-family: var(--font-mono); font-size: 10px; background: var(--bg-overlay); padding: 1px 4px; border-radius: 3px; }
.cp-content-wrap pre code { background: none; padding: 0; }
.cp-content-wrap table { border-collapse: collapse; margin: 6px 0; width: 100%; font-size: 10px; }
.cp-content-wrap td { border: 1px solid var(--border-subtle); padding: 3px 6px; }
.cp-content-wrap ul { margin: 4px 0; padding-left: 16px; }
.cp-content-wrap li { margin: 2px 0; }
.cp-content-wrap strong { color: var(--text-primary); }
.cp-toggle { font-size: 11px; color: var(--accent); background: none; border: none; cursor: pointer; padding: 0; }
.cp-toggle:hover { text-decoration: underline; }
.cp-actions { display: flex; gap: 10px; margin-top: 4px; }
.cp-btn {
  padding: 8px 18px; font-family: var(--font-ui); font-size: 12px;
  font-weight: 600; border-radius: 5px; cursor: pointer; border: none;
  transition: background 0.15s, transform 0.1s;
}
.cp-btn:active { transform: scale(0.97); }
.cp-btn.approve { background: var(--st-complete); color: #fff; }
.cp-btn.approve:hover { background: #2ad46a; }
.cp-btn.reject { background: var(--st-error); color: #fff; }
.cp-btn.reject:hover { background: #f55; }
.cp-btn.override { background: #F59E0B; color: #000; }
.cp-btn.override:hover { background: #fbba3c; }

/* ═══════════════════════════════════════════
   SETTINGS & FEEDBACK MODALS
   ═══════════════════════════════════════════ */
.modal-overlay {
  position: fixed; inset: 0; background: rgba(8,11,20,0.7);
  backdrop-filter: blur(8px); z-index: 200;
  display: flex; justify-content: center; padding-top: 15vh;
}
.modal-box {
  width: 480px; max-height: 70vh; background: var(--bg-raised);
  border: 1px solid var(--border-default); border-radius: 10px;
  overflow: hidden; box-shadow: 0 16px 48px rgba(0,0,0,0.4);
  display: flex; flex-direction: column;
}
.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 14px 18px; border-bottom: 1px solid var(--border-subtle);
  font-size: 14px; font-weight: 600;
}
.modal-close { background: none; border: none; color: var(--text-muted); font-size: 18px; cursor: pointer; }
.modal-close:hover { color: var(--text-primary); }
.modal-body { flex: 1; overflow-y: auto; padding: 16px 18px; }
.modal-footer { padding: 12px 18px; border-top: 1px solid var(--border-subtle); display: flex; justify-content: flex-end; gap: 8px; }
.modal-btn {
  padding: 7px 16px; font-family: var(--font-ui); font-size: 12px; font-weight: 500;
  border-radius: 5px; cursor: pointer; border: 1px solid var(--border-subtle);
  background: transparent; color: var(--text-secondary);
}
.modal-btn:hover { border-color: var(--border-active); color: var(--text-primary); }
.modal-btn.primary { background: var(--accent); color: #fff; border-color: transparent; }
.modal-btn.primary:hover { background: #4B92FF; }

/* Settings toggles */
.setting-section { margin-bottom: 16px; }
.setting-section-title { font-size: 10px; font-weight: 600; letter-spacing: 0.06em; color: var(--text-muted); text-transform: uppercase; margin-bottom: 8px; }
.setting-row {
  display: flex; justify-content: space-between; align-items: center;
  padding: 6px 0; font-size: 12px; color: var(--text-secondary);
}
.toggle-switch {
  position: relative; width: 36px; height: 20px; background: var(--border-subtle);
  border-radius: 10px; cursor: pointer; transition: background 0.2s;
}
.toggle-switch.on { background: var(--accent); }
.toggle-switch::after {
  content: ''; position: absolute; top: 2px; left: 2px;
  width: 16px; height: 16px; background: #fff; border-radius: 50%;
  transition: transform 0.2s;
}
.toggle-switch.on::after { transform: translateX(16px); }
.setting-select {
  background: var(--bg-base); border: 1px solid var(--border-subtle);
  color: var(--text-primary); font-family: var(--font-mono); font-size: 11px;
  padding: 3px 6px; border-radius: 4px;
}

/* Feedback textarea */
.feedback-textarea {
  width: 100%; min-height: 100px; background: var(--bg-base);
  border: 1px solid var(--border-subtle); border-radius: 6px;
  color: var(--text-primary); font-family: var(--font-ui); font-size: 13px;
  padding: 10px; outline: none; resize: vertical;
}
.feedback-textarea:focus { border-color: var(--accent); }

/* ═══════════════════════════════════════════
   VIEW TOGGLE
   ═══════════════════════════════════════════ */
.view-toggle { display: flex; gap: 2px; background: var(--bg-raised); border-radius: 5px; padding: 2px; border: 1px solid var(--border-subtle); }
.view-btn {
  padding: 3px 10px; font-family: var(--font-ui); font-size: 10px; font-weight: 500;
  color: var(--text-muted); background: transparent; border: none; border-radius: 3px;
  cursor: pointer; transition: all 0.15s;
}
.view-btn:hover { color: var(--text-secondary); }
.view-btn.active { background: var(--accent); color: #fff; }

/* ═══════════════════════════════════════════
   KANBAN BOARD
   ═══════════════════════════════════════════ */
.kanban-view { display: none; grid-area: canvas; overflow: hidden; flex-direction: column; }
.kanban-view.active { display: flex; }
.topology-view { display: contents; }
.topology-view.hidden .roster { display: none; }
.topology-view.hidden .canvas-area { display: none; }

.kanban-toolbar {
  display: flex; align-items: center; gap: 8px; padding: 8px 16px;
  border-bottom: 1px solid var(--border-subtle); background: var(--bg-base); flex-shrink: 0;
}
.kanban-filter {
  background: var(--bg-raised); border: 1px solid var(--border-subtle);
  color: var(--text-primary); font-family: var(--font-ui); font-size: 11px;
  padding: 4px 8px; border-radius: 4px;
}
.kanban-search {
  flex: 1; max-width: 200px; background: var(--bg-raised); border: 1px solid var(--border-subtle);
  color: var(--text-primary); font-family: var(--font-ui); font-size: 11px;
  padding: 4px 8px; border-radius: 4px; outline: none;
}
.kanban-search:focus { border-color: var(--accent); }
.kanban-search::placeholder { color: var(--text-muted); }
.kanban-sync-btn {
  background: transparent; border: 1px solid var(--border-subtle); color: var(--text-secondary);
  font-family: var(--font-ui); font-size: 10px; padding: 4px 8px; border-radius: 4px; cursor: pointer;
}
.kanban-sync-btn:hover { border-color: var(--accent); color: var(--accent); }

.kanban-board {
  display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px;
  padding: 12px 16px; flex: 1; overflow-y: auto; align-content: start;
}
.kanban-col {
  background: var(--bg-base); border: 1px solid var(--border-subtle);
  border-radius: 8px; display: flex; flex-direction: column;
  min-height: 200px; max-height: calc(100vh - 260px);
}
.kanban-col-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 10px 12px; border-bottom: 1px solid var(--border-subtle);
  font-size: 11px; font-weight: 600; letter-spacing: 0.04em;
  color: var(--text-secondary); text-transform: uppercase;
}
.kanban-col-count {
  background: var(--bg-raised); padding: 1px 6px; border-radius: 8px;
  font-size: 10px; font-weight: 500; color: var(--text-muted);
}
.kanban-col-body {
  flex: 1; overflow-y: auto; padding: 8px; display: flex;
  flex-direction: column; gap: 6px; min-height: 60px;
}
.kanban-col-body.drag-over { background: rgba(59,130,246,0.05); border-radius: 0 0 8px 8px; }

/* Kanban Card */
.kanban-card {
  background: var(--bg-raised); border: 1px solid var(--border-subtle);
  border-radius: 6px; padding: 10px; cursor: grab;
  border-left: 3px solid var(--text-muted);
  transition: transform 0.15s, box-shadow 0.15s, opacity 0.15s;
}
.kanban-card:hover { border-color: var(--border-active); transform: translateY(-1px); box-shadow: 0 4px 12px rgba(0,0,0,0.2); }
.kanban-card.dragging { opacity: 0.5; }
.kanban-card.blocked { border-left-style: dashed; opacity: 0.6; }

.kc-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 6px; }
.kc-id { font-family: var(--font-mono); font-size: 9px; color: var(--text-muted); }
.kc-agent {
  width: 20px; height: 20px; border-radius: 50%; display: flex;
  align-items: center; justify-content: center; font-size: 8px;
  font-weight: 700; border: 1px solid;
}
.kc-title { font-size: 12px; font-weight: 500; line-height: 1.4; margin-bottom: 6px; color: var(--text-primary); }
.kc-bar { height: 3px; background: var(--bg-overlay); border-radius: 2px; overflow: hidden; margin-bottom: 4px; }
.kc-bar-fill { height: 100%; border-radius: 2px; transition: width 0.5s; }
.kc-footer { display: flex; justify-content: space-between; font-size: 9px; color: var(--text-muted); }
.kc-priority {
  padding: 1px 5px; border-radius: 3px; font-size: 8px; font-weight: 600;
  letter-spacing: 0.04em; text-transform: uppercase;
}
.kc-priority.critical { background: rgba(239,68,68,0.15); color: #EF4444; }
.kc-priority.high { background: rgba(245,158,11,0.15); color: #F59E0B; }
.kc-priority.medium { background: rgba(59,130,246,0.15); color: #3B82F6; }
.kc-priority.low { background: rgba(72,79,88,0.15); color: #8B949E; }

/* Add task button */
.kanban-add-btn {
  width: 100%; padding: 6px; background: transparent; border: 1px dashed var(--border-subtle);
  border-radius: 4px; color: var(--text-muted); font-size: 11px; cursor: pointer;
  font-family: var(--font-ui);
}
.kanban-add-btn:hover { border-color: var(--accent); color: var(--accent); }

/* Inline card editor */
.kanban-new-card {
  background: var(--bg-raised); border: 1px solid var(--accent);
  border-radius: 6px; padding: 10px; display: flex; flex-direction: column; gap: 6px;
}
.kanban-new-card input, .kanban-new-card select, .kanban-new-card textarea {
  background: var(--bg-base); border: 1px solid var(--border-subtle);
  color: var(--text-primary); font-family: var(--font-ui); font-size: 11px;
  padding: 4px 8px; border-radius: 4px; outline: none;
}
.kanban-new-card textarea { min-height: 40px; resize: vertical; }
.kanban-new-card input:focus, .kanban-new-card textarea:focus { border-color: var(--accent); }
.kanban-card-actions { display: flex; gap: 6px; justify-content: flex-end; }
.kanban-card-actions button {
  padding: 4px 10px; font-family: var(--font-ui); font-size: 10px;
  border-radius: 3px; cursor: pointer; border: 1px solid var(--border-subtle);
  background: transparent; color: var(--text-secondary);
}
.kanban-card-actions .save { background: var(--accent); color: #fff; border-color: transparent; }

/* ═══════════════════════════════════════════
   COMMAND PALETTE
   ═══════════════════════════════════════════ */
.cmd-overlay {
  position: fixed; inset: 0; background: rgba(8,11,20,0.7);
  backdrop-filter: blur(8px); z-index: 200;
  display: flex; justify-content: center; padding-top: 20vh;
}
.cmd-box {
  width: 520px; max-height: 360px; background: var(--bg-raised);
  border: 1px solid var(--border-default); border-radius: 10px;
  overflow: hidden; box-shadow: 0 16px 48px rgba(0,0,0,0.4);
  display: flex; flex-direction: column;
}
.cmd-input {
  width: 100%; padding: 14px 18px; background: transparent; border: none;
  border-bottom: 1px solid var(--border-subtle);
  color: var(--text-primary); font-family: var(--font-mono); font-size: 14px;
  outline: none;
}
.cmd-input::placeholder { color: var(--text-muted); }
.cmd-results { flex: 1; overflow-y: auto; padding: 6px 0; }
.cmd-item {
  padding: 8px 18px; font-size: 13px; cursor: pointer;
  display: flex; align-items: center; gap: 10px;
  color: var(--text-secondary);
}
.cmd-item:hover, .cmd-item.active { background: var(--bg-overlay); color: var(--text-primary); }
.cmd-item-icon { font-size: 14px; width: 20px; text-align: center; }
.cmd-item-label { flex: 1; }
.cmd-item-hint { font-family: var(--font-mono); font-size: 10px; color: var(--text-muted); }
.cmd-hidden { display: none; }

/* ═══════════════════════════════════════════
   SVG ARC ANIMATIONS
   ═══════════════════════════════════════════ */
.arc-path {
  fill: none; stroke-width: 2; opacity: 0;
  stroke-dasharray: 600; stroke-dashoffset: 600;
}
.arc-path.active {
  opacity: 0.7;
  animation: arc-flow 0.6s ease forwards;
}
.arc-path.trail {
  opacity: 0.15;
  stroke-dashoffset: 0;
  transition: opacity 8s linear;
}
.arc-path.fading { opacity: 0; }

@keyframes arc-flow {
  to { stroke-dashoffset: 0; opacity: 0.7; }
}

/* ═══════════════════════════════════════════
   SCROLLBAR
   ═══════════════════════════════════════════ */
::-webkit-scrollbar { width: 5px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: var(--border-subtle); border-radius: 3px; }

/* ═══════════════════════════════════════════
   RESPONSIVE
   ═══════════════════════════════════════════ */
@media (max-width: 1200px) {
  :root { --roster-w: 56px; --panel-w: 300px; }
  .roster-info, .roster-title { display: none; }
  .roster-card { justify-content: center; padding: 8px; }
  .kpi-strip { display: none; }
}
@media (max-width: 900px) {
  .app { grid-template-columns: 0 1fr 0; }
  .roster { display: none; }
  .pipeline { display: none; }
}
  </style>
</head>
<body>
<div class="app" id="app">

  <!-- ════════ HEADER ════════ -->
  <div class="header">
    <div class="logo">Agent<span>House</span></div>

    <div class="pipeline" id="pipeline"></div>

    <div class="kpi-strip" id="kpiStrip">
      <div class="kpi"><div class="kpi-dot live" id="wsDot"></div> <span class="kpi-val" id="kpiLive">LIVE</span></div>
      <div class="kpi">Phase <span class="kpi-val" id="kpiPhase">-</span></div>
      <div class="kpi">Agents <span class="kpi-val" id="kpiAgents">0/8</span></div>
      <div class="kpi">Msgs <span class="kpi-val" id="kpiMsgs">0</span></div>
      <div class="kpi">Time <span class="kpi-val" id="kpiTime">00:00</span></div>
    </div>

    <div class="view-toggle" id="viewToggle">
      <button class="view-btn active" data-view="topology">Topology</button>
      <button class="view-btn" data-view="kanban">Kanban</button>
    </div>

    <div class="hdr-controls">
      <select class="project-sel" id="projectSel"><option value="default">default</option></select>
      <button class="hdr-btn" id="newProjBtn">+ New</button>
      <button class="hdr-btn" id="settingsBtn" title="Workflow Settings">&#9881;</button>
      <button class="hdr-btn" id="cmdBtn" title="Cmd+K">&#8984;K</button>
      <a class="hdr-btn" href="/live">Live</a>
      <a class="hdr-btn" href="/office">8-Bit</a>
      <a class="hdr-btn" href="/">Dashboard</a>
    </div>
  </div>

  <!-- ════════ ROSTER ════════ -->
  <div class="roster" id="roster">
    <div class="roster-title">Team</div>
  </div>

  <!-- ════════ CANVAS ════════ -->
  <div class="canvas-area" id="canvasArea">
    <!-- Checkpoint Banner -->
    <div class="checkpoint-banner" id="cpBanner" style="display:none">
      <div class="cp-badge">Waiting for Approval</div>
      <div class="cp-title" id="cpTitle"></div>
      <div class="cp-summary" id="cpSummary"></div>
      <div class="cp-countdown" id="cpCountdown" style="display:none"></div>
      <button class="cp-toggle" id="cpToggle">Show details</button>
      <div class="cp-content-wrap" id="cpContent" style="display:none"></div>
      <div class="cp-actions">
        <button class="cp-btn approve" id="cpApprove">Approve &amp; Continue</button>
        <button class="cp-btn reject" id="cpReject">Reject with Feedback</button>
        <button class="cp-btn override" id="cpOverride">Override</button>
      </div>
    </div>
    <svg class="arc-layer" id="arcLayer"></svg>
    <div id="nodesLayer"></div>
  </div>

  <!-- ════════ KANBAN VIEW ════════ -->
  <div class="kanban-view" id="kanbanView">
    <div class="kanban-toolbar">
      <select class="kanban-filter" id="kbFilterAgent"><option value="">All Agents</option></select>
      <select class="kanban-filter" id="kbFilterPhase"><option value="">All Phases</option></select>
      <select class="kanban-filter" id="kbFilterPriority"><option value="">All Priorities</option><option value="critical">Critical</option><option value="high">High</option><option value="medium">Medium</option><option value="low">Low</option></select>
      <input class="kanban-search" id="kbSearch" placeholder="Search tasks...">
      <button class="kanban-sync-btn" id="kbSyncBtn" title="Sync from development plan">Sync Plan</button>
    </div>
    <div class="kanban-board" id="kanbanBoard">
      <div class="kanban-col" data-status="backlog">
        <div class="kanban-col-header">Backlog <span class="kanban-col-count" id="kbCountBacklog">0</span></div>
        <div class="kanban-col-body" id="kbColBacklog"></div>
      </div>
      <div class="kanban-col" data-status="in_progress">
        <div class="kanban-col-header">In Progress <span class="kanban-col-count" id="kbCountInProgress">0</span></div>
        <div class="kanban-col-body" id="kbColInProgress"></div>
      </div>
      <div class="kanban-col" data-status="in_review">
        <div class="kanban-col-header">In Review <span class="kanban-col-count" id="kbCountInReview">0</span></div>
        <div class="kanban-col-body" id="kbColInReview"></div>
      </div>
      <div class="kanban-col" data-status="done">
        <div class="kanban-col-header">Done <span class="kanban-col-count" id="kbCountDone">0</span></div>
        <div class="kanban-col-body" id="kbColDone"></div>
      </div>
    </div>
  </div>

  <!-- ════════ DETAIL PANEL ════════ -->
  <div class="detail-panel" id="detailPanel">
    <div class="dp-header">
      <div class="dp-avatar" id="dpAvatar"></div>
      <div class="dp-info"><div class="dp-name" id="dpName"></div><div class="dp-role" id="dpRole"></div></div>
      <button class="dp-close" id="dpClose">&times;</button>
    </div>
    <div class="dp-tabs">
      <div class="dp-tab active" data-tab="activity">Activity</div>
      <div class="dp-tab" data-tab="tasks">Tasks</div>
      <div class="dp-tab" data-tab="chat">Chat</div>
      <div class="dp-tab" data-tab="metrics">Metrics</div>
    </div>
    <div class="dp-body" id="dpBody"></div>
  </div>

  <!-- ════════ LOG ════════ -->
  <div class="log-area" id="logArea">
    <div class="log-toolbar">
      <span class="log-toolbar-title log-tab active" data-logtab="log" id="logToggle">Activity Log</span>
      <span class="log-toolbar-title log-tab" data-logtab="project" id="projectToggle">Project</span>
      <div class="task-input-row">
        <input type="text" class="task-input" id="taskInput" placeholder="Describe what you want to build...">
        <button class="task-submit" id="taskSubmit">Launch</button>
        <button class="task-submit" id="injectBtn" style="background:var(--c-architect);font-size:10px;padding:5px 10px" title="Inject a task into the running pipeline">&#x1F4CC; Inject</button>
      </div>
    </div>
    <div class="log-stream" id="logStream"></div>
    <div class="log-stream" id="projectPanel" style="display:none"></div>
  </div>
</div>

<!-- ════════ INJECT TASK DIALOG ════════ -->
<div class="cmd-overlay cmd-hidden" id="injectOverlay" onclick="if(event.target===this)closeInject()">
  <div class="cmd-box" style="max-width:500px;padding:20px">
    <div style="font-size:14px;font-weight:600;margin-bottom:16px">&#x1F4CC; Inject Task</div>
    <div style="font-size:11px;color:var(--text-secondary);margin-bottom:12px">Send a task directly to an agent during the current pipeline run. It will execute between phases.</div>
    <textarea id="injectTaskInput" style="width:100%;min-height:80px;background:var(--bg-raised);border:1px solid var(--border-default);color:var(--text-primary);padding:10px;border-radius:6px;font-family:var(--font-ui);font-size:12px;resize:vertical;margin-bottom:12px" placeholder="What should the agent do?"></textarea>
    <div style="display:flex;gap:10px;align-items:center;margin-bottom:16px">
      <label style="font-size:11px;color:var(--text-muted)">Agent:</label>
      <select id="injectAgentSel" style="background:var(--bg-raised);border:1px solid var(--border-default);color:var(--text-primary);padding:4px 8px;border-radius:4px;font-size:11px">
        <option value="senior_dev">Senior Dev</option>
        <option value="junior_dev">Junior Dev</option>
        <option value="architect">Architect</option>
        <option value="pm">PM</option>
        <option value="ceo">CEO</option>
        <option value="ux">UX</option>
        <option value="ui">UI</option>
        <option value="security">Security</option>
      </select>
      <label style="font-size:11px;color:var(--text-muted);margin-left:8px">Priority:</label>
      <select id="injectPrioritySel" style="background:var(--bg-raised);border:1px solid var(--border-default);color:var(--text-primary);padding:4px 8px;border-radius:4px;font-size:11px">
        <option value="normal">Normal</option>
        <option value="high">High</option>
      </select>
    </div>
    <div style="display:flex;gap:8px;justify-content:flex-end">
      <button onclick="closeInject()" style="background:transparent;border:1px solid var(--border-default);color:var(--text-secondary);padding:6px 16px;border-radius:6px;font-size:12px;cursor:pointer">Cancel</button>
      <button onclick="submitInject()" style="background:var(--c-architect);color:#fff;border:none;padding:6px 16px;border-radius:6px;font-size:12px;font-weight:600;cursor:pointer">Inject Task</button>
    </div>
    <div id="injectStatus" style="font-size:11px;color:var(--text-muted);margin-top:8px"></div>
  </div>
</div>

<!-- ════════ COMMAND PALETTE ════════ -->
<div class="cmd-overlay cmd-hidden" id="cmdOverlay">
  <div class="cmd-box">
    <input class="cmd-input" id="cmdInput" placeholder="Search agents, tasks, or type a command...">
    <div class="cmd-results" id="cmdResults"></div>
  </div>
</div>

<!-- Settings Modal -->
<div class="modal-overlay" id="settingsOverlay" style="display:none">
  <div class="modal-box">
    <div class="modal-header"><span>Workflow Settings</span><button class="modal-close" id="settingsClose">&times;</button></div>
    <div class="modal-body" id="settingsBody"></div>
    <div class="modal-footer"><button class="modal-btn" id="settingsCancel">Cancel</button><button class="modal-btn primary" id="settingsSave">Save</button></div>
  </div>
</div>

<!-- Feedback Modal -->
<div class="modal-overlay" id="feedbackOverlay" style="display:none">
  <div class="modal-box" style="width:420px">
    <div class="modal-header"><span id="feedbackTitle">Reject with Feedback</span><button class="modal-close" id="feedbackClose">&times;</button></div>
    <div class="modal-body">
      <textarea class="feedback-textarea" id="feedbackText" placeholder="What needs to change?"></textarea>
    </div>
    <div class="modal-footer"><button class="modal-btn" id="feedbackCancel">Cancel</button><button class="modal-btn primary" id="feedbackSubmit">Submit</button></div>
  </div>
</div>

<!-- Toast Container -->
<div class="toast-container" id="toastContainer"></div>

<!-- Content Viewer Modal -->
<div class="viewer-overlay" id="viewerOverlay" style="display:none">
  <div class="viewer-box">
    <div class="viewer-header"><span id="viewerTitle">Message</span><button class="modal-close" id="viewerClose">&times;</button></div>
    <div class="viewer-body" id="viewerBody"></div>
  </div>
</div>

<script>
// ═══════════════════════════════════════════
// AGENT DEFINITIONS
// ═══════════════════════════════════════════
const AGENTS = {
  ceo:         { key:'ceo',        name:'CEO',         full:'Chief Executive Officer', color:'#60B4D4', icon:'crown',    x:50,  y:10 },
  pm:          { key:'pm',         name:'PM',          full:'Product Manager',         color:'#8B7EF8', icon:'gantt',    x:30,  y:32 },
  architect:   { key:'architect',  name:'Architect',   full:'Software Architect',      color:'#A87CF5', icon:'blueprint',x:70,  y:32 },
  ux:          { key:'ux',         name:'UX',          full:'UX Designer',             color:'#F47F7F', icon:'cursor',   x:18,  y:56 },
  ui:          { key:'ui',         name:'UI',          full:'UI Designer',             color:'#4FD1C5', icon:'palette',  x:40,  y:56 },
  security:    { key:'security',   name:'Security',    full:'Security Expert',         color:'#F6A623', icon:'shield',   x:62,  y:56 },
  'senior-dev':{ key:'senior-dev', name:'Sr. Dev',     full:'Senior Developer',        color:'#34D478', icon:'terminal', x:35,  y:82 },
  'junior-dev':{ key:'junior-dev', name:'Jr. Dev',     full:'Junior Developer',        color:'#4BA3E3', icon:'code',     x:65,  y:82 },
};

// SVG icons per role (simple geometric shapes)
const ICONS = {
  crown:     '<svg viewBox="0 0 36 36"><path d="M6 26h24l-3-14-6 6-3-8-3 8-6-6z" fill="none" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/><rect x="6" y="26" width="24" height="3" rx="1" fill="currentColor" opacity="0.3"/></svg>',
  gantt:     '<svg viewBox="0 0 36 36"><rect x="8" y="9" width="14" height="3" rx="1.5" fill="currentColor" opacity="0.6"/><rect x="12" y="15" width="16" height="3" rx="1.5" fill="currentColor"/><rect x="8" y="21" width="10" height="3" rx="1.5" fill="currentColor" opacity="0.6"/><rect x="14" y="27" width="12" height="3" rx="1.5" fill="currentColor" opacity="0.4"/></svg>',
  cursor:    '<svg viewBox="0 0 36 36"><path d="M12 6v22l5-5h10z" fill="none" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/></svg>',
  palette:   '<svg viewBox="0 0 36 36"><rect x="8" y="8" width="9" height="9" rx="2" fill="currentColor" opacity="0.4"/><rect x="19" y="8" width="9" height="9" rx="2" fill="currentColor" opacity="0.7"/><rect x="8" y="19" width="9" height="9" rx="2" fill="currentColor" opacity="0.9"/><rect x="19" y="19" width="9" height="9" rx="2" fill="currentColor" opacity="0.3"/></svg>',
  shield:    '<svg viewBox="0 0 36 36"><path d="M18 4L6 10v8c0 8 5.4 13.4 12 16 6.6-2.6 12-8 12-16v-8z" fill="none" stroke="currentColor" stroke-width="2"/><path d="M14 18l3 3 6-6" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>',
  blueprint: '<svg viewBox="0 0 36 36"><rect x="6" y="6" width="24" height="24" rx="2" fill="none" stroke="currentColor" stroke-width="2"/><line x1="18" y1="6" x2="18" y2="30" stroke="currentColor" stroke-width="1.5" opacity="0.4"/><line x1="6" y1="18" x2="30" y2="18" stroke="currentColor" stroke-width="1.5" opacity="0.4"/><circle cx="18" cy="18" r="5" fill="none" stroke="currentColor" stroke-width="1.5"/></svg>',
  terminal:  '<svg viewBox="0 0 36 36"><rect x="5" y="8" width="26" height="20" rx="3" fill="none" stroke="currentColor" stroke-width="2"/><path d="M11 16l4 3-4 3" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="18" y1="22" x2="24" y2="22" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>',
  code:      '<svg viewBox="0 0 36 36"><path d="M14 12l-6 6 6 6" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M22 12l6 6-6 6" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>',
};

const PHASES = ['template_selection','research','planning','discussion','development'];
const PHASE_LABELS = { template_selection:'Template', research:'Research', planning:'Planning', discussion:'Discussion', development:'Development' };

// ═══════════════════════════════════════════
// STATE
// ═══════════════════════════════════════════
const S = {
  ws: null, projectId: 'default', running: false,
  agents: {}, messages: [], phase: '', selected: null,
  startTime: null, msgCount: 0, activeTab: 'activity',
  checkpoint: null, cpInterval: null, settings: null,
};

// ═══════════════════════════════════════════
// RENDER: PIPELINE
// ═══════════════════════════════════════════
function renderPipeline() {
  const el = document.getElementById('pipeline');
  const ci = PHASES.indexOf(S.phase);
  el.innerHTML = PHASES.map((p,i) => {
    const cls = i < ci ? 'done' : i === ci ? 'active' : '';
    const connCls = i < ci ? 'done' : i === ci ? 'active' : '';
    const dot = '<span class="dot"></span>';
    const step = '<div class="pipe-step '+cls+'">'+dot+PHASE_LABELS[p]+'</div>';
    return (i>0 ? '<div class="pipe-connector '+connCls+'"></div>' : '') + step;
  }).join('');
}

// ═══════════════════════════════════════════
// RENDER: ROSTER
// ═══════════════════════════════════════════
function renderRoster() {
  const el = document.getElementById('roster');
  el.innerHTML = '<div class="roster-title">Team</div>';
  Object.values(AGENTS).forEach(ag => {
    const as = S.agents[ag.key] || {};
    const state = as.state || 'idle';
    const stColor = state==='working'?'var(--st-working)':state==='completed'?'var(--st-complete)':state==='error'?'var(--st-error)':'var(--st-idle)';
    const sel = S.selected === ag.key ? ' selected' : '';
    const card = document.createElement('div');
    card.className = 'roster-card' + sel;
    if(state==='working') card.style.borderLeftColor = ag.color;
    const mode = (S.agentModes||{})[ag.key] || 'session';
    const modeBadge = mode==='oneshot' ? '<span class="mode-badge" style="font-size:8px;padding:1px 4px;border-radius:3px;background:rgba(246,166,35,0.2);color:#F6A623;margin-left:4px;cursor:pointer" title="Click to switch to Session mode" onclick="event.stopPropagation();toggleAgentMode(\''+ag.key+'\')">&#x26A1;API</span>'
      : '<span class="mode-badge" style="font-size:8px;padding:1px 4px;border-radius:3px;background:rgba(59,130,246,0.2);color:#3B82F6;margin-left:4px;cursor:pointer" title="Click to switch to Oneshot mode" onclick="event.stopPropagation();toggleAgentMode(\''+ag.key+'\')">&#x1F517;CC</span>';
    card.innerHTML =
      '<div class="roster-avatar" style="border-color:'+ag.color+';color:'+ag.color+'">'+
        ag.name.charAt(0)+
        '<div class="state-dot" style="background:'+stColor+'"></div>'+
      '</div>'+
      '<div class="roster-info">'+
        '<div class="roster-name" style="color:'+ag.color+'">'+ag.name+modeBadge+'</div>'+
        '<div class="roster-role">'+ag.full+'</div>'+
        (as.taskTitle ? '<div class="roster-task">'+esc(as.taskTitle)+'</div>' : '')+
      '</div>';
    card.addEventListener('click', () => selectAgent(ag.key));
    el.appendChild(card);
  });
}

// ═══════════════════════════════════════════
// RENDER: CANVAS NODES
// ═══════════════════════════════════════════
function renderNodes() {
  const layer = document.getElementById('nodesLayer');
  layer.innerHTML = '';
  Object.values(AGENTS).forEach(ag => {
    const as = S.agents[ag.key] || {};
    const state = as.state || 'idle';
    const node = document.createElement('div');
    node.className = 'canvas-node';
    node.dataset.agent = ag.key;
    node.dataset.state = state;
    node.style.left = ag.x + '%';
    node.style.top = ag.y + '%';
    if(S.selected && S.selected !== ag.key) node.classList.add('dimmed');
    node.innerHTML =
      '<div class="node-ring" style="border-color:'+(state==='idle'?'var(--st-idle)':ag.color)+'">'+
        '<div class="icon-shape" style="color:'+ag.color+'">'+ICONS[ag.icon]+'</div>'+
      '</div>'+
      '<div class="node-label">'+ag.name+'</div>'+
      '<div class="node-status">'+state+'</div>'+
      (as.taskTitle && state==='working' ? '<div class="node-task">'+esc(as.taskTitle)+'</div>' : '');
    node.addEventListener('click', (e) => { e.stopPropagation(); selectAgent(ag.key); });
    layer.appendChild(node);
  });
}

// ═══════════════════════════════════════════
// RENDER: KPIs
// ═══════════════════════════════════════════
function updateKPIs() {
  const active = Object.values(S.agents).filter(a => a.state==='working').length;
  document.getElementById('kpiAgents').textContent = active+'/8';
  document.getElementById('kpiMsgs').textContent = S.msgCount;
  document.getElementById('kpiPhase').textContent = S.phase ? PHASE_LABELS[S.phase]||S.phase : '-';
  if(S.startTime) {
    const s = Math.floor((Date.now()-S.startTime)/1000);
    const m = Math.floor(s/60), sec = s%60;
    document.getElementById('kpiTime').textContent = String(m).padStart(2,'0')+':'+String(sec).padStart(2,'0');
  }
}

// ═══════════════════════════════════════════
// DETAIL PANEL
// ═══════════════════════════════════════════
function selectAgent(key) {
  S.selected = key;
  document.getElementById('app').classList.add('panel-open');
  const ag = AGENTS[key]; if(!ag) return;
  const as = S.agents[key] || {};
  document.getElementById('dpAvatar').textContent = ag.name.charAt(0);
  document.getElementById('dpAvatar').style.borderColor = ag.color;
  document.getElementById('dpAvatar').style.color = ag.color;
  document.getElementById('dpName').textContent = ag.name;
  document.getElementById('dpName').style.color = ag.color;
  document.getElementById('dpRole').textContent = ag.full;
  renderPanelTab(S.activeTab);
  renderNodes(); renderRoster();
}

function closePanel() {
  S.selected = null;
  document.getElementById('app').classList.remove('panel-open');
  renderNodes(); renderRoster();
}

function renderPanelTab(tab) {
  S.activeTab = tab;
  document.querySelectorAll('.dp-tab').forEach(t => t.classList.toggle('active', t.dataset.tab===tab));
  const body = document.getElementById('dpBody');
  const key = S.selected; if(!key) return;
  const ag = AGENTS[key]; const as = S.agents[key] || {};

  if(tab==='activity') {
    // Show live session activity (tool calls, thinking, etc.)
    const activity = (as.activity || []).slice(-50);
    const toolCount = as.toolCount || 0;
    const stDot = as.state==='working'?'var(--st-working)':as.state==='completed'?'var(--st-complete)':as.state==='error'?'var(--st-error)':'var(--st-idle)';

    let html = '<div class="dp-section"><div class="dp-section-title">Status</div>'+
      '<div class="dp-status-row"><div class="dp-status-dot" style="background:'+stDot+'"></div>'+
      '<span style="text-transform:uppercase;letter-spacing:0.06em;font-size:12px">'+(as.state||'idle')+'</span>'+
      (toolCount ? '<span style="margin-left:12px;font-size:11px;color:var(--text-muted)">'+toolCount+' tool calls</span>':'')+
      '</div></div>';

    if(activity.length) {
      html += '<div class="dp-section"><div class="dp-section-title">Live Activity</div>';
      html += activity.map(a => {
        const time = new Date(a.t).toLocaleTimeString('en-US',{hour12:false,hour:'2-digit',minute:'2-digit',second:'2-digit'});
        return '<div style="display:flex;gap:8px;padding:4px 0;border-bottom:1px solid var(--border-subtle);font-size:11px;font-family:var(--font-mono)">'+
          '<span style="flex-shrink:0">'+a.icon+'</span>'+
          '<span style="flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;color:var(--text-secondary)">'+esc(a.text)+'</span>'+
          '<span style="flex-shrink:0;color:var(--text-muted);font-size:10px">'+time+'</span></div>';
      }).join('');
      html += '</div>';
    }

    // Also show message-based activity
    const msgs = S.messages.filter(m => normRole(m.from)===key || normRole(m.to)===key).slice(-10);
    if(msgs.length) {
      html += '<div class="dp-section"><div class="dp-section-title">Messages</div>'+
        msgs.map((m,idx) => {
          const content = m.content||'';
          const truncated = content.length > 200;
          return '<div class="dp-msg"><div class="dp-msg-header"><span class="dp-msg-from" style="color:'+(AGENTS[normRole(m.from)]?.color||'#888')+'">'+(m.from||'sys')+'</span>'+
          '<span class="dp-msg-time">'+fmtTime(m.timestamp)+'</span></div>'+
          '<div class="dp-msg-text">'+renderMd(content.substring(0,200))+(truncated?'...':'')+'</div></div>';
        }).join('')+'</div>';
    }

    if(!activity.length && !msgs.length) {
      html += '<div style="font-size:12px;color:var(--text-muted);padding:20px;text-align:center">No activity yet</div>';
    }
    body.innerHTML = html;
  } else if(tab==='tasks') {
    const tasks = as.tasks || [];
    body.innerHTML = '<div class="dp-section"><div class="dp-section-title">Assigned Tasks</div>'+
      (tasks.length ? tasks.map(t =>
        '<div class="dp-task"><div class="dp-task-title">'+esc(t.title||'Task')+'</div>'+
        '<div class="dp-task-bar"><div class="dp-task-fill" style="width:'+(t.progress||0)+'%;background:'+ag.color+'"></div></div>'+
        '<div class="dp-task-meta"><span>'+(t.status||'pending')+'</span><span>'+(t.progress||0)+'%</span></div></div>'
      ).join('') : '<div style="font-size:12px;color:var(--text-muted)">No tasks</div>')+
      '</div>';
  } else if(tab==='chat') {
    body.innerHTML =
      '<div class="dp-chat">'+
      '<div class="dp-chat-messages" id="chatMessages"></div>'+
      '<div class="dp-chat-input">'+
      '<textarea class="chat-input" id="chatInput" placeholder="Ask '+ag.name+' something..." rows="1"></textarea>'+
      '<button class="chat-send" id="chatSend">Send</button>'+
      '</div></div>';
    loadChatHistory(key);
    document.getElementById('chatSend').addEventListener('click', () => sendChat(key));
    document.getElementById('chatInput').addEventListener('keydown', e => {
      if(e.key==='Enter' && !e.shiftKey) { e.preventDefault(); sendChat(key); }
    });
  } else if(tab==='metrics') {
    const tc = as.toolCount||0;
    body.innerHTML = '<div class="dp-section"><div class="dp-section-title">Session Metrics</div>'+
      '<div style="font-size:12px;color:var(--text-secondary);line-height:2.2">'+
      'State: <strong style="color:'+ag.color+'">'+(as.state||'idle')+'</strong><br>'+
      'Tool calls: <strong>'+tc+'</strong><br>'+
      'Messages sent: <strong>'+(S.messages.filter(m=>normRole(m.from)===key).length)+'</strong><br>'+
      'Messages received: <strong>'+(S.messages.filter(m=>normRole(m.to)===key).length)+'</strong><br>'+
      'Tasks completed: <strong>'+(as.tasks?as.tasks.filter(t=>t.status==='completed').length:0)+'</strong>'+
      '</div></div>'+
      '<div class="dp-section"><div class="dp-section-title">Activity Summary</div>'+
      '<div style="font-size:11px;color:var(--text-muted)">'+(as.activity||[]).length+' events recorded</div></div>';
  }
}

// ═══════════════════════════════════════════
// COMMUNICATION ARCS (SVG)
// ═══════════════════════════════════════════
function sendArc(fromKey, toKey) {
  const f = AGENTS[fromKey], t = AGENTS[toKey];
  if(!f || !t) return;
  const svg = document.getElementById('arcLayer');
  const rect = svg.getBoundingClientRect();
  const sx = f.x/100*rect.width, sy = f.y/100*rect.height;
  const ex = t.x/100*rect.width, ey = t.y/100*rect.height;
  const mx = (sx+ex)/2, my = (sy+ey)/2 - 40;
  const d = 'M'+sx+' '+sy+' Q'+mx+' '+my+' '+ex+' '+ey;

  const path = document.createElementNS('http://www.w3.org/2000/svg','path');
  path.setAttribute('d', d);
  path.setAttribute('stroke', f.color);
  path.setAttribute('class','arc-path active');
  // Set dasharray to path length
  svg.appendChild(path);
  const len = path.getTotalLength();
  path.style.strokeDasharray = len;
  path.style.strokeDashoffset = len;
  path.style.animation = 'arc-flow 0.6s ease forwards';

  // After flow, become trail
  setTimeout(() => {
    path.classList.remove('active');
    path.classList.add('trail');
    path.style.animation = 'none';
    path.style.strokeDashoffset = '0';
    path.style.opacity = '0.15';
    // Fade and remove
    setTimeout(() => { path.classList.add('fading'); }, 100);
    setTimeout(() => { path.remove(); }, 10000);
  }, 700);
}

// ═══════════════════════════════════════════
// LOG STREAM
// ═══════════════════════════════════════════
function addLog(msg) {
  const el = document.getElementById('logStream');
  const row = document.createElement('div');
  row.className = 'log-entry';
  const c = AGENTS[normRole(msg.from)]?.color || '#555';
  row.innerHTML =
    '<span class="log-time">'+fmtTime(msg.timestamp)+'</span>'+
    '<span class="log-who" style="color:'+c+'">'+(msg.from||'SYS')+'</span>'+
    '<span class="log-text">'+esc((msg.content||'').substring(0,200))+'</span>';
  el.appendChild(row);
  el.scrollTop = el.scrollHeight;
  while(el.children.length > 500) el.removeChild(el.firstChild);
}

// ═══════════════════════════════════════════
// COMMAND PALETTE
// ═══════════════════════════════════════════
function openCmd() {
  document.getElementById('cmdOverlay').classList.remove('cmd-hidden');
  const inp = document.getElementById('cmdInput');
  inp.value = ''; inp.focus();
  renderCmdResults('');
}
function closeCmd() { document.getElementById('cmdOverlay').classList.add('cmd-hidden'); }

function renderCmdResults(q) {
  const el = document.getElementById('cmdResults');
  const items = [];
  // Agents
  Object.values(AGENTS).forEach(ag => {
    if(!q || ag.name.toLowerCase().includes(q) || ag.full.toLowerCase().includes(q))
      items.push({icon: ag.name.charAt(0), label:'Focus: '+ag.full, hint:ag.key, action:()=>{ closeCmd(); selectAgent(ag.key); }});
  });
  // Actions
  if(!q || 'submit task'.includes(q))
    items.push({icon:'&#9654;', label:'Submit new task', hint:'Enter', action:()=>{ closeCmd(); document.getElementById('taskInput').focus(); }});
  if(!q || 'toggle log'.includes(q))
    items.push({icon:'&#9776;', label:'Toggle log panel', hint:'Ctrl+J', action:()=>{ closeCmd(); toggleLog(); }});

  el.innerHTML = items.slice(0,10).map((it,i) =>
    '<div class="cmd-item'+(i===0?' active':'')+'" data-idx="'+i+'">'+
    '<span class="cmd-item-icon">'+it.icon+'</span>'+
    '<span class="cmd-item-label">'+it.label+'</span>'+
    '<span class="cmd-item-hint">'+it.hint+'</span></div>'
  ).join('');
  el.querySelectorAll('.cmd-item').forEach((item,i) => {
    item.addEventListener('click', () => items[i]?.action());
  });
}

// ═══════════════════════════════════════════
// CHECKPOINT SYSTEM
// ═══════════════════════════════════════════
const CP_TITLES = {
  template_approval: 'Template Selection Complete',
  plan_approval: 'Development Plan Ready for Review',
  phase_gate: 'Phase Complete — Proceed?',
  final_acceptance: 'All Phases Complete — Final Review',
  research_review: 'Research Phase Complete',
  spec_review: 'Specifications Ready',
  pre_qa_review: 'Pre-QA Review',
};

function handleCheckpointEvt(evt) {
  switch(evt.event_type) {
    case 'checkpoint_reached':
      showCheckpointBanner(evt);
      showToast('Approval Required', CP_TITLES[evt.checkpoint_type]||'Checkpoint reached', 'warning');
      break;
    case 'checkpoint_resolved':
      hideCheckpointBanner();
      showToast('Checkpoint Approved', evt.data?.artifact_summary||'Resolved', 'success');
      addLog({from:'system', content:'Checkpoint resolved: '+(evt.data?.action||''), timestamp:new Date().toISOString()});
      break;
    case 'checkpoint_auto_delegated':
      hideCheckpointBanner();
      showToast('CEO Auto-Delegated', (evt.data?.action||'')+' — '+(evt.data?.feedback||'').substring(0,80), 'info');
      addLog({from:'ceo', content:'Auto-delegated: '+(evt.data?.action||'')+' — '+(evt.data?.feedback||''), timestamp:new Date().toISOString()});
      break;
  }
}

async function showCheckpointBanner(evt) {
  S.checkpoint = evt;
  document.getElementById('cpTitle').textContent = CP_TITLES[evt.checkpoint_type] || evt.checkpoint_type;
  document.getElementById('cpSummary').textContent = evt.data?.artifact_summary || '';
  document.getElementById('cpContent').style.display = 'none';
  document.getElementById('cpContent').innerHTML = '';
  document.getElementById('cpBanner').style.display = '';

  // Fetch full artifact
  try {
    const res = await fetch('/api/checkpoints/'+evt.checkpoint_id+'?project='+S.projectId);
    const d = await res.json();
    if(d.artifact_content) {
      document.getElementById('cpContent').innerHTML = renderMd(d.artifact_content.substring(0,500));
      document.getElementById('cpContent').dataset.raw = d.artifact_content;
    }
  } catch(e){}

  // Start countdown if auto-delegate configured
  if(S.settings?.auto_delegate_minutes > 0) {
    startCpCountdown(S.settings.auto_delegate_minutes * 60);
  } else {
    document.getElementById('cpCountdown').style.display = 'none';
  }
}

function hideCheckpointBanner() {
  document.getElementById('cpBanner').style.display = 'none';
  S.checkpoint = null;
  if(S.cpInterval) { clearInterval(S.cpInterval); S.cpInterval=null; }
}

async function resolveCheckpoint(action, feedback, overrideData) {
  if(!S.checkpoint) return;
  if(S.cpInterval) { clearInterval(S.cpInterval); S.cpInterval=null; }
  try {
    const res = await fetch('/api/checkpoints/'+S.checkpoint.checkpoint_id+'/decide?project='+S.projectId, {
      method:'POST', headers:{'Content-Type':'application/json'},
      body:JSON.stringify({action, feedback:feedback||'', override_data:overrideData||''})
    });
    const d = await res.json();
    if(d.success) hideCheckpointBanner();
    else addLog({from:'system', content:'Checkpoint resolve failed: '+(d.error||''), timestamp:new Date().toISOString()});
  } catch(e) { console.error('[CP]', e); }
}

function startCpCountdown(seconds) {
  const end = Date.now() + seconds*1000;
  const el = document.getElementById('cpCountdown');
  el.style.display = '';
  S.cpInterval = setInterval(() => {
    const rem = Math.max(0, Math.ceil((end - Date.now())/1000));
    const m = Math.floor(rem/60), s = rem%60;
    el.textContent = 'Auto-delegating to CEO in '+m+':'+String(s).padStart(2,'0');
    if(rem<=0) clearInterval(S.cpInterval);
  }, 1000);
}

// ── SETTINGS MODAL ──
async function loadSettings() {
  try {
    const d = await (await fetch('/api/settings/workflow?project='+S.projectId)).json();
    S.settings = d;
  } catch(e) { S.settings = {}; }
}

function openSettings() {
  const s = S.settings || {};
  const body = document.getElementById('settingsBody');
  const toggles = [
    {key:'require_template_approval', label:'Template Selection', def:true, mandatory:true},
    {key:'require_plan_approval', label:'Development Plan', def:true, mandatory:true},
    {key:'require_phase_gate', label:'Phase Gates (between dev phases)', def:true, mandatory:true},
    {key:'require_final_acceptance', label:'Final Acceptance', def:true, mandatory:true},
    {key:'require_research_review', label:'Research Review', def:false},
    {key:'require_spec_review', label:'Spec Review', def:false},
    {key:'require_pre_qa_review', label:'Pre-QA Review', def:false},
  ];
  body.innerHTML =
    '<div class="setting-section"><div class="setting-section-title">Approval Checkpoints</div>'+
    toggles.map(t => {
      const on = s[t.key] !== undefined ? s[t.key] : t.def;
      return '<div class="setting-row"><span>'+t.label+(t.mandatory?' *':'')+'</span>'+
        '<div class="toggle-switch'+(on?' on':'')+'" data-key="'+t.key+'"></div></div>';
    }).join('')+'</div>'+
    '<div class="setting-section"><div class="setting-section-title">Auto-Delegation</div>'+
    '<div class="setting-row"><span>Delegate to CEO after</span>'+
    '<select class="setting-select" id="autoDelegateSelect">'+
    [0,2,5,10,15,30].map(v => '<option value="'+v+'"'+(v===(s.auto_delegate_minutes||0)?' selected':'')+'>'+
      (v===0?'Off':v+' min')+'</option>').join('')+
    '</select></div>'+
    '<div style="font-size:10px;color:var(--text-muted);margin-top:4px;line-height:1.5">'+
    'When enabled, the CEO agent will review checkpoints<br>on your behalf if you don\'t respond in time.</div></div>';

  body.querySelectorAll('.toggle-switch').forEach(el => {
    el.addEventListener('click', () => el.classList.toggle('on'));
  });
  document.getElementById('settingsOverlay').style.display = '';
}

function closeSettings() { document.getElementById('settingsOverlay').style.display = 'none'; }

async function saveSettings() {
  const settings = {project_id: S.projectId};
  document.querySelectorAll('#settingsBody .toggle-switch').forEach(el => {
    settings[el.dataset.key] = el.classList.contains('on');
  });
  settings.auto_delegate_minutes = parseInt(document.getElementById('autoDelegateSelect')?.value||'0');
  try {
    await fetch('/api/settings/workflow?project='+S.projectId, {
      method:'PUT', headers:{'Content-Type':'application/json'},
      body:JSON.stringify(settings)
    });
    S.settings = settings;
  } catch(e) {}
  closeSettings();
}

// ── FEEDBACK MODAL ──
let feedbackMode = 'reject'; // 'reject' or 'override'
function openFeedback(mode) {
  feedbackMode = mode;
  document.getElementById('feedbackTitle').textContent = mode==='override' ? 'Override Content' : 'Reject with Feedback';
  document.getElementById('feedbackText').placeholder = mode==='override' ? 'Enter the corrected content...' : 'What needs to change?';
  document.getElementById('feedbackText').value = '';
  document.getElementById('feedbackOverlay').style.display = '';
  document.getElementById('feedbackText').focus();
}
function closeFeedback() { document.getElementById('feedbackOverlay').style.display = 'none'; }
function submitFeedback() {
  const text = document.getElementById('feedbackText').value.trim();
  if(!text) return;
  if(feedbackMode==='override') resolveCheckpoint('overridden', '', text);
  else resolveCheckpoint('rejected', text, '');
  closeFeedback();
}

// ═══════════════════════════════════════════
// KANBAN BOARD
// ═══════════════════════════════════════════
let currentView = 'topology';
let kanbanTasks = [];

function switchView(view) {
  currentView = view;
  document.querySelectorAll('.view-btn').forEach(b => b.classList.toggle('active', b.dataset.view===view));

  const roster = document.getElementById('roster');
  const canvas = document.getElementById('canvasArea');
  const kb = document.getElementById('kanbanView');

  if(view==='kanban') {
    roster.style.display = 'none';
    canvas.style.display = 'none';
    kb.classList.add('active');
    document.getElementById('app').style.gridTemplateColumns = '0 1fr 0';
    if(document.getElementById('app').classList.contains('panel-open'))
      document.getElementById('app').style.gridTemplateColumns = '0 1fr var(--panel-w)';
    fetchKanbanTasks();
  } else {
    roster.style.display = '';
    canvas.style.display = '';
    kb.classList.remove('active');
    document.getElementById('app').style.gridTemplateColumns = '';
    if(document.getElementById('app').classList.contains('panel-open'))
      document.getElementById('app').style.gridTemplateColumns = 'var(--roster-w) 1fr var(--panel-w)';
  }
}

async function fetchKanbanTasks() {
  try {
    const params = new URLSearchParams({project: S.projectId});
    const fa = document.getElementById('kbFilterAgent').value;
    const fp = document.getElementById('kbFilterPhase').value;
    const fpr = document.getElementById('kbFilterPriority').value;
    if(fa) params.set('agent', fa);
    if(fp) params.set('phase', fp);
    if(fpr) params.set('priority', fpr);

    const d = await (await fetch('/api/kanban/tasks?'+params)).json();
    kanbanTasks = d.tasks || [];
    renderKanban();
  } catch(e) { console.error('[KB]', e); }
}

function renderKanban() {
  const search = (document.getElementById('kbSearch').value||'').toLowerCase();
  const cols = {backlog:[], in_progress:[], in_review:[], done:[]};

  kanbanTasks.forEach(t => {
    if(search && !t.title.toLowerCase().includes(search) && !(t.assigned_agent||'').toLowerCase().includes(search)) return;
    const s = t.status === 'blocked' ? 'backlog' : t.status;
    if(cols[s]) cols[s].push(t);
  });

  // Sort each column by position
  Object.values(cols).forEach(arr => arr.sort((a,b) => a.position - b.position));

  // Update counts
  document.getElementById('kbCountBacklog').textContent = cols.backlog.length;
  document.getElementById('kbCountInProgress').textContent = cols.in_progress.length;
  document.getElementById('kbCountInReview').textContent = cols.in_review.length;
  document.getElementById('kbCountDone').textContent = cols.done.length;

  // Render cards
  renderKbColumn('kbColBacklog', cols.backlog, true);
  renderKbColumn('kbColInProgress', cols.in_progress, false);
  renderKbColumn('kbColInReview', cols.in_review, false);
  renderKbColumn('kbColDone', cols.done, false);
}

function renderKbColumn(elId, tasks, showAddBtn) {
  const el = document.getElementById(elId);
  el.innerHTML = '';
  tasks.forEach(t => {
    const ag = AGENTS[normRole(t.assigned_agent)];
    const color = ag?.color || '#888';
    const card = document.createElement('div');
    card.className = 'kanban-card' + (t.status==='blocked'?' blocked':'');
    card.style.borderLeftColor = color;
    card.draggable = true;
    card.dataset.taskId = t.id;

    const phaseLabel = t.phase_name || ('P'+t.phase_index);
    const priCls = t.priority || 'medium';

    card.innerHTML =
      '<div class="kc-header">'+
        '<span class="kc-id">'+esc(t.id.slice(-8))+'</span>'+
        '<div class="kc-agent" style="color:'+color+';border-color:'+color+'">'+(ag?.name?.charAt(0)||'?')+'</div>'+
      '</div>'+
      '<div class="kc-title">'+esc(t.title)+'</div>'+
      (t.progress > 0 ? '<div class="kc-bar"><div class="kc-bar-fill" style="width:'+t.progress+'%;background:'+color+'"></div></div>' : '')+
      '<div class="kc-footer">'+
        '<span class="kc-priority '+priCls+'">'+priCls+'</span>'+
        '<span>'+phaseLabel+'</span>'+
      '</div>';

    // Drag events
    card.addEventListener('dragstart', e => {
      e.dataTransfer.setData('text/plain', t.id);
      card.classList.add('dragging');
    });
    card.addEventListener('dragend', () => card.classList.remove('dragging'));

    // Click to open detail
    card.addEventListener('click', () => openTaskDetail(t));

    el.appendChild(card);
  });

  if(showAddBtn) {
    const btn = document.createElement('button');
    btn.className = 'kanban-add-btn';
    btn.textContent = '+ Add Task';
    btn.addEventListener('click', () => showNewTaskForm(elId));
    el.appendChild(btn);
  }
}

function showNewTaskForm(colId) {
  const el = document.getElementById(colId);
  // Remove existing form if any
  const existing = el.querySelector('.kanban-new-card');
  if(existing) { existing.remove(); return; }

  const form = document.createElement('div');
  form.className = 'kanban-new-card';
  form.innerHTML =
    '<input type="text" placeholder="Task title" class="nc-title">'+
    '<textarea placeholder="Description (optional)" class="nc-desc"></textarea>'+
    '<select class="nc-agent"><option value="">Assign agent</option>'+
    Object.values(AGENTS).map(a => '<option value="'+a.key+'">'+a.name+'</option>').join('')+'</select>'+
    '<select class="nc-priority"><option value="medium">Medium</option><option value="critical">Critical</option><option value="high">High</option><option value="low">Low</option></select>'+
    '<div class="kanban-card-actions"><button class="cancel">Cancel</button><button class="save">Create</button></div>';

  form.querySelector('.cancel').addEventListener('click', () => form.remove());
  form.querySelector('.save').addEventListener('click', async () => {
    const title = form.querySelector('.nc-title').value.trim();
    if(!title) return;
    await fetch('/api/kanban/tasks?project='+S.projectId, {
      method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({
        title,
        description: form.querySelector('.nc-desc').value,
        assigned_agent: form.querySelector('.nc-agent').value,
        priority: form.querySelector('.nc-priority').value,
      })
    });
    form.remove();
    fetchKanbanTasks();
  });

  // Insert before add button
  const addBtn = el.querySelector('.kanban-add-btn');
  if(addBtn) el.insertBefore(form, addBtn);
  else el.appendChild(form);
  form.querySelector('.nc-title').focus();
}

function openTaskDetail(t) {
  // Reuse detail panel for task view
  const ag = AGENTS[normRole(t.assigned_agent)];
  const color = ag?.color || '#888';
  S.selected = null; // Not an agent selection
  document.getElementById('app').classList.add('panel-open');
  if(currentView==='kanban') document.getElementById('app').style.gridTemplateColumns = '0 1fr var(--panel-w)';

  document.getElementById('dpAvatar').textContent = ag?.name?.charAt(0) || '?';
  document.getElementById('dpAvatar').style.borderColor = color;
  document.getElementById('dpAvatar').style.color = color;
  document.getElementById('dpName').textContent = t.title;
  document.getElementById('dpName').style.color = color;
  document.getElementById('dpRole').textContent = (ag?.full||'Unassigned') + ' | ' + (t.status||'backlog');

  const body = document.getElementById('dpBody');
  const criteria = (t.criteria||[]).map((c,i) =>
    '<label style="display:flex;align-items:center;gap:6px;font-size:11px;padding:3px 0;cursor:pointer">'+
    '<input type="checkbox" '+(c.checked?'checked':'')+' data-idx="'+i+'">'+esc(c.text)+'</label>'
  ).join('') || '<span style="font-size:11px;color:var(--text-muted)">No criteria</span>';

  body.innerHTML =
    '<div class="dp-section"><div class="dp-section-title">Details</div>'+
    '<div style="font-size:12px;color:var(--text-secondary);line-height:1.8">'+
    'ID: <span style="font-family:var(--font-mono)">'+esc(t.id)+'</span><br>'+
    'Status: <strong>'+esc(t.status)+'</strong><br>'+
    'Priority: <span class="kc-priority '+t.priority+'">'+esc(t.priority)+'</span><br>'+
    'Phase: '+(t.phase_name||'P'+t.phase_index)+'<br>'+
    'Progress: '+t.progress+'%<br>'+
    'Source: '+esc(t.source||'manual')+
    '</div></div>'+
    (t.description ? '<div class="dp-section"><div class="dp-section-title">Description</div><div style="font-size:12px;color:var(--text-secondary);line-height:1.5">'+renderMd(t.description)+'</div></div>' : '')+
    '<div class="dp-section"><div class="dp-section-title">Acceptance Criteria</div>'+criteria+'</div>'+
    '<div class="dp-section"><div class="dp-section-title">Review</div>'+
    '<div style="display:flex;gap:6px;margin-top:4px">'+
    '<button class="cp-btn approve" style="font-size:10px;padding:5px 12px" onclick="reviewKbTask(\''+t.id+'\',\'approved\')">Approve</button>'+
    '<button class="cp-btn reject" style="font-size:10px;padding:5px 12px" onclick="reviewKbTask(\''+t.id+'\',\'rejected\')">Reject</button>'+
    '</div></div>';
}

async function reviewKbTask(taskId, action) {
  const feedback = action==='rejected' ? prompt('Feedback:') : '';
  if(action==='rejected' && !feedback) return;
  await fetch('/api/kanban/tasks/'+taskId+'/review?project='+S.projectId, {
    method:'POST', headers:{'Content-Type':'application/json'},
    body: JSON.stringify({action, feedback: feedback||''})
  });
  fetchKanbanTasks();
  closePanel();
}

async function syncKanban() {
  const res = await fetch('/api/kanban/sync?project='+S.projectId, {method:'POST'});
  const d = await res.json();
  if(d.success) {
    addLog({from:'system', content:'Synced '+d.created+' tasks from development plan', timestamp:new Date().toISOString()});
    fetchKanbanTasks();
  }
}

// Drag & Drop
function initKanbanDragDrop() {
  document.querySelectorAll('.kanban-col-body').forEach(col => {
    col.addEventListener('dragover', e => { e.preventDefault(); col.classList.add('drag-over'); });
    col.addEventListener('dragleave', () => col.classList.remove('drag-over'));
    col.addEventListener('drop', async e => {
      e.preventDefault(); col.classList.remove('drag-over');
      const taskId = e.dataTransfer.getData('text/plain');
      const newStatus = col.closest('.kanban-col').dataset.status;
      if(!taskId || !newStatus) return;
      await fetch('/api/kanban/tasks/'+taskId+'?project='+S.projectId, {
        method:'PATCH', headers:{'Content-Type':'application/json'},
        body: JSON.stringify({status: newStatus})
      });
      fetchKanbanTasks();
    });
  });
}

// Populate agent filter
function populateKbFilters() {
  const sel = document.getElementById('kbFilterAgent');
  sel.innerHTML = '<option value="">All Agents</option>';
  Object.values(AGENTS).forEach(a => {
    sel.innerHTML += '<option value="'+a.key+'">'+a.name+'</option>';
  });
}

// ═══════════════════════════════════════════
// WEBSOCKET
// ═══════════════════════════════════════════
function connectWS() {
  const proto = location.protocol==='https:'?'wss:':'ws:';
  S.ws = new WebSocket(proto+'//'+location.host+'/ws');
  S.ws.onopen = () => { document.getElementById('wsDot').className='kpi-dot live'; };
  S.ws.onclose = () => { document.getElementById('wsDot').className='kpi-dot off'; setTimeout(connectWS,3000); };
  S.ws.onerror = () => { S.ws.close(); };
  S.ws.onmessage = (e) => { try { handleWS(JSON.parse(e.data)); } catch(x){ console.error('[WS]',x); } };
}

function handleWS(data) {
  if(data.type==='message' && data.message) {
    const m = data.message;
    S.messages.push(m); if(S.messages.length>500) S.messages.shift();
    S.msgCount++;
    addLog(m);
    handleMsg(m);
  }
  if(data.type==='agent_task_event' && data.event) handleTaskEvt(data.event);
  if(data.type==='agent_event' && data.event) handleAgentSessionEvt(data.event);
  if(data.type==='message' && data.message && data.message.type==='lifecycle') handleLifecycleEvt(data.message);
  if(data.type==='checkpoint' && data.event) {
    if(data.event.event_type==='chat_message' || data.event.event_type==='chat_response') handleChatWS(data.event);
    else handleCheckpointEvt(data.event);
  }
}

function handleMsg(msg) {
  const fk = normRole(msg.from), tk = normRole(msg.to);
  // Phase detection
  if(msg.type==='system'||msg.type==='phase') {
    const c = (msg.content||'').toLowerCase();
    PHASES.forEach(p => { if(c.includes(p) && c.includes('phase')) { S.phase=p; renderPipeline(); } });
  }
  if(fk && AGENTS[fk]) {
    if(!S.agents[fk]) S.agents[fk] = {state:'idle',tasks:[],taskTitle:''};
    S.agents[fk].state = 'working';
    if(msg.type==='delegate' && tk && AGENTS[tk]) sendArc(fk, tk);
  }
  if(S.selected===fk||S.selected===tk) renderPanelTab(S.activeTab);
  renderNodes(); renderRoster(); updateKPIs();
}

function handleTaskEvt(evt) {
  const k = normRole(evt.agent_role); if(!k) return;
  if(!S.agents[k]) S.agents[k] = {state:'idle',tasks:[],taskTitle:''};
  const a = S.agents[k]; a.tasks = a.tasks||[];
  switch(evt.event_type) {
    case 'task_started':
      a.state='working'; a.taskTitle=evt.data?.title||'';
      a.tasks.push({id:evt.task_id,title:evt.data?.title||'',status:'in_progress',progress:0});
      showToast(AGENTS[k]?.name+' Started', evt.data?.title||'New task', 'info');
      break;
    case 'task_progress':
      a.state='working'; a.progress=evt.data?.progress||0;
      const pt = a.tasks.find(t=>t.id===evt.task_id); if(pt) pt.progress=evt.data?.progress||0;
      break;
    case 'task_completed':
      a.state='completed';
      const ct = a.tasks.find(t=>t.id===evt.task_id); if(ct){ct.status='completed';ct.progress=100;}
      showToast(AGENTS[k]?.name+' Done', (ct?.title||'Task')+' completed', 'success');
      setTimeout(()=>{ if(S.agents[k]?.state==='completed'){S.agents[k].state='idle';S.agents[k].taskTitle='';renderNodes();renderRoster();} },5000);
      break;
    case 'task_failed':
      a.state='error';
      showToast(AGENTS[k]?.name+' Error', 'Task failed', 'error');
      const ft = a.tasks.find(t=>t.id===evt.task_id); if(ft) ft.status='failed';
      break;
  }
  renderNodes(); renderRoster(); updateKPIs();
  if(S.selected===k) renderPanelTab(S.activeTab);
}

// ═══════════════════════════════════════════
// AGENT SESSION EVENTS (Claude Code live activity)
// ═══════════════════════════════════════════
const TOOL_ICONS = {Read:'📖',Write:'📝',Edit:'✏️',Bash:'💻',Grep:'🔍',Glob:'📂',Agent:'🤖',WebSearch:'🌐',WebFetch:'🌐',default:'🔧'};
function fmtToolInput(name, inputStr) {
  try {
    const inp = JSON.parse(inputStr);
    switch(name) {
      case 'Read': return inp.file_path||'?';
      case 'Write': return (inp.file_path||'?')+' ('+((inp.content||'').length)+' chars)';
      case 'Edit': return (inp.file_path||'?')+' (edit)';
      case 'Bash': return '$ '+(inp.command||'?');
      case 'Grep': return 'grep "'+((inp.pattern||''))+'" '+(inp.path||'.');
      case 'Glob': return 'glob "'+(inp.pattern||'')+'"';
      default: return inputStr.substring(0,80);
    }
  } catch(e) { return (inputStr||'').substring(0,80); }
}

function handleAgentSessionEvt(ev) {
  const k = normRole(ev.agent_role); if(!k) return;
  if(!S.agents[k]) S.agents[k] = {state:'idle',tasks:[],taskTitle:'',activity:[]};
  const a = S.agents[k];
  if(!a.activity) a.activity=[];

  switch(ev.type) {
    case 'text_delta':
      a.state='working';
      a.taskTitle=ev.content?(ev.content.substring(0,60)):'Working...';
      break;
    case 'thinking_delta':
      a.state='working';
      a.taskTitle='💭 Thinking...';
      a.activity.push({t:Date.now(),icon:'💭',text:'Thinking...'});
      break;
    case 'tool_use':
      a.state='working';
      const icon = TOOL_ICONS[ev.tool_name]||TOOL_ICONS.default;
      const desc = fmtToolInput(ev.tool_name, ev.input||'');
      a.taskTitle=icon+' '+ev.tool_name+': '+desc.substring(0,50);
      a.activity.push({t:Date.now(),icon:icon,text:ev.tool_name+': '+desc});
      if(!a.toolCount) a.toolCount=0;
      a.toolCount++;
      break;
    case 'tool_result':
      if(ev.is_error) a.activity.push({t:Date.now(),icon:'❌',text:'Error: '+(ev.output||'').substring(0,100)});
      break;
    case 'turn_complete':
      a.state='completed';
      a.taskTitle='✅ Turn complete';
      a.activity.push({t:Date.now(),icon:'✅',text:'Turn complete — tokens: '+(ev.input_tokens||0)+' in / '+(ev.output_tokens||0)+' out, $'+(ev.cost_usd||0).toFixed(4)});
      setTimeout(()=>{ if(S.agents[k]?.state==='completed'){S.agents[k].state='idle';S.agents[k].taskTitle='';renderNodes();renderRoster();} },5000);
      break;
    case 'session_meta':
      a.state='working';
      a.taskTitle='Session started ('+(ev.model||'unknown')+')';
      a.activity.push({t:Date.now(),icon:'⚡',text:'Session started: '+(ev.model||'unknown')});
      break;
    case 'error':
      a.state='error';
      a.taskTitle='Error';
      a.activity.push({t:Date.now(),icon:'⚠️',text:'Error: '+(ev.content||'unknown')});
      break;
  }
  // Keep activity log bounded
  if(a.activity.length>200) a.activity=a.activity.slice(-100);

  renderNodes(); renderRoster(); updateKPIs();
  if(S.selected===k) renderPanelTab(S.activeTab);
}

// Handle lifecycle events (phase/task completion)
function handleLifecycleEvt(msg) {
  const extra = msg.metadata?.extra || {};
  const evt = extra.event_type;
  if(evt==='task_completed') {
    S.running=false;
    document.getElementById('taskSubmit').disabled=false;
    const fc = extra.files_count||0;
    const turns = extra.turns||0;
    addLog({from:'system',content:'Task completed: '+turns+' turns, '+fc+' files created',timestamp:new Date().toISOString()});
    // Reset all agents
    Object.keys(S.agents).forEach(k => {
      S.agents[k].state='idle'; S.agents[k].taskTitle='';
    });
    renderNodes(); renderRoster(); updateKPIs();
    // Auto-switch to project tab and load details
    switchLogTab('project');
    fetchProjectDetail();
  }
  if(evt==='phase_started') {
    addLog({from:'system',content:'Phase: '+extra.phase_name,timestamp:new Date().toISOString()});
  }
}

// ═══════════════════════════════════════════
// PROJECT PANEL (integrated in log area)
// ═══════════════════════════════════════════
let activeLogTab = 'log';
const AGENT_COLORS_P = {ceo:'#60B4D4',pm:'#8B7EF8',ux:'#F47F7F',ui:'#4FD1C5',security:'#F6A623',architect:'#A87CF5',senior_dev:'#34D478',junior_dev:'#4BA3E3'};
const FILE_ICONS_P = {html:'&#x1F310;',css:'&#x1F3A8;',js:'&#x26A1;',ts:'&#x1F535;',json:'&#x1F4CB;',md:'&#x1F4DD;',go:'&#x1F439;',default:'&#x1F4C4;'};

function switchLogTab(tab) {
  activeLogTab = tab;
  document.querySelectorAll('.log-tab').forEach(t => t.classList.toggle('active', t.dataset.logtab===tab));
  document.getElementById('logStream').style.display = tab==='log' ? '' : 'none';
  document.getElementById('projectPanel').style.display = tab==='project' ? '' : 'none';
  if(tab==='project') fetchProjectDetail();
  // Expand log area if collapsed
  document.getElementById('logArea').classList.remove('collapsed');
}

async function fetchProjectDetail() {
  const panel = document.getElementById('projectPanel');
  panel.innerHTML = '<div style="padding:16px;color:var(--text-muted);font-size:12px">Loading project details...</div>';
  try {
    const res = await fetch('/api/projects/'+encodeURIComponent(S.projectId));
    const data = await res.json();
    renderProjectPanel(data);
  } catch(e) {
    panel.innerHTML = '<div style="padding:16px;color:var(--text-muted);font-size:12px">Failed to load: '+e.message+'</div>';
  }
}

function renderProjectPanel(data) {
  const panel = document.getElementById('projectPanel');
  const files = data.files || [];
  const sessions = data.sessions || {};
  const agents = sessions.agent_metrics || {};
  const plan = data.dev_plan;

  let html = '<div class="proj-grid">';

  // Files card
  html += '<div class="proj-card"><div class="proj-card-title">Files ('+files.length+')</div>';
  if(files.length) {
    html += files.slice(0,20).map(f => {
      const ext = (f.path||'').split('.').pop().toLowerCase();
      const icon = FILE_ICONS_P[ext]||FILE_ICONS_P.default;
      const sz = f.size<1024 ? f.size+'B' : (f.size/1024).toFixed(1)+'KB';
      return '<div class="proj-file" onclick="viewProjectFile(\''+esc(f.path)+'\')">'+icon+' <span style="flex:1;overflow:hidden;text-overflow:ellipsis">'+esc(f.path)+'</span> <span style="color:var(--text-muted)">'+sz+'</span></div>';
    }).join('');
    if(files.length>20) html += '<div style="font-size:10px;color:var(--text-muted);padding-top:4px">+'+(files.length-20)+' more files</div>';
  } else {
    html += '<div style="color:var(--text-muted)">No files yet</div>';
  }
  html += '</div>';

  // Agents card
  html += '<div class="proj-card"><div class="proj-card-title">Agent Metrics</div>';
  const agentEntries = Object.entries(agents).sort((a,b) => a[0].localeCompare(b[0]));
  if(agentEntries.length) {
    html += agentEntries.map(([role,m]) => {
      const color = AGENT_COLORS_P[role]||'#888';
      const toolStr = Object.entries(m.tools||{}).slice(0,4).map(([k,v])=>k+':'+v).join(' ') || '-';
      return '<div class="proj-agent">'+
        '<div class="proj-agent-dot" style="background:'+color+'"></div>'+
        '<div class="proj-agent-name">'+(m.name||role)+'</div>'+
        '<div class="proj-agent-val">'+(m.turns||0)+'t '+(m.tool_calls||0)+'tc $'+(m.cost_usd||0).toFixed(2)+'</div>'+
        '<div class="proj-agent-val" style="flex:1;overflow:hidden;text-overflow:ellipsis">'+esc(toolStr)+'</div></div>';
    }).join('');
    html += '<div class="proj-agent" style="font-weight:600;border-top:1px solid var(--border-subtle);margin-top:4px;padding-top:6px">'+
      '<div class="proj-agent-dot" style="background:transparent"></div>'+
      '<div class="proj-agent-name">Total</div>'+
      '<div class="proj-agent-val" style="color:var(--st-complete)">$'+(sessions.total_cost||'0')+'</div>'+
      '<div class="proj-agent-val">'+(sessions.total_tokens||0)+' tokens</div></div>';
  } else {
    html += '<div style="color:var(--text-muted)">No session data yet</div>';
  }
  html += '</div>';

  // Dev plan card (if exists)
  if(plan && plan.phases) {
    html += '<div class="proj-card"><div class="proj-card-title">Development Plan</div>';
    plan.phases.forEach(phase => {
      const stColors = {pending:'var(--text-muted)',in_progress:'var(--accent)',completed:'var(--st-complete)',needs_revision:'#F6A623'};
      html += '<div style="margin-bottom:6px"><span style="display:inline-block;width:6px;height:6px;border-radius:50%;background:'+(stColors[phase.status]||'var(--text-muted)')+';margin-right:4px"></span>'+
        '<span style="font-weight:500">'+esc(phase.name)+'</span> <span style="color:var(--text-muted)">('+phase.status+')</span></div>';
      if(phase.subtasks) phase.subtasks.forEach(st => {
        const stC = st.status==='completed'?'var(--st-complete)':st.status==='in_progress'?'var(--accent)':'var(--text-muted)';
        html += '<div style="display:flex;gap:6px;padding-left:14px;font-size:10px;color:var(--text-secondary)">'+
          '<span style="display:inline-block;width:6px;height:6px;border-radius:50%;background:'+stC+';margin-top:4px;flex-shrink:0"></span>'+
          esc(st.title||'Subtask')+'</div>';
      });
    });
    html += '</div>';
  }

  // Task history card
  const tasks = data.tasks||[];
  if(tasks.length) {
    html += '<div class="proj-card"><div class="proj-card-title">Task History ('+tasks.length+')</div>';
    tasks.slice(-5).forEach(t => {
      const status = t.Status||t.status||'?';
      const task = t.Task||t.task||'';
      const isOk = status==='completed'||status==='success';
      html += '<div style="padding:4px 0;border-bottom:1px solid var(--border-subtle);font-size:10px">'+
        '<span style="display:inline-block;padding:1px 5px;border-radius:3px;font-size:9px;background:'+(isOk?'rgba(34,197,94,0.15)':'rgba(239,68,68,0.15)')+';color:'+(isOk?'var(--st-complete)':'var(--st-error)')+'">'+status+'</span> '+
        esc(task.substring(0,80))+'</div>';
    });
    html += '</div>';
  }

  html += '</div>';
  panel.innerHTML = html;
}

// ═══════════════════════════════════════════
// AGENT MODES (oneshot / session switching)
// ═══════════════════════════════════════════
async function fetchAgentModes() {
  try {
    const d = await (await fetch('/api/agents/modes')).json();
    S.agentModes = d.modes || {};
    renderRoster();
  } catch(e) { console.error('fetchAgentModes:', e); }
}
async function toggleAgentMode(role) {
  const current = (S.agentModes||{})[role] || 'session';
  const next = current==='oneshot' ? 'session' : 'oneshot';
  try {
    await fetch('/api/agents/modes', {method:'POST', headers:{'Content-Type':'application/json'},
      body:JSON.stringify({role,mode:next})});
    if(!S.agentModes) S.agentModes={};
    S.agentModes[role] = next;
    renderRoster();
    addLog({from:'system',content:role+' switched to '+next+' mode',timestamp:new Date().toISOString()});
  } catch(e) { console.error('toggleAgentMode:', e); }
}

async function viewProjectFile(relPath) {
  try {
    const res = await fetch('/api/file-content?path='+encodeURIComponent(S.projectId+'/'+relPath));
    if(res.ok) {
      const d = await res.json();
      openViewerRaw(relPath, d.content||'(empty)');
    }
  } catch(e) { console.error(e); }
}

// ═══════════════════════════════════════════
// API
// ═══════════════════════════════════════════
function openInject() {
  document.getElementById('injectOverlay').classList.remove('cmd-hidden');
  document.getElementById('injectTaskInput').focus();
}
function closeInject() {
  document.getElementById('injectOverlay').classList.add('cmd-hidden');
  document.getElementById('injectStatus').textContent='';
}
async function submitInject() {
  const task = document.getElementById('injectTaskInput').value.trim();
  if(!task) return;
  const agentRole = document.getElementById('injectAgentSel').value;
  const priority = document.getElementById('injectPrioritySel').value;
  const statusEl = document.getElementById('injectStatus');
  statusEl.textContent='Injecting...';
  try {
    const res = await fetch('/api/inject',{method:'POST',headers:{'Content-Type':'application/json'},
      body:JSON.stringify({task,agent_role:agentRole,project_id:S.projectId,priority})});
    const d = await res.json();
    if(d.success) {
      statusEl.innerHTML='&#x2705; '+d.message;
      document.getElementById('injectTaskInput').value='';
      addLog({from:'human',content:'Injected task for '+agentRole+': '+task,timestamp:new Date().toISOString()});
      setTimeout(closeInject,2000);
    } else {
      statusEl.textContent='Error: '+(d.error||'Failed');
    }
  } catch(e) { statusEl.textContent='Error: '+e.message; }
}

async function submitTask() {
  const inp = document.getElementById('taskInput');
  const task = inp.value.trim(); if(!task) return;
  const btn = document.getElementById('taskSubmit');
  btn.disabled = true; S.running = true; S.startTime = Date.now();
  S.agents = {}; S.messages = []; S.msgCount = 0;
  document.getElementById('logStream').innerHTML = '';
  renderNodes(); renderRoster();
  try {
    const res = await fetch('/api/task',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({task,project_id:S.projectId})});
    const d = await res.json();
    if(d.success) { inp.value=''; addLog({from:'system',content:'Task: '+task,timestamp:new Date().toISOString()}); }
    else { addLog({from:'system',content:'Error: '+(d.error||'Failed'),timestamp:new Date().toISOString()}); btn.disabled=false; S.running=false; }
  } catch(e) { addLog({from:'system',content:'Connection error',timestamp:new Date().toISOString()}); btn.disabled=false; S.running=false; }
}

async function fetchProjects() {
  try {
    const d = await (await fetch('/api/projects')).json();
    const sel = document.getElementById('projectSel'); sel.innerHTML='';
    (d.projects||[]).forEach(p => { const o=document.createElement('option'); o.value=p.id; o.textContent=p.id; if(p.id===S.projectId)o.selected=true; sel.appendChild(o); });
  } catch(e){}
}

async function fetchStatus() {
  try {
    const d = await (await fetch('/api/status')).json();
    const was = S.running;

    // Check if OUR project is the one running
    const ourProjectRunning = (d.running_projects||[]).includes(S.projectId);
    S.running = ourProjectRunning;

    // Enable/disable submit based on our project's state
    document.getElementById('taskSubmit').disabled = ourProjectRunning;

    if(was && !ourProjectRunning) {
      // Task just finished — move working agents to completed, then idle
      showToast('Build Complete', 'All agents have finished', 'success');
      Object.keys(S.agents).forEach(k => {
        if(S.agents[k].state==='working') S.agents[k].state='completed';
      });
      renderNodes(); renderRoster(); updateKPIs();

      // After 5s, reset all to idle
      const completionTime = Date.now();
      S._lastCompletion = completionTime;
      setTimeout(()=>{
        // Only reset if no new task started since
        if(S._lastCompletion === completionTime && !S.running) {
          Object.keys(S.agents).forEach(k => { S.agents[k].state='idle'; S.agents[k].taskTitle=''; });
          renderNodes(); renderRoster();
        }
      }, 5000);
    }
  } catch(e){}
}

async function fetchAgentTasks() {
  try {
    const d = await (await fetch('/api/agent-tasks?project='+S.projectId)).json();
    if(!d.tasks) return;

    // Reset agent states from tasks — derive state from actual task statuses
    const agentHasInProgress = {};

    d.tasks.forEach(task => {
      const k = normRole(task.agent_role); if(!k||!AGENTS[k]) return;
      if(!S.agents[k]) S.agents[k] = {state:'idle',tasks:[],taskTitle:''};
      S.agents[k].tasks = S.agents[k].tasks||[];
      const td = {id:task.id,title:task.title,status:task.status,progress:task.progress};
      const idx = S.agents[k].tasks.findIndex(t=>t.id===task.id);
      if(idx>=0) S.agents[k].tasks[idx]=td; else S.agents[k].tasks.push(td);

      // Track which agents have in-progress work
      if(task.status==='in_progress') { agentHasInProgress[k]=task.title; }
    });

    // Set state based on actual task statuses
    Object.keys(AGENTS).forEach(k => {
      if(!S.agents[k]) S.agents[k] = {state:'idle',tasks:[],taskTitle:''};
      if(agentHasInProgress[k]) {
        S.agents[k].state = 'working';
        S.agents[k].taskTitle = agentHasInProgress[k];
      } else if(!S.running && S.agents[k].state === 'working') {
        // No in-progress tasks and system not running — agent is idle
        S.agents[k].state = 'idle';
        S.agents[k].taskTitle = '';
      }
    });
    renderNodes(); renderRoster(); updateKPIs();
  } catch(e){}
}

// ═══════════════════════════════════════════
// UTILS
// ═══════════════════════════════════════════
function normRole(r) {
  if(!r) return null;
  const m={ceo:'ceo',pm:'pm',ux:'ux',ui:'ui',security:'security',architect:'architect',
    'senior-dev':'senior-dev','senior_dev':'senior-dev',seniordev:'senior-dev',
    'junior-dev':'junior-dev','junior_dev':'junior-dev',juniordev:'junior-dev'};
  return m[r.toLowerCase().replace(/_/g,'-')]||null;
}
function esc(s) { return s?s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'):''; }
function renderMd(s) {
  if(!s) return '';
  var BT = String.fromCharCode(96);
  // Escape HTML first
  s = esc(s);
  // Code blocks (triple backtick)
  var cbRe = new RegExp(BT+BT+BT+'[\\s\\S]*?'+BT+BT+BT, 'g');
  s = s.replace(cbRe, function(m) {
    var inner = m.slice(3).replace(/^[^\n]*\n?/,'');
    inner = inner.slice(0, inner.length - 3).replace(/\n$/, '');
    return '<pre><code>'+inner+'</code></pre>';
  });
  // Inline code
  var icRe = new RegExp(BT+'([^'+BT+']+)'+BT, 'g');
  s = s.replace(icRe, '<code>$1</code>');
  // Headings
  s = s.replace(/^### (.+)$/gm, '<h5>$1</h5>');
  s = s.replace(/^## (.+)$/gm, '<h4>$1</h4>');
  s = s.replace(/^# (.+)$/gm, '<h3>$1</h3>');
  // Bold & italic
  s = s.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
  s = s.replace(/\*(.+?)\*/g, '<em>$1</em>');
  // Tables (lines with pipes)
  s = s.replace(/((?:^\|.+\|$\n?)+)/gm, function(block) {
    var rows = block.trim().split('\n').filter(function(r){ return !/^[\|\s\-:]+$/.test(r); });
    return '<table>'+rows.map(function(r){
      var cells = r.split('|').filter(function(c,i,a){ return i>0 && i<a.length-1; });
      return '<tr>'+cells.map(function(c){ return '<td>'+c.trim()+'</td>'; }).join('')+'</tr>';
    }).join('')+'</table>';
  });
  // Unordered lists
  s = s.replace(/((?:^- .+$\n?)+)/gm, function(block) {
    var items = block.trim().split('\n');
    return '<ul>'+items.map(function(li){ return '<li>'+li.replace(/^- /,'')+'</li>'; }).join('')+'</ul>';
  });
  // Newlines (but not inside pre/table/ul)
  s = s.replace(/\n/g, '<br>');
  // Clean up extra <br> after block elements
  s = s.replace(/<\/(pre|h3|h4|h5|table|ul|li)><br>/g, '</$1>');
  s = s.replace(/<br><(pre|h3|h4|h5|table|ul)/g, '<$1');
  return s;
}
function fmtTime(t) { return t?new Date(t).toLocaleTimeString([],{hour:'2-digit',minute:'2-digit',second:'2-digit'}):''; }

// ── Content Viewer Modal ──
function openViewer(title, msgIdx, agentKey) {
  const msgs = S.messages.filter(m => normRole(m.from)===agentKey || normRole(m.to)===agentKey);
  const m = msgs[parseInt(msgIdx)];
  if(!m) return;
  document.getElementById('viewerTitle').textContent = (m.from||'Message') + ' — ' + fmtTime(m.timestamp);
  document.getElementById('viewerBody').innerHTML = renderMd(m.content||'');
  document.getElementById('viewerOverlay').style.display = '';
}

function openViewerRaw(title, content) {
  document.getElementById('viewerTitle').textContent = title;
  document.getElementById('viewerBody').innerHTML = renderMd(content);
  document.getElementById('viewerOverlay').style.display = '';
}

function closeViewer() {
  document.getElementById('viewerOverlay').style.display = 'none';
}

// ── Agent Chat ──
async function loadChatHistory(agentKey) {
  const el = document.getElementById('chatMessages');
  if(!el) return;
  const ag = AGENTS[agentKey]; if(!ag) return;

  try {
    const d = await (await fetch('/api/agents/'+agentKey+'/chat?project='+S.projectId)).json();
    const msgs = d.messages || [];
    el.innerHTML = msgs.map(m => {
      const isUser = m.direction === 'user_to_agent';
      return '<div class="chat-bubble '+(isUser?'user':'agent')+'">'+
        '<div class="cb-from" style="color:'+(isUser?'var(--accent)':ag.color)+'">'+(isUser?'You':ag.name)+'</div>'+
        '<div class="cb-text">'+renderMd(m.content)+'</div>'+
        '<div class="cb-time">'+fmtTime(m.timestamp)+'</div></div>';
    }).join('');
    el.scrollTop = el.scrollHeight;
  } catch(e) {
    el.innerHTML = '<div style="color:var(--text-muted);font-size:11px;padding:8px">Start a conversation with '+ag.name+'</div>';
  }
}

async function sendChat(agentKey) {
  const inp = document.getElementById('chatInput');
  const btn = document.getElementById('chatSend');
  const content = inp.value.trim();
  if(!content || !agentKey) return;

  inp.value = ''; btn.disabled = true;

  // Add user bubble immediately
  const el = document.getElementById('chatMessages');
  if(el) {
    const ag = AGENTS[agentKey];
    el.innerHTML +=
      '<div class="chat-bubble user">'+
      '<div class="cb-from" style="color:var(--accent)">You</div>'+
      '<div class="cb-text">'+renderMd(content)+'</div>'+
      '<div class="cb-time">'+fmtTime(new Date().toISOString())+'</div></div>';

    // Add typing indicator
    el.innerHTML +=
      '<div class="chat-bubble agent typing" id="chatTyping">'+
      '<div class="cb-from" style="color:'+(ag?.color||'#888')+'">'+(ag?.name||agentKey)+'</div>'+
      '<div class="cb-text">Thinking</div></div>';
    el.scrollTop = el.scrollHeight;
  }

  try {
    await fetch('/api/agents/'+agentKey+'/chat?project='+S.projectId, {
      method: 'POST', headers: {'Content-Type':'application/json'},
      body: JSON.stringify({content, context: S.running ? 'mid_workflow' : 'idle_chat'})
    });
  } catch(e) { console.error('[CHAT]', e); }

  btn.disabled = false;
  inp.focus();
}

function handleChatWS(evt) {
  if(!evt.data) return;
  const agentKey = evt.checkpoint_type; // We reuse checkpoint_type field for agent role
  const direction = evt.data.direction;
  const content = evt.data.content;

  // Only update if we're viewing this agent's chat
  if(S.selected !== agentKey || S.activeTab !== 'chat') return;

  const el = document.getElementById('chatMessages');
  if(!el) return;

  // Remove typing indicator
  const typing = document.getElementById('chatTyping');
  if(typing) typing.remove();

  if(direction === 'agent_to_user') {
    const ag = AGENTS[agentKey];
    el.innerHTML +=
      '<div class="chat-bubble agent">'+
      '<div class="cb-from" style="color:'+(ag?.color||'#888')+'">'+(ag?.name||agentKey)+'</div>'+
      '<div class="cb-text">'+renderMd(content)+'</div>'+
      '<div class="cb-time">'+fmtTime(new Date().toISOString())+'</div></div>';
    el.scrollTop = el.scrollHeight;
  }
}

// ── Toast Notifications ──
function showToast(title, text, type) {
  type = type || 'info';
  const icons = {info:'&#9432;', success:'&#10003;', warning:'&#9888;', error:'&#10007;'};
  const container = document.getElementById('toastContainer');
  const el = document.createElement('div');
  el.className = 'toast ' + type;
  el.innerHTML = '<span class="toast-icon">'+icons[type]+'</span><div class="toast-body"><div class="toast-title">'+esc(title)+'</div><div class="toast-text">'+esc(text)+'</div></div>';
  container.appendChild(el);
  setTimeout(() => { el.style.animation='toast-out 0.3s ease forwards'; setTimeout(()=>el.remove(),300); }, 5000);
}

function toggleLog() {
  document.getElementById('logArea').classList.toggle('collapsed');
}

// ═══════════════════════════════════════════
// NEW PROJECT MODAL (reuse pattern from office)
// ═══════════════════════════════════════════
function openNewProject() {
  const name = prompt('New project name (lowercase, hyphens OK):');
  if(!name) return;
  const clean = name.toLowerCase().replace(/\\s+/g,'-').replace(/[^a-z0-9\\-_]/g,'');
  if(!clean) return;
  const sel = document.getElementById('projectSel');
  if(!Array.from(sel.options).some(o=>o.value===clean)) {
    const opt = document.createElement('option'); opt.value=clean; opt.textContent=clean; sel.appendChild(opt);
  }
  sel.value=clean; S.projectId=clean;
  S.agents={}; S.messages=[]; S.msgCount=0;
  document.getElementById('logStream').innerHTML='';
  renderNodes(); renderRoster(); updateKPIs(); fetchAgentTasks();
}

// ═══════════════════════════════════════════
// INIT
// ═══════════════════════════════════════════
function init() {
  renderPipeline(); renderNodes(); renderRoster(); updateKPIs(); loadSettings();
  connectWS(); fetchProjects(); fetchAgentModes();

  // Event listeners
  document.getElementById('taskSubmit').addEventListener('click', submitTask);
  document.getElementById('taskInput').addEventListener('keydown', e => { if(e.key==='Enter') submitTask(); });
  document.getElementById('dpClose').addEventListener('click', closePanel);
  document.getElementById('logToggle').addEventListener('click', () => switchLogTab('log'));
  document.getElementById('projectToggle').addEventListener('click', () => switchLogTab('project'));
  document.getElementById('injectBtn').addEventListener('click', openInject);
  document.getElementById('newProjBtn').addEventListener('click', openNewProject);

  // Checkpoint handlers
  document.getElementById('cpApprove').addEventListener('click', () => resolveCheckpoint('approved'));
  document.getElementById('cpReject').addEventListener('click', () => openFeedback('reject'));
  document.getElementById('cpOverride').addEventListener('click', () => openFeedback('override'));
  document.getElementById('cpToggle').addEventListener('click', () => {
    // Open full content in the viewer modal
    const content = document.getElementById('cpContent').dataset.raw || '';
    const title = document.getElementById('cpTitle').textContent || 'Checkpoint Details';
    if(content) openViewerRaw(title, content);
  });

  // Settings handlers
  document.getElementById('settingsBtn').addEventListener('click', openSettings);
  document.getElementById('settingsClose').addEventListener('click', closeSettings);
  document.getElementById('settingsCancel').addEventListener('click', closeSettings);
  document.getElementById('settingsSave').addEventListener('click', saveSettings);

  // Feedback handlers
  document.getElementById('feedbackClose').addEventListener('click', closeFeedback);
  document.getElementById('feedbackCancel').addEventListener('click', closeFeedback);
  document.getElementById('feedbackSubmit').addEventListener('click', submitFeedback);
  document.getElementById('viewerClose').addEventListener('click', closeViewer);
  document.getElementById('viewerOverlay').addEventListener('click', e => { if(e.target.id==='viewerOverlay') closeViewer(); });

  // View toggle
  document.querySelectorAll('.view-btn').forEach(btn => {
    btn.addEventListener('click', () => switchView(btn.dataset.view));
  });

  // Kanban
  document.getElementById('kbSyncBtn').addEventListener('click', syncKanban);
  document.getElementById('kbFilterAgent').addEventListener('change', fetchKanbanTasks);
  document.getElementById('kbFilterPhase').addEventListener('change', fetchKanbanTasks);
  document.getElementById('kbFilterPriority').addEventListener('change', fetchKanbanTasks);
  document.getElementById('kbSearch').addEventListener('input', renderKanban);
  populateKbFilters();
  initKanbanDragDrop();

  document.getElementById('cmdBtn').addEventListener('click', openCmd);
  document.getElementById('projectSel').addEventListener('change', e => {
    S.projectId=e.target.value; S.agents={}; S.messages=[]; S.msgCount=0;
    S.running=false; S.startTime=null;
    document.getElementById('logStream').innerHTML='';
    closePanel();
    renderNodes(); renderRoster(); updateKPIs();
    // Fetch current state for the new project
    fetchStatus(); fetchAgentTasks(); loadSettings();
    if(currentView==='kanban') fetchKanbanTasks();
  });
  document.getElementById('canvasArea').addEventListener('click', e => {
    if(!e.target.closest('.canvas-node')) closePanel();
  });

  // Panel tabs
  document.querySelectorAll('.dp-tab').forEach(tab => {
    tab.addEventListener('click', () => renderPanelTab(tab.dataset.tab));
  });

  // Command palette
  document.getElementById('cmdOverlay').addEventListener('click', e => { if(e.target.id==='cmdOverlay') closeCmd(); });
  document.getElementById('cmdInput').addEventListener('input', e => renderCmdResults(e.target.value.toLowerCase()));
  document.getElementById('cmdInput').addEventListener('keydown', e => { if(e.key==='Escape') closeCmd(); });

  // Keyboard shortcuts
  document.addEventListener('keydown', e => {
    if((e.metaKey||e.ctrlKey) && e.key==='k') { e.preventDefault(); openCmd(); }
    if(e.key==='Escape') { closeCmd(); closePanel(); closeViewer(); }
    if((e.metaKey||e.ctrlKey) && e.key==='j') { e.preventDefault(); toggleLog(); }
    // 1-8 focus agents
    const num = parseInt(e.key);
    if(num>=1 && num<=8 && !e.metaKey && !e.ctrlKey && document.activeElement.tagName!=='INPUT') {
      const keys = Object.keys(AGENTS);
      if(keys[num-1]) selectAgent(keys[num-1]);
    }
  });

  // Polling
  setInterval(fetchStatus, 3000);
  setInterval(fetchAgentTasks, 5000);
  setInterval(updateKPIs, 1000);
}

document.addEventListener('DOMContentLoaded', init);
</script>
</body>
</html>` + "\n"
