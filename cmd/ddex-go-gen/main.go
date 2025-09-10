package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func repoRootDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func downloadSchemas(spec struct{ name, version, mainFile string }) error {
	repoRoot := repoRootDir()
	schemasDir := filepath.Join(repoRoot, "tmp", spec.name+"v"+spec.version)
	os.MkdirAll(schemasDir, 0755)
	
	// Track what we've downloaded to avoid duplicates
	seen := make(map[string]bool)
	
	// Start with the main schema
	mainURL := fmt.Sprintf("https://service.ddex.net/xml/%s/%s/%s", spec.name, spec.version, spec.mainFile)
	return downloadRecursive(mainURL, schemasDir, seen)
}

func downloadRecursive(urlStr, schemasDir string, seen map[string]bool) error {
	if seen[urlStr] {
		return nil
	}
	seen[urlStr] = true
	
	// Download this file
	filename := filepath.Base(urlStr)
	filename = strings.ReplaceAll(filename, "-", "_") // Convert to Go-friendly names
	localPath := filepath.Join(schemasDir, filename)
	
	if err := downloadFile(urlStr, localPath); err != nil {
		return err
	}
	
	// Read and parse for dependencies
	content, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}
	
	// Extract schema locations from imports/includes
	imports := extractSchemaLocations(string(content))
	
	// Download dependencies recursively
	baseURL := urlStr[:strings.LastIndex(urlStr, "/")]
	for _, schemaLoc := range imports {
		var depURL string
		if strings.HasPrefix(schemaLoc, "http") {
			// Absolute URL
			depURL = schemaLoc
		} else {
			// Relative URL - resolve against base
			depURL = resolveURL(baseURL, schemaLoc)
		}
		
		if err := downloadRecursive(depURL, schemasDir, seen); err != nil {
			log.Printf("Warning: failed to download dependency %s: %v", depURL, err)
		}
	}
	
	// Fix schema locations in this file to point to local files
	fixedContent := fixSchemaLocations(string(content))
	return os.WriteFile(localPath, []byte(fixedContent), 0644)
}

func extractSchemaLocations(content string) []string {
	var locations []string
	
	// Simple regex-like extraction of schemaLocation attributes
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "schemaLocation=") {
			start := strings.Index(line, `schemaLocation="`)
			if start == -1 {
				continue
			}
			start += len(`schemaLocation="`)
			end := strings.Index(line[start:], `"`)
			if end == -1 {
				continue
			}
			schemaLoc := line[start : start+end]
			if schemaLoc != "" {
				locations = append(locations, schemaLoc)
			}
		}
	}
	
	return locations
}

func resolveURL(baseURL, relativePath string) string {
	// Handle relative paths like ../../../../ddex.net/xml/...
	if strings.HasPrefix(relativePath, "../") {
		// For now, just handle the common ddex.net case
		if strings.Contains(relativePath, "ddex.net") {
			return "http://" + relativePath[strings.Index(relativePath, "ddex.net"):]
		}
	}
	return baseURL + "/" + relativePath
}

func fixSchemaLocations(content string) string {
	// Replace remote schema locations with local filenames
	fixed := content
	
	// Extract all schemaLocation values and replace with just the filename
	locations := extractSchemaLocations(content)
	for _, loc := range locations {
		if strings.Contains(loc, "/") {
			// Get just the filename from the path
			filename := filepath.Base(loc)
			filename = strings.ReplaceAll(filename, "-", "_") // Convert to Go-friendly names
			fixed = strings.ReplaceAll(fixed, `"`+loc+`"`, `"`+filename+`"`)
		}
	}
	
	return fixed
}

func downloadFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = io.Copy(file, resp.Body)
	log.Printf("Downloaded %s", filepath.Base(path))
	return err
}

func cleanupTempFiles(spec struct{ name, version, mainFile string }) error {
	repoRoot := repoRootDir()
	schemasDir := filepath.Join(repoRoot, "tmp", spec.name+"v"+spec.version)
	
	if err := os.RemoveAll(schemasDir); err != nil {
		return fmt.Errorf("failed to cleanup temp directory %s: %v", schemasDir, err)
	}
	
	log.Printf("Cleaned up temp directory: %s", schemasDir)
	return nil
}

func generatePackage(spec struct{ name, version, mainFile string }) error {
	repoRoot := repoRootDir()
	schemasDir := filepath.Join(repoRoot, "tmp", spec.name+"v"+spec.version)
	packageName := spec.name + "v" + spec.version
	outDir := filepath.Join(repoRoot, packageName)
	
	// Clean and create output directory
	os.RemoveAll(outDir)
	os.MkdirAll(outDir, 0755)
	
	// Generate with xgen using directory input
	cmd := exec.Command("xgen", "-i", schemasDir, "-o", outDir, "-l", "Go", "-p", packageName)
	cmd.Dir = repoRoot
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
	
	fmt.Printf("Generated %s package:\n", packageName)
	files, _ := os.ReadDir(outDir)
	for _, file := range files {
		fmt.Printf("  %s/%s\n", packageName, file.Name())
	}
	
	return nil
}

func main() {
	for _, spec := range specs {
		log.Printf("Processing %s v%s...", spec.name, spec.version)
		
		if err := downloadSchemas(spec); err != nil {
			log.Fatalf("Failed to download %s v%s: %v", spec.name, spec.version, err)
		}
		
		if err := generatePackage(spec); err != nil {
			log.Fatalf("Failed to generate %s v%s: %v", spec.name, spec.version, err)
		}
		
		if err := cleanupTempFiles(spec); err != nil {
			log.Printf("Warning: failed to cleanup temp files for %s v%s: %v", spec.name, spec.version, err)
		}
	}
}