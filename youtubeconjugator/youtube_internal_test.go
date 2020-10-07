package youtubeconjugator

import (
	"testing"

	"github.com/kn100/conjugate/conjugator"
	"google.golang.org/api/youtube/v3"
	"gotest.tools/assert"
)

func TestExtractIDFromLink(t *testing.T) {
	testCases := []struct {
		name string
		link string
		ok   bool
		res  string
	}{
		{
			"Standard Youtube link, unlicensed",
			"https://www.youtube.com/watch?v=m3wzpC2o42I",
			true,
			"m3wzpC2o42I",
		},
		{
			"Youtube Music link, with tracking garbage, licensed",
			"https://music.youtube.com/watch?v=-kSn0kgn1rY&feature=share",
			true,
			"-kSn0kgn1rY",
		},
		{
			"Youtube link, with tracking garbage, licensed",
			"https://youtube.com/watch?v=9dVHT_AQYdM?youare=theproduct",
			true,
			"9dVHT_AQYdM",
		},
		{
			"Youtube shorthand link, licensed",
			"https://youtu.be/yN-1Z-yE0EA",
			true,
			"yN-1Z-yE0EA",
		},
		{
			"Youtube Music link, with tracking garbage, unlicensed",
			"https://music.youtube.com/watch?v=B77wrds5-vo&feature=share",
			true,
			"B77wrds5-vo",
		},
		{
			"Garbage link, wrong domain",
			"https://notyoutube.com/watch?v=B77wrds5",
			false,
			"",
		},
		{
			"Garbage link, youtube profile page",
			"https://www.youtube.com/channel/UCSnULYPo-BA8zo5IcG4doIQ",
			false,
			"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, ok := extractIDFromLink(tc.link)
			assert.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.res, id)
		})
	}
}

func MockYoutubeVideo(licensed bool, title, artist string) *youtube.Video {
	return &youtube.Video{
		ContentDetails: &youtube.VideoContentDetails{
			LicensedContent: licensed,
		},
		Snippet: &youtube.VideoSnippet{
			Title: title,
			Tags:  []string{artist},
		},
	}
}

func TestMakeTrack(t *testing.T) {
	testCases := []struct {
		name  string
		video *youtube.Video
		track conjugator.Track
	}{
		{
			"Perfectly normal, licensed content",
			MockYoutubeVideo(true, "Go with me (Original Mix)", "Boston"),
			conjugator.Track{
				Title:     "Go with me (Original Mix)",
				Artists:   []string{"Boston"},
				FullTitle: "Go with me (Original Mix)",
			},
		},
		{
			"Perfectly normal, unlicensed content",
			MockYoutubeVideo(false, "Crab Rave [Monstercat Release]", "Moo"),
			conjugator.Track{
				FullTitle: "Crab Rave [Monstercat Release]",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			track := makeTrack(tc.video)
			assert.DeepEqual(t, tc.track, track)
		})
	}
}
