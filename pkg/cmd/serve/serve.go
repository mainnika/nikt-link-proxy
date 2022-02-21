package serve

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/api"
)

func createListener(tcpAddr, unixAddr string) (listener net.Listener, err error) {

	var netw, addr string

	switch {
	case tcpAddr != "":
		netw = "tcp"
		addr = tcpAddr
	case unixAddr != "":
		netw = "unix"
		addr = unixAddr
	default:
		err = fmt.Errorf("no address given")
		return
	}

	return net.Listen(netw, addr)
}

func stopListener(signals <-chan os.Signal, listener net.Listener) {
	_ = <-signals
	_ = listener.Close()
}

// Serve command runs http server and waits for the termination signal
func Serve(cmd *cobra.Command, args []string) {

	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		logrus.Warnf("Cannot unmarshal config, %v", err)
	}

	apiConfig := api.Config{
		Base: config.Base,
	}
	apiHandler := api.New(apiConfig)

	httpServer := fasthttp.Server{
		Logger:  logrus.StandardLogger(),
		Handler: apiHandler.Handler(),
	}
	httpListener, err := createListener(config.Addr, config.Unix)
	if err != nil {
		logrus.Fatal(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go stopListener(signals, httpListener)

	logrus.Infof("Version: %s", cmd.Version)
	logrus.Debugf("Listen: %s", httpListener.Addr().String())
	logrus.Debugf("Config: %#v", config)

	err = httpServer.Serve(httpListener)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Debugf("Stopped")
}
