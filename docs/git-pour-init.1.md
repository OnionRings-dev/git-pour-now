# GIT-POUR-INIT(1)

## NAME

git-pour-init - Initialize git-pour in a repository

## SYNOPSIS

**git-pour init** [**-f**|**--force**] [**--init**] [**--preset**=*preset*] [**--custom**] [**--defaults**] [**--shared**|**--local**|**--global**|**--system**|**--file**=*path*] [*options*]

## DESCRIPTION

Initialize git-pour configuration in the current Git repository. The default workflow uses main, feature/, and temporary prerelease/ branches. Stable releases are SemVer tags on main; prereleases are tagged on the matching prerelease/<version> branch; no develop branch is created. Use --preset=classic explicitly for the traditional workflow.

**git-pour init** supports three initialization modes:

1. **Interactive Mode** (default) - Prompts for branch names and prefixes one question at a time, each showing its default in brackets that Enter accepts
2. **Preset Mode** - Automatically applies a predefined workflow configuration
3. **Custom Mode** - Sets up only the trunk branch and shows configuration commands

## OPTIONS

### General Options

**-f**, **--force**
: Force reconfiguration of git-pour even if already initialized. Without this option, **git pour init** will fail if configuration already exists (in non-interactive mode) or prompt for confirmation (in interactive mode).

**--init**
: Create a git repository in the current directory when there is none, then initialize git-pour in it. Without this option, running **git pour init** outside a repository fails (exit status **3**) when stdin is not an interactive terminal, and prompts **No git repository here. Create one? [y/N]** when it is; declining creates nothing. Inside an existing repository **--init** is a no-op. The created repository's initial branch is the resolved git-pour trunk (for example **main** with **--defaults**, or **trunk** with **--main trunk**), overriding any ambient **init.defaultBranch**. **--init** governs only repository creation — it does not change how configuration is selected, and it applies with every configuration scope option, including **--global**, **--system**, **--file** and **--shared**.

### Preset Options

**--preset**=*preset*, **-p** *preset*
: Apply a predefined workflow preset. Valid values: **classic**, **github**, **gitlab**

**--custom**
: Enable custom configuration mode. Prompts for trunk branch and displays configuration commands.

**--defaults**, **-d**
: Configure main and feature/ without prompting. No develop, release, or hotfix branch is created.

**--no-create-branches**
: Don't create branches even if they don't exist in the repository.

### Configuration Scope Options

These options control where git-pour configuration is stored. Only one scope option may be specified at a time. When no scope option is given, git-pour reads from merged config (local > global > system precedence) and writes to local config.

**--shared**
: Author the configuration into a committable **.gitpour** file at the repository top level, then copy the **gitpour.*** keys into the repository's local **.git/config**. Committing **.gitpour** lets teammates share one git-pour configuration: on a fresh clone, git-pour offers to activate it (see **gitpour-config**(5), FIRST-RUN ACTIVATION). **--shared** is mutually exclusive with **--local**, **--global**, **--system**, and **--file**. Without **--force**, running **--shared** again when a **.gitpour** already exists fails and leaves the file untouched; with **--force** the file is rewritten and re-copied (removing any stale managed keys from local config).

**--local**
: Read and write configuration only in the repository's **.git/config** file. This is the default for writes when no scope option is specified.

**--global**
: Read and write configuration in the user's global **~/.gitconfig** file. Useful for setting up defaults that apply to all repositories.

**--system**
: Read and write configuration in the system-wide **/etc/gitconfig** file. Typically requires administrator privileges.

**--file**=*path*
: Read and write configuration in the specified file. The parent directory must exist and be writable. Paths may be absolute or relative to the current working directory. Useful for managing shared configuration files.

### Branch Name Overrides

**--main**=*name*, **-m** *name*
: Override main branch name (default: main)

**--develop**=*name*, **-e** *name*
: Classic preset only: override develop branch name (default: develop)

### Prefix Overrides

**--feature**=*prefix*
: Override feature branch prefix (default: feature/)

**--bugfix**=*prefix*, **-b** *prefix*
: Override bugfix branch prefix (default: bugfix/)

**--release**=*prefix*, **-r** *prefix*
: Override release branch prefix (default: release/)

**--hotfix**=*prefix*, **-x** *prefix*
: Override hotfix branch prefix (default: hotfix/)

**--support**=*prefix*, **-s** *prefix*
: Override support branch prefix (default: support/)

**--tag**=*prefix*, **-t** *prefix*
: Override version tag prefix (default: none)

A prefix is stored exactly as given; no separator is appended. The trailing slash in the defaults is a naming convention, not a requirement — `feature/`, `feature_` and `feature-` are all valid and produce the branches `feature/login`, `feature_login` and `feature-login` respectively.

## PRESETS

### Classic GitFlow

Traditional git-pour workflow with the following structure:

- **main** - Production releases (trunk)
- **develop** - Integration branch (auto-updates from main)
- **feature/** - New features (parent: develop)
- **release/** - Release preparation (parent: main, starts from develop, creates tags)
- **hotfix/** - Emergency fixes (parent: main, creates tags)
- **support/** - Long-term support (parent: main)

### GitHub Flow

Simplified workflow for continuous deployment:

- **main** - Production branch (trunk)
- **feature/** - All development work (parent: main)

### GitLab Flow

Multi-environment workflow for staged deployments:

- **production** - Production environment (trunk)
- **staging** - Staging environment (parent: production)
- **main** - Development integration (parent: staging)
- **feature/** - Development work (parent: main)
- **hotfix/** - Production fixes (parent: production)

## INTERACTIVE MODE

When run without options, **git-pour init** asks:

```
Main branch [main]:
Feature branch prefix [feature/]:
```

Answers are trimmed of surrounding whitespace and then stored verbatim. No separator is appended to a prefix, so answering the feature prompt with `feature_` stores exactly that prefix, and `git pour feature start login` then creates the branch `feature_login`.

An answer that is empty after trimming leaves the default in place.

If a git-flow-avh configuration is already present, **git-pour init** imports it instead of prompting.

## CUSTOM MODE

With **--custom**, only prompts for the trunk branch:

```
? What's your trunk branch (holds production code)? [main] production
✓ Trunk branch: production

Configuration commands:
  git-pour config add base <name> [<parent>] [options...]
  git-pour config add topic <name> <parent> [options...]
  [... full command reference displayed ...]
```

## EXAMPLES

Initialize with Classic GitFlow using defaults:
```bash
git pour init --preset=classic
```

Initialize with defaults without prompting:
```bash
git pour init --defaults
```

Initialize with preset and defaults:
```bash
git pour init --preset=classic --defaults
```

Initialize GitHub Flow with custom main branch:
```bash
git pour init --preset=github --main=master
```

Initialize Classic GitFlow with custom branch names:
```bash
git pour init --preset=classic --main=master --develop=dev --feature=feat/
```

Initialize with short flags:
```bash
git pour init -p classic -d -m master -b bug/ -r rel/
```

Custom configuration mode:
```bash
git pour init --custom
```

Interactive initialization:
```bash
git pour init
```

Reconfigure git-pour with new settings:
```bash
git pour init --force --feature=feat/
```

Force reinitialize with github preset:
```bash
git pour init --preset=github --force
```

Reconfigure with short flag:
```bash
git pour init -f --defaults
```

Initialize with global scope (user-wide defaults):
```bash
git pour init --defaults --global
```

Initialize with local scope (repository-specific):
```bash
git pour init --defaults --local
```

Initialize with configuration file:
```bash
git pour init --defaults --file=/path/to/custom-gitpour.config
```

Initialize a shared, committable configuration:
```bash
git pour init --defaults --shared
git add .gitpour && git commit -m "Add shared git-pour configuration"
```

Create local config when global config already exists:
```bash
git pour init --defaults --local
```

## CONFIGURATION

By default, git-pour stores configuration in the repository's **.git/config** file under the **gitpour.*** namespace. The **--global**, **--system**, or **--file** options can be used to store configuration in alternate locations.

```
[gitpour]
    version = 1.0
    initialized = true
[gitpour "branch.main"]
    type = base
[gitpour "branch.develop"]
    type = base
    parent = main
    autoupdate = true
[gitpour "branch.feature"]
    type = topic
    parent = develop
    prefix = feature/
```

## EXIT STATUS

**0**
: Successful initialization

**1**
: Usage error — an unknown option or an unexpected argument

**2**
: Invalid options — for example multiple configuration scope options, or an invalid branch name

**3**
: Git operation failed (for example the current directory is not a git repository and repository creation was not authorized — no **--init**, and either no interactive terminal or the prompt was declined)

**6**
: A required precondition failed — git-pour is already initialized and **--force** was not given, or the repository has no commits and branch creation is requested but no git identity (**user.name** and **user.email**) is configured

## SEE ALSO

**git-pour**(1), **git-pour-config**(1), **gitpour-config**(5)

## NOTES

- Outside a git repository, **git-pour init** never creates one implicitly; creation requires **--init** or an affirmative answer to the interactive prompt. This preserves the safety of typing **git pour init** in the wrong empty directory
- When **--init** creates the repository, the identity precondition described below still applies. The repository is created before the check runs, so a failure leaves an initialized-but-unconfigured repository behind; configure **user.name**/**user.email** and re-run, which then takes the ordinary "already a repository" path
- **git-pour init** requires **--force** to reconfigure an already initialized repository in non-interactive mode
- In interactive mode without **--force**, users are prompted for confirmation before reconfiguring
- Existing branches are preserved during initialization
- When initializing a repository with no existing commits, **git-pour init** creates an empty initial commit to enable branch creation. No files are added to the working directory
- Creating that initial commit requires a configured git identity. When the repository has no commits and branch creation is requested, **git-pour init** verifies that both **user.name** and **user.email** are set (in local, global, or system config) before writing any configuration. If the identity is missing, init fails fast with an actionable error (exit status **6**) and writes no configuration, so it can be re-run after configuring the identity. A repository that **--init** created moments earlier is not rolled back and stays in place, as described above. Repositories that already have commits, or runs with **--no-create-branches**, do not require an identity
- Compatible with repositories previously initialized with git-flow-avh
- Configuration scope options (**--local**, **--global**, **--system**, **--file**) only affect the **init** command. All other git-pour commands (start, finish, update, etc.) always read from merged config using Git's standard precedence (local > global > system)
- When checking initialization status without an explicit scope flag, git-pour checks merged config and reports which scope the configuration was found in
- When initialized via global or system config, attempting to initialize again without a scope flag will suggest using **--local** to create repository-specific config
