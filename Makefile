# DDEX Go Library Makefile

.PHONY: test testdata clean generate-go generate-proto generate-all help

# Default target
help:
	@echo "DDEX Go Library - Makefile targets:"
	@echo ""
	@echo "Generation:"
	@echo "  generate-go    - Generate pure Go structs from XSD (ddex/ directory)"
	@echo "  generate-proto - Generate .proto files from XSD (proto/ directory)"  
	@echo "  generate-all   - Generate both Go and proto files"
	@echo ""
	@echo "Testing:"
	@echo "  test          - Run all tests (downloads testdata if needed)"
	@echo "  test-roundtrip - Test XML roundtrip between ddex/ and gen/"
	@echo "  testdata      - Download DDEX sample files"
	@echo ""
	@echo "Maintenance:"
	@echo "  clean         - Clean generated files and test data"
	@echo "  testdata-refresh - Force re-download test data"

# Generate pure Go structs from XSD
generate-go:
	@echo "Generating pure Go structs from XSD..."
	go run tools/xsd2go/main.go

# Generate proto files from XSD
generate-proto:
	@echo "Generating proto files from XSD..."
	go run tools/xsd2proto/main.go

# Generate everything
generate-all: generate-go generate-proto
	@echo "All generation complete!"

# Run tests, ensuring testdata exists
test:
	go test -v ./...

# Test roundtrip compatibility between pure Go and proto-generated Go
test-roundtrip:
	go test -v ./test/roundtrip/...

# Clean up generated files and test data
clean:
	rm -rf ddex/ernv* ddex/meadv* ddex/piev*
	rm -rf gen/ernv* gen/meadv* gen/piev*  
	rm -rf proto/ernv*/*.proto proto/meadv*/*.proto proto/piev*/*.proto
	rm -rf testdata/
	rm -rf tmp/
