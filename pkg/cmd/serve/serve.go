package serve

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/api"
	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data"
	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data/datasource"
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

	serveCtx, cancelCtx := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelCtx()

	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		logrus.Warnf("Cannot unmarshal config, %v", err)
	}

	logrus.Infof("Version: %s", cmd.Version)
	logrus.Debugf("Config: %#v", config)

	redisClient := redis.NewUniversalClient(&config.Redis.UniversalOptions)
	redisSource := datasource.NewRedisSource(redisClient, datasource.WithRedisPQ(config.Redis.P, config.Redis.Q))
	apiData := data.NewData(data.WithDataSource(redisSource))

	err = redisSource.Sync(serveCtx)
	if err != nil {
		logrus.Fatalf("Cannot sync redis source, %v", err)
	}

	apiHandler := api.New(api.Config{
		RootRedirect: config.RootRedirect,
		Base:         config.Base,
		Data:         apiData,
	})

	httpServer := fasthttp.Server{
		Name:    "stm32f103c8t6",
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

	logrus.Debugf("Listen: %s", httpListener.Addr().String())

	err = httpServer.Serve(httpListener)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Debugf("Stopped")
}
