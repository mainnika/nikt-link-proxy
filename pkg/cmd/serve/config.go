package serve

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/pflag"
)

// Config contains serve configuration
type Config struct {
	RootRedirect string `mapstructure:"rootRedirect"`
	Base         string `mapstructure:"base"`
	Addr         string `mapstructure:"addr"`
	Unix         string `mapstructure:"unix"`
	Redis        struct {
		redis.UniversalOptions `mapstructure:",squash"`

		P string
		Q uint64
	} `mapstructure:"redis"`
}

func InitFlags(flags *pflag.FlagSet) {
	flags.String("rootRedirect", "", "redirect the api root to this url")
	flags.String("base", "", "prefix for the api base URI")
	flags.String("addr", "127.0.0.1:8080", "tcp addr to listen")
	flags.String("unix", "", "unix socket path to listen")
	flags.StringToString("redis", map[string]string{}, "redis configuration map as described here https://pkg.go.dev/github.com/go-redis/redis/v8#UniversalOptions")
}
