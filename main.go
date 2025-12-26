package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

func main() {
	mode := flag.String("mode", "print", "Mode: 'print' or 'interactive'")
	flag.Parse()

	switch *mode {
	case "print":
		testPrintMode()
	case "interactive":
		testInteractiveMode()
	default:
		fmt.Printf("Unknown mode: %s\n", *mode)
		os.Exit(1)
	}
}

func testPrintMode() {
	fmt.Println("=== Claude --print Mode Test ===")

	prompt := "Say 'Hello PTY test!' and nothing else"

	output, err := runClaudePrint(prompt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("--- Captured Output ---")
	fmt.Println(output)
	fmt.Println("--- End ---")

	saveLog(output, "print")
}

func testInteractiveMode() {
	fmt.Println("=== Claude Interactive Mode Test ===")
	fmt.Println("Starting Claude session...")

	session, err := NewInteractiveSession()
	if err != nil {
		fmt.Printf("Error starting session: %v\n", err)
		os.Exit(1)
	}
	defer session.Close()

	// Wait for Claude to initialize - look for the input prompt indicator ">"
	fmt.Println("\n[Waiting for Claude to be ready...]")
	ready := session.WaitForPattern("> ", 10*time.Second)
	if !ready {
		fmt.Println("[Warning: Did not detect ready state, proceeding anyway]")
	}
	session.Wait(500 * time.Millisecond) // Extra buffer for UI to stabilize

	// Send our prompt - typing character by character may be more reliable
	prompt := "Say 'Hello interactive test!' and nothing else"
	fmt.Printf("\n[Sending prompt: %s]\n", prompt)

	// Type the prompt
	for _, ch := range prompt {
		session.SendInput(string(ch))
		time.Sleep(10 * time.Millisecond) // Small delay between characters
	}

	// Press Enter to submit
	time.Sleep(100 * time.Millisecond)
	session.SendInput("\r") // Carriage return (Enter key)

	// Wait for response - look for completion indicator
	fmt.Println("\n[Waiting for response...]")
	session.Wait(15 * time.Second)

	// Exit Claude with Escape key then /exit
	fmt.Println("\n[Exiting...]")
	session.SendKey(27) // Escape key to clear any prompts
	session.Wait(200 * time.Millisecond)
	session.SendLine("/exit")
	session.Wait(500 * time.Millisecond)
	session.SendInput("\r") // Confirm exit
	session.Wait(2 * time.Second)

	// Get final output
	output := session.GetOutput()
	fmt.Println("\n--- Session Complete ---")

	saveLog(output, "interactive")
}

func runClaudePrint(prompt string) (string, error) {
	cmd := exec.Command("claude", "--print", prompt)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return "", fmt.Errorf("pty start: %w", err)
	}
	defer ptmx.Close()

	var output bytes.Buffer
	done := make(chan error, 1)

	go func() {
		_, err := io.Copy(&output, ptmx)
		done <- err
	}()

	cmdErr := cmd.Wait()
	<-done

	if cmdErr != nil {
		return output.String(), fmt.Errorf("command error: %w (output: %s)", cmdErr, output.String())
	}

	return output.String(), nil
}

func saveLog(content, mode string) {
	timestamp := time.Now().Format("20060102_150405")

	// Save raw log (with ANSI codes)
	rawFile := fmt.Sprintf("logs/session_raw_%s_%s.log", mode, timestamp)
	if err := os.WriteFile(rawFile, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing raw log: %v\n", err)
	} else {
		fmt.Printf("Raw log saved to: %s\n", rawFile)
	}

	// Save clean log (ANSI stripped)
	cleanContent := StripANSI(content)
	cleanFile := fmt.Sprintf("logs/session_clean_%s_%s.log", mode, timestamp)
	if err := os.WriteFile(cleanFile, []byte(cleanContent), 0644); err != nil {
		fmt.Printf("Error writing clean log: %v\n", err)
	} else {
		fmt.Printf("Clean log saved to: %s\n", cleanFile)
	}

	// Also extract and show the response
	response := ExtractClaudeResponse(content)
	fmt.Println("\n--- Extracted Response ---")
	fmt.Println(response)
}
