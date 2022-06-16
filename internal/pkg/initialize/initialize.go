package initialize

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	pb_account "gitlab.bianjie.ai/avata/chains/api/pb/account"
	pb_business "gitlab.bianjie.ai/avata/chains/api/pb/buy"
	pb_class "gitlab.bianjie.ai/avata/chains/api/pb/class"
	pb_msgs "gitlab.bianjie.ai/avata/chains/api/pb/msgs"
	pb_nft "gitlab.bianjie.ai/avata/chains/api/pb/nft"
	pb_tx "gitlab.bianjie.ai/avata/chains/api/pb/tx"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/redis"
	"gitlab.bianjie.ai/avata/open-api/pkg/logs"
)

var RedisClient *redis.RedisClient
var MysqlDB *gorm.DB
var GrpcConnMap map[string]*grpc.ClientConn
var AccountClientMap map[string]pb_account.AccountClient
var BusineessClientMap map[string]pb_business.BuyClient
var MsgsClientMap map[string]pb_msgs.MSGSClient
var NftClientMap map[string]pb_nft.NFTClient
var ClassClientMap map[string]pb_class.ClassClient
var TxClientMap map[string]pb_tx.TxClient

func Logger(cfg *configs.Config) *log.Logger {
	if cfg.App.Env == constant.EnvPro {
		log.SetFormatter(&log.JSONFormatter{})
	}
	switch cfg.App.LogLevel {
	case constant.LogLevelDebug:
		log.SetLevel(log.DebugLevel)
	case constant.LogLevelWarn:
		log.SetLevel(log.WarnLevel)
	case constant.LogLevelError:
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	return log.StandardLogger()
}

func InitMysqlDB(cfg *configs.Config, logger *log.Logger) {
	gormLogger := logs.NewGormLogger(logger)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Mysql.Username,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.DB)

	mysqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.Fatal("init mysqlDB failed: ", err.Error())
	}
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		log.Fatal("init sqlDB failed: ", err.Error())
	}
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(cfg.Mysql.MaxOpenConns)
	// 设置数据库缓存池大小
	sqlDB.SetMaxIdleConns(cfg.Mysql.MaxIdleConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(cfg.Mysql.MaxLifetime))

	MysqlDB = mysqlDB
}

func InitGrpcClient(cfg *configs.Config, logger *log.Logger) {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}
	GrpcConnMap = make(map[string]*grpc.ClientConn)
	myLogger := NewMyLoggerV2(99, logger)
	grpclog.SetLoggerV2(&myLogger)
	wenNativeConn, err := grpc.Dial(cfg.GrpcClient.WenchangchainNativeAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		logger.Fatal("get wenchangchain-native grpc connect failed, err: ", err.Error())
	}
	GrpcConnMap[constant.WenchangNative] = wenNativeConn
	wenDDcConn, err := grpc.Dial(cfg.GrpcClient.WenchangchainDDCAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		logger.Fatal("get wenchangchain-ddc grpc connect failed, err: ", err.Error())
	}
	GrpcConnMap[constant.WenchangDDC] = wenDDcConn

	//初始化Account grpc client
	AccountClientMap = make(map[string]pb_account.AccountClient)
	AccountClientMap[constant.WenchangDDC] = pb_account.NewAccountClient(GrpcConnMap[constant.WenchangDDC])
	AccountClientMap[constant.WenchangNative] = pb_account.NewAccountClient(GrpcConnMap[constant.WenchangNative])
	//初始化business grpc client
	BusineessClientMap = make(map[string]pb_business.BuyClient)
	BusineessClientMap[constant.WenchangDDC] = pb_business.NewBuyClient(GrpcConnMap[constant.WenchangDDC])
	BusineessClientMap[constant.WenchangNative] = pb_business.NewBuyClient(GrpcConnMap[constant.WenchangNative])
	//初始化msgs grpc client
	MsgsClientMap = make(map[string]pb_msgs.MSGSClient)
	MsgsClientMap[constant.WenchangDDC] = pb_msgs.NewMSGSClient(GrpcConnMap[constant.WenchangDDC])
	MsgsClientMap[constant.WenchangNative] = pb_msgs.NewMSGSClient(GrpcConnMap[constant.WenchangNative])
	//初始化nft grpc client
	NftClientMap = make(map[string]pb_nft.NFTClient)
	NftClientMap[constant.WenchangDDC] = pb_nft.NewNFTClient(GrpcConnMap[constant.WenchangDDC])
	NftClientMap[constant.WenchangNative] = pb_nft.NewNFTClient(GrpcConnMap[constant.WenchangNative])
	//初始化nft class grpc client
	ClassClientMap = make(map[string]pb_class.ClassClient)
	ClassClientMap[constant.WenchangDDC] = pb_class.NewClassClient(GrpcConnMap[constant.WenchangDDC])
	ClassClientMap[constant.WenchangNative] = pb_class.NewClassClient(GrpcConnMap[constant.WenchangNative])
	//初始化tx grpc client
	TxClientMap = make(map[string]pb_tx.TxClient)
	TxClientMap[constant.WenchangDDC] = pb_tx.NewTxClient(GrpcConnMap[constant.WenchangDDC])
	TxClientMap[constant.WenchangNative] = pb_tx.NewTxClient(GrpcConnMap[constant.WenchangNative])
}

func InitRedisClient(cfg *configs.Config, logger *log.Logger) {
	RedisClient = redis.NewRedisClient(cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB, logger)
}
