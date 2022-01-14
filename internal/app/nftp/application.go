package nftp

import (
	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/controller"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/mysql"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/redis"
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
	mysql.Connect(
		conf.Mysql.User,
		conf.Mysql.Password,
		conf.Mysql.Host,
		conf.Mysql.Port,
		conf.Mysql.DBName,
		mysql.DebugOption(conf.Server.LogLevel == log.DebugLevel),
		mysql.MaxIdleConnsOption(conf.Mysql.MaxIdleConns),
		mysql.MaxOpenConnsOption(conf.Mysql.MaxOpenConns),
		mysql.MaxLifetimeOption(conf.Mysql.MaxLifeTime),
		mysql.WriteOption(log.Log),
	)
}

func (s NFTPServer) Stop() {
	mysql.Close()
}
