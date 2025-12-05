package runtime

import (
	"os"
	"os/exec"
	"testing"
)

func TestRuntimeJS(t *testing.T) {
	// Check if bun is installed
	bunPath, err := exec.LookPath("bun")
	if err != nil {
		t.Skip("bun not found in PATH, skipping JS runtime tests")
	}

	cmd := exec.Command(bunPath, "test", "chan.test.js", "stdlib.test.js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("bun test failed: %v", err)
	}
}
