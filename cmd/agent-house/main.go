package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
	"pty-claude-test/internal/orchestrator"
	"pty-claude-test/internal/web"
)

func main() {
	// CLI flags
	agentRole := flag.String("agent", "ceo", "Agent role to use (single agent mode)")
	task := flag.String("task", "", "Task to give the agent")
	chain := flag.Bool("chain", false, "Enable agent chaining (follow delegations)")
	projectDir := flag.String("project", "./projects", "Base directory for project files")
	orchestrate := flag.Bool("orchestrate", false, "Use full orchestrator (recommended)")
	maxDepth := flag.Int("max-depth", 3, "Maximum delegation depth")
	maxTurns := flag.Int("max-turns", 8, "Maximum total turns")
	verbose := flag.Bool("verbose", true, "Enable verbose output")
	serve := flag.Bool("serve", false, "Start web dashboard server")
	port := flag.Int("port", 8080, "Port for web server")
	flag.Parse()

	// Ensure project directory exists
	if err := os.MkdirAll(*projectDir, 0755); err != nil {
		fmt.Printf("Error creating project directory: %v\n", err)
		os.Exit(1)
	}

	// Ensure logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("Error creating logs directory: %v\n", err)
		os.Exit(1)
	}

	// Create message store
	store := message.NewStore()

	// Start web server if requested
	if *serve {
		runWebServer(*port, *projectDir, store)
		return
	}

	if *task == "" {
		fmt.Println("Error: --task is required (or use --serve for web dashboard)")
		fmt.Println("\nUsage:")
		fmt.Println("  Single agent:  agent-house --agent=ceo --task=\"Review our strategy\"")
		fmt.Println("  With chaining: agent-house --agent=ceo --task=\"Build a todo app\" --chain")
		fmt.Println("  Orchestrated:  agent-house --task=\"Build a todo app\" --orchestrate")
		fmt.Println("  Web dashboard: agent-house --serve --port=8080")
		os.Exit(1)
	}

	if *orchestrate {
		// Use full orchestrator
		runOrchestrator(*task, *projectDir, *maxDepth, *maxTurns, *verbose, store)
	} else {
		// Simple single/chain mode
		role := agent.Role(*agentRole)
		runAgent(role, *task, store, *chain, 0, *projectDir)
	}
}

func runOrchestrator(task, projectDir string, maxDepth, maxTurns int, verbose bool, store *message.Store) {
	config := orchestrator.Config{
		ProjectDir:  projectDir,
		MaxDepth:    maxDepth,
		MaxTurns:    maxTurns,
		EnableFiles: true,
		Verbose:     verbose,
	}

	orch := orchestrator.New(config, store)

	// Generate project ID
	projectID := filepath.Base(projectDir)

	// Process the task
	result, err := orch.ProcessTask(projectID, task)
	if err != nil {
		fmt.Printf("\n❌ Orchestration error: %v\n", err)
	}

	// Save conversation
	logFile := fmt.Sprintf("logs/conversation_%s.json", projectID)
	if err := store.SaveToFile(logFile); err != nil {
		fmt.Printf("Warning: Could not save conversation: %v\n", err)
	} else {
		fmt.Printf("💾 Conversation saved to %s (%d messages)\n", logFile, store.Count())
	}

	// Summary
	if result != nil {
		fmt.Printf("\n📊 Summary:\n")
		fmt.Printf("   Messages: %d\n", len(result.Messages))
		fmt.Printf("   Files created: %d\n", len(result.Files))
		fmt.Printf("   Duration: %s\n", result.Duration.Round(1e8))

		if len(result.Files) > 0 {
			fmt.Printf("\n📁 Files:\n")
			for _, f := range result.Files {
				if f.Success {
					fmt.Printf("   ✅ %s (by %s)\n", f.Path, f.Agent)
				} else {
					fmt.Printf("   ❌ %s - %v\n", f.Path, f.Error)
				}
			}
		}
	}
}

func runAgent(role agent.Role, task string, store *message.Store, chain bool, depth int, projectDir string) {
	if depth > 5 {
		fmt.Println("\n⚠️  Max delegation depth reached")
		return
	}

	// Create the agent
	a, err := agent.NewAgent(role)
	if err != nil {
		fmt.Printf("Error creating agent %s: %v\n", role, err)
		return
	}

	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s🤖 %s is processing...\n\n", indent, a.Name)

	// Process the task with file operations
	response, fileResults, err := a.ProcessWithFileOps(task, projectDir)
	if err != nil {
		fmt.Printf("%sError: %v\n", indent, err)
		return
	}

	// Store the message
	msg := message.NewMessage(message.TypeResponse, string(role), "user", response.Content)
	store.Add(msg)

	// Display the response
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n", indent)
	fmt.Printf("%s📋 %s Response:\n", indent, a.Name)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n", indent)
	fmt.Println(response.Content)
	fmt.Println()

	// Show file operations
	if len(fileResults) > 0 {
		fmt.Printf("%s📁 File Operations:\n", indent)
		for _, fr := range fileResults {
			if fr.Success {
				fmt.Printf("%s   ✅ Created: %s\n", indent, fr.Path)
			} else {
				fmt.Printf("%s   ❌ Failed: %s - %v\n", indent, fr.Path, fr.Error)
			}
		}
		fmt.Println()
	}

	// Show delegations
	if len(response.DelegateTo) > 0 {
		fmt.Printf("%s📨 Delegating to: %v\n\n", indent, response.DelegateTo)

		if chain {
			// Run each delegated agent
			for _, delegateRole := range response.DelegateTo {
				// Create context for the next agent
				context := fmt.Sprintf(`You are receiving a task from the %s.

Original task: %s

%s's analysis:
%s

Now provide your specialized input.`, a.Name, task, a.Name, response.Content)

				runAgent(delegateRole, context, store, chain, depth+1, projectDir)
			}
		}
	}

	// Save conversation
	if depth == 0 {
		store.SaveToFile("logs/conversation.json")
		fmt.Printf("\n💾 Conversation saved to logs/conversation.json (%d messages)\n", store.Count())
	}
}

func runWebServer(port int, projectDir string, store *message.Store) {
	fmt.Println("🏠 Agent House - Web Dashboard")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := web.Config{
		Port:       port,
		ProjectDir: projectDir,
		Store:      store,
	}

	server := web.NewServer(config)

	fmt.Printf("📁 Project directory: %s\n", projectDir)
	fmt.Printf("🌐 Dashboard: http://localhost:%d\n", port)
	fmt.Printf("📡 API: http://localhost:%d/api/\n\n", port)

	if err := server.Start(port); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
