#!/usr/bin/env bash
# Every check this repository has to pass, both languages. Run identically by
# .github/workflows/ci.yml and by a person:
#
#     docker run --rm -v "$PWD:/w" -w /w ontavu/cams-ci:latest \
#         bash -c 'git config --global --add safe.directory /w && ci/check.sh'
#
# In that image rather than with whatever a workstation has: protoc stamps its
# own version into the Go bindings, and the image's protoc is the one that
# generated the committed ones.

set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

log() { printf '\n=== %s\n' "$*"; }
die() { printf 'error: %s\n' "$*" >&2; exit 1; }

for tool in go gofmt git protoc grpc_cpp_plugin cmake pkg-config; do
    command -v "$tool" >/dev/null || die "$tool is not on PATH"
done

# Never a workspace: this module is checked as its consumers fetch it.
export GOWORK=off

scratch="$(mktemp -d)"
trap 'rm -rf "$scratch"' EXIT

# Compared before and after rather than required clean, so that the checks can
# run on a change that is not committed yet. The diff is in it so that a file
# already modified and modified again does not go unnoticed.
tree_state() { git status --porcelain | LC_ALL=C sort; git diff; }
tree_before="$(tree_state)"

log "toolchain"
go version
protoc --version

log "go mod download"
go mod download

# gofmt -l exits 0 whether or not it found anything.
log "gofmt"
unformatted="$(gofmt -l .)"
[ -z "$unformatted" ] || { gofmt -d .; die "these files are not gofmt'd: $unformatted"; }

# The two plugins come from the tool directives in go.mod, so that pinning
# them is pinning the runtimes; protoc finds them on PATH by name.
log "go generate reproduces the committed codec"
export GOBIN="$scratch/bin"
PATH="$GOBIN:$PATH"
go install tool
cp pb/hub.pb.go pb/hub_grpc.pb.go "$scratch/"
go generate ./...
for f in hub.pb.go hub_grpc.pb.go; do
    if ! cmp -s "$scratch/$f" "pb/$f"; then
        diff -u "$scratch/$f" "pb/$f" || true
        die "pb/$f does not match hub.proto; commit what 'go generate ./...' produced"
    fi
done

log "go vet"
go vet ./...

# -race: wire is small, and the habit is cheaper than the exception.
log "go test -race"
go test -race -count=1 ./...

log "go.mod is tidy"
cp go.mod go.sum "$scratch/"
go mod tidy
cmp -s go.mod "$scratch/go.mod" && cmp -s go.sum "$scratch/go.sum" \
    || die "go.mod or go.sum is not tidy, run 'go mod tidy'"

# The C++ half: built, installed, and loaded by a consumer through the package
# configuration, which is the part nothing else would ever exercise.
log "cmake: build and install cams::api-cpp"
cmake -S . -B "$scratch/build" -DCMAKE_BUILD_TYPE=Release
cmake --build "$scratch/build" --parallel
cmake --install "$scratch/build" --prefix "$scratch/prefix"

log "cmake: find_package(cams-api 1) from a consumer"
cmake -S test/consumer -B "$scratch/consumer" -DCMAKE_PREFIX_PATH="$scratch/prefix"
cmake --build "$scratch/consumer"
"$scratch/consumer/cams-api-consumer"

log "the working tree is unchanged"
if [ "$(tree_state)" != "$tree_before" ]; then
    diff <(printf '%s\n' "$tree_before") <(tree_state) || true
    die "the checks changed the working tree"
fi

log "OK"
