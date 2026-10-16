// Copyright (c) 2022-2024 The authors (see the AUTHORS file)
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

/**
 * gRPC codec for the communication toward the hub.
 */
package pb

//go:generate protoc --proto_path=.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative hub.proto
