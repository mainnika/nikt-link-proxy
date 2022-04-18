package api

import (
	"encoding/json"
	"net/http"
	"path"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data"
)

// Root creates a root redirect.
func (api *API) Root(c *routing.Context) (httpError error) {

	rootRedirect := api.Config.RootRedirect
	if rootRedirect == "" {
		return c.Write(json.RawMessage{'n', 'u', 'l', 'l'})
	}

	shortLink, err := api.Data.MakeLink(c, rootRedirect, data.Meta{
		Ref: string(c.Referer()),
	})
	if err != nil {
		return NewHTTPError(http.StatusForbidden, err.Error())
	}

	baseRedirect := path.Join("/", api.Base, shortLink.ID)

	c.Redirect(baseRedirect, http.StatusFound)

	return c.Write(json.RawMessage{'n', 'u', 'l', 'l'})
}
