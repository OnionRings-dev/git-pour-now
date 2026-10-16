// Validate release tags for the build and CI using the same rules as the CLI.
package main

import (
	"fmt"
	"github.com/OnionRings-dev/git-pour-now/internal/semver"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "expected one SemVer version or tag")
		os.Exit(1)
	}
	value := strings.TrimPrefix(os.Args[1], "v")
	if !semver.Valid(value) {
		fmt.Fprintf(os.Stderr, "invalid SemVer release: %q\n", os.Args[1])
		os.Exit(1)
	}
	fmt.Println(value)
}
