package conjugator

type Conjugator interface {
	FriendlyName() string
	Extract(link string) (Track, *ExtractionError)
	CanExtract(link string) bool
	RequiredConfigurationOptions() []string
	ImFeelingLucky(Track) (Result, *SearchError)
	Search(Track) ([]Result, *SearchError)
	Help() string
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
