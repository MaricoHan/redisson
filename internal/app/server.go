package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/mw"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/helper"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/metric"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//Server define a http Server
type Server struct {
	svr         *http.Server
	router      *mux.Router
	middlewares []mw.Middleware
}

//Start a instance of the http server
func Start() {
	app := nftp.NewNFTPServer()
	log.Info("Initialize nftp server ")
	app.Initialize()
	s := NewServer()
	for _, rt := range app.GetEndpoints() {
		s.RegisterEndpoint(rt)
	}

	// metric
	metric.NewPrometheus().InitPrometheus()
	//time
	mw.ProcessTimer()

	lis, err := net.Listen("tcp", config.Get().Server.Address)
	if err != nil {
		panic(err)
	}
	log.Info("Start nftp server")
	//启动http服务
	go func() {
		_ = s.svr.Serve(lis)
	}()

	helper.GoSafe(func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(config.Get().Server.PrometheusAddr, nil)
	}, func(err error) {
		log.Info("http server listenidg error: %s", err)
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	//释放所有资源
	app.Stop()
	log.Info("Stop the nftp server", "sig", int(sig.(syscall.Signal))+128)
}

func NewServer() Server {
	middlewares := []mw.Middleware{
		//should be last one
		mw.RecoverMiddleware,
		mw.IdempotentMiddleware,
		mw.AuthMiddleware,
	}

	r := mux.NewRouter()
	svr := http.Server{Handler: r}

	return Server{
		svr:         &svr,
		router:      r,
		middlewares: middlewares,
	}
}

//RegisterEndpoint registers the handler for the given pattern.
func (s *Server) RegisterEndpoint(end kit.Endpoint) {
	var h = end.Handler
	for _, m := range s.middlewares {
		h = m(h)
	}
	s.router.Handle(fmt.Sprintf("/v1beta1%s", end.URI), h).Methods(end.Method)
}
