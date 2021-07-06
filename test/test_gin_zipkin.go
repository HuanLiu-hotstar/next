package main

import (
	"fmt"
	"log"
	// "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	openzipkin "github.com/openzipkin/zipkin-go"
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
	//exporter := zipkin.NewExporter(r.reporter, endpoint)
	r.tracer, err = openzipkin.NewTracer(r.reporter, openzipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		panic(fmt.Sprintf("err:%s", err))
	}
	opentracing.SetGlobalTracer(zipkinot.Wrap(r.tracer))
	//trace.RegisterExporter(exporter)
	//trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return nil
}
func Destroy() {
	defer func() {
		r.reporter.Close()
	}()
}

func dohttp(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "dohttp")
	defer span.Finish()
	time.Sleep(1 * time.Second)
}
func main() {

	Init()
	defer Destroy()
	rgin := gin.Default()

	rgin.Use(ginhttp.Middleware(opentracing.GlobalTracer()))
	rgin.GET("/",
		func(c *gin.Context) {
			time.Sleep(500 * time.Millisecond)
			c.JSON(200, gin.H{"code": 200, "msg": "OK1"})
		})
	rgin.GET("/list",
		func(c *gin.Context) {
			// time.Sleep(500 * time.Millisecond)
			span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "list")
			defer span.Finish()
			dohttp(c)
			c.JSON(200, gin.H{"code": 200, "msg": "OK2"})
		})
	rgin.Run(":8080")

}
