package api_test

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/api"
)

func TestAPI_GetHealthz(t *testing.T) {
	t.Run("expect no error if database ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		dataw := NewMockDataWriter(ctrl)

		defer ctrl.Finish()

		api := api.API{}
		c := routing.NewContext(&fasthttp.RequestCtx{}, api.GetHealthz)

		dataw.EXPECT().SetHeader(gomock.Any())
		dataw.EXPECT().Write(gomock.Any(), json.RawMessage(`{"ok":true}`))

		c.SetDataWriter(dataw)

		assert.NoError(t, c.Next())
	})
}
