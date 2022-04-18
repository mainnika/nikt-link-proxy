package api

import (
	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data"
)

// Config contains externally configurable data
type Config struct {
	RootRedirect string
	Base         string
	Data         data.DataInterface
}
