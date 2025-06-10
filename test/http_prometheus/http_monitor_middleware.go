package main

// package prometheus_middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto" // auto register metrics with init
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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

var (
	totalRequests  *prometheus.CounterVec
	responseStatus *prometheus.CounterVec
	httpDuration   *prometheus.HistogramVec
	namespace      = "namespace"
	subsystem      = "subsystem"
	reqLabels      = []string{"status", "endpoint", "method"}
	one            sync.Once
)

func getCounterVecOpt(name, help string) prometheus.CounterOpts {
	fmt.Println(namespace, subsystem, name, help)
	return prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	}
}
func newMetric() {
	totalRequests = prometheus.NewCounterVec(
		getCounterVecOpt("requests_count", "Number of get requests."),
		reqLabels,
	)

	responseStatus = prometheus.NewCounterVec(
		getCounterVecOpt("response_status", "Status of HTTP response."),
		reqLabels,
	)

	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "response_time_seconds",
		Help:      "Duration of HTTP requests.",
	}, reqLabels)
}
func registerMetric() {
	list := []prometheus.Collector{totalRequests /*responseStatus,*/, httpDuration}
	for _, c := range list {
		fmt.Println(c)
		if err := prometheus.Register(c); err != nil {
			panic(fmt.Sprintf("err:%s", err))
		}
	}
}

// Init Namespace: should be the project name
// SubSystem: should be the server name
// using alphabet and _ for string name
func Init(Namespace string, SubSystem string) func(next http.Handler) http.Handler {
	one.Do(func() {
		namespace, subsystem = Namespace, SubSystem
		newMetric()
		registerMetric()
	})
	return PrometheusMiddleware
}
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// route := mux.CurrentRoute(r)
		// path, _ := route.GetPathTemplate()
		path := r.URL.Path
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
	listen()
}
func main2() {
	router := mux.NewRouter()
	middleware := Init("chat_infra", "router")
	Init("chat_infra", "router")
	router.Use(middleware)
	// router.Use(PrometheusMiddleware)
	// Prometheus endpoint
	router.Path("/prometheus").Handler(promhttp.Handler())
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world http2\n"))
	}
	// Serving static files
	router.PathPrefix("/hello").Handler(http.HandlerFunc(f))

	fmt.Println("Serving requests on port 9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}
func F(h http.Handler) http.HandlerFunc {
	// return http.HandleFunc(h.ServeHTTP)
	return nil
}
func listen() {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world http2\n"))
	}
	fmt.Println("hello world2222")
	// f is type of  HandlerFunc
	middleware := Init("chat_infra", "router")
	// f2 := Func(http.HandleFunc(f))
	f2 := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello owrld")
		middleware(http.Handler(http.HandlerFunc(f))).ServeHTTP(w, r)
	}
	// middleware()
	// func(next http.Handler) http.Handler
	// HandlerFunc
	// middleware(http.Handler(f))
	// http.HandlerFunc(f) -> is a http.Handler
	// middleware(http.Handler(http.HandlerFunc(f)))
	// f3 := func(w http.ResponseWriter, r *http.Request) {
	// 	promhttp.Handler().ServeHTTP(w, r)
	// }

	// http.HandleFunc("/prometheus", f3)

	http.HandleFunc("/world", f2)
	http.ListenAndServe(":9001", nil)
}

// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }
// type HandlerFunc func(ResponseWriter, *Request)
// HandleFunc is a type of Handler
// ServeHTTP calls f(w, r).
// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
// 	f(w, r)
// }
