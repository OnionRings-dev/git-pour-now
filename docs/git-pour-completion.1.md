# GIT-POUR-COMPLETION(1)

## NAME

git-pour-completion - Generate shell completion scripts

## SYNOPSIS

**git-pour completion** *shell*

## DESCRIPTION

Generate shell completion scripts for git-pour. For bash, zsh, and fish the generated scripts provide tab completion for both **git-pour** (direct invocation) and **git pour** (git subcommand) command forms. PowerShell completion supports **git-pour** only.

Supported shells: **bash**, **zsh**, **fish**, **powershell**.

## SHELL-SPECIFIC BEHAVIOUR

### Bash

The generated script includes Cobra's standard completion for the **git-pour** binary plus a `_git_pour()` bridge function. Git's bash completion system calls `_git_<subcommand>()` when completing **git \<subcommand\>**, so this bridge remaps the completion context and delegates to Cobra's `__start_git-pour`.

### Zsh

Zsh's built-in **_git** completion auto-discovers `_git-<subcommand>` functions in **fpath**, so both **git-pour** and **git pour** invocations reach the generated `_git-pour` function. When **_git** dispatches for **git pour**, it shifts the command words so that `words[1]` is `flow` rather than `git-pour`. Cobra's generated body uses `words[1]` to invoke the binary, which would fail. The generated script includes a bridge that rewrites `words[1]` back to `git-pour` before Cobra's logic runs.

### Fish

The generated script includes Cobra's standard completion for **git-pour** plus a bridge that registers completions for the **git** command when the subcommand is **flow**. Deeper completions are delegated to git-pour's built-in `__complete` mechanism.

### PowerShell

Uses Cobra's standard PowerShell completion for the **git-pour** binary. No subcommand bridge is provided, so **git pour** completion is not supported in PowerShell.

## INSTALLATION

### Bash

```bash
# Load in the current session
source <(git pour completion bash)

# Persist for all sessions (Linux)
git pour completion bash > /etc/bash_completion.d/git-pour

# Persist for all sessions (macOS with Homebrew)
git pour completion bash > $(brew --prefix)/etc/bash_completion.d/git-pour
```

### Zsh

```bash
# Enable completion if not already configured
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Install the completion function
git pour completion zsh > "${fpath[1]}/_git-pour"

# Restart your shell
```

### Fish

```bash
# Load in the current session
git pour completion fish | source

# Persist for all sessions
git pour completion fish > ~/.config/fish/completions/git-pour.fish
```

### PowerShell

```powershell
# Load in the current session
git pour completion powershell | Out-String | Invoke-Expression

# Persist: add the above to your PowerShell profile
```

## EXAMPLES

### Generate and preview bash completion

```bash
git pour completion bash
```

### Install bash completion system-wide

```bash
sudo git pour completion bash > /etc/bash_completion.d/git-pour
```

### Install zsh completion for the current user

```bash
mkdir -p ~/.zsh/completions
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
git pour completion zsh > ~/.zsh/completions/_git-pour
```

## EXIT STATUS

**0**
: Completion script generated successfully

## SEE ALSO

**git-pour**(1), **bash**(1), **zshcompsys**(1), **fish**(1)
