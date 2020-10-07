package spotifyconjugator

import (
	"context"
	"fmt"

	"github.com/kn100/conjugate/conjugator"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyConjugator struct {
	SpotifyClientID     string
	SpotifyClientSecret string
}

func (spe SpotifyConjugator) Search(track conjugator.Track) ([]conjugator.Result, *conjugator.SearchError) {
	if !spe.isConfigured() {
		return nil, &conjugator.SearchError{Name: FriendlyName(), Details: "requires configuration."}
	}
	config := &clientcredentials.Config{
		ClientID:     spe.SpotifyClientID,
		ClientSecret: spe.SpotifyClientSecret,
		TokenURL:     spotify.TokenURL,
	}
	// TODO: store this token and reuse it until expiry
	token, err := config.Token(context.Background())
	if err != nil {
		return nil, &conjugator.SearchError{Name: FriendlyName(), Details: fmt.Sprintf("unable to get a token, maybe the credentials are wrong: %s", err.Error())}
	}

	client := spotify.Authenticator{}.NewClient(token)
	res, err := client.Search(track.FullTitle, spotify.SearchTypeTrack)
	if err != nil {
		return nil, &conjugator.SearchError{Name: FriendlyName(), Details: fmt.Sprintf("search failed. API limit hit? %s", err.Error())}
	}

	var results []conjugator.Result = []conjugator.Result{}
	for _, st := range res.Tracks.Tracks {
		result := conjugator.Result{
			FoundTrack: conjugator.Track{
				FullTitle: fmt.Sprintf("%s - %s (%s)", st.Name, st.Artists[0].Name, st.Album.Name),
				Title:     st.Name,
				Artists:   []string{st.Artists[0].Name}, // TODO: make this return all the artists, not just one
				Album:     st.Album.Name,
				Year:      fmt.Sprint(st.Album.ReleaseDateTime().Year()),
			},
			URI:    fmt.Sprintf("https://open.spotify.com/track/%s", st.ID.String()),
			Source: FriendlyName(),
		}

		results = append(results, result)
	}

	return results, nil
}

func (spe SpotifyConjugator) ImFeelingLucky(track conjugator.Track) (conjugator.Result, *conjugator.SearchError) {
	results, err := spe.Search(track)
	if err != nil {
		return conjugator.Result{}, err
	}
	if len(results) > 0 {
		return results[0], nil
	} else {
		return conjugator.Result{}, nil
	}

}

func (spe SpotifyConjugator) isConfigured() bool {
	return spe.SpotifyClientID != "" && spe.SpotifyClientSecret != ""
}

func FriendlyName() string {
	return "Spotify Data Conjugator"
}

func Help() string {
	return "Some help text revolving around Spotify"
}

func (spe SpotifyConjugator) RequiredConfigurationOptions() []string {
	return []string{"SpotifyClientID", "SpotifyClientSecret"}
}
