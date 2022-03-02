package api

import (
	"net/http"
	"net/url"
	"path"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data"
)

// GetMakeLinkWithSID creates a new short link or returns an existed one for the full URL provided as query.
// Metadata values are taken into account to check full URL uniqueness:
// -- Referer header URL
// -- SID path value
// It is guaranteed that for the same pair of fullURL+metadata the api returns the same short link.
// The query returns a short link as a http-302 redirect with the link, whom will be automatically resolved to full URL.
func (api *API) GetMakeLinkWithSID(c *routing.Context) (httpError error) {

	queryString := c.QueryArgs().String()
	queryUnescaped, err := url.QueryUnescape(queryString)
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err.Error())
	}

	urlParsed, err := url.Parse(queryUnescaped)
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fullURL := urlParsed.String()
	shortLink, err := api.Data.MakeLink(c, fullURL, data.Meta{
		Ref: string(c.Referer()),
		SID: c.Param(pathKeySID),
	})
	if err != nil {
		return NewHTTPError(http.StatusForbidden, err.Error())
	}

	baseRedirect := path.Join("/", api.Base, shortLink.ID)

	c.Redirect(baseRedirect, http.StatusFound)

	return
}
