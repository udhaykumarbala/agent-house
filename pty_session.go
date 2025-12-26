package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/creack/pty"
)

// PTYSession wraps a PTY-based Claude session
type PTYSession struct {
	cmd    *exec.Cmd
	ptmx   *os.File
	output bytes.Buffer
	mu     sync.Mutex
	done   chan struct{}
}

// NewInteractiveSession starts Claude in interactive mode
func NewInteractiveSession() (*PTYSession, error) {
	cmd := exec.Command("claude")

	// Set terminal size
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: 40,
		Cols: 120,
	})
	if err != nil {
		return nil, fmt.Errorf("pty start: %w", err)
	}

	session := &PTYSession{
		cmd:  cmd,
		ptmx: ptmx,
		done: make(chan struct{}),
	}

	// Start reading output in background
	go session.readOutput()

	return session, nil
}

// readOutput continuously reads from PTY
func (s *PTYSession) readOutput() {
	defer close(s.done)

	buf := make([]byte, 4096)
	for {
		n, err := s.ptmx.Read(buf)
		if n > 0 {
			s.mu.Lock()
			s.output.Write(buf[:n])
			s.mu.Unlock()

			// Print to console for debugging
			fmt.Print(string(buf[:n]))
		}
		if err != nil {
			if err != io.EOF {
				fmt.Printf("\n[Read error: %v]\n", err)
			}
			return
		}
	}
}

// SendInput writes text to the PTY
func (s *PTYSession) SendInput(text string) error {
	_, err := s.ptmx.Write([]byte(text))
	return err
}

// SendLine sends text followed by Enter
func (s *PTYSession) SendLine(text string) error {
	return s.SendInput(text + "\n")
}

// SendKey sends a special key (like Escape, Ctrl+C)
func (s *PTYSession) SendKey(key byte) error {
	_, err := s.ptmx.Write([]byte{key})
	return err
}

// GetOutput returns the captured output so far
func (s *PTYSession) GetOutput() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.output.String()
}

// WaitForPattern waits until output contains pattern or timeout
func (s *PTYSession) WaitForPattern(pattern string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if bytes.Contains(s.output.Bytes(), []byte(pattern)) {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

// Wait waits for a duration
func (s *PTYSession) Wait(d time.Duration) {
	time.Sleep(d)
}

// Close terminates the session
func (s *PTYSession) Close() error {
	// Send Ctrl+C then Ctrl+D to exit
	s.SendKey(3) // Ctrl+C
	time.Sleep(200 * time.Millisecond)
	s.SendKey(4) // Ctrl+D
	time.Sleep(200 * time.Millisecond)

	s.ptmx.Close()
	s.cmd.Process.Kill()
	s.cmd.Wait()
	return nil
}

// WaitForExit waits for the session to finish
func (s *PTYSession) WaitForExit(timeout time.Duration) error {
	select {
	case <-s.done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for exit")
	}
}
