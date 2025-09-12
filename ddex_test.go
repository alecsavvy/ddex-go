package ddex

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecsavvy/ddex-go/ddex/ernv432"
	"github.com/alecsavvy/ddex-go/ddex/meadv11"
	"github.com/alecsavvy/ddex-go/ddex/piev10"
)

func TestERNUnmarshaling(t *testing.T) {
	testFiles := []struct {
		filename string
		description string
	}{
		{"1 Audio.xml", "Audio Album"},
		{"2 Video.xml", "Video Album"},
		{"4 SimpleAudioSingle.xml", "Simple Audio Single"},
		{"5 SimpleVideoSingle.xml", "Simple Video Single"},
		{"6 Ringtone.xml", "Ringtone"},
		{"8 DjMix.xml", "DJ Mix"},
	}

	for _, testFile := range testFiles {
		t.Run(testFile.description, func(t *testing.T) {
			// Read the sample XML file
			xmlPath := filepath.Join("testdata", "ernv432", "Samples43", testFile.filename)
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", xmlPath, err)
			}

			// Test NewReleaseMessage unmarshaling
			var newRelease ernv432.NewReleaseMessage
			err = xml.Unmarshal(xmlData, &newRelease)
			if err != nil {
				t.Fatalf("Failed to unmarshal %s as NewReleaseMessage: %v", testFile.filename, err)
			}

			// Basic validation
			if newRelease.MessageHeader == nil {
				t.Errorf("MessageHeader is nil in %s", testFile.filename)
			}
			
			if newRelease.ReleaseList == nil {
				t.Errorf("ReleaseList is nil in %s", testFile.filename)
				return
			}

			// Count total releases
			releaseCount := 0
			if newRelease.ReleaseList.Release != nil {
				releaseCount++
			}
			releaseCount += len(newRelease.ReleaseList.TrackRelease)
			releaseCount += len(newRelease.ReleaseList.ClipRelease)
			
			if releaseCount == 0 {
				t.Errorf("No releases found in %s", testFile.filename)
			}

			t.Logf("Successfully parsed %s: %d release(s)", testFile.filename, releaseCount)
		})
	}
}

func TestPurgeReleaseMessage(t *testing.T) {
	// Note: Purge messages typically use the same structure but different root element
	// This test demonstrates how to handle different message types
	xmlPath := filepath.Join("testdata", "ernv432", "Samples43", "1 Audio.xml")
	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Skipf("Sample file not found: %v", err)
	}

	// Try parsing as different message types to test flexibility
	var newRelease ernv432.NewReleaseMessage
	err = xml.Unmarshal(xmlData, &newRelease)
	if err != nil {
		t.Fatalf("Failed to unmarshal as NewReleaseMessage: %v", err)
	}

	// Verify we can access type aliases
	var aliased NewReleaseMessageV432 = newRelease
	if aliased.MessageHeader == nil {
		t.Error("Type alias not working correctly")
	}
}

func TestXMLRoundTrip(t *testing.T) {
	xmlPath := filepath.Join("testdata", "ernv432", "Samples43", "5 SimpleVideoSingle.xml")
	
	// Read original
	originalData, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Skipf("Sample file not found: %v", err)
	}

	// Unmarshal
	var message ernv432.NewReleaseMessage
	err = xml.Unmarshal(originalData, &message)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Marshal back
	marshaledData, err := xml.MarshalIndent(message, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Basic validation - should be valid XML
	var roundTrip ernv432.NewReleaseMessage
	err = xml.Unmarshal(marshaledData, &roundTrip)
	if err != nil {
		t.Fatalf("Round trip failed: %v", err)
	}

	t.Logf("Round trip successful for %s", xmlPath)
}

func TestPIEUnmarshaling(t *testing.T) {
	testFiles := []struct {
		filename    string
		description string
	}{
		{"pie_award_example.xml", "PIE Award Example"},
	}

	for _, testFile := range testFiles {
		t.Run(testFile.description, func(t *testing.T) {
			// Read the sample XML file
			xmlPath := filepath.Join("testdata", "piev10", testFile.filename)
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", xmlPath, err)
			}

			// Test PieMessage unmarshaling
			var pieMessage piev10.PieMessage
			err = xml.Unmarshal(xmlData, &pieMessage)
			if err != nil {
				t.Fatalf("Failed to unmarshal %s as PieMessage: %v", testFile.filename, err)
			}

			// Basic validation
			if pieMessage.MessageHeader == nil {
				t.Errorf("MessageHeader is nil in %s", testFile.filename)
			}

			if pieMessage.PartyList == nil {
				t.Errorf("PartyList is nil in %s", testFile.filename)
				return
			}

			// Count parties
			partyCount := len(pieMessage.PartyList.Party)
			if partyCount == 0 {
				t.Errorf("No parties found in %s", testFile.filename)
			}

			// Check for awards
			hasAwards := false
			for _, party := range pieMessage.PartyList.Party {
				if len(party.Award) > 0 {
					hasAwards = true
					break
				}
			}

			if !hasAwards {
				t.Errorf("No awards found in %s", testFile.filename)
			}

			t.Logf("Successfully parsed %s: %d partie(s)", testFile.filename, partyCount)
		})
	}
}

func TestMEADUnmarshaling(t *testing.T) {
	testFiles := []struct {
		filename    string
		description string
	}{
		{"mead_award_example.xml", "MEAD Award Example"},
	}

	for _, testFile := range testFiles {
		t.Run(testFile.description, func(t *testing.T) {
			// Read the sample XML file
			xmlPath := filepath.Join("testdata", "meadv11", testFile.filename)
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", xmlPath, err)
			}

			// Test MeadMessage unmarshaling
			var meadMessage meadv11.MeadMessage
			err = xml.Unmarshal(xmlData, &meadMessage)
			if err != nil {
				t.Fatalf("Failed to unmarshal %s as MeadMessage: %v", testFile.filename, err)
			}

			// Basic validation
			if meadMessage.MessageHeader == nil {
				t.Errorf("MessageHeader is nil in %s", testFile.filename)
			}

			if meadMessage.ReleaseInformationList == nil {
				t.Errorf("ReleaseInformationList is nil in %s", testFile.filename)
				return
			}

			// Count releases
			releaseCount := len(meadMessage.ReleaseInformationList.ReleaseInformation)
			if releaseCount == 0 {
				t.Errorf("No release information found in %s", testFile.filename)
			}

			// Check for basic release data
			hasReleaseData := false
			for _, releaseInfo := range meadMessage.ReleaseInformationList.ReleaseInformation {
				if releaseInfo.ReleaseSummary != nil && releaseInfo.ReleaseSummary.ReleaseId != nil {
					hasReleaseData = true
					break
				}
			}

			if !hasReleaseData {
				t.Errorf("No valid release data found in %s", testFile.filename)
			}

			t.Logf("Successfully parsed %s: %d release(s)", testFile.filename, releaseCount)
		})
	}
}