package cmd

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.bianjie.ai/avata/open-api/internal/app"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"io/ioutil"
)

var (
	localConfig                       string

	startCmd = &cobra.Command{
		Use:     "start",
		Example: "start openapi server",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	startCmd.Flags().StringVarP(&localConfig, "config", "c", "", "config path: /opt/local.toml")
	rootCmd.AddCommand(startCmd)
}


func run() {
	data, err := ioutil.ReadFile(localConfig)
	if err != nil {
		log.Error("Error: get config data error: ", err)
		return
	}
	config, err := readConfig(data)
	if err != nil {
		log.Error("Error: read config error: ", err)
		return
	}
	configs.Cfg = *config
	app.Start(&configs.Context{
		Logger: log.New(),
		Config: config,
	})
}

func readConfig(data []byte) (*configs.Config, error) {
	v := viper.New()
	v.SetConfigType("toml")
	reader := bytes.NewReader(data)
	err := v.ReadConfig(reader)
	if err != nil {
		log.Error("Error: viper read config error: ", err)
		return nil, err
	}
	cfg := &configs.Config{}
	if err := v.Unmarshal(cfg); err != nil {
		log.Error("Error: unmarshal config error: ", err)
		return nil, err
	}
	return cfg, nil
}
