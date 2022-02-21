package serve_test

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/cmd/serve"
)

func TestServe(t *testing.T) {
	t.Run("expect serve creates tcp http server", func(t *testing.T) {
		// allocate random tcp listener and reuse its addr after close
		listener, _ := net.Listen("tcp", "127.0.0.1:0")
		listener.Close()

		viper.Set("addr", listener.Addr().String())
		viper.Set("dataURI", "postgres://localhost/postgres")
		viper.Set("unix", "")
		viper.Set("base", "")

		go serve.Serve(&cobra.Command{Version: "testing"}, nil)

		t.Run("expect accept the connection", func(t *testing.T) {
			// due to async Serve waiting for the listener is necessary
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			for {
				_, err := net.Dial("tcp", listener.Addr().String())
				if err == nil {
					break
				}
				if ctx.Err() != nil {
					t.Fatal(ctx.Err())
				}
			}
		})
		t.Run("expect http server works", func(t *testing.T) {
			httpexpect.New(t, "http://"+listener.Addr().String()).GET("/").Expect().Headers().NotEmpty()
		})
		t.Run("expect shutdown by signal", func(t *testing.T) {
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		})
	})
}
