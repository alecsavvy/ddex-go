package testdata

import (
	ernv432 "github.com/alecsavvy/ddex-go/gen/ddex/ern/v432"
	meadv11 "github.com/alecsavvy/ddex-go/gen/ddex/mead/v11"
	piev10 "github.com/alecsavvy/ddex-go/gen/ddex/pie/v10"
)

// SimpleERNTest returns an ERN message for "Dark Side of the Moon" by "Pink Floyd"
func SimpleERNTest() *ernv432.NewReleaseMessage {
	return &ernv432.NewReleaseMessage{
		ReleaseProfileVersionId: "CommonReleaseTypes/14",
		LanguageAndScriptCode:   "en",
		AvsVersionId:            "4.3.2",
		MessageHeader: &ernv432.MessageHeader{
			MessageThreadId: "MSG_THREAD_001",
			MessageId:       "ERN_001_DARKSIDEOFTHEMOON",
			MessageSender: &ernv432.MessagingPartyWithoutCode{
				PartyId: "HARVEST_RECORDS_001",
				PartyName: &ernv432.PartyNameWithoutCode{
					FullName: "Harvest Records",
				},
			},
			MessageRecipient: []*ernv432.MessagingPartyWithoutCode{{
				PartyId: "SPOTIFY_001",
				PartyName: &ernv432.PartyNameWithoutCode{
					FullName: "Spotify Technology S.A.",
				},
			}},
			MessageCreatedDateTime: "2023-03-01T12:00:00.000Z",
			MessageControlType:     "LiveMessage",
		},
		PartyList: &ernv432.PartyList{
			Party: []*ernv432.Party{
				{
					PartyReference: "PINK_FLOYD_001",
				},
			},
		},
		ReleaseList: &ernv432.ReleaseList{
			Release: &ernv432.Release{
				ReleaseReference: "DSOTM_RELEASE_001",
				ReleaseId: &ernv432.ReleaseId{
					GRid: "A1HARVEST73DARKSIDEOFTHEMOON",
				},
				DisplayTitleText: []*ernv432.DisplayTitleText{
					{
						Value:                 "The Dark Side of the Moon",
						LanguageAndScriptCode: "en",
					},
				},
				DisplayArtistName: []*ernv432.DisplayArtistNameWithOriginalLanguage{
					{
						Value:                 "Pink Floyd",
						LanguageAndScriptCode: "en",
					},
				},
			},
		},
	}
}

// SimpleMEADTest returns a MEAD message for "Dark Side of the Moon" Grammy award
func SimpleMEADTest() *meadv11.MeadMessage {
	return &meadv11.MeadMessage{
		AvsVersionId:          "1.1",
		LanguageAndScriptCode: "en",
		MessageHeader: &meadv11.MessageHeader{
			MessageId: "MEAD_001_DSOTM_GRAMMY",
			MessageSender: &meadv11.MessagingPartyWithoutCode{
				PartyId: "NARAS_001",
				PartyName: &meadv11.PartyNameWithoutCode{
					FullName: "Recording Academy",
				},
			},
			MessageRecipient: []*meadv11.MessagingPartyWithoutCode{{
				PartyId: "HARVEST_RECORDS_001",
				PartyName: &meadv11.PartyNameWithoutCode{
					FullName: "Harvest Records",
				},
			}},
			MessageCreatedDateTime: "2023-02-05T18:00:00+00:00",
		},
		ReleaseInformationList: &meadv11.ReleaseInformationList{
			ReleaseInformation: []*meadv11.ReleaseInformation{
				{
					ReleaseSummary: &meadv11.ReleaseSummary{
						ReleaseId: &meadv11.ReleaseId{
							GRid: "A1HARVEST73DARKSIDEOFTHEMOON",
						},
						DisplayTitle: []*meadv11.DisplayTitle{
							{
								TitleText: &meadv11.TitleText{
									Title: "The Dark Side of the Moon",
								},
							},
						},
					},
				},
			},
		},
	}
}

// SimplePIETest returns a PIE message for Pink Floyd Grammy award
func SimplePIETest() *piev10.PieMessage {
	return &piev10.PieMessage{
		AvsVersionId:          "1.0",
		LanguageAndScriptCode: "en",
		MessageHeader: &piev10.MessageHeader{
			MessageId: "PIE_001_PINKFLOYD_GRAMMY",
			MessageSender: &piev10.MessagingPartyWithoutCode{
				PartyId: "NARAS_001",
				PartyName: &piev10.PartyNameWithoutCode{
					FullName: "Recording Academy",
				},
			},
			MessageRecipient: []*piev10.MessagingPartyWithoutCode{{
				PartyId: "HARVEST_RECORDS_001",
				PartyName: &piev10.PartyNameWithoutCode{
					FullName: "Harvest Records",
				},
			}},
			MessageCreatedDateTime: "2023-02-05T18:00:00+00:00",
		},
		PartyList: &piev10.PartyList{
			Party: []*piev10.Party{
				{
					PartyReference: "PINK_FLOYD_001",
					PartyName: []*piev10.PartyName{
						{
							FullName: &piev10.NameWithScriptCode{
								Name: &piev10.Name{
									Value: "Pink Floyd",
								},
							},
						},
					},
					Award: []*piev10.Award{
						{
							AwardName: &piev10.NameWithPronunciationAndScriptCode{
								Name: &piev10.Name{
									Value: "Grammy Award for Best Engineered Album, Non-Classical",
								},
							},
							Date: &piev10.EventDate{
								Value: "1973",
							},
							IsWinner: true,
						},
					},
				},
			},
		},
	}
}
