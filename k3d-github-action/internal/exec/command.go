package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// RunCommand executes a command with arguments and returns its output
func RunCommand(ctx context.Context, command string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()
	if err != nil {
		if stderr.Len() > 0 {
			return output, fmt.Errorf("%w: %s", err, stderr.String())
		}
		return output, err
	}

	return output, nil
}

// CommandExists checks if a command is available in PATH
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
