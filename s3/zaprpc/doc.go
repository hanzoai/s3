// Package zaprpc is the ZAP-native transport + codec foundation for Hanzo S3,
// replacing the upstream gRPC + protobuf stack.
//
// Design (decomplected: shape / codec / transport):
//
//   - Shape: message types stay plain Go structs (field set + GetX getters
//     preserved) so the ~988 consumer files are untouched by the rip.
//   - Codec: each message gains MarshalZap()/UnmarshalZap(), hand-rolled over
//     the luxfi/zap wire format (Builder on write, Object view on read). No
//     protobuf wire, no reflection. Variable-length fields (string/bytes/list/
//     nested) use zap's tail-pointer layout; scalars sit inline in the object's
//     fixed payload at codegen-assigned offsets.
//   - Transport: gRPC services become ZAP procedures — server registers
//     handlers, client issues Call(ctx, procedure, req). Streaming RPCs map to
//     a ZAP stream abstraction (built in F4).
//
// The per-message Marshal/Unmarshal + per-service stubs are emitted by the
// proto->zap generator (cmd/protoc-gen-zap). This file documents the contract
// the generator targets; codec_proof_test.go is the hand-written reference the
// generator output must match byte-for-byte.
package zaprpc
