package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	// Disable Cobra's default completion command; we provide our own
	// that adds git-subcommand integration for "git pour".
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(completionCmd)
}

var completionCmd = &cobra.Command{
	Use:       "completion {bash|zsh|fish|powershell}",
	Short:     "Generate shell completion scripts",
	Long:      completionLong,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return genBashCompletion(os.Stdout)
		case "zsh":
			return genZshCompletion(os.Stdout)
		case "fish":
			return genFishCompletion(os.Stdout)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

const completionLong = `Generate shell completion scripts for git-pour.

For bash, zsh, and fish the generated scripts provide tab completion for
both "git-pour" (direct) and "git pour" (git subcommand) invocation styles.
PowerShell completion supports "git-pour" only.

For zsh, the "git pour" form works because zsh's built-in _git completion
dispatches to any _git-<name> function it finds in fpath.

Bash:

  source <(git pour completion bash)

  # To load completions for each session, execute once:
  # Linux:
  git pour completion bash > /etc/bash_completion.d/git-pour
  # macOS (Homebrew):
  git pour completion bash > $(brew --prefix)/etc/bash_completion.d/git-pour

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it first:
  echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  git pour completion zsh > "${fpath[1]}/_git-pour"

  # You will need to start a new shell for this setup to take effect.

Fish:

  git pour completion fish | source

  # To load completions for each session, execute once:
  git pour completion fish > ~/.config/fish/completions/git-pour.fish

PowerShell:

  git pour completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, add the output of
  # the above command to your PowerShell profile.
`

// genBashCompletion generates bash completion with a bridge for "git pour" usage.
// Git's bash completion calls _git_<subcommand>() when completing "git <subcommand>",
// so we append a _git_pour() function that remaps the completion context and
// delegates to Cobra's __start_git-pour.
func genBashCompletion(w io.Writer) error {
	if err := rootCmd.GenBashCompletionV2(w, true); err != nil {
		return err
	}
	_, err := fmt.Fprint(w, bashBridge)
	return err
}

const bashBridge = `
# Bridge for "git pour" subcommand completion.
# Git's bash completion calls _git_pour() when the user types "git pour <tab>".
# This remaps the completion context so Cobra sees "git-pour ..." and delegates.
_git_pour()
{
    # Rewrite "git pour" to "git-pour" so Cobra recognises the binary name.
    # Both forms are 8 characters, so COMP_POINT stays unchanged.
    COMP_LINE="${COMP_LINE/git pour/git-pour}"
    COMP_WORDS=("git-pour" "${COMP_WORDS[@]:2}")
    ((COMP_CWORD--))

    # Reset compopt to match the defaults Cobra expects when compopt is available
    # ("complete -o default", no nospace). Git's completion context may have set
    # options that interfere with trailing-space behaviour.
    if [[ $(type -t compopt) = "builtin" ]]; then
        compopt -o default +o nospace
    fi

    __start_git-pour
}
`

// genZshCompletion generates zsh completion with a bridge for "git pour" usage.
// Zsh's _git dispatcher calls _git-pour() to complete "git pour <...>" and
// shifts words so words[1] is "pour" (not "git-pour"). Cobra's generated body
// uses "${words[1]} __complete ..." to invoke the binary, so without the
// rewrite it would try to run a nonexistent command called "pour".
func genZshCompletion(w io.Writer) error {
	var buf bytes.Buffer
	if err := rootCmd.GenZshCompletion(&buf); err != nil {
		return err
	}
	const anchor = "_git-pour()\n{\n"
	const replacement = `_git-pour()
{
    # Bridge for "git pour" subcommand: zsh's _git dispatcher shifts words so
    # words[1] is "pour" when this function is invoked via "git pour ...".
    # Rewrite it back to "git-pour" so the Cobra logic below invokes the
    # correct binary via __complete.
    if [[ ${words[1]} == pour ]]; then
        words[1]=git-pour
    fi

`
	output := strings.Replace(buf.String(), anchor, replacement, 1)
	if !strings.Contains(output, "words[1]=git-pour") {
		return fmt.Errorf("zsh completion: failed to inject git-subcommand bridge (Cobra output format may have changed)")
	}
	_, err := fmt.Fprint(w, output)
	return err
}

// genFishCompletion generates fish completion with support for "git pour" usage.
// Fish has no automatic git-subcommand discovery, so we append registrations
// that delegate to git-pour's __complete when the user types "git pour <tab>".
func genFishCompletion(w io.Writer) error {
	if err := rootCmd.GenFishCompletion(w, true); err != nil {
		return err
	}
	_, err := fmt.Fprint(w, fishBridge)
	return err
}

const fishBridge = `
# Bridge for "git pour" subcommand completion.
# Registers completions for "git pour ..." by delegating to git-pour's
# Cobra-generated __complete command.
function __fish_git_pour_complete
    set -l tokens (commandline -opc)
    set -l current (commandline -ct)
    # Skip "git" and "pour", pass the rest to git-pour __complete
    set -l args $tokens[3..-1]

    set -l results (GIT_POUR_ACTIVE_HELP=0 git-pour __complete $args $current 2>/dev/null)
    if test (count $results) -eq 0
        return
    end

    # The last line is the completion directive; strip it
    set -e results[-1]
    string join \n -- $results
end

# Register "pour" as a git subcommand
complete -c git -f -n '__fish_use_subcommand' -a pour -d 'Manage branches with git-pour'

# Delegate deeper completions to git-pour
complete -c git -f -n '__fish_seen_subcommand_from pour' -a '(__fish_git_pour_complete)'
`
