package ddex

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	// Proto-generated implementations
	piev10 "github.com/alecsavvy/ddex-go/gen/ddex/pie/v10"
)

// PIE test files from DDEX samples
var pieTestFiles = map[string]string{
	"Award Example": "pie_award_example.xml",
}

// TestPIEConformance validates parsing against DDEX PIE sample files
func TestPIEConformance(t *testing.T) {
	for testName, filename := range pieTestFiles {
		t.Run("PIE_"+testName, func(t *testing.T) {
			xmlPath := filepath.Join("testdata", "piev10", filename)

			// Read the sample XML file
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Skipf("Sample file not found: %s", xmlPath)
			}

			// Test proto-generated PIE parsing
			var protoPIE piev10.PieMessage
			err = xml.Unmarshal(xmlData, &protoPIE)
			if err != nil {
				t.Fatalf("Failed to parse %s with proto structs: %v", filename, err)
			}

			// Validate structure
			validatePIEStructure(t, &protoPIE, filename)

			t.Logf("✓ Successfully parsed %s (%d bytes)", filename, len(xmlData))
		})
	}
}

// TestPIERoundTrip tests XML -> struct -> XML round-trip functionality
func TestPIERoundTrip(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
	}{
		{"Award Example", "pie_award_example.xml"},
	}

	for _, tc := range testCases {
		t.Run("PIE_RoundTrip_"+tc.name, func(t *testing.T) {
			xmlPath := filepath.Join("testdata", "piev10", tc.filename)

			// Read original XML
			originalData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Skipf("Sample file not found: %s", xmlPath)
			}

			// Parse original
			var originalMsg piev10.PieMessage
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
			var roundTripMsg piev10.PieMessage
			err = xml.Unmarshal(fullXML, &roundTripMsg)
			if err != nil {
				t.Fatalf("Round trip parsing failed: %v", err)
			}

			// Semantic comparison
			if !semanticallyEqualPIE(&originalMsg, &roundTripMsg) {
				t.Errorf("Round trip changed semantic content for %s", tc.filename)

				// Debug output
				t.Logf("Original MessageId: %s", getPIEMessageId(&originalMsg))
				t.Logf("RoundTrip MessageId: %s", getPIEMessageId(&roundTripMsg))
			} else {
				t.Logf("✓ Round trip successful for %s", tc.filename)
			}
		})
	}
}

// TestPIEFieldCompleteness validates that critical fields are present and populated
func TestPIEFieldCompleteness(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
	}{
		{"Award Example", "pie_award_example.xml"},
	}

	for _, tc := range testCases {
		t.Run("PIE_Completeness_"+tc.name, func(t *testing.T) {
			xmlPath := filepath.Join("testdata", "piev10", tc.filename)

			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Skipf("Sample file not found: %s", xmlPath)
			}

			var msg piev10.PieMessage
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
				{"PartyList", msg.PartyList},
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

			// Test PartyList completeness
			if msg.PartyList != nil {
				partyCount := len(msg.PartyList.Party)
				if partyCount == 0 {
					t.Error("PartyList contains no parties")
				} else {
					t.Logf("✓ Found %d party(ies) in %s", partyCount, tc.filename)

					// Check for awards in parties
					totalAwards := 0
					for _, party := range msg.PartyList.Party {
						totalAwards += len(party.Award)
					}
					if totalAwards > 0 {
						t.Logf("✓ Found %d total award(s) across all parties", totalAwards)
					}
				}
			}
		})
	}
}

// TestPIEXMLTagsEffectiveness tests that XML tags are working correctly
func TestPIEXMLTagsEffectiveness(t *testing.T) {
	xmlPath := filepath.Join("testdata", "piev10", "pie_award_example.xml")

	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Skip("Sample file not found")
	}

	var msg piev10.PieMessage
	err = xml.Unmarshal(xmlData, &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Test that we can marshal it back
	_, err = xml.Marshal(&msg)
	if err != nil {
		t.Errorf("XML tags not working - marshal failed: %v", err)
	} else {
		t.Log("✓ PIE XML tags working correctly - marshal/unmarshal successful")
	}
}

// Helper functions

func validatePIEStructure(t *testing.T, msg *piev10.PieMessage, filename string) {
	if msg.MessageHeader == nil {
		t.Errorf("MessageHeader is nil in %s", filename)
		return
	}

	if msg.MessageHeader.MessageId == "" {
		t.Errorf("MessageId is empty in %s", filename)
	}

	if msg.PartyList == nil {
		t.Errorf("PartyList is nil in %s", filename)
		return
	}

	partyCount := len(msg.PartyList.Party)
	if partyCount == 0 {
		t.Errorf("No parties found in %s", filename)
	}
}

func semanticallyEqualPIE(msg1, msg2 *piev10.PieMessage) bool {
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

	// Compare party counts
	if (msg1.PartyList == nil) != (msg2.PartyList == nil) {
		return false
	}

	if msg1.PartyList != nil && msg2.PartyList != nil {
		count1 := len(msg1.PartyList.Party)
		count2 := len(msg2.PartyList.Party)
		if count1 != count2 {
			return false
		}
	}

	return true
}

func getPIEMessageId(msg *piev10.PieMessage) string {
	if msg.MessageHeader != nil {
		return msg.MessageHeader.MessageId
	}
	return ""
}

// BenchmarkPIEParsing benchmarks PIE XML parsing performance
func BenchmarkPIEParsing(b *testing.B) {
	xmlPath := filepath.Join("testdata", "piev10", "pie_award_example.xml")
	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		b.Skip("Sample file not found")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var msg piev10.PieMessage
		err := xml.Unmarshal(xmlData, &msg)
		if err != nil {
			b.Fatalf("Unmarshal failed: %v", err)
		}
	}
}

func BenchmarkPIEMarshal(b *testing.B) {
	xmlPath := filepath.Join("testdata", "piev10", "pie_award_example.xml")
	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		b.Skip("Sample file not found")
	}

	var msg piev10.PieMessage
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