package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	conf Config
)

type (
	// Config define a struct for starting the http server
	Config struct {
		Mysql  Mysql  `mapstructure:"mysql"`
		Server Server `mapstructure:"server"`
		Redis  Redis  `mapstructure:"redis"`
		Chain  Chain  `mapstructure:"chain"`
		DDC    DDC    `mapstructure:"ddc"`
		BSN    BSN    `mapstructure:"bsn"`
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

		GasCoefficient float64 `mapstructure:"gas_coefficient"`
		Gas            uint64  `mapstructure:"gas"`
		Denom          string  `mapstructure:"denom"`
		Amount         int64   `mapstructure:"amount"`
		AccoutGas      int64   `mapstructure:"account_gas"`

		ChainEncryption string `mapstructure:"chain_encryption"`
	}

	// Server define a struct for http server
	Server struct {
		Address            string `mapstructure:"address"`
		PrometheusAddr     string `mapstructure:"prometheus_addr"`
		LogLevel           string `mapstructure:"log_level"`
		LogFormat          string `mapstructure:"log_format"`
		Env                string `mapstructure:"app_env"`
		RouterPrefix       string `mapstructure:"router_prefix"`
		SignatureAuth      bool   `mapstructure:"signature_auth"`
		DefaultKeyPassword string `mapstructure:"default_key_password"`
		AccountWhiteList   string `mapstructure:"account_white_list"`
		AccountCount       string `mapstructure:"account_count"`
	}

	DDC struct {
		DDCAuthorityAddress string `mapstructure:"ddc_authority_address"`
		DDCChargeAddress    string `mapstructure:"ddc_charge_address"`
		DDC721Address       string `mapstructure:"ddc_721_address"`
		DDC1155Address      string `mapstructure:"ddc_1155_address"`
		DDCGatewayUrl       string `mapstructure:"ddc_gateway_url"`

		RcpAddr          string `mapstructure:"rpc_address"`
		GrpcAddr         string `mapstructure:"grpc_address"`
		WsAddr           string `mapstructure:"ws_addr"`
		ChainID          string `mapstructure:"chain_id"`
		ProjectID        string `mapstructure:"project_id"`
		ProjectKey       string `mapstructure:"project_key"`
		ChainAccountAddr string `mapstructure:"chain_account_addr"`

		GasCoefficient float64 `mapstructure:"gas_coefficient"`
		Gas            uint64  `mapstructure:"gas"`
		Denom          string  `mapstructure:"denom"`
		Amount         int64   `mapstructure:"amount"`
		AccoutGas      int64   `mapstructure:"account_gas"`

		ChainEncryption string `mapstructure:"ddc_encryption"`
	}

	BSN struct {
		BSNUrl             string `mapstructure:"bsn_url"`
		APIAddress         string `mapstructure:"api_address"`
		APIToken           string `mapstructure:"api_token"`
		OPBChainClientType string `mapstructure:"opb_chain_client_type"`
		OPBChainID         string `mapstructure:"opb_chain_id"`
		OPBKeyType         string `mapstructure:"opb_key_type"`
		OpenDDC            string `mapstructure:"open_ddc"`
		Proof              string `mapstructure:"proof"`
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
