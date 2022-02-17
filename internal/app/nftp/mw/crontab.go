package mw

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/volatiletech/null/v8"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/chain"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/metric"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/redis"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

func ProcessTimer() {
	crontab := cron.New(cron.WithSeconds())
	spec := "*/5 * * * * ?" //每五秒一次
	task := func() {
		txsPending, _ := models.TTXS(
			models.TTXWhere.Status.EQ(models.TTXSStatusPending),
		).CountG(context.Background())
		metric.NewPrometheus().SyncTxPendingTotal.With().Set(float64(txsPending)) //系统未完成的交易总量

		txsFailed, _ := models.TTXS(
			models.TTXWhere.Status.EQ(models.TTXSStatusFailed),
		).CountG(context.Background())
		metric.NewPrometheus().SyncTxFailedTotal.With().Set(float64(txsFailed)) //系统失败的交易总量

		nftLocked, _ := models.TNFTS(
			models.TNFTWhere.LockedBy.NEQ(null.Uint64FromPtr(nil)),
		).CountG(context.Background())
		metric.NewPrometheus().SyncNftLockedTotal.With().Set(float64(nftLocked)) //系统锁定的nft总量

		nft, _ := models.TNFTS().CountG(context.Background())
		metric.NewPrometheus().SyncNftTotal.With().Set(float64(nft)) //系统创建的nft总量

		nftClassLocked, _ := models.TClasses(
			models.TClassWhere.LockedBy.NEQ(null.Uint64FromPtr(nil)),
		).CountG(context.Background())
		metric.NewPrometheus().SyncClassLockedTotal.With().Set(float64(nftClassLocked)) //系统锁定的class总量

		nftClass, _ := models.TClasses().CountG(context.Background())
		metric.NewPrometheus().SyncClassTotal.With().Set(float64(nftClass)) //系统创建的class总量

		if err := orm.GetDB().Ping(); err != nil {
			metric.NewPrometheus().SyncMysqlException.With().Set(-1) //监控mysql连接
		} else {
			metric.NewPrometheus().SyncMysqlException.With().Set(0) //监控mysql连接
		}

		if !redis.RedisPing() {
			metric.NewPrometheus().SyncRedisException.With().Set(-1) //监控redis连接
		} else {
			metric.NewPrometheus().SyncRedisException.With().Set(0) //监控redis连接
		}

		if _, err := chain.GetSdkClient().GenConn(); err != nil {
			metric.NewPrometheus().SyncChainConnError.With().Set(-1) //访问链异常
		} else {
			metric.NewPrometheus().SyncChainConnError.With().Set(0) //访问链异常
		}
	}
	crontab.AddFunc(spec, task)
	crontab.Start()
}
