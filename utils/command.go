package utils

import (
    "os"
    "os/exec"
)

// RunCommand executes a system command and streams its output.
func RunCommand(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}