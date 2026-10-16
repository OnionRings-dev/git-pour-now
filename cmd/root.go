package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-pour",
	Short: "git-pour-now manages features on main and releases through SemVer tags",
	Long: `git-pour-now keeps development on main, with short-lived feature branches.
Prepare prereleases on temporary prerelease/<version> branches.
Release a committed main revision by creating an immutable SemVer tag.`,
	Example: `  git pour init --defaults
  git pour feature start my-feature
  git pour feature finish my-feature
  git pour prerelease start 0.1.0
  git pour release 0.1.0-rc.1
  git pour prerelease finish 0.1.0
  git pour release 0.1.0`,
	// First-run activation of a committable .gitpour file runs before any
	// command. It is a no-op unless the repo is unconfigured and a valid .gitpour
	// is present (see firstRunActivation).
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return firstRunActivation(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, print help
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
