package mw

import (
	"context"
	"database/sql"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"github.com/volatiletech/null/v8"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/chain"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/metric"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

func ProcessTimer() {
	crontab := cron.New(cron.WithSeconds())
	spec := "*/10 * * * * ?" //每十秒一次
	task := func() {
		txsPending, err := models.TTXS(
			models.TTXWhere.Status.EQ(models.TTXSStatusPending),
		).AllG(context.Background())
		if err != nil && errors.Cause(err) == sql.ErrNoRows {
			//404
		} else if err != nil {
			//500
			log.Error("query tx pending total", "query tx error:", err.Error())
		}
		for _, tx := range txsPending {
			metric.NewPrometheus().SyncTxPendingTotal.With().Add(1) //系统未完成的交易总量
			interval := time.Now().Sub(tx.CreateAt.Time)
			metric.NewPrometheus().SyncTxPendingSeconds.With([]string{"tx_hash", tx.Hash, "interval", interval.String()}...).Set(-1) //系统未完成的交易
		}

		txsFailed, err := models.TTXS(
			models.TTXWhere.Status.EQ(models.TTXSStatusFailed),
		).AllG(context.Background())
		if err != nil && errors.Cause(err) == sql.ErrNoRows {
			//404
		} else if err != nil {
			//500
			log.Error("query tx failed total", "query tx error:", err.Error())
		}
		for _ = range txsFailed {
			metric.NewPrometheus().SyncTxFailedTotal.With().Add(1) //系统失败的交易总量
		}

		nftLocked, err := models.TNFTS(
			models.TNFTWhere.LockedBy.NEQ(null.Uint64FromPtr(nil)),
		).AllG(context.Background())
		if err != nil && errors.Cause(err) == sql.ErrNoRows {
			//404
		} else if err != nil {
			//500
			log.Error("query nft locked total", "query nft error:", err.Error())
		}
		for _ = range nftLocked {
			metric.NewPrometheus().SyncNftLockedTotal.With().Add(1) //系统锁定的nft总量
		}

		nft, err := models.TNFTS().AllG(context.Background())
		if err != nil && errors.Cause(err) == sql.ErrNoRows {
			//404
		} else if err != nil {
			//500
			log.Error("query nft total", "query nft error:", err.Error())
		}
		for _ = range nft {
			metric.NewPrometheus().SyncNftTotal.With().Add(1) //系统创建的nft总量
		}

		nftClassLocked, err := models.TClasses(
			models.TClassWhere.LockedBy.NEQ(null.Uint64FromPtr(nil)),
		).AllG(context.Background())
		if err != nil && errors.Cause(err) == sql.ErrNoRows {
			//404
		} else if err != nil {
			//500
			log.Error("query nft class locked total", "query nft class error:", err.Error())
		}
		for _ = range nftClassLocked {
			metric.NewPrometheus().SyncClassLockedTotal.With().Add(1) //系统锁定的class总量
		}

		nftClass, err := models.TClasses().AllG(context.Background())
		if err != nil && errors.Cause(err) == sql.ErrNoRows {
			//404
		} else if err != nil {
			//500
			log.Error("query nft class total", "query nft class error:", err.Error())
		}
		for _ = range nftClass {
			metric.NewPrometheus().SyncClassTotal.With().Add(1) //系统创建的class总量
		}

		if err := orm.GetDB().Ping(); err != nil {
			metric.NewPrometheus().SyncMysqlException.With().Set(-1) //监控mysql连接
		}

		if err := redis.Client.Ping(redis.Client{}, context.Background()); err != nil {
			metric.NewPrometheus().SyncRedisException.With().Set(-1) //监控redis连接
		}

		if _, err := chain.GetSdkClient().GenConn(); err != nil {
			metric.NewPrometheus().SyncChainConnError.With().Set(-1) //访问链异常
		}
	}
	crontab.AddFunc(spec, task)
	crontab.Start()
}
