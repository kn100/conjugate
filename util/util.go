package util

import (
	"fmt"
	"regexp"

	"github.com/kn100/conjugate/conjugator"
)

func Prompt(prompt string) string {
	fmt.Printf(prompt + "\n")
	var input string
	fmt.Printf("> ")
	fmt.Scanln(&input)
	return input
}

func Successful(prompt string) string {
	return fmt.Sprintf("%s ✨\n", prompt)
}

func Failed(prompt string) string {
	return fmt.Sprintf("%s 😩\n", prompt)
}

// Gravy is nice
func Gravify(track string) string {
	reg := regexp.MustCompile(`(.+)\(feat.+\)(.+)`)
	if reg.MatchString(track) {
		return reg.ReplaceAllString(track, "${1}${2}")
	}
	return track
}

func FormatOutput(result conjugator.Result, raw bool) string {
	if !raw {
		if result.URI == "" {
			return Failed("No match found")
		}
		return Successful(fmt.Sprintf("Matched track: %s\nLink: %s", result.FoundTrack.FullTitle, result.URI))
	}
	return result.URI
}
