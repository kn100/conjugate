package spotifyconjugator

import (
	"context"
	"fmt"
	"time"

	"github.com/kn100/conjugate/configmanager"
	"github.com/kn100/conjugate/conjugator"
	"github.com/kn100/conjugate/util"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyConjugator struct {
	config configmanager.ConfigManager
}

func (spe *SpotifyConjugator) Search(track conjugator.Track) ([]conjugator.Result, *conjugator.SearchError) {
	if !spe.IsConfigured() {
		return nil, &conjugator.SearchError{Name: spe.FriendlyName(), Details: "requires configuration."}
	}
	token, err := spe.getToken()
	if err != nil {
		return nil, &conjugator.SearchError{Name: spe.FriendlyName(), Details: fmt.Sprintf("unable to get token. %s", err.Error())}
	}

	client := spotify.Authenticator{}.NewClient(token)
	res, err := client.Search(track.FullTitle, spotify.SearchTypeTrack)
	if err != nil {
		return nil, &conjugator.SearchError{Name: spe.FriendlyName(), Details: fmt.Sprintf("search failed. API limit hit? %s", err.Error())}
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
			Source: spe.FriendlyName(),
		}

		results = append(results, result)
	}

	return results, nil
}

func (spe *SpotifyConjugator) ImFeelingLucky(track conjugator.Track) (conjugator.Result, *conjugator.SearchError) {
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

func (spe *SpotifyConjugator) IsConfigured() bool {
	if (spe.config == configmanager.ConfigManager{}) {
		cfgm := configmanager.ConfigManager{}
		err := cfgm.Init(spe.FriendlyName())
		if err != nil {
			util.Failed(err.Error())
			return false
		}
		spe.config = cfgm
	}
	return spe.config.Get("spotify-client-id") != "" && spe.config.Get("spotify-client-secret") != ""
}

func (spe *SpotifyConjugator) FriendlyName() string {
	return "Spotify Data Conjugator"
}

func (spe *SpotifyConjugator) Configure() bool {
	cfgm := configmanager.ConfigManager{}

	err := cfgm.Init(spe.FriendlyName())
	if err != nil {
		util.Failed(err.Error())
		return false
	}

	fmt.Println("You will need to have a (free or paid) Spotify account in order to get API Access to Spotify. Visit https://developer.spotify.com/documentation/general/guides/app-settings/ to find out how to generate these credentials")
	clientID := util.Prompt("Enter your Spotify clientID")
	clientSecret := util.Prompt("Enter your Spotify clientSecret")

	err = cfgm.Set("spotify-client-id", clientID)
	if err != nil {
		util.Failed("Failed to save client ID")
		return false
	}

	err = cfgm.Set("spotify-client-secret", clientSecret)
	if err != nil {
		util.Failed("Failed to save client secret")
		return false
	}

	util.Successful("API Key set!")
	spe.config = cfgm
	return true
}

func (spe *SpotifyConjugator) getToken() (*oauth2.Token, error) {
	var expiresAt time.Time
	err := expiresAt.UnmarshalText([]byte(spe.config.Get("spotify-client-access-token-expires")))
	if err != nil || time.Now().After(expiresAt) {
		config := &clientcredentials.Config{
			ClientID:     spe.config.Get("spotify-client-id"),
			ClientSecret: spe.config.Get("spotify-client-secret"),
			TokenURL:     spotify.TokenURL,
		}
		token, err := config.Token(context.Background())
		if err != nil {
			return nil, &conjugator.SearchError{Name: spe.FriendlyName(), Details: fmt.Sprintf("unable to get a token, maybe the credentials are wrong: %s", err.Error())}
		}
		spe.config.Set("spotify-client-access-token", token.AccessToken)
		expiry, err := token.Expiry.MarshalText()
		if err != nil {
			return nil, &conjugator.SearchError{Name: spe.FriendlyName(), Details: fmt.Sprintf("unable to get expiry from token: %s", err.Error())}
		}
		spe.config.Set("spotify-client-access-token-expires", string(expiry))
	}

	return &oauth2.Token{
		AccessToken: spe.config.Get("spotify-client-access-token"),
		Expiry:      expiresAt,
		TokenType:   "Bearer",
	}, nil
}
