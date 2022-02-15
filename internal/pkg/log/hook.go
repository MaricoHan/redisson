package log

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/metric"
)

type DefaultFieldHook struct {
}

func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	err, _ := entry.Data["error"]
	switch err {
	case sql.ErrConnDone:
		metric.NewPrometheus().ApiMysqlException.With([]string{}...).Set(-1)
		fmt.Println(metric.NewPrometheus().ApiMysqlException)
	case "redis connect":
		metric.NewPrometheus().ApiRedisException.With([]string{}...).Set(-1)
	}
	return nil
}

func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
