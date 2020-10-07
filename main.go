package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/kn100/conjugate/conjugator"
	"github.com/kn100/conjugate/spotifyconjugator"
	"github.com/kn100/conjugate/youtubeconjugator"
	"github.com/spf13/viper"
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Conjugate - converts Youtube or Youtube Music links into Spotify links.\n")
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	flag.PrintDefaults()
}

func main() {
	viper.SetConfigName(".conjugator-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config")
	viper.SafeWriteConfig() // write config if it doesn't exist

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf(failed(fmt.Sprintf("Couldn't read config from $HOME/.config/.conjugator-config, err: %s", err.Error())))
	}

	if viper.GetBool("configured") == false {
		configure()
	}

	youtubeAPIKey := viper.GetString("youtube-data-api-key")
	spotifyClientID := viper.GetString("spotify-client-id")
	spotifyClientSecret := viper.GetString("spotify-client-secret")

	videoLink := flag.String("y", "", "A Youtube link (For example https://music.youtube.com/watch?v=oHg5SJYRHA0")
	raw := flag.Bool("raw", false, "Raw output (no formatting) - useful for piping")
	reset := flag.Bool("reconfigure", false, "Reconfigure conjugate")

	flag.Parse()

	if *reset == true {
		configure()
	}

	if *videoLink == "" {
		Usage()
		return
	}

	yte := youtubeconjugator.YoutubeConjugator{YoutubeDataAPIKey: youtubeAPIKey}

	track, yerr := yte.Extract(*videoLink)
	if yerr != nil {
		fmt.Printf(failed(yerr.Error()))
		return
	}

	track.Title = gravify(track.Title)
	track.FullTitle = gravify(track.FullTitle)

	spe := spotifyconjugator.SpotifyConjugator{SpotifyClientID: spotifyClientID, SpotifyClientSecret: spotifyClientSecret}

	result, serr := spe.ImFeelingLucky(track)
	if serr != nil {
		fmt.Printf(failed(serr.Error()))
	}

	fmt.Print(formatOutput(result, *raw))
}

// Gravy is nice
func gravify(track string) string {
	reg := regexp.MustCompile(`(.+)\(feat.+\)(.+)`)
	if reg.MatchString(track) {
		return reg.ReplaceAllString(track, "${1}${2}")
	}
	return track
}

func formatOutput(result conjugator.Result, raw bool) string {
	if !raw {
		if result.URI == "" {
			return failed("No match found")
		}
		return successful(fmt.Sprintf("Matched track: %s\nLink: %s", result.FoundTrack.FullTitle, result.URI))
	}
	return result.URI
}

func configureYoutube() bool {
	fmt.Println("You will need to have a Google account in order to get an API key for Youtube Data. See https://developers.google.com/youtube/registering_an_application for more info")
	apiKey := prompt("Enter your Youtube Data API key")
	viper.Set("youtube-data-api-key", apiKey)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf(failed(fmt.Sprintf("Something went wrong setting that API key: %v", err)))
		return false
	}
	successful("API Key set!")
	return true
}

func configure() bool {
	fmt.Println("Before you continue, some configuration is required.")
	if configureSpotify() && configureYoutube() {
		viper.Set("configured", true)
		viper.WriteConfig()
		return true
	}
	return false
}

func configureSpotify() bool {
	fmt.Println("You will need to have a (free or paid) Spotify account in order to get API Access to Spotify. Visit https://developer.spotify.com/documentation/general/guides/app-settings/ to find out how to generate these credentials")
	clientID := prompt("Enter your Spotify clientID")
	clientSecret := prompt("Enter your Spotify clientSecret")
	viper.Set("spotify-client-id", clientID)
	viper.Set("spotify-client-secret", clientSecret)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf(failed(fmt.Sprintf("Something went wrong setting that API key: %v", err)))
		return false
	}
	successful("API Key set!")
	return true
}

func prompt(prompt string) string {
	fmt.Printf(prompt + "\n")
	var input string
	fmt.Printf("> ")
	fmt.Scanln(&input)
	return input
}

func successful(prompt string) string {
	return fmt.Sprintf("%s ✨\n", prompt)
}

func failed(prompt string) string {
	return fmt.Sprintf("%s 😩\n", prompt)
}
