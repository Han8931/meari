#!/usr/bin/env bash
#
# Meari install script: builds the binary and puts it on your PATH.
#
#   ./install.sh                 # install to ~/.local/bin
#   ./install.sh --clean         # also remove old binaries and regenerable
#                                # app state (learning progress is kept)
#   BIN_DIR=/usr/local/bin ./install.sh
#
# Meari roots its files (config.toml, vault/, data/, workspace/, ...) at the
# directory you run it from, so after installing, launch it from the folder
# you want as your Meari home.
#
# --clean never touches: config.toml, vault/ (your notes), data/ (progress),
# workspace/drafts/ (in-progress challenge solutions), meari-course/
# (progress references these), meari-publish/ (your sharing repo).

set -euo pipefail

BIN_DIR="${BIN_DIR:-$HOME/.local/bin}"
GO_MIN_MINOR=26 # requires go 1.26+

repo_dir="$(cd "$(dirname "$0")" && pwd)"

clean=false
case "${1:-}" in
    --clean) clean=true ;;
    '') ;;
    *) printf 'usage: %s [--clean]\n' "$0" >&2; exit 2 ;;
esac

info() { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33mnote:\033[0m %s\n' "$*"; }
fail() { printf '\033[1;31merror:\033[0m %s\n' "$*" >&2; exit 1; }

# --- check Go toolchain ------------------------------------------------------

command -v go >/dev/null 2>&1 || fail "Go is not installed. Get it from https://go.dev/dl/ (need 1.$GO_MIN_MINOR+)."

go_version="$(go env GOVERSION)" # e.g. go1.26.1
go_minor="$(printf '%s' "$go_version" | sed -E 's/^go1\.([0-9]+).*/\1/')"
case "$go_minor" in
    ''|*[!0-9]*) fail "could not parse Go version from '$go_version'" ;;
esac
if [ "$go_minor" -lt "$GO_MIN_MINOR" ]; then
    fail "Go 1.$GO_MIN_MINOR+ required, found $go_version."
fi

# --- build -------------------------------------------------------------------

# Build to a temp location so no stale ./meari artifact is left in the repo
# to shadow the installed binary.
build_dir="$(mktemp -d)"
trap 'rm -rf "$build_dir"' EXIT

info "Building meari ($go_version)"
(cd "$repo_dir" && go build -trimpath -ldflags="-s -w" -o "$build_dir/meari" .)

# --- install -----------------------------------------------------------------

info "Installing to $BIN_DIR/meari"
mkdir -p "$BIN_DIR"
install -m 0755 "$build_dir/meari" "$BIN_DIR/meari"

case ":$PATH:" in
    *":$BIN_DIR:"*) ;;
    *)
        warn "$BIN_DIR is not on your PATH. Add this to your shell profile:"
        printf '      export PATH="%s:$PATH"\n' "$BIN_DIR"
        ;;
esac

# --- clean: previous binaries + regenerable app state ------------------------

if $clean; then
    # Remove every other meari binary on PATH so nothing stale can shadow
    # the one just installed.
    removed_old=false
    IFS=':' read -ra path_dirs <<<"$PATH"
    for dir in "${path_dirs[@]}"; do
        [ -n "$dir" ] || continue
        candidate="$dir/meari"
        [ -f "$candidate" ] && [ -x "$candidate" ] || continue
        [ "$candidate" -ef "$BIN_DIR/meari" ] && continue
        if rm -f "$candidate" 2>/dev/null; then
            info "Removed old binary: $candidate"
            removed_old=true
        else
            warn "could not remove $candidate — run: sudo rm \"$candidate\""
        fi
    done
    if $removed_old; then
        warn "your shell may have cached the old location — run: hash -r"
    fi

    # The build artifact from `go build -o meari .` shadows the install when
    # run as ./meari.
    if [ -f "$repo_dir/meari" ]; then
        rm -f "$repo_dir/meari"
        info "Removed build artifact: $repo_dir/meari"
    fi

    # Regenerable app state. Progress is kept: data/ (progress.json) and
    # workspace/drafts/ (in-progress challenge solutions survive).
    if [ -d "$repo_dir/workspace" ]; then
        find "$repo_dir/workspace" -mindepth 1 -maxdepth 1 ! -name drafts \
            -exec rm -rf {} +
        info "Cleared workspace/ (kept workspace/drafts/)"
    fi
    if [ -d "$repo_dir/exports" ]; then
        find "$repo_dir/exports" -mindepth 1 -maxdepth 1 -exec rm -rf {} +
        info "Cleared exports/"
    fi
fi

# --- done --------------------------------------------------------------------

cat <<EOF

Installed. Meari keeps its notes and data in whatever directory you run it
from, so pick (or create) a home for it first:

    mkdir -p ~/meari && cd ~/meari
    cp "$repo_dir/config.example.toml" config.toml   # optional; runs fine without

    meari -vault    # the vault, in your terminal
    meari serve     # the same vault, in your browser
    meari           # the tutor (launch dashboard)
    meari check     # verify config / AI provider setup

EOF
