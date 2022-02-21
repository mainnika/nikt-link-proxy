package api_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/api"
)

func TestJSONWriter_SetHeader(t *testing.T) {
	t.Run("expect writer set json content type header", func(t *testing.T) {
		jsonWriter := api.JSONWriter{}

		rh := fasthttp.AcquireResponse()
		jsonWriter.SetHeader(&rh.Header)

		assert.Equal(t, "application/json", string(rh.Header.ContentType()))
	})
}

func TestJSONWriter_Write(t *testing.T) {
	t.Run("expect json writer encodes easydata to json", func(t *testing.T) {
		jsonWriter := api.JSONWriter{}

		buf := bytes.NewBuffer(nil)

		assert.NoError(t, jsonWriter.Write(buf, api.NewHTTPError(http.StatusTeapot)))
		assert.JSONEq(t, "{\"code\":418,\"message\":\"I'm a teapot\"}", buf.String())
	})
	t.Run("expect json writer encodes any data to json", func(t *testing.T) {
		jsonWriter := api.JSONWriter{}

		buf := bytes.NewBuffer(nil)

		assert.NoError(t, jsonWriter.Write(buf, map[string]string{"foo": "bar"}))
		assert.JSONEq(t, "{\"foo\":\"bar\"}", buf.String())
	})
}
