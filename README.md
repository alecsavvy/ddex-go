# DDEX Go

A comprehensive Go implementation of DDEX (Digital Data Exchange) standards with support for both pure XML and Protocol Buffer serialization.

## What is DDEX?

DDEX is a consortium of leading media companies, music licensing organizations, digital service providers and technical intermediaries that develop and promote the adoption of global standards for the exchange of information and rights data along the digital supply chain.

## Features

This library provides Go structs with both Protocol Buffer and XML serialization support for:

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

    // Convert back to XML with proper header
    regeneratedXML, err := xml.MarshalIndent(&release, "", "  ")
    if err != nil {
        panic(err)
    }

    // Add XML declaration for complete DDEX document
    fullXML := xml.Header + string(regeneratedXML)
    fmt.Println(fullXML)

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

The `examples/` directory contains tools for testing DDEX file parsing:

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

The example will automatically detect the message type (ERN, MEAD, or PIE) and dump the parsed structure using `spew.Dump()` for easy inspection.

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
# Generate everything
make generate

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
│   └── xsd2proto/          # XSD to Proto converter with namespace-aware imports
│
├── examples/                # Usage examples and documentation
│   └── proto/              # Example using protobuf-generated structs (gen/)
│
└── xsd/                     # Original DDEX XSD schema files
    ├── ernv432/            # ERN v4.3.2 XSD files
    ├── meadv11/            # MEAD v1.1 XSD files
    └── piev10/             # PIE v1.0 XSD files
```

## Protocol Buffer Implementation

The library uses Protocol Buffer definitions with XML tag annotations to provide:
- Both binary Protocol Buffer and XML serialization support
- Full compatibility with gRPC and ConnectRPC
- Efficient binary serialization while maintaining XML compatibility
- Shared `avs/` package for common enum types across all DDEX specifications
- Complete bidirectional conversion between proto structs and valid DDEX XML

## License

This library is for working with DDEX standards. DDEX specifications are developed by the DDEX consortium.