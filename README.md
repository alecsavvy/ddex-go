# DDEX Go

A comprehensive Go implementation of DDEX (Digital Data Exchange) standards with full Protocol Buffer, JSON, and XML serialization support.

## What is DDEX?

DDEX is a consortium of leading media companies, music licensing organizations, digital service providers and technical intermediaries that develop and promote the adoption of global standards for the exchange of information and rights data along the digital supply chain.

## Features

This library provides Go structs with Protocol Buffer, JSON, and XML serialization support for:

- **ERN v4.3.2** (Electronic Release Notification) - For communicating music release information
- **MEAD v1.1** (Media Enrichment and Description) - For enriching media metadata
- **PIE v1.0** (Party Identification and Enrichment) - For party/artist information and awards

### Key Capabilities

- **Triple serialization**: Full support for Protocol Buffer binary, JSON, and XML formats
- **XML compliance**: Generated structs maintain complete DDEX XSD compliance
- **gRPC/ConnectRPC ready**: Protocol Buffer definitions work seamlessly with RPC frameworks
- **Bidirectional conversion**: Convert between protobuf, JSON, and XML without data loss
- **Performance optimized**: Efficient binary serialization for high-throughput applications
- **Type safety**: Strong typing with comprehensive test coverage

## Installation

```bash
go get github.com/alecsavvy/ddex-go@latest
```

## Quick Start

### Basic XML Parsing

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

    // Use type aliases for convenience
    var typedRelease ddex.NewReleaseMessageV432 = release
    fmt.Printf("Release Count: %d\n", len(typedRelease.ReleaseList.TrackRelease))
}
```

### Protocol Buffer and JSON Serialization

```go
package main

import (
    "encoding/json"
    "fmt"
    ernv432 "github.com/alecsavvy/ddex-go/gen/ddex/ern/v432"
    "google.golang.org/protobuf/proto"
)

func main() {
    // Create a new release message
    release := &ernv432.NewReleaseMessage{
        MessageHeader: &ernv432.MessageHeader{
            MessageId: "MSG-12345",
        },
    }

    // Serialize to Protocol Buffer binary format
    protoData, err := proto.Marshal(release)
    if err != nil {
        panic(err)
    }

    // Serialize to JSON
    jsonData, err := json.Marshal(release)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Proto size: %d bytes\n", len(protoData))
    fmt.Printf("JSON: %s\n", string(jsonData))

    // Deserialize from binary format
    var decoded ernv432.NewReleaseMessage
    err = proto.Unmarshal(protoData, &decoded)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Message ID: %s\n", decoded.MessageHeader.MessageId)
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

### Testing with Real DDEX Files

The `examples/proto/` directory contains a comprehensive tool for parsing and validating DDEX files:

```bash
# Parse any DDEX file - automatically detects message type
go run examples/proto/main.go -file path/to/your/ddex-file.xml

# Examples with different message types
go run examples/proto/main.go -file testdata/ernv432/Samples43/1\ Audio.xml
go run examples/proto/main.go -file testdata/meadv11/mead_award_example.xml
go run examples/proto/main.go -file testdata/piev10/pie_award_example.xml
```

For safely storing real DDEX files for testing, create a `test-files/` or `ddex-samples/` directory (gitignored):

```bash
mkdir test-files
# Copy your DDEX files here
go run examples/proto/main.go -file test-files/sample.xml
```

The example automatically detects the message type (ERN, MEAD, or PIE) and provides detailed output using `spew.Dump()` for easy inspection.

## Development

### Running Tests

```bash
# Run all tests including comprehensive validation
make test

# Run specific test suites
make test-comprehensive  # Conformance, roundtrip, and completeness tests
make test-roundtrip     # XML bidirectional conversion tests
make benchmark          # Performance benchmarks
```

**Test Coverage:**
- **Conformance tests**: Validate against official DDEX sample files
- **Roundtrip tests**: Ensure XML ↔ protobuf conversion without data loss
- **Field completeness**: Verify all XSD fields are properly mapped
- **Performance benchmarks**: Memory and speed optimization validation

**Test Data:**
- **ERN test files**: Official DDEX consortium sample files (complete accuracy)
- **MEAD/PIE test files**: Manually created examples (representative but not exhaustive)

### Code Generation

The library uses a sophisticated generation pipeline:

```bash
# Complete generation workflow
make generate           # XSD → proto → Go with XML tags

# Individual steps
make generate-proto     # XSD schemas → Protocol Buffer definitions
make generate-proto-go  # Proto files → Go structs with XML tags
make buf-generate      # Alternative: use buf for Go generation
make buf-lint          # Lint protobuf files

# See all available commands
make help
```

### Generation Pipeline Details

1. **XSD → Proto**: `tools/xsd2proto/` converts DDEX XSD schemas to protobuf with XML annotations
2. **Proto → Go**: `buf generate` creates Go structs with protobuf support
3. **XML Tag Injection**: `protoc-go-inject-tag` adds XML struct tags for compatibility
4. **Enum Strings**: `tools/generate-enum-strings/` generates human-readable enum methods

### Manual Commands

```bash
# Run tests without generation
go test -v ./...

# Clean generated files and test data
make clean

# Force refresh of test data
make testdata-refresh
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
│   ├── xsd2proto/          # XSD to Proto converter with namespace-aware imports
│   └── generate-enum-strings/ # Enum string method generator
│
├── examples/                # Usage examples and documentation
│   └── proto/              # Comprehensive parsing example (supports all message types)
│
├── testdata/                # Test files for validation
│   ├── ernv432/           # Official DDEX consortium sample files
│   ├── meadv11/           # MEAD test examples
│   └── piev10/            # PIE test examples
│
├── xsd/                     # Original DDEX XSD schema files
│   ├── ernv432/           # ERN v4.3.2 XSD files
│   ├── meadv11/           # MEAD v1.1 XSD files
│   └── piev10/            # PIE v1.0 XSD files
│
├── buf.yaml                 # Protocol Buffer configuration
├── buf.gen.yaml            # Code generation configuration
├── Makefile                # Build automation
└── ddex.go                 # Main package with type aliases
```

## Protocol Buffer Implementation

This library implements a sophisticated triple-serialization approach:

### Core Architecture
- **Protocol Buffer definitions** with XML tag annotations via `[(tagger.tags)]`
- **Generated Go structs** support `proto.Marshal/Unmarshal`, `json.Marshal/Unmarshal`, and `xml.Marshal/Unmarshal`
- **Shared enum types** in `ddex/avs/` package used across all DDEX specifications
- **Namespace-aware imports** ensure proper XSD compliance and proto organization

### Benefits
- **Performance**: Binary protobuf serialization for high-throughput applications
- **Compatibility**: Full XML compliance for DDEX standard requirements
- **Interoperability**: Native gRPC/ConnectRPC support for microservices
- **Type Safety**: Strong typing with comprehensive validation
- **Bidirectional**: Perfect XML ↔ JSON ↔ protobuf conversion without data loss

### Usage Patterns
- Use **XML** for DDEX standard compliance and external integrations
- Use **JSON** for REST APIs, web services, and JavaScript interoperability
- Use **Protocol Buffers** for internal APIs, microservices, and performance-critical applications
- Convert seamlessly between all three formats as needed

## License

This library is for working with DDEX standards. DDEX specifications are developed by the DDEX consortium.