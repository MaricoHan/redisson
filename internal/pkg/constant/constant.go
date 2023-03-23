package constant

var GrpcTimeout int

const (
	// ENV
	EnvPro   = "prod"
	EnvDev   = "dev"
	EnvLocal = "local"

	// LogLevel
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"

	// Time
	TimeLayout = "2006-01-02 15:04:05"

	Delete = "DELETE"

	// orderType
	OrderTypeGas      = "gas"
	OrderTypeBusiness = "business"

	// chain map
	WenchangNative = "wenchangchain-native"
	WenchangDDC    = "wenchangchain-ddc"
	DatangNative   = "datangchain-native"
	IritaOPBEVM    = "irita-opb-evm"
	IritaOPBNative = "irita-opb-native"
	IrisHubNative  = "irishub-native"

	// wallet
	WalletServer = "wallet-server"
	Wallet       = "wallet"
	Server       = "server"

	// stage
	IritaOPB = "irita-opb"
	Native   = "native"
	EVM      = "evm"

	// rights map
	JiangSu = "jiangsu"
	Guizhou = "guizhou"
)

// redis key
const (
	RedisPrefix           = "open-api"
	KeyProjectApikey      = "project:apikey:"
	KeyChain              = "chain:"
	KeyExistWalletService = "project:wallet:"
)

var RightsMap = map[uint64]string{
	1: JiangSu,
	2: Guizhou,
}
