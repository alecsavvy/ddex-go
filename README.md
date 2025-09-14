# DDEX Go

A comprehensive Go implementation of DDEX (Digital Data Exchange) standards with support for both pure XML and Protocol Buffer serialization.

## What is DDEX?

DDEX is a consortium of leading media companies, music licensing organizations, digital service providers and technical intermediaries that develop and promote the adoption of global standards for the exchange of information and rights data along the digital supply chain.

## Features

This library provides Go structs and unmarshaling support for:

- **ERN v4.3.2** (Electronic Release Notification) - For communicating music release information
- **MEAD v1.1** (Media Enrichment and Description) - For enriching media metadata
- **PIE v1.0** (Party Identification and Enrichment) - For party/artist information and awards

## Installation

```bash
go get github.com/alecsavvy/ddex-go@latest
```

## Quick Start

```go
package main

import (
    "encoding/xml"
    "fmt"
    "os"
    
    "github.com/alecsavvy/ddex-go"
    ernv432 "github.com/alecsavvy/ddex-go/gen/ddex/ern/v432"
)

func main() {
    // Read DDEX XML file
    xmlData, err := os.ReadFile("release.xml")
    if err != nil {
        panic(err)
    }

    // Unmarshal into typed struct
    var release ernv432.NewReleaseMessage
    err = xml.Unmarshal(xmlData, &release)
    if err != nil {
        panic(err)
    }

    // Access structured data
    fmt.Printf("Message ID: %s\n", release.MessageHeader.MessageId)
    
    // Or use type aliases for convenience
    var typedRelease ddex.NewReleaseMessageV432 = release
    fmt.Printf("Release Count: %d\n", len(typedRelease.ReleaseList.TrackRelease))
}
```

## Supported Message Types

### ERN (Electronic Release Notification) v4.3.2
- `NewReleaseMessage` - New music releases
- `PurgeReleaseMessage` - Release removal notifications

### MEAD (Media Enrichment and Description) v1.1  
- `MeadMessage` - Media metadata enrichment

### PIE (Party Identification and Enrichment) v1.0
- `PieMessage` - Party/artist information
- `PieRequestMessage` - Party information requests

## Type Aliases

For convenience, the main package exports versioned type aliases:

```go
// ERN v4.3.2 - Main message types
type NewReleaseMessageV432   = ernv432.NewReleaseMessage
type PurgeReleaseMessageV432 = ernv432.PurgeReleaseMessage

// MEAD v1.1 types
type MeadMessageV11 = meadv11.MeadMessage

// PIE v1.0 types
type PieMessageV10        = piev10.PieMessage
type PieRequestMessageV10 = piev10.PieRequestMessage
```

## Examples

### Manual Testing with Real DDEX Files

The `examples/` directory contains tools for testing DDEX file parsing with both pure Go and protobuf-generated structs:

#### Using Pure Go Structs (XSD-generated)
```bash
# Parse with XSD-generated Go structs (ddex/ package)
go run examples/xsd/main.go -file path/to/your/ddex-file.xml
```

#### Using Protocol Buffer Structs (with XML tags)
```bash
# Parse with protobuf-generated Go structs that include XML unmarshaling (gen/ package)
go run examples/proto/main.go -file path/to/your/ddex-file.xml
```

For safely storing real DDEX files for testing, create a `test-files/` or `ddex-samples/` directory (already gitignored):

```bash
mkdir test-files
# Copy your DDEX files here
go run examples/proto/main.go -file test-files/sample.xml
```

Both examples will automatically detect the message type (ERN, MEAD, or PIE) and dump the parsed structure using `spew.Dump()` for easy inspection.

## Development

### Running Tests

```bash
make test
```

This runs comprehensive unmarshaling tests for all supported DDEX message types.

**Note on Test Data:**
- **ERN test files** are official DDEX sample files downloaded directly from the DDEX consortium
- **PIE and MEAD test files** are manually created examples and may not be fully accurate representations of real-world usage

### Regenerating Code

The repository provides two approaches for generating Go structs from DDEX XSD schemas:

```bash
# Generate everything (both pure Go and protobuf approaches)
make generate

# Generate only pure Go structs from XSD
make generate-go

# Generate only protobuf definitions and Go code
make generate-proto generate-proto-go
```

### Manual Commands

```bash
# Run tests without generation
go test -v ./...

# Clean generated files and test data  
make clean

# See all available targets
make help
```

## Repository Structure

```
ddex-go/
├── ddex/                    # Pure Go DDEX structs (generated from XSD)
│   ├── ernv432/            # ERN v4.3.2 - Electronic Release Notification
│   ├── meadv11/            # MEAD v1.1 - Media Enrichment and Description
│   └── piev10/             # PIE v1.0 - Party Identification and Enrichment
│
├── proto/                   # Protocol Buffer definitions with XML tags
│   └── ddex/               # Namespace-aware proto organization
│       ├── avs/            # Allowed Value Sets (enums shared across specs)
│       ├── ern/v432/       # ERN v4.3.2 .proto files
│       ├── mead/v11/       # MEAD v1.1 .proto files
│       └── pie/v10/        # PIE v1.0 .proto files
│
├── gen/                     # Generated Go code from proto files
│   └── ddex/               # Mirrors proto structure
│       ├── avs/            # Shared enum types with proper XML tags
│       ├── ern/v432/       # ERN Go code with protobuf + XML support
│       ├── mead/v11/       # MEAD Go code with protobuf + XML support
│       └── pie/v10/        # PIE Go code with protobuf + XML support
│
├── tools/                   # Generation and conversion tools
│   ├── xsd2go/             # XSD to Go generator (using xgen)
│   └── xsd2proto/          # XSD to Proto converter with namespace-aware imports
│
├── examples/                # Usage examples and documentation
│   ├── proto/              # Example using protobuf-generated structs (gen/)
│   └── xsd/                # Example using XSD-generated structs (ddex/)
│
└── xsd/                     # Original DDEX XSD schema files
    ├── ernv432/            # ERN v4.3.2 XSD files
    ├── meadv11/            # MEAD v1.1 XSD files
    └── piev10/             # PIE v1.0 XSD files
```

## Two Approaches

### 1. Pure XML Package (`ddex/`)
- Generated directly from DDEX XSD schemas using xgen
- Native Go XML marshaling/unmarshaling
- Full DDEX XSD compliance
- Lightweight, no external dependencies

### 2. Protocol Buffer Package (`gen/`)
- Generated from `.proto` files with XML tag annotations  
- Supports both binary Protocol Buffer and XML serialization
- Compatible with gRPC and ConnectRPC
- Enables efficient binary serialization while maintaining XML compatibility
- Includes shared `avs/` package for common enum types across all DDEX specifications

## License

This library is for working with DDEX standards. DDEX specifications are developed by the DDEX consortium.