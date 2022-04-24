package configs

var Cfg Config

type (
	Config struct {
		App   App   `mapstructure:"app"`
		Mysql Mysql `mapstructure:"mysql"`
		Redis Redis `mapstructure:"redis"`
		GrpcClient GrpcClient `mapstructure:"grpc_client"`
	}

	App struct {
		ServerName         string `mapstructure:"name"`
		Version            string `mapstructure:"version"`
		LogLevel           string `mapstructure:"log_level"`
		LogFormat          string `mapstructure:"log_format"`
		Addr               string `mapstructure:"addr"`
		Env                string `mapstructure:"env"`
		RouterPrefix       string `mapstructure:"router_prefix"`
		SignatureAuth      bool   `mapstructure:"signature_auth"`
		TimestampAuth      bool   `mapstructure:"timestamp_auth"`
		DefaultKeyPassword string `mapstructure:"default_key_password"`
		PrometheusAddr     string `mapstructure:"prometheus_addr"`
		GprcTimeout        int    `mapstructure:"grpc_timeout"`
	}

	Mysql struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		DB       string `mapstructure:"db"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}

	Redis struct {
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"password"`
		DB       int64  `mapstructure:"db"`
	}


	GrpcClient struct {
		WenchangchainDDCAddr string `mapstructure:"wenchangchain_ddc_addr"`
		WenchangchainNativeAddr string `mapstructure:"wenchangchain_native_addr"`
	}
)
