package web

// dashboardHTML contains the embedded dashboard HTML with WebSocket support
var dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Agent House - Multi-Agent Dashboard</title>
    <style>
        /* ===================================
           DESIGN SYSTEM
        =================================== */
        :root {
            /* Backgrounds */
            --color-bg-primary: #0a0a1a;
            --color-bg-secondary: #12122a;
            --color-bg-tertiary: #1a1a3e;
            --color-bg-elevated: #252550;

            /* Text */
            --color-text-primary: #f0f0f5;
            --color-text-secondary: #a0a0b8;
            --color-text-tertiary: #6a6a80;

            /* Borders */
            --color-border: #2a2a50;
            --color-border-focus: #4ECDC4;

            /* Status Colors */
            --color-status-success: #2ECC71;
            --color-status-warning: #F39C12;
            --color-status-error: #E74C3C;
            --color-status-info: #3498DB;

            /* Agent Colors */
            --color-agent-ceo: #4A90A4;
            --color-agent-pm: #7B68EE;
            --color-agent-ux: #FF6B6B;
            --color-agent-ui: #4ECDC4;
            --color-agent-security: #F39C12;
            --color-agent-architect: #9B59B6;
            --color-agent-senior-dev: #2ECC71;
            --color-agent-junior-dev: #3498DB;

            /* Phase Colors */
            --color-phase-template: #FF9F43;
            --color-phase-research: #54A0FF;
            --color-phase-planning: #5F27CD;
            --color-phase-discussion: #00D2D3;
            --color-phase-development: #10AC84;

            /* Accent */
            --color-accent: #4ECDC4;
            --color-accent-hover: #45b7aa;

            /* Typography Scale */
            --text-xs: 0.64rem;
            --text-sm: 0.8rem;
            --text-base: 1rem;
            --text-lg: 1.25rem;
            --text-xl: 1.563rem;
            --text-2xl: 1.953rem;

            /* Spacing (8px base) */
            --space-1: 0.25rem;
            --space-2: 0.5rem;
            --space-3: 0.75rem;
            --space-4: 1rem;
            --space-6: 1.5rem;
            --space-8: 2rem;

            /* Border Radius */
            --radius-sm: 4px;
            --radius-md: 8px;
            --radius-lg: 12px;

            /* Transitions */
            --transition-fast: 0.15s ease;
            --transition-normal: 0.25s ease;

            /* Shadows */
            --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.3);
            --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.4);
            --shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.5);

            /* Layout */
            --sidebar-width: 280px;
            --context-width: 300px;
            --header-height: 60px;
        }

        /* Wide screens */
        @media (min-width: 1280px) {
            :root {
                --sidebar-width: 320px;
                --context-width: 360px;
            }
        }

        /* ===================================
           RESET & BASE STYLES
        =================================== */
        *, *::before, *::after {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        html {
            font-size: 16px;
            scroll-behavior: smooth;
        }

        body {
            font-family: 'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
            background: var(--color-bg-primary);
            color: var(--color-text-primary);
            min-height: 100vh;
            line-height: 1.5;
        }

        /* ===================================
           ACCESSIBILITY
        =================================== */
        .skip-link {
            position: absolute;
            top: -40px;
            left: 0;
            background: var(--color-accent);
            color: var(--color-bg-primary);
            padding: var(--space-2) var(--space-4);
            z-index: 1000;
            border-radius: 0 0 var(--radius-md) 0;
            font-weight: 600;
            text-decoration: none;
        }

        .skip-link:focus {
            top: 0;
        }

        .sr-only {
            position: absolute;
            width: 1px;
            height: 1px;
            padding: 0;
            margin: -1px;
            overflow: hidden;
            clip: rect(0, 0, 0, 0);
            white-space: nowrap;
            border: 0;
        }

        /* Focus visible styles */
        :focus-visible {
            outline: 2px solid var(--color-border-focus);
            outline-offset: 2px;
        }

        /* Reduced motion */
        @media (prefers-reduced-motion: reduce) {
            *, *::before, *::after {
                animation-duration: 0.01ms !important;
                animation-iteration-count: 1 !important;
                transition-duration: 0.01ms !important;
            }
        }

        /* ===================================
           LAYOUT
        =================================== */
        .app-container {
            display: flex;
            flex-direction: column;
            min-height: 100vh;
        }

        .main-layout {
            display: grid;
            grid-template-columns: 1fr;
            flex: 1;
            overflow: hidden;
        }

        /* Tablet: sidebar + main */
        @media (min-width: 768px) {
            .main-layout {
                grid-template-columns: var(--sidebar-width) 1fr;
            }
        }

        /* Desktop: 3-column */
        @media (min-width: 1024px) {
            .main-layout {
                grid-template-columns: var(--sidebar-width) 1fr var(--context-width);
            }
        }

        /* ===================================
           DEVELOPMENT PHASE TIMELINE
        =================================== */
        .dev-phases-section {
            background: var(--color-bg-secondary);
            border: 1px solid var(--color-border);
            border-radius: var(--radius-md);
            padding: var(--space-6);
            margin: var(--space-4) 0;
        }

        .dev-phases-section h2 {
            margin: 0 0 var(--space-4) 0;
            color: var(--color-text-primary);
            font-size: var(--text-lg);
            font-weight: 600;
        }

        .dev-phases-timeline {
            display: flex;
            gap: var(--space-4);
            align-items: center;
            overflow-x: auto;
            padding: var(--space-2) 0;
        }

        .phase-card {
            flex: 1;
            min-width: 200px;
            background: var(--color-bg-tertiary);
            border: 2px solid var(--color-border);
            border-radius: var(--radius-md);
            padding: var(--space-4);
            cursor: pointer;
            transition: all var(--transition-normal);
        }

        .phase-card:hover {
            border-color: var(--color-border-focus);
            transform: translateY(-2px);
        }

        .phase-card.completed {
            border-color: var(--color-status-success);
            background: rgba(46, 204, 113, 0.1);
        }

        .phase-card.in-progress {
            border-color: var(--color-status-info);
            background: rgba(52, 152, 219, 0.1);
            animation: pulse 2s ease-in-out infinite;
        }

        .phase-card.rejected {
            border-color: var(--color-status-error);
            background: rgba(231, 76, 60, 0.1);
        }

        .phase-card.pending {
            opacity: 0.6;
        }

        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.7; }
        }

        .phase-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: var(--space-3);
        }

        .phase-badge {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: bold;
            font-size: 14px;
        }

        .phase-card.completed .phase-badge {
            background: var(--color-status-success);
            color: white;
        }

        .phase-card.in-progress .phase-badge {
            background: var(--color-status-info);
            color: white;
        }

        .phase-card.rejected .phase-badge {
            background: var(--color-status-error);
            color: white;
        }

        .phase-card.pending .phase-badge {
            background: var(--color-bg-elevated);
            color: var(--color-text-secondary);
        }

        .phase-name {
            font-size: var(--text-base);
            font-weight: 600;
            color: var(--color-text-primary);
            margin-bottom: var(--space-2);
        }

        .phase-progress {
            width: 100%;
            height: 4px;
            background: var(--color-bg-elevated);
            border-radius: 2px;
            margin: var(--space-3) 0;
            overflow: hidden;
        }

        .phase-progress-fill {
            height: 100%;
            background: var(--color-status-success);
            transition: width var(--transition-normal);
        }

        .phase-card.in-progress .phase-progress-fill {
            background: var(--color-status-info);
        }

        .phase-card.rejected .phase-progress-fill {
            background: var(--color-status-error);
        }

        .phase-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: var(--text-xs);
            color: var(--color-text-secondary);
            margin-top: var(--space-2);
        }

        .qa-status {
            display: inline-flex;
            align-items: center;
            gap: var(--space-1);
            padding: var(--space-1) var(--space-2);
            border-radius: var(--radius-sm);
            font-size: 11px;
            font-weight: 600;
        }

        .qa-status.approved {
            background: rgba(46, 204, 113, 0.2);
            color: var(--color-status-success);
        }

        .qa-status.rejected {
            background: rgba(231, 76, 60, 0.2);
            color: var(--color-status-error);
        }

        .qa-status.pending {
            background: rgba(160, 160, 184, 0.2);
            color: var(--color-text-secondary);
        }

        .iteration-badge {
            background: var(--color-status-warning);
            color: white;
            padding: 2px 6px;
            border-radius: var(--radius-sm);
            font-size: 10px;
            font-weight: bold;
        }

        .iteration-badge.warning {
            background: var(--color-status-error);
            animation: blink 1s ease-in-out infinite;
        }

        @keyframes blink {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }

        .timeline-connector {
            width: 24px;
            height: 2px;
            background: var(--color-border);
            flex-shrink: 0;
        }

        /* ===================================
           QA FEEDBACK PANEL
        =================================== */
        .qa-feedback-panel {
            background: var(--color-bg-secondary);
            border: 1px solid var(--color-border);
            border-radius: var(--radius-md);
            padding: var(--space-6);
            margin: var(--space-4) 0;
        }

        .qa-feedback-panel .panel-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: var(--space-4);
        }

        .qa-feedback-panel .panel-header h3 {
            margin: 0;
            color: var(--color-text-primary);
            font-size: var(--text-lg);
            font-weight: 600;
        }

        .qa-status-banner {
            display: flex;
            align-items: center;
            gap: var(--space-3);
            padding: var(--space-4);
            border-radius: var(--radius-md);
            margin-bottom: var(--space-4);
            border: 2px solid transparent;
        }

        .qa-status-banner.approved {
            background: rgba(46, 204, 113, 0.15);
            border-color: var(--color-status-success);
        }

        .qa-status-banner.rejected {
            background: rgba(231, 76, 60, 0.15);
            border-color: var(--color-status-error);
        }

        .qa-status-banner .status-icon {
            font-size: 32px;
            line-height: 1;
        }

        .qa-status-banner .status-text {
            font-size: var(--text-xl);
            font-weight: bold;
            color: var(--color-text-primary);
            flex: 1;
        }

        .qa-status-banner .reviewer-info {
            font-size: var(--text-sm);
            color: var(--color-text-secondary);
        }

        .qa-feedback-content {
            color: var(--color-text-secondary);
            line-height: 1.6;
        }

        .failed-criteria {
            margin-top: var(--space-4);
        }

        .failed-criteria h4 {
            color: var(--color-text-primary);
            margin-bottom: var(--space-3);
            font-size: var(--text-base);
            font-weight: 600;
        }

        .criterion-item {
            display: flex;
            gap: var(--space-3);
            padding: var(--space-3);
            background: var(--color-bg-tertiary);
            border-left: 3px solid var(--color-status-error);
            border-radius: var(--radius-sm);
            margin-bottom: var(--space-2);
        }

        .criterion-item .icon {
            font-size: 20px;
            line-height: 1;
            flex-shrink: 0;
        }

        .criterion-item .content {
            flex: 1;
        }

        .criterion-item strong {
            color: var(--color-text-primary);
            display: block;
            margin-bottom: var(--space-1);
        }

        .criterion-item p {
            margin: var(--space-1) 0;
            font-size: var(--text-sm);
        }

        .file-link {
            color: var(--color-status-info);
            text-decoration: none;
            font-family: 'Fira Code', 'Consolas', monospace;
            font-size: var(--text-xs);
            display: inline-block;
            margin-top: var(--space-1);
            padding: 2px 6px;
            background: rgba(52, 152, 219, 0.1);
            border-radius: var(--radius-sm);
        }

        .file-link:hover {
            background: rgba(52, 152, 219, 0.2);
            text-decoration: underline;
        }

        .required-fixes {
            margin-top: var(--space-4);
        }

        .required-fixes h4 {
            color: var(--color-text-primary);
            margin-bottom: var(--space-3);
            font-size: var(--text-base);
            font-weight: 600;
        }

        .required-fixes ol {
            margin: 0;
            padding-left: var(--space-6);
        }

        .required-fixes li {
            color: var(--color-text-secondary);
            margin-bottom: var(--space-2);
            line-height: 1.6;
        }

        .feedback-actions {
            margin-top: var(--space-4);
            display: flex;
            gap: var(--space-3);
        }

        .btn-secondary {
            background: var(--color-bg-tertiary);
            border: 1px solid var(--color-border);
            color: var(--color-text-primary);
            padding: var(--space-2) var(--space-4);
            border-radius: var(--radius-sm);
            cursor: pointer;
            font-size: var(--text-sm);
            transition: all var(--transition-fast);
        }

        .btn-secondary:hover {
            background: var(--color-bg-elevated);
            border-color: var(--color-border-focus);
        }

        /* ===================================
           PROGRESS DASHBOARD
        =================================== */
        .progress-dashboard {
            background: var(--color-bg-secondary);
            border: 1px solid var(--color-border);
            border-radius: var(--radius-md);
            padding: var(--space-6);
            margin-bottom: var(--space-6);
        }

        .progress-dashboard h2 {
            font-size: var(--text-xl);
            font-weight: 600;
            color: var(--color-text-primary);
            margin: 0 0 var(--space-5) 0;
        }

        .metrics-grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: var(--space-4);
        }

        .metric-card {
            background: var(--color-bg-primary);
            border: 1px solid var(--color-border);
            border-radius: var(--radius-sm);
            padding: var(--space-5);
            display: flex;
            flex-direction: column;
            gap: var(--space-3);
            transition: all var(--transition-base);
        }

        .metric-card:hover {
            border-color: var(--color-border-focus);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }

        .metric-header {
            display: flex;
            align-items: center;
            gap: var(--space-2);
        }

        .metric-icon {
            font-size: var(--text-2xl);
            line-height: 1;
        }

        .metric-label {
            font-size: var(--text-sm);
            color: var(--color-text-secondary);
            font-weight: 500;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        .metric-value {
            font-size: var(--text-3xl);
            font-weight: 700;
            color: var(--color-text-primary);
            line-height: 1.2;
        }

        .metric-value.success {
            color: var(--color-status-success);
        }

        .metric-value.warning {
            color: var(--color-status-warning);
        }

        .metric-value.error {
            color: var(--color-status-error);
        }

        .metric-value.info {
            color: var(--color-status-info);
        }

        .metric-description {
            font-size: var(--text-sm);
            color: var(--color-text-secondary);
            line-height: 1.5;
        }

        .metric-badge {
            display: inline-flex;
            align-items: center;
            gap: var(--space-2);
            padding: var(--space-1) var(--space-3);
            border-radius: var(--radius-full);
            font-size: var(--text-xs);
            font-weight: 600;
            background: var(--color-bg-tertiary);
            color: var(--color-text-secondary);
        }

        .metric-badge.success {
            background: rgba(46, 204, 113, 0.15);
            color: var(--color-status-success);
        }

        .metric-badge.warning {
            background: rgba(241, 196, 15, 0.15);
            color: var(--color-status-warning);
        }

        .metric-badge.error {
            background: rgba(231, 76, 60, 0.15);
            color: var(--color-status-error);
        }

        .metric-badge.info {
            background: rgba(52, 152, 219, 0.15);
            color: var(--color-status-info);
        }

        .progress-ring {
            width: 80px;
            height: 80px;
            margin: 0 auto;
        }

        .progress-ring-circle {
            fill: transparent;
            stroke: var(--color-border);
            stroke-width: 8;
        }

        .progress-ring-fill {
            fill: transparent;
            stroke: var(--color-status-success);
            stroke-width: 8;
            stroke-linecap: round;
            transform: rotate(-90deg);
            transform-origin: 50% 50%;
            transition: stroke-dashoffset 0.5s ease;
        }

        /* Mobile responsive */
        @media (max-width: 640px) {
            .metrics-grid {
                grid-template-columns: 1fr;
            }

            .metric-value {
                font-size: var(--text-2xl);
            }
        }

        /* ===================================
           HEADER
        =================================== */
        .header {
            background: var(--color-bg-secondary);
            border-bottom: 1px solid var(--color-border);
            height: var(--header-height);
            display: flex;
            flex-direction: column;
            padding: 0 var(--space-4);
        }

        .header-main {
            display: flex;
            justify-content: space-between;
            align-items: center;
            height: 50%;
        }

        .logo {
            display: flex;
            align-items: center;
            gap: var(--space-2);
        }

        .logo h1 {
            font-size: var(--text-lg);
            font-weight: 700;
            color: var(--color-accent);
            letter-spacing: -0.02em;
        }

        .logo-icon {
            width: 28px;
            height: 28px;
            background: linear-gradient(135deg, var(--color-accent), var(--color-agent-pm));
            border-radius: var(--radius-sm);
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: var(--text-sm);
        }

        .status-bar {
            display: flex;
            gap: var(--space-6);
            font-size: var(--text-sm);
        }

        .status-item {
            display: flex;
            align-items: center;
            gap: var(--space-2);
            color: var(--color-text-secondary);
        }

        .status-dot {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background: var(--color-text-tertiary);
        }

        .status-dot.connected { background: var(--color-status-success); }
        .status-dot.running {
            background: var(--color-status-warning);
            animation: pulse 1s infinite;
        }

        @keyframes pulse {
            0%, 100% { opacity: 1; transform: scale(1); }
            50% { opacity: 0.6; transform: scale(0.95); }
        }

        /* ===================================
           PHASE STEPPER
        =================================== */
        .phase-stepper {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: var(--space-1);
            height: 50%;
            padding: 0 var(--space-4);
            overflow-x: auto;
        }

        .phase-step {
            display: flex;
            align-items: center;
            gap: var(--space-2);
            padding: var(--space-1) var(--space-3);
            border-radius: var(--radius-lg);
            font-size: var(--text-xs);
            background: transparent;
            color: var(--color-text-tertiary);
            transition: var(--transition-fast);
            white-space: nowrap;
        }

        .phase-step.completed {
            color: var(--color-status-success);
        }

        .phase-step.active {
            background: var(--color-bg-elevated);
            color: var(--color-text-primary);
            box-shadow: 0 0 12px rgba(78, 205, 196, 0.3);
        }

        .phase-icon {
            width: 22px;
            height: 22px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: var(--text-xs);
            background: var(--color-bg-tertiary);
            border: 2px solid var(--color-border);
        }

        .phase-step.completed .phase-icon {
            background: var(--color-status-success);
            border-color: var(--color-status-success);
            color: white;
        }

        .phase-step.active .phase-icon {
            border-color: var(--color-accent);
            box-shadow: 0 0 8px var(--color-accent);
        }

        .phase-step[data-phase="template_selection"] .phase-icon { background: var(--color-phase-template); border-color: var(--color-phase-template); }
        .phase-step[data-phase="research"] .phase-icon { background: var(--color-phase-research); border-color: var(--color-phase-research); }
        .phase-step[data-phase="planning"] .phase-icon { background: var(--color-phase-planning); border-color: var(--color-phase-planning); }
        .phase-step[data-phase="discussion"] .phase-icon { background: var(--color-phase-discussion); border-color: var(--color-phase-discussion); }
        .phase-step[data-phase="development"] .phase-icon { background: var(--color-phase-development); border-color: var(--color-phase-development); }

        .phase-step.active .phase-icon,
        .phase-step.completed .phase-icon {
            color: white;
        }

        .phase-connector {
            width: 24px;
            height: 2px;
            background: var(--color-border);
        }

        .phase-connector.completed {
            background: var(--color-status-success);
        }

        /* Hide phase labels on small screens */
        @media (max-width: 640px) {
            .phase-label { display: none; }
            .phase-step { padding: var(--space-1); }
        }

        /* ===================================
           SIDEBAR
        =================================== */
        .sidebar {
            background: var(--color-bg-secondary);
            border-right: 1px solid var(--color-border);
            display: none;
            flex-direction: column;
            overflow: hidden;
        }

        @media (min-width: 768px) {
            .sidebar { display: flex; }
        }

        .sidebar-header {
            padding: var(--space-4);
            border-bottom: 1px solid var(--color-border);
        }

        .sidebar-title {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: var(--space-3);
        }

        .sidebar-title h2 {
            font-size: var(--text-sm);
            font-weight: 600;
            color: var(--color-text-secondary);
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        .project-select {
            width: 100%;
            padding: var(--space-2) var(--space-3);
            background: var(--color-bg-tertiary);
            border: 1px solid var(--color-border);
            border-radius: var(--radius-md);
            color: var(--color-text-primary);
            font-size: var(--text-sm);
            cursor: pointer;
        }

        .project-select:focus {
            border-color: var(--color-accent);
        }

        .search-input {
            display: flex;
            align-items: center;
            gap: var(--space-2);
            padding: var(--space-2) var(--space-3);
            background: var(--color-bg-tertiary);
            border: 1px solid var(--color-border);
            border-radius: var(--radius-md);
            margin-top: var(--space-3);
        }

        .search-input input {
            flex: 1;
            background: transparent;
            border: none;
            color: var(--color-text-primary);
            font-size: var(--text-sm);
        }

        .search-input input:focus { outline: none; }

        .search-input input::placeholder {
            color: var(--color-text-tertiary);
        }

        .search-shortcut {
            font-size: var(--text-xs);
            color: var(--color-text-tertiary);
            background: var(--color-bg-secondary);
            padding: 2px 6px;
            border-radius: var(--radius-sm);
        }

        .sidebar-content {
            flex: 1;
            overflow-y: auto;
            padding: var(--space-3);
        }

        .new-task-btn {
            width: 100%;
            padding: var(--space-3);
            border: 1px dashed var(--color-border);
            border-radius: var(--radius-md);
            background: transparent;
            color: var(--color-accent);
            cursor: pointer;
            font-size: var(--text-sm);
            font-weight: 500;
            margin-bottom: var(--space-3);
            transition: var(--transition-fast);
            display: flex;
            align-items: center;
            justify-content: center;
            gap: var(--space-2);
        }

        .new-task-btn:hover {
            background: var(--color-bg-tertiary);
            border-style: solid;
        }

        /* ===================================
           TASK CARDS
        =================================== */
        .task-list {
            display: flex;
            flex-direction: column;
            gap: var(--space-2);
        }

        .task-card {
            background: var(--color-bg-tertiary);
            border-radius: var(--radius-md);
            padding: var(--space-3);
            border-left: 3px solid var(--color-border);
            cursor: pointer;
            transition: var(--transition-fast);
        }

        .task-card:hover {
            background: var(--color-bg-elevated);
        }

        .task-card.selected {
            border-left-color: var(--color-accent);
            background: var(--color-bg-elevated);
        }

        .task-card.running {
            border-left-color: var(--color-status-warning);
        }

        .task-card.running::before {
            content: '';
            position: absolute;
            top: var(--space-3);
            right: var(--space-3);
            width: 8px;
            height: 8px;
            background: var(--color-status-warning);
            border-radius: 50%;
            animation: pulse 1s infinite;
        }

        .task-date {
            font-size: var(--text-xs);
            color: var(--color-text-tertiary);
            margin-bottom: var(--space-1);
        }

        .task-summary {
            font-size: var(--text-sm);
            color: var(--color-text-primary);
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            margin-bottom: var(--space-2);
        }

        .task-meta {
            display: flex;
            flex-wrap: wrap;
            gap: var(--space-2);
            font-size: var(--text-xs);
            color: var(--color-text-tertiary);
        }

        .task-badge {
            padding: 2px 8px;
            border-radius: 10px;
            font-weight: 500;
            text-transform: uppercase;
            font-size: 10px;
            letter-spacing: 0.03em;
        }

        .task-badge.completed { background: var(--color-status-success); color: #000; }
        .task-badge.running { background: var(--color-status-warning); color: #000; animation: pulse 1s infinite; }
        .task-badge.failed { background: var(--color-status-error); color: #fff; }
        .task-badge.pending { background: var(--color-text-tertiary); color: #fff; }

        /* ===================================
           MAIN CONTENT
        =================================== */
        .main-content {
            display: flex;
            flex-direction: column;
            overflow: hidden;
            background: var(--color-bg-primary);
        }

        /* Task Input Area */
        .task-input-area {
            padding: var(--space-4);
            background: var(--color-bg-secondary);
            border-bottom: 1px solid var(--color-border);
        }

        .task-form {
            display: flex;
            flex-direction: column;
            gap: var(--space-3);
        }

        @media (min-width: 640px) {
            .task-form {
                flex-direction: row;
            }
        }

        .task-form input {
            flex: 1;
            padding: var(--space-3) var(--space-4);
            border: 1px solid var(--color-border);
            border-radius: var(--radius-md);
            background: var(--color-bg-tertiary);
            color: var(--color-text-primary);
            font-size: var(--text-base);
        }

        .task-form input:focus {
            outline: none;
            border-color: var(--color-accent);
        }

        .task-form input::placeholder {
            color: var(--color-text-tertiary);
        }

        .project-input-wrapper {
            display: flex;
            align-items: center;
            gap: var(--space-2);
        }

        @media (min-width: 640px) {
            .project-input-wrapper {
                max-width: 200px;
            }
        }

        .task-form button[type="submit"] {
            padding: var(--space-3) var(--space-6);
            border: none;
            border-radius: var(--radius-md);
            background: var(--color-accent);
            color: var(--color-bg-primary);
            font-weight: 600;
            font-size: var(--text-base);
            cursor: pointer;
            transition: var(--transition-fast);
            white-space: nowrap;
        }

        .task-form button[type="submit"]:hover {
            background: var(--color-accent-hover);
        }

        .task-form button[type="submit"]:disabled {
            background: var(--color-text-tertiary);
            cursor: not-allowed;
        }

        /* History Banner */
        .history-banner {
            background: var(--color-bg-tertiary);
            padding: var(--space-3) var(--space-4);
            border-bottom: 1px solid var(--color-border);
            display: none;
        }

        .history-banner.visible { display: block; }

        .history-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: var(--space-3);
        }

        .history-task-name {
            font-weight: 600;
            color: var(--color-text-primary);
        }

        .back-btn {
            background: transparent;
            color: var(--color-text-secondary);
            border: 1px solid var(--color-border);
            padding: var(--space-2) var(--space-3);
            border-radius: var(--radius-sm);
            cursor: pointer;
            font-size: var(--text-sm);
            transition: var(--transition-fast);
        }

        .back-btn:hover {
            background: var(--color-bg-secondary);
        }

        .continue-form {
            display: flex;
            gap: var(--space-3);
        }

        .continue-form input {
            flex: 1;
            padding: var(--space-2) var(--space-3);
            border: 1px solid var(--color-border);
            border-radius: var(--radius-md);
            background: var(--color-bg-secondary);
            color: var(--color-text-primary);
            font-size: var(--text-sm);
        }

        .continue-form input:focus {
            outline: none;
            border-color: var(--color-accent);
        }

        .continue-form button {
            padding: var(--space-2) var(--space-4);
            border: none;
            border-radius: var(--radius-md);
            background: var(--color-accent);
            color: var(--color-bg-primary);
            font-weight: 600;
            cursor: pointer;
            white-space: nowrap;
        }

        /* Delegation Breadcrumb */
        .delegation-breadcrumb {
            display: flex;
            align-items: center;
            gap: var(--space-2);
            padding: var(--space-2) var(--space-4);
            background: var(--color-bg-secondary);
            border-bottom: 1px solid var(--color-border);
            font-size: var(--text-sm);
            overflow-x: auto;
        }

        .delegation-breadcrumb:empty { display: none; }

        .breadcrumb-item {
            display: flex;
            align-items: center;
            gap: var(--space-1);
            padding: var(--space-1) var(--space-2);
            background: var(--color-bg-tertiary);
            border-radius: var(--radius-sm);
            color: var(--color-text-secondary);
            white-space: nowrap;
        }

        .breadcrumb-item.current {
            background: var(--color-accent);
            color: var(--color-bg-primary);
            font-weight: 600;
        }

        .breadcrumb-arrow {
            color: var(--color-text-tertiary);
        }

        .breadcrumb-arrow.review {
            color: var(--color-status-info);
            transform: rotate(180deg);
        }

        /* Messages Area */
        .messages-area {
            flex: 1;
            overflow-y: auto;
            padding: var(--space-4);
        }

        /* ===================================
           MESSAGE CARDS
        =================================== */
        .message-group {
            margin-bottom: var(--space-6);
        }

        .message-group-header {
            font-size: var(--text-xs);
            color: var(--color-text-tertiary);
            text-transform: uppercase;
            letter-spacing: 0.05em;
            padding: var(--space-2) 0;
            border-bottom: 1px solid var(--color-border);
            margin-bottom: var(--space-3);
        }

        .message {
            background: var(--color-bg-secondary);
            border-radius: var(--radius-md);
            padding: var(--space-4);
            margin-bottom: var(--space-3);
            border-left: 4px solid var(--color-border);
            animation: slideIn 0.3s ease;
        }

        @keyframes slideIn {
            from { opacity: 0; transform: translateY(-10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        /* Agent color borders */
        .message[data-agent="ceo"] { border-left-color: var(--color-agent-ceo); }
        .message[data-agent="pm"] { border-left-color: var(--color-agent-pm); }
        .message[data-agent="ux"] { border-left-color: var(--color-agent-ux); }
        .message[data-agent="ui"] { border-left-color: var(--color-agent-ui); }
        .message[data-agent="security"] { border-left-color: var(--color-agent-security); }
        .message[data-agent="architect"] { border-left-color: var(--color-agent-architect); }
        .message[data-agent="senior_dev"] { border-left-color: var(--color-agent-senior-dev); }
        .message[data-agent="junior_dev"] { border-left-color: var(--color-agent-junior-dev); }

        .message-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: var(--space-3);
            gap: var(--space-3);
        }

        .message-sender {
            display: flex;
            align-items: center;
            gap: var(--space-2);
        }

        .agent-avatar {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: var(--text-sm);
            background: var(--color-bg-tertiary);
        }

        .agent-avatar[data-agent="ceo"] { background: var(--color-agent-ceo); }
        .agent-avatar[data-agent="pm"] { background: var(--color-agent-pm); }
        .agent-avatar[data-agent="ux"] { background: var(--color-agent-ux); }
        .agent-avatar[data-agent="ui"] { background: var(--color-agent-ui); }
        .agent-avatar[data-agent="security"] { background: var(--color-agent-security); }
        .agent-avatar[data-agent="architect"] { background: var(--color-agent-architect); }
        .agent-avatar[data-agent="senior_dev"] { background: var(--color-agent-senior-dev); }
        .agent-avatar[data-agent="junior_dev"] { background: var(--color-agent-junior-dev); }

        .sender-info {
            display: flex;
            flex-direction: column;
        }

        .sender-name {
            font-weight: 600;
            font-size: var(--text-sm);
            color: var(--color-text-primary);
        }

        .sender-flow {
            font-size: var(--text-xs);
            color: var(--color-text-tertiary);
        }

        .message-meta {
            display: flex;
            align-items: center;
            gap: var(--space-2);
        }

        .message-type-badge {
            padding: 2px 8px;
            border-radius: 10px;
            font-size: 10px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.03em;
        }

        .message-type-badge.response { background: var(--color-bg-tertiary); color: var(--color-text-secondary); }
        .message-type-badge.delegate { background: var(--color-status-info); color: #fff; }
        .message-type-badge.file_create { background: var(--color-status-success); color: #000; }
        .message-type-badge.error { background: var(--color-status-error); color: #fff; }
        .message-type-badge.system { background: var(--color-status-warning); color: #000; }
        .message-type-badge.task { background: var(--color-agent-pm); color: #fff; }
        .message-type-badge.complete { background: var(--color-status-success); color: #000; }

        .message-time {
            font-size: var(--text-xs);
            color: var(--color-text-tertiary);
        }

        .message-content {
            white-space: pre-wrap;
            line-height: 1.6;
            font-size: var(--text-sm);
            color: var(--color-text-secondary);
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
            height: 60px;
            background: linear-gradient(transparent, var(--color-bg-secondary));
            pointer-events: none;
        }

        .expand-btn {
            background: none;
            border: none;
            color: var(--color-accent);
            cursor: pointer;
            font-size: var(--text-sm);
            margin-top: var(--space-2);
            padding: 0;
        }

        .expand-btn:hover {
            text-decoration: underline;
        }

        /* File message styling */
        .message.file-message {
            background: rgba(46, 204, 113, 0.1);
        }

        .file-path {
            font-family: 'JetBrains Mono', 'Fira Code', monospace;
            background: var(--color-bg-tertiary);
            padding: var(--space-1) var(--space-2);
            border-radius: var(--radius-sm);
            font-size: var(--text-sm);
        }

        /* ===================================
           CONTEXT PANEL
        =================================== */
        .context-panel {
            background: var(--color-bg-secondary);
            border-left: 1px solid var(--color-border);
            display: none;
            flex-direction: column;
            overflow: hidden;
        }

        @media (min-width: 1024px) {
            .context-panel { display: flex; }
        }

        .context-tabs {
            display: flex;
            border-bottom: 1px solid var(--color-border);
        }

        .context-tab {
            flex: 1;
            padding: var(--space-3);
            background: transparent;
            border: none;
            color: var(--color-text-tertiary);
            font-size: var(--text-sm);
            font-weight: 500;
            cursor: pointer;
            transition: var(--transition-fast);
            position: relative;
        }

        .context-tab:hover {
            color: var(--color-text-secondary);
        }

        .context-tab.active {
            color: var(--color-accent);
        }

        .context-tab.active::after {
            content: '';
            position: absolute;
            bottom: -1px;
            left: 0;
            right: 0;
            height: 2px;
            background: var(--color-accent);
        }

        .context-content {
            flex: 1;
            overflow-y: auto;
            padding: var(--space-4);
        }

        .context-section {
            display: none;
        }

        .context-section.active {
            display: block;
        }

        /* ===================================
           AGENT ACTIVITY PANEL
        =================================== */
        .agents-section h3 {
            font-size: var(--text-sm);
            color: var(--color-text-secondary);
            text-transform: uppercase;
            letter-spacing: 0.05em;
            margin-bottom: var(--space-4);
        }

        .agent-org-chart {
            display: flex;
            flex-direction: column;
            gap: var(--space-3);
        }

        .agent-tier {
            display: flex;
            flex-wrap: wrap;
            gap: var(--space-2);
            justify-content: center;
        }

        .agent-card {
            display: flex;
            align-items: center;
            gap: var(--space-2);
            padding: var(--space-2) var(--space-3);
            background: var(--color-bg-tertiary);
            border-radius: var(--radius-md);
            border: 1px solid var(--color-border);
            min-width: 100px;
        }

        .agent-card.active {
            border-color: var(--color-status-success);
            box-shadow: 0 0 8px rgba(46, 204, 113, 0.3);
        }

        .agent-card.working {
            border-color: var(--color-status-warning);
            animation: pulseGlow 1.5s infinite;
        }

        @keyframes pulseGlow {
            0%, 100% { box-shadow: 0 0 4px rgba(243, 156, 18, 0.3); }
            50% { box-shadow: 0 0 12px rgba(243, 156, 18, 0.5); }
        }

        .agent-card-avatar {
            width: 24px;
            height: 24px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
        }

        .agent-card-name {
            font-size: var(--text-xs);
            font-weight: 500;
        }

        .agent-status-indicator {
            width: 6px;
            height: 6px;
            border-radius: 50%;
            background: var(--color-text-tertiary);
            margin-left: auto;
        }

        .agent-card.active .agent-status-indicator { background: var(--color-status-success); }
        .agent-card.working .agent-status-indicator {
            background: var(--color-status-warning);
            animation: pulse 1s infinite;
        }

        .tier-connector {
            display: flex;
            justify-content: center;
            padding: var(--space-1) 0;
        }

        .tier-connector::after {
            content: '';
            width: 2px;
            height: 12px;
            background: var(--color-border);
        }

        /* ===================================
           FILES PANEL
        =================================== */
        .files-section h3 {
            font-size: var(--text-sm);
            color: var(--color-text-secondary);
            text-transform: uppercase;
            letter-spacing: 0.05em;
            margin-bottom: var(--space-3);
        }

        .file-tree {
            font-size: var(--text-sm);
        }

        .file-tree-folder {
            margin-bottom: var(--space-2);
        }

        .folder-header {
            display: flex;
            align-items: center;
            gap: var(--space-2);
            padding: var(--space-2);
            border-radius: var(--radius-sm);
            cursor: pointer;
            color: var(--color-text-secondary);
        }

        .folder-header:hover {
            background: var(--color-bg-tertiary);
        }

        .folder-icon {
            font-size: var(--text-base);
        }

        .folder-name {
            font-weight: 500;
        }

        .folder-children {
            margin-left: var(--space-4);
            border-left: 1px solid var(--color-border);
            padding-left: var(--space-2);
        }

        .file-item {
            display: flex;
            align-items: center;
            gap: var(--space-2);
            padding: var(--space-2);
            border-radius: var(--radius-sm);
            cursor: pointer;
        }

        .file-item:hover {
            background: var(--color-bg-tertiary);
        }

        .file-icon {
            opacity: 0.7;
        }

        .file-name {
            flex: 1;
            color: var(--color-text-primary);
        }

        .file-agent {
            font-size: var(--text-xs);
            padding: 2px 6px;
            border-radius: var(--radius-sm);
            background: var(--color-bg-tertiary);
            color: var(--color-text-tertiary);
        }

        /* ===================================
           TOAST NOTIFICATIONS
        =================================== */
        .toast-container {
            position: fixed;
            top: var(--space-4);
            right: var(--space-4);
            z-index: 1000;
            display: flex;
            flex-direction: column;
            gap: var(--space-2);
            pointer-events: none;
        }

        .toast {
            background: var(--color-bg-elevated);
            border: 1px solid var(--color-border);
            border-radius: var(--radius-md);
            padding: var(--space-3) var(--space-4);
            box-shadow: var(--shadow-lg);
            display: flex;
            align-items: center;
            gap: var(--space-3);
            min-width: 300px;
            max-width: 450px;
            pointer-events: auto;
            animation: slideInRight 0.3s ease;
        }

        @keyframes slideInRight {
            from { opacity: 0; transform: translateX(100px); }
            to { opacity: 1; transform: translateX(0); }
        }

        .toast.fade-out {
            animation: fadeOut 0.3s ease forwards;
        }

        @keyframes fadeOut {
            from { opacity: 1; }
            to { opacity: 0; }
        }

        .toast-icon {
            font-size: var(--text-lg);
        }

        .toast.success { border-left: 4px solid var(--color-status-success); }
        .toast.error { border-left: 4px solid var(--color-status-error); }
        .toast.warning { border-left: 4px solid var(--color-status-warning); }
        .toast.info { border-left: 4px solid var(--color-status-info); }

        .toast-content {
            flex: 1;
        }

        .toast-title {
            font-weight: 600;
            font-size: var(--text-sm);
            margin-bottom: var(--space-1);
        }

        .toast-message {
            font-size: var(--text-xs);
            color: var(--color-text-secondary);
        }

        .toast-close {
            background: transparent;
            border: none;
            color: var(--color-text-tertiary);
            cursor: pointer;
            font-size: var(--text-lg);
            padding: var(--space-1);
        }

        .toast-close:hover {
            color: var(--color-text-primary);
        }

        /* ===================================
           EMPTY STATES
        =================================== */
        .empty-state {
            text-align: center;
            padding: var(--space-8);
            color: var(--color-text-tertiary);
        }

        .empty-state-icon {
            font-size: 3rem;
            margin-bottom: var(--space-4);
            opacity: 0.5;
        }

        .empty-state-title {
            font-size: var(--text-lg);
            font-weight: 600;
            color: var(--color-text-secondary);
            margin-bottom: var(--space-2);
        }

        .empty-state-description {
            font-size: var(--text-sm);
            max-width: 300px;
            margin: 0 auto;
        }

        /* ===================================
           LOADING STATES
        =================================== */
        .skeleton {
            background: linear-gradient(90deg, var(--color-bg-tertiary) 25%, var(--color-bg-elevated) 50%, var(--color-bg-tertiary) 75%);
            background-size: 200% 100%;
            animation: shimmer 1.5s infinite;
            border-radius: var(--radius-sm);
        }

        @keyframes shimmer {
            0% { background-position: -200% 0; }
            100% { background-position: 200% 0; }
        }

        .skeleton-text {
            height: 14px;
            margin-bottom: var(--space-2);
        }

        .skeleton-title {
            height: 20px;
            width: 60%;
            margin-bottom: var(--space-3);
        }

        .skeleton-avatar {
            width: 32px;
            height: 32px;
            border-radius: 50%;
        }

        /* ===================================
           MOBILE SIDEBAR TOGGLE
        =================================== */
        .mobile-sidebar-toggle {
            display: flex;
            position: fixed;
            bottom: var(--space-4);
            left: var(--space-4);
            z-index: 100;
            width: 48px;
            height: 48px;
            border-radius: 50%;
            background: var(--color-accent);
            color: var(--color-bg-primary);
            border: none;
            align-items: center;
            justify-content: center;
            font-size: var(--text-lg);
            cursor: pointer;
            box-shadow: var(--shadow-lg);
        }

        @media (min-width: 768px) {
            .mobile-sidebar-toggle { display: none; }
        }

        .sidebar.mobile-open {
            display: flex;
            position: fixed;
            top: var(--header-height);
            left: 0;
            bottom: 0;
            width: 100%;
            max-width: var(--sidebar-width);
            z-index: 50;
        }

        .mobile-overlay {
            display: none;
            position: fixed;
            inset: 0;
            background: rgba(0, 0, 0, 0.5);
            z-index: 40;
        }

        .mobile-overlay.visible {
            display: block;
        }
    </style>
</head>
<body>
    <!-- Skip Link for Accessibility -->
    <a href="#main-content" class="skip-link">Skip to main content</a>

    <!-- Screen Reader Announcements -->
    <div role="status" aria-live="polite" class="sr-only" id="announcer"></div>

    <div class="app-container">
        <!-- Header -->
        <header class="header" role="banner">
            <div class="header-main">
                <div class="logo">
                    <div class="logo-icon">AH</div>
                    <h1>Agent House</h1>
                </div>
                <div class="status-bar">
                    <div class="status-item">
                        <div class="status-dot" id="wsStatus" aria-hidden="true"></div>
                        <span id="wsStatusText">Connecting...</span>
                    </div>
                    <div class="status-item">
                        <div class="status-dot" id="taskStatus" aria-hidden="true"></div>
                        <span id="taskStatusText">Idle</span>
                    </div>
                    <div class="status-item" aria-label="Message count">
                        <span id="messageCount">0 messages</span>
                    </div>
                </div>
            </div>

            <!-- Phase Progress Stepper -->
            <nav class="phase-stepper" role="navigation" aria-label="Workflow phases" id="phaseStepper">
                <div class="phase-step" data-phase="template_selection">
                    <div class="phase-icon" aria-hidden="true">1</div>
                    <span class="phase-label">Template</span>
                </div>
                <div class="phase-connector" data-after="template_selection"></div>
                <div class="phase-step" data-phase="research">
                    <div class="phase-icon" aria-hidden="true">2</div>
                    <span class="phase-label">Research</span>
                </div>
                <div class="phase-connector" data-after="research"></div>
                <div class="phase-step" data-phase="planning">
                    <div class="phase-icon" aria-hidden="true">3</div>
                    <span class="phase-label">Planning</span>
                </div>
                <div class="phase-connector" data-after="planning"></div>
                <div class="phase-step" data-phase="discussion">
                    <div class="phase-icon" aria-hidden="true">4</div>
                    <span class="phase-label">Discussion</span>
                </div>
                <div class="phase-connector" data-after="discussion"></div>
                <div class="phase-step" data-phase="development">
                    <div class="phase-icon" aria-hidden="true">5</div>
                    <span class="phase-label">Development</span>
                </div>
            </nav>
        </header>

        <div class="main-layout">
            <!-- Sidebar -->
            <aside class="sidebar" id="sidebar" role="navigation" aria-label="Task History">
                <div class="sidebar-header">
                    <div class="sidebar-title">
                        <h2>Projects</h2>
                    </div>
                    <select class="project-select" id="projectSelect" aria-label="Select project">
                        <option value="default">Default Project</option>
                    </select>
                    <div class="search-input">
                        <span aria-hidden="true">&#128269;</span>
                        <input type="text" id="taskSearch" placeholder="Search tasks..." aria-label="Search tasks">
                        <span class="search-shortcut" aria-hidden="true">/</span>
                    </div>
                </div>
                <div class="sidebar-content">
                    <button class="new-task-btn" id="newTaskBtn" aria-label="Create new task">
                        <span aria-hidden="true">+</span> New Task
                    </button>
                    <div class="task-list" id="tasksList" role="list" aria-label="Task history">
                        <div class="empty-state">
                            <p>No tasks yet</p>
                        </div>
                    </div>
                </div>
            </aside>

            <!-- Main Content -->
            <main class="main-content" id="main-content" role="main">
                <!-- History Banner -->
                <div class="history-banner" id="historyBanner" role="region" aria-label="Viewing historical task">
                    <div class="history-header">
                        <span>Viewing: <strong class="history-task-name" id="historyTaskName"></strong></span>
                        <button class="back-btn" id="backToLiveBtn">&#8592; Back to Live</button>
                    </div>
                    <form class="continue-form" id="continueForm">
                        <input type="text" id="continueInput" placeholder="Continue with: add feature, fix bug, improve..." autocomplete="off" aria-label="Continue task with instructions">
                        <button type="submit" id="continueBtn">Continue Task</button>
                    </form>
                </div>

                <!-- Development Phase Timeline -->
                <section class="dev-phases-section" id="devPhasesSection" style="display: none;" role="region" aria-label="Development phases progress">
                    <h2>Development Phases</h2>
                    <div class="dev-phases-timeline" id="devPhasesTimeline" role="list">
                        <!-- Phases will be inserted here by JavaScript -->
                    </div>
                </section>

                <!-- Progress Dashboard -->
                <section class="progress-dashboard" id="progressDashboard" style="display: none;" role="region" aria-label="Project progress overview">
                    <h2>Progress Overview</h2>
                    <div class="metrics-grid">
                        <!-- Total Progress -->
                        <div class="metric-card">
                            <div class="metric-header">
                                <span class="metric-icon" aria-hidden="true">📊</span>
                                <span class="metric-label">Total Progress</span>
                            </div>
                            <div class="metric-value info" id="totalProgress">0%</div>
                            <div class="metric-description" id="progressDescription">0 of 0 subtasks complete</div>
                        </div>

                        <!-- Current Phase -->
                        <div class="metric-card">
                            <div class="metric-header">
                                <span class="metric-icon" aria-hidden="true">⚡</span>
                                <span class="metric-label">Current Phase</span>
                            </div>
                            <div class="metric-value" id="currentPhaseName">—</div>
                            <div class="metric-description">
                                <span class="metric-badge" id="currentPhaseIteration" style="display: none;"></span>
                            </div>
                        </div>

                        <!-- Phases Complete -->
                        <div class="metric-card">
                            <div class="metric-header">
                                <span class="metric-icon" aria-hidden="true">✅</span>
                                <span class="metric-label">Phases Complete</span>
                            </div>
                            <div class="metric-value success" id="phasesComplete">0 / 0</div>
                            <div class="metric-description" id="phasesDescription">No phases yet</div>
                        </div>

                        <!-- QA Status -->
                        <div class="metric-card">
                            <div class="metric-header">
                                <span class="metric-icon" aria-hidden="true">🔍</span>
                                <span class="metric-label">QA Status</span>
                            </div>
                            <div class="metric-value" id="qaStatusSummary">—</div>
                            <div class="metric-description">
                                <span class="metric-badge" id="qaStatusBadge" style="display: none;"></span>
                            </div>
                        </div>
                    </div>
                </section>

                <!-- QA Feedback Panel -->
                <section class="qa-feedback-panel" id="qaFeedbackPanel" style="display: none;" role="region" aria-label="QA review feedback">
                    <div class="panel-header">
                        <h3 id="qaFeedbackTitle">QA Review</h3>
                        <span id="qaIterationBadge" class="iteration-badge" style="display: none;"></span>
                    </div>

                    <div class="qa-status-banner" id="qaStatusBanner" role="status">
                        <div class="status-icon" id="qaStatusIcon" aria-hidden="true"></div>
                        <div class="status-text" id="qaStatusText"></div>
                        <div class="reviewer-info" id="qaReviewerInfo"></div>
                    </div>

                    <div class="qa-feedback-content" id="qaFeedbackContent">
                        <!-- Feedback will be inserted here by JavaScript -->
                    </div>

                    <div class="feedback-actions" id="qaFeedbackActions" style="display: none;">
                        <button class="btn-secondary" onclick="showFullFeedback()">View Full Feedback</button>
                        <button class="btn-secondary" onclick="showIterationHistory()">View Iteration History</button>
                    </div>
                </section>

                <!-- Task Input -->
                <div class="task-input-area">
                    <form class="task-form" id="taskForm">
                        <div class="project-input-wrapper">
                            <input type="text" id="projectInput" placeholder="Project" autocomplete="off" aria-label="Project name">
                        </div>
                        <input type="text" id="taskInput" placeholder="Enter a task for the agent team..." autocomplete="off" aria-label="Task description">
                        <button type="submit" id="submitBtn">Start Task</button>
                    </form>
                </div>

                <!-- Delegation Breadcrumb -->
                <div class="delegation-breadcrumb" id="delegationBreadcrumb" role="navigation" aria-label="Current delegation chain">
                </div>

                <!-- Messages -->
                <div class="messages-area" id="messagesArea" role="log" aria-live="polite" aria-label="Agent messages">
                    <div class="empty-state" id="emptyState">
                        <div class="empty-state-icon" aria-hidden="true">&#128172;</div>
                        <div class="empty-state-title">No messages yet</div>
                        <p class="empty-state-description">Submit a task above to start collaborating with the agent team.</p>
                    </div>
                </div>
            </main>

            <!-- Context Panel -->
            <aside class="context-panel" id="contextPanel" role="complementary" aria-label="Context information">
                <div class="context-tabs" role="tablist">
                    <button class="context-tab active" data-tab="agents" role="tab" aria-selected="true" aria-controls="agentsPanel">Agents</button>
                    <button class="context-tab" data-tab="files" role="tab" aria-selected="false" aria-controls="filesPanel">Files</button>
                </div>
                <div class="context-content">
                    <!-- Agents Panel -->
                    <div class="context-section agents-section active" id="agentsPanel" role="tabpanel" aria-labelledby="Agents">
                        <h3>Agent Hierarchy</h3>
                        <div class="agent-org-chart" id="agentOrgChart">
                            <!-- Executive Tier -->
                            <div class="agent-tier">
                                <div class="agent-card" data-agent="ceo">
                                    <div class="agent-card-avatar" data-agent="ceo" aria-hidden="true">&#128084;</div>
                                    <span class="agent-card-name">CEO</span>
                                    <div class="agent-status-indicator"></div>
                                </div>
                            </div>
                            <div class="tier-connector"></div>
                            <!-- Management Tier -->
                            <div class="agent-tier">
                                <div class="agent-card" data-agent="pm">
                                    <div class="agent-card-avatar" data-agent="pm" aria-hidden="true">&#128203;</div>
                                    <span class="agent-card-name">PM</span>
                                    <div class="agent-status-indicator"></div>
                                </div>
                            </div>
                            <div class="tier-connector"></div>
                            <!-- Design Tier -->
                            <div class="agent-tier">
                                <div class="agent-card" data-agent="ux">
                                    <div class="agent-card-avatar" data-agent="ux" aria-hidden="true">&#127912;</div>
                                    <span class="agent-card-name">UX</span>
                                    <div class="agent-status-indicator"></div>
                                </div>
                                <div class="agent-card" data-agent="ui">
                                    <div class="agent-card-avatar" data-agent="ui" aria-hidden="true">&#128444;</div>
                                    <span class="agent-card-name">UI</span>
                                    <div class="agent-status-indicator"></div>
                                </div>
                            </div>
                            <div class="tier-connector"></div>
                            <!-- Architecture Tier -->
                            <div class="agent-tier">
                                <div class="agent-card" data-agent="architect">
                                    <div class="agent-card-avatar" data-agent="architect" aria-hidden="true">&#127959;</div>
                                    <span class="agent-card-name">Architect</span>
                                    <div class="agent-status-indicator"></div>
                                </div>
                                <div class="agent-card" data-agent="security">
                                    <div class="agent-card-avatar" data-agent="security" aria-hidden="true">&#128274;</div>
                                    <span class="agent-card-name">Security</span>
                                    <div class="agent-status-indicator"></div>
                                </div>
                            </div>
                            <div class="tier-connector"></div>
                            <!-- Development Tier -->
                            <div class="agent-tier">
                                <div class="agent-card" data-agent="senior_dev">
                                    <div class="agent-card-avatar" data-agent="senior_dev" aria-hidden="true">&#128104;&#8205;&#128187;</div>
                                    <span class="agent-card-name">Senior Dev</span>
                                    <div class="agent-status-indicator"></div>
                                </div>
                                <div class="agent-card" data-agent="junior_dev">
                                    <div class="agent-card-avatar" data-agent="junior_dev" aria-hidden="true">&#128105;&#8205;&#128187;</div>
                                    <span class="agent-card-name">Junior Dev</span>
                                    <div class="agent-status-indicator"></div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Files Panel -->
                    <div class="context-section files-section" id="filesPanel" role="tabpanel" aria-labelledby="Files">
                        <h3>Project Files</h3>
                        <div class="file-tree" id="fileTree">
                            <div class="empty-state">
                                <p>No files created yet</p>
                            </div>
                        </div>
                    </div>
                </div>
            </aside>
        </div>
    </div>

    <!-- Mobile Sidebar Toggle -->
    <button class="mobile-sidebar-toggle" id="mobileSidebarToggle" aria-label="Toggle sidebar" aria-expanded="false">
        &#9776;
    </button>

    <!-- Mobile Overlay -->
    <div class="mobile-overlay" id="mobileOverlay"></div>

    <!-- Toast Container -->
    <div class="toast-container" id="toastContainer" role="alert" aria-live="assertive"></div>

    <!-- Agent Detail Modal -->
    <div class="modal" id="agentDetailModal">
        <div class="modal-content">
            <div class="modal-header">
                <h2 id="agentDetailTitle">Agent Details</h2>
                <button class="modal-close" onclick="closeAgentDetailModal()">&times;</button>
            </div>
            <div class="modal-body">
                <div class="agent-status-section">
                    <h3>Status</h3>
                    <div class="agent-status-card" id="agentStatusCard">
                        <div class="status-row">
                            <span>State:</span>
                            <span id="agentState" class="status-badge">idle</span>
                        </div>
                        <div class="status-row">
                            <span>Tasks Completed:</span>
                            <span id="agentTasksCompleted">0</span>
                        </div>
                        <div class="status-row">
                            <span>Tasks Failed:</span>
                            <span id="agentTasksFailed">0</span>
                        </div>
                    </div>
                </div>
                <div class="agent-tasks-section">
                    <h3>Current Tasks</h3>
                    <div id="agentTasksList" class="tasks-list"></div>
                </div>
                <div class="agent-outputs-section">
                    <h3>Outputs</h3>
                    <div id="agentOutputsList" class="outputs-list"></div>
                </div>
            </div>
        </div>
    </div>

    <!-- Output Viewer Modal -->
    <div class="modal" id="outputViewerModal">
        <div class="modal-content">
            <div class="modal-header">
                <h2 id="outputViewerTitle">Output</h2>
                <button class="modal-close" onclick="closeOutputViewerModal()">&times;</button>
            </div>
            <div class="modal-body">
                <div id="outputViewerContent" class="output-content"></div>
            </div>
        </div>
    </div>

    <style>
        /* Modal Styles */
        .modal {
            display: none;
            position: fixed;
            z-index: 9999;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0, 0, 0, 0.7);
            animation: fadeIn 0.2s ease;
        }

        .modal.active {
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .modal-content {
            background: var(--color-bg-secondary);
            border-radius: var(--radius-lg);
            max-width: 800px;
            max-height: 90vh;
            width: 90%;
            overflow-y: auto;
            box-shadow: var(--shadow-lg);
        }

        .modal-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: var(--space-4);
            border-bottom: 1px solid var(--color-border);
        }

        .modal-header h2 {
            margin: 0;
            color: var(--color-text-primary);
        }

        .modal-close {
            background: none;
            border: none;
            font-size: 2rem;
            color: var(--color-text-secondary);
            cursor: pointer;
            padding: 0;
            width: 30px;
            height: 30px;
            display: flex;
            align-items: center;
            justify-content: center;
            transition: var(--transition-fast);
        }

        .modal-close:hover {
            color: var(--color-text-primary);
        }

        .modal-body {
            padding: var(--space-4);
        }

        .agent-status-card {
            background: var(--color-bg-tertiary);
            padding: var(--space-4);
            border-radius: var(--radius-md);
            margin-bottom: var(--space-4);
        }

        .status-row {
            display: flex;
            justify-content: space-between;
            padding: var(--space-2) 0;
            border-bottom: 1px solid var(--color-border);
        }

        .status-row:last-child {
            border-bottom: none;
        }

        .status-badge {
            padding: var(--space-1) var(--space-3);
            border-radius: var(--radius-sm);
            font-size: var(--text-sm);
            font-weight: 600;
        }

        .status-badge.working {
            background: var(--color-status-info);
            color: white;
        }

        .status-badge.idle {
            background: var(--color-text-tertiary);
            color: white;
        }

        .status-badge.completed {
            background: var(--color-status-success);
            color: white;
        }

        .status-badge.error {
            background: var(--color-status-error);
            color: white;
        }

        .tasks-list, .outputs-list {
            display: flex;
            flex-direction: column;
            gap: var(--space-2);
        }

        .task-item {
            background: var(--color-bg-tertiary);
            padding: var(--space-3);
            border-radius: var(--radius-md);
            border-left: 3px solid var(--color-accent);
        }

        .task-item-header {
            display: flex;
            justify-content: space-between;
            margin-bottom: var(--space-2);
        }

        .task-title {
            font-weight: 600;
            color: var(--color-text-primary);
        }

        .task-progress {
            font-size: var(--text-sm);
            color: var(--color-text-secondary);
        }

        .progress-bar {
            height: 4px;
            background: var(--color-bg-primary);
            border-radius: 2px;
            overflow: hidden;
            margin-top: var(--space-2);
        }

        .progress-bar-fill {
            height: 100%;
            background: var(--color-accent);
            transition: width 0.3s ease;
        }

        .output-item {
            background: var(--color-bg-tertiary);
            padding: var(--space-3);
            border-radius: var(--radius-md);
            cursor: pointer;
            transition: var(--transition-fast);
        }

        .output-item:hover {
            background: var(--color-bg-elevated);
        }

        .output-title {
            font-weight: 600;
            color: var(--color-text-primary);
            margin-bottom: var(--space-1);
        }

        .output-meta {
            font-size: var(--text-sm);
            color: var(--color-text-secondary);
        }

        .output-content {
            background: var(--color-bg-tertiary);
            padding: var(--space-4);
            border-radius: var(--radius-md);
            white-space: pre-wrap;
            font-family: monospace;
            color: var(--color-text-primary);
            max-height: 60vh;
            overflow-y: auto;
        }

        /* Markdown content styling */
        .markdown-content {
            line-height: 1.6;
            color: var(--color-text-primary);
        }

        .markdown-content h1,
        .markdown-content h2,
        .markdown-content h3,
        .markdown-content h4,
        .markdown-content h5,
        .markdown-content h6 {
            margin-top: var(--space-6);
            margin-bottom: var(--space-3);
            color: var(--color-text-primary);
            font-weight: 600;
        }

        .markdown-content h1 {
            font-size: 2em;
            border-bottom: 2px solid var(--color-border);
            padding-bottom: var(--space-2);
        }

        .markdown-content h2 {
            font-size: 1.5em;
            border-bottom: 1px solid var(--color-border);
            padding-bottom: var(--space-2);
        }

        .markdown-content h3 {
            font-size: 1.25em;
        }

        .markdown-content p {
            margin-bottom: var(--space-4);
        }

        .markdown-content ul,
        .markdown-content ol {
            margin-bottom: var(--space-4);
            padding-left: var(--space-6);
        }

        .markdown-content li {
            margin-bottom: var(--space-2);
        }

        .markdown-content code {
            background: var(--color-bg-elevated);
            padding: 2px 6px;
            border-radius: var(--radius-sm);
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 0.9em;
            color: var(--color-status-info);
        }

        .markdown-content pre {
            background: var(--color-bg-elevated);
            padding: var(--space-4);
            border-radius: var(--radius-md);
            overflow-x: auto;
            margin-bottom: var(--space-4);
            border-left: 3px solid var(--color-accent);
        }

        .markdown-content pre code {
            background: none;
            padding: 0;
            color: var(--color-text-primary);
            font-size: var(--text-sm);
        }

        .markdown-content blockquote {
            border-left: 4px solid var(--color-accent);
            padding-left: var(--space-4);
            margin: var(--space-4) 0;
            color: var(--color-text-secondary);
            font-style: italic;
        }

        .markdown-content table {
            border-collapse: collapse;
            width: 100%;
            margin-bottom: var(--space-4);
        }

        .markdown-content th,
        .markdown-content td {
            border: 1px solid var(--color-border);
            padding: var(--space-2) var(--space-3);
            text-align: left;
        }

        .markdown-content th {
            background: var(--color-bg-elevated);
            font-weight: 600;
        }

        .markdown-content a {
            color: var(--color-accent);
            text-decoration: none;
        }

        .markdown-content a:hover {
            text-decoration: underline;
        }

        .markdown-content img {
            max-width: 100%;
            height: auto;
            border-radius: var(--radius-md);
            margin: var(--space-4) 0;
        }

        .markdown-content hr {
            border: none;
            border-top: 1px solid var(--color-border);
            margin: var(--space-6) 0;
        }

        /* Agent card click */
        .agent-card {
            cursor: pointer;
            transition: var(--transition-fast);
        }

        .agent-card:hover {
            transform: scale(1.05);
        }

        /* Pulse animation for working state */
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }

        .agent-status-indicator.working {
            animation: pulse 2s infinite;
        }

        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }
    </style>

    <!-- Marked.js for markdown rendering -->
    <script src="https://cdn.jsdelivr.net/npm/marked@11.1.1/marked.min.js"></script>

    <script>
        // ===================================
        // STATE MANAGEMENT
        // ===================================
        const state = {
            ws: null,
            messages: [],
            tasks: [],
            files: [],
            isTaskRunning: false,
            viewingHistory: false,
            selectedTaskId: null,
            currentTaskId: null,
            currentPhase: null,
            delegationChain: [],
            activeAgent: null,
            projectId: 'default',
            currentPlan: null // Development plan data
        };

        // ===================================
        // AGENT CONFIGURATION
        // ===================================
        const agentConfig = {
            ceo: { emoji: '&#128084;', name: 'CEO', color: '#4A90A4' },
            pm: { emoji: '&#128203;', name: 'PM', color: '#7B68EE' },
            ux: { emoji: '&#127912;', name: 'UX', color: '#FF6B6B' },
            ui: { emoji: '&#128444;', name: 'UI', color: '#4ECDC4' },
            security: { emoji: '&#128274;', name: 'Security', color: '#F39C12' },
            architect: { emoji: '&#127959;', name: 'Architect', color: '#9B59B6' },
            senior_dev: { emoji: '&#128104;&#8205;&#128187;', name: 'Senior Dev', color: '#2ECC71' },
            junior_dev: { emoji: '&#128105;&#8205;&#128187;', name: 'Junior Dev', color: '#3498DB' }
        };

        const phaseConfig = {
            template_selection: { name: 'Template', icon: '1' },
            research: { name: 'Research', icon: '2' },
            planning: { name: 'Planning', icon: '3' },
            discussion: { name: 'Discussion', icon: '4' },
            development: { name: 'Development', icon: '5' }
        };

        // ===================================
        // DOM ELEMENTS
        // ===================================
        const elements = {
            wsStatus: document.getElementById('wsStatus'),
            wsStatusText: document.getElementById('wsStatusText'),
            taskStatus: document.getElementById('taskStatus'),
            taskStatusText: document.getElementById('taskStatusText'),
            messageCount: document.getElementById('messageCount'),
            messagesArea: document.getElementById('messagesArea'),
            emptyState: document.getElementById('emptyState'),
            tasksList: document.getElementById('tasksList'),
            fileTree: document.getElementById('fileTree'),
            taskForm: document.getElementById('taskForm'),
            projectInput: document.getElementById('projectInput'),
            taskInput: document.getElementById('taskInput'),
            submitBtn: document.getElementById('submitBtn'),
            newTaskBtn: document.getElementById('newTaskBtn'),
            historyBanner: document.getElementById('historyBanner'),
            historyTaskName: document.getElementById('historyTaskName'),
            backToLiveBtn: document.getElementById('backToLiveBtn'),
            continueForm: document.getElementById('continueForm'),
            continueInput: document.getElementById('continueInput'),
            phaseStepper: document.getElementById('phaseStepper'),
            delegationBreadcrumb: document.getElementById('delegationBreadcrumb'),
            toastContainer: document.getElementById('toastContainer'),
            announcer: document.getElementById('announcer'),
            projectSelect: document.getElementById('projectSelect'),
            taskSearch: document.getElementById('taskSearch'),
            sidebar: document.getElementById('sidebar'),
            mobileSidebarToggle: document.getElementById('mobileSidebarToggle'),
            mobileOverlay: document.getElementById('mobileOverlay'),
            contextTabs: document.querySelectorAll('.context-tab'),
            contextSections: document.querySelectorAll('.context-section'),
            agentOrgChart: document.getElementById('agentOrgChart'),
            devPhasesSection: document.getElementById('devPhasesSection'),
            devPhasesTimeline: document.getElementById('devPhasesTimeline'),
            qaFeedbackPanel: document.getElementById('qaFeedbackPanel'),
            qaFeedbackTitle: document.getElementById('qaFeedbackTitle'),
            qaIterationBadge: document.getElementById('qaIterationBadge'),
            qaStatusBanner: document.getElementById('qaStatusBanner'),
            qaStatusIcon: document.getElementById('qaStatusIcon'),
            qaStatusText: document.getElementById('qaStatusText'),
            qaReviewerInfo: document.getElementById('qaReviewerInfo'),
            qaFeedbackContent: document.getElementById('qaFeedbackContent'),
            qaFeedbackActions: document.getElementById('qaFeedbackActions'),
            progressDashboard: document.getElementById('progressDashboard'),
            totalProgress: document.getElementById('totalProgress'),
            progressDescription: document.getElementById('progressDescription'),
            currentPhaseName: document.getElementById('currentPhaseName'),
            currentPhaseIteration: document.getElementById('currentPhaseIteration'),
            phasesComplete: document.getElementById('phasesComplete'),
            phasesDescription: document.getElementById('phasesDescription'),
            qaStatusSummary: document.getElementById('qaStatusSummary'),
            qaStatusBadge: document.getElementById('qaStatusBadge')
        };

        // ===================================
        // DEVELOPMENT PHASE TIMELINE
        // ===================================
        async function updateDevelopmentPhases(projectId) {
            if (!projectId || !elements.devPhasesSection) return;

            try {
                const response = await fetch('/api/phase-status?project=' + projectId);
                const status = await response.json();

                if (!status.has_plan) {
                    elements.devPhasesSection.style.display = 'none';
                    return;
                }

                // Show section
                elements.devPhasesSection.style.display = 'block';

                // Fetch full plan
                const planResponse = await fetch('/api/development-plan?project=' + projectId);
                state.currentPlan = await planResponse.json();

                renderDevelopmentPhases(state.currentPlan, status);
            } catch (error) {
                console.error('Failed to fetch development phases:', error);
                elements.devPhasesSection.style.display = 'none';
            }
        }

        function renderDevelopmentPhases(plan, status) {
            if (!elements.devPhasesTimeline || !plan || !plan.phases) return;

            elements.devPhasesTimeline.innerHTML = '';

            plan.phases.forEach((phase, index) => {
                // Create phase card
                const card = document.createElement('div');
                card.className = 'phase-card ' + getPhaseStatusClass(phase);
                card.dataset.phase = phase.index;
                card.setAttribute('role', 'listitem');

                const completedSubtasks = phase.subtasks.filter(st => st.status === 'completed').length;
                const totalSubtasks = phase.subtasks.length;
                const progress = totalSubtasks > 0 ? (completedSubtasks / totalSubtasks) * 100 : 0;

                var iterationHTML = '';
                if (phase.iteration > 1) {
                    var badgeClass = 'iteration-badge' + (phase.iteration >= 3 ? ' warning' : '');
                    iterationHTML = '<span class="' + badgeClass + '" aria-label="Iteration ' + phase.iteration + ' of 3">Iteration ' + phase.iteration + '/3</span>';
                }

                card.innerHTML =
                    '<div class="phase-header">' +
                        '<div class="phase-badge" aria-label="Phase ' + phase.index + '">' + getPhaseIcon(phase) + '</div>' +
                        iterationHTML +
                    '</div>' +
                    '<div class="phase-name">' + escapeHtml(phase.name) + '</div>' +
                    '<div class="phase-progress" role="progressbar" aria-valuenow="' + Math.round(progress) + '" aria-valuemin="0" aria-valuemax="100">' +
                        '<div class="phase-progress-fill" style="width: ' + progress + '%"></div>' +
                    '</div>' +
                    '<div class="phase-meta">' +
                        '<span class="qa-status ' + getQAStatusClass(phase.qa_status) + '">' + getQAStatusText(phase.qa_status) + '</span>' +
                        '<span aria-label="' + completedSubtasks + ' of ' + totalSubtasks + ' subtasks complete">' + completedSubtasks + '/' + totalSubtasks + ' subtasks</span>' +
                    '</div>';

                card.addEventListener('click', () => showPhaseDetails(phase));
                card.setAttribute('tabindex', '0');
                card.addEventListener('keypress', (e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault();
                        showPhaseDetails(phase);
                    }
                });

                elements.devPhasesTimeline.appendChild(card);

                // Add connector if not last
                if (index < plan.phases.length - 1) {
                    const connector = document.createElement('div');
                    connector.className = 'timeline-connector';
                    connector.setAttribute('aria-hidden', 'true');
                    elements.devPhasesTimeline.appendChild(connector);
                }
            });
        }

        function getPhaseStatusClass(phase) {
            if (phase.qa_status === 'approved' && phase.status === 'completed') {
                return 'completed';
            } else if (phase.status === 'needs_revision' || phase.qa_status === 'rejected') {
                return 'rejected';
            } else if (phase.status === 'in_progress') {
                return 'in-progress';
            } else {
                return 'pending';
            }
        }

        function getPhaseIcon(phase) {
            if (phase.qa_status === 'approved' && phase.status === 'completed') {
                return '✓';
            } else if (phase.status === 'needs_revision') {
                return '⚠';
            } else {
                return phase.index;
            }
        }

        function getQAStatusClass(status) {
            return status || 'pending';
        }

        function getQAStatusText(status) {
            const icons = {
                approved: '✅ Approved',
                rejected: '❌ Rejected',
                pending: '⏳ Pending'
            };
            return icons[status] || '⏳ Pending';
        }

        function showPhaseDetails(phase) {
            // TODO: Show modal or sidebar with phase details, subtasks, QA feedback
            console.log('Show details for phase:', phase);
            const detailsMsg = 'Phase ' + phase.index + ': ' + phase.name + '\n\nSubtasks: ' + phase.subtasks.length + '\nStatus: ' + phase.status + '\nQA: ' + phase.qa_status;
            announce(detailsMsg);
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        // ===================================
        // QA FEEDBACK PANEL
        // ===================================
        async function updateQAFeedback(projectId) {
            if (!projectId || !elements.qaFeedbackPanel) return;

            try {
                const response = await fetch('/api/qa-reviews?project=' + projectId);
                const history = await response.json();

                if (!history.reviews || history.reviews.length === 0) {
                    elements.qaFeedbackPanel.style.display = 'none';
                    return;
                }

                // Get latest review
                const latest = history.reviews[history.reviews.length - 1];

                // Show panel
                elements.qaFeedbackPanel.style.display = 'block';

                // Update title
                elements.qaFeedbackTitle.textContent = 'QA Review - Phase ' + latest.phase_index + ': ' + latest.phase_name;

                // Update iteration badge
                if (latest.iteration > 1) {
                    elements.qaIterationBadge.textContent = 'Iteration ' + latest.iteration + '/3';
                    elements.qaIterationBadge.className = 'iteration-badge' + (latest.iteration >= 3 ? ' warning' : '');
                    elements.qaIterationBadge.style.display = 'inline-block';
                } else {
                    elements.qaIterationBadge.style.display = 'none';
                }

                // Update status banner
                if (latest.status === 'approved') {
                    elements.qaStatusBanner.className = 'qa-status-banner approved';
                    elements.qaStatusIcon.textContent = '✅';
                    elements.qaStatusText.textContent = 'QA APPROVED';
                } else if (latest.status === 'rejected') {
                    elements.qaStatusBanner.className = 'qa-status-banner rejected';
                    elements.qaStatusIcon.textContent = '❌';
                    elements.qaStatusText.textContent = 'QA REJECTED';
                } else {
                    elements.qaStatusBanner.className = 'qa-status-banner';
                    elements.qaStatusIcon.textContent = '⏳';
                    elements.qaStatusText.textContent = 'QA PENDING';
                }

                // Update reviewer info
                const reviewTime = timeAgo(latest.reviewed_at);
                elements.qaReviewerInfo.textContent = 'Reviewer: ' + latest.reviewer.toUpperCase() + ' | ' + reviewTime;

                // Update feedback content
                elements.qaFeedbackContent.innerHTML = formatQAFeedback(latest);

                // Show action buttons if there's detailed feedback
                if (latest.feedback || (latest.failed_criteria && latest.failed_criteria.length > 0)) {
                    elements.qaFeedbackActions.style.display = 'flex';
                }

            } catch (error) {
                console.error('Failed to fetch QA feedback:', error);
                elements.qaFeedbackPanel.style.display = 'none';
            }
        }

        function formatQAFeedback(review) {
            var html = '';

            // Failed criteria section
            if (review.failed_criteria && review.failed_criteria.length > 0) {
                html += '<div class="failed-criteria">';
                html += '<h4>Failed Criteria:</h4>';

                review.failed_criteria.forEach(function(criterion) {
                    // Try to parse criterion as "title: description"
                    var parts = criterion.split(':');
                    var title = parts[0] || criterion;
                    var description = parts.slice(1).join(':').trim();

                    html += '<div class="criterion-item">';
                    html += '<span class="icon">❌</span>';
                    html += '<div class="content">';
                    html += '<strong>' + escapeHtml(title) + '</strong>';
                    if (description) {
                        html += '<p>' + escapeHtml(description) + '</p>';
                    }

                    // Try to extract file references from description
                    var fileMatch = description.match(/(\w+\.(js|css|html|ts|tsx|go|md))(:(\d+))?/);
                    if (fileMatch) {
                        html += '<a href="#" class="file-link" onclick="openFile(\'' + fileMatch[1] + '\', ' + (fileMatch[4] || 'null') + '); return false;">' + fileMatch[0] + '</a>';
                    }

                    html += '</div>';
                    html += '</div>';
                });

                html += '</div>';
            }

            // Required fixes section (parse from feedback text)
            if (review.feedback) {
                var fixes = extractRequiredFixes(review.feedback);
                if (fixes.length > 0) {
                    html += '<div class="required-fixes">';
                    html += '<h4>Required Fixes:</h4>';
                    html += '<ol>';
                    fixes.forEach(function(fix) {
                        html += '<li>' + escapeHtml(fix) + '</li>';
                    });
                    html += '</ol>';
                    html += '</div>';
                }

                // Show raw feedback if no structured data
                if (!review.failed_criteria || review.failed_criteria.length === 0) {
                    html += '<div style="margin-top: 16px;">';
                    html += '<h4>Feedback:</h4>';
                    html += '<pre style="white-space: pre-wrap; font-family: inherit; color: var(--color-text-secondary);">' + escapeHtml(review.feedback) + '</pre>';
                    html += '</div>';
                }
            }

            return html || '<p style="color: var(--color-text-secondary);">No detailed feedback available.</p>';
        }

        function extractRequiredFixes(feedback) {
            var fixes = [];
            if (!feedback) return fixes;

            // Look for numbered list items (1. or 1) format)
            var lines = feedback.split('\n');
            var inFixesSection = false;

            lines.forEach(function(line) {
                line = line.trim();

                // Detect "REQUIRED FIXES" section
                if (line.match(/REQUIRED FIXES|FIXES NEEDED|TO FIX/i)) {
                    inFixesSection = true;
                    return;
                }

                // Stop at next section
                if (inFixesSection && line.match(/^[A-Z\s]{3,}:/)) {
                    inFixesSection = false;
                }

                // Extract numbered items
                if (inFixesSection) {
                    var match = line.match(/^(\d+[\.\)])\s*(.+)/);
                    if (match) {
                        fixes.push(match[2]);
                    }
                }
            });

            return fixes;
        }

        function timeAgo(timestamp) {
            if (!timestamp) return 'recently';

            var date = new Date(timestamp);
            var seconds = Math.floor((new Date() - date) / 1000);

            if (seconds < 60) return seconds + 's ago';
            if (seconds < 3600) return Math.floor(seconds / 60) + 'm ago';
            if (seconds < 86400) return Math.floor(seconds / 3600) + 'h ago';
            return Math.floor(seconds / 86400) + 'd ago';
        }

        function openFile(filename, lineNumber) {
            // Placeholder for file viewer integration
            console.log('Open file:', filename, 'line:', lineNumber);
            announce('Opening file: ' + filename + (lineNumber ? ' line ' + lineNumber : ''));
            // TODO: Integrate with file tree or code viewer
        }

        function showFullFeedback() {
            // Placeholder for modal with complete feedback
            console.log('Show full feedback modal');
            // TODO: Implement modal with full QA review details
        }

        function showIterationHistory() {
            // Placeholder for iteration history view
            console.log('Show iteration history');
            // TODO: Implement iteration timeline modal
        }

        // ===================================
        // PROGRESS DASHBOARD
        // ===================================
        async function updateProgressDashboard(projectId) {
            if (!projectId || !elements.progressDashboard) return;

            try {
                // Fetch phase status
                const statusResponse = await fetch('/api/phase-status?project=' + projectId);
                const status = await statusResponse.json();

                if (!status.has_plan) {
                    elements.progressDashboard.style.display = 'none';
                    return;
                }

                // Show dashboard
                elements.progressDashboard.style.display = 'block';

                // Fetch full plan
                const planResponse = await fetch('/api/development-plan?project=' + projectId);
                const plan = await planResponse.json();

                // Calculate metrics
                var totalSubtasks = 0;
                var completedSubtasks = 0;
                var phasesCompleted = 0;
                var currentPhase = null;
                var latestQAStatus = '—';

                if (plan && plan.phases) {
                    plan.phases.forEach(function(phase) {
                        // Count subtasks
                        if (phase.sub_tasks) {
                            totalSubtasks = totalSubtasks + phase.sub_tasks.length;
                            phase.sub_tasks.forEach(function(subtask) {
                                if (subtask.status === 'completed') {
                                    completedSubtasks = completedSubtasks + 1;
                                }
                            });
                        }

                        // Count completed phases
                        if (phase.status === 'completed' && phase.qa_status === 'approved') {
                            phasesCompleted = phasesCompleted + 1;
                        }

                        // Find current phase
                        if (phase.status === 'in_progress' || phase.status === 'needs_revision') {
                            currentPhase = phase;
                        }

                        // Get latest QA status
                        if (phase.qa_status && phase.qa_status !== 'pending') {
                            latestQAStatus = phase.qa_status;
                        }
                    });

                    // Update Total Progress
                    var progressPercent = totalSubtasks > 0 ? Math.round((completedSubtasks / totalSubtasks) * 100) : 0;
                    elements.totalProgress.textContent = progressPercent + '%';
                    elements.progressDescription.textContent = completedSubtasks + ' of ' + totalSubtasks + ' subtasks complete';

                    // Color code based on progress
                    elements.totalProgress.className = 'metric-value';
                    if (progressPercent === 100) {
                        elements.totalProgress.classList.add('success');
                    } else if (progressPercent >= 50) {
                        elements.totalProgress.classList.add('info');
                    } else if (progressPercent > 0) {
                        elements.totalProgress.classList.add('warning');
                    }

                    // Update Current Phase
                    if (currentPhase) {
                        elements.currentPhaseName.textContent = currentPhase.name || 'Phase ' + currentPhase.index;
                        elements.currentPhaseName.className = 'metric-value';

                        // Show iteration badge if iteration > 1
                        if (currentPhase.iteration && currentPhase.iteration > 1) {
                            elements.currentPhaseIteration.textContent = 'Iteration ' + currentPhase.iteration + '/3';
                            elements.currentPhaseIteration.style.display = 'inline-flex';
                            elements.currentPhaseIteration.className = 'metric-badge';
                            if (currentPhase.iteration >= 3) {
                                elements.currentPhaseIteration.classList.add('error');
                            } else {
                                elements.currentPhaseIteration.classList.add('warning');
                            }
                        } else {
                            elements.currentPhaseIteration.style.display = 'none';
                        }

                        // Color code based on status
                        if (currentPhase.status === 'needs_revision') {
                            elements.currentPhaseName.classList.add('warning');
                        } else {
                            elements.currentPhaseName.classList.add('info');
                        }
                    } else {
                        elements.currentPhaseName.textContent = '—';
                        elements.currentPhaseName.className = 'metric-value';
                        elements.currentPhaseIteration.style.display = 'none';
                    }

                    // Update Phases Complete
                    var totalPhases = plan.phases.length;
                    elements.phasesComplete.textContent = phasesCompleted + ' / ' + totalPhases;
                    elements.phasesDescription.textContent = (totalPhases - phasesCompleted) + ' phases remaining';

                    // Color code phases complete
                    elements.phasesComplete.className = 'metric-value';
                    if (phasesCompleted === totalPhases) {
                        elements.phasesComplete.classList.add('success');
                    } else if (phasesCompleted > 0) {
                        elements.phasesComplete.classList.add('info');
                    }

                    // Update QA Status
                    if (latestQAStatus === 'approved') {
                        elements.qaStatusSummary.textContent = 'Approved';
                        elements.qaStatusSummary.className = 'metric-value success';
                        elements.qaStatusBadge.textContent = '✓ All checks passed';
                        elements.qaStatusBadge.className = 'metric-badge success';
                        elements.qaStatusBadge.style.display = 'inline-flex';
                    } else if (latestQAStatus === 'rejected') {
                        elements.qaStatusSummary.textContent = 'Rejected';
                        elements.qaStatusSummary.className = 'metric-value error';
                        elements.qaStatusBadge.textContent = '✗ Revisions needed';
                        elements.qaStatusBadge.className = 'metric-badge error';
                        elements.qaStatusBadge.style.display = 'inline-flex';
                    } else {
                        elements.qaStatusSummary.textContent = '—';
                        elements.qaStatusSummary.className = 'metric-value';
                        elements.qaStatusBadge.style.display = 'none';
                    }
                }
            } catch (error) {
                console.error('Error updating progress dashboard:', error);
                elements.progressDashboard.style.display = 'none';
            }
        }

        // ===================================
        // WEBSOCKET CONNECTION
        // ===================================
        function connectWebSocket() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            state.ws = new WebSocket(protocol + '//' + window.location.host + '/ws');

            state.ws.onopen = () => {
                elements.wsStatus.classList.add('connected');
                elements.wsStatusText.textContent = 'Connected';
                announce('WebSocket connected');
            };

            state.ws.onclose = () => {
                elements.wsStatus.classList.remove('connected');
                elements.wsStatusText.textContent = 'Disconnected';
                setTimeout(connectWebSocket, 3000);
            };

            state.ws.onerror = () => {
                elements.wsStatus.classList.remove('connected');
                elements.wsStatusText.textContent = 'Error';
            };

            state.ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    if (data.type === 'message') {
                        handleIncomingMessage(data.message);
                    } else if (data.type === 'agent_task_event') {
                        handleAgentTaskEvent(data.event);
                    }
                } catch (e) {
                    console.error('Failed to parse WebSocket message:', e);
                }
            };
        }

        function handleIncomingMessage(msg) {
            // Update delegation chain
            if (msg.from && msg.from !== 'orchestrator' && msg.from !== 'system') {
                updateDelegationChain(msg.from, msg.to);
                updateActiveAgent(msg.from);
            }

            // Handle phase change
            if (msg.content && msg.content.includes('Entering') && msg.content.includes('phase')) {
                const phaseMatch = msg.content.match(/Entering (.*) phase/);
                if (phaseMatch) {
                    const phaseName = phaseMatch[1].toLowerCase().replace(' ', '_');
                    updatePhase(phaseName);
                }
            }

            // Handle task events
            if (msg.metadata && msg.metadata.tags) {
                if (msg.metadata.tags.includes('task_started')) {
                    state.currentTaskId = msg.metadata.task_id;
                    state.delegationChain = [];
                    fetchTasks();
                    fetchStatus();
                    fetchPhase();
                    showToast('Task Started', 'Agent team is working on your task', 'info');
                }
                if (msg.metadata.tags.includes('task_completed')) {
                    state.currentTaskId = null;
                    state.activeAgent = null;
                    state.delegationChain = [];
                    resetAgentStates();
                    resetPhase();
                    fetchTasks();
                    fetchStatus();
                    fetchFiles();
                    showToast('Task Completed', 'The agent team has finished the task', 'success');
                }
            }

            // Add message if in live view
            if (!state.viewingHistory) {
                addMessage(msg);
            }
        }

        // ===================================
        // MESSAGE HANDLING
        // ===================================
        function addMessage(msg) {
            state.messages.push(msg);
            renderMessages();
            updateMessageCount();

            // Scroll to bottom
            elements.messagesArea.scrollTop = elements.messagesArea.scrollHeight;
        }

        function renderMessages() {
            if (state.messages.length === 0) {
                elements.emptyState.style.display = 'block';
                return;
            }
            elements.emptyState.style.display = 'none';

            // Group messages by time
            const groups = groupMessagesByTime(state.messages);

            let html = '';
            for (const [groupName, msgs] of Object.entries(groups)) {
                html += '<div class="message-group">';
                html += '<div class="message-group-header">' + groupName + '</div>';
                html += msgs.map(renderMessage).join('');
                html += '</div>';
            }

            elements.messagesArea.innerHTML = html;
        }

        function groupMessagesByTime(messages) {
            const groups = {};
            const now = new Date();
            const today = now.toDateString();
            const yesterday = new Date(now - 86400000).toDateString();

            messages.forEach(msg => {
                const msgDate = new Date(msg.timestamp);
                const dateStr = msgDate.toDateString();
                let groupName;

                if (dateStr === today) {
                    groupName = 'Today';
                } else if (dateStr === yesterday) {
                    groupName = 'Yesterday';
                } else {
                    groupName = msgDate.toLocaleDateString('en-US', { weekday: 'long', month: 'short', day: 'numeric' });
                }

                if (!groups[groupName]) {
                    groups[groupName] = [];
                }
                groups[groupName].push(msg);
            });

            return groups;
        }

        function renderMessage(msg) {
            const time = new Date(msg.timestamp).toLocaleTimeString();
            const agent = agentConfig[msg.from] || { emoji: '&#129302;', name: msg.from, color: '#888' };
            const messageType = msg.type || 'response';
            const isFile = messageType === 'file_create';

            // Render markdown for message content (except file paths)
            let content;
            if (isFile) {
                content = escapeHtml(msg.content);
            } else if (typeof marked !== 'undefined') {
                // Configure marked
                marked.setOptions({
                    breaks: true,
                    gfm: true,
                    headerIds: false,
                    mangle: false
                });
                // Render markdown to HTML
                content = marked.parse(msg.content || '');
            } else {
                // Fallback to escaped HTML if marked is not available
                content = escapeHtml(msg.content);
            }

            const isLong = content.length > 500;

            return '<div class="message ' + (isFile ? 'file-message' : '') + '" data-agent="' + msg.from + '">' +
                '<div class="message-header">' +
                '<div class="message-sender">' +
                '<div class="agent-avatar" data-agent="' + msg.from + '">' + agent.emoji + '</div>' +
                '<div class="sender-info">' +
                '<span class="sender-name">' + agent.name + '</span>' +
                '<span class="sender-flow">' + (msg.to ? '&#8594; ' + msg.to : '') + '</span>' +
                '</div>' +
                '</div>' +
                '<div class="message-meta">' +
                '<span class="message-type-badge ' + messageType + '">' + formatMessageType(messageType) + '</span>' +
                '<span class="message-time">' + time + '</span>' +
                '</div>' +
                '</div>' +
                '<div class="message-content markdown-content' + (isLong ? ' collapsed' : '') + '">' +
                (isFile ? '<span class="file-path">' + content + '</span>' : content) +
                '</div>' +
                (isLong ? '<button class="expand-btn" onclick="toggleExpand(this)">Show more</button>' : '') +
                '</div>';
        }

        function formatMessageType(type) {
            const types = {
                'response': 'Response',
                'delegate': 'Delegate',
                'file_create': 'File',
                'error': 'Error',
                'system': 'System',
                'task': 'Task',
                'complete': 'Complete'
            };
            return types[type] || type;
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
            const count = state.messages.length;
            elements.messageCount.textContent = count + ' message' + (count !== 1 ? 's' : '');
        }

        // ===================================
        // PHASE MANAGEMENT
        // ===================================
        function updatePhase(phase) {
            state.currentPhase = phase;
            const phases = ['template_selection', 'research', 'planning', 'discussion', 'development'];
            const currentIndex = phases.indexOf(phase);

            document.querySelectorAll('.phase-step').forEach((step, index) => {
                step.classList.remove('active', 'completed');
                if (index < currentIndex) {
                    step.classList.add('completed');
                } else if (index === currentIndex) {
                    step.classList.add('active');
                }
            });

            document.querySelectorAll('.phase-connector').forEach((connector, index) => {
                connector.classList.remove('completed');
                if (index < currentIndex) {
                    connector.classList.add('completed');
                }
            });

            announce('Workflow entered ' + (phaseConfig[phase]?.name || phase) + ' phase');
        }

        function resetPhase() {
            state.currentPhase = null;
            document.querySelectorAll('.phase-step').forEach(step => {
                step.classList.remove('active', 'completed');
            });
            document.querySelectorAll('.phase-connector').forEach(connector => {
                connector.classList.remove('completed');
            });
        }

        // ===================================
        // DELEGATION CHAIN
        // ===================================
        function updateDelegationChain(from, to) {
            if (!state.delegationChain.includes(from)) {
                state.delegationChain.push(from);
            }
            renderDelegationBreadcrumb();
        }

        function renderDelegationBreadcrumb() {
            if (state.delegationChain.length === 0) {
                elements.delegationBreadcrumb.innerHTML = '';
                return;
            }

            const html = state.delegationChain.map((agent, index) => {
                const config = agentConfig[agent] || { emoji: '&#129302;', name: agent };
                const isCurrent = index === state.delegationChain.length - 1;
                return '<div class="breadcrumb-item' + (isCurrent ? ' current' : '') + '">' +
                    '<span>' + config.emoji + '</span>' +
                    '<span>' + config.name + '</span>' +
                    '</div>' +
                    (index < state.delegationChain.length - 1 ? '<span class="breadcrumb-arrow">&#8594;</span>' : '');
            }).join('');

            elements.delegationBreadcrumb.innerHTML = html;
        }

        // ===================================
        // AGENT ACTIVITY
        // ===================================
        function updateActiveAgent(agent) {
            state.activeAgent = agent;

            document.querySelectorAll('.agent-card').forEach(card => {
                const cardAgent = card.dataset.agent;
                card.classList.remove('working', 'active');

                if (cardAgent === agent) {
                    card.classList.add('working');
                } else if (state.delegationChain.includes(cardAgent)) {
                    card.classList.add('active');
                }
            });
        }

        function resetAgentStates() {
            document.querySelectorAll('.agent-card').forEach(card => {
                card.classList.remove('working', 'active');
            });
        }

        // ===================================
        // TASK MANAGEMENT
        // ===================================
        async function fetchTasks() {
            try {
                const res = await fetch('/api/tasks?project=' + encodeURIComponent(state.projectId));
                const data = await res.json();
                state.tasks = data.tasks || [];
                renderTasks();
            } catch (e) {
                console.error('Failed to fetch tasks:', e);
            }
        }

        function renderTasks() {
            const filteredTasks = filterTasks(state.tasks);

            if (filteredTasks.length === 0) {
                elements.tasksList.innerHTML = '<div class="empty-state"><p>No tasks found</p></div>';
                return;
            }

            elements.tasksList.innerHTML = filteredTasks.map(t => {
                const date = new Date(t.createdAt).toLocaleDateString();
                const statusClass = t.status || 'completed';
                const isSelected = state.selectedTaskId === t.taskId;
                const isRunning = statusClass === 'running';

                return '<div class="task-card' + (isSelected ? ' selected' : '') + (isRunning ? ' running' : '') + '" ' +
                    'data-task-id="' + t.taskId + '" ' +
                    'onclick="selectTask(\'' + t.taskId + '\')" ' +
                    'role="listitem" ' +
                    'tabindex="0" ' +
                    'aria-label="Task: ' + escapeHtml(t.summary) + '">' +
                    '<div class="task-date">' + date + '</div>' +
                    '<div class="task-summary">' + escapeHtml(t.summary) + '</div>' +
                    '<div class="task-meta">' +
                    '<span class="task-badge ' + statusClass + '">' + statusClass + '</span>' +
                    '<span>' + t.turns + ' turns</span>' +
                    '<span>' + t.filesCreated + ' files</span>' +
                    '</div>' +
                    '</div>';
            }).join('');
        }

        function filterTasks(tasks) {
            const searchTerm = elements.taskSearch.value.toLowerCase();
            if (!searchTerm) return tasks;
            return tasks.filter(t => t.summary.toLowerCase().includes(searchTerm));
        }

        async function selectTask(taskId) {
            try {
                const res = await fetch('/api/task/' + taskId + '?project=' + encodeURIComponent(state.projectId));
                const data = await res.json();

                if (data.success) {
                    state.selectedTaskId = taskId;
                    state.viewingHistory = true;
                    state.messages = data.messages || [];
                    renderMessages();
                    updateMessageCount();

                    elements.historyBanner.classList.add('visible');
                    elements.historyTaskName.textContent = data.task.task.substring(0, 50) + (data.task.task.length > 50 ? '...' : '');

                    renderTasks();
                    announce('Viewing historical task');
                }
            } catch (e) {
                console.error('Failed to fetch task:', e);
                showToast('Error', 'Failed to load task details', 'error');
            }
        }

        function backToLive() {
            state.viewingHistory = false;
            state.selectedTaskId = null;
            elements.historyBanner.classList.remove('visible');
            fetchMessages();
            renderTasks();
            announce('Returned to live view');
        }

        async function submitTask(e) {
            e.preventDefault();
            const task = elements.taskInput.value.trim();
            state.projectId = elements.projectInput.value.trim().toLowerCase().replace(/[^a-z0-9-_]/g, '-') || 'default';

            if (!task || state.isTaskRunning) return;

            if (state.viewingHistory) {
                backToLive();
            }

            state.messages = [];
            state.delegationChain = [];
            renderMessages();
            renderDelegationBreadcrumb();

            try {
                const response = await fetch('/api/task', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ task: task, project_id: state.projectId })
                });
                const data = await response.json();

                if (data.success) {
                    state.currentTaskId = data.task_id;
                    elements.taskInput.value = '';
                    fetchStatus();
                    fetchTasks();
                    announce('Task submitted');
                }
            } catch (e) {
                console.error('Failed to submit task:', e);
                showToast('Error', 'Failed to submit task', 'error');
            }
        }

        async function continueTask(e) {
            e.preventDefault();
            const task = elements.continueInput.value.trim();

            if (!task || state.isTaskRunning || !state.selectedTaskId) return;

            try {
                const response = await fetch('/api/task', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        task: task,
                        project_id: state.projectId,
                        continue_from: state.selectedTaskId
                    })
                });
                const data = await response.json();

                if (data.success) {
                    state.currentTaskId = data.task_id;
                    elements.continueInput.value = '';
                    backToLive();
                    fetchStatus();
                    fetchTasks();
                    announce('Continuing task');
                }
            } catch (e) {
                console.error('Failed to continue task:', e);
                showToast('Error', 'Failed to continue task', 'error');
            }
        }

        // ===================================
        // STATUS & FILES
        // ===================================
        async function fetchStatus() {
            try {
                const res = await fetch('/api/status');
                const data = await res.json();
                state.isTaskRunning = data.task_running;

                if (state.isTaskRunning) {
                    elements.taskStatus.classList.add('running');
                    elements.taskStatusText.textContent = 'Running';
                    elements.submitBtn.disabled = true;
                    elements.submitBtn.textContent = 'Running...';
                } else {
                    elements.taskStatus.classList.remove('running');
                    elements.taskStatusText.textContent = 'Idle';
                    elements.submitBtn.disabled = false;
                    elements.submitBtn.textContent = 'Start Task';
                }
            } catch (e) {
                console.error('Failed to fetch status:', e);
            }
        }

        async function fetchMessages() {
            try {
                const res = await fetch('/api/messages');
                const data = await res.json();
                state.messages = data.messages || [];
                renderMessages();
                updateMessageCount();
            } catch (e) {
                console.error('Failed to fetch messages:', e);
            }
        }

        async function fetchFiles() {
            try {
                const res = await fetch('/api/files?project=' + encodeURIComponent(state.projectId));
                const data = await res.json();
                state.files = data.files || [];
                renderFileTree();
            } catch (e) {
                console.error('Failed to fetch files:', e);
            }
        }

        async function fetchPhase() {
            try {
                const res = await fetch('/api/phase?project=' + encodeURIComponent(state.projectId));
                const data = await res.json();
                if (data.current_phase) {
                    updatePhase(data.current_phase);
                }
            } catch (e) {
                console.error('Failed to fetch phase:', e);
            }
        }

        function renderFileTree() {
            if (state.files.length === 0) {
                elements.fileTree.innerHTML = '<div class="empty-state"><p>No files created yet</p></div>';
                return;
            }

            // Build tree structure
            const tree = buildFileTree(state.files);
            elements.fileTree.innerHTML = renderTreeNode(tree, '');
        }

        function buildFileTree(files) {
            const tree = {};
            files.forEach(f => {
                const parts = f.path.split('/');
                let current = tree;
                parts.forEach((part, i) => {
                    if (!current[part]) {
                        current[part] = i === parts.length - 1 ? { _file: f } : {};
                    }
                    current = current[part];
                });
            });
            return tree;
        }

        function renderTreeNode(node, path) {
            let html = '';
            for (const [name, child] of Object.entries(node)) {
                if (child._file) {
                    const file = child._file;
                    const agent = file.created_by || '';
                    html += '<div class="file-item" data-path="' + file.path + '">' +
                        '<span class="file-icon">&#128196;</span>' +
                        '<span class="file-name">' + name + '</span>' +
                        (agent ? '<span class="file-agent">' + agent + '</span>' : '') +
                        '</div>';
                } else {
                    html += '<div class="file-tree-folder">' +
                        '<div class="folder-header">' +
                        '<span class="folder-icon">&#128193;</span>' +
                        '<span class="folder-name">' + name + '</span>' +
                        '</div>' +
                        '<div class="folder-children">' + renderTreeNode(child, path + name + '/') + '</div>' +
                        '</div>';
                }
            }
            return html;
        }

        // ===================================
        // TOAST NOTIFICATIONS
        // ===================================
        function showToast(title, message, type = 'info') {
            const toast = document.createElement('div');
            toast.className = 'toast ' + type;

            const icons = {
                success: '&#10004;',
                error: '&#10006;',
                warning: '&#9888;',
                info: '&#8505;'
            };

            toast.innerHTML =
                '<span class="toast-icon">' + (icons[type] || icons.info) + '</span>' +
                '<div class="toast-content">' +
                '<div class="toast-title">' + escapeHtml(title) + '</div>' +
                '<div class="toast-message">' + escapeHtml(message) + '</div>' +
                '</div>' +
                '<button class="toast-close" onclick="this.parentElement.remove()" aria-label="Close">&times;</button>';

            elements.toastContainer.appendChild(toast);

            // Auto dismiss after 5 seconds
            setTimeout(() => {
                toast.classList.add('fade-out');
                setTimeout(() => toast.remove(), 300);
            }, 5000);
        }

        // ===================================
        // ACCESSIBILITY
        // ===================================
        function announce(message) {
            elements.announcer.textContent = message;
            setTimeout(() => elements.announcer.textContent = '', 1000);
        }

        // ===================================
        // UI INTERACTIONS
        // ===================================
        function setupContextTabs() {
            elements.contextTabs.forEach(tab => {
                tab.addEventListener('click', () => {
                    const targetTab = tab.dataset.tab;

                    elements.contextTabs.forEach(t => {
                        t.classList.remove('active');
                        t.setAttribute('aria-selected', 'false');
                    });
                    tab.classList.add('active');
                    tab.setAttribute('aria-selected', 'true');

                    elements.contextSections.forEach(section => {
                        section.classList.remove('active');
                        if (section.id === targetTab + 'Panel') {
                            section.classList.add('active');
                        }
                    });
                });
            });
        }

        function setupMobileSidebar() {
            elements.mobileSidebarToggle.addEventListener('click', () => {
                const isOpen = elements.sidebar.classList.toggle('mobile-open');
                elements.mobileOverlay.classList.toggle('visible', isOpen);
                elements.mobileSidebarToggle.setAttribute('aria-expanded', isOpen);
            });

            elements.mobileOverlay.addEventListener('click', () => {
                elements.sidebar.classList.remove('mobile-open');
                elements.mobileOverlay.classList.remove('visible');
                elements.mobileSidebarToggle.setAttribute('aria-expanded', 'false');
            });
        }

        function setupKeyboardShortcuts() {
            document.addEventListener('keydown', (e) => {
                // Focus search with /
                if (e.key === '/' && document.activeElement.tagName !== 'INPUT') {
                    e.preventDefault();
                    elements.taskSearch.focus();
                }

                // Close mobile sidebar with Escape
                if (e.key === 'Escape') {
                    elements.sidebar.classList.remove('mobile-open');
                    elements.mobileOverlay.classList.remove('visible');
                }
            });

            // Task card keyboard navigation
            elements.tasksList.addEventListener('keydown', (e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                    const taskCard = e.target.closest('.task-card');
                    if (taskCard) {
                        e.preventDefault();
                        selectTask(taskCard.dataset.taskId);
                    }
                }
            });
        }

        // ===================================
        // PROJECT MANAGEMENT
        // ===================================
        async function fetchProjects() {
            try {
                const res = await fetch('/api/projects');
                const data = await res.json();

                const projectSelect = elements.projectSelect;

                // Clear existing options except default
                projectSelect.innerHTML = '<option value="default">Default Project</option>';

                // Add projects to dropdown
                if (data.projects && data.projects.length > 0) {
                    data.projects.forEach(project => {
                        const option = document.createElement('option');
                        option.value = project.id;
                        option.textContent = project.name + ' (' + project.task_count + ' tasks)';
                        projectSelect.appendChild(option);
                    });
                }

                // Set current project as selected
                if (state.projectId) {
                    projectSelect.value = state.projectId;
                }
            } catch (e) {
                console.error('Failed to fetch projects:', e);
            }
        }

        function setupProjectSelector() {
            elements.projectSelect.addEventListener('change', () => {
                state.projectId = elements.projectSelect.value || 'default';

                // Update project input field to match
                elements.projectInput.value = state.projectId;

                // Reload data for new project
                fetchFiles();
                fetchTasks();
                fetchMessages();
                updateDevelopmentPhases(state.projectId);
                updateQAFeedback(state.projectId);
                updateProgressDashboard(state.projectId);

                // Show toast
                showToast('Project Switched', 'Now viewing: ' + state.projectId, 'info');
            });
        }

        // ===================================
        // AGENT TASK MANAGEMENT
        // ===================================
        async function openAgentDetailModal(agentRole) {
            const modal = document.getElementById('agentDetailModal');
            const title = document.getElementById('agentDetailTitle');

            title.textContent = agentRole.toUpperCase() + ' Agent';

            // Fetch agent status
            try {
                const statusRes = await fetch('/api/agents/' + agentRole + '/status');
                const statusData = await statusRes.json();

                document.getElementById('agentState').textContent = statusData.state;
                document.getElementById('agentState').className = 'status-badge ' + statusData.state;
                document.getElementById('agentTasksCompleted').textContent = statusData.tasks_completed;
                document.getElementById('agentTasksFailed').textContent = statusData.tasks_failed;

                // Fetch agent tasks
                const tasksRes = await fetch('/api/agents/' + agentRole + '/tasks');
                const tasksData = await tasksRes.json();

                const tasksList = document.getElementById('agentTasksList');
                if (tasksData.tasks && tasksData.tasks.length > 0) {
                    tasksList.innerHTML = tasksData.tasks.map(task =>
                        '<div class="task-item">' +
                        '<div class="task-item-header">' +
                        '<span class="task-title">' + escapeHtml(task.title) + '</span>' +
                        '<span class="task-progress">' + task.progress + '%</span>' +
                        '</div>' +
                        '<div class="progress-bar">' +
                        '<div class="progress-bar-fill" style="width: ' + task.progress + '%"></div>' +
                        '</div>' +
                        '<div class="task-status">' + task.status + '</div>' +
                        '</div>'
                    ).join('');
                } else {
                    tasksList.innerHTML = '<p style="color: var(--color-text-secondary);">No tasks yet</p>';
                }

                // Fetch agent outputs
                const outputsRes = await fetch('/api/agents/' + agentRole + '/outputs');
                const outputsData = await outputsRes.json();

                const outputsList = document.getElementById('agentOutputsList');
                if (outputsData.outputs && outputsData.outputs.length > 0) {
                    outputsList.innerHTML = outputsData.outputs.map((output, index) =>
                        '<div class="output-item" data-output-index="' + index + '" style="cursor: pointer;">' +
                        '<div class="output-title">' + escapeHtml(output.title) + '</div>' +
                        '<div class="output-meta">' + output.type + ' • ' + output.format + '</div>' +
                        '</div>'
                    ).join('');

                    // Add click handlers
                    document.querySelectorAll('#agentOutputsList .output-item').forEach((item, index) => {
                        item.addEventListener('click', () => {
                            const output = outputsData.outputs[index];
                            viewOutput(output.file_path, output.title);
                        });
                    });
                } else {
                    outputsList.innerHTML = '<p style="color: var(--color-text-secondary);">No outputs yet</p>';
                }

            } catch (e) {
                console.error('Failed to load agent details:', e);
                showToast('Error', 'Failed to load agent details', 'error');
            }

            modal.classList.add('active');
        }

        function closeAgentDetailModal() {
            document.getElementById('agentDetailModal').classList.remove('active');
        }

        async function viewOutput(filePath, title) {
            console.log('viewOutput called with:', { filePath, title });

            const modal = document.getElementById('outputViewerModal');
            const titleEl = document.getElementById('outputViewerTitle');
            const content = document.getElementById('outputViewerContent');

            titleEl.textContent = title;
            content.innerHTML = '<div style="padding: var(--space-4); color: var(--color-text-secondary);">Loading...</div>';

            modal.classList.add('active');

            try {
                // Fetch file content from API
                const url = '/api/file-content?path=' + encodeURIComponent(filePath);
                console.log('Fetching:', url);

                const response = await fetch(url);
                console.log('Response status:', response.status);

                if (!response.ok) {
                    throw new Error('Failed to load file: ' + response.status + ' ' + response.statusText);
                }

                const fileContent = await response.text();
                console.log('File content length:', fileContent.length);
                console.log('First 100 chars:', fileContent.substring(0, 100));

                // Check if it's a markdown file
                const isMarkdown = filePath.endsWith('.md');
                console.log('Is markdown:', isMarkdown);
                console.log('marked available:', typeof marked !== 'undefined');

                if (isMarkdown && typeof marked !== 'undefined') {
                    // Configure marked for better code rendering
                    marked.setOptions({
                        breaks: true,
                        gfm: true,
                        headerIds: false,
                        mangle: false
                    });

                    // Render markdown as HTML
                    const htmlContent = marked.parse(fileContent);
                    console.log('Rendered HTML length:', htmlContent.length);
                    content.innerHTML = '<div class="markdown-content">' + htmlContent + '</div>';
                } else {
                    console.log('Showing as plain text (markdown:', isMarkdown, ', marked:', typeof marked !== 'undefined', ')');
                    // For non-markdown files, show as preformatted text
                    const pre = document.createElement('pre');
                    pre.style.cssText = 'margin: 0; padding: var(--space-4); overflow-x: auto; font-size: var(--text-sm); line-height: 1.6;';
                    pre.textContent = fileContent;
                    content.innerHTML = '';
                    content.appendChild(pre);
                }
            } catch (error) {
                console.error('Error in viewOutput:', error);
                content.innerHTML = '<div style="padding: var(--space-4); color: var(--color-status-error);">Error loading file: ' + error.message + '</div>';
            }
        }

        function closeOutputViewerModal() {
            document.getElementById('outputViewerModal').classList.remove('active');
        }

        function updateAgentStatus(agentRole, state) {
            const agentCard = document.querySelector('[data-agent="' + agentRole + '"]');
            if (!agentCard) return;

            const indicator = agentCard.querySelector('.agent-status-indicator');
            indicator.className = 'agent-status-indicator ' + state;
        }

        function handleAgentTaskEvent(event) {
            const { event_type, agent_role, data } = event;

            switch (event_type) {
                case 'task_started':
                    updateAgentStatus(agent_role, 'working');
                    showToast('Task Started', agent_role.toUpperCase() + ': ' + data.title, 'info');
                    break;

                case 'task_progress':
                    // Update progress if modal is open
                    break;

                case 'task_completed':
                    updateAgentStatus(agent_role, 'completed');
                    showToast('Task Completed', agent_role.toUpperCase() + ' finished in ' + data.duration, 'success');
                    break;

                case 'task_failed':
                    updateAgentStatus(agent_role, 'error');
                    showToast('Task Failed', agent_role.toUpperCase() + ': ' + data.error, 'error');
                    break;
            }
        }

        // ===================================
        // INITIALIZATION
        // ===================================
        function init() {
            connectWebSocket();

            fetchStatus();
            fetchMessages();
            fetchTasks();
            fetchFiles();
            fetchPhase();
            fetchProjects();

            // Event listeners
            elements.taskForm.addEventListener('submit', submitTask);
            elements.continueForm.addEventListener('submit', continueTask);
            elements.newTaskBtn.addEventListener('click', () => {
                if (state.viewingHistory) backToLive();
                elements.taskInput.focus();
            });
            elements.backToLiveBtn.addEventListener('click', backToLive);

            elements.projectInput.addEventListener('change', () => {
                state.projectId = elements.projectInput.value.trim().toLowerCase().replace(/[^a-z0-9-_]/g, '-') || 'default';
                fetchFiles();
                fetchTasks();
            });

            elements.taskSearch.addEventListener('input', renderTasks);

            // Agent card click handlers
            document.querySelectorAll('.agent-card').forEach(card => {
                card.addEventListener('click', () => {
                    const agentRole = card.dataset.agent;
                    openAgentDetailModal(agentRole);
                });
            });

            // Setup UI
            setupProjectSelector();
            setupContextTabs();
            setupMobileSidebar();
            setupKeyboardShortcuts();

            // Polling
            setInterval(fetchStatus, 2000);
            setInterval(fetchFiles, 5000);
            setInterval(fetchTasks, 10000);
            setInterval(fetchPhase, 3000);
            setInterval(fetchProjects, 30000); // Refresh projects every 30 seconds
            setInterval(function() { updateDevelopmentPhases(state.projectId); }, 5000); // Refresh dev phases every 5 seconds
            setInterval(function() { updateQAFeedback(state.projectId); }, 5000); // Refresh QA feedback every 5 seconds
            setInterval(function() { updateProgressDashboard(state.projectId); }, 5000); // Refresh progress dashboard every 5 seconds
        }

        // Start app
        init();
    </script>
</body>
</html>`
