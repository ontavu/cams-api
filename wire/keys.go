// Copyright (c) 2022-2024 The authors (see the AUTHORS file)
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package wire is what both ends of a call have to agree on and that
// hub.proto cannot say: the metadata keys a call is tagged with, the session
// id it carries, and the keepalive each side expects of the other.
//
// It lives beside the schema rather than in either party because the two ends
// of a call are built and upgraded separately. A key spelled differently on one
// side compiles everywhere and is refused at runtime.
package wire

// The metadata keys an agent tags its calls with. They are the contract as
// much as the schema is: a party in C++, which cannot import them, spells the
// same strings, and renaming one breaks every deployed agent.
const (
	KeyUser   string = "user"
	KeyStream        = "stream"

	// KeySession carries SessionID, the join column between the logs of
	// the agent, the control plane and the data plane; see session.go.
	KeySession = "session-id"
)
