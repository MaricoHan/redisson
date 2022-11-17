package initialize

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb_account "gitlab.bianjie.ai/avata/chains/api/pb/account"
	pb_business "gitlab.bianjie.ai/avata/chains/api/pb/buy"
	pb_class "gitlab.bianjie.ai/avata/chains/api/pb/class"
	pb_msgs "gitlab.bianjie.ai/avata/chains/api/pb/msgs"
	pb_mt "gitlab.bianjie.ai/avata/chains/api/pb/mt"
	pb_mt_class "gitlab.bianjie.ai/avata/chains/api/pb/mt_class"
	pb_mt_msgs "gitlab.bianjie.ai/avata/chains/api/pb/mt_msgs"
	pb_nft "gitlab.bianjie.ai/avata/chains/api/pb/nft"
	pb_tx "gitlab.bianjie.ai/avata/chains/api/pb/tx"
	pb_tx_queue "gitlab.bianjie.ai/avata/chains/api/pb/tx_queue"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/redis"
	"gitlab.bianjie.ai/avata/open-api/pkg/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
var MTClientMap map[string]pb_mt.MTClient
var MTClassClientMap map[string]pb_mt_class.MTClassClient
var MTMsgsClientMap map[string]pb_mt_msgs.MTMSGSClient

var StateGatewayServer *grpc.ClientConn
var TxQueueClient pb_tx_queue.TxQueueClient

var Log = new(log.Logger)

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
	Log = log.StandardLogger()
	return Log
}

func InitMysqlDB(cfg *configs.Config, logger *log.Logger) {
	logger.Info("connecting mysql ...")
	gormLogger := logs.NewGormLogger(logger)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Mysql.Username,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.DB)

	mysqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormLogger, NamingStrategy: schema.NamingStrategy{
		TablePrefix:   "t_",  // 表前缀
		SingularTable: false, // 复数形式
	}})
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
	logger.Info("connecting wenchangchain-native ...")
	wenNativeConn, err := grpc.DialContext(context.Background(), cfg.GrpcClient.WenchangchainNativeAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp), grpc.WithBlock(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		logger.Fatal("get wenchangchain-native grpc connect failed, err: ", err.Error())
	}
	GrpcConnMap[constant.WenchangNative] = wenNativeConn

	logger.Info("connecting wenchangchain-ddc ...")
	wenDDcConn, err := grpc.DialContext(context.Background(), cfg.GrpcClient.WenchangchainDDCAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp), grpc.WithBlock(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		logger.Fatal("get wenchangchain-ddc grpc connect failed, err: ", err.Error())
	}
	GrpcConnMap[constant.WenchangDDC] = wenDDcConn

	logger.Info("connecting irita-opb-native ...")
	IritaOPBNativeConn, err := grpc.DialContext(context.Background(), cfg.GrpcClient.IritaOPBNativeAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp), grpc.WithBlock(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		logger.Fatal("get irita-opb-native grpc connect failed, err: ", err.Error())
	}
	GrpcConnMap[constant.IritaOPBNative] = IritaOPBNativeConn

	logger.Info("connecting state-gateway-server ...")
	StateGatewayServer, err = grpc.DialContext(context.Background(), cfg.GrpcClient.StateGatewayAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp), grpc.WithBlock(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		logger.Fatal("get state-gateway-server grpc connect failed, err: ", err.Error())
	}

	// 初始化Account grpc client
	AccountClientMap = make(map[string]pb_account.AccountClient)
	AccountClientMap[constant.WenchangDDC] = pb_account.NewAccountClient(GrpcConnMap[constant.WenchangDDC])
	AccountClientMap[constant.WenchangNative] = pb_account.NewAccountClient(GrpcConnMap[constant.WenchangNative])
	AccountClientMap[constant.IritaOPBNative] = pb_account.NewAccountClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化business grpc client
	BusineessClientMap = make(map[string]pb_business.BuyClient)
	BusineessClientMap[constant.WenchangDDC] = pb_business.NewBuyClient(GrpcConnMap[constant.WenchangDDC])
	BusineessClientMap[constant.WenchangNative] = pb_business.NewBuyClient(GrpcConnMap[constant.WenchangNative])
	BusineessClientMap[constant.IritaOPBNative] = pb_business.NewBuyClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化msgs grpc client
	MsgsClientMap = make(map[string]pb_msgs.MSGSClient)
	MsgsClientMap[constant.WenchangDDC] = pb_msgs.NewMSGSClient(GrpcConnMap[constant.WenchangDDC])
	MsgsClientMap[constant.WenchangNative] = pb_msgs.NewMSGSClient(GrpcConnMap[constant.WenchangNative])
	MsgsClientMap[constant.IritaOPBNative] = pb_msgs.NewMSGSClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化nft grpc client
	NftClientMap = make(map[string]pb_nft.NFTClient)
	NftClientMap[constant.WenchangDDC] = pb_nft.NewNFTClient(GrpcConnMap[constant.WenchangDDC])
	NftClientMap[constant.WenchangNative] = pb_nft.NewNFTClient(GrpcConnMap[constant.WenchangNative])
	NftClientMap[constant.IritaOPBNative] = pb_nft.NewNFTClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化nft class grpc client
	ClassClientMap = make(map[string]pb_class.ClassClient)
	ClassClientMap[constant.WenchangDDC] = pb_class.NewClassClient(GrpcConnMap[constant.WenchangDDC])
	ClassClientMap[constant.WenchangNative] = pb_class.NewClassClient(GrpcConnMap[constant.WenchangNative])
	ClassClientMap[constant.IritaOPBNative] = pb_class.NewClassClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化tx grpc client
	TxClientMap = make(map[string]pb_tx.TxClient)
	TxClientMap[constant.WenchangDDC] = pb_tx.NewTxClient(GrpcConnMap[constant.WenchangDDC])
	TxClientMap[constant.WenchangNative] = pb_tx.NewTxClient(GrpcConnMap[constant.WenchangNative])
	TxClientMap[constant.IritaOPBNative] = pb_tx.NewTxClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化mt
	MTClientMap = make(map[string]pb_mt.MTClient)
	MTClientMap[constant.WenchangDDC] = pb_mt.NewMTClient(GrpcConnMap[constant.WenchangDDC])
	MTClientMap[constant.WenchangNative] = pb_mt.NewMTClient(GrpcConnMap[constant.WenchangNative])
	MTClientMap[constant.IritaOPBNative] = pb_mt.NewMTClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化mt_class
	MTClassClientMap = make(map[string]pb_mt_class.MTClassClient)
	MTClassClientMap[constant.WenchangDDC] = pb_mt_class.NewMTClassClient(GrpcConnMap[constant.WenchangDDC])
	MTClassClientMap[constant.WenchangNative] = pb_mt_class.NewMTClassClient(GrpcConnMap[constant.WenchangNative])
	MTClassClientMap[constant.IritaOPBNative] = pb_mt_class.NewMTClassClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化mt_msgs
	MTMsgsClientMap = make(map[string]pb_mt_msgs.MTMSGSClient)
	MTMsgsClientMap[constant.WenchangNative] = pb_mt_msgs.NewMTMSGSClient(GrpcConnMap[constant.WenchangNative])
	MTMsgsClientMap[constant.IritaOPBNative] = pb_mt_msgs.NewMTMSGSClient(GrpcConnMap[constant.IritaOPBNative])

	// 初始化tx_queue
	TxQueueClient = pb_tx_queue.NewTxQueueClient(StateGatewayServer)
}

func InitRedisClient(cfg *configs.Config, logger *log.Logger) {
	logger.Info("connecting redis ...")
	RedisClient = redis.NewRedisClient(cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB, logger)
}
