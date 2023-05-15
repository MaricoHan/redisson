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
	pb_evm_class "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/class"
	pb_evm_contract "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/contract"
	pb_evm_dict "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/dict"

	pb_evm_msgs "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/msgs"
	pb_evm_nft "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/nft"
	pb_evm_ns "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/ns"
	pb_evm_tx "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/tx"
	pb_l2_dict "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/l2/dict"
	pb_l2_nft "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/l2/nft"
	pb_native_nft_class "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/class"
	pb_native_dict "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/dict"
	pb_native_msgs "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/msgs"
	pb_native_mt "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/mt"
	pb_native_mt_class "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/mt_class"
	pb_native_nft "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/nft"
	pb_native_notice "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/notice"
	pb_native_record "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/record"
	pb_native_tx "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/tx"
	pb_native_tx_queue "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/tx_queue"
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
var SignClient pb_account.AccountClient
var BusineessClientMap map[string]pb_business.BuyClient

var EvmMsgsClientMap map[string]pb_evm_msgs.MSGSClient
var EvmNftClientMap map[string]pb_evm_nft.NFTClient
var EvmClassClientMap map[string]pb_evm_class.ClassClient
var EvmTxClientMap map[string]pb_evm_tx.TxClient
var EvmNsClientMap map[string]pb_evm_ns.NSClient
var EvmContractClientMap map[string]pb_evm_contract.ContractClient
var EvmDictClientMap map[string]pb_evm_dict.DictClient

var NativeRecordClientMap map[string]pb_native_record.RecordClient
var NativeMTClientMap map[string]pb_native_mt.MTClient
var NativeMTClassClientMap map[string]pb_native_mt_class.MTClassClient
var NativeNFTClientMap map[string]pb_native_nft.NFTClient
var NativeNFTClassClientMap map[string]pb_native_nft_class.ClassClient
var NativeMsgClientMap map[string]pb_native_msgs.MSGSClient
var NativeTxClientMap map[string]pb_native_tx.TxClient
var NativeTxQueueClientMap map[string]pb_native_tx_queue.TxQueueClient
var NativeNoticeClientMap map[string]pb_native_notice.NoticeClient
var NativeDictClientMap map[string]pb_native_dict.DictClient

var StateGatewayServer *grpc.ClientConn

//var GrpcConnRightsMap map[string]*grpc.ClientConn
//var RightsClientMap map[string]rights.RightsClient

var WalletClientMap map[string]pb_wallet.WalletClient

var L2NftClientMap map[string]pb_l2_nft.NFTClient
var L2NftClassClientMap map[string]pb_l2_nft.ClassClient
var L2DictClientMap map[string]pb_l2_dict.DictClient

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
	tianzhouEvmConn, err := grpc.DialContext(
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
	GrpcConnMap[constant.IritaOPBNative] = tianzhouEvmConn

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

	logger.Info("connecting sign-server ...")
	signServer, err := grpc.DialContext(
		context.Background(),
		cfg.GrpcClient.SignServer,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(middleware.NewGrpcInterceptorMiddleware().Interceptor()))
	if err != nil {
		logger.Fatal("get sign-server grpc connect failed, err: ", err.Error())
	}

	// 初始化sign grpc client
	SignClient = pb_account.NewAccountClient(signServer)
	// 初始化msgs grpc client
	EvmMsgsClientMap = make(map[string]pb_evm_msgs.MSGSClient)
	EvmMsgsClientMap[constant.IritaOPBNative] = pb_evm_msgs.NewMSGSClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化nft grpc client
	EvmNftClientMap = make(map[string]pb_evm_nft.NFTClient)
	EvmNftClientMap[constant.IritaOPBNative] = pb_evm_nft.NewNFTClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化nft class grpc client
	EvmClassClientMap = make(map[string]pb_evm_class.ClassClient)
	EvmClassClientMap[constant.IritaOPBNative] = pb_evm_class.NewClassClient(GrpcConnMap[constant.IritaOPBNative])
	// 初始化tx grpc client
	EvmTxClientMap = make(map[string]pb_evm_tx.TxClient)
	EvmTxClientMap[constant.IritaOPBNative] = pb_evm_tx.NewTxClient(GrpcConnMap[constant.IritaOPBNative])

	NativeNFTClientMap = make(map[string]pb_native_nft.NFTClient)
	NativeNFTClientMap[constant.TianzhouNative] = pb_native_nft.NewNFTClient(GrpcConnMap[constant.TianzhouNative])
	NativeNFTClassClientMap = make(map[string]pb_native_nft_class.ClassClient)
	NativeNFTClassClientMap[constant.TianzhouNative] = pb_native_nft_class.NewClassClient(GrpcConnMap[constant.TianzhouNative])
	NativeTxClientMap = make(map[string]pb_native_tx.TxClient)
	NativeTxClientMap[constant.TianzhouNative] = pb_native_tx.NewTxClient(GrpcConnMap[constant.TianzhouNative])
	NativeTxQueueClientMap = make(map[string]pb_native_tx_queue.TxQueueClient)
	NativeTxQueueClientMap[constant.TianzhouNative] = pb_native_tx_queue.NewTxQueueClient(StateGatewayServer)
	NativeNoticeClientMap = make(map[string]pb_native_notice.NoticeClient)
	NativeNoticeClientMap[constant.TianzhouNative] = pb_native_notice.NewNoticeClient(GrpcConnMap[constant.TianzhouNative])
	// 初始化mt
	NativeMTClientMap = make(map[string]pb_native_mt.MTClient)
	NativeMTClientMap[constant.TianzhouNative] = pb_native_mt.NewMTClient(GrpcConnMap[constant.TianzhouNative])
	// 初始化mt_class
	NativeMTClassClientMap = make(map[string]pb_native_mt_class.MTClassClient)
	NativeMTClassClientMap[constant.TianzhouNative] = pb_native_mt_class.NewMTClassClient(GrpcConnMap[constant.TianzhouNative])
	// 初始化msgs
	NativeMsgClientMap = make(map[string]pb_native_msgs.MSGSClient)
	NativeMsgClientMap[constant.TianzhouNative] = pb_native_msgs.NewMSGSClient(GrpcConnMap[constant.TianzhouNative])

	// 初始化tx_queue
	// TxQueueClient = pb_tx_queue.NewTxQueueClient(StateGatewayServer)

	// 初始化rights_jiangsu
	//RightsClientMap = make(map[string]rights.RightsClient)
	//RightsClientMap[constant.JiangSu] = rights.NewRightsClient(GrpcConnRightsMap[constant.JiangSu])

	//初始化record grpc client
	NativeRecordClientMap = make(map[string]pb_native_record.RecordClient)
	//NativeRecordClientMap[constant.WenchangDDC] = pb_native_record.NewRecordClient(GrpcConnMap[constant.WenchangDDC])
	//NativeRecordClientMap[constant.WenchangNative] = pb_native_record.NewRecordClient(GrpcConnMap[constant.WenchangNative])
	NativeRecordClientMap[constant.IritaOPBNative] = pb_native_record.NewRecordClient(GrpcConnMap[constant.IritaOPBNative])
	//NativeRecordClientMap[constant.IrisHubNative] = pb_native_record.NewRecordClient(GrpcConnMap[constant.IrisHubNative])

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
	EvmNsClientMap = make(map[string]pb_evm_ns.NSClient)
	EvmNsClientMap[constant.IritaOPBNative] = pb_evm_ns.NewNSClient(GrpcConnMap[constant.IritaOPBNative])

	// 初始化contract grpc client
	EvmContractClientMap = make(map[string]pb_evm_contract.ContractClient)
	EvmContractClientMap[constant.IritaOPBNative] = pb_evm_contract.NewContractClient(GrpcConnMap[constant.IritaOPBNative])

	// 初始化 l2 NftClass grpc client
	L2NftClassClientMap = make(map[string]pb_l2_nft.ClassClient)
	L2NftClassClientMap[constant.IritaOPBNative] = pb_l2_nft.NewClassClient(GrpcConnMap[constant.IritaLayer2])

	// 初始化 l2 nft grpc client
	L2NftClientMap = make(map[string]pb_l2_nft.NFTClient)
	L2NftClientMap[constant.IritaOPBNative] = pb_l2_nft.NewNFTClient(GrpcConnMap[constant.IritaLayer2])

}

func InitRedisClient(cfg *configs.Config, logger *log.Logger) {
	logger.Info("connecting redis ...")
	RedisClient = redis.NewRedisClient(cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB, logger)
}
