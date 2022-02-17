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
	spec := "*/10 * * * * ?" //每十秒一次
	task := func() {
		metric.NewPrometheus().SyncTxPendingTotal.With().Set(0)
		metric.NewPrometheus().SyncTxFailedTotal.With().Set(0)
		metric.NewPrometheus().SyncNftLockedTotal.With().Set(0)
		metric.NewPrometheus().SyncNftTotal.With().Set(0)
		metric.NewPrometheus().SyncClassLockedTotal.With().Set(0)
		metric.NewPrometheus().SyncClassTotal.With().Set(0)
		metric.NewPrometheus().SyncMysqlException.With().Set(0)
		metric.NewPrometheus().SyncRedisException.With().Set(0)
		metric.NewPrometheus().SyncChainConnError.With().Set(0)

		txsPending, _ := models.TTXS(
			models.TTXWhere.Status.EQ(models.TTXSStatusPending),
		).AllG(context.Background())
		for _ = range txsPending {
			metric.NewPrometheus().SyncTxPendingTotal.With().Add(1) //系统未完成的交易总量
		}

		txsFailed, _ := models.TTXS(
			models.TTXWhere.Status.EQ(models.TTXSStatusFailed),
		).AllG(context.Background())
		for _ = range txsFailed {
			metric.NewPrometheus().SyncTxFailedTotal.With().Add(1) //系统失败的交易总量
		}

		nftLocked, _ := models.TNFTS(
			models.TNFTWhere.LockedBy.NEQ(null.Uint64FromPtr(nil)),
		).AllG(context.Background())
		for _ = range nftLocked {
			metric.NewPrometheus().SyncNftLockedTotal.With().Add(1) //系统锁定的nft总量
		}

		nft, _ := models.TNFTS().AllG(context.Background())
		for _ = range nft {
			metric.NewPrometheus().SyncNftTotal.With().Add(1) //系统创建的nft总量
		}

		nftClassLocked, _ := models.TClasses(
			models.TClassWhere.LockedBy.NEQ(null.Uint64FromPtr(nil)),
		).AllG(context.Background())
		for _ = range nftClassLocked {
			metric.NewPrometheus().SyncClassLockedTotal.With().Add(1) //系统锁定的class总量
		}

		nftClass, _ := models.TClasses().AllG(context.Background())
		for _ = range nftClass {
			metric.NewPrometheus().SyncClassTotal.With().Add(1) //系统创建的class总量
		}

		if err := orm.GetDB().Ping(); err != nil {
			metric.NewPrometheus().SyncMysqlException.With().Set(-1) //监控mysql连接
		}

		if !redis.RedisPing() {
			metric.NewPrometheus().SyncRedisException.With().Set(-1) //监控redis连接
		}

		if _, err := chain.GetSdkClient().GenConn(); err != nil {
			metric.NewPrometheus().SyncChainConnError.With().Set(-1) //访问链异常
		}
	}
	crontab.AddFunc(spec, task)
	crontab.Start()
}
