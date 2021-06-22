package main

import (
	// "os"
	"log"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	openzipkin "github.com/openzipkin/zipkin-go"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	// zipkin middleware
	//zipkinmw "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

type ConfigOpt func(c *Config)
type Config struct {
	ReporterAddr string
	LocalAddr    string
	LocalSvrName string
}
type Result struct {
	zkTracer opentracing.Tracer
	reporter reporter.Reporter
	tracer   *openzipkin.Tracer
}

var r Result

func Init(opt ...ConfigOpt) func(c *gin.Context) {
	c := Config{
		ReporterAddr: "http://localhost:9411/api/v2/spans",
		LocalAddr:    "localhost:80",
		LocalSvrName: "LocalSvr",
	}
	for _, o := range opt {
		o(&c)
	}

	endpoint, err := openzipkin.NewEndpoint(c.LocalSvrName, c.LocalAddr)
	if err != nil {
		log.Fatalf("unable to create local endpoint:%s  %+v\n", c.LocalAddr, err)
	}
	r.reporter = zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	exporter := zipkin.NewExporter(r.reporter, endpoint)
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return nil
	// return ZipKinMiddleware
}
func Destroy() {
	defer func() {
		r.reporter.Close()
	}()
}

// func ZipKinMiddleware(c *gin.Context) {
// 	span := r.zkTracer.StartSpan(c.FullPath())
// 	defer span.Finish()
// 	c.Next()
// }
func dohttp(c *gin.Context) {
	//span := r.zkTracer.StartSpan("dohttp")
	_, span := trace.StartSpan(c.Request.Context(), "dohttp")
	defer span.End()
	time.Sleep(1 * time.Second)
}
func main() {

	//zipkinMiddle := Init()
	Init()
	defer Destroy()
	// spanName := "test-svr"

	// f := func(c *gin.Context) {
	// 	c.Next()
	// }
	rgin := gin.Default()
	// 第三步: 添加一个 middleWare, 为每一个请求添加span
	// r.Use(zipkinMiddle)
	rgin.GET("/",
		func(c *gin.Context) {
			time.Sleep(500 * time.Millisecond)
			c.JSON(200, gin.H{"code": 200, "msg": "OK1"})
		})
	rgin.GET("/list",
		func(c *gin.Context) {
			// time.Sleep(500 * time.Millisecond)
			dohttp(c)
			c.JSON(200, gin.H{"code": 200, "msg": "OK2"})
		})
	// r.Run(":8080")
	// handler := zipkinmw.NewServerMiddleware(
	// 	r.tracer,
	// 	zipkinmw.SpanName(spanName),
	// 	zipkinmw.TagResponseSize(true),
	// 	// mw.ServerTags(tags),
	// )(rgin)
	h := &ochttp.Handler{Handler: rgin}
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		log.Fatal("Failed to register ochttp.DefaultServerViews")
	}
	http.Handle("/", h)
	port := ":8080"
	http.ListenAndServe(port, nil)
}
