package middlewares

import (
	cron "github.com/robfig/cron/v3"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/metric"
)

func ProcessTimer() {
	crontab := cron.New(cron.WithSeconds())
	spec := "*/5 * * * * ?" //每五秒一次

	task := func() {
		if !initialize.RedisClient.Ping() {
			metric.NewPrometheus().ApiRedisException.With().Set(-1) //监控redis连接
		} else {
			metric.NewPrometheus().ApiRedisException.With().Set(0) //监控redis连接
		}
	}
	crontab.AddFunc(spec, task)
	crontab.Start()
}
