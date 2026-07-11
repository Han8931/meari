#!/usr/bin/env bash
#
# Meari install script: builds the binary and puts it on your PATH.
#
#   ./install.sh                 # install to ~/.local/bin; also removes any
#                                # stale meari binaries on your PATH so nothing
#                                # old can shadow the one just installed
#   ./install.sh --keep-old      # install but leave other meari binaries alone
#   ./install.sh --clean         # also clear regenerable app state
#                                # (learning progress is kept)
#   BIN_DIR=/usr/local/bin ./install.sh
#
# Meari roots its files (config.toml, vault/, data/, workspace/, ...) at the
# directory you run it from, so after installing, launch it from the folder
# you want as your Meari home.
#
# Cleanup never touches: config.toml, vault/ (your notes), data/ (progress),
# workspace/drafts/ (in-progress challenge solutions), meari-course/
# (progress references these), meari-publish/ (your sharing repo).

set -euo pipefail

BIN_DIR="${BIN_DIR:-$HOME/.local/bin}"
GO_MIN_MINOR=26 # requires go 1.26+

repo_dir="$(cd "$(dirname "$0")" && pwd)"

clean=false      # --clean: also clear regenerable app state
keep_old=false   # --keep-old: skip removing stale meari binaries
for arg in "$@"; do
    case "$arg" in
        --clean)    clean=true ;;
        --keep-old) keep_old=true ;;
        *) printf 'usage: %s [--clean] [--keep-old]\n' "$0" >&2; exit 2 ;;
    esac
done

info() { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33mnote:\033[0m %s\n' "$*"; }
fail() { printf '\033[1;31merror:\033[0m %s\n' "$*" >&2; exit 1; }

# --- ensure user-owned course storage ---------------------------------------

# Course files are ordinary user data and must never require sudo to edit.
# A directory copied from a container or sandbox can occasionally be readable
# but owned by a different UID (commonly `nobody`). Because the project parent
# is user-owned, repair that state without chown: atomically preserve the old
# tree, create a new user-owned directory, and copy the readable content back.
course_dir="$repo_dir/meari-course"
course_needs_repair=false
if [ -d "$course_dir" ]; then
    while IFS= read -r -d '' course_path; do
        if [ ! -w "$course_path" ]; then
            course_needs_repair=true
            break
        fi
    done < <(find "$course_dir" -print0)
fi
if [ ! -e "$course_dir" ]; then
    mkdir -p "$course_dir"
elif [ ! -d "$course_dir" ]; then
    fail "$course_dir exists but is not a directory"
elif $course_needs_repair; then
    backup="$repo_dir/meari-course.ownership-backup-$(date +%Y%m%d%H%M%S)"
    [ ! -e "$backup" ] || fail "ownership backup already exists: $backup"

    warn "$course_dir is not writable; preserving it at $backup"
    if ! mv "$course_dir" "$backup"; then
        fail "could not preserve the unwritable course directory"
    fi
    if ! mkdir -p "$course_dir" || ! cp -R "$backup"/. "$course_dir"/; then
        rm -rf "$course_dir"
        mv "$backup" "$course_dir" 2>/dev/null || true
        fail "could not create a writable course directory; original restored"
    fi
    info "Repaired course directory ownership without sudo"
    warn "original course files remain at $backup"
fi

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

# --- remove stale binaries so nothing old can shadow the new install ---------

if ! $keep_old; then
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
fi

# --- clean: regenerable app state (opt-in) -----------------------------------

if $clean; then
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
    # optional (Meari runs fine without a config); never overwrites an existing one:
    [ -f config.toml ] || cp "$repo_dir/config.example.toml" config.toml

    meari -vault    # the vault, in your terminal
    meari serve     # the same vault, in your browser
    meari           # the tutor (launch dashboard)
    meari check     # verify config / AI provider setup

EOF
