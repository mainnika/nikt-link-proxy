package api

import (
	"bytes"
	"html/template"
	"io"
	"mime"
	"net/http"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/valyala/fasthttp"
)

var _ routing.DataWriter = (*TemplateWriter)(nil)

// tmplWriter is thread-safe static instance of template writer
var tmplWriter = &TemplateWriter{}

// TemplateWriter is the fasthttp data writer that loads and executes template using the content
type TemplateWriter struct{}

type TemplateData struct {
	*template.Template
	Content interface{}
}

// SetHeader sets the header for the response
func (tw *TemplateWriter) SetHeader(rh *fasthttp.ResponseHeader) {
	// nothing to do
}

// Write executes the template and writes result to the response writer
func (tw *TemplateWriter) Write(w io.Writer, data interface{}) error {

	td, hasTemplate := data.(TemplateData)
	if !hasTemplate {
		return NewHTTPError(http.StatusInternalServerError, "cannot write non template result")
	}

	return td.Execute(w, td.Content)
}

// UseTemplateWriter is the routing middleware to set the default data writer
func (api *API) UseTemplateWriter(c *routing.Context) (_ error) {

	contentType := routing.MIME_HTML

	pth := c.Path()
	extIndex := bytes.LastIndexByte(pth, '.')
	if extIndex > 0 {
		contentType = mime.TypeByExtension(string(pth[extIndex:]))
	}

	c.SetContentType(contentType)
	c.SetDataWriter(tmplWriter)

	return
}
