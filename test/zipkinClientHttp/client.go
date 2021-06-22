package zipkinClientHttp

import (
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"time"
	//"os"

	zzipkin "contrib.go.opencensus.io/exporter/zipkin"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	// zipkinreport "github.com/openzipkin/zipkin-go/reporter"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
)

func Init(opts ...ConfigOpt) {
	ip := GetLocalIP()
	c := Config{
		LocalAddr:       ip,
		LocalServerName: ip,
		ServerAddr:      ZipkinServerAddr,
		ServiceName:     ZipkinServerName,
	}
	for _, o := range opts {
		o(&c)
	}
	fmt.Println(c)
	localEndpoint, err := zipkin.NewEndpoint(c.LocalServerName, c.LocalAddr)
	//fmt.Println(localEndpoint)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Zipkin exporter: %v", err))
	}
	res.reporter = zipkinHTTP.NewReporter(c.ServerAddr)
	exporter := zzipkin.NewExporter(res.reporter, localEndpoint)
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	// initialize our tracer
	res.tracer, err = zipkin.NewTracer(res.reporter, zipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		panic(fmt.Sprintf("unable to create tracer: %+v\n", err))
	}
	//return reporter
}

func NewClient(opt ...zipkinhttp.ClientOption) (*zipkinhttp.Client, error) {
	if res.tracer == nil {
		panic(fmt.Sprintf("err not init tracer"))
	}
	opts := []zipkinhttp.ClientOption{WithClient(&http.Client{Timeout: 5 * time.Second}), zipkinhttp.ClientTrace(true)}
	opts = append(opts, opt...)
	fmt.Println(len(opts))
	client, err := zipkinhttp.NewClient(res.tracer, opts...)
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}
	return client, err
}
