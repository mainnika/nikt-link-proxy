package api

import (
	"encoding/json"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

// GetHealthz return http-ok if service is healthy
func (api *API) GetHealthz(c *routing.Context) (httpError error) {
	return c.Write(json.RawMessage(`{"ok":true}`))
}
