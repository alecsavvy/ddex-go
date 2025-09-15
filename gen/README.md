# Generated Go Code from Protocol Buffers

This directory contains Go code generated from the Protocol Buffer definitions in `proto/`.

## Packages

- `ddex/avs/` - Shared Allowed Value Sets (enums) used across all DDEX specifications
- `ddex/ern/v432/` - Generated Go code for ERN v4.3.2 with native XML + protobuf support
- `ddex/mead/v11/` - Generated Go code for MEAD v1.1 with native XML + protobuf support
- `ddex/pie/v10/` - Generated Go code for PIE v1.0 with native XML + protobuf support

## Features

Generated code includes:
- **Native XML support**: Full DDEX XSD-compliant XML marshal/unmarshal
- **Protocol Buffer serialization**: High-performance binary format for microservices
- **JSON serialization**: Standard Go JSON support for REST APIs
- **gRPC/ConnectRPC support**: Ready for use in microservice architectures
- **Type safety**: Strong typing with comprehensive validation

## Generation Pipeline

Generated using:
- `buf generate` with protoc-gen-go for Protocol Buffer support
- `protoc-go-inject-tag` for XML tag injection into generated structs
- `tools/generate-go-extensions/` for enum strings and XML marshaling methods

Run `make generate-proto-go` or `make buf-generate` to regenerate.