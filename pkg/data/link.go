package data

import (
	"fmt"
	"net/url"
)

// Link represents a short link
type Link struct {
	// ID is the base of a short link
	ID string
	// URL is the full url to resolve the short link
	URL *url.URL
	// Meta
	Meta Meta
}

// NewLink creates a new link
func NewLink(id string, fullURL string) (l Link, err error) {

	l.ID = id
	l.URL, err = url.Parse(fullURL)

	return
}

// Valid checks if link valid or not
func (l Link) Valid() (err error) {

	if l.URL == nil {
		err = fmt.Errorf("url is not provided")
	}

	switch {
	case err != nil:
	case len(l.ID) == 0:
		err = fmt.Errorf("id is empty")
	case !IsAllowedScheme(l.URL.Scheme):
		err = fmt.Errorf("scheme is not allowed")
	case !IsAllowedHost(l.URL.Host):
		err = fmt.Errorf("host is not allowed")
	}

	return
}
