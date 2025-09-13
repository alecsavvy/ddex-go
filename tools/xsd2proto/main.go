package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//
// =======================
// Spec list (entry schemas)
// =======================
//

var specs = []struct {
	name     string
	version  string
	mainFile string
}{
	{"ern", "432", "release-notification.xsd"},
	{"mead", "11", "media-enrichment-and-description.xsd"},
	{"pie", "10", "party-identification-and-enrichment.xsd"},
}

//
// =======================
// XSD Models (extended)
// =======================
//

type XSDSchema struct {
	XMLName         xml.Name         `xml:"schema"`
	TargetNamespace string           `xml:"targetNamespace,attr"`
	Elements        []XSDElement     `xml:"element"`
	ComplexTypes    []XSDComplexType `xml:"complexType"`
	SimpleTypes     []XSDSimpleType  `xml:"simpleType"`

	// NEW: follow schema structure
	Imports  []XSDImport  `xml:"import"`
	Includes []XSDInclude `xml:"include"`
}

type XSDImport struct {
	Namespace      string `xml:"namespace,attr"`
	SchemaLocation string `xml:"schemaLocation,attr"`
}

type XSDInclude struct {
	SchemaLocation string `xml:"schemaLocation,attr"`
}

type XSDElement struct {
	Name        string          `xml:"name,attr"`
	Type        string          `xml:"type,attr"`
	MinOccurs   string          `xml:"minOccurs,attr"`
	MaxOccurs   string          `xml:"maxOccurs,attr"`
	ComplexType *XSDComplexType `xml:"complexType"`
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
	Name        string          `xml:"name,attr"`
	Restriction *XSDRestriction `xml:"restriction"`
}

type XSDRestriction struct {
	Base         string           `xml:"base,attr"`
	Enumerations []XSDEnumeration `xml:"enumeration"`
}

type XSDEnumeration struct {
	Value string `xml:"value,attr"`
}

//
// =======================
// Aggregation by namespace
// =======================
//

type NamespaceBundle struct {
	TargetNamespace string

	// Aggregated components (includes merged here).
	Elements     []XSDElement
	ComplexTypes []XSDComplexType
	SimpleTypes  []XSDSimpleType

	// Cross-namespace dependencies discovered via xs:import
	// and via seeing qname types with foreign prefixes (best-effort).
	Imports map[string]struct{} // set of targetNamespace strings
}

// Graph loader state
type loadState struct {
	visitedFiles map[string]struct{} // absolute paths visited
	// Map of targetNamespace → bundle
	nsBundles map[string]*NamespaceBundle
	// file path → schema's ns (helpful for relative includes)
	fileToNS map[string]string
}

func newLoadState() *loadState {
	return &loadState{
		visitedFiles: make(map[string]struct{}),
		nsBundles:    make(map[string]*NamespaceBundle),
		fileToNS:     make(map[string]string),
	}
}

//
// =======================
// Entry
// =======================
//

// 1) Put this near the top-level (package scope), not inside convertSpec:

type protoPkgInfo struct {
	pkgName   string
	goPackage string
	filePath  string // relative to proto root
}

func main() {
	for _, spec := range specs {
		log.Printf("Converting %s v%s to protobuf (namespace-aware)...", spec.name, spec.version)

		if err := validateSchemas(spec); err != nil {
			log.Fatalf("Schema validation failed for %s v%s: %v", spec.name, spec.version, err)
		}

		if err := convertSpec(spec); err != nil {
			log.Fatalf("Failed to convert %s v%s: %v", spec.name, spec.version, err)
		}
	}
}

func validateSchemas(spec struct{ name, version, mainFile string }) error {
	schemasDir := filepath.Join("xsd", spec.name+"v"+spec.version)
	if _, err := os.Stat(schemasDir); os.IsNotExist(err) {
		return fmt.Errorf("schema directory %s does not exist", schemasDir)
	}

	entry := filepath.Join(schemasDir, spec.mainFile)
	if _, err := os.Stat(entry); os.IsNotExist(err) {
		alt := filepath.Join(schemasDir, strings.ReplaceAll(spec.mainFile, "-", "_"))
		if _, err2 := os.Stat(alt); os.IsNotExist(err2) {
			return fmt.Errorf("main schema not found; tried %s and %s", entry, alt)
		}
	}
	return nil
}

//
// =======================
// Conversion pipeline
// =======================
//

func convertSpec(spec struct{ name, version, mainFile string }) error {
	schemasDir := filepath.Join("xsd", spec.name+"v"+spec.version)
	entryPath := filepath.Join(schemasDir, spec.mainFile)
	if _, err := os.Stat(entryPath); os.IsNotExist(err) {
		entryPath = filepath.Join(schemasDir, strings.ReplaceAll(spec.mainFile, "-", "_"))
	}

	st := newLoadState()
	if err := loadSchemaGraph(st, entryPath); err != nil {
		return fmt.Errorf("load graph: %w", err)
	}

	// Create output dir: proto/<spec or inferred>/*
	outRoot := filepath.Join("proto")
	if err := os.MkdirAll(outRoot, 0755); err != nil {
		return err
	}

	// Emit one .proto per namespace bundle.
	// We need deterministic order for stable builds.
	var namespaces []string
	for ns := range st.nsBundles {
		namespaces = append(namespaces, ns)
	}
	sort.Strings(namespaces)

	// Pre-compute package + file paths for imports
	pkgs := make(map[string]protoPkgInfo) // ns → info
	for _, ns := range namespaces {
		pkg := namespaceToProtoPackage(ns, spec)
		goPkg := namespaceToGoPackage(ns, spec)
		path := packageToPath(pkg)
		pkgs[ns] = protoPkgInfo{pkgName: pkg, goPackage: goPkg, filePath: path}
	}

	for _, ns := range namespaces {
		b := st.nsBundles[ns]
		info := pkgs[ns]

		// Ensure directory exists
		dir := filepath.Join(outRoot, filepath.Dir(info.filePath))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Build file content
		content, err := generateProtoForBundle(b, info.pkgName, info.goPackage, pkgs)
		if err != nil {
			return fmt.Errorf("generate for ns %s: %w", ns, err)
		}

		outFile := filepath.Join(outRoot, info.filePath)
		if err := os.WriteFile(outFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("write %s: %w", outFile, err)
		}
		log.Printf("Generated %s", outFile)
	}

	return nil
}

//
// =======================
// Graph loader (includes/imports)
// =======================
//

func loadSchemaGraph(st *loadState, filePath string) error {
	abs, _ := filepath.Abs(filePath)
	if _, ok := st.visitedFiles[abs]; ok {
		return nil
	}
	st.visitedFiles[abs] = struct{}{}

	data, err := os.ReadFile(abs)
	if err != nil {
		return fmt.Errorf("read %s: %w", abs, err)
	}

	var schema XSDSchema
	if err := xml.Unmarshal(data, &schema); err != nil {
		return fmt.Errorf("parse %s: %w", abs, err)
	}

	if schema.TargetNamespace == "" {
		return fmt.Errorf("schema %s missing targetNamespace", abs)
	}

	st.fileToNS[abs] = schema.TargetNamespace

	// Get or create namespace bundle
	b := st.nsBundles[schema.TargetNamespace]
	if b == nil {
		b = &NamespaceBundle{
			TargetNamespace: schema.TargetNamespace,
			Imports:         make(map[string]struct{}),
		}
		st.nsBundles[schema.TargetNamespace] = b
	}

	// Merge components (includes naturally collapse here)
	b.Elements = append(b.Elements, schema.Elements...)
	b.ComplexTypes = append(b.ComplexTypes, schema.ComplexTypes...)
	b.SimpleTypes = append(b.SimpleTypes, schema.SimpleTypes...)

	// Track declared imports by namespace
	for _, imp := range schema.Imports {
		if imp.Namespace != "" && imp.Namespace != schema.TargetNamespace {
			b.Imports[imp.Namespace] = struct{}{}
		}
	}

	// Follow xs:include (same-namespace; relative to this file)
	baseDir := filepath.Dir(abs)
	for _, inc := range schema.Includes {
		if inc.SchemaLocation == "" {
			continue
		}
		next := filepath.Join(baseDir, inc.SchemaLocation)
		if err := loadSchemaGraph(st, next); err != nil {
			return err
		}
	}

	// Follow xs:import where schemaLocation is present; if absent, we still recorded the ns in Imports
	for _, imp := range schema.Imports {
		if imp.SchemaLocation == "" {
			continue
		}
		next := filepath.Join(baseDir, imp.SchemaLocation)
		// If the imported file has a different targetNamespace, it will get its own bundle.
		if err := loadSchemaGraph(st, next); err != nil {
			return err
		}
	}

	return nil
}

//
// =======================
// Codegen per namespace
// =======================
//

func generateProtoForBundle(
	b *NamespaceBundle,
	packageName string,
	goPackage string,
	all map[string]protoPkgInfo,
) (string, error) {

	var sb strings.Builder

	// Header
	sb.WriteString(`syntax = "proto3";` + "\n\n")
	sb.WriteString(fmt.Sprintf("package %s;\n\n", packageName))
	sb.WriteString(fmt.Sprintf("option go_package = \"%s\";\n\n", goPackage))
	sb.WriteString(fmt.Sprintf("// Target namespace: %s\n\n", b.TargetNamespace))

	// Imports (protobuf)
	// Sort for determinism
	var deps []string
	for ns := range b.Imports {
		if ns == b.TargetNamespace {
			continue
		}
		if info, ok := all[ns]; ok {
			deps = append(deps, info.filePath)
		}
	}
	sort.Strings(deps)
	for _, f := range deps {
		// Normalize to POSIX paths in import statements
		sb.WriteString(fmt.Sprintf("import \"%s\";\n", toPosixPath(f)))
	}
	if len(deps) > 0 {
		sb.WriteString("\n")
	}

	// Track generated type names (message & enum in one space) for this package
	generated := make(map[string]struct{})

	// Top-level elements with inline complex types → message
	for _, el := range b.Elements {
		if el.ComplexType != nil {
			name := toProtoMessageName(el.Name)
			if _, exists := generated[name]; !exists {
				msg, err := generateComplexTypeMessage(el.Name, el.ComplexType, all)
				if err != nil {
					return "", err
				}
				sb.WriteString(msg)
				sb.WriteString("\n\n")
				generated[name] = struct{}{}
			}
		}
	}

	// Named complex types → message
	for _, ct := range b.ComplexTypes {
		if ct.Name == "" {
			continue
		}
		name := toProtoMessageName(ct.Name)
		if _, exists := generated[name]; exists {
			continue
		}
		msg, err := generateComplexTypeMessage(ct.Name, &ct, all)
		if err != nil {
			return "", err
		}
		sb.WriteString(msg)
		sb.WriteString("\n\n")
		generated[name] = struct{}{}
	}

	// Simple types with enumerations → enum
	for _, st := range b.SimpleTypes {
		if st.Name == "" || st.Restriction == nil || len(st.Restriction.Enumerations) == 0 {
			continue
		}
		en := toProtoEnumName(st.Name)
		if _, exists := generated[en]; exists {
			continue
		}
		sb.WriteString(generateEnum(st))
		sb.WriteString("\n\n")
		generated[en] = struct{}{}
	}

	return strings.TrimSpace(sb.String()) + "\n", nil
}

//
// =======================
// Field/message/enum codegen (your original logic, kept)
// =======================
//

func generateComplexTypeMessage(name string, complexType *XSDComplexType, allPkgs map[string]protoPkgInfo) (string, error) {
	var builder strings.Builder

	messageName := toProtoMessageName(name)
	builder.WriteString(fmt.Sprintf("message %s {\n", messageName))

	fieldNum := 1

	// sequence → fields
	if complexType.Sequence != nil {
		for _, element := range complexType.Sequence.Elements {
			field, err := generateField(element, fieldNum, allPkgs)
			if err != nil {
				return "", fmt.Errorf("failed to generate field for element %s: %v", element.Name, err)
			}
			builder.WriteString(field + "\n")
			fieldNum++
		}
	}

	// choice → oneof
	if complexType.Choice != nil && len(complexType.Choice.Elements) > 0 {
		builder.WriteString("  oneof choice {\n")
		for _, element := range complexType.Choice.Elements {
			field, err := generateChoiceField(element, fieldNum, allPkgs)
			if err != nil {
				return "", fmt.Errorf("failed to generate choice field for element %s: %v", element.Name, err)
			}
			builder.WriteString(field + "\n")
			fieldNum++
		}
		builder.WriteString("  }\n")
	}

	// simpleContent extension → value + attributes
	if complexType.SimpleContent != nil && complexType.SimpleContent.Extension != nil {
		// chardata value
		injectComment := "  // @gotags: xml:\",chardata\""
		builder.WriteString(fmt.Sprintf("%s\n  string value = %d;\n", injectComment, fieldNum))
		fieldNum++

		// attributes
		for _, attr := range complexType.SimpleContent.Extension.Attributes {
			field := generateAttributeField(attr, fieldNum, allPkgs)
			builder.WriteString(field + "\n")
			fieldNum++
		}
	}

	// attributes on the complexType itself
	for _, attr := range complexType.Attributes {
		field := generateAttributeField(attr, fieldNum, allPkgs)
		builder.WriteString(field + "\n")
		fieldNum++
	}

	builder.WriteString("}")
	return builder.String(), nil
}

func generateField(element XSDElement, fieldNum int, allPkgs map[string]protoPkgInfo) (string, error) {
	fieldName := toProtoFieldName(element.Name)

	// Type mapping
	fieldType := "string" // default
	if element.Type != "" {
		fieldType = xsdTypeToProto(element.Type, allPkgs)
	}

	// Cardinality
	repeated := ""
	if element.MaxOccurs == "unbounded" {
		repeated = "repeated "
	}

	// gotags for xml element name
	injectComment := fmt.Sprintf("  // @gotags: xml:\"%s\"", element.Name)

	return fmt.Sprintf("%s\n  %s%s %s = %d;", injectComment, repeated, fieldType, fieldName, fieldNum), nil
}

func generateChoiceField(element XSDElement, fieldNum int, allPkgs map[string]protoPkgInfo) (string, error) {
	fieldName := toProtoFieldName(element.Name)

	fieldType := "string"
	if element.Type != "" {
		fieldType = xsdTypeToProto(element.Type, allPkgs)
	}

	injectComment := fmt.Sprintf("  // @gotags: xml:\"%s\"", element.Name)
	return fmt.Sprintf("%s\n    %s %s = %d;", injectComment, fieldType, fieldName, fieldNum), nil
}

func generateAttributeField(attr XSDAttribute, fieldNum int, allPkgs map[string]protoPkgInfo) string {
	fieldName := toProtoFieldName(attr.Name)

	fieldType := "string"
	if attr.Type != "" {
		fieldType = xsdTypeToProto(attr.Type, allPkgs)
	}

	injectComment := fmt.Sprintf("  // @gotags: xml:\"%s,attr\"", attr.Name)
	return fmt.Sprintf("%s\n  %s %s = %d;", injectComment, fieldType, fieldName, fieldNum)
}

func generateEnum(simpleType XSDSimpleType) string {
	var builder strings.Builder

	enumName := strings.ReplaceAll(toProtoMessageName(simpleType.Name), "_", "")
	builder.WriteString(fmt.Sprintf("enum %s {\n", enumName))

	enumPrefix := toProtoFieldName(simpleType.Name)
	enumPrefix = strings.ToUpper(enumPrefix)
	for strings.Contains(enumPrefix, "__") {
		enumPrefix = strings.ReplaceAll(enumPrefix, "__", "_")
	}

	builder.WriteString(fmt.Sprintf("  %s_UNSPECIFIED = 0;\n", enumPrefix))

	for i, enum := range simpleType.Restriction.Enumerations {
		rawValue := toProtoEnumValue(enum.Value)
		enumValue := enumPrefix + "_" + rawValue
		builder.WriteString(fmt.Sprintf("  %s = %d;\n", enumValue, i+1))
	}

	builder.WriteString("}")
	return builder.String()
}

//
// =======================
// Type/name helpers (kept + small improvements)
// =======================
//

func xsdTypeToProto(xsdType string, allPkgs map[string]protoPkgInfo) string {
	originalType := xsdType
	var prefix string
	
	// Extract prefix if present (xs:, avs:, ern:, etc.)
	if idx := strings.Index(xsdType, ":"); idx != -1 {
		prefix = xsdType[:idx]
		xsdType = xsdType[idx+1:]
	}

	switch xsdType {
	case "string", "normalizedString", "token", "anyURI", "NMTOKEN":
		return "string"
	case "int", "integer", "positiveInteger", "PositiveInteger":
		return "int32"
	case "long":
		return "int64"
	case "boolean":
		return "bool"
	case "decimal", "float":
		return "string" // preserve precision for decimals
	case "double":
		return "double"
	case "dateTime", "date", "time", "duration", "gYear", "GYear", "ddex_IsoDate", "Ddex_IsoDate":
		return "string" // ISO8601 strings
	case "base64Binary":
		return "bytes"
	default:
		// Handle namespace prefixes for custom types
		if prefix == "avs" {
			// AVS types go to ddex.avs package
			return "ddex.avs." + strings.ReplaceAll(toProtoMessageName(xsdType), "_", "")
		}
		
		// For other prefixes, try to map to known packages
		if prefix != "" && prefix != "xs" {
			// Look for a namespace that might match this prefix
			for ns, pkg := range allPkgs {
				if strings.Contains(strings.ToLower(ns), prefix) {
					packageName := pkg.pkgName
					return packageName + "." + strings.ReplaceAll(toProtoMessageName(xsdType), "_", "")
				}
			}
		}
		
		// Assume custom type → Proto message in local package
		if xsdType != originalType {
			log.Printf("Unmapped XSD type: %s (original: %s) -> treating as custom message", xsdType, originalType)
		} else {
			log.Printf("Unmapped XSD type: %s -> treating as custom message", xsdType)
		}
		return strings.ReplaceAll(toProtoMessageName(xsdType), "_", "")
	}
}

func toProtoFieldName(name string) string {
	var b strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteByte(byte(strings.ToLower(string(r))[0]))
	}
	return b.String()
}

func toProtoMessageName(name string) string {
	if name == "" {
		return ""
	}
	// Basic PascalCase
	return strings.ToUpper(name[:1]) + name[1:]
}

func toProtoEnumName(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func toProtoEnumValue(value string) string {
	result := strings.ToUpper(value)
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, " ", "_")
	result = strings.ReplaceAll(result, ".", "_")
	result = strings.ReplaceAll(result, "/", "_")
	result = strings.ReplaceAll(result, "+", "_PLUS_")
	result = strings.ReplaceAll(result, "(", "_")
	result = strings.ReplaceAll(result, ")", "_")
	result = strings.ReplaceAll(result, "'", "_")
	result = strings.ReplaceAll(result, "\"", "_")
	result = strings.ReplaceAll(result, "&", "_AND_")

	var cleaned strings.Builder
	for _, r := range result {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			cleaned.WriteRune(r)
		}
	}
	out := cleaned.String()
	if out == "" {
		out = "UNKNOWN"
	}
	if out[0] >= '0' && out[0] <= '9' {
		out = "E_" + out
	}
	for strings.Contains(out, "__") {
		out = strings.ReplaceAll(out, "__", "_")
	}
	return strings.Trim(out, "_")
}

func toPosixPath(p string) string {
	return strings.ReplaceAll(p, string(os.PathSeparator), "/")
}

//
// =======================
// Namespace → package mapping
// =======================
//

// For DDEX, we want:
//
//	ddex.xml/ern/43 → package "ddex.ern.v43"
//	ddex.xml/mead/11 → "ddex.mead.v11"
//	ddex.xml/pie/10 → "ddex.pie.v10"
//	ddex.xml/avs/avs → "ddex.avs"
func namespaceToProtoPackage(ns string, spec struct{ name, version, mainFile string }) string {
	host, pathParts := splitNS(ns)

	// DDEX-friendly mapping
	if host == "ddex.net" && len(pathParts) >= 2 && pathParts[0] == "xml" {
		// AVS in the wild appears as: /xml/avs/avs, /xml/allowed-value-sets, /xml/allowed_value_sets
		if pathParts[1] == "avs" ||
			pathParts[1] == "allowed-value-sets" ||
			pathParts[1] == "allowed_value_sets" {
			return "ddex.avs"
		}
		// Normal versioned families: /xml/{ern|mead|pie}/{digits}
		if len(pathParts) >= 3 && isDigits(pathParts[2]) {
			return fmt.Sprintf("ddex.%s.v%s", pathParts[1], pathParts[2])
		}
	}

	// Fallback: if this namespace matches the entry spec (ern/mead/pie…), pin it
	if looksLikeEntry(ns, spec) {
		return fmt.Sprintf("ddex.%s.v%s", spec.name, spec.version)
	}

	// Generic fallback: reverse host + path; add v<digits> suffix where appropriate
	revHost := reverseHost(host)
	if len(pathParts) > 0 && isDigits(pathParts[len(pathParts)-1]) {
		last := "v" + pathParts[len(pathParts)-1]
		return sanitizePackage(strings.Join(append([]string{revHost}, append(pathParts[:len(pathParts)-1], last)...), "."))
	}
	return sanitizePackage(strings.Join(append([]string{revHost}, pathParts...), "."))
}

func looksLikeEntry(ns string, spec struct{ name, version, mainFile string }) bool {
	// ddex.net/xml/<spec>/<versionDigits>, but never treat AVS as an entry package
	host, parts := splitNS(ns)
	if host != "ddex.net" || len(parts) < 3 || parts[0] != "xml" {
		return false
	}
	if parts[1] == "avs" || parts[1] == "allowed-value-sets" || parts[1] == "allowed_value_sets" {
		return false
	}
	return parts[1] == spec.name && isDigits(parts[2]) && parts[2] == stripLeadingV(spec.version)
}

func namespaceToGoPackage(ns string, spec struct{ name, version, mainFile string }) string {
	// Put Go package paths under your repo. Mirror the proto package path as directories.
	pkg := namespaceToProtoPackage(ns, spec)
	path := strings.ReplaceAll(pkg, ".", "/")
	return "github.com/alecsavvy/ddex-go/gen/" + path
}

func packageToPath(pkg string) string {
	parts := strings.Split(pkg, ".")
	if len(parts) == 0 {
		return "unknown.proto"
	}
	dir := strings.Join(parts, "/")
	filename := parts[len(parts)-1] + ".proto"
	return filepath.Join(dir, filename)
}

func splitNS(ns string) (host string, parts []string) {
	u, err := url.Parse(ns)
	if err != nil {
		return "", nil
	}
	host = u.Host
	parts = strings.Split(strings.Trim(u.Path, "/"), "/")
	return
}

func reverseHost(h string) string {
	if h == "" {
		return "unknown"
	}
	parts := strings.Split(h, ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, ".")
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func sanitizePackage(p string) string {
	// very light cleanup: no invalid idents
	p = strings.ReplaceAll(p, "-", "_")
	p = strings.ReplaceAll(p, " ", "_")
	return p
}

func stripLeadingV(s string) string {
	return strings.TrimPrefix(strings.ToLower(s), "v")
}
