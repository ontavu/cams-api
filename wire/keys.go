// Copyright (c) 2022-2024 The authors (see the AUTHORS file)
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package wire is what both ends of a call have to agree on and that
// hub.proto cannot say: the metadata keys a call is tagged with, and the
// keepalive each side expects of the other.
//
// It lives beside the schema rather than in either party because the agent and
// the hub are built and upgraded separately. A key spelled differently on one
// side compiles everywhere and is refused at runtime -- which is how every
// upload was once refused, when the hub expected "camera" and the agent sent
// "stream".
package wire

// The metadata keys an agent tags its calls with. The C++ hub mirrors them in
// cams-hub/cpp/libupstream/MediaService.hpp, in ontavu/cams, and
// cams-hub/cpp/tests/service_test.cpp there asserts the literal strings, not
// the hub's constants.
const (
	KeyUser   string = "user"
	KeyStream        = "stream"

	// KeySession carries SessionID, the join column between the logs of
	// the agent, the control plane and the data plane. cams-ctrl hands it to
	// its access log as utils.ServerConfig.SessionKey; see session.go.
	KeySession = "session-id"
)
