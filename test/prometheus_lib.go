package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusClient struct {
	Namespace string
	Subsystem string
	m         map[string]prometheus.Counter
}

func NewPrometheusClient(namespace, subsystem string) *PrometheusClient {
	return &PrometheusClient{
		Namespace: namespace,
		Subsystem: subsystem,
		m:         make(map[string]prometheus.Counter),
	}
}
func (p *PrometheusClient) GetCounter(name string) *prometheus.Counter {
	if c, ok := p.m[name]; ok {
		return &c
	}
	t := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: p.Namespace,
		Subsystem: p.Subsystem,
		Name:      name,
		Help:      "count of " + name,
	})
	p.m[name] = t
	return &t
}
func (p *PrometheusClient) Add(name string, del float64) {
	if _, ok := p.m[name]; ok {
		p.m[name].Add(del)
		return
	}
	p.GetCounter(name)
	p.Add(name, del)
}

const (
	Namespace = "chat-infra"
	Subsystem = "router"
)

var (
	p *PrometheusClient
)

func init() {
	p = NewPrometheusClient(Namespace, Subsystem)
}
func getFriend(pid string) []string {
	(*p.GetCounter("get_friend")).Add(1)
	list := []string{}
	return list
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

var valueCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "count_value",
	Help: "count some value.",
}, []string{"name"})

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}

type Func func(w http.ResponseWriter, r *http.Request)

func (f Func) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}
func main() {
	router := mux.NewRouter()
	router.Use(prometheusMiddleware)

	// Prometheus endpoint
	router.Path("/prometheus").Handler(promhttp.Handler())
	// if err := prometheus.Register(valueCount); err != nil {
	// 	fmt.Println(err)
	// }
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world http2\n"))
	}
	// Serving static files
	router.PathPrefix("/hello").Handler(Func(f))

	fmt.Println("Serving requests on port 9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}
