package session

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// SessionManager manages all agent sessions across projects.
type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*AgentSession // key: "{projectID}:{agentRole}"

	// Global event bus — all agent events flow here for WebSocket broadcast
	eventBus   chan AgentEvent
	eventSubs  []chan AgentEvent
	eventSubMu sync.Mutex
}

// NewSessionManager creates a new session manager.
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*AgentSession),
		eventBus: make(chan AgentEvent, 1024),
	}
	go sm.eventRouter()
	return sm
}

// eventRouter distributes events from the eventBus to all subscribers.
func (sm *SessionManager) eventRouter() {
	for ev := range sm.eventBus {
		sm.eventSubMu.Lock()
		for _, sub := range sm.eventSubs {
			select {
			case sub <- ev:
			default:
				// Slow subscriber, drop
			}
		}
		sm.eventSubMu.Unlock()
	}
}

// sessionKey builds the map key for a session.
func sessionKey(projectID, agentRole string) string {
	return projectID + ":" + agentRole
}

// GetOrCreate returns an existing session or creates a new one.
func (sm *SessionManager) GetOrCreate(config SessionConfig) (*AgentSession, error) {
	key := sessionKey(config.ProjectID, config.AgentRole)

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sess, ok := sm.sessions[key]; ok {
		// Check if the session is still usable
		if sess.State != StateTerminated {
			return sess, nil
		}
		// Terminated session — clean up and create new
		delete(sm.sessions, key)
	}

	// Create new session
	sess := NewAgentSession(config)

	// Wire session events to the global event bus
	evCh, _ := sess.Subscribe()
	go func() {
		for ev := range evCh {
			select {
			case sm.eventBus <- ev:
			default:
				// Event bus full, drop
			}
		}
	}()

	sm.sessions[key] = sess
	log.Printf("[SESSION-MGR] created session for %s in project %s", config.AgentRole, config.ProjectID)

	return sess, nil
}

// Get returns an existing session (nil if not found).
func (sm *SessionManager) Get(projectID, agentRole string) *AgentSession {
	key := sessionKey(projectID, agentRole)
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.sessions[key]
}

// ListByProject returns all sessions for a project.
func (sm *SessionManager) ListByProject(projectID string) []*AgentSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var sessions []*AgentSession
	prefix := projectID + ":"
	for key, sess := range sm.sessions {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			sessions = append(sessions, sess)
		}
	}
	return sessions
}

// ListAll returns all sessions.
func (sm *SessionManager) ListAll() []SessionInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	list := make([]SessionInfo, 0, len(sm.sessions))
	for _, sess := range sm.sessions {
		list = append(list, sess.GetInfo())
	}
	return list
}

// TerminateProject terminates all sessions for a project.
func (sm *SessionManager) TerminateProject(projectID string) {
	sm.mu.Lock()
	var toTerminate []*AgentSession
	prefix := projectID + ":"
	for key, sess := range sm.sessions {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			toTerminate = append(toTerminate, sess)
			delete(sm.sessions, key)
		}
	}
	sm.mu.Unlock()

	for _, sess := range toTerminate {
		sess.Terminate()
	}
	log.Printf("[SESSION-MGR] terminated %d sessions for project %s", len(toTerminate), projectID)
}

// Shutdown terminates all sessions.
func (sm *SessionManager) Shutdown() {
	sm.mu.Lock()
	sessions := make([]*AgentSession, 0, len(sm.sessions))
	for _, sess := range sm.sessions {
		sessions = append(sessions, sess)
	}
	sm.sessions = make(map[string]*AgentSession)
	sm.mu.Unlock()

	for _, sess := range sessions {
		sess.Terminate()
	}
	log.Printf("[SESSION-MGR] shut down %d sessions", len(sessions))
}

// SubscribeAll returns a channel that receives ALL agent events across all sessions.
// This is used by the WebSocket hub to broadcast to dashboard clients.
func (sm *SessionManager) SubscribeAll() <-chan AgentEvent {
	ch := make(chan AgentEvent, 512)
	sm.eventSubMu.Lock()
	sm.eventSubs = append(sm.eventSubs, ch)
	sm.eventSubMu.Unlock()
	return ch
}

// GetProjectMetrics returns aggregated metrics for a project.
func (sm *SessionManager) GetProjectMetrics(projectID string) map[string]interface{} {
	sessions := sm.ListByProject(projectID)

	totalCost := 0.0
	totalTokens := 0
	agentMetrics := make(map[string]interface{})

	for _, sess := range sessions {
		stats := sess.GetStats()
		totalCost += stats.TotalCostUSD
		totalTokens += stats.TotalInputTokens + stats.TotalOutputTokens
		agentMetrics[sess.Config.AgentRole] = map[string]interface{}{
			"name":          sess.Config.AgentName,
			"state":         sess.State,
			"input_tokens":  stats.TotalInputTokens,
			"output_tokens": stats.TotalOutputTokens,
			"cost_usd":      stats.TotalCostUSD,
			"turns":         stats.TotalTurns,
			"tool_calls":    stats.TotalToolCalls,
			"tools":         stats.ToolBreakdown,
			"created_at":    sess.CreatedAt.Format(time.RFC3339),
		}
	}

	return map[string]interface{}{
		"project_id":    projectID,
		"total_cost":    fmt.Sprintf("%.4f", totalCost),
		"total_tokens":  totalTokens,
		"agent_count":   len(sessions),
		"agent_metrics": agentMetrics,
	}
}
