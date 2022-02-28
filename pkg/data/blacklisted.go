package data

import "strings"

var (
	allowedschemes   = [...]string{"http", "https", "mailto"}
	blacklistedhosts = [...]string{"localhost", "", "nikt.tk"}
)

// IsAllowedScheme returns ok if schema is allowed to be proceeded
func IsAllowedScheme(scheme string) (ok bool) {

	for i := range allowedschemes {
		if scheme == allowedschemes[i] {
			ok = true
			break
		}
	}

	return
}

// IsAllowedHost returns ok if host is allowed to be proceeded
func IsAllowedHost(host string) (ok bool) {

	ok = true
	host = strings.TrimSpace(host)

	for i := range blacklistedhosts {
		if host == blacklistedhosts[i] {
			ok = false
			break
		}
	}

	return
}
