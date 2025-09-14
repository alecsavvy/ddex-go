package ddex

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	// Proto-generated implementations
	ernv432 "github.com/alecsavvy/ddex-go/gen/ddex/ern/v432"
)

// ERN test files from official DDEX samples
var ernTestFiles = map[string]string{
	"Audio Album":             "1 Audio.xml",
	"Video Album":             "2 Video.xml",
	"Mixed Media Bundle":      "3 MixedMedia.xml",
	"Simple Audio Single":     "4 SimpleAudioSingle.xml",
	"Simple Video Single":     "5 SimpleVideoSingle.xml",
	"Ringtone":               "6 Ringtone.xml",
	"Longform Musical Work":   "7 LongformMusicalWorkVideo.xml",
	"DJ Mix":                 "8 DjMix.xml",
	"Classical Variant":       "Variant Classical.xml",
}

// TestConformance validates parsing against official DDEX sample files
func TestConformance(t *testing.T) {
	for testName, filename := range ernTestFiles {
		t.Run("ERN_"+testName, func(t *testing.T) {
			xmlPath := filepath.Join("testdata", "ernv432", "Samples43", filename)

			// Read the sample XML file
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Skipf("Sample file not found: %s", xmlPath)
			}

			// Test proto-generated ERN parsing
			var protoERN ernv432.NewReleaseMessage
			err = xml.Unmarshal(xmlData, &protoERN)
			if err != nil {
				t.Fatalf("Failed to parse %s with proto structs: %v", filename, err)
			}

			// Validate structure
			validateERNStructure(t, &protoERN, filename)

			t.Logf("✓ Successfully parsed %s (%d bytes)", filename, len(xmlData))
		})
	}
}

// TestRoundTrip tests XML -> struct -> XML round-trip functionality
func TestRoundTrip(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
	}{
		{"Audio Album", "1 Audio.xml"},
		{"Simple Video Single", "5 SimpleVideoSingle.xml"},
		{"DJ Mix", "8 DjMix.xml"},
	}

	for _, tc := range testCases {
		t.Run("RoundTrip_"+tc.name, func(t *testing.T) {
			xmlPath := filepath.Join("testdata", "ernv432", "Samples43", tc.filename)

			// Read original XML
			originalData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Skipf("Sample file not found: %s", xmlPath)
			}

			// Parse original
			var originalMsg ernv432.NewReleaseMessage
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
			var roundTripMsg ernv432.NewReleaseMessage
			err = xml.Unmarshal(fullXML, &roundTripMsg)
			if err != nil {
				t.Fatalf("Round trip parsing failed: %v", err)
			}

			// Semantic comparison
			if !semanticallyEqualERN(&originalMsg, &roundTripMsg) {
				t.Errorf("Round trip changed semantic content for %s", tc.filename)

				// Debug output
				t.Logf("Original MessageId: %s", getMessageId(&originalMsg))
				t.Logf("RoundTrip MessageId: %s", getMessageId(&roundTripMsg))
			} else {
				t.Logf("✓ Round trip successful for %s", tc.filename)
			}
		})
	}
}

// TestFieldCompleteness validates that critical fields are present and populated
func TestFieldCompleteness(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
	}{
		{"Audio Album", "1 Audio.xml"},
		{"Simple Audio Single", "4 SimpleAudioSingle.xml"},
	}

	for _, tc := range testCases {
		t.Run("Completeness_"+tc.name, func(t *testing.T) {
			xmlPath := filepath.Join("testdata", "ernv432", "Samples43", tc.filename)

			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Skipf("Sample file not found: %s", xmlPath)
			}

			var msg ernv432.NewReleaseMessage
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
				{"ReleaseList", msg.ReleaseList},
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

			// Test ReleaseList completeness
			if msg.ReleaseList != nil {
				releaseCount := countReleases(msg.ReleaseList)
				if releaseCount == 0 {
					t.Error("ReleaseList contains no releases")
				} else {
					t.Logf("✓ Found %d release(s) in %s", releaseCount, tc.filename)
				}
			}
		})
	}
}

// TestXMLTagsEffectiveness tests that XML tags are working correctly
func TestXMLTagsEffectiveness(t *testing.T) {
	xmlPath := filepath.Join("testdata", "ernv432", "Samples43", "5 SimpleVideoSingle.xml")

	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Skip("Sample file not found")
	}

	var msg ernv432.NewReleaseMessage
	err = xml.Unmarshal(xmlData, &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Test that we can marshal it back
	_, err = xml.Marshal(&msg)
	if err != nil {
		t.Errorf("XML tags not working - marshal failed: %v", err)
	} else {
		t.Log("✓ XML tags working correctly - marshal/unmarshal successful")
	}
}

// Helper functions

func validateERNStructure(t *testing.T, msg *ernv432.NewReleaseMessage, filename string) {
	if msg.MessageHeader == nil {
		t.Errorf("MessageHeader is nil in %s", filename)
		return
	}

	if msg.MessageHeader.MessageId == "" {
		// DJ Mix sample has intentionally empty MessageId - this is valid DDEX
		if filename != "8 DjMix.xml" {
			t.Errorf("MessageId is empty in %s", filename)
		} else {
			t.Logf("Note: MessageId is intentionally empty in %s (valid DDEX format)", filename)
		}
	}

	if msg.ReleaseList == nil {
		t.Errorf("ReleaseList is nil in %s", filename)
		return
	}

	releaseCount := countReleases(msg.ReleaseList)
	if releaseCount == 0 {
		t.Errorf("No releases found in %s", filename)
	}
}

func countReleases(releaseList *ernv432.ReleaseList) int {
	count := 0
	if releaseList.Release != nil {
		count++
	}
	count += len(releaseList.TrackRelease)
	count += len(releaseList.ClipRelease)
	return count
}

func semanticallyEqualERN(msg1, msg2 *ernv432.NewReleaseMessage) bool {
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

	// Compare release counts
	if (msg1.ReleaseList == nil) != (msg2.ReleaseList == nil) {
		return false
	}

	if msg1.ReleaseList != nil && msg2.ReleaseList != nil {
		count1 := countReleases(msg1.ReleaseList)
		count2 := countReleases(msg2.ReleaseList)
		if count1 != count2 {
			return false
		}
	}

	return true
}

func getMessageId(msg *ernv432.NewReleaseMessage) string {
	if msg.MessageHeader != nil {
		return msg.MessageHeader.MessageId
	}
	return ""
}

// BenchmarkParsing benchmarks XML parsing performance
func BenchmarkERNParsing(b *testing.B) {
	xmlPath := filepath.Join("testdata", "ernv432", "Samples43", "1 Audio.xml")
	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		b.Skip("Sample file not found")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var msg ernv432.NewReleaseMessage
		err := xml.Unmarshal(xmlData, &msg)
		if err != nil {
			b.Fatalf("Unmarshal failed: %v", err)
		}
	}
}

func BenchmarkERNMarshal(b *testing.B) {
	xmlPath := filepath.Join("testdata", "ernv432", "Samples43", "1 Audio.xml")
	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		b.Skip("Sample file not found")
	}

	var msg ernv432.NewReleaseMessage
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