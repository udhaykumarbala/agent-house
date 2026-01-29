package web

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"pty-claude-test/internal/agentask"
	"pty-claude-test/internal/message"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Hub manages WebSocket connections
type Hub struct {
	clients         map[*Client]bool
	broadcast       chan *message.Message
	agentTaskEvents chan *AgentTaskEvent
	register        chan *Client
	unregister      chan *Client
	mu              sync.RWMutex
}

// AgentTaskEvent represents an event related to agent tasks
type AgentTaskEvent struct {
	EventType string                  `json:"event_type"` // task_started, task_progress, task_completed, task_failed
	TaskID    string                  `json:"task_id"`
	AgentRole string                  `json:"agent_role"`
	Data      map[string]interface{}  `json:"data"`
}

// Client represents a WebSocket client
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		broadcast:       make(chan *message.Message),
		agentTaskEvents: make(chan *AgentTaskEvent, 100),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket client connected (%d total)", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected (%d remaining)", len(h.clients))

		case msg := <-h.broadcast:
			data, err := json.Marshal(map[string]interface{}{
				"type":    "message",
				"message": msg,
			})
			if err != nil {
				continue
			}

			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()

		case event := <-h.agentTaskEvents:
			data, err := json.Marshal(map[string]interface{}{
				"type":  "agent_task_event",
				"event": event,
			})
			if err != nil {
				continue
			}

			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(msg *message.Message) {
	h.broadcast <- msg
}

// BroadcastAgentTaskEvent sends an agent task event to all connected clients
func (h *Hub) BroadcastAgentTaskEvent(event *AgentTaskEvent) {
	select {
	case h.agentTaskEvents <- event:
	default:
		log.Printf("Warning: Agent task event channel full, dropping event")
	}
}

// BroadcastTaskStarted broadcasts a task started event
func (h *Hub) BroadcastTaskStarted(task *agentask.AgentTask) {
	h.BroadcastAgentTaskEvent(&AgentTaskEvent{
		EventType: "task_started",
		TaskID:    task.ID,
		AgentRole: string(task.AgentRole),
		Data: map[string]interface{}{
			"title":       task.Title,
			"task_type":   task.TaskType,
			"description": task.Description,
		},
	})
}

// BroadcastTaskProgress broadcasts a task progress event
func (h *Hub) BroadcastTaskProgress(task *agentask.AgentTask) {
	h.BroadcastAgentTaskEvent(&AgentTaskEvent{
		EventType: "task_progress",
		TaskID:    task.ID,
		AgentRole: string(task.AgentRole),
		Data: map[string]interface{}{
			"progress": task.Progress,
		},
	})
}

// BroadcastTaskCompleted broadcasts a task completed event
func (h *Hub) BroadcastTaskCompleted(task *agentask.AgentTask) {
	h.BroadcastAgentTaskEvent(&AgentTaskEvent{
		EventType: "task_completed",
		TaskID:    task.ID,
		AgentRole: string(task.AgentRole),
		Data: map[string]interface{}{
			"duration": task.Duration.String(),
			"outputs":  task.Outputs,
		},
	})
}

// BroadcastTaskFailed broadcasts a task failed event
func (h *Hub) BroadcastTaskFailed(task *agentask.AgentTask) {
	h.BroadcastAgentTaskEvent(&AgentTaskEvent{
		EventType: "task_failed",
		TaskID:    task.ID,
		AgentRole: string(task.AgentRole),
		Data: map[string]interface{}{
			"error": task.Error,
		},
	})
}

// ClientCount returns the number of connected clients
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// readPump handles incoming messages from the client
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// We don't process incoming messages for now
		// Could be used for task submission or other commands
	}
}

// writePump handles outgoing messages to the client
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// ServeWS handles WebSocket upgrade requests
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}

	h.register <- client

	// Send welcome message
	welcome, _ := json.Marshal(map[string]interface{}{
		"type":    "connected",
		"message": "Connected to Agent House",
	})
	client.send <- welcome

	go client.writePump()
	go client.readPump()
}

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}
