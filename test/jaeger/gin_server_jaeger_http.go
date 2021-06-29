package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/opentracing-contrib/go-stdlib/nethttp"

	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go"
)

var (
	port  = ":8080"
	Lport = flag.Int("port", 8090, "listen port")
)

func Init(serviceName string) func() {
	localReportEndpoint := "localhost:6831"
	cfg := jaegercfg.Configuration{
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: localReportEndpoint,
		},
		Sampler: &jaegercfg.SamplerConfig{
			//Type:  jaeger.SamplerTypeConst, //sample type
			Type:  jaeger.SamplerTypeProbabilistic, //sample type
			Param: 1,                               // sample prob
		},
	}

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		s := fmt.Sprintf("Could not initialize jaeger tracer: %s", err.Error())
		panic(s)
	}
	return func() {
		defer closer.Close()
	}
}
func main() {
	flag.Parse()
	if Lport != nil {
		port = fmt.Sprintf(":%d", *Lport)
	}
	// Initialize tracer with a logger and a metrics factory
	serviceName := "gin-jaeger"
	f := Init(serviceName)
	defer f()
	rgin := gin.Default()
	// http.HandleFunc("/list", handler)
	log.Printf("listen %s", port)
	m := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("hehlow")
	}
	rgin.Use(func(c *gin.Context) {
		x := nethttp.Middleware(opentracing.GlobalTracer(), http.HandlerFunc(m))
		x.ServeHTTP(c.Writer, c.Request)
		c.Next()
	})
	rgin.GET("/list", func(c *gin.Context) {
		handler(c.Writer, c.Request)
		c.Next()
	})
	rgin.Run(port)
	http.ListenAndServe(
		port,
		// use nethttp.Middleware to enable OpenTracing for server
		nethttp.Middleware(opentracing.GlobalTracer(), http.DefaultServeMux))

	// if err := http.ListenAndServe(port, nil); err != nil {
	// 	log.Fatal("err:%s", err)
	// }
}
func handler(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "jaeger")
	defer span.Finish()
	span.SetTag("key", "world")
	x := rand.Intn(100) + 50
	time.Sleep(time.Duration(x) * time.Millisecond)
	data := fmt.Sprintf("sleep:%d", x)
	log.Printf("data:%s", data)
	callclient(ctx)
	w.Write([]byte(data))
}
func callclient(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "child-client")
	defer span.Finish()
	x := rand.Intn(100) + 50
	time.Sleep(time.Duration(x) * time.Millisecond)
}
