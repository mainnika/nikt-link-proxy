package api_test

import (
	"fmt"
	"net/http"
	"testing"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/api"
)

func TestError_Error(t *testing.T) {
	t.Run("expect error returns message", func(t *testing.T) {
		assert.Equal(t, "this is message", api.NewHTTPError(http.StatusForbidden, "this is message").Error())
	})
	t.Run("expect error creates message if empty", func(t *testing.T) {
		assert.Equal(t, "I'm a teapot", api.NewHTTPError(http.StatusTeapot).Error())
	})
	t.Run("expect error returns status code", func(t *testing.T) {
		assert.Equal(t, 451, api.NewHTTPError(http.StatusUnavailableForLegalReasons).StatusCode())
	})
}

func TestAPI_ErrorNotFound(t *testing.T) {
	t.Run("expect handler returns not found error", func(t *testing.T) {
		err := (&api.API{}).ErrorNotFound(&routing.Context{})
		assert.Equal(t, "Not Found", err.Error())
	})
}

func TestAPI_UseErrorHandler(t *testing.T) {
	t.Run("expect works without errors", func(t *testing.T) {

		a := api.API{}
		m := mock.Mock{}
		c := routing.NewContext(&fasthttp.RequestCtx{},
			a.UseErrorHandler,
			func(context *routing.Context) error {
				return m.MethodCalled("handler").Error(0)
			},
		)

		m.Test(t)
		m.On("handler").Return(nil)

		assert.NoError(t, c.Next())
		assert.Equal(t, 200, c.Response.StatusCode())
	})
	t.Run("expect catch error and set status", func(t *testing.T) {

		a := api.API{}
		m := mock.Mock{}
		c := routing.NewContext(&fasthttp.RequestCtx{},
			a.UseErrorHandler,
			func(context *routing.Context) error {
				return m.MethodCalled("handler").Error(0)
			},
		)

		m.Test(t)
		m.On("handler").Return(api.NewHTTPError(http.StatusBadRequest))

		assert.NoError(t, c.Next())
		assert.Equal(t, 400, c.Response.StatusCode())
	})
	t.Run("expect catch panic and set status", func(t *testing.T) {

		a := api.API{}
		c := routing.NewContext(&fasthttp.RequestCtx{},
			a.UseErrorHandler,
			func(context *routing.Context) error {
				panic(fmt.Errorf("volcano eruption"))
			},
		)

		assert.NoError(t, c.Next())
		assert.Equal(t, 500, c.Response.StatusCode())
	})
}
