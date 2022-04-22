package app

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/app/middlewares"
	"gitlab.bianjie.ai/avata/open-api/internal/app/server"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/metric"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/safe"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//Start a instance of the http server
func Start(ctx *configs.Context) {

	if ctx.Config.App.GprcTimeout == 0 {
		constant.GrpcTimeout = 10
	} else {
		constant.GrpcTimeout = ctx.Config.App.GprcTimeout
	}

	//初始化logger
	ctx.Logger = initialize.Logger(ctx.Config)
	//初始化redis
	initialize.InitRedisClient(ctx.Config, ctx.Logger)
	//初始化mysql
	initialize.InitMysqlDB(ctx.Config, ctx.Logger)
	//初始化etcd解析器
	initialize.InitEtcdResolver(ctx.Config, ctx.Logger)

	app := server.NewApplication(ctx.Logger)
	s := server.NewServer(ctx)
	for _, rt := range app.GetEndpoints() {
		s.RegisterEndpoint(rt)
	}

	// metric
	metric.NewPrometheus().InitPrometheus()
	//crontab
	middlewares.ProcessTimer()

	lis, err := net.Listen("tcp", ctx.Config.App.Addr)
	if err != nil {
		panic(err)
	}

	//启动http服务
	log.WithFields(log.Fields{
		"transport": "http",
		"http_addr": ctx.Config.App.Addr,
	}).Info("start")
	go func() {
		_ = s.Svr.Serve(lis)
	}()

	safe.GoSafe(func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(ctx.Config.App.PrometheusAddr, nil)
	}, func(err error) {
		log.Info("http server listenidg error: %s", err)
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	//释放所有资源
	initialize.GrpcConnMap[constant.WenchangNative].Close()
	initialize.GrpcConnMap[constant.WenchangDDC].Close()
	initialize.RedisClient.Close()
	log.WithFields(log.Fields{
		"transport": "http",
		"http_addr": ctx.Config.App.Addr,
	}).Info("stop the openapi server, sig:", int(sig.(syscall.Signal))+128)
}
