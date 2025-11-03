package utils

import (
	"runtime"
	"testing"
)

func TestRunCommand(t *testing.T) {
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/C", "echo", "test"}
	} else {
		cmd = "echo"
		args = []string{"test"}
	}

	err := RunCommand(cmd, args...)
	if err != nil {
		t.Errorf("RunCommand() error = %v", err)
	}
}

func TestRunCommand_InvalidCommand(t *testing.T) {
	err := RunCommand("nonexistent-command-xyz-12345")
	if err == nil {
		t.Error("RunCommand() expected error for invalid command, got nil")
	}
}

func TestRunCommand_WithMultipleArgs(t *testing.T) {
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/C", "echo", "hello", "world"}
	} else {
		cmd = "echo"
		args = []string{"hello", "world"}
	}

	err := RunCommand(cmd, args...)
	if err != nil {
		t.Errorf("RunCommand() with multiple args error = %v", err)
	}
}
