package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	metricsprometheus "github.com/go-kit/kit/metrics/prometheus"
)

type prometheusModel struct {
	ApiMysqlException       *metricsprometheus.Gauge     // api_mysql_exception 监控mysql连接
	ApiRedisException       *metricsprometheus.Gauge     // api_redis_exception 监控redis连接
	ApiHttpRequestCount     *metricsprometheus.Counter   // api_http_request_total api调用次数
	ApiHttpRequestRtSeconds *metricsprometheus.Histogram // api_http_request_rt_seconds api响应时间
	ApiRootBalanceAmount    *metricsprometheus.Gauge     // api_root_balance_amount 系统root账户余额
	SyncTxFailedTotal       *metricsprometheus.Gauge     // sync_tx_failed_total 系统失败的交易总量
	SyncTxPendingTotal      *metricsprometheus.Gauge     // sync_tx_pending_total 系统未完成的交易总量
	SyncChainConnError      *metricsprometheus.Gauge     // sync_chain_conn_error 访问链异常
	SyncMysqlException      *metricsprometheus.Gauge     // sync_mysql_exception 监控mysql连接
	SyncRedisException      *metricsprometheus.Gauge     // sync_redis_exception 监控redis连接
	SyncNftLockedTotal      *metricsprometheus.Gauge     // sync_nft_locked_total 系统锁定的nft总量
	SyncClassLockedTotal    *metricsprometheus.Gauge     // sync_class_locked_total 系统锁定的class总量
	SyncNftTotal            *metricsprometheus.Gauge     // sync_nft_total 系统创建的nft总量
	SyncClassTotal          *metricsprometheus.Gauge     // sync_class_total 系统创建的class总量
	SyncTxUndoSeconds       *metricsprometheus.Gauge     // sync_tx_undo_seconds 还未广播的交易
	SyncTxPendingSeconds    *metricsprometheus.Gauge     // sync_tx_pending_seconds 广播完成，但是还未同步到交易状态的交易
}

var (
	prometheusCache *prometheusModel
	prometheusOnce  sync.Once
)

// NewPrometheus 单列模式
func NewPrometheus() *prometheusModel {
	prometheusOnce.Do(func() {
		prometheusCache = &prometheusModel{}
	})
	return prometheusCache
}

// InitPrometheus 注册prometheus配置
func (p *prometheusModel) InitPrometheus() {
	// api_mysql_exception
	apiMysqlException := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "api",
		Name:      "mysql_exception",
		Help:      "mysql exception",
	}, []string{})

	// api_redis_exception
	apiRedisException := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "api",
		Name:      "redis_exception",
		Help:      "redis exception",
	}, []string{})

	// api_http_request_total
	apiHttpRequestTotal := metricsprometheus.NewCounterFrom(prometheus.CounterOpts{
		Subsystem: "api",
		Name:      "http_request_total",
		Help:      "http request total",
	}, []string{"code", "method", "uri"})

	// api_http_request_rt_seconds
	apiHttpRequestRtSeconds := metricsprometheus.NewHistogramFrom(prometheus.HistogramOpts{
		Subsystem: "api",
		Name:      "http_request_rt_seconds",
		Help:      "http request rt seconds",
	}, []string{"method", "uri"})

	// api_root_balance_amount
	apiRootBalanceAmout := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "api",
		Name:      "root_balance_amount",
		Help:      "root balance amount",
	}, []string{})

	// sync_tx_failed_total
	syncTxFailedTotal := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "tx_failed_total",
		Help:      "tx failed total",
	}, []string{})

	// sync_tx_pending_total
	syncTxPendingTotal := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "tx_pending_total",
		Help:      "tx pending total",
	}, []string{})

	// sync_chain_conn_error
	syncChainConnError := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "chain_conn_error",
		Help:      "chain conn error",
	}, []string{})

	// sync_mysql_exception
	syncMysqlException := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "mysql_exception",
		Help:      "mysql exception",
	}, []string{})

	// sync_redis_exception
	syncRedisException := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "redis_exception",
		Help:      "redis exception",
	}, []string{})

	// sync_nft_locked_total
	syncNftLockedTotal := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "nft_locked_total",
		Help:      "nft locked total",
	}, []string{})

	// sync_class_locked_total
	syncClassLockedTotal := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "class_locked_total",
		Help:      "class locked total",
	}, []string{})

	// sync_nft_total
	syncNftTotal := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "nft_total",
		Help:      "nft total",
	}, []string{})

	// sync_class_total
	syncClassTotal := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "class_total",
		Help:      "class total",
	}, []string{})

	// sync_tx_undo_seconds
	syncTxUndoSeconds := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "tx_undo_seconds",
		Help:      "tx undo seconds",
	}, []string{})

	// sync_tx_pending_seconds
	syncTxPendingSeconds := metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "sync",
		Name:      "tx_pending_seconds",
		Help:      "tx pending seconds",
	}, []string{})

	apiMysqlException.With([]string{}...).Set(0)
	apiRedisException.With([]string{}...).Set(0)
	syncTxFailedTotal.With([]string{}...).Set(0)
	syncTxPendingTotal.With([]string{}...).Set(0)
	syncChainConnError.With([]string{}...).Set(0)
	syncMysqlException.With([]string{}...).Set(0)
	syncRedisException.With([]string{}...).Set(0)
	syncNftLockedTotal.With([]string{}...).Set(0)
	syncClassLockedTotal.With([]string{}...).Set(0)
	syncNftTotal.With([]string{}...).Set(0)
	syncClassTotal.With([]string{}...).Set(0)
	syncTxUndoSeconds.With([]string{}...).Set(0)
	syncTxPendingSeconds.With([]string{}...).Set(0)

	p.ApiMysqlException = apiMysqlException
	p.ApiRedisException = apiRedisException
	p.ApiHttpRequestCount = apiHttpRequestTotal
	p.ApiHttpRequestRtSeconds = apiHttpRequestRtSeconds
	p.ApiRootBalanceAmount = apiRootBalanceAmout
	p.SyncTxFailedTotal = syncTxFailedTotal
	p.SyncTxPendingTotal = syncTxPendingTotal
	p.SyncChainConnError = syncChainConnError
	p.SyncMysqlException = syncMysqlException
	p.SyncRedisException = syncRedisException
	p.SyncNftLockedTotal = syncNftLockedTotal
	p.SyncClassLockedTotal = syncClassLockedTotal
	p.SyncNftTotal = syncNftTotal
	p.SyncClassTotal = syncClassTotal
	p.SyncTxUndoSeconds = syncTxUndoSeconds
	p.SyncTxPendingSeconds = syncTxPendingSeconds
}
