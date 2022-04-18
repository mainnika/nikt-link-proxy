package api

import (
	"encoding/json"
	"io"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jwriter"
	"github.com/valyala/fasthttp"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/utils"
)

var _ routing.DataWriter = (*JSONWriter)(nil)

// jsonWriter is thread-safe static instance of JSON writer
var jsonWriter = &JSONWriter{}

// null ready-to-use response
var nullBytes = []byte("null")

// JSONWriter is the fasthttp data writer that marshals the output to JSON
type JSONWriter struct{}

// SetHeader sets the content type to JSON
func (jw *JSONWriter) SetHeader(rh *fasthttp.ResponseHeader) {
	rh.SetContentType(routing.MIME_JSON)
}

// Write marshals the content to JSON and writes it to out writer
func (jw *JSONWriter) Write(w io.Writer, content interface{}) (err error) {

	if utils.IsNilInterface(content) {
		_, err = w.Write(nullBytes)
		return
	}

	easyMarshaler, isEasyMarshaler := content.(easyjson.Marshaler)
	if isEasyMarshaler {
		easyWriter := &jwriter.Writer{Flags: jwriter.NilMapAsEmpty | jwriter.NilSliceAsEmpty}
		easyMarshaler.MarshalEasyJSON(easyWriter)
		_, err = easyWriter.DumpTo(w)
		return
	}

	err = json.NewEncoder(w).Encode(content)

	return
}

// UseJSONWriter is the routing middleware to set the default data writer
func (api *API) UseJSONWriter(c *routing.Context) (_ error) {
	c.SetDataWriter(jsonWriter)
	return
}
