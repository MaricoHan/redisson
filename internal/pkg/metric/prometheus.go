package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	metricsprometheus "github.com/go-kit/kit/metrics/prometheus"
)

type prometheusModel struct {
	ApiMysqlException *metricsprometheus.Gauge // api_mysql_exception 监控mysql连接
	ApiRedisException *metricsprometheus.Gauge // api_redis_exception 监控redis连接
	// ApiHttpRequestCount     *metricsprometheus.Counter   // api_http_request_total api调用次数
	ApiHttpRequestRtSeconds *metricsprometheus.Histogram // api_http_request_rt_seconds api响应时间
	ApiRootBalanceAmount    *metricsprometheus.Gauge     // api_root_balance_amount 系统root账户余额
	ApiServiceRequests      *metricsprometheus.Counter   // api_service_requests
}

var (
	prometheusCache *prometheusModel
	prometheusOnce  sync.Once
)

// NewPrometheus 单例模式
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
	// apiHttpRequestTotal := metricsprometheus.NewCounterFrom(prometheus.CounterOpts{
	//	Subsystem: "api",
	//	Name:      "http_request_total",
	//	Help:      "http request total",
	// }, []string{"code", "method", "uri"})

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
	}, []string{"address", "denom"})

	// api_service_requests
	apiServiceRequests := metricsprometheus.NewCounterFrom(prometheus.CounterOpts{
		Subsystem: "api",
		Name:      "api_service_requests",
		Help:      "api service request",
	}, []string{"status", "method", "name"})

	apiMysqlException.With([]string{}...).Set(0)
	apiRedisException.With([]string{}...).Set(0)

	p.ApiMysqlException = apiMysqlException
	p.ApiRedisException = apiRedisException
	// p.ApiHttpRequestCount = apiHttpRequestTotal
	p.ApiHttpRequestRtSeconds = apiHttpRequestRtSeconds
	p.ApiRootBalanceAmount = apiRootBalanceAmout
	p.ApiServiceRequests = apiServiceRequests
}
