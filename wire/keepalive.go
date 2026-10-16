// Copyright (c) 2022-2024 The authors (see the AUTHORS file)
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package wire

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// The keepalive of a call between an agent and the cloud, both halves of it.
//
// They are one contract, not two settings: the client pings, and the server's
// enforcement policy decides whether a ping is welcome. A policy that refuses
// what the client sends answers with GOAWAY, which turns the keepalive into a
// reconnect loop -- and the agent and the control plane are built from
// different trees, so the only place both halves can be read together is here.

// KeepaliveClient makes gRPC prove the peer is still there.
//
// This is not a nicety. The upstream agent's whole reconnect logic hangs off
// stream.Recv returning an error, and the hub sends nothing on the control
// stream unless a viewer asks for a camera -- so an idle stream carries no
// traffic at all. Without a keepalive, a NAT between an agent and the hub drops
// the connection after its idle timeout, nobody is told, Recv parks forever,
// and the agent is a zombie that still reports its link as up.
//
// PermitWithoutStream is on because the interesting case is exactly the one
// where no RPC is in flight.
var KeepaliveClient = keepalive.ClientParameters{
	Time:                30 * time.Second,
	Timeout:             10 * time.Second,
	PermitWithoutStream: true,
}

// KeepaliveServer lets the hub notice an agent that has gone away without
// saying so, which is the mirror of the client-side problem: a control stream
// is idle most of the time, so nothing else would reveal it.
var KeepaliveServer = keepalive.ServerParameters{
	Time:    30 * time.Second,
	Timeout: 10 * time.Second,
}

// KeepalivePolicy has to admit the pings the agents send. The gRPC default
// refuses a ping from a client with no active stream and answers with GOAWAY,
// which would turn KeepaliveClient into a reconnect loop.
var KeepalivePolicy = keepalive.EnforcementPolicy{
	MinTime:             15 * time.Second,
	PermitWithoutStream: true,
}

// ServerOptions applies KeepaliveServer and KeepalivePolicy to a server. Every
// server an agent dials is built with them.
func ServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveParams(KeepaliveServer),
		grpc.KeepaliveEnforcementPolicy(KeepalivePolicy),
	}
}

// DialOptions applies KeepaliveClient to a client. Every dial between an agent,
// the cloud and a viewer process passes it: the helpers that dial know nothing
// of this contract and apply no keepalive of their own.
func DialOptions() []grpc.DialOption {
	return []grpc.DialOption{grpc.WithKeepaliveParams(KeepaliveClient)}
}
