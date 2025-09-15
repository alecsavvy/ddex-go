package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Find all generated protobuf packages
	err := filepath.Walk("gen", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".pb.go") {
			packageDir := filepath.Dir(path)
			packageName := filepath.Base(packageDir)

			// Parse the .pb.go file to find enum types and message types
			enums, err := findEnumTypes(path)
			if err != nil {
				return fmt.Errorf("parsing %s: %w", path, err)
			}

			messages, err := findMessageTypes(path)
			if err != nil {
				return fmt.Errorf("parsing messages %s: %w", path, err)
			}

			// Generate enum strings file if there are enums
			if len(enums) > 0 {
				err = generateEnumStringsFile(packageDir, packageName, enums)
				if err != nil {
					return fmt.Errorf("generating enum strings file for %s: %w", packageDir, err)
				}
				log.Printf("Generated enum_strings.go for package %s with %d enums", packageName, len(enums))
			}

			// Generate single XML file for all messages in the package
			if len(messages) > 0 {
				err = generatePackageXMLFile(packageDir, packageName, messages)
				if err != nil {
					return fmt.Errorf("generating XML file for package %s: %w", packageDir, err)
				}
				log.Printf("Generated %s.xml.go for package %s with %d messages", packageName, packageName, len(messages))
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// findEnumTypes parses a .pb.go file and extracts enum type information
func findEnumTypes(filename string) ([]EnumInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var enums []EnumInfo

	// Look for enum type definitions and their constants
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if ident, ok := ts.Type.(*ast.Ident); ok && ident.Name == "int32" {
							// Found an enum type - now find its constants
							enumName := ts.Name.Name
							constants := findEnumConstants(node, enumName)
							if len(constants) > 0 {
								enums = append(enums, EnumInfo{
									Name:      enumName,
									Constants: constants,
								})
							}
						}
					}
				}
			}
		}
	}

	return enums, nil
}

// findEnumConstants finds all constants for a given enum type
func findEnumConstants(node *ast.File, enumTypeName string) []string {
	var constants []string

	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					// Check if this constant is of our enum type
					if ident, ok := valueSpec.Type.(*ast.Ident); ok && ident.Name == enumTypeName {
						for _, name := range valueSpec.Names {
							constants = append(constants, name.Name)
						}
					}
				}
			}
		}
	}

	return constants
}

type EnumInfo struct {
	Name      string
	Constants []string
}

type MessageInfo struct {
	Name string
}

// findMessageTypes parses a .pb.go file and extracts main message types
func findMessageTypes(filename string) ([]MessageInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var messages []MessageInfo

	// Look for main message type definitions (ones ending with "Message")
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if _, ok := ts.Type.(*ast.StructType); ok {
							// Found a struct type - check if it's a main message type
							messageName := ts.Name.Name
							if strings.HasSuffix(messageName, "Message") {
								messages = append(messages, MessageInfo{
									Name: messageName,
								})
							}
						}
					}
				}
			}
		}
	}

	return messages, nil
}

// generateEnumStringsFile creates an enum_strings.go file with String() methods and parsers
func generateEnumStringsFile(packageDir, packageName string, enums []EnumInfo) error {
	content := generateEnumStringsContent(packageName, enums)

	enumStringsPath := filepath.Join(packageDir, "enum_strings.go")
	return os.WriteFile(enumStringsPath, []byte(content), 0644)
}

// generatePackageXMLFile creates a single XML file for all messages in a package
func generatePackageXMLFile(packageDir, packageName string, messages []MessageInfo) error {
	content := generatePackageXMLContent(packageDir, packageName, messages)

	xmlFileName := packageName + ".xml.go"
	xmlPath := filepath.Join(packageDir, xmlFileName)
	return os.WriteFile(xmlPath, []byte(content), 0644)
}

// generateEnumStringsContent creates the content for enum_strings.go
func generateEnumStringsContent(packageName string, enums []EnumInfo) string {
	var sb strings.Builder

	// Package header
	sb.WriteString(fmt.Sprintf("// Code generated by generate-go-extensions. DO NOT EDIT.\n\n"))
	sb.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	if len(enums) > 0 {
		sb.WriteString("import \"strings\"\n\n")
	}

	// Generate String() methods and parsers for each enum
	// These allow developers to use type-safe enum constants with string fields
	for _, enum := range enums {
		sb.WriteString(generateEnumStringMethod(enum))
		sb.WriteString("\n\n")
		sb.WriteString(generateEnumParser(enum))
		sb.WriteString("\n\n")
	}

	return sb.String()
}

// NamespaceInfo holds namespace configuration for a package
type NamespaceInfo struct {
	Namespace       string
	NamespacePrefix string
	SchemaFile      string
}

// deriveNamespaceInfo extracts namespace info from package directory path
func deriveNamespaceInfo(packageDir string) *NamespaceInfo {
	// packageDir is something like "gen/ddex/ern/v432"
	// We want to extract: ddex type (ern/mead/pie), version (432/43/11/10)

	parts := strings.Split(filepath.Clean(packageDir), string(filepath.Separator))
	if len(parts) < 4 {
		return nil // Not a DDEX package
	}

	// Look for the ddex directory and extract type/version
	ddexIndex := -1
	for i, part := range parts {
		if part == "ddex" {
			ddexIndex = i
			break
		}
	}

	if ddexIndex == -1 || ddexIndex+2 >= len(parts) {
		return nil // Not found or not enough parts
	}

	messageType := parts[ddexIndex+1] // ern, mead, pie
	version := parts[ddexIndex+2]     // v432, v43, v11, etc.

	// Remove 'v' prefix from version
	versionNumber := strings.TrimPrefix(version, "v")

	info := &NamespaceInfo{
		NamespacePrefix: messageType,
	}

	// Set namespace and schema file based on type
	switch messageType {
	case "ern":
		info.Namespace = fmt.Sprintf("http://ddex.net/xml/ern/%s", versionNumber)
		info.SchemaFile = "release-notification.xsd"
	case "mead":
		info.Namespace = fmt.Sprintf("http://ddex.net/xml/mead/%s", versionNumber)
		info.SchemaFile = "media-enrichment-and-description.xsd"
	case "pie":
		info.Namespace = fmt.Sprintf("http://ddex.net/xml/pie/%s", versionNumber)
		info.SchemaFile = "party-identification-and-enrichment.xsd"
	default:
		return nil
	}

	return info
}

// generatePackageXMLContent creates the content for a package XML file
func generatePackageXMLContent(packageDir, packageName string, messages []MessageInfo) string {
	var sb strings.Builder

	// Package header
	sb.WriteString(fmt.Sprintf("// Code generated by generate-go-extensions. DO NOT EDIT.\n\n"))
	sb.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	sb.WriteString("import \"encoding/xml\"\n\n")

	// Derive namespace info from package path
	nsInfo := deriveNamespaceInfo(packageDir)
	if nsInfo != nil {
		sb.WriteString("// Package-level namespace constants\n")
		sb.WriteString("const (\n")
		sb.WriteString(fmt.Sprintf("\tNamespace = \"%s\"\n", nsInfo.Namespace))
		sb.WriteString(fmt.Sprintf("\tNamespacePrefix = \"%s\"\n", nsInfo.NamespacePrefix))
		sb.WriteString(fmt.Sprintf("\tSchemaLocation = \"%s %s/%s\"\n", nsInfo.Namespace, nsInfo.Namespace, nsInfo.SchemaFile))
		sb.WriteString("\tNamespaceXSI = \"http://www.w3.org/2001/XMLSchema-instance\"\n")
		sb.WriteString(")\n\n")
	}

	// Generate XML marshaling methods for all messages in the package
	for i, message := range messages {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(generateXMLMarshalingMethods(message, nsInfo))
	}

	return sb.String()
}

// generateEnumStringMethod creates a String() method for the enum type
func generateEnumStringMethod(enum EnumInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// XMLString returns the XML string representation of %s\n", enum.Name))
	sb.WriteString(fmt.Sprintf("func (e %s) XMLString() string {\n", enum.Name))
	sb.WriteString("\tswitch e {\n")

	// Generate cases for each constant
	for _, constant := range enum.Constants {
		if strings.HasSuffix(constant, "_UNSPECIFIED") {
			continue // Skip UNSPECIFIED values
		}

		// Extract the meaningful part of the constant name
		upperName := strings.ToUpper(enum.Name)
		idx := strings.LastIndex(constant, upperName+"_")
		if idx >= 0 {
			afterPrefix := constant[idx+len(upperName)+1:]
			if afterPrefix != "" && afterPrefix != "UNSPECIFIED" {
				sb.WriteString(fmt.Sprintf("\tcase %s:\n", constant))
				sb.WriteString(fmt.Sprintf("\t\treturn \"%s\"\n", afterPrefix))
			}
		}
	}

	sb.WriteString("\tdefault:\n")
	sb.WriteString("\t\treturn \"\"\n")
	sb.WriteString("\t}\n")
	sb.WriteString("}")

	return sb.String()
}

// generateEnumParser creates the parser function for an enum
func generateEnumParser(enum EnumInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// Parse%sString parses a string value to %s enum (case-insensitive)\n", enum.Name, enum.Name))
	sb.WriteString(fmt.Sprintf("func Parse%sString(s string) (%s, bool) {\n", enum.Name, enum.Name))
	sb.WriteString("\ts = strings.ToUpper(s)\n")
	sb.WriteString("\tswitch s {\n")

	// Generate cases for each constant
	for _, constant := range enum.Constants {
		if strings.HasSuffix(constant, "_UNSPECIFIED") {
			continue // Skip UNSPECIFIED values
		}

		// Extract the meaningful part of the constant name
		// Try to find the enum pattern: EnumName_ENUM_NAME_VALUE
		// We'll look for the last occurrence of the enum name in uppercase
		upperName := strings.ToUpper(enum.Name)

		// Find the pattern EnumName_..._VALUE
		idx := strings.LastIndex(constant, upperName+"_")
		if idx >= 0 {
			// Skip past "EnumName_..._" to get the value part
			afterPrefix := constant[idx+len(upperName)+1:]
			if afterPrefix != "" && afterPrefix != "UNSPECIFIED" {
				sb.WriteString(fmt.Sprintf("\tcase \"%s\":\n", afterPrefix))
				sb.WriteString(fmt.Sprintf("\t\treturn %s, true\n", constant))
			}
		}
	}

	sb.WriteString("\tdefault:\n")
	sb.WriteString(fmt.Sprintf("\t\treturn %s(0), false\n", enum.Name))
	sb.WriteString("\t}\n")
	sb.WriteString("}")

	return sb.String()
}

// generateXMLMarshalingMethods creates MarshalXML and UnmarshalXML methods for message types
func generateXMLMarshalingMethods(message MessageInfo, nsInfo *NamespaceInfo) string {
	var sb strings.Builder

	// Generate MarshalXML method
	sb.WriteString(fmt.Sprintf("// MarshalXML implements xml.Marshaler for %s\n", message.Name))
	sb.WriteString(fmt.Sprintf("func (m *%s) MarshalXML(e *xml.Encoder, start xml.StartElement) error {\n", message.Name))

	// Add namespace population for root message types if we have namespace info
	if nsInfo != nil && isRootMessage(message.Name) {
		sb.WriteString("\t// Set default namespace values if empty\n")

		// Generate field name based on prefix (XmlnsErn, XmlnsMead, XmlnsPie)
		fieldName := fmt.Sprintf("Xmlns%s", strings.Title(nsInfo.NamespacePrefix))
		sb.WriteString(fmt.Sprintf("\tif m.%s == \"\" {\n", fieldName))
		sb.WriteString("\t\tm." + fieldName + " = Namespace\n")
		sb.WriteString("\t}\n")

		sb.WriteString("\tif m.XmlnsXsi == \"\" {\n")
		sb.WriteString("\t\tm.XmlnsXsi = NamespaceXSI\n")
		sb.WriteString("\t}\n")

		sb.WriteString("\tif m.XsiSchemaLocation == \"\" {\n")
		sb.WriteString("\t\tm.XsiSchemaLocation = SchemaLocation\n")
		sb.WriteString("\t}\n\n")
	}

	sb.WriteString("\t// Create an alias type to avoid infinite recursion\n")
	sb.WriteString(fmt.Sprintf("\ttype alias %s\n", message.Name))
	sb.WriteString("\treturn e.EncodeElement((*alias)(m), start)\n")
	sb.WriteString("}\n\n")

	// Generate UnmarshalXML method
	sb.WriteString(fmt.Sprintf("// UnmarshalXML implements xml.Unmarshaler for %s\n", message.Name))
	sb.WriteString(fmt.Sprintf("func (m *%s) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {\n", message.Name))
	sb.WriteString("\t// Create an alias type to avoid infinite recursion\n")
	sb.WriteString(fmt.Sprintf("\ttype alias %s\n", message.Name))
	sb.WriteString("\treturn d.DecodeElement((*alias)(m), &start)\n")
	sb.WriteString("}")

	return sb.String()
}

// isRootMessage determines if a message type is a root message that needs namespace handling
func isRootMessage(messageName string) bool {
	switch messageName {
	case "NewReleaseMessage", "PurgeReleaseMessage", "MeadMessage", "PieMessage", "PieRequestMessage":
		return true
	default:
		return false
	}
}
