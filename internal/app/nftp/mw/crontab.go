package mw

import (
	"github.com/robfig/cron/v3"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/metric"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/redis"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
)

func ProcessTimer() {
	crontab := cron.New(cron.WithSeconds())
	spec := "*/5 * * * * ?" //每五秒一次
	task := func() {
		if err := orm.GetDB().Ping(); err != nil {
			metric.NewPrometheus().ApiMysqlException.With().Set(-1) //监控mysql连接
		} else {
			metric.NewPrometheus().ApiMysqlException.With().Set(0) //监控mysql连接
		}

		if !redis.RedisPing() {
			metric.NewPrometheus().ApiRedisException.With().Set(-1) //监控redis连接
		} else {
			metric.NewPrometheus().ApiRedisException.With().Set(0) //监控redis连接
		}
	}
	crontab.AddFunc(spec, task)
	crontab.Start()
}
