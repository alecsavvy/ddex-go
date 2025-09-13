package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	ernv432 "github.com/alecsavvy/ddex-go/gen/ddex/ern/v432"
	meadv11 "github.com/alecsavvy/ddex-go/gen/ddex/mead/v11"
	piev10 "github.com/alecsavvy/ddex-go/gen/ddex/pie/v10"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "Path to DDEX XML file")
	flag.Parse()

	if filePath == "" {
		fmt.Println("Usage: go run main.go -file <path-to-ddex-file>")
		fmt.Println("\nExample:")
		fmt.Println("  go run main.go -file ../../test-files/sample.xml")
		os.Exit(1)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	fileName := filepath.Base(filePath)
	fmt.Printf("Processing: %s\n\n", fileName)

	// Try ERN NewReleaseMessage
	var newRelease ernv432.NewReleaseMessage
	if err := xml.Unmarshal(data, &newRelease); err == nil && newRelease.MessageHeader != nil {
		fmt.Println("✓ Parsed as ERN v4.3.2 NewReleaseMessage (protobuf)")
		spew.Dump(&newRelease)
		return
	}

	// Try ERN PurgeReleaseMessage
	var purgeRelease ernv432.PurgeReleaseMessage
	if err := xml.Unmarshal(data, &purgeRelease); err == nil && purgeRelease.MessageHeader != nil {
		fmt.Println("✓ Parsed as ERN v4.3.2 PurgeReleaseMessage (protobuf)")
		spew.Dump(&purgeRelease)
		return
	}

	// Try MEAD
	var mead meadv11.MeadMessage
	if err := xml.Unmarshal(data, &mead); err == nil && mead.MessageHeader != nil {
		fmt.Println("✓ Parsed as MEAD v1.1 MeadMessage (protobuf)")
		spew.Dump(&mead)
		return
	}

	// Try PIE Message
	var pie piev10.PieMessage
	if err := xml.Unmarshal(data, &pie); err == nil && pie.MessageHeader != nil {
		fmt.Println("✓ Parsed as PIE v1.0 PieMessage (protobuf)")
		spew.Dump(&pie)
		return
	}

	// Try PIE Request
	var pieRequest piev10.PieRequestMessage
	if err := xml.Unmarshal(data, &pieRequest); err == nil && pieRequest.MessageHeader != nil {
		fmt.Println("✓ Parsed as PIE v1.0 PieRequestMessage (protobuf)")
		spew.Dump(&pieRequest)
		return
	}

	fmt.Println("❌ Could not parse file as any supported DDEX message type (protobuf)")
	fmt.Println("\nSupported types:")
	fmt.Println("  - ERN v4.3.2 (NewReleaseMessage, PurgeReleaseMessage)")
	fmt.Println("  - MEAD v1.1 (MeadMessage)")
	fmt.Println("  - PIE v1.0 (PieMessage, PieRequestMessage)")
	fmt.Println("\nNote: This example uses protobuf-generated structs to parse XML data.")
	fmt.Println("Compare with examples/xsd/main.go which uses XSD-generated structs.")
}
