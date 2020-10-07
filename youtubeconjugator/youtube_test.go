package youtubeconjugator_test

import (
	"testing"

	"github.com/kn100/conjugate/youtubeconjugator"
	"gotest.tools/assert"
)

func TestCanExtract(t *testing.T) {
	testCases := []struct {
		name string
		link string
		res  bool
	}{
		{
			"Standard Youtube link, unlicensed",
			"https://www.youtube.com/watch?v=m3wzpC2o42I",
			true,
		},
		{
			"Youtube Music link, with tracking garbage, licensed",
			"https://music.youtube.com/watch?v=-kSn0kgn1rY&feature=share",
			true,
		},
		{
			"Youtube link, with tracking garbage, licensed",
			"https://youtube.com/watch?v=9dVHT_AQYdM?youare=theproduct",
			true,
		},
		{
			"Youtube shorthand link, licensed",
			"https://youtu.be/yN-1Z-yE0EA",
			true,
		},
		{
			"Youtube Music link, with tracking garbage, unlicensed",
			"https://music.youtube.com/watch?v=B77wrds5-vo&feature=share",
			true,
		},
		{
			"Garbage link, wrong domain",
			"https://notyoutube.com/watch?v=B77wrds5",
			false,
		},
		{
			"Garbage link, youtube profile page",
			"https://www.youtube.com/channel/UCSnULYPo-BA8zo5IcG4doIQ",
			false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			yte := youtubeconjugator.YoutubeConjugator{}
			assert.Equal(t, tc.res, yte.CanExtract(tc.link))
		})
	}
}
