package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/OnionRings-dev/git-pour-now/internal/selfupdate"
	"github.com/spf13/cobra"
)

func newSelfUpdateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "self-update",
		Short: "Check for updates or update the git-pour executable",
		Long: `Check the latest Go module version, or build and install it over this executable.
Requires Go 1.23.6 or newer and network access. Installation supports macOS and Linux.
Works outside Git repositories and does not change repository configuration.`,
		Example: "  git pour self-update --check\n  git pour self-update",
		Args:    cobra.NoArgs,
		// Updating the CLI must never activate a repository's shared config.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
		RunE: func(cmd *cobra.Command, args []string) error {
			checkOnly, _ := cmd.Flags().GetBool("check")
			checkCtx, cancel := context.WithTimeout(cmd.Context(), time.Minute)
			defer cancel()
			status, err := selfupdate.Check(checkCtx, installedModuleVersion())
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Installed: %s\nLatest: %s\n", status.Current, status.Latest)
			switch {
			case !status.Known:
				fmt.Fprintln(out, "The installed development build cannot be compared with the latest module version.")
			case !status.Available:
				fmt.Fprintln(out, "Already up to date (installed version is equal to or newer than latest).")
				return nil
			default:
				fmt.Fprintln(out, "An update is available.")
			}
			if checkOnly {
				fmt.Fprintln(out, "Run 'git pour self-update' to install the latest module version.")
				return nil
			}
			executable, err := os.Executable()
			if err != nil {
				return fmt.Errorf("locate executable: %w", err)
			}
			fmt.Fprintf(out, "Installing %s...\n", status.Latest)
			buildCtx, cancelBuild := context.WithTimeout(cmd.Context(), 10*time.Minute)
			defer cancelBuild()
			if err := selfupdate.Install(buildCtx, status.Latest, executable, cmd.ErrOrStderr()); err != nil {
				return err
			}
			fmt.Fprintf(out, "Updated git-pour to %s.\n", status.Latest)
			return nil
		},
	}
	command.Flags().BoolP("check", "c", false, "Check for updates without installing")
	return command
}

// Go-installed binaries retain their module tag or pseudo-version even when
// release ldflags were not supplied. Local untagged builds still report dev.
func installedModuleVersion() string {
	if Version != "dev" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Path == selfupdate.Module && info.Main.Version != "(devel)" && info.Main.Version != "" {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.modified" && setting.Value == "true" {
				return Version
			}
		}
		return info.Main.Version
	}
	return Version
}

func init() {
	rootCmd.AddCommand(newSelfUpdateCommand())
}
