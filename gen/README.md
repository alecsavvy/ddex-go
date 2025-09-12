# Generated Go Code from Protocol Buffers

This directory contains Go code generated from the Protocol Buffer definitions in `proto/`.

## Packages

- `ernv432/` - Generated Go code for ERN v4.3.2 with protobuf + XML support
- `meadv11/` - Generated Go code for MEAD v1.1 with protobuf + XML support  
- `piev10/` - Generated Go code for PIE v1.0 with protobuf + XML support

## Features

Generated code includes:
- Protocol Buffer binary serialization
- XML marshaling/unmarshaling via struct tags
- ConnectRPC/gRPC service support
- Full DDEX XSD compliance in XML mode

## Generation

Generated using buf with protoc-gen-go and protoc-gen-gotag for XML tag injection.