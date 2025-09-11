# DDEX Go

A Go library for parsing and working with DDEX (Digital Data Exchange) XML messages, providing type-safe access to music industry standard data formats.

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
    "github.com/alecsavvy/ddex-go/ernv432"
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
type NewReleaseMessageV432 = ernv432.NewReleaseMessage
type MeadMessageV11 = meadv11.MeadMessage
type PieMessageV10 = piev10.PieMessage
```

## Examples

### Manual Testing with Real DDEX Files

The `examples/` directory contains a simple tool for testing DDEX file parsing:

```bash
# Parse any DDEX file and dump the parsed structure
go run examples/main.go -file path/to/your/ddex-file.xml
```

For safely storing real DDEX files for testing, create a `test-files/` or `ddex-samples/` directory (already gitignored):

```bash
mkdir test-files
# Copy your DDEX files here
go run examples/main.go -file test-files/sample.xml
```

The example will automatically detect the message type (ERN, MEAD, or PIE) and output the parsed structure.

## Development

### Running Tests

```bash
make test
```

This will automatically download DDEX sample files and run comprehensive unmarshaling tests.

**Note on Test Data:**
- **ERN test files** are official DDEX sample files downloaded directly from the DDEX consortium
- **PIE and MEAD test files** are manually created examples and may not be fully accurate representations of real-world usage

### Regenerating Code

The Go structs are generated from official DDEX XSD schemas using [xgen](https://github.com/xuri/xgen):

```bash
make generate
```

### Manual Commands

```bash
# Download test data only
make testdata

# Run tests without downloading
go test -v ./...

# Clean generated files and test data  
make clean
```

## Architecture

- **Generated packages** (`ernv432/`, `meadv11/`, `piev10/`) contain auto-generated Go structs from XSD schemas
- **Main package** (`ddex.go`) provides versioned type aliases for discoverability
- **Generator** (`cmd/ddex-go-gen/`) downloads schemas and generates Go code
- **Tests** validate unmarshaling against official DDEX sample files

## License

This library is for working with DDEX standards. DDEX specifications are developed by the DDEX consortium.