package cmd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/OnionRings-dev/git-pour-now/test/testutil"
)

// TestMain builds the git-pour binary before running the command tests so
// they never execute a stale or missing binary. Set GIT_POUR_PATH to test a
// specific prebuilt binary instead.
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
