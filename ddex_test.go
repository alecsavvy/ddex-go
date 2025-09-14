package ddex

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	// Proto-generated implementations
	ernv432 "github.com/alecsavvy/ddex-go/gen/ddex/ern/v432"
	meadv11 "github.com/alecsavvy/ddex-go/gen/ddex/mead/v11"
	piev10 "github.com/alecsavvy/ddex-go/gen/ddex/pie/v10"
	"github.com/alecsavvy/ddex-go/testdata"
)

// Test data maps for each message type
var (
	ernTestFiles = map[string]string{
		"Audio Album":           "1 Audio.xml",
		"Video Album":           "2 Video.xml",
		"Mixed Media Bundle":    "3 MixedMedia.xml",
		"Simple Audio Single":   "4 SimpleAudioSingle.xml",
		"Simple Video Single":   "5 SimpleVideoSingle.xml",
		"Ringtone":              "6 Ringtone.xml",
		"Longform Musical Work": "7 LongformMusicalWorkVideo.xml",
		"DJ Mix":                "8 DjMix.xml",
		"Classical Variant":     "Variant Classical.xml",
	}

	meadTestFiles = map[string]string{
		"Award Example": "mead_award_example.xml",
	}

	pieTestFiles = map[string]string{
		"Award Example": "pie_award_example.xml",
	}
)

// TestDDEXMessages is the main test suite for all DDEX message types
func TestDDEXMessages(t *testing.T) {
	// Test ERN messages
	t.Run("ERN", func(t *testing.T) {
		t.Parallel()

		t.Run("Conformance", func(t *testing.T) {
			for testName, filename := range ernTestFiles {
				t.Run(testName, func(t *testing.T) {
					xmlPath := filepath.Join("testdata", "ernv432", "Samples43", filename)
					xmlData, err := os.ReadFile(xmlPath)
					if err != nil {
						t.Skipf("Sample file not found: %s", xmlPath)
					}

					var msg ernv432.NewReleaseMessage
					err = xml.Unmarshal(xmlData, &msg)
					if err != nil {
						t.Fatalf("Failed to parse %s: %v", filename, err)
					}

					validateERNStructure(t, &msg, filename)
					t.Logf("✓ Successfully parsed %s (%d bytes)", filename, len(xmlData))
				})
			}
		})

		t.Run("RoundTrip", func(t *testing.T) {
			testProtoToXMLToProtoRoundTrip(t, "ERN", testdata.SimpleERNTest)
		})

		t.Run("FieldCompleteness", func(t *testing.T) {
			testCases := []struct {
				name     string
				filename string
			}{
				{"Audio Album", "1 Audio.xml"},
				{"Simple Audio Single", "4 SimpleAudioSingle.xml"},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
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

					// Test required fields
					validateRequiredFields(t, []fieldCheck{
						{"MessageHeader", msg.MessageHeader},
						{"ReleaseList", msg.ReleaseList},
					})

					// ERN-specific validations
					if msg.MessageHeader != nil {
						if msg.MessageHeader.MessageId == "" && tc.filename != "8 DjMix.xml" {
							t.Error("MessageHeader.MessageId is empty")
						}
						if msg.MessageHeader.MessageSender == nil {
							t.Error("MessageHeader.MessageSender is nil")
						}
					}

					if msg.ReleaseList != nil {
						releaseCount := countERNReleases(msg.ReleaseList)
						if releaseCount == 0 {
							t.Error("ReleaseList contains no releases")
						} else {
							t.Logf("✓ Found %d release(s) in %s", releaseCount, tc.filename)
						}
					}
				})
			}
		})
	})

	// Test MEAD messages
	t.Run("MEAD", func(t *testing.T) {
		t.Parallel()

		t.Run("Conformance", func(t *testing.T) {
			for testName, filename := range meadTestFiles {
				t.Run(testName, func(t *testing.T) {
					xmlPath := filepath.Join("testdata", "meadv11", filename)
					xmlData, err := os.ReadFile(xmlPath)
					if err != nil {
						t.Skipf("Sample file not found: %s", xmlPath)
					}

					var msg meadv11.MeadMessage
					err = xml.Unmarshal(xmlData, &msg)
					if err != nil {
						t.Fatalf("Failed to parse %s: %v", filename, err)
					}

					validateMEADStructure(t, &msg, filename)
					t.Logf("✓ Successfully parsed %s (%d bytes)", filename, len(xmlData))
				})
			}
		})

		t.Run("RoundTrip", func(t *testing.T) {
			testProtoToXMLToProtoRoundTrip(t, "MEAD", testdata.SimpleMEADTest)
		})

		t.Run("FieldCompleteness", func(t *testing.T) {
			for testName, filename := range meadTestFiles {
				t.Run(testName, func(t *testing.T) {
					xmlPath := filepath.Join("testdata", "meadv11", filename)
					xmlData, err := os.ReadFile(xmlPath)
					if err != nil {
						t.Skipf("Sample file not found: %s", xmlPath)
					}

					var msg meadv11.MeadMessage
					err = xml.Unmarshal(xmlData, &msg)
					if err != nil {
						t.Fatalf("Failed to unmarshal: %v", err)
					}

					// Test required fields
					validateRequiredFields(t, []fieldCheck{
						{"MessageHeader", msg.MessageHeader},
						{"ReleaseInformationList", msg.ReleaseInformationList},
					})

					// MEAD-specific validations
					if msg.MessageHeader != nil {
						if msg.MessageHeader.MessageId == "" {
							t.Error("MessageHeader.MessageId is empty")
						}
						if msg.MessageHeader.MessageSender == nil {
							t.Error("MessageHeader.MessageSender is nil")
						}
					}

					if msg.ReleaseInformationList != nil {
						releaseCount := len(msg.ReleaseInformationList.ReleaseInformation)
						if releaseCount == 0 {
							t.Error("ReleaseInformationList contains no releases")
						} else {
							t.Logf("✓ Found %d release(s) in %s", releaseCount, filename)
						}
					}
				})
			}
		})
	})

	// Test PIE messages
	t.Run("PIE", func(t *testing.T) {
		t.Parallel()

		t.Run("Conformance", func(t *testing.T) {
			for testName, filename := range pieTestFiles {
				t.Run(testName, func(t *testing.T) {
					xmlPath := filepath.Join("testdata", "piev10", filename)
					xmlData, err := os.ReadFile(xmlPath)
					if err != nil {
						t.Skipf("Sample file not found: %s", xmlPath)
					}

					var msg piev10.PieMessage
					err = xml.Unmarshal(xmlData, &msg)
					if err != nil {
						t.Fatalf("Failed to parse %s: %v", filename, err)
					}

					validatePIEStructure(t, &msg, filename)
					t.Logf("✓ Successfully parsed %s (%d bytes)", filename, len(xmlData))
				})
			}
		})

		t.Run("RoundTrip", func(t *testing.T) {
			testProtoToXMLToProtoRoundTrip(t, "PIE", testdata.SimplePIETest)
		})

		t.Run("FieldCompleteness", func(t *testing.T) {
			for testName, filename := range pieTestFiles {
				t.Run(testName, func(t *testing.T) {
					xmlPath := filepath.Join("testdata", "piev10", filename)
					xmlData, err := os.ReadFile(xmlPath)
					if err != nil {
						t.Skipf("Sample file not found: %s", xmlPath)
					}

					var msg piev10.PieMessage
					err = xml.Unmarshal(xmlData, &msg)
					if err != nil {
						t.Fatalf("Failed to unmarshal: %v", err)
					}

					// Test required fields
					validateRequiredFields(t, []fieldCheck{
						{"MessageHeader", msg.MessageHeader},
						{"PartyList", msg.PartyList},
					})

					// PIE-specific validations
					if msg.MessageHeader != nil {
						if msg.MessageHeader.MessageId == "" {
							t.Error("MessageHeader.MessageId is empty")
						}
						if msg.MessageHeader.MessageSender == nil {
							t.Error("MessageHeader.MessageSender is nil")
						}
					}

					if msg.PartyList != nil {
						partyCount := len(msg.PartyList.Party)
						if partyCount == 0 {
							t.Error("PartyList contains no parties")
						} else {
							t.Logf("✓ Found %d party(ies) in %s", partyCount, filename)

							// Count awards
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
		})
	})
}


// TestXMLTagsEffectiveness validates XML marshaling/unmarshaling for all message types
func TestXMLTagsEffectiveness(t *testing.T) {
	t.Run("ERN", func(t *testing.T) {
		t.Parallel()
		testXMLTags(t, "testdata/ernv432/Samples43/5 SimpleVideoSingle.xml", &ernv432.NewReleaseMessage{}, "ERN")
	})

	t.Run("MEAD", func(t *testing.T) {
		t.Parallel()
		testXMLTags(t, "testdata/meadv11/mead_award_example.xml", &meadv11.MeadMessage{}, "MEAD")
	})

	t.Run("PIE", func(t *testing.T) {
		t.Parallel()
		testXMLTags(t, "testdata/piev10/pie_award_example.xml", &piev10.PieMessage{}, "PIE")
	})
}

// Benchmark tests
func BenchmarkDDEX(b *testing.B) {
	b.Run("ERN", func(b *testing.B) {
		b.Run("Parse", func(b *testing.B) {
			xmlPath := filepath.Join("testdata", "ernv432", "Samples43", "1 Audio.xml")
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				b.Skip("Sample file not found")
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var msg ernv432.NewReleaseMessage
				xml.Unmarshal(xmlData, &msg)
			}
		})

		b.Run("Marshal", func(b *testing.B) {
			xmlPath := filepath.Join("testdata", "ernv432", "Samples43", "1 Audio.xml")
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				b.Skip("Sample file not found")
			}

			var msg ernv432.NewReleaseMessage
			xml.Unmarshal(xmlData, &msg)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				xml.Marshal(&msg)
			}
		})
	})

	b.Run("MEAD", func(b *testing.B) {
		b.Run("Parse", func(b *testing.B) {
			xmlPath := filepath.Join("testdata", "meadv11", "mead_award_example.xml")
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				b.Skip("Sample file not found")
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var msg meadv11.MeadMessage
				xml.Unmarshal(xmlData, &msg)
			}
		})

		b.Run("Marshal", func(b *testing.B) {
			xmlPath := filepath.Join("testdata", "meadv11", "mead_award_example.xml")
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				b.Skip("Sample file not found")
			}

			var msg meadv11.MeadMessage
			xml.Unmarshal(xmlData, &msg)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				xml.Marshal(&msg)
			}
		})
	})

	b.Run("PIE", func(b *testing.B) {
		b.Run("Parse", func(b *testing.B) {
			xmlPath := filepath.Join("testdata", "piev10", "pie_award_example.xml")
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				b.Skip("Sample file not found")
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var msg piev10.PieMessage
				xml.Unmarshal(xmlData, &msg)
			}
		})

		b.Run("Marshal", func(b *testing.B) {
			xmlPath := filepath.Join("testdata", "piev10", "pie_award_example.xml")
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				b.Skip("Sample file not found")
			}

			var msg piev10.PieMessage
			xml.Unmarshal(xmlData, &msg)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				xml.Marshal(&msg)
			}
		})
	})
}

// Helper functions

type fieldCheck struct {
	name  string
	value interface{}
}

func validateRequiredFields(t *testing.T, fields []fieldCheck) {
	for _, field := range fields {
		if field.value == nil {
			t.Errorf("Required field %s is nil", field.name)
		} else if reflect.ValueOf(field.value).IsZero() {
			t.Errorf("Required field %s is zero value", field.name)
		}
	}
}

func testProtoToXMLToProtoRoundTrip(t *testing.T, msgType string, constructor interface{}) {
	switch msgType {
	case "ERN":
		constructor := constructor.(func() *ernv432.NewReleaseMessage)
		original := constructor()

		// Marshal to XML
		xmlData, err := xml.MarshalIndent(original, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal to XML: %v", err)
		}

		// Add XML header for proper parsing
		fullXML := []byte(xml.Header + string(xmlData))

		// Unmarshal back to proto struct
		var roundTrip ernv432.NewReleaseMessage
		err = xml.Unmarshal(fullXML, &roundTrip)
		if err != nil {
			t.Fatalf("Failed to unmarshal from XML: %v", err)
		}

		// Compare structs using reflect for deep equality
		if !reflect.DeepEqual(original, &roundTrip) {
			t.Errorf("Round trip failed: original and unmarshaled structs are not equal")
			t.Logf("Original: %+v", original)
			t.Logf("RoundTrip: %+v", &roundTrip)
		} else {
			t.Log("✓ ERN proto->XML->proto round trip successful")
		}

	case "MEAD":
		constructor := constructor.(func() *meadv11.MeadMessage)
		original := constructor()

		// Marshal to XML
		xmlData, err := xml.MarshalIndent(original, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal to XML: %v", err)
		}

		// Add XML header for proper parsing
		fullXML := []byte(xml.Header + string(xmlData))

		// Unmarshal back to proto struct
		var roundTrip meadv11.MeadMessage
		err = xml.Unmarshal(fullXML, &roundTrip)
		if err != nil {
			t.Fatalf("Failed to unmarshal from XML: %v", err)
		}

		// Compare structs using reflect for deep equality
		if !reflect.DeepEqual(original, &roundTrip) {
			t.Errorf("Round trip failed: original and unmarshaled structs are not equal")
			t.Logf("Original: %+v", original)
			t.Logf("RoundTrip: %+v", &roundTrip)
		} else {
			t.Log("✓ MEAD proto->XML->proto round trip successful")
		}

	case "PIE":
		constructor := constructor.(func() *piev10.PieMessage)
		original := constructor()

		// Marshal to XML
		xmlData, err := xml.MarshalIndent(original, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal to XML: %v", err)
		}

		// Add XML header for proper parsing
		fullXML := []byte(xml.Header + string(xmlData))

		// Unmarshal back to proto struct
		var roundTrip piev10.PieMessage
		err = xml.Unmarshal(fullXML, &roundTrip)
		if err != nil {
			t.Fatalf("Failed to unmarshal from XML: %v", err)
		}

		// Compare structs using reflect for deep equality
		if !reflect.DeepEqual(original, &roundTrip) {
			t.Errorf("Round trip failed: original and unmarshaled structs are not equal")
			t.Logf("Original: %+v", original)
			t.Logf("RoundTrip: %+v", &roundTrip)
		} else {
			t.Log("✓ PIE proto->XML->proto round trip successful")
		}

	default:
		t.Fatalf("Unknown message type: %s", msgType)
	}
}

func testRoundTrip(t *testing.T, xmlPath string, msgType interface{}, filename string) {
	originalData, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Skipf("Sample file not found: %s", xmlPath)
	}

	// Parse original
	err = xml.Unmarshal(originalData, msgType)
	if err != nil {
		t.Fatalf("Failed to unmarshal original: %v", err)
	}

	// Marshal back
	regenerated, err := xml.MarshalIndent(msgType, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Parse regenerated
	msgCopy := reflect.New(reflect.TypeOf(msgType).Elem()).Interface()
	fullXML := []byte(xml.Header + string(regenerated))
	err = xml.Unmarshal(fullXML, msgCopy)
	if err != nil {
		t.Fatalf("Round trip parsing failed: %v", err)
	}

	// Type-specific semantic comparison
	equal := false
	switch msg := msgType.(type) {
	case *ernv432.NewReleaseMessage:
		equal = semanticallyEqualERN(msg, msgCopy.(*ernv432.NewReleaseMessage))
	case *meadv11.MeadMessage:
		equal = semanticallyEqualMEAD(msg, msgCopy.(*meadv11.MeadMessage))
	case *piev10.PieMessage:
		equal = semanticallyEqualPIE(msg, msgCopy.(*piev10.PieMessage))
	}

	if !equal {
		t.Errorf("Round trip changed semantic content for %s", filename)
	} else {
		t.Logf("✓ Round trip successful for %s", filename)
	}
}

func testXMLTags(t *testing.T, xmlPath string, msgType interface{}, msgName string) {
	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Skip("Sample file not found")
	}

	err = xml.Unmarshal(xmlData, msgType)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	_, err = xml.Marshal(msgType)
	if err != nil {
		t.Errorf("%s XML tags not working - marshal failed: %v", msgName, err)
	} else {
		t.Logf("✓ %s XML tags working correctly", msgName)
	}
}

// Validation functions for each message type

func validateERNStructure(t *testing.T, msg *ernv432.NewReleaseMessage, filename string) {
	if msg.MessageHeader == nil {
		t.Errorf("MessageHeader is nil in %s", filename)
		return
	}

	if msg.MessageHeader.MessageId == "" {
		// DJ Mix sample has intentionally empty MessageId
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

	releaseCount := countERNReleases(msg.ReleaseList)
	if releaseCount == 0 {
		t.Errorf("No releases found in %s", filename)
	}
}

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

// Semantic equality functions

func semanticallyEqualERN(msg1, msg2 *ernv432.NewReleaseMessage) bool {
	if (msg1.MessageHeader == nil) != (msg2.MessageHeader == nil) {
		return false
	}

	if msg1.MessageHeader != nil && msg2.MessageHeader != nil {
		if msg1.MessageHeader.MessageId != msg2.MessageHeader.MessageId {
			return false
		}
	}

	if (msg1.ReleaseList == nil) != (msg2.ReleaseList == nil) {
		return false
	}

	if msg1.ReleaseList != nil && msg2.ReleaseList != nil {
		count1 := countERNReleases(msg1.ReleaseList)
		count2 := countERNReleases(msg2.ReleaseList)
		if count1 != count2 {
			return false
		}
	}

	return true
}

func semanticallyEqualMEAD(msg1, msg2 *meadv11.MeadMessage) bool {
	if (msg1.MessageHeader == nil) != (msg2.MessageHeader == nil) {
		return false
	}

	if msg1.MessageHeader != nil && msg2.MessageHeader != nil {
		if msg1.MessageHeader.MessageId != msg2.MessageHeader.MessageId {
			return false
		}
	}

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

func semanticallyEqualPIE(msg1, msg2 *piev10.PieMessage) bool {
	if (msg1.MessageHeader == nil) != (msg2.MessageHeader == nil) {
		return false
	}

	if msg1.MessageHeader != nil && msg2.MessageHeader != nil {
		if msg1.MessageHeader.MessageId != msg2.MessageHeader.MessageId {
			return false
		}
	}

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

// Utility functions

func countERNReleases(releaseList *ernv432.ReleaseList) int {
	count := 0
	if releaseList.Release != nil {
		count++
	}
	count += len(releaseList.TrackRelease)
	count += len(releaseList.ClipRelease)
	return count
}

// Protobuf to XML test functions

func testProtobufToXMLERN(t *testing.T, filename string, constructor func() *ernv432.NewReleaseMessage) {
	// Read the original XML file
	xmlPath := filepath.Join("testdata", filename)
	originalData, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Skipf("Sample file not found: %s", xmlPath)
	}

	// Parse original XML into protobuf struct
	var originalMsg ernv432.NewReleaseMessage
	err = xml.Unmarshal(originalData, &originalMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal original XML: %v", err)
	}

	// Create a new protobuf message manually constructed to match the original
	constructedMsg := constructor()

	// Marshal both to XML
	originalXML, err := xml.MarshalIndent(&originalMsg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal original message: %v", err)
	}

	constructedXML, err := xml.MarshalIndent(constructedMsg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal constructed message: %v", err)
	}

	// Compare semantic equality
	if !semanticallyEqualERN(&originalMsg, constructedMsg) {
		t.Errorf("Constructed protobuf message does not match original for %s", filename)
		t.Logf("Original XML length: %d", len(originalXML))
		t.Logf("Constructed XML length: %d", len(constructedXML))
	} else {
		t.Logf("✓ Protobuf construction matches original for %s", filename)
	}
}

func testProtobufToXMLMEAD(t *testing.T, filename string, constructor func() *meadv11.MeadMessage) {
	// Read the original XML file
	xmlPath := filepath.Join("testdata", filename)
	originalData, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Skipf("Sample file not found: %s", xmlPath)
	}

	// Parse original XML into protobuf struct
	var originalMsg meadv11.MeadMessage
	err = xml.Unmarshal(originalData, &originalMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal original XML: %v", err)
	}

	// Create a new protobuf message manually constructed to match the original
	constructedMsg := constructor()

	// Marshal both to XML
	originalXML, err := xml.MarshalIndent(&originalMsg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal original message: %v", err)
	}

	constructedXML, err := xml.MarshalIndent(constructedMsg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal constructed message: %v", err)
	}

	// Compare semantic equality
	if !semanticallyEqualMEAD(&originalMsg, constructedMsg) {
		t.Errorf("Constructed protobuf message does not match original for %s", filename)
		t.Logf("Original XML length: %d", len(originalXML))
		t.Logf("Constructed XML length: %d", len(constructedXML))
	} else {
		t.Logf("✓ Protobuf construction matches original for %s", filename)
	}
}

func testProtobufToXMLPIE(t *testing.T, filename string, constructor func() *piev10.PieMessage) {
	// Read the original XML file
	xmlPath := filepath.Join("testdata", filename)
	originalData, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Skipf("Sample file not found: %s", xmlPath)
	}

	// Parse original XML into protobuf struct
	var originalMsg piev10.PieMessage
	err = xml.Unmarshal(originalData, &originalMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal original XML: %v", err)
	}

	// Create a new protobuf message manually constructed to match the original
	constructedMsg := constructor()

	// Marshal both to XML
	originalXML, err := xml.MarshalIndent(&originalMsg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal original message: %v", err)
	}

	constructedXML, err := xml.MarshalIndent(constructedMsg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal constructed message: %v", err)
	}

	// Compare semantic equality
	if !semanticallyEqualPIE(&originalMsg, constructedMsg) {
		t.Errorf("Constructed protobuf message does not match original for %s", filename)
		t.Logf("Original XML length: %d", len(originalXML))
		t.Logf("Constructed XML length: %d", len(constructedXML))
	} else {
		t.Logf("✓ Protobuf construction matches original for %s", filename)
	}
}

