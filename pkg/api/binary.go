package api

import (
	"bytes"
	"fmt"
	"net/http"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"

	"overridable/loader"
)

// GetBinary serves the binary template for the given path.
func (api *API) GetBinary(c *routing.Context) (httpError error) {

	sid := c.Param(pathKeySID)
	hasSID := len(sid) > 0

	path := c.Path()
	path = bytes.TrimLeft(path, "/")

	if hasSID {
		path = path[len(sid):]
		path = bytes.TrimLeft(path, "/")
	}

	tmpl, err := loader.LookupBytes(path)
	if err != nil {
		return NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("no template: %s", path),
		)
	}

	return c.Write(TemplateData{
		Template: tmpl,
		Content: map[string]string{
			"TargetURL": "",
			"SID":       sid,
		},
	})
}
