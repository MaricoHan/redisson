package cmd

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.bianjie.ai/avata/open-api/internal/app"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	"io/ioutil"
)

var (
	localConfig                       string
	host, port, username, pwd, prefix string

	testCmd = &cobra.Command{ // test
		Use:   "test",
		Short: "start test openapi server.",
		Run: func(cmd *cobra.Command, args []string) {
			test()
		},
	}
	startCmd = &cobra.Command{
		Use:     "start",
		Example: "start openapi server",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	testCmd.Flags().StringVarP(&localConfig, "CONFIG", "c", "", "config path: /opt/local.toml")
	rootCmd.AddCommand(testCmd)
	startCmd.PersistentFlags().StringVarP(&host, "HOST", "e", "", "etcd host address")
	startCmd.PersistentFlags().StringVarP(&port, "PORT", "p", "", "etcd prod")
	startCmd.PersistentFlags().StringVarP(&username, "USERNAME", "u", "", "etcd username")
	startCmd.PersistentFlags().StringVarP(&pwd, "PWD", "P", "", "etcd password")
	startCmd.PersistentFlags().StringVarP(&prefix, "PREFIX", "n", "", "etcd prefix")

	rootCmd.AddCommand(startCmd)
}

func run() {
	var config *configs.Config
	if host != "" && port != "" && prefix != "" {
		if username == "" || pwd == "" {
			log.Error("Error: Username Or Password null")
			return
		}
		conn, err := initialize.InitEtcdClient(host, port, username, pwd)
		if err != nil {
			log.Error("Error: Conn Etcd Error: ", err)
			return
		}
		resp, err := conn.GetEntries(prefix)
		if err != nil {
			log.Error("Error: Get Config Data Error: ", err)
			return
		}
		config, err = readConfig([]byte(resp[0]))
		if err != nil {
			log.Error("Error: Get Config Data Error: ", err)
			return
		}
	}
	configs.Cfg = *config
	app.Start(&configs.Context{
		Logger: log.New(),
		Config: config,
	})
}

func test() {
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
