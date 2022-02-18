package nftp

import (
	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/controller"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/chain"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/redis"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
)

type NFTPServer struct {
}

func NewNFTPServer() kit.Application {
	return NFTPServer{}
}

//GetEndpoints return all the endpoints for http server
func (s NFTPServer) GetEndpoints() []kit.Endpoint {
	var rs []kit.Endpoint

	ctls := controller.GetAllControllers()
	for _, c := range ctls {
		rs = append(rs, c.GetEndpoints()...)
	}
	return rs
}

func (s NFTPServer) Initialize() {
	conf := config.Get()
	//初始化logger连接
	log.SetLogger(conf.Server.LogFormat, conf.Server.LogLevel)
	//初始化redis连接
	redis.Connect(conf.Redis.Address, conf.Redis.Password, conf.Redis.DB)
	//初始化mysql连接
	orm.Connect(
		conf.Mysql.User,
		conf.Mysql.Password,
		conf.Mysql.Host,
		conf.Mysql.Port,
		conf.Mysql.DBName,
		orm.DebugOption(conf.Server.LogLevel == log.DebugLevel),
		orm.MaxIdleConnsOption(conf.Mysql.MaxIdleConns),
		orm.MaxOpenConnsOption(conf.Mysql.MaxOpenConns),
		orm.MaxLifetimeOption(conf.Mysql.MaxLifeTime),
		orm.WriteOption(log.Log),
	)
	// 链客户端初始化
	chain.NewSdkClient(conf.Server.Env, conf.Chain, orm.GetDB())
}

func (s NFTPServer) Stop() {
	orm.CloseDB()
}
