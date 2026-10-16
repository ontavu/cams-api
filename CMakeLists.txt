# Copyright (c) 2022-2024 The authors (see the AUTHORS file)
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

# cams-api, the C++ half: hub.proto turned into a library, cams::api-cpp.
#
#     cmake -S . -B build
#     cmake --build build
#     cmake --install build --prefix <prefix>
#     # then, elsewhere: find_package(cams-api 1 CONFIG REQUIRED)
#
# The bindings are generated into the build tree and never committed: the
# output of protoc-gen-grpc-cpp is tied to the installed gRPC, so a committed
# copy is guaranteed to disagree with some other machine's toolchain. The Go
# bindings beside this file are committed because their generator is pinned in
# go.mod; these cannot be.
#
# The cams hub builds this project as an ExternalProject, out of the release
# asset of its tag (cams-hub/cpp/cmake/CamsApi.cmake in ontavu/cams), and hands
# it its own protoc, plugin and compiler through the cache. That is why the
# names below are the hub's: a value already in the cache is what find_program
# keeps.

cmake_minimum_required(VERSION 3.25)
project(cams-api VERSION 1.0.0 LANGUAGES CXX
        DESCRIPTION "The gRPC contract between the cams agents and the hub"
        HOMEPAGE_URL "https://github.com/ontavu/cams-api")

include(GNUInstallDirs)
include(CMakePackageConfigHelpers)

# A default, never an override: a consumer that forces its own standard through
# the cache must win, and a normal variable here would hide it.
if (NOT DEFINED CMAKE_CXX_STANDARD)
    set(CMAKE_CXX_STANDARD 17)
endif ()
set(CMAKE_CXX_STANDARD_REQUIRED ON)
if (NOT DEFINED CMAKE_CXX_EXTENSIONS)
    set(CMAKE_CXX_EXTENSIONS OFF)
endif ()

find_package(PkgConfig REQUIRED)

find_program(PROTOBUF_PROTOC protoc REQUIRED)
find_program(GRPC_CPP_PLUGIN_EXECUTABLE grpc_cpp_plugin REQUIRED)

# Prefixes of our own, so that a project including this one -- or loading the
# package configuration it installs -- keeps its GRPCPP_* and PROTOBUF_*.
pkg_check_modules(CAMS_API_GRPCPP REQUIRED IMPORTED_TARGET grpc++)
pkg_check_modules(CAMS_API_PROTOBUF REQUIRED IMPORTED_TARGET protobuf)

set(CAMS_API_HEADERS
    "${CMAKE_CURRENT_BINARY_DIR}/hub.pb.h"
    "${CMAKE_CURRENT_BINARY_DIR}/hub.grpc.pb.h")
set(CAMS_API_SOURCES
    "${CMAKE_CURRENT_BINARY_DIR}/hub.pb.cc"
    "${CMAKE_CURRENT_BINARY_DIR}/hub.grpc.pb.cc")

# The tools are dependencies too: a different protoc produces different code,
# and the headers it writes refuse to compile against a mismatched runtime.
add_custom_command(
      OUTPUT ${CAMS_API_HEADERS} ${CAMS_API_SOURCES}
      COMMAND "${PROTOBUF_PROTOC}"
        --grpc_out "${CMAKE_CURRENT_BINARY_DIR}"
        --cpp_out "${CMAKE_CURRENT_BINARY_DIR}"
        -I "${CMAKE_CURRENT_SOURCE_DIR}"
        "--plugin=protoc-gen-grpc=${GRPC_CPP_PLUGIN_EXECUTABLE}"
        "${CMAKE_CURRENT_SOURCE_DIR}/hub.proto"
      DEPENDS "${CMAKE_CURRENT_SOURCE_DIR}/hub.proto"
        "${PROTOBUF_PROTOC}" "${GRPC_CPP_PLUGIN_EXECUTABLE}"
      COMMENT "Generating the C++ gRPC codec from hub.proto"
      VERBATIM)

# No warning flags and no -Werror: generated code is not ours to keep
# warning-free, and a build hostage to the installed protoc is not a build.
add_library(cams-api-cpp STATIC ${CAMS_API_SOURCES})
add_library(cams::api-cpp ALIAS cams-api-cpp)
set_target_properties(cams-api-cpp PROPERTIES EXPORT_NAME api-cpp)

# Flat, under include/cams-api/: protoc writes `#include "hub.pb.h"` into
# hub.grpc.pb.h, so the two have to sit side by side and be reached by name.
target_sources(cams-api-cpp PUBLIC
        FILE_SET HEADERS
        BASE_DIRS "${CMAKE_CURRENT_BINARY_DIR}"
        FILES ${CAMS_API_HEADERS})

target_link_libraries(cams-api-cpp PUBLIC
        PkgConfig::CAMS_API_GRPCPP
        PkgConfig::CAMS_API_PROTOBUF)

#----------------------------------------------------------------------------#
# What an installed cams-api is: the library, its headers, the schema it was
# generated from, and a package configuration for find_package.

set(CAMS_API_CMAKEDIR "${CMAKE_INSTALL_LIBDIR}/cmake/cams-api")

install(TARGETS cams-api-cpp EXPORT cams-api-targets
        ARCHIVE DESTINATION "${CMAKE_INSTALL_LIBDIR}"
        FILE_SET HEADERS DESTINATION "${CMAKE_INSTALL_INCLUDEDIR}/cams-api")

# The schema, so that a consumer in a third language generates from the very
# file this library was generated from.
install(FILES hub.proto DESTINATION "${CMAKE_INSTALL_DATADIR}/cams-api")

# The MPL-2.0, beside what it covers: whoever ships a binary linking this
# library has to tell its recipients how to obtain the source of this part
# (section 3.2 of the licence).
install(FILES LICENSE DESTINATION "${CMAKE_INSTALL_DOCDIR}")

install(EXPORT cams-api-targets
        NAMESPACE cams::
        DESTINATION "${CAMS_API_CMAKEDIR}")

configure_package_config_file(cmake/cams-apiConfig.cmake.in
        "${CMAKE_CURRENT_BINARY_DIR}/cams-apiConfig.cmake"
        INSTALL_DESTINATION "${CAMS_API_CMAKEDIR}")

# A major version is a contract that broke. Within one, fields are added and
# never renumbered (see AGENTS.md), so any minor release will do.
write_basic_package_version_file(
        "${CMAKE_CURRENT_BINARY_DIR}/cams-apiConfigVersion.cmake"
        COMPATIBILITY SameMajorVersion)

install(FILES
        "${CMAKE_CURRENT_BINARY_DIR}/cams-apiConfig.cmake"
        "${CMAKE_CURRENT_BINARY_DIR}/cams-apiConfigVersion.cmake"
        DESTINATION "${CAMS_API_CMAKEDIR}")
