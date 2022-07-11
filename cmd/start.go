package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.bianjie.ai/avata/open-api/internal/app"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
)

var (
	localConfig string
	serverPort  string
	startCmd    = &cobra.Command{
		Use:     "start",
		Example: "start order server",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	startCmd.Flags().StringVarP(&localConfig, "config", "c", "", "config path: /opt/local.toml")
	startCmd.Flags().StringVarP(&serverPort, "port", "p", "", "config path: /opt/local.toml")
	rootCmd.AddCommand(startCmd)
}

func run() {
	v := viper.New()
	// Find home directory.
	v.AddConfigPath(localConfig)
	v.SetConfigName("config")
	v.SetConfigType("toml")

	// Find and read the config file
	if err := v.ReadInConfig(); err != nil { // Handle errors reading the config file
		log.Errorf("read config err:%s", err.Error())
		return
	}
	var config configs.Config
	if err := v.Unmarshal(&config); err != nil {
		log.Errorf("unmarshal config err:%s", err.Error())
		return
	}
	if serverPort != "" {
		config.App.Addr = ":" + serverPort
	}
	configs.Cfg = config
	app.Start(&configs.Context{
		Logger: log.New(),
		Config: &config,
	})
}
