// Copyright (c) 2022-2024 The authors (see the AUTHORS file)
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Reaches both halves of the codec, so that neither can be missing from what
// the package links: a message round trip, which is the generated protobuf code
// and libprotobuf, and a service skeleton, which is the generated gRPC code and
// grpc++.

#include <iostream>

#include "hub.grpc.pb.h"
#include "hub.pb.h"

int main() {
    cams::api::hub::StreamId sent;
    sent.set_user("someone");

    cams::api::hub::StreamId received;
    if (!received.ParseFromString(sent.SerializeAsString()) || received.user() != "someone") {
        std::cerr << "a StreamId did not survive its own codec\n";
        return 1;
    }

    cams::api::hub::Uploader::Service service;
    std::cout << cams::api::hub::Uploader::service_full_name() << '\n';
    return 0;
}
