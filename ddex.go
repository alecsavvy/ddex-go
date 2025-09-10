package ddex

import (
	"github.com/alecsavvy/ddex-go/ernv432"
	"github.com/alecsavvy/ddex-go/meadv11"
	"github.com/alecsavvy/ddex-go/piev10"
)

// Versioned type aliases for discoverability
type (
	// ERN v4.3.2 types
	NewReleaseMessageV432   = ernv432.NewReleaseMessage
	PurgeReleaseMessageV432 = ernv432.PurgeReleaseMessage

	// MEAD v1.1 types
	MeadMessageV11 = meadv11.MeadMessage

	// PIE v1.0 types
	PieMessageV10        = piev10.PieMessage
	PieRequestMessageV10 = piev10.PieRequestMessage
)
