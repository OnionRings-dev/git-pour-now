package worktree_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/OnionRings-dev/git-pour-now/test/testutil"
)

// TestMain builds the git-pour binary before running the worktree path tests,
// which initialize git-pour via the binary to obtain a realistic configuration.
// Set GIT_POUR_PATH to test a specific prebuilt binary instead.
func TestMain(m *testing.M) {
	cleanup, err := testutil.BuildGitFlow()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build git-pour binary: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()
	cleanup()
	os.Exit(code)
}
