# DDEX Go Library Makefile

.PHONY: test testdata clean generate help

# Default target
help:
	@echo "Available targets:"
	@echo "  test      - Run tests (downloads testdata if needed)"
	@echo "  testdata  - Download and extract DDEX sample files"
	@echo "  generate  - Generate Go code from DDEX XSD schemas"
	@echo "  clean     - Clean generated files and test data"

# Run tests, ensuring testdata exists
test: testdata
	go test -v ./...

# Download and extract DDEX sample files for testing (only if missing)
testdata:
	@if [ ! -d "testdata/ernv432/Samples43" ]; then \
		echo "Downloading DDEX sample files..."; \
		mkdir -p testdata/ernv432 testdata/meadv11 testdata/piev10; \
		cd testdata/ernv432 && \
		curl -L -o samples.zip "https://service.ddex.net/doc/Standards/ERN43/Samples43.zip" && \
		unzip -q samples.zip && \
		rm samples.zip; \
		echo "Test data ready in testdata/"; \
	else \
		echo "Test data already exists"; \
	fi

# Generate Go code from DDEX schemas
generate:
	go run ./cmd/ddx-go-gen

# Clean up generated files and test data
clean:
	rm -rf testdata/
	rm -rf tmp/

# Force regenerate test data
testdata-refresh: clean testdata