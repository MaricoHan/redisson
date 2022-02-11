package log

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/metric"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/redis"
)

type DefaultFieldHook struct {
}

func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	err, ok := entry.Data["error"]
	if !ok {
		if metric.NewPrometheus().ApiMysqlException != nil {
			metric.NewPrometheus().ApiMysqlException.With([]string{}...).Set(0)
			metric.NewPrometheus().ApiRedisException.With([]string{}...).Set(0)
		}
		return nil
	}
	switch err {
	case sql.ErrConnDone:
		metric.NewPrometheus().ApiMysqlException.With([]string{}...).Set(-1)
		fmt.Println(metric.NewPrometheus().ApiMysqlException)
	case redis.ErrNotObtained:
		metric.NewPrometheus().ApiRedisException.With([]string{}...).Set(-1)
	default:
		metric.NewPrometheus().ApiMysqlException.With([]string{}...).Set(0)
		metric.NewPrometheus().ApiRedisException.With([]string{}...).Set(0)
	}
	return nil
}

func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
