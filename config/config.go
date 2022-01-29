package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var conf Config

type (
	// Config define a struct for starting the http server
	Config struct {
		Mysql  Mysql  `mapstructure:"mysql"`
		Server Server `mapstructure:"server"`
		Redis  Redis  `mapstructure:"redis"`
		Chain  Chain  `mapstructure:"chain"`
	}

	// Mysql define a struct for mysql connect
	Mysql struct {
		User         string `mapstructure:"user"`
		Password     string `mapstructure:"password"`
		Host         string `mapstructure:"host"`
		Port         string `mapstructure:"port"`
		DBName       string `mapstructure:"db_name"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
		MaxLifeTime  string `mapstructure:"max_life_time"`
	}

	// Redis define a struct for redis server
	Redis struct {
		Address  string `mapstructure:"address"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}

	// Chain define a struct for Chain server
	Chain struct {
		RcpAddr          string `mapstructure:"rpc_address"`
		GrpcAddr         string `mapstructure:"grpc_address"`
		WsAddr           string `mapstructure:"ws_addr"`
		ChainID          string `mapstructure:"chain_id"`
		ProjectID        string `mapstructure:"project_id"`
		ProjectKey       string `mapstructure:"project_key"`
		ChainAccountAddr string `mapstructure:"chain_account_addr"`

		Gas       uint64 `mapstructure:"gas"`
		Denom     string `mapstructure:"denom"`
		Amount    int64  `mapstructure:"amount"`
		AccoutGas int64  `mapstructure:"account_gas"`
	}

	// Server define a struct for http server
	Server struct {
		Address   string `mapstructure:"address"`
		LogLevel  string `mapstructure:"log_level"`
		LogFormat string `mapstructure:"log_format"`
	}
)

func Get() Config {
	return conf
}

func Load(cmd *cobra.Command, home string) error {
	rootViper := viper.New()
	_ = rootViper.BindPFlags(cmd.Flags())
	// Find home directory.
	rootViper.AddConfigPath(rootViper.GetString(home))
	rootViper.SetConfigName("config")
	rootViper.SetConfigType("toml")

	// Find and read the config file
	if err := rootViper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return err
	}

	if err := rootViper.Unmarshal(&conf); err != nil {
		return err
	}
	return nil
}
