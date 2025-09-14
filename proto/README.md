# DDEX Protocol Buffer Definitions

This directory contains Protocol Buffer definitions generated from DDEX XSD schemas.

## Structure

- `ddex/avs/` - Shared Allowed Value Sets (enums) used across all DDEX specifications
- `ddex/ern/v432/` - ERN v4.3.2 .proto files with XML tag annotations
- `ddex/mead/v11/` - MEAD v1.1 .proto files with XML tag annotations  
- `ddex/pie/v10/` - PIE v1.0 .proto files with XML tag annotations

## Features

Each .proto file includes:
- Standard protobuf message definitions
- XML tag annotations via `[(tagger.tags)]` for XML compatibility
- Support for both binary protobuf and XML serialization

## Generation

Generated using:
- `tools/xsd2proto` - Converts XSD schemas to protobuf with namespace-aware imports
- XML tag annotations via `[(tagger.tags)]` for XML compatibility  

Run `make generate-proto` to regenerate from XSD schemas.