package initialize

import (
	log "github.com/sirupsen/logrus"
)

type myLoggerV2 struct {
	log *log.Entry
	v   int
}

func (c *myLoggerV2) V(l int) bool {
	return l < c.v
}

func NewMyLoggerV2(v int, logger *log.Logger) myLoggerV2 {
	return myLoggerV2{log: logger.WithField("service", "grpc"), v: v}
}

func (c *myLoggerV2) Info(args ...interface{}) {
	c.log.Info(args)
}

func (c *myLoggerV2) Warning(args ...interface{}) {
	c.log.Warning(args)
}

func (c *myLoggerV2) Error(args ...interface{}) {
	c.log.Error(args)
}

func (c *myLoggerV2) Fatal(args ...interface{}) {
	c.log.Fatal(args)
}

func (c *myLoggerV2) Infof(format string, args ...interface{}) {
	c.log.Infof(format, args)
}

func (c *myLoggerV2) Warningf(format string, args ...interface{}) {
	c.log.Warningf(format, args)
}

func (c *myLoggerV2) Errorf(format string, args ...interface{}) {
	c.log.Errorf(format, args)
}

func (c *myLoggerV2) Fatalf(format string, args ...interface{}) {
	c.log.Fatalf(format, args)
}

func (c *myLoggerV2) Infoln(args ...interface{}) {
	c.log.Infoln(args)
}

func (c *myLoggerV2) Warningln(args ...interface{}) {
	c.log.Warningln(args)
}

func (c *myLoggerV2) Errorln(args ...interface{}) {
	c.log.Errorln(args)
}

func (c *myLoggerV2) Fatalln(args ...interface{}) {
	c.log.Fatalln(args)
}
