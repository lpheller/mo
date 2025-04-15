package utils

import (
    "fmt"
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

func RunRemoteCommand(user, host, cmd string) error {
    sshCmd := exec.Command("ssh", fmt.Sprintf("%s@%s", user, host), cmd)
    sshCmd.Stdout = os.Stdout
    sshCmd.Stderr = os.Stderr
    return sshCmd.Run()
}