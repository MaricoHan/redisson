package configs

import log "github.com/sirupsen/logrus"

type Context struct {
	Logger *log.Logger
	Config *Config
}
