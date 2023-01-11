package configs

var Cfg Config

type (
	Config struct {
		App        App        `mapstructure:"app"`
		Mysql      Mysql      `mapstructure:"mysql"`
		Redis      Redis      `mapstructure:"redis"`
		GrpcClient GrpcClient `mapstructure:"grpc_client"`
		Project    Project    `mapstructure:"project"`
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
		HttpTimeout        int    `mapstructure:"http_timeout"`
	}

	Mysql struct {
		Host         string `mapstructure:"host"`
		Port         int    `mapstructure:"port"`
		DB           string `mapstructure:"db"`
		Username     string `mapstructure:"username"`
		Password     string `mapstructure:"password"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
		MaxLifetime  int    `mapstructure:"max_life_time"`
	}

	Redis struct {
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"password"`
		DB       int64  `mapstructure:"db"`
	}

	GrpcClient struct {
		WenchangchainDDCAddr    string `mapstructure:"wenchangchain_ddc_addr"`
		WenchangchainNativeAddr string `mapstructure:"wenchangchain_native_addr"`
		IritaOPBNativeAddr      string `mapstructure:"irita_opb_native_addr"`
		IrisHubNativeAddr       string `mapstructure:"iris_hub_native_addr"`
		StateGatewayAddr        string `mapstructure:"state_gateway_addr"`
		RightsJiangSu           string `mapstructure:"rights_jiangsu"`
	}

	Project struct {
		SecretPwd string `mapstructure:"secret_pwd"`
	}
)
