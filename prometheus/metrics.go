package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	// Initial count.
	currentCount = 0

	// The Prometheus metric that will be exposed.
	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ars2023_http_hit_total",
			Help: "Total number of http hits.",
		},
	)

	createConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "create_config_http_hit_total",
			Help: "Total number of create config hits.",
		},
	)

	getAllConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_all_config_http_hit_total",
			Help: "Total number of get all config hits.",
		},
	)

	getConfigVersionsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_config_version_http_hit_total",
			Help: "Total number of get config version hits.",
		},
	)

	delConfigVersionsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_config_version_http_hit_total",
			Help: "Total number of del config version hits.",
		},
	)

	getConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_config_http_hit_total",
			Help: "Total number of get config hits.",
		},
	)

	delConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_config_http_hit_total",
			Help: "Total number of del config hits.",
		},
	)

	createGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "create_group_http_hit_total",
			Help: "Total number of create group hits.",
		},
	)

	getAllGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_all_group_http_hit_total",
			Help: "Total number of get all group hits.",
		},
	)

	getGroupVersionHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_group_version_http_hit_total",
			Help: "Total number of get group versions hits.",
		},
	)

	delGroupVersionHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_group_version_http_hit_total",
			Help: "Total number of del group versions hits.",
		},
	)

	getGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_group_http_hit_total",
			Help: "Total number of get group hits.",
		},
	)

	delGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_group_http_hit_total",
			Help: "Total number of del group hits.",
		},
	)

	appendGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "append_group_http_hit_total",
			Help: "Total number of append group hits.",
		},
	)

	getConfigByLabelsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_config_by_labels_http_hit_total",
			Help: "Total number of get config by labels hits.",
		},
	)

	delConfigByLabelsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_config_by_labels_http_hit_total",
			Help: "Total number of del config by labels hits.",
		},
	)

	// Add all metrics that will be resisted
	metricsList = []prometheus.Collector{
		httpHits,
		createConfigHits,
		getAllConfigHits,
		getConfigVersionsHits,
		delConfigVersionsHits,
		getConfigHits,
		delConfigHits,
		createGroupHits,
		getAllGroupHits,
		delGroupVersionHits,
		getGroupHits,
		delGroupHits,
		appendGroupHits,
		getConfigByLabelsHits,
		delConfigByLabelsHits,
	}

	// Prometheus Registry to register metrics.
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	// Register metrics that will be exposed.
	prometheusRegistry.MustRegister(metricsList...)
}

func MetricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func CountCreateConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		createConfigHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetAllConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getAllConfigHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetConfigVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigVersionsHits.Inc()
		f(w, r) // original function call
	}
}

func CountDelConfigVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delConfigVersionsHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigHits.Inc()
		f(w, r) // original function call
	}
}

func CountDelConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delConfigHits.Inc()
		f(w, r) // original function call
	}
}

func CountCreateGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		createGroupHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetAllGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getAllGroupHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetGroupVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getGroupVersionHits.Inc()
		f(w, r) // original function call
	}
}

func CountDelGroupVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delGroupVersionHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getGroupHits.Inc()
		f(w, r) // original function call
	}
}

func CountDelGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delGroupHits.Inc()
		f(w, r) // original function call
	}
}

func CountAppendGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		appendGroupHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetConfigByLabels(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigByLabelsHits.Inc()
		f(w, r) // original function call
	}
}

func CountDelConfigByLabels(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delConfigByLabelsHits.Inc()
		f(w, r) // original function call
	}
}
