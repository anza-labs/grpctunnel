package internal_test

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestTunnelIntegration(t *testing.T) {
	// Define test cases with different flags for server and client
	testCases := []struct {
		name        string
		serverFlags []string
		clientFlags []string
	}{
		{
			name:        "Default configuration",
			serverFlags: []string{},
			clientFlags: []string{},
		},
	}

	// Build the server and client binaries
	serverBinary := filepath.Join(t.TempDir(), "tunneltestsvr")
	clientBinary := filepath.Join(t.TempDir(), "tunneltestclient")

	buildBinary(t, "./cmd/tunneltestsvr", serverBinary)
	buildBinary(t, "./cmd/tunneltestclient", clientBinary)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a context with timeout for the test
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()

			// Start the server with flags
			serverCmd := exec.CommandContext(ctx, serverBinary, tc.serverFlags...)
			serverCmd.Stdout = t.Output()
			serverCmd.Stderr = t.Output()

			if err := serverCmd.Start(); err != nil {
				t.Fatalf("Failed to start server: %v", err)
			}

			// Ensure the server process is cleaned up
			t.Cleanup(func() {
				if err := serverCmd.Process.Kill(); err != nil {
					t.Logf("Failed to kill server process: %v", err)
				}
			})

			// Give the server some time to start
			time.Sleep(2 * time.Second)

			// Start the client with flags
			clientCmd := exec.CommandContext(ctx, clientBinary, tc.clientFlags...)
			clientCmd.Stdout = t.Output()
			clientCmd.Stderr = t.Output()

			if err := clientCmd.Start(); err != nil {
				t.Fatalf("Failed to start client: %v", err)
			}

			// Ensure the client process is cleaned up
			t.Cleanup(func() {
				if err := clientCmd.Process.Kill(); err != nil {
					t.Logf("Failed to kill client process: %v", err)
				}
			})

			// Wait for the client to finish
			clientErr := clientCmd.Wait()
			if clientErr != nil {
				t.Errorf("Client process exited with error: %v", clientErr)
			}

			// Since the server never finishes, we log its termination but ignore errors
			if err := serverCmd.Process.Kill(); err != nil {
				t.Logf("Failed to kill server process: %v", err)
			}
		})
	}
}

// buildBinary builds the Go binary from the given source directory and outputs it to the specified path.
func buildBinary(t *testing.T, sourceDir, outputPath string) {
	t.Helper()

	cmd := exec.CommandContext(t.Context(), "go", "build", "-o", outputPath, sourceDir)
	cmd.Stderr = t.Output()

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary for %s: %v", sourceDir, err)
	}
}
