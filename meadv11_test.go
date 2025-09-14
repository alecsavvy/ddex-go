package ddex

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	// Proto-generated implementations
	meadv11 "github.com/alecsavvy/ddex-go/gen/ddex/mead/v11"
)

// MEAD test files from DDEX samples
var meadTestFiles = map[string]string{
	"Award Example": "mead_award_example.xml",
}

// TestMEADConformance validates parsing against DDEX MEAD sample files
func TestMEADConformance(t *testing.T) {
	for testName, filename := range meadTestFiles {
		t.Run("MEAD_"+testName, func(t *testing.T) {
			xmlPath := filepath.Join("testdata", "meadv11", filename)

			// Read the sample XML file
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Skipf("Sample file not found: %s", xmlPath)
			}

			// Test proto-generated MEAD parsing
			var protoMEAD meadv11.MeadMessage
			err = xml.Unmarshal(xmlData, &protoMEAD)
			if err != nil {
				t.Fatalf("Failed to parse %s with proto structs: %v", filename, err)
			}

			// Validate structure
			validateMEADStructure(t, &protoMEAD, filename)

			t.Logf("✓ Successfully parsed %s (%d bytes)", filename, len(xmlData))
		})
	}
}

// TestMEADRoundTrip tests XML -> struct -> XML round-trip functionality
func TestMEADRoundTrip(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
	}{
		{"Award Example", "mead_award_example.xml"},
	}

	for _, tc := range testCases {
		t.Run("MEAD_RoundTrip_"+tc.name, func(t *testing.T) {
			xmlPath := filepath.Join("testdata", "meadv11", tc.filename)

			// Read original XML
			originalData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Skipf("Sample file not found: %s", xmlPath)
			}

			// Parse original
			var originalMsg meadv11.MeadMessage
			err = xml.Unmarshal(originalData, &originalMsg)
			if err != nil {
				t.Fatalf("Failed to unmarshal original: %v", err)
			}

			// Marshal back to XML
			regenerated, err := xml.MarshalIndent(&originalMsg, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal back to XML: %v", err)
			}

			// Add XML header for valid XML
			fullXML := []byte(xml.Header + string(regenerated))

			// Parse regenerated XML
			var roundTripMsg meadv11.MeadMessage
			err = xml.Unmarshal(fullXML, &roundTripMsg)
			if err != nil {
				t.Fatalf("Round trip parsing failed: %v", err)
			}

			// Semantic comparison
			if !semanticallyEqualMEAD(&originalMsg, &roundTripMsg) {
				t.Errorf("Round trip changed semantic content for %s", tc.filename)

				// Debug output
				t.Logf("Original MessageId: %s", getMEADMessageId(&originalMsg))
				t.Logf("RoundTrip MessageId: %s", getMEADMessageId(&roundTripMsg))
			} else {
				t.Logf("✓ Round trip successful for %s", tc.filename)
			}
		})
	}
}

// TestMEADFieldCompleteness validates that critical fields are present and populated
func TestMEADFieldCompleteness(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
	}{
		{"Award Example", "mead_award_example.xml"},
	}

	for _, tc := range testCases {
		t.Run("MEAD_Completeness_"+tc.name, func(t *testing.T) {
			xmlPath := filepath.Join("testdata", "meadv11", tc.filename)

			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Skipf("Sample file not found: %s", xmlPath)
			}

			var msg meadv11.MeadMessage
			err = xml.Unmarshal(xmlData, &msg)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			// Test required top-level fields
			checkRequired := []struct {
				name  string
				value interface{}
			}{
				{"MessageHeader", msg.MessageHeader},
				{"ReleaseInformationList", msg.ReleaseInformationList},
			}

			for _, check := range checkRequired {
				if check.value == nil {
					t.Errorf("Required field %s is nil", check.name)
				} else if reflect.ValueOf(check.value).IsZero() {
					t.Errorf("Required field %s is zero value", check.name)
				}
			}

			// Test MessageHeader completeness
			if msg.MessageHeader != nil {
				if msg.MessageHeader.MessageId == "" {
					t.Error("MessageHeader.MessageId is empty")
				}
				if msg.MessageHeader.MessageSender == nil {
					t.Error("MessageHeader.MessageSender is nil")
				}
			}

			// Test ReleaseInformationList completeness
			if msg.ReleaseInformationList != nil {
				releaseCount := len(msg.ReleaseInformationList.ReleaseInformation)
				if releaseCount == 0 {
					t.Error("ReleaseInformationList contains no releases")
				} else {
					t.Logf("✓ Found %d release(s) in %s", releaseCount, tc.filename)
				}
			}
		})
	}
}

// TestMEADXMLTagsEffectiveness tests that XML tags are working correctly
func TestMEADXMLTagsEffectiveness(t *testing.T) {
	xmlPath := filepath.Join("testdata", "meadv11", "mead_award_example.xml")

	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Skip("Sample file not found")
	}

	var msg meadv11.MeadMessage
	err = xml.Unmarshal(xmlData, &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Test that we can marshal it back
	_, err = xml.Marshal(&msg)
	if err != nil {
		t.Errorf("XML tags not working - marshal failed: %v", err)
	} else {
		t.Log("✓ MEAD XML tags working correctly - marshal/unmarshal successful")
	}
}

// Helper functions

func validateMEADStructure(t *testing.T, msg *meadv11.MeadMessage, filename string) {
	if msg.MessageHeader == nil {
		t.Errorf("MessageHeader is nil in %s", filename)
		return
	}

	if msg.MessageHeader.MessageId == "" {
		t.Errorf("MessageId is empty in %s", filename)
	}

	if msg.ReleaseInformationList == nil {
		t.Errorf("ReleaseInformationList is nil in %s", filename)
		return
	}

	releaseCount := len(msg.ReleaseInformationList.ReleaseInformation)
	if releaseCount == 0 {
		t.Errorf("No release information found in %s", filename)
	}
}

func semanticallyEqualMEAD(msg1, msg2 *meadv11.MeadMessage) bool {
	// Compare critical fields for semantic equality

	// Both nil or both non-nil
	if (msg1.MessageHeader == nil) != (msg2.MessageHeader == nil) {
		return false
	}

	if msg1.MessageHeader != nil && msg2.MessageHeader != nil {
		if msg1.MessageHeader.MessageId != msg2.MessageHeader.MessageId {
			return false
		}
	}

	// Compare release information counts
	if (msg1.ReleaseInformationList == nil) != (msg2.ReleaseInformationList == nil) {
		return false
	}

	if msg1.ReleaseInformationList != nil && msg2.ReleaseInformationList != nil {
		count1 := len(msg1.ReleaseInformationList.ReleaseInformation)
		count2 := len(msg2.ReleaseInformationList.ReleaseInformation)
		if count1 != count2 {
			return false
		}
	}

	return true
}

func getMEADMessageId(msg *meadv11.MeadMessage) string {
	if msg.MessageHeader != nil {
		return msg.MessageHeader.MessageId
	}
	return ""
}

// BenchmarkMEADParsing benchmarks MEAD XML parsing performance
func BenchmarkMEADParsing(b *testing.B) {
	xmlPath := filepath.Join("testdata", "meadv11", "mead_award_example.xml")
	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		b.Skip("Sample file not found")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var msg meadv11.MeadMessage
		err := xml.Unmarshal(xmlData, &msg)
		if err != nil {
			b.Fatalf("Unmarshal failed: %v", err)
		}
	}
}

func BenchmarkMEADMarshal(b *testing.B) {
	xmlPath := filepath.Join("testdata", "meadv11", "mead_award_example.xml")
	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		b.Skip("Sample file not found")
	}

	var msg meadv11.MeadMessage
	err = xml.Unmarshal(xmlData, &msg)
	if err != nil {
		b.Fatalf("Setup failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := xml.Marshal(&msg)
		if err != nil {
			b.Fatalf("Marshal failed: %v", err)
		}
	}
}