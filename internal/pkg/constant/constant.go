package constant

var GrpcTimeout int

const (
	//ENV
	EnvPro   = "prod"
	EnvDev   = "dev"
	EnvLocal = "local"

	//LogLevel
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"

	//Time
	TimeLayout = "2006-01-02 15:04:05"

	Delete = "DELETE"

	//orderType
	OrderTypeGas      = "gas"
	OrderTypeBusiness = "business"

	//chain map
	WenchangNative = "wenchangchain-native"
	WenchangDDC    = "wenchangchain-ddc"
	DatangNative   = "datangchain-native"

	//etcdSchema
	Schema = "avata"

	//etcdEndpoint
	WenchangNativeEndpoint = "services/chains/wenchangchain-native"
	WenchangDDCEndpoint    = "services/chains/wenchangchain-ddc"
	DatangNatveEndpoint    = "services/chains/datangchain-native"

	// Enum values for NFTSStatus
	NFTSStatusActive = "active"
	NFTSStatusBurned = "burned"
)

// redis key
const (
	RedisPrefix      = "open-api"
	KeyProjectApikey = "project:apikey:"
	KeyChain         = "chain:"
)
