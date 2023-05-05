package initialize

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	pb_account "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/account"
	pb_business "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/buy"
	pb_class "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/class"
	pb_msgs "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/msgs"
	pb_nft "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/nft"
	pb_record "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/record"
	pb_tx "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/tx"
	//pb_notice "gitlab.bianjie.ai/avata/chains/api/pb/v2/notice"
	pb_contract "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/contract"
	pb_l2_nft "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/l2/nft"
	pb_ns "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/ns"
	//pb_tx_queue "gitlab.bianjie.ai/avata/chains/api/pb/v2/tx_queue"
	pb_wallet "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/wallet"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/middleware"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/redis"
	"gitlab.bianjie.ai/avata/open-api/pkg/logs"
	trace_log "gitlab.bianjie.ai/avata/utils/commons/trace/log"
)

var RedisClient *redis.RedisClient
var MysqlDB *gorm.DB
var GrpcConnMap map[string]*grpc.ClientConn
var AccountClientMap map[string]pb_account.AccountClient
var MsgsClientMap map[string]pb_msgs.MSGSClient
var NftClientMap map[string]pb_nft.NFTClient

var RecordClientMap map[string]pb_record.RecordClient
var ClassClientMap map[string]pb_class.ClassClient
var TxClientMap map[string]pb_tx.TxClient
var BusineessClientMap map[string]pb_business.BuyClient

//var MTClientMap map[string]pb_mt.MTClient
//var MTClassClientMap map[string]pb_mt_class.MTClassClient
//var MTMsgsClientMap map[string]pb_mt_msgs.MTMSGSClient

var StateGatewayServer *grpc.ClientConn

//var GrpcConnRightsMap map[string]*grpc.ClientConn
//var RightsClientMap map[string]rights.RightsClient

var WalletClientMap map[string]pb_wallet.WalletClient

var NsClientMap map[string]pb_ns.NSClient

var ContractClientMap map[string]pb_contract.ContractClient

var L2NftClientMap map[string]pb_l2_nft.NFTClient

var L2NftClassClientMap map[string]pb_l2_nft.ClassClient

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
	Log.AddHook(&trace_log.TraceHook{})
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

	logger.Info("connecting tianzhou-evm ...")
	iritaOpbNativeConn, err := grpc.DialContext(
		context.Background(),
		cfg.GrpcClient.TianZhouEVM,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(middleware.NewGrpcInterceptorMiddleware().Interceptor()))
	if err != nil {
		logger.Fatal("get tianzhou-evm grpc connect failed, err: ", err.Error())
	}
	GrpcConnMap[constant.IritaOPBNative] = iritaOpbNativeConn

	logger.Info("connecting state-gateway-server ...")
	StateGatewayServer, err = grpc.DialContext(
		context.Background(),
		cfg.GrpcClient.StateGateway,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(middleware.NewGrpcInterceptorMiddleware().Interceptor()))
	if err != nil {
		logger.Fatal("get state-gateway-server grpc connect failed, err: ", err.Error())
	}

	logger.Info("connecting wallet-server ...")
	walletServer, err := grpc.DialContext(
		context.Background(),
		cfg.GrpcClient.WalletServer,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(middleware.NewGrpcInterceptorMiddleware().Interceptor()))
	if err != nil {
		logger.Fatal("get wallet-server grpc connect failed, err: ", err.Error())
	}
	GrpcConnMap[constant.WalletServer] = walletServer

	logger.Info("connecting irita-layer2 ...")
	iritaLayer2, err := grpc.DialContext(
		context.Background(),
		cfg.GrpcClient.IritaLayer2,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(middleware.NewGrpcInterceptorMiddleware().Interceptor()))
	if err != nil {
		logger.Fatal("get irita-layer2 grpc connect failed, err: ", err.Error())
	}
	GrpcConnMap[constant.IritaLayer2] = iritaLayer2

	// 初始化Account grpc client
	AccountClientMap = make(map[string]pb_account.AccountClient)
	AccountClientMap[constant.IritaOPBNative] = pb_account.NewAccountClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化msgs grpc client
	MsgsClientMap = make(map[string]pb_msgs.MSGSClient)
	MsgsClientMap[constant.IritaOPBNative] = pb_msgs.NewMSGSClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化nft grpc client
	NftClientMap = make(map[string]pb_nft.NFTClient)
	NftClientMap[constant.IritaOPBNative] = pb_nft.NewNFTClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化nft class grpc client
	ClassClientMap = make(map[string]pb_class.ClassClient)
	ClassClientMap[constant.IritaOPBNative] = pb_class.NewClassClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化tx grpc client
	TxClientMap = make(map[string]pb_tx.TxClient)
	TxClientMap[constant.IritaOPBNative] = pb_tx.NewTxClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化mt
	//MTClientMap = make(map[string]pb_mt.MTClient)
	//MTClientMap[constant.WenchangDDC] = pb_mt.NewMTClient(GrpcConnMap[constant.WenchangDDC])
	//MTClientMap[constant.WenchangNative] = pb_mt.NewMTClient(GrpcConnMap[constant.WenchangNative])
	//MTClientMap[constant.IritaOPBNative] = pb_mt.NewMTClient(GrpcConnMap[constant.IritaOPBNative])
	//MTClientMap[constant.IrisHubNative] = pb_mt.NewMTClient(GrpcConnMap[constant.IrisHubNative])
	// 初始化mt_class
	//MTClassClientMap = make(map[string]pb_mt_class.MTClassClient)
	//MTClassClientMap[constant.WenchangDDC] = pb_mt_class.NewMTClassClient(GrpcConnMap[constant.WenchangDDC])
	//MTClassClientMap[constant.WenchangNative] = pb_mt_class.NewMTClassClient(GrpcConnMap[constant.WenchangNative])
	//MTClassClientMap[constant.IritaOPBNative] = pb_mt_class.NewMTClassClient(GrpcConnMap[constant.IritaOPBNative])
	//MTClassClientMap[constant.IrisHubNative] = pb_mt_class.NewMTClassClient(GrpcConnMap[constant.IrisHubNative])
	// 初始化mt_msgs
	//MTMsgsClientMap = make(map[string]pb_mt_msgs.MTMSGSClient)
	//MTMsgsClientMap[constant.WenchangNative] = pb_mt_msgs.NewMTMSGSClient(GrpcConnMap[constant.WenchangNative])
	//MTMsgsClientMap[constant.IritaOPBNative] = pb_mt_msgs.NewMTMSGSClient(GrpcConnMap[constant.IritaOPBNative])
	//MTMsgsClientMap[constant.IrisHubNative] = pb_mt_msgs.NewMTMSGSClient(GrpcConnMap[constant.IrisHubNative])

	// 初始化tx_queue
	// TxQueueClient = pb_tx_queue.NewTxQueueClient(StateGatewayServer)

	// 初始化rights_jiangsu
	//RightsClientMap = make(map[string]rights.RightsClient)
	//RightsClientMap[constant.JiangSu] = rights.NewRightsClient(GrpcConnRightsMap[constant.JiangSu])

	//初始化record grpc client
	RecordClientMap = make(map[string]pb_record.RecordClient)
	//RecordClientMap[constant.WenchangDDC] = pb_record.NewRecordClient(GrpcConnMap[constant.WenchangDDC])
	//RecordClientMap[constant.WenchangNative] = pb_record.NewRecordClient(GrpcConnMap[constant.WenchangNative])
	RecordClientMap[constant.IritaOPBNative] = pb_record.NewRecordClient(GrpcConnMap[constant.IritaOPBNative])
	//RecordClientMap[constant.IrisHubNative] = pb_record.NewRecordClient(GrpcConnMap[constant.IrisHubNative])

	// 初始化notice
	//NoticeClientMap = make(map[string]pb_notice.NoticeClient)
	//NoticeClientMap[constant.WenchangNative] = pb_notice.NewNoticeClient(GrpcConnMap[constant.WenchangNative])
	//NoticeClientMap[constant.IritaOPBNative] = pb_notice.NewNoticeClient(GrpcConnMap[constant.IritaOPBNative])
	//NoticeClientMap[constant.WenchangDDC] = pb_notice.NewNoticeClient(GrpcConnMap[constant.WenchangDDC])
	//NoticeClientMap[constant.IrisHubNative] = pb_notice.NewNoticeClient(GrpcConnMap[constant.IrisHubNative])

	// 初始化business grpc client
	BusineessClientMap = make(map[string]pb_business.BuyClient)
	BusineessClientMap[constant.IritaOPBNative] = pb_business.NewBuyClient(GrpcConnMap[constant.IritaOPBNative])

	// 初始化wallet grpc client
	WalletClientMap = make(map[string]pb_wallet.WalletClient)
	WalletClientMap[constant.WalletServer] = pb_wallet.NewWalletClient(GrpcConnMap[constant.WalletServer])

	// 初始化ns grpc client
	NsClientMap = make(map[string]pb_ns.NSClient)
	NsClientMap[constant.IritaOPBNative] = pb_ns.NewNSClient(GrpcConnMap[constant.IritaOPBNative])

	// 初始化contract grpc client
	ContractClientMap = make(map[string]pb_contract.ContractClient)
	ContractClientMap[constant.IritaOPBNative] = pb_contract.NewContractClient(GrpcConnMap[constant.IritaOPBNative])

	// 初始化 l2 NftClass grpc client
	L2NftClassClientMap = make(map[string]pb_l2_nft.ClassClient)
	L2NftClassClientMap[constant.IritaOPBNative] = pb_l2_nft.NewClassClient(GrpcConnMap[constant.IritaLayer2])

}

func InitRedisClient(cfg *configs.Config, logger *log.Logger) {
	logger.Info("connecting redis ...")
	RedisClient = redis.NewRedisClient(cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB, logger)
}
