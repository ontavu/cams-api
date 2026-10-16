# cams-api

Read [README.md](README.md) first: it says what the contract holds, in which languages,
and how each is consumed. This file holds the rules for changing it.

## Commands

```sh
ci/check.sh     # everything CI checks, both languages
```

Run it in the image of `ci/Dockerfile`, because that image's `protoc` is the one whose
version is stamped into the committed Go bindings:

```sh
docker build -t cams-api-ci ci
docker run --rm --user "$(id -u):$(id -g)" -e HOME=/tmp -v "$PWD:/w" -w /w \
    cams-api-ci ci/check.sh
```

`--user` keeps the files the run writes, `.git` included, owned by you rather than by
root. Without a command, and with `-it`, the same `docker run` opens a shell in the image.
`.github/workflows/ci.yml` runs these two commands. It builds the image on every run rather
than pulling one, so a change to `ci/Dockerfile` is checked by the run of that change. Add a
check to the script, never to the workflow alone.

The Ubuntu release `ci/Dockerfile` starts from is what pins that `protoc`: moving to
another one means regenerating `pb/` in the same commit.

## `hub.proto` is the source of truth

Both ends of every wire are generated from that single file: the Go bindings under `pb/`,
which are committed, and the C++ ones, which `CMakeLists.txt` generates into its build tree.
The schema is the only place where the two languages meet, which makes it the most
expensive file in the project to get wrong. Its consumers generate nothing: the Go side
imports `pb`, and the hub links `cams::api-cpp`, built from a release of this repository
with its own `protoc`, plugin and compiler, so that the library and the hub cannot be
built against two gRPCs. The keys a call is tagged with are in `wire/` for the same reason:
the agent and the hub are built and upgraded separately, and a key spelled differently on
one side compiles everywhere.

1. **A protocol change starts in `hub.proto`, never in generated code.** Editing
   `pb/*.pb.go` by hand is always wrong: the next generation silently reverts it, and the
   two languages drift apart.
2. **Regenerate the Go bindings in the same commit.** `ci/check.sh` fails when the
   committed bindings no longer match the schema.
3. **Field numbers are the contract, field names are not.** Renaming a field costs nothing
   on the wire; renumbering one, or recycling a retired number, breaks every deployed agent.
   `reserved` what you retire.
4. **Agents and hub are upgraded separately.** Assume a version skew of at least one
   release in both directions, and add fields rather than repurposing them.
5. **The hub's two languages move to a new release together.** In Cams, the requirement
   of `cams-hub/go/go.mod` on this module and the asset `cams-hub/cpp/cmake/CamsApi.cmake`
   downloads name one tag. Its CMake configuration refuses a disagreement, because a
   schema bumped on one side only still compiles on both and fails at runtime, which is
   the most expensive class of mistake here. The agent, in its own repository, takes a
   release on its own schedule, which is what rule 4 is for.

### Regenerating the Go bindings

```sh
GOBIN="$PWD/.bin" go install tool      # the two plugins, at the versions go.mod pins
PATH="$PWD/.bin:$PATH" go generate ./...
```

Run them in a shell in the image of `ci/Dockerfile`, like the checks. With any other
`protoc`, the version stamped into `pb/` changes and `ci/check.sh` fails on that line.

`go generate` runs the directive in `pb/doc.go`:

```sh
protoc --proto_path=.. \
       --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       hub.proto
```

One trap, verified:

* **The generator and the runtime must move together.** The committed bindings come from
  `protoc-gen-go` v1.32.0, which matches the pinned `google.golang.org/protobuf` v1.32.0 and
  `google.golang.org/grpc` v1.61.1. A newer `protoc-gen-go-grpc` emits code that needs
  `grpc.SupportPackageIsVersion9` and the generic stream types, and v1.61.1 has neither.
  Regenerating with newer plugins therefore requires bumping the runtimes in `go.mod` in the
  same change, or the module stops compiling. That is why the plugins are `tool` directives
  in `go.mod` rather than whatever is on `PATH`: `protoc-gen-go` lives in the module of the
  protobuf runtime, so pinning the plugin pins the runtime.

## Releases

A release is a tag, `vX.Y.Z`, and it names both languages at once:

1. Set `project(cams-api VERSION X.Y.Z ...)` in `CMakeLists.txt`, and commit.
2. Tag `vX.Y.Z` and push the tag. `.github/workflows/release.yml` runs `ci/check.sh`,
   refuses a tag that disagrees with the CMake version, and creates the release. The release
   carries the asset `cams-api-X.Y.Z.tar.gz`, built with `git archive`, and its `.sha256`.
3. In Cams, in one commit: `go get github.com/ontavu/cams-api@vX.Y.Z` and `go mod tidy` in
   `cams-hub/go/`, and the new version and SHA-256 in `cams-hub/cpp/cmake/CamsApi.cmake`.
   In cams-agent, `go get` and `go mod tidy`, when the agent needs the release.

When the release run of a tag fails for a reason outside the tag -- the CI image, a
runner -- start the **Release** workflow by hand from the Actions tab, with that tag as
input. It runs the workflow, and builds the image from `ci/Dockerfile`, as they are on the
branch it is started from, against the content of the tag. So the fix lands on `main`, and
the tag stays where it is.

**Never move or recreate a tag, and never replace an asset.** The Go checksum database
remembers the first content it saw for a version, and the hub pins the asset's hash, so a
changed release fails everywhere that already used it. A mistake is fixed by the next
version.

A major version is a contract that broke: the Go module path then takes a `/vN` suffix,
and `find_package` with `SameMajorVersion` refuses the old major.

## Licence header on every file

Every source file, the schema and the CMake files included, starts with:

```
Copyright (c) 2022-2024 The authors (see the AUTHORS file)

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
```

in that file's comment syntax. The header is Exhibit A and nothing more. **Never add
Exhibit B** ("Incompatible With Secondary Licenses"): leaving it out is what keeps this
module compatible with the AGPL as a secondary licence, which is how the Cams agent links
it. The generated `pb/*.pb.go` carry the header too, because `hub.proto` does and protoc
copies its leading comment. They say which generator produced them from which input,
which is a fact about the build, not attribution. In Go, keep a blank line between the
notice and `package`, or Go takes the licence for the package's doc comment.

## Attribution: `Assisted-By`, never `Co-Authored-By`

A model is a tool, not an author. Work an assistant helped with ends its commit message,
and its pull request description, with one line: `Assisted-By: Claude Code`,
`Assisted-By: Codex`, and so on.

- **The tool, not the model.** Model names change every few months; the tool is what was
  actually run.
- **No address**, and no `Co-Authored-By`. Git and GitHub read that trailer as a person.
  In a copyleft project the authors are the copyright holders, and they are listed in
  `AUTHORS`.
- **Nothing in the sources.** No `// generated by <model>`, and no model name anywhere.

## What may be imported

Nothing of Cams: no `cams-go-utils`, no `cams-agent` and no `cams-hub`. Every other module
depends on this one, and a dependency in the other direction would be a cycle between
repositories. `go list -deps ./... | grep ontavu` lists only this module.
