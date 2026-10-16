#!/bin/bash

# Abort on the first error so a failed build never reaches the packaging step,
# which would otherwise archive an incomplete set of binaries.
set -e

# Release versions come from an exact tag; untagged builds identify as dev.
VERSION=${1:-$(git describe --tags --exact-match 2>/dev/null || echo dev)}
if [[ "$VERSION" != "dev" ]]; then
    VERSION=$(go run ./scripts/release-version "$VERSION")
fi
BINARY_NAME="git-pour"
PACKAGE_NAME="git-pour-now"

# Build directory
BUILD_DIR="dist"

# Per-platform staging directories, shared with package.sh.
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/platforms.sh"

# Get build information
GIT_COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')

# Build flags
# -s: Strip symbol table
# -w: Strip DWARF debug info
# Combined with -trimpath and CGO_ENABLED=0 for minimal binary size
BUILD_FLAGS="-s -w -X 'github.com/OnionRings-dev/git-pour-now/version.Version=${VERSION}' -X 'github.com/OnionRings-dev/git-pour-now/cmd.BuildDate=${BUILD_TIME}' -X 'github.com/OnionRings-dev/git-pour-now/cmd.GitCommit=${GIT_COMMIT}' -X 'github.com/OnionRings-dev/git-pour-now/version.BuildTime=${BUILD_TIME}' -X 'github.com/OnionRings-dev/git-pour-now/version.GitCommit=${GIT_COMMIT}'"

# Start from a clean build directory so stale archives from an earlier run
# never leak into the checksums and re-archiving can't retain an obsolete
# binary entry. Then create the per-platform staging directories.
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"
for dir in "${PLATFORM_DIRS[@]}"; do
    mkdir -p "$BUILD_DIR/$dir"
done

# Build for each platform/architecture
echo "Building $PACKAGE_NAME version $VERSION..."

# macOS (both Intel and Apple Silicon)
echo "Building darwin/amd64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "${BUILD_FLAGS}" -o "$BUILD_DIR/darwin-amd64/${BINARY_NAME}" main.go
echo "Building darwin/arm64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "${BUILD_FLAGS}" -o "$BUILD_DIR/darwin-arm64/${BINARY_NAME}" main.go

# Linux
echo "Building linux/amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "${BUILD_FLAGS}" -o "$BUILD_DIR/linux-amd64/${BINARY_NAME}" main.go
echo "Building linux/arm64..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "${BUILD_FLAGS}" -o "$BUILD_DIR/linux-arm64/${BINARY_NAME}" main.go
echo "Building linux/386..."
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -trimpath -ldflags "${BUILD_FLAGS}" -o "$BUILD_DIR/linux-386/${BINARY_NAME}" main.go

# Windows
echo "Building windows/amd64..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "${BUILD_FLAGS}" -o "$BUILD_DIR/windows-amd64/${BINARY_NAME}.exe" main.go
echo "Building windows/386..."
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -trimpath -ldflags "${BUILD_FLAGS}" -o "$BUILD_DIR/windows-386/${BINARY_NAME}.exe" main.go
echo "Building windows/arm64..."
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -trimpath -ldflags "${BUILD_FLAGS}" -o "$BUILD_DIR/windows-arm64/${BINARY_NAME}.exe" main.go

# Verify all binaries were created
echo "Verifying binaries..."
MISSING_BINARIES=()

for dir in "${PLATFORM_DIRS[@]}"; do
    binary="$(platform_binary "$dir" "$BINARY_NAME")"
    if [[ ! -f "$BUILD_DIR/$binary" ]]; then
        MISSING_BINARIES+=("$binary")
    fi
done

if [[ ${#MISSING_BINARIES[@]} -gt 0 ]]; then
    echo "Error: The following binaries were not created:"
    for binary in "${MISSING_BINARIES[@]}"; do
        echo "  - $binary"
    done
    exit 1
fi

# Archiving and checksums live in package.sh so the release workflow can sign
# the Windows binaries in between. SKIP_PACKAGE=1 leaves the staging
# directories in place for that signing step; a plain build still packages, so
# a local run produces the same artifacts as before.
if [[ "${SKIP_PACKAGE:-}" == "1" ]]; then
    echo "Build complete! Binaries are in the $BUILD_DIR directory"
else
    "$SCRIPT_DIR/package.sh" "$VERSION"
fi
