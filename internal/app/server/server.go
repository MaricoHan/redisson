package server

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/app/controller"
	"gitlab.bianjie.ai/avata/open-api/internal/app/middlewares"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
)

//Server define a http Server
type Server struct {
	Svr         *http.Server
	router      *mux.Router
	middlewares []func(http.Handler) http.Handler
}

func NewServer(ctx *configs.Context) Server {
	middlewares := []func(http.Handler) http.Handler{
		//should be last one
		middlewares.RecoverMiddleware,
		middlewares.IdempotentMiddleware,
		middlewares.AuthMiddleware,
	}

	r := mux.NewRouter()
	r = r.PathPrefix(fmt.Sprintf("/%s", ctx.Config.App.RouterPrefix)).Subrouter()
	svr := http.Server{Handler: r}

	return Server{
		Svr:         &svr,
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

type Application struct {
	logger *log.Logger
}

func NewApplication(logger *log.Logger) kit.Application {
	return &Application{logger: logger}
}

func (s *Application) GetEndpoints() []kit.Endpoint {
	var rs []kit.Endpoint

	ctls := controller.GetAllControllers(s.logger)
	for _, c := range ctls {
		rs = append(rs, c.GetEndpoints()...)
	}
	return rs
}
