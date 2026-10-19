#!/bin/sh

# Install from the Go module without requiring a clone or administrator access.
set -eu

fail() {
    printf 'Error: %s\n' "$*" >&2
    exit 1
}

case "$(uname -s)" in
    Darwin|Linux) ;;
    *) fail 'This installer supports macOS and Linux.' ;;
esac

for dependency in go git; do
    command -v "$dependency" >/dev/null 2>&1 ||
        fail "$dependency is required. Install Go 1.23.6 or newer and Git, then retry."
done

: "${HOME:?HOME must be set}"
bin_dir="$HOME/.local/bin"
shell_path=${SHELL:-sh}
shell_name=${shell_path##*/}

# Use a temporary GOBIN so a failed build cannot replace an existing install.
temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' 0
trap 'exit 1' HUP INT TERM

printf 'Building the latest git-pour-now with Go...\n'
GOBIN="$temp_dir" go install github.com/OnionRings-dev/git-pour-now@latest ||
    fail 'Go installation failed; the existing git-pour installation was not changed.'

mkdir -p "$bin_dir"
install -m 755 "$temp_dir/git-pour-now" "$bin_dir/git-pour"

add_path() {
    profile=$1
    path_line='export PATH="$HOME/.local/bin:$PATH"'
    if [ ! -f "$profile" ] || ! grep -Fqx "$path_line" "$profile"; then
        printf '\n# git-pour\n%s\n' "$path_line" >> "$profile"
        printf 'Added PATH setup to %s\n' "$profile"
    fi
}

case "$shell_name" in
    zsh)
        mkdir -p "${ZDOTDIR:-$HOME}"
        add_path "${ZDOTDIR:-$HOME}/.zshrc"
        ;;
    bash)
        add_path "$HOME/.bashrc"
        # Login shells read the first existing file in this order. Configure
        # that file too, for terminals that start bash as a login shell.
        if [ -f "$HOME/.bash_profile" ]; then
            add_path "$HOME/.bash_profile"
        elif [ -f "$HOME/.bash_login" ]; then
            add_path "$HOME/.bash_login"
        else
            add_path "$HOME/.profile"
        fi
        ;;
    sh|dash|ksh|'') add_path "$HOME/.profile" ;;
    *) printf 'Add %s to PATH in your shell configuration.\n' "$bin_dir" ;;
esac

printf '\nInstalled %s\n' "$bin_dir/git-pour"
"$bin_dir/git-pour" version
printf '\nTo use it in this terminal, run:\n  export PATH="$HOME/.local/bin:$PATH"\n  git pour version\n'
