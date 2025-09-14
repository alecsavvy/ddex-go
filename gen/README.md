# Generated Go Code from Protocol Buffers

This directory contains Go code generated from the Protocol Buffer definitions in `proto/`.

## Packages

- `ddex/avs/` - Shared Allowed Value Sets (enums) used across all DDEX specifications
- `ddex/ern/v432/` - Generated Go code for ERN v4.3.2 with protobuf + XML support
- `ddex/mead/v11/` - Generated Go code for MEAD v1.1 with protobuf + XML support  
- `ddex/pie/v10/` - Generated Go code for PIE v1.0 with protobuf + XML support

## Features

Generated code includes:
- Protocol Buffer binary serialization
- XML marshaling/unmarshaling via struct tags
- ConnectRPC/gRPC service support
- Full DDEX XSD compliance in XML mode

## Generation

Generated using:
- `buf generate` with protoc-gen-go for Protocol Buffer support
- `protoc-go-inject-tag` for XML tag injection into generated structs  
- Custom enum string generator for human-readable enum values

Run `make generate-proto-go` or `make buf-generate` to regenerate.