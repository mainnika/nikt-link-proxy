package serve

import (
	"github.com/spf13/pflag"
)

// Config contains serve configuration
type Config struct {
	Base string `mapstructure:"base"`
	Addr string `mapstructure:"addr"`
	Unix string `mapstructure:"unix"`
}

func InitFlags(flags *pflag.FlagSet) {
	flags.String("base", "", "http URI prefix")
	flags.String("addr", "127.0.0.1:8080", "tcp addr to listen")
	flags.String("unix", "", "unix socket path to listen")
}
