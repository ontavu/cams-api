// Copyright (c) 2022-2024 The authors (see the AUTHORS file)
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package wire

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"
)

// SessionID names this process's run in the logs of everything it talks to.
//
// The hub's logging interceptor has read a "session-id" metadata key since it
// was written, and nothing ever set it. So the three logs a single camera
// produces -- the agent's, the control plane's, the data plane's -- had no
// column to join on, and following one stream across them meant guessing from
// timestamps. This is that column.
//
// One value per process, so a restart produces a new one. That is deliberate:
// a restart is one of the things an operator is trying to see, and an id that
// survived it would hide the seam.
//
// It is diagnostic only. It is asserted by the caller, like every other
// identity on this wire today, and must never come to mean anything more --
// see the authentication section of ARCHITECTURE.md.
var SessionID = newSessionID()

func newSessionID() string {
	var raw [8]byte
	if _, err := rand.Read(raw[:]); err != nil {
		// Not worth refusing to start over. A clock reading is still unique
		// enough to join three logs from one run, which is all this is for.
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(raw[:])
}

// WithSession tags an outgoing call with SessionID.
//
// It appends, so it composes with metadata already on the context. Where a
// caller builds the metadata block from scratch with NewOutgoingContext, this
// has to be applied to the result rather than before it, or the block replaces
// the tag.
func WithSession(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, KeySession, SessionID)
}
