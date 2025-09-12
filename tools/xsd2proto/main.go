package main

import (
	"encoding/xml"
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

// XSD Schema parsing structures - tailored for DDEX patterns
type XSDSchema struct {
	XMLName           xml.Name         `xml:"schema"`
	TargetNamespace   string           `xml:"targetNamespace,attr"`
	Elements          []XSDElement     `xml:"element"`
	ComplexTypes      []XSDComplexType `xml:"complexType"`
	SimpleTypes       []XSDSimpleType  `xml:"simpleType"`
}

type XSDElement struct {
	Name        string           `xml:"name,attr"`
	Type        string           `xml:"type,attr"`
	MinOccurs   string           `xml:"minOccurs,attr"`
	MaxOccurs   string           `xml:"maxOccurs,attr"`
	ComplexType *XSDComplexType  `xml:"complexType"`
}

type XSDComplexType struct {
	Name          string            `xml:"name,attr"`
	Sequence      *XSDSequence      `xml:"sequence"`
	Choice        *XSDChoice        `xml:"choice"`
	SimpleContent *XSDSimpleContent `xml:"simpleContent"`
	Attributes    []XSDAttribute    `xml:"attribute"`
}

type XSDSequence struct {
	Elements []XSDElement `xml:"element"`
}

type XSDChoice struct {
	MinOccurs string       `xml:"minOccurs,attr"`
	MaxOccurs string       `xml:"maxOccurs,attr"`
	Elements  []XSDElement `xml:"element"`
}

type XSDSimpleContent struct {
	Extension *XSDExtension `xml:"extension"`
}

type XSDExtension struct {
	Base       string         `xml:"base,attr"`
	Attributes []XSDAttribute `xml:"attribute"`
}

type XSDAttribute struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Use  string `xml:"use,attr"`
}

type XSDSimpleType struct {
	Name        string         `xml:"name,attr"`
	Restriction *XSDRestriction `xml:"restriction"`
}

type XSDRestriction struct {
	Base         string           `xml:"base,attr"`
	Enumerations []XSDEnumeration `xml:"enumeration"`
}

type XSDEnumeration struct {
	Value string `xml:"value,attr"`
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
	log.Printf("Converting %s v%s schemas to .proto files...", spec.name, spec.version)
	
	// Parse the main XSD schema
	schemasDir := filepath.Join("xsd", spec.name+"v"+spec.version)
	schemaFile := filepath.Join(schemasDir, spec.mainFile)
	
	schema, err := parseXSDFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to parse XSD file %s: %v", schemaFile, err)
	}
	
	// Parse the allowed-value-sets.xsd file to get enum definitions
	allowedValueSetsFile := filepath.Join("xsd", "allowed-value-sets.xsd")
	if _, err := os.Stat(allowedValueSetsFile); err == nil {
		avsSchema, err := parseXSDFile(allowedValueSetsFile)
		if err != nil {
			return fmt.Errorf("failed to parse allowed-value-sets.xsd: %v", err)
		}
		// Merge the schemas - append types from allowed-value-sets
		schema.ComplexTypes = append(schema.ComplexTypes, avsSchema.ComplexTypes...)
		schema.SimpleTypes = append(schema.SimpleTypes, avsSchema.SimpleTypes...)
		schema.Elements = append(schema.Elements, avsSchema.Elements...)
		log.Printf("Merged allowed-value-sets.xsd: +%d complex types, +%d simple types, +%d elements", 
			len(avsSchema.ComplexTypes), len(avsSchema.SimpleTypes), len(avsSchema.Elements))
	}
	
	// Create output directory
	protoDir := filepath.Join("proto", spec.name+"v"+spec.version)
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		return fmt.Errorf("failed to create proto directory: %v", err)
	}
	
	// Generate proto file
	protoFile := filepath.Join(protoDir, spec.name+".proto")
	protoContent, err := generateProtoFromXSD(schema, spec)
	if err != nil {
		return fmt.Errorf("failed to generate proto content: %v", err)
	}
	
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		return fmt.Errorf("failed to write proto file: %v", err)
	}
	
	log.Printf("Generated proto file: %s", protoFile)
	return nil
}

func parseXSDFile(filePath string) (*XSDSchema, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read XSD file: %v", err)
	}
	
	var schema XSDSchema
	if err := xml.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse XSD: %v", err)
	}
	
	log.Printf("Parsed XSD: %d elements, %d complex types, %d simple types", 
		len(schema.Elements), len(schema.ComplexTypes), len(schema.SimpleTypes))
	
	return &schema, nil
}

func generateProtoFromXSD(schema *XSDSchema, spec struct{ name, version, mainFile string }) (string, error) {
	var builder strings.Builder
	
	// Proto header
	builder.WriteString(`syntax = "proto3";

`)
	builder.WriteString(fmt.Sprintf("package ddex.%s.v%s;\n\n", spec.name, spec.version))
	builder.WriteString(fmt.Sprintf("option go_package = \"github.com/alecsavvy/ddex-go/gen/%sv%s\";\n\n", spec.name, spec.version))
	builder.WriteString("import \"tagger/tagger.proto\";\n\n")
	
	builder.WriteString(fmt.Sprintf("// Generated from %s\n", spec.mainFile))
	builder.WriteString(fmt.Sprintf("// Target namespace: %s\n\n", schema.TargetNamespace))
	
	// Generate messages for top-level elements
	for _, element := range schema.Elements {
		if element.ComplexType != nil {
			// Inline complex type
			msgContent, err := generateComplexTypeMessage(element.Name, element.ComplexType)
			if err != nil {
				return "", fmt.Errorf("failed to generate message for element %s: %v", element.Name, err)
			}
			builder.WriteString(msgContent)
			builder.WriteString("\n")
		}
	}
	
	// Generate messages for named complex types
	for _, complexType := range schema.ComplexTypes {
		msgContent, err := generateComplexTypeMessage(complexType.Name, &complexType)
		if err != nil {
			return "", fmt.Errorf("failed to generate message for complex type %s: %v", complexType.Name, err)
		}
		builder.WriteString(msgContent)
		builder.WriteString("\n")
	}
	
	// Generate enums for simple types with restrictions
	for _, simpleType := range schema.SimpleTypes {
		if simpleType.Restriction != nil && len(simpleType.Restriction.Enumerations) > 0 {
			enumContent := generateEnum(simpleType)
			builder.WriteString(enumContent)
			builder.WriteString("\n")
		}
	}
	
	return builder.String(), nil
}

func generateComplexTypeMessage(name string, complexType *XSDComplexType) (string, error) {
	var builder strings.Builder
	
	messageName := toProtoMessageName(name)
	builder.WriteString(fmt.Sprintf("message %s {\n", messageName))
	
	fieldNum := 1
	
	// Handle sequence elements
	if complexType.Sequence != nil {
		for _, element := range complexType.Sequence.Elements {
			field, err := generateField(element, fieldNum)
			if err != nil {
				return "", fmt.Errorf("failed to generate field for element %s: %v", element.Name, err)
			}
			builder.WriteString("  " + field + "\n")
			fieldNum++
		}
	}
	
	// Handle choice elements (use oneof)
	if complexType.Choice != nil {
		builder.WriteString("  oneof choice {\n")
		for _, element := range complexType.Choice.Elements {
			field, err := generateField(element, fieldNum)
			if err != nil {
				return "", fmt.Errorf("failed to generate choice field for element %s: %v", element.Name, err)
			}
			builder.WriteString("    " + field + "\n")
			fieldNum++
		}
		builder.WriteString("  }\n")
	}
	
	// Handle simple content with extension (attributes + value)
	if complexType.SimpleContent != nil && complexType.SimpleContent.Extension != nil {
		// Add value field for the simple content
		xmlTag := generateXMLTag("", false, true) // chardata
		builder.WriteString(fmt.Sprintf("  string value = %d %s;\n", fieldNum, xmlTag))
		fieldNum++
		
		// Add attribute fields
		for _, attr := range complexType.SimpleContent.Extension.Attributes {
			field := generateAttributeField(attr, fieldNum)
			builder.WriteString("  " + field + "\n")
			fieldNum++
		}
	}
	
	// Handle attributes
	for _, attr := range complexType.Attributes {
		field := generateAttributeField(attr, fieldNum)
		builder.WriteString("  " + field + "\n")
		fieldNum++
	}
	
	builder.WriteString("}")
	return builder.String(), nil
}

func generateField(element XSDElement, fieldNum int) (string, error) {
	fieldName := toProtoFieldName(element.Name)
	
	// Handle elements without explicit type (inline types)
	fieldType := "string" // Default for elements with inline types
	if element.Type != "" {
		fieldType = xsdTypeToProto(element.Type)
	}
	
	// Handle cardinality
	repeated := ""
	if element.MaxOccurs == "unbounded" {
		repeated = "repeated "
	}
	
	// Generate XML tag
	xmlTag := generateXMLTag(element.Name, false, false)
	
	return fmt.Sprintf("%s%s %s = %d %s;", repeated, fieldType, fieldName, fieldNum, xmlTag), nil
}

func generateAttributeField(attr XSDAttribute, fieldNum int) string {
	fieldName := toProtoFieldName(attr.Name)
	
	// If no type is specified, default to string (XSD default for attributes)
	fieldType := "string" // XSD default
	if attr.Type != "" {
		fieldType = xsdTypeToProto(attr.Type)
	}
	
	xmlTag := generateXMLTag(attr.Name, true, false) // is attribute
	
	return fmt.Sprintf("%s %s = %d %s;", fieldType, fieldName, fieldNum, xmlTag)
}

func generateXMLTag(name string, isAttr, isCharData bool) string {
	if isCharData {
		return `[(tagger.tags) = "xml:\",chardata\""]`
	}
	if isAttr {
		return fmt.Sprintf(`[(tagger.tags) = "xml:\"%s,attr\""]`, name)
	}
	return fmt.Sprintf(`[(tagger.tags) = "xml:\"%s\""]`, name)
}

func generateEnum(simpleType XSDSimpleType) string {
	var builder strings.Builder
	
	enumName := toProtoMessageName(simpleType.Name)
	builder.WriteString(fmt.Sprintf("enum %s {\n", enumName))
	
	for i, enum := range simpleType.Restriction.Enumerations {
		enumValue := toProtoEnumValue(enum.Value)
		builder.WriteString(fmt.Sprintf("  %s = %d;\n", enumValue, i))
	}
	
	builder.WriteString("}")
	return builder.String()
}

// Type mapping functions
func xsdTypeToProto(xsdType string) string {
	// Remove namespace prefix if present
	if idx := strings.Index(xsdType, ":"); idx != -1 {
		xsdType = xsdType[idx+1:]
	}
	
	switch xsdType {
	case "string", "normalizedString", "token", "anyURI":
		return "string"
	case "int", "integer":
		return "int32"
	case "long":
		return "int64"
	case "boolean":
		return "bool"
	case "decimal", "float":
		return "string" // Preserve precision for decimal
	case "double":
		return "double"
	case "dateTime", "date", "time", "duration":
		return "string" // ISO 8601 format strings
	case "base64Binary":
		return "bytes"
	default:
		// Assume it's a custom type/message
		return toProtoMessageName(xsdType)
	}
}

func toProtoFieldName(name string) string {
	// Convert to snake_case
	result := ""
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		result += strings.ToLower(string(r))
	}
	return result
}

func toProtoMessageName(name string) string {
	// Ensure PascalCase
	if name == "" {
		return ""
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func toProtoEnumValue(value string) string {
	// Convert to UPPER_SNAKE_CASE
	result := strings.ToUpper(value)
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, " ", "_")
	return result
}