package data

// Link represents a short link
type Link struct {
	// ID is the base of a short link
	ID string
	// FullURL is the url to resolve the short link
	FullURL string
}
