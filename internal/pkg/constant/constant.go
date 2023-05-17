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
	OrderTypeGas = 1

	// chain map
	TianheNative = "wenchangchain-native"
	TianheDDC    = "wenchangchain-ddc"

	TianzhouEVM    = "irita-tianzhou-evm"
	TianzhouNative = "irita-tianzhou-native"
	IritaLayer2    = "irita-tianzhou-layer2"

	// wallet
	WalletServer = "wallet-server"
	Wallet       = "wallet"
	Server       = "server"

	// stage
	Tianzhou = "irita-tianzhou"
	IritaOPB = "irita-opb"
	Native   = "native"
	Evm      = "evm"
	Layer2   = "layer2"

	// rights map
	JiangSu = "jiangsu"
	Guizhou = "guizhou"
)

// redis key
const (
	RedisPrefix      = "open-api"
	KeyProjectApikey = "project:apikey:"
	KeyChain         = "chain:"
	KeyAuth          = "auth"
)

var RightsMap = map[uint64]string{
	1: JiangSu,
	2: Guizhou,
}

const (
	MysqlProjectXServicesTable   = "t_project_x_services"
	MysqlServicesTable           = "t_services"
	MysqlServiceXPermissoinTable = "t_service_x_permissions"
	MysqlPermissoinTable         = "t_permissions"
)

//项目状态
const (
	ProjectStatusEnable  int = iota + 1 //启用
	ProjectStatusDisable                //禁用
	ProjectStatusCancel                 //注销
)

// 是否删除
const (
	// 是
	IsDelete = 1
	// 否
	IsNotDelete = 2
)

// 权限操作
const (
	ActionAllow  = 1 // 允许
	ActionReject = 2 // 拒绝
)

// 服务id
const (
	ServiceTypeWallet = 1 // 钱包服务
)
