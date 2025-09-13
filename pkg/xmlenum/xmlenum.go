// Package xmlenum provides a generic wrapper for protobuf enums that handles XML marshaling with string values
package xmlenum

import (
	"encoding/xml"
	"strings"
)

// XMLEnum wraps protobuf enum types to handle XML marshaling/unmarshaling with string values
// while preserving the original string for perfect round-trip fidelity
type XMLEnum[T ~int32] struct {
	Value    T      // The protobuf enum value
	RawValue string // Original string from XML (for round-trip fidelity)
}

// UnmarshalXML implements xml.Unmarshaler interface
func (e *XMLEnum[T]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	
	e.RawValue = s
	
	// Try to parse using case-insensitive matching
	if val, ok := parseEnumString[T](strings.ToUpper(s)); ok {
		e.Value = val
	} else {
		e.Value = T(0) // UNSPECIFIED value
	}
	return nil
}

// MarshalXML implements xml.Marshaler interface
func (e XMLEnum[T]) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	// Always use original string for perfect round-trip
	return enc.EncodeElement(e.RawValue, start)
}

// parseEnumString is a placeholder - actual implementations will be generated per-package
func parseEnumString[T ~int32](s string) (T, bool) {
	return T(0), false
}