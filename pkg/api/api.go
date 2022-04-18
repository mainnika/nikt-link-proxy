package api

import (
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/valyala/fasthttp"
)

const pathKeySID = "sid"
const pathKeyLinkID = "linkID"

// api schema
const (
	URLRoot    = "/"
	URLHealthz = "/healthz"

	URLSIDBinary   = "/<" + pathKeySID + ">/binary/*"
	URLSIDMakeLink = "/<" + pathKeySID + ">/go"

	URLBinary   = "/binary/*"
	URLMakeLink = "/go"

	URLResolveLinkID = "/<" + pathKeyLinkID + ">"
)

// API is the main handler that contains all routes handlers
type API struct {
	Config
	Router *routing.Router
}

// New creates a new api handler instance
func New(config Config) (api *API) {

	api = &API{Config: config}
	api.Router = routing.New()

	api.Router.Use(api.UseJSONWriter)
	api.Router.Use(api.UseErrorHandler)
	api.Router.NotFound(api.ErrorNotFound)

	apiBase := api.Router.Group(api.Base)

	apiBase.Get(URLRoot, api.Root)

	apiBase.Get(URLHealthz, api.GetHealthz)

	apiBase.Get(URLBinary, api.UseTemplateWriter, api.GetBinary)
	apiBase.Get(URLSIDBinary, api.UseTemplateWriter, api.GetBinary)

	apiBase.To("GET,HEAD", URLMakeLink, api.GetMakeLinkWithSID)
	apiBase.To("GET,HEAD", URLSIDMakeLink, api.GetMakeLinkWithSID)

	apiBase.To("GET,HEAD", URLResolveLinkID, api.GetResolveLinkID)

	return
}

// Handler returns the api router request handler
func (api *API) Handler() (handler fasthttp.RequestHandler) {
	handler = api.Router.HandleRequest
	return
}
