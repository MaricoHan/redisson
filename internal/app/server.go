package app

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
)

//Server define a http Server
type Server struct {
	svr         *http.Server
	router      *mux.Router
	middlewares []Middleware
}

//Start a instance of the http server
func Start() {
	app := nftp.NewNFTPServer()
	app.Initialize()

	s := NewServer()
	for _, rt := range app.GetEndpoints() {
		s.RegisterEndpoint(rt)
	}

	lis, err := net.Listen("tcp", config.Get().Server.Address)
	if err != nil {
		panic(err)
	}

	//启动http服务
	go func() {
		_ = s.svr.Serve(lis)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	//释放所有资源
	app.Stop()
	log.Info("Stop the dapp server", "sig", int(sig.(syscall.Signal))+128)
}

func NewServer() Server {
	middlewares := []Middleware{
		//should be last one
		RecoverMiddleware,
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
	s.router.Handle(end.URI, h).Methods(end.Method)
}
