# DDEX Protocol Buffer Definitions

This directory contains Protocol Buffer definitions generated from DDEX XSD schemas.

## Structure

- `ernv432/` - ERN v4.3.2 .proto files with XML tag annotations
- `meadv11/` - MEAD v1.1 .proto files with XML tag annotations
- `piev10/` - PIE v1.0 .proto files with XML tag annotations

## Features

Each .proto file includes:
- Standard protobuf message definitions
- XML tag annotations via `[(tagger.tags)]` for XML compatibility
- Support for both binary protobuf and XML serialization

## Generation

These files are generated using `tools/xsd2proto` which converts XSD schemas to protobuf with appropriate XML tags.