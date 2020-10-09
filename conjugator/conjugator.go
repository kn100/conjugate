package conjugator

type Conjugator interface {
	FriendlyName() string
	IsConfigured() bool
	Configure() bool
}

type SearchConjugator interface {
	Conjugator
	Search(Track) ([]Result, *SearchError)
	ImFeelingLucky(Track) (Result, *SearchError)
}

type ExtractConjugator interface {
	Conjugator
	Extract(link string) (Track, *ExtractionError)
	CanExtract(link string) bool
}
type Result struct {
	FoundTrack Track
	URI        string
	Source     string
}

type Track struct {
	FullTitle string // If it is difficult to determine the title and artists
	Title     string
	Artists   []string
	Album     string
	Year      string
}

type ExtractionError struct {
	Name    string
	Details string
}

func (e *ExtractionError) Error() string {
	return e.Name + " encountered a problem during extraction: " + e.Details
}

type SearchError struct {
	Name    string
	Details string
}

func (e *SearchError) Error() string {
	return e.Name + " encountered a problem during searching: " + e.Details
}
