# DDEX Protocol Buffer Definitions

This directory contains Protocol Buffer definitions generated from DDEX XSD schemas with XML tag annotations for native XML support.

## Structure

- `ddex/avs/` - Shared Allowed Value Sets (enums) used across all DDEX specifications
- `ddex/ern/v432/` - ERN v4.3.2 .proto files with XML tag annotations for native XML support
- `ddex/mead/v11/` - MEAD v1.1 .proto files with XML tag annotations for native XML support
- `ddex/pie/v10/` - PIE v1.0 .proto files with XML tag annotations for native XML support

## Features

Each .proto file includes:
- Standard protobuf message definitions for high-performance binary serialization
- XML tag annotations via `@gotags:` comments for DDEX XSD-compliant XML support
- Namespace-aware imports ensuring proper schema organization
- Support for bidirectional conversion between protobuf, XML, and JSON

## Generation Pipeline

Generated using:
- `tools/xsd2proto/` - Converts XSD schemas to protobuf with namespace-aware imports
- XML tag annotations via `@gotags:` comments processed by `protoc-go-inject-tag`
- Comprehensive type mapping from XSD primitives to protobuf types

Run `make generate-proto` to regenerate from XSD schemas.