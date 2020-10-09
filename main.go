package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kn100/conjugate/conjugator"
	"github.com/kn100/conjugate/spotifyconjugator"
	"github.com/kn100/conjugate/util"
	"github.com/kn100/conjugate/youtubeconjugator"
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Conjugate - converts Youtube or Youtube Music links into Spotify links.\n")
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	videoLink := flag.String("y", "", "A Youtube link (For example https://music.youtube.com/watch?v=oHg5SJYRHA0")
	raw := flag.Bool("raw", false, "Raw output (no formatting) - useful for piping")
	reset := flag.Bool("reconfigure", false, "Reconfigure conjugate")

	flag.Parse()

	conjugators := make(map[string]conjugator.Conjugator)
	conjugators["youtube"] = &youtubeconjugator.YoutubeConjugator{}
	conjugators["spotify"] = &spotifyconjugator.SpotifyConjugator{}

	for _, c := range conjugators {
		conj := c.(conjugator.Conjugator)
		if !conj.IsConfigured() || *reset == true {
			conj.Configure()
		}
	}

	if *videoLink == "" {
		Usage()
		return
	}

	track, yerr := conjugators["youtube"].(conjugator.ExtractConjugator).Extract(*videoLink)
	if yerr != nil {
		fmt.Printf(util.Failed(yerr.Error()))
		return
	}

	track.Title = util.Gravify(track.Title)
	track.FullTitle = util.Gravify(track.FullTitle)

	result, serr := conjugators["spotify"].(conjugator.SearchConjugator).ImFeelingLucky(track)
	if serr != nil {
		fmt.Printf(util.Failed(serr.Error()))
	}

	fmt.Print(util.FormatOutput(result, *raw))
}
