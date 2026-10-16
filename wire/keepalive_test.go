// Copyright (c) 2022-2024 The authors (see the AUTHORS file)
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package wire

import "testing"

// TestKeepalivePolicy_AdmitsTheClient holds the contract the two halves make:
// a server enforcing KeepalivePolicy never sends GOAWAY to a client pinging
// with KeepaliveClient.
func TestKeepalivePolicy_AdmitsTheClient(t *testing.T) {
	if KeepaliveClient.Time < KeepalivePolicy.MinTime {
		t.Fatalf("the client pings every %v and the server refuses pings closer than %v",
			KeepaliveClient.Time, KeepalivePolicy.MinTime)
	}
	if KeepaliveClient.PermitWithoutStream && !KeepalivePolicy.PermitWithoutStream {
		t.Fatal("the client pings an idle connection and the server refuses that")
	}
}
