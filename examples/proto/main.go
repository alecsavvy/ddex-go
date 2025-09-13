package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alecsavvy/ddex-go/gen/ernv432"
	"github.com/alecsavvy/ddex-go/gen/meadv11"
	"github.com/alecsavvy/ddex-go/gen/piev10"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	var filePath string
	var outputFormat string
	flag.StringVar(&filePath, "file", "", "Path to DDEX XML file")
	flag.StringVar(&outputFormat, "format", "dump", "Output format: dump, json, text")
	flag.Parse()

	if filePath == "" {
		fmt.Println("Usage: go run main.go -file <path-to-ddex-file> [-format dump|json|text]")
		fmt.Println("\nExample:")
		fmt.Println("  go run main.go -file ../../test-files/sample.xml")
		fmt.Println("  go run main.go -file ../../test-files/sample.xml -format json")
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
		outputMessage(&newRelease, outputFormat)
		return
	}

	// Try ERN PurgeReleaseMessage
	var purgeRelease ernv432.PurgeReleaseMessage
	if err := xml.Unmarshal(data, &purgeRelease); err == nil && purgeRelease.MessageHeader != nil {
		fmt.Println("✓ Parsed as ERN v4.3.2 PurgeReleaseMessage (protobuf)")
		outputMessage(&purgeRelease, outputFormat)
		return
	}

	// Try MEAD
	var mead meadv11.MeadMessage
	if err := xml.Unmarshal(data, &mead); err == nil && mead.MessageHeader != nil {
		fmt.Println("✓ Parsed as MEAD v1.1 MeadMessage (protobuf)")
		outputMessage(&mead, outputFormat)
		return
	}

	// Try PIE Message
	var pie piev10.PieMessage
	if err := xml.Unmarshal(data, &pie); err == nil && pie.MessageHeader != nil {
		fmt.Println("✓ Parsed as PIE v1.0 PieMessage (protobuf)")
		outputMessage(&pie, outputFormat)
		return
	}

	// Try PIE Request
	var pieRequest piev10.PieRequestMessage
	if err := xml.Unmarshal(data, &pieRequest); err == nil && pieRequest.MessageHeader != nil {
		fmt.Println("✓ Parsed as PIE v1.0 PieRequestMessage (protobuf)")
		outputMessage(&pieRequest, outputFormat)
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

func outputMessage(msg proto.Message, format string) {
	switch format {
	case "json":
		jsonData, err := protojson.Marshal(msg)
		if err != nil {
			log.Printf("Failed to marshal to JSON: %v", err)
			fallback(msg)
			return
		}
		fmt.Println(string(jsonData))
	case "text":
		textData, err := prototext.Marshal(msg)
		if err != nil {
			log.Printf("Failed to marshal to text: %v", err)
			fallback(msg)
			return
		}
		fmt.Println(string(textData))
	default:
		fallback(msg)
	}
}

func fallback(msg proto.Message) {
	spew.Dump(msg)
}