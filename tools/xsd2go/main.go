package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DDEX specifications to generate
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
		log.Printf("Processing %s v%s...", spec.name, spec.version)
		
		if err := validateSchemas(spec); err != nil {
			log.Fatalf("Schema validation failed for %s v%s: %v", spec.name, spec.version, err)
		}
		
		if err := generatePackage(spec); err != nil {
			log.Fatalf("Failed to generate %s v%s: %v", spec.name, spec.version, err)
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


func generatePackage(spec struct{ name, version, mainFile string }) error {
	schemasDir := filepath.Join("xsd", spec.name+"v"+spec.version)
	packageName := spec.name + "v" + spec.version
	outDir := filepath.Join("ddex", packageName)
	
	// Clean and create output directory
	os.RemoveAll(outDir)
	os.MkdirAll(outDir, 0755)
	
	// Generate with xgen using directory input
	cmd := exec.Command("xgen", "-i", schemasDir, "-o", outDir, "-l", "Go", "-p", packageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("xgen output: %s", string(output))
		return fmt.Errorf("xgen failed: %v", err)
	}
	
	// Keep only schema files, remove xgen artifacts and allowed_value_sets
	outEntries, _ := os.ReadDir(outDir)
	for _, entry := range outEntries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".xsd.go") || name == "allowed_value_sets.xsd.go" {
			os.RemoveAll(filepath.Join(outDir, name))
		}
	}
	
	fmt.Printf("Generated %s package in ddex/%s:\n", packageName, packageName)
	files, _ := os.ReadDir(outDir)
	for _, file := range files {
		fmt.Printf("  %s\n", file.Name())
	}
	
	return nil
}