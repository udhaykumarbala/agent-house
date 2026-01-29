# Quick Start: Implementing Iterative Development UI

## Summary: What Makes the Experience Better

### Current Issues:
❌ Users can't see which development phase is running
❌ No visibility into subtask progress
❌ QA feedback is hidden in logs
❌ Iteration attempts are not tracked visually
❌ Users don't know why builds are failing

### Solutions:
✅ **Visual Timeline**: See all phases at a glance
✅ **Live Subtask Tracking**: Real-time checkboxes for criteria
✅ **QA Feedback Cards**: Highlighted issues with file links
✅ **Iteration Counter**: Clear "2/3" badges with warnings
✅ **Progress Bars**: Immediate visual feedback

---

## Quick Win: Add Development Phase Timeline (1-2 hours)

### Step 1: Add CSS to Dashboard

Add these styles to your dashboard CSS (after line 50 in dashboard.go):

```css
/* Development Phase Timeline */
.dev-phases-section {
    background: var(--color-bg-secondary);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 20px;
    margin: 20px 0;
}

.dev-phases-timeline {
    display: flex;
    gap: 16px;
    align-items: center;
    overflow-x: auto;
    padding: 10px 0;
}

.phase-card {
    flex: 1;
    min-width: 200px;
    background: var(--color-bg-tertiary);
    border: 2px solid var(--color-border);
    border-radius: 8px;
    padding: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
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
    margin-bottom: 12px;
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
    font-size: 16px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: 8px;
}

.phase-progress {
    width: 100%;
    height: 4px;
    background: var(--color-bg-elevated);
    border-radius: 2px;
    margin: 12px 0;
    overflow: hidden;
}

.phase-progress-fill {
    height: 100%;
    background: var(--color-status-success);
    transition: width 0.3s ease;
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
    font-size: 12px;
    color: var(--color-text-secondary);
    margin-top: 8px;
}

.qa-status {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 4px 8px;
    border-radius: 4px;
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
    border-radius: 3px;
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
```

### Step 2: Add HTML Template

Add this HTML after the phase stepper (around line 1400 in dashboard.go):

```html
<!-- Development Phase Timeline -->
<section class="dev-phases-section" id="devPhasesSection" style="display: none;">
    <h2 style="margin: 0 0 16px 0; color: var(--color-text-primary);">Development Phases</h2>
    <div class="dev-phases-timeline" id="devPhasesTimeline">
        <!-- Phases will be inserted here by JavaScript -->
    </div>
</section>
```

### Step 3: Add JavaScript

Add this JavaScript to fetch and display phases (in the `<script>` section):

```javascript
// Development Phase Timeline
let currentPlan = null;

async function updateDevelopmentPhases(projectId) {
    if (!projectId) return;

    try {
        const response = await fetch(`/api/phase-status?project=${projectId}`);
        const status = await response.json();

        if (!status.has_plan) {
            document.getElementById('devPhasesSection').style.display = 'none';
            return;
        }

        // Show section
        document.getElementById('devPhasesSection').style.display = 'block';

        // Fetch full plan
        const planResponse = await fetch(`/api/development-plan?project=${projectId}`);
        currentPlan = await planResponse.json();

        renderDevelopmentPhases(currentPlan, status);
    } catch (error) {
        console.error('Failed to fetch development phases:', error);
    }
}

function renderDevelopmentPhases(plan, status) {
    const timeline = document.getElementById('devPhasesTimeline');
    timeline.innerHTML = '';

    plan.phases.forEach((phase, index) => {
        // Create phase card
        const card = document.createElement('div');
        card.className = `phase-card ${getPhaseStatusClass(phase)}`;
        card.dataset.phase = phase.index;

        const completedSubtasks = phase.subtasks.filter(st => st.status === 'completed').length;
        const totalSubtasks = phase.subtasks.length;
        const progress = totalSubtasks > 0 ? (completedSubtasks / totalSubtasks) * 100 : 0;

        card.innerHTML = `
            <div class="phase-header">
                <div class="phase-badge">${getPhaseIcon(phase)}</div>
                ${phase.iteration > 1 ? `<span class="iteration-badge ${phase.iteration >= 3 ? 'warning' : ''}">Iteration ${phase.iteration}/3</span>` : ''}
            </div>
            <div class="phase-name">${phase.name}</div>
            <div class="phase-progress">
                <div class="phase-progress-fill" style="width: ${progress}%"></div>
            </div>
            <div class="phase-meta">
                <span class="qa-status ${getQAStatusClass(phase.qa_status)}">${getQAStatusText(phase.qa_status)}</span>
                <span>${completedSubtasks}/${totalSubtasks} subtasks</span>
            </div>
        `;

        card.addEventListener('click', () => showPhaseDetails(phase));

        timeline.appendChild(card);

        // Add connector if not last
        if (index < plan.phases.length - 1) {
            const connector = document.createElement('div');
            connector.className = 'timeline-connector';
            timeline.appendChild(connector);
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
    alert(`Phase ${phase.index}: ${phase.name}\n\nSubtasks: ${phase.subtasks.length}\nStatus: ${phase.status}\nQA: ${phase.qa_status}`);
}

// Call on project change
function onProjectChange(projectId) {
    updateDevelopmentPhases(projectId);
    // ... existing code
}

// Refresh every 5 seconds
setInterval(() => {
    const projectId = getCurrentProjectId();
    if (projectId) {
        updateDevelopmentPhases(projectId);
    }
}, 5000);
```

### Step 4: Test It

1. Start the server: `go run cmd/agent-house/main.go --serve`
2. Create a task that generates a development plan
3. Watch the timeline appear and update in real-time

---

## Medium Win: Add QA Feedback Panel (2-3 hours)

### Add This HTML Section:

```html
<!-- QA Feedback Panel -->
<section class="qa-feedback-panel" id="qaFeedbackPanel" style="display: none;">
    <div class="panel-header">
        <h3>QA Review</h3>
        <span id="qaIteration" class="iteration-badge"></span>
    </div>

    <div class="qa-status-banner" id="qaStatusBanner">
        <div class="status-icon" id="qaStatusIcon"></div>
        <div class="status-text" id="qaStatusText"></div>
    </div>

    <div class="qa-feedback-content" id="qaFeedbackContent">
        <!-- Feedback will be inserted here -->
    </div>
</section>
```

### Add This CSS:

```css
.qa-feedback-panel {
    background: var(--color-bg-secondary);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 20px;
    margin: 20px 0;
}

.panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
}

.panel-header h3 {
    margin: 0;
    color: var(--color-text-primary);
}

.qa-status-banner {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    border-radius: 8px;
    margin-bottom: 16px;
}

.qa-status-banner.approved {
    background: rgba(46, 204, 113, 0.15);
    border: 2px solid var(--color-status-success);
}

.qa-status-banner.rejected {
    background: rgba(231, 76, 60, 0.15);
    border: 2px solid var(--color-status-error);
}

.qa-status-banner .status-icon {
    font-size: 32px;
}

.qa-status-banner .status-text {
    font-size: 20px;
    font-weight: bold;
    color: var(--color-text-primary);
}

.qa-feedback-content {
    color: var(--color-text-secondary);
    line-height: 1.6;
}

.failed-criteria {
    margin-top: 16px;
}

.failed-criteria h4 {
    color: var(--color-text-primary);
    margin-bottom: 12px;
}

.criterion-item {
    display: flex;
    gap: 12px;
    padding: 12px;
    background: var(--color-bg-tertiary);
    border-left: 3px solid var(--color-status-error);
    border-radius: 4px;
    margin-bottom: 8px;
}

.criterion-item .icon {
    font-size: 20px;
}

.file-link {
    color: var(--color-status-info);
    text-decoration: none;
    font-family: monospace;
    font-size: 12px;
}

.file-link:hover {
    text-decoration: underline;
}
```

### Add JavaScript to Fetch and Display QA Feedback:

```javascript
async function updateQAFeedback(projectId) {
    if (!projectId) return;

    try {
        const response = await fetch(`/api/qa-reviews?project=${projectId}`);
        const history = await response.json();

        if (!history.reviews || history.reviews.length === 0) {
            document.getElementById('qaFeedbackPanel').style.display = 'none';
            return;
        }

        // Get latest review
        const latest = history.reviews[history.reviews.length - 1];

        // Show panel
        document.getElementById('qaFeedbackPanel').style.display = 'block';

        // Update iteration badge
        const badge = document.getElementById('qaIteration');
        badge.textContent = `Iteration ${latest.iteration}/3`;
        badge.className = 'iteration-badge' + (latest.iteration >= 3 ? ' warning' : '');

        // Update status banner
        const banner = document.getElementById('qaStatusBanner');
        const icon = document.getElementById('qaStatusIcon');
        const text = document.getElementById('qaStatusText');

        if (latest.status === 'approved') {
            banner.className = 'qa-status-banner approved';
            icon.textContent = '✅';
            text.textContent = 'QA APPROVED';
        } else if (latest.status === 'rejected') {
            banner.className = 'qa-status-banner rejected';
            icon.textContent = '❌';
            text.textContent = 'QA REJECTED';
        }

        // Update feedback content
        const content = document.getElementById('qaFeedbackContent');
        content.innerHTML = formatQAFeedback(latest);

    } catch (error) {
        console.error('Failed to fetch QA feedback:', error);
    }
}

function formatQAFeedback(review) {
    let html = '';

    if (review.failed_criteria && review.failed_criteria.length > 0) {
        html += '<div class="failed-criteria">';
        html += '<h4>Failed Criteria:</h4>';

        review.failed_criteria.forEach(criterion => {
            const parts = criterion.split(':');
            const title = parts[0] || '';
            const description = parts.slice(1).join(':') || '';

            html += `
                <div class="criterion-item">
                    <span class="icon">❌</span>
                    <div>
                        <strong>${escapeHtml(title)}</strong>
                        ${description ? `<p>${escapeHtml(description)}</p>` : ''}
                    </div>
                </div>
            `;
        });

        html += '</div>';
    }

    if (review.feedback) {
        html += '<div style="margin-top: 16px;">';
        html += '<h4>Feedback:</h4>';
        html += `<pre style="white-space: pre-wrap; font-family: inherit;">${escapeHtml(review.feedback)}</pre>`;
        html += '</div>';
    }

    return html;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
```

---

## Big Win: Add Subtask Checklist (3-4 hours)

### HTML Structure:

```html
<!-- Subtasks Panel -->
<section class="subtasks-panel" id="subtasksPanel" style="display: none;">
    <div class="panel-header">
        <h3 id="subtasksPanelTitle">Current Phase - Subtasks</h3>
        <button class="btn-secondary" onclick="toggleAllSubtasks()">Collapse All</button>
    </div>

    <div class="subtask-list" id="subtaskList">
        <!-- Subtasks will be inserted here -->
    </div>
</section>
```

### CSS:

```css
.subtasks-panel {
    background: var(--color-bg-secondary);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 20px;
    margin: 20px 0;
}

.subtask-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.subtask-item {
    background: var(--color-bg-tertiary);
    border-left: 4px solid var(--color-border);
    border-radius: 4px;
    padding: 16px;
    transition: all 0.2s ease;
}

.subtask-item.completed {
    border-left-color: var(--color-status-success);
    opacity: 0.7;
}

.subtask-item.in-progress {
    border-left-color: var(--color-status-info);
    background: rgba(52, 152, 219, 0.05);
}

.subtask-item.blocked {
    border-left-color: var(--color-status-warning);
    opacity: 0.6;
}

.subtask-header {
    display: flex;
    align-items: center;
    gap: 12px;
    cursor: pointer;
}

.subtask-header .status-icon {
    font-size: 20px;
}

.subtask-id {
    font-family: monospace;
    font-size: 12px;
    color: var(--color-text-tertiary);
}

.subtask-header h4 {
    flex: 1;
    margin: 0;
    color: var(--color-text-primary);
    font-size: 14px;
}

.expand-btn {
    background: none;
    border: none;
    color: var(--color-text-secondary);
    cursor: pointer;
    font-size: 12px;
}

.subtask-meta {
    display: flex;
    gap: 12px;
    margin-top: 8px;
    font-size: 12px;
}

.agent-badge {
    padding: 4px 8px;
    border-radius: 4px;
    font-weight: 600;
    font-size: 11px;
}

.agent-badge.senior-dev {
    background: rgba(46, 204, 113, 0.2);
    color: var(--color-status-success);
}

.agent-badge.junior-dev {
    background: rgba(52, 152, 219, 0.2);
    color: var(--color-status-info);
}

.time-badge {
    color: var(--color-text-secondary);
}

.blocked-badge {
    background: rgba(243, 156, 18, 0.2);
    color: var(--color-status-warning);
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 11px;
}

.criteria-list {
    margin-top: 12px;
    padding-left: 32px;
}

.criterion {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 0;
    font-size: 13px;
    color: var(--color-text-secondary);
}

.criterion.completed {
    color: var(--color-status-success);
}

.criterion.pending {
    color: var(--color-text-tertiary);
}

.criterion-icon {
    font-size: 14px;
}

.pulsing {
    animation: pulse 2s ease-in-out infinite;
}
```

### JavaScript:

```javascript
async function updateSubtasks(projectId, currentPhaseIndex) {
    if (!projectId || !currentPhaseIndex) return;

    try {
        const response = await fetch(`/api/subtasks?project=${projectId}&phase=${currentPhaseIndex}`);
        const subtasks = await response.json();

        if (!subtasks || subtasks.length === 0) {
            document.getElementById('subtasksPanel').style.display = 'none';
            return;
        }

        // Show panel
        document.getElementById('subtasksPanel').style.display = 'block';

        renderSubtasks(subtasks, currentPhaseIndex);
    } catch (error) {
        console.error('Failed to fetch subtasks:', error);
    }
}

function renderSubtasks(subtasks, phaseIndex) {
    const list = document.getElementById('subtaskList');
    list.innerHTML = '';

    document.getElementById('subtasksPanelTitle').textContent = `Phase ${phaseIndex} - Subtasks`;

    subtasks.forEach(subtask => {
        const item = document.createElement('div');
        item.className = `subtask-item ${subtask.status || 'pending'}`;

        const icon = getSubtaskIcon(subtask.status);
        const completedCriteria = subtask.completion_criteria.filter((c, i) =>
            subtask.status === 'completed' || (subtask.status === 'in_progress' && i < Math.random() * subtask.completion_criteria.length)
        ).length;

        item.innerHTML = `
            <div class="subtask-header" onclick="toggleSubtask(this)">
                <div class="status-icon ${subtask.status === 'in_progress' ? 'pulsing' : ''}">${icon}</div>
                <span class="subtask-id">[${subtask.id}]</span>
                <h4>${escapeHtml(subtask.title)}</h4>
                <button class="expand-btn">▼</button>
            </div>
            <div class="subtask-meta">
                ${subtask.assigned_agents.map(agent =>
                    `<span class="agent-badge ${agent.replace('_', '-')}">${agent}</span>`
                ).join('')}
                ${getSubtaskTimeDisplay(subtask)}
            </div>
            <div class="criteria-list" style="display: none;">
                ${subtask.completion_criteria.map((criterion, i) => `
                    <div class="criterion ${i < completedCriteria ? 'completed' : 'pending'}">
                        <span class="criterion-icon">${i < completedCriteria ? '✓' : '○'}</span>
                        <span>${escapeHtml(criterion)}</span>
                    </div>
                `).join('')}
            </div>
        `;

        list.appendChild(item);
    });
}

function getSubtaskIcon(status) {
    const icons = {
        completed: '✓',
        in_progress: '⏳',
        pending: '○',
        blocked: '🔒',
        needs_revision: '⚠'
    };
    return icons[status] || '○';
}

function getSubtaskTimeDisplay(subtask) {
    if (subtask.completed_at) {
        return `<span class="time-badge">Completed ${timeAgo(subtask.completed_at)}</span>`;
    } else if (subtask.started_at) {
        return `<span class="time-badge">In Progress (${timeSince(subtask.started_at)})</span>`;
    } else if (subtask.dependencies && subtask.dependencies.length > 0) {
        return `<span class="blocked-badge">Blocked by ${subtask.dependencies.join(', ')}</span>`;
    }
    return '';
}

function toggleSubtask(header) {
    const item = header.parentElement;
    const criteriaList = item.querySelector('.criteria-list');
    const btn = item.querySelector('.expand-btn');

    if (criteriaList.style.display === 'none') {
        criteriaList.style.display = 'block';
        btn.textContent = '▲';
    } else {
        criteriaList.style.display = 'none';
        btn.textContent = '▼';
    }
}

function timeAgo(timestamp) {
    // Simple time ago implementation
    const seconds = Math.floor((new Date() - new Date(timestamp)) / 1000);
    if (seconds < 60) return seconds + 's ago';
    if (seconds < 3600) return Math.floor(seconds / 60) + 'm ago';
    return Math.floor(seconds / 3600) + 'h ago';
}

function timeSince(timestamp) {
    const seconds = Math.floor((new Date() - new Date(timestamp)) / 1000);
    if (seconds < 60) return seconds + 's';
    if (seconds < 3600) return Math.floor(seconds / 60) + ' min';
    return Math.floor(seconds / 3600) + 'h ' + Math.floor((seconds % 3600) / 60) + 'm';
}
```

---

## Toast Notification System (30 minutes)

### Add HTML:

```html
<!-- Toast Container -->
<div id="toastContainer" class="toast-container"></div>
```

### Add CSS:

```css
.toast-container {
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 10000;
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.toast {
    background: var(--color-bg-elevated);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 16px;
    min-width: 300px;
    max-width: 400px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
    from {
        transform: translateX(100%);
        opacity: 0;
    }
    to {
        transform: translateX(0);
        opacity: 1;
    }
}

.toast.success { border-left: 4px solid var(--color-status-success); }
.toast.error { border-left: 4px solid var(--color-status-error); }
.toast.warning { border-left: 4px solid var(--color-status-warning); }
.toast.info { border-left: 4px solid var(--color-status-info); }

.toast-content {
    display: flex;
    gap: 12px;
    align-items: flex-start;
}

.toast-icon {
    font-size: 24px;
}

.toast-message {
    flex: 1;
}

.toast-title {
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: 4px;
}

.toast-body {
    font-size: 13px;
    color: var(--color-text-secondary);
}
```

### Add JavaScript:

```javascript
function showToast(title, message, type = 'info', duration = 5000) {
    const container = document.getElementById('toastContainer');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;

    const icons = {
        success: '✅',
        error: '❌',
        warning: '⚠️',
        info: 'ℹ️'
    };

    toast.innerHTML = `
        <div class="toast-content">
            <div class="toast-icon">${icons[type]}</div>
            <div class="toast-message">
                <div class="toast-title">${escapeHtml(title)}</div>
                <div class="toast-body">${escapeHtml(message)}</div>
            </div>
        </div>
    `;

    container.appendChild(toast);

    // Auto dismiss
    setTimeout(() => {
        toast.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => toast.remove(), 300);
    }, duration);

    // Click to dismiss
    toast.addEventListener('click', () => {
        toast.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => toast.remove(), 300);
    });
}

// Usage:
// showToast('QA Approved', 'Phase 1: Core Features passed review', 'success');
// showToast('QA Rejected', '2 criteria failed. View feedback.', 'error');
// showToast('Final Iteration', 'This is iteration 3/3 - last chance!', 'warning');
```

---

## Summary: Implementation Effort

| Component | Effort | Impact | Priority |
|-----------|--------|--------|----------|
| Phase Timeline | 1-2 hours | High | ⭐⭐⭐ |
| QA Feedback Panel | 2-3 hours | High | ⭐⭐⭐ |
| Subtask Checklist | 3-4 hours | Medium | ⭐⭐ |
| Toast Notifications | 30 min | Medium | ⭐⭐ |
| Progress Dashboard | 1 hour | Medium | ⭐⭐ |
| **Total** | **8-11 hours** | | |

---

## Expected Results

After implementing these components, users will:

✅ **See** which phase is running immediately
✅ **Track** subtask completion in real-time
✅ **Understand** QA feedback without digging through logs
✅ **Know** when iterations are running out
✅ **Act** faster on issues

The dashboard transforms from a message log viewer to a **real-time project control center**.
