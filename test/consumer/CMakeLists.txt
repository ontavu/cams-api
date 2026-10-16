# Copyright (c) 2022-2024 The authors (see the AUTHORS file)
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

# A consumer of an installed cams-api, and nothing else: what CI configures
# against the prefix it has just installed into, so that the package
# configuration is checked by somebody.
#
#     cmake -S test/consumer -B build-consumer -DCMAKE_PREFIX_PATH=<prefix>
#     cmake --build build-consumer && build-consumer/cams-api-consumer

cmake_minimum_required(VERSION 3.25)
project(cams-api-consumer LANGUAGES CXX)

find_package(cams-api 1 CONFIG REQUIRED)

add_executable(cams-api-consumer main.cpp)
target_link_libraries(cams-api-consumer PRIVATE cams::api-cpp)
