# The image ci/check.sh runs in: everything it takes to build both halves of
# cams-api -- the Go bindings regenerated, vetted and tested, the C++ library
# built, installed and loaded by test/consumer -- and nothing else.
#
# Cut down from cams-hub's docker/ci/Dockerfile, which builds the whole hub and
# so also carries FFmpeg, libarchive, GoogleTest, Google Benchmark, Envoy,
# ninja, clang and gdb. Nothing here uses any of them. What is kept is the base,
# the gRPC and protobuf packages and the Go toolchain, at the same versions.
#
# Nothing of the repository is copied in, so this directory is the whole build
# context. The workflows build it on every run rather than pull it, with the
# same two commands as a person:
#
#     docker build -t cams-api-ci ci
#     docker run --rm --user "$(id -u):$(id -g)" -e HOME=/tmp -v "$PWD:/w" -w /w \
#         cams-api-ci ci/check.sh

# The Go toolchain, from the official image rather than from the distribution:
# it has to be at or above the go line of go.mod, which no distribution
# guarantees. Copied out rather than built on, the way cams-hub takes Envoy:
# the toolchain is statically linked and runs as is on the base below.
FROM golang:1.27.1 AS go

# The release of Ubuntu is not a free choice. Its protobuf-compiler is the
# protoc that generated the committed Go bindings, and protoc stamps its own
# version into them ("protoc v3.21.12" in pb/hub.pb.go). Another release brings
# another protoc, and ci/check.sh then fails on that line: moving the base
# means regenerating pb/ in the same change.
FROM ubuntu:26.04

# An ARG and not an ENV, so that it holds while apt runs and does not follow
# the image into the containers it starts.
ARG DEBIAN_FRONTEND=noninteractive

# What ci/check.sh calls, and what it calls them for:
#  * ca-certificates: `go mod download` and `go install tool` fetch modules over
#    HTTPS;
#  * git: the working tree is compared before and after the checks;
#  * build-essential: the C++ compiler and make for the CMake half, and the C
#    compiler cgo needs, which `go test -race` cannot do without;
#  * cmake and pkg-config: CMakeLists.txt finds grpc++ and protobuf through
#    pkg-config;
#  * protobuf-compiler is protoc, protobuf-compiler-grpc carries
#    grpc_cpp_plugin, and the two -dev packages are what the library links.
RUN set -ex \
 && apt-get update \
 && apt-get dist-upgrade -y \
 && apt-get install -y --no-install-recommends \
        ca-certificates \
        git \
        build-essential \
        cmake \
        pkg-config \
        protobuf-compiler \
        protobuf-compiler-grpc \
        libprotobuf-dev \
        libgrpc++-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

# On PATH through /usr/local/bin rather than through ENV PATH: a CI system that
# runs its steps with a PATH of its own would lose the latter.
COPY --from=go /usr/local/go /usr/local/go
RUN ln -s /usr/local/go/bin/go /usr/local/go/bin/gofmt /usr/local/bin/

# The toolchain in this image is the one that runs. With the default, auto, a
# go.mod whose go line is ahead of it downloads another toolchain silently
# instead of failing.
ENV GOTOOLCHAIN=local

# Fail here, at image build time, rather than halfway through a CI run. --exists
# resolves the Requires of both modules, so a -dev package the two above do not
# pull in is caught too. The versions are printed for the build log.
RUN set -e \
 && for tool in go gofmt git protoc grpc_cpp_plugin cmake pkg-config c++; do \
        command -v "$tool" >/dev/null || { echo "missing: $tool" >&2; exit 1; }; \
    done \
 && pkg-config --exists --print-errors grpc++ protobuf \
 && go version \
 && protoc --version \
 && printf 'grpc++ %s\n' "$(pkg-config --modversion grpc++)" \
 && cmake --version | head -1

CMD ["/bin/bash", "-l"]
