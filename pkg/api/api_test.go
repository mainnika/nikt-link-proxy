package api_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/api"
)

func TestNew(t *testing.T) {
	t.Run("expect api just works", func(t *testing.T) {
		api := api.New(api.Config{Base: "/foobar"})

		expect := httpexpect.WithConfig(httpexpect.Config{
			Reporter: httpexpect.NewAssertReporter(t),
			Client: &http.Client{
				Transport: httpexpect.NewFastBinder(api.Handler()),
				Jar:       httpexpect.NewJar(),
			},
		})

		expect.GET("/foobar").Expect().JSON()
	})
}
