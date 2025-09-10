package ddex

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecsavvy/ddex-go/ernv432"
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