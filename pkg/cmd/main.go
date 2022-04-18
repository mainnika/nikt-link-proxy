package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "overridable/loader"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/cmd/serve"
)

var (
	CfgFile string = "config.yaml"
	Verbose bool   = false
	Version string = "raw"
)

var root = &cobra.Command{
	Short:   "The Nikt link proxy shortener",
	Use:     "nikt-link-proxy",
	Version: Version,
	Run:     serve.Serve,
}

func init() {

	cobra.OnInitialize(initConfig)

	flags := root.PersistentFlags()
	flags.StringP("config", "c", CfgFile, "config file path")
	flags.BoolP("verbose", "v", Verbose, "enable verbose logging")

	viper.BindPFlags(flags)
	serve.InitFlags(root.Flags())
	viper.BindPFlags(root.Flags())

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func initConfig() {

	viper.SetConfigFile(viper.GetString("config"))

	err := viper.ReadInConfig()

	if viper.GetBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Verbose mode enabled")
	}

	if err != nil {
		logrus.Debugf("Skip invalid config file %s, %v", viper.ConfigFileUsed(), err)
	} else {
		logrus.Debugf("Config file in use %s", viper.ConfigFileUsed())
	}
}

func main() {
	err := root.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}
