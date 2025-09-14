package testdata

import (
	ernv432 "github.com/alecsavvy/ddex-go/gen/ddex/ern/v432"
	meadv11 "github.com/alecsavvy/ddex-go/gen/ddex/mead/v11"
	piev10 "github.com/alecsavvy/ddex-go/gen/ddex/pie/v10"
)

// SimpleERNTest returns a comprehensive ERN message for "Dark Side of the Moon" by "Pink Floyd"
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
				{
					PartyReference: "HARVEST_RECORDS_001",
				},
				{
					PartyReference: "ABBEY_ROAD_STUDIOS_001",
				},
			},
		},
		ResourceList: &ernv432.ResourceList{
			SoundRecording: []*ernv432.SoundRecording{
				{
					ResourceReference: "DSOTM_TRACK_001",
					DisplayTitleText: []*ernv432.DisplayTitleText{
						{
							Value:                 "Money",
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
				{
					ResourceReference: "DSOTM_TRACK_002",
					DisplayTitleText: []*ernv432.DisplayTitleText{
						{
							Value:                 "Time",
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
			Image: []*ernv432.Image{
				{
					ResourceReference: "DSOTM_ARTWORK_001",
					DisplayTitleText: []*ernv432.DisplayTitleText{
						{
							Value:                 "Dark Side of the Moon Album Cover",
							LanguageAndScriptCode: "en",
						},
					},
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
		DealList: &ernv432.DealList{
			ReleaseDeal: []*ernv432.ReleaseDeal{
				{
					DealReleaseReference: []string{"DSOTM_RELEASE_001"},
				},
			},
		},
	}
}

// SimpleMEADTest returns a comprehensive MEAD message for "Dark Side of the Moon"
func SimpleMEADTest() *meadv11.MeadMessage {
	return &meadv11.MeadMessage{
		AvsVersionId:          "1.1",
		LanguageAndScriptCode: "en",
		SubscriptionId:        "MEAD_SUB_001",
		MessageHeader: &meadv11.MessageHeader{
			MessageThreadId: "MEAD_THREAD_001",
			MessageId:       "MEAD_001_DSOTM_COMPREHENSIVE",
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
			}, {
				PartyId: "PINK_FLOYD_MANAGEMENT_001",
				PartyName: &meadv11.PartyNameWithoutCode{
					FullName: "Pink Floyd Management",
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
							ICPN: "0724354132028",
						},
						DisplayTitle: []*meadv11.DisplayTitle{
							{
								TitleText: &meadv11.TitleText{
									Title: "The Dark Side of the Moon",
								},
							},
							{
								TitleText: &meadv11.TitleText{
									Title: "La Face Cachée de la Lune",
								},
							},
						},
					},
					GenreCategory: []*meadv11.GenreCategory{
						{
							Value: &meadv11.GenreCategoryValue{
								Value: "Rock",
							},
						},
					},
				},
			},
		}}
}

// SimplePIETest returns a comprehensive PIE message for Pink Floyd with multiple awards
func SimplePIETest() *piev10.PieMessage {
	return &piev10.PieMessage{
		AvsVersionId:          "1.0",
		LanguageAndScriptCode: "en",
		MessageHeader: &piev10.MessageHeader{
			MessageThreadId: "PIE_THREAD_001",
			MessageId:       "PIE_001_PINKFLOYD_COMPREHENSIVE",
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
			}, {
				PartyId: "PINK_FLOYD_MANAGEMENT_001",
				PartyName: &piev10.PartyNameWithoutCode{
					FullName: "Pink Floyd Management",
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
									Value:                 "Pink Floyd",
									LanguageAndScriptCode: "en",
								},
							},
						},
						{
							FullName: &piev10.NameWithScriptCode{
								Name: &piev10.Name{
									Value:                 "Пинк Флойд",
									LanguageAndScriptCode: "ru-Cyrl",
								},
							},
						},
					},
					Award: []*piev10.Award{
						{
							AwardName: &piev10.NameWithPronunciationAndScriptCode{
								Name: &piev10.Name{
									Value:                 "Grammy Award for Best Engineered Album, Non-Classical",
									LanguageAndScriptCode: "en",
								},
							},
							Date: &piev10.EventDate{
								Value: "1973",
							},
							IsWinner: true,
						},
						{
							AwardName: &piev10.NameWithPronunciationAndScriptCode{
								Name: &piev10.Name{
									Value:                 "Rock and Roll Hall of Fame Inductee",
									LanguageAndScriptCode: "en",
								},
							},
							Date: &piev10.EventDate{
								Value: "1996",
							},
							IsWinner: true,
						},
						{
							AwardName: &piev10.NameWithPronunciationAndScriptCode{
								Name: &piev10.Name{
									Value:                 "UK Music Hall of Fame Inductee",
									LanguageAndScriptCode: "en",
								},
							},
							Date: &piev10.EventDate{
								Value: "2005",
							},
							IsWinner: true,
						},
					},
				},
				{
					PartyReference: "DAVID_GILMOUR_001",
					PartyName: []*piev10.PartyName{
						{
							FullName: &piev10.NameWithScriptCode{
								Name: &piev10.Name{
									Value:                 "David Jon Gilmour",
									LanguageAndScriptCode: "en",
								},
							},
						},
					},
					Award: []*piev10.Award{
						{
							AwardName: &piev10.NameWithPronunciationAndScriptCode{
								Name: &piev10.Name{
									Value:                 "Commander of the Order of the British Empire",
									LanguageAndScriptCode: "en",
								},
							},
							Date: &piev10.EventDate{
								Value: "2003",
							},
							IsWinner: true,
						},
					},
				},
				{
					PartyReference: "ROGER_WATERS_001",
					PartyName: []*piev10.PartyName{
						{
							FullName: &piev10.NameWithScriptCode{
								Name: &piev10.Name{
									Value:                 "George Roger Waters",
									LanguageAndScriptCode: "en",
								},
							},
						},
					},
					Award: []*piev10.Award{
						{
							AwardName: &piev10.NameWithPronunciationAndScriptCode{
								Name: &piev10.Name{
									Value:                 "Polar Music Prize",
									LanguageAndScriptCode: "en",
								},
							},
							Date: &piev10.EventDate{
								Value: "2008",
							},
							IsWinner: true,
						},
					},
				},
			},
		},
	}
}
