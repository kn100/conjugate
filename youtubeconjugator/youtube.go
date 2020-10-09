package youtubeconjugator

import (
	"context"
	"fmt"
	"regexp"

	"github.com/kn100/conjugate/configmanager"
	"github.com/kn100/conjugate/conjugator"
	"github.com/kn100/conjugate/util"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeConjugator struct {
	config configmanager.ConfigManager
}

func (yte *YoutubeConjugator) FriendlyName() string {
	return "Youtube Data Conjugator"
}

func (yte *YoutubeConjugator) Extract(link string) (conjugator.Track, *conjugator.ExtractionError) {
	if !yte.IsConfigured() {
		return conjugator.Track{}, &conjugator.ExtractionError{Name: yte.FriendlyName(), Details: "requires configuration."}
	}

	id, ok := extractIDFromLink(link)
	if !ok {
		return conjugator.Track{}, &conjugator.ExtractionError{Name: yte.FriendlyName(), Details: "couldn't extract the ID from that URL - is it valid?"}
	}

	yts, err := youtube.NewService(context.Background(), option.WithAPIKey(yte.config.Get("youtube-data-api-key")))
	if err != nil {
		return conjugator.Track{}, &conjugator.ExtractionError{Name: yte.FriendlyName(), Details: err.Error()}
	}

	resp, err := yts.Videos.List([]string{"contentDetails", "snippet"}).Id(id).Do()
	if err != nil {
		return conjugator.Track{}, &conjugator.ExtractionError{Name: yte.FriendlyName(), Details: fmt.Sprintf("Youtube responded with an error: %s", err.Error())}
	}
	if len(resp.Items) < 1 {
		return conjugator.Track{}, &conjugator.ExtractionError{Name: yte.FriendlyName(), Details: "no results found"}
	}

	return makeTrack(resp.Items[0]), nil
}

func (yte *YoutubeConjugator) CanExtract(link string) bool {
	_, ok := extractIDFromLink(link)
	return ok
}

func extractIDFromLink(link string) (string, bool) {
	reg := regexp.MustCompile(`(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*?[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})`)
	if reg.MatchString(link) {
		substrings := reg.FindStringSubmatch(link)
		return substrings[1], true
	}
	return "", false
}

func (yte *YoutubeConjugator) IsConfigured() bool {
	if (yte.config == configmanager.ConfigManager{}) {
		cfgm := configmanager.ConfigManager{}
		err := cfgm.Init(yte.FriendlyName())
		if err != nil {
			fmt.Println(err)
			util.Failed(err.Error())
			return false
		}
		yte.config = cfgm
	}
	return yte.config.Get("youtube-data-api-key") != ""
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

func (yte *YoutubeConjugator) Configure() bool {
	cfgm := configmanager.ConfigManager{}

	err := cfgm.Init(yte.FriendlyName())
	if err != nil {
		util.Failed(err.Error())
		return false
	}

	fmt.Println("You will need to have a Google account in order to get an API key for Youtube Data. See https://developers.google.com/youtube/registering_an_application for more info")
	apiKey := util.Prompt("Enter your Youtube Data API key")

	err = cfgm.Set("youtube-data-api-key", apiKey)
	if err != nil {
		util.Failed("Failed to save api key")
		return false
	}

	util.Successful("API Key set!")
	yte.config = cfgm
	return true
}
