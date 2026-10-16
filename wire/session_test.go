// Copyright (c) 2022-2024 The authors (see the AUTHORS file)
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package wire

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestSessionID_IsSetAndStable(t *testing.T) {
	if SessionID == "" {
		t.Fatal("the session id is empty, so every log line joins on nothing")
	}
	// Two draws must differ, or a fleet of agents shares one id and the
	// column stops distinguishing anything.
	if a, b := newSessionID(), newSessionID(); a == b {
		t.Fatalf("two session ids are identical: %q", a)
	}
}

func TestWithSession_TagsAnOutgoingCall(t *testing.T) {
	md, ok := metadata.FromOutgoingContext(WithSession(context.Background()))
	if !ok {
		t.Fatal("no outgoing metadata at all")
	}
	got := md.Get(KeySession)
	if len(got) != 1 || got[0] != SessionID {
		t.Fatalf("session metadata is %v, want [%s]", got, SessionID)
	}
}

// TestWithSession_KeepsWhatWasThere is the ordering trap: the registrar builds
// its metadata block with NewOutgoingContext, which replaces rather than adds,
// so the tag has to be applied to the result.
func TestWithSession_KeepsWhatWasThere(t *testing.T) {
	ctx := metadata.NewOutgoingContext(context.Background(),
		metadata.New(map[string]string{KeyUser: "someone"}))

	md, ok := metadata.FromOutgoingContext(WithSession(ctx))
	if !ok {
		t.Fatal("no outgoing metadata at all")
	}
	if u := md.Get(KeyUser); len(u) != 1 || u[0] != "someone" {
		t.Fatalf("the user was lost: %v", u)
	}
	if s := md.Get(KeySession); len(s) != 1 || s[0] != SessionID {
		t.Fatalf("session metadata is %v, want [%s]", s, SessionID)
	}
}
