package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// DDEX specifications to convert
var specs = []struct {
	name     string
	version  string
	mainFile string
}{
	{"ern", "432", "release-notification.xsd"},
	{"mead", "11", "media-enrichment-and-description.xsd"},
	{"pie", "10", "party-identification-and-enrichment.xsd"},
}

func main() {
	for _, spec := range specs {
		log.Printf("Converting %s v%s to protobuf...", spec.name, spec.version)
		
		if err := validateSchemas(spec); err != nil {
			log.Fatalf("Schema validation failed for %s v%s: %v", spec.name, spec.version, err)
		}
		
		if err := convertToProto(spec); err != nil {
			log.Fatalf("Failed to convert %s v%s: %v", spec.name, spec.version, err)
		}
	}
}

func validateSchemas(spec struct{ name, version, mainFile string }) error {
	schemasDir := filepath.Join("xsd", spec.name+"v"+spec.version)
	
	// Check if schema directory exists
	if _, err := os.Stat(schemasDir); os.IsNotExist(err) {
		return fmt.Errorf("schema directory %s does not exist - please ensure XSD files are placed in xsd/ directory", schemasDir)
	}
	
	// Check if main schema file exists (try both original name and underscore version)
	mainSchemaPath := filepath.Join(schemasDir, spec.mainFile)
	if _, err := os.Stat(mainSchemaPath); os.IsNotExist(err) {
		// Try with underscores
		mainSchemaPath = filepath.Join(schemasDir, strings.ReplaceAll(spec.mainFile, "-", "_"))
		if _, err := os.Stat(mainSchemaPath); os.IsNotExist(err) {
			return fmt.Errorf("main schema file not found (tried %s and %s)", 
				filepath.Join(schemasDir, spec.mainFile),
				filepath.Join(schemasDir, strings.ReplaceAll(spec.mainFile, "-", "_")))
		}
	}
	
	log.Printf("Found schemas in %s", schemasDir)
	return nil
}

func convertToProto(spec struct{ name, version, mainFile string }) error {
	// TODO: Implement XSD parsing and proto generation
	log.Printf("Converting %s v%s schemas to .proto files...", spec.name, spec.version)
	
	// Create output directory
	protoDir := filepath.Join("proto", spec.name+"v"+spec.version)
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		return fmt.Errorf("failed to create proto directory: %v", err)
	}
	
	// For now, just create a placeholder
	protoFile := filepath.Join(protoDir, spec.name+".proto")
	placeholder := fmt.Sprintf(`syntax = "proto3";

package ddex.%s.v%s;

option go_package = "github.com/alecsavvy/ddex-go/gen/%sv%s";

// TODO: Generated from %s
// This is a placeholder - actual XSD parsing not yet implemented

message Placeholder {
  string message = 1;
}
`, spec.name, spec.version, spec.name, spec.version, spec.mainFile)
	
	if err := os.WriteFile(protoFile, []byte(placeholder), 0644); err != nil {
		return fmt.Errorf("failed to write proto file: %v", err)
	}
	
	log.Printf("Created placeholder proto file: %s", protoFile)
	return nil
}