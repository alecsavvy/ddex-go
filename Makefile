# DDEX Go Library Makefile

.PHONY: test testdata clean generate-go generate-proto generate buf-lint buf-generate buf-all help

# Default target
help:
	@echo "DDEX Go Library - Makefile targets:"
	@echo ""
	@echo "Generation:"
	@echo "  generate-go    - Generate pure Go structs from XSD (ddex/ directory)"
	@echo "  generate-proto - Generate .proto files from XSD (proto/ directory)"  
	@echo "  generate       - Generate both Go and proto files"
	@echo "  buf-lint      - Lint protobuf files with buf"
	@echo "  buf-generate  - Generate Go code from .proto files with buf"
	@echo "  buf-all       - Generate protos from XSD, then Go code from protos"
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
generate: generate-go generate-proto
	@echo "All generation complete!"

# Lint protobuf files with buf
buf-lint:
	@echo "Linting protobuf files..."
	buf lint

# Generate Go code from protobuf files
buf-generate: 
	@echo "Generating Go code from protobuf files..."
	buf generate
	@echo "Injecting XML tags with protoc-go-inject-tag..."
	@$(MAKE) inject-tags
	@echo "Generating XML marshaling for enums..."
	@$(MAKE) generate-marshals

# Inject XML tags into generated protobuf structs using protoc-go-inject-tag
inject-tags:
	@echo "Injecting tags into generated Go files..."
	@for file in $$(find gen -name "*.pb.go" 2>/dev/null); do \
		echo "  Processing $$file..."; \
		protoc-go-inject-tag -input=$$file 2>/dev/null || true; \
	done
	@echo "XML tags injected successfully!"

# Generate XML marshaling wrappers for enum types
generate-marshals:
	@echo "Generating marshal.go files for enum XML handling..."
	go run tools/generate-marshals.go
	@echo "Marshal generation complete!"

# Complete protobuf workflow: XSD -> proto -> Go with XML tags
buf-all: generate-proto buf-lint buf-generate
	@echo "Complete protobuf generation workflow complete!"

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
