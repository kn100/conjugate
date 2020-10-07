package youtubeconjugator

import (
	"context"
	"fmt"
	"regexp"

	"github.com/kn100/conjugate/conjugator"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeConjugator struct {
	YoutubeDataAPIKey string
}

func FriendlyName() string {
	return "Youtube Data Conjugator"
}

func (yte YoutubeConjugator) Extract(link string) (conjugator.Track, *conjugator.ExtractionError) {
	if !yte.isConfigured() {
		return conjugator.Track{}, &conjugator.ExtractionError{Name: FriendlyName(), Details: "requires configuration."}
	}

	id, ok := extractIDFromLink(link)
	if !ok {
		return conjugator.Track{}, &conjugator.ExtractionError{Name: FriendlyName(), Details: "couldn't extract the ID from that URL - is it valid?"}
	}

	yts, err := youtube.NewService(context.Background(), option.WithAPIKey(yte.YoutubeDataAPIKey))
	if err != nil {
		return conjugator.Track{}, &conjugator.ExtractionError{Name: FriendlyName(), Details: err.Error()}
	}

	resp, err := yts.Videos.List([]string{"contentDetails", "snippet"}).Id(id).Do()

	if err != nil {
		return conjugator.Track{}, &conjugator.ExtractionError{Name: FriendlyName(), Details: fmt.Sprintf("Youtube responded with an error: %s", err.Error())}
	}
	if len(resp.Items) < 1 {
		return conjugator.Track{}, &conjugator.ExtractionError{Name: FriendlyName(), Details: "no results found"}
	}

	return makeTrack(resp.Items[0]), nil
}

func (yte YoutubeConjugator) CanExtract(link string) bool {
	_, ok := extractIDFromLink(link)
	return ok
}

func (yte YoutubeConjugator) RequiredConfigurationOptions() []string {
	return []string{"YoutubeDataAPIKey"}
}

func (yte YoutubeConjugator) Help() string {
	return "SomeHelpText"
}

func extractIDFromLink(link string) (string, bool) {
	reg := regexp.MustCompile(`(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*?[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})`)
	if reg.MatchString(link) {
		substrings := reg.FindStringSubmatch(link)
		return substrings[1], true
	}
	return "", false
}

func (yte YoutubeConjugator) isConfigured() bool {
	return yte.YoutubeDataAPIKey != ""
}

func makeTrack(video *youtube.Video) conjugator.Track {
	track := conjugator.Track{}
	//optimisticRegex := regexp.MustCompile(`^Provided to YouTube by.+\\n\\n(.+) · (.+)\\n\\n(.+)\\n\\n℗ (\d{4})`)
	if video.ContentDetails.LicensedContent && len(video.Snippet.Tags) > 0 {
		track.Title = video.Snippet.Title
		track.Artists = []string{video.Snippet.Tags[0]}
	}
	track.FullTitle = video.Snippet.Title
	return track
}
