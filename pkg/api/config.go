package api

import (
	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data"
)

// Config contains externally configurable data
type Config struct {
	Base string
	Data data.DataInterface
}
