package cmd

import (
	"fmt"
	"strings"

	"github.com/OnionRings-dev/git-pour-now/internal/config"
	"github.com/OnionRings-dev/git-pour-now/internal/git"
	"github.com/OnionRings-dev/git-pour-now/internal/semver"
	"github.com/spf13/cobra"
)

func newReleaseCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "release <version>",
		Short:   "Create a SemVer tag on main or a prerelease branch",
		Long:    "Create an annotated v<version> tag. Stable releases require the configured trunk;\npre-releases require prerelease/<MAJOR.MINOR.PATCH>.\nThe working tree must be clean. Existing release versions cannot be reused.\nVersions must increase according to SemVer precedence. Tags are not pushed.",
		Example: "  git pour release 0.1.0\n  git pour release 1.0.0-rc.1",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			value := strings.TrimPrefix(args[0], "v")
			if !semver.Valid(value) {
				return fmt.Errorf("invalid SemVer version %q; expected MAJOR.MINOR.PATCH with optional pre-release and build metadata", args[0])
			}
			repo, err := openRepo()
			if err != nil {
				return err
			}
			initialized, err := config.IsInitialized(repo)
			if err != nil {
				return err
			}
			if !initialized {
				return fmt.Errorf("git-pour-now is not initialized; run 'git pour init --defaults'")
			}
			cfg, err := config.Load(repo)
			if err != nil {
				return err
			}
			branch, err := repo.GetCurrentBranch()
			if err != nil {
				return err
			}
			trunk := cfg.TrunkBranch()
			target := trunk
			core := strings.SplitN(strings.SplitN(value, "+", 2)[0], "-", 2)
			if len(core) == 2 {
				pre, ok := cfg.Branches["prerelease"]
				if !ok || pre.Type != string(config.BranchTypeTopic) || pre.Parent != trunk {
					return fmt.Errorf("configure a prerelease topic type with the trunk as its parent first")
				}
				target = pre.Prefix + core[0]
			}
			if target == "" || branch != target {
				return fmt.Errorf("version %s must be tagged from %s", value, target)
			}
			if err := repo.RequireCleanReleaseTree(); err != nil {
				return err
			}
			tags, err := repo.ListTags()
			if err != nil {
				return err
			}
			for _, tag := range tags {
				existing := strings.TrimPrefix(tag, "v")
				if semver.Valid(existing) && semver.Compare(value, existing) <= 0 {
					return fmt.Errorf("release %s must be newer than existing tag %s", value, tag)
				}
			}
			tag := "v" + value
			message, _ := cmd.Flags().GetString("message")
			if message == "" {
				message = "Release " + tag
			}
			sign, _ := cmd.Flags().GetBool("sign")
			if err := repo.CreateTag(tag, &git.TagOptions{Message: message, Sign: sign, FailIfExists: true}); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created release %s on %s\n", tag, branch)
			return nil
		},
	}
	command.Flags().StringP("message", "m", "", "Annotated tag message")
	command.Flags().BoolP("sign", "s", false, "Sign the release tag")
	return command
}
