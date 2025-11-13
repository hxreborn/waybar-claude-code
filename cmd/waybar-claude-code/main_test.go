package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestJSONOutput(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	cmd.Env = append(os.Environ(), "CLAUDE_INTERVAL_SEC=300")

	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start binary: %v", err)
	}

	time.Sleep(2 * time.Second)
	cmd.Process.Kill()

	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		t.Fatal("No output from binary")
	}

	line := scanner.Bytes()

	var result struct {
		Text    string `json:"text"`
		Tooltip string `json:"tooltip"`
	}

	if err := json.Unmarshal(line, &result); err != nil {
		t.Fatalf("Invalid JSON output: %v\nOutput: %s", err, line)
	}

	if result.Text == "" {
		t.Error("Expected non-empty text field")
	}

	if result.Tooltip == "" {
		t.Error("Expected non-empty tooltip field")
	}

	t.Logf("Valid JSON output: %s", line)
}
