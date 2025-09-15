# DDEX XSD Schema Files

This directory contains the official DDEX XSD schema files used for code generation.

## Schema Sources

All schemas are downloaded directly from the DDEX consortium:

### ERN (Electronic Release Notification) v4.3.2
- **Source**: https://service.ddex.net/xml/ern/432/
- **Main Schema**: release-notification.xsd
- **Downloaded**: 2025-01-15

### MEAD (Media Enrichment and Description) v1.1  
- **Source**: https://service.ddex.net/xml/mead/11/
- **Main Schema**: media-enrichment-and-description.xsd
- **Downloaded**: 2025-01-15

### PIE (Party Identification and Enrichment) v1.0
- **Source**: https://service.ddex.net/xml/pie/10/
- **Main Schema**: party-identification-and-enrichment.xsd
- **Downloaded**: 2025-01-15

## Directory Structure

```
xsd/
├── ernv432/                    # ERN v4.3.2 schemas
│   ├── release-notification.xsd
│   ├── avs.xsd                 # Allowed Value Sets
│   ├── ddex.xsd               # Common DDEX types
│   └── ... (other dependencies)
│
├── meadv11/                    # MEAD v1.1 schemas  
│   ├── media-enrichment-and-description.xsd
│   └── ... (dependencies)
│
└── piev10/                     # PIE v1.0 schemas
    ├── party-identification-and-enrichment.xsd
    └── ... (dependencies)
```

## Schema Processing

- **Filename normalization**: Hyphens converted to underscores for Go compatibility
- **Local references**: Schema location attributes updated to reference local files
- **No modifications**: Schemas are kept as close to original as possible

## Updating Schemas

To update to newer versions:

1. Download new schemas from DDEX service URLs
2. Update version directories (e.g., `ernv433/` for ERN v4.3.3)
3. Update this README with new sources and dates
4. Update `tools/xsd2proto/main.go` specs array for new versions
5. Regenerate code with `make generate`

## License

These XSD files are property of the DDEX consortium. See https://ddex.net for license terms.