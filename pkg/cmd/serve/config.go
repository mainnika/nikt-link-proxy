package serve

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/pflag"
)

// Config contains serve configuration
type Config struct {
	Base  string                  `mapstructure:"base"`
	Addr  string                  `mapstructure:"addr"`
	Unix  string                  `mapstructure:"unix"`
	Redis *redis.UniversalOptions `mapstructure:"redis"`
}

func InitFlags(flags *pflag.FlagSet) {
	flags.String("base", "", "http URI prefix")
	flags.String("addr", "127.0.0.1:8080", "tcp addr to listen")
	flags.String("unix", "", "unix socket path to listen")
	flags.StringToString("redis", map[string]string{}, "redis configuration map as described here https://pkg.go.dev/github.com/go-redis/redis/v8#UniversalOptions")
}
