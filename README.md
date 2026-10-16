# cams-api

The gRPC contract of Cams: what its agents, its control plane, its data plane and its
viewer processes say to one another. It comes in two languages, both generated from one
file, [`hub.proto`](hub.proto). It is under the MPL-2.0, so that code under any licence can
build against it, copyleft or proprietary.

| Path | What it is |
|---|---|
| `hub.proto` | The contract. **Source of truth** for both languages. |
| `pb/` | The Go bindings, generated and committed. Never hand-edited. |
| `wire/` | What both ends of a call agree on that the schema cannot say: the metadata keys (`user`, `stream`, `session-id`), the session tag, and the keepalive each side expects of the other, with `ServerOptions` and `DialOptions`, the gRPC options that apply it. |
| `CMakeLists.txt` | The C++ half: generates the bindings into the build tree and installs them as the static library `cams::api-cpp`, with a package configuration. |

## Go

```go
import (
    "github.com/ontavu/cams-api/pb"
    "github.com/ontavu/cams-api/wire"
)
```

## C++

The C++ bindings are **never committed**. The output of `grpc_cpp_plugin` is tied to
the gRPC it is compiled against, so a committed copy would disagree with somebody's
toolchain. Build and install them with the gRPC you link:

```shell
cmake -S . -B build && cmake --build build
cmake --install build --prefix <prefix>
```

then, in a consumer, write `find_package(cams-api 1 CONFIG REQUIRED)` and link
`cams::api-cpp`. The headers are flat, so a consumer includes `"hub.pb.h"` and
`"hub.grpc.pb.h"`. `test/consumer/` is such a consumer, and CI builds it. The build needs
CMake 3.25 or newer, a C++17 compiler, `protoc`, `grpc_cpp_plugin`, and `pkg-config` with
the modules `grpc++` and `protobuf`. The image of [`ci/Dockerfile`](ci/Dockerfile) has all
of them.

A project can also build this one as an `ExternalProject`, out of the source asset of a
release, pinned by the SHA-256 published beside it. It then hands it its own `protoc` and
plugin through the cache, as `PROTOBUF_PROTOC` and `GRPC_CPP_PLUGIN_EXECUTABLE`.

## Versions

One tag versions both languages: `vX.Y.Z` is the Go module version, the CMake `project()`
version, and the release whose asset is `cams-api-X.Y.Z.tar.gz`. A major version is a
contract that broke; within one, fields are added and never renumbered.

A project that uses both languages takes them from the same tag: a schema bumped on one
side only still compiles on both, and fails at runtime.

## Licence

[MPL-2.0](LICENSE). The authors are listed in [AUTHORS](AUTHORS). `cmake --install`
installs the licence with the library, in `share/doc/cams-api/`, because whoever ships a
binary that links it has to tell its recipients where the source of this part is
(section 3.2).
