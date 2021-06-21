package zipkinClientHttp

import (
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"time"
	//"os"

	//"contrib.go.opencensus.io/exporter/prometheus"
	//"contrib.go.opencensus.io/exporter/zipkin"
	// openzipkin "github.com/openzipkin/zipkin-go"
	zzipkin "contrib.go.opencensus.io/exporter/zipkin"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinreport "github.com/openzipkin/zipkin-go/reporter"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	//logreporter "github.com/openzipkin/zipkin-go/reporter/log"
	"go.opencensus.io/trace"
)

type Config struct {
	LocalAddr       string
	LocalServerName string
	ServerAddr      string
	ServiceName     string
}
type ConfigOpt func(c *Config)

const (
	ZipkinServerAddr = "http://localhost:9411/api/v2/spans"
	ZipkinServerName = "ZipkinServer"
)

var (
	reporter zipkinreport.Reporter //zipkinHTTP.NewReporter(c.ServerAddr)
	tracer   *zipkin.Tracer
)

func WithLocalAddr(localAddr string) ConfigOpt {
	return func(c *Config) {
		c.LocalAddr = localAddr
	}
}
func WithLocalServerName(serverName string) ConfigOpt {
	return func(c *Config) {
		c.LocalServerName = serverName
	}
}
func WithServerAddr(serverAddr string) ConfigOpt {
	return func(c *Config) {
		c.ServerAddr = serverAddr
	}
}
func WithServiceName(serviceName string) ConfigOpt {
	return func(c *Config) {
		c.ServiceName = serviceName
	}
}
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
	reporter = zipkinHTTP.NewReporter(c.ServerAddr)
	exporter := zzipkin.NewExporter(reporter, localEndpoint)
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	// initialize our tracer
	tracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		panic(fmt.Sprintf("unable to create tracer: %+v\n", err))
	}
	//return reporter
}
func Destroy() {
	fmt.Println("destroy")
	defer reporter.Close()
}
func WithClient(client *http.Client) zipkinhttp.ClientOption {
	return zipkinhttp.WithClient(client)
}
func NewClient(opt ...zipkinhttp.ClientOption) (*zipkinhttp.Client, error) {
	if tracer == nil {
		panic(fmt.Sprintf("err not init tracer"))
	}
	opts := []zipkinhttp.ClientOption{WithClient(&http.Client{Timeout: 5 * time.Second}), zipkinhttp.ClientTrace(true)}
	opts = append(opts, opt...)
	fmt.Println(len(opts))
	// create global zipkin traced http client
	//client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	client, err := zipkinhttp.NewClient(tracer, opts...)
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}
	return client, err
}

// func registerZipkin() zipkinreport.Reporter {
// 	localEndpoint, err := zipkin.NewEndpoint("golangsvc", "localhost:8081")
// 	if err != nil {
// 		log.Fatalf("Failed to create Zipkin exporter: %v", err)
// 		return nil
// 	}
// 	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
// 	exporter := zzipkin.NewExporter(reporter, localEndpoint)
// 	trace.RegisterExporter(exporter)
// 	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
// 	return reporter
// }

/*
func TestMain(testing test.Test) {
	// set up a span reporter
	reporter1 := logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags))
	defer reporter1.Close()

	reporter := registerZipkin()
	defer reporter.Close()

	addr := "127.0.0.1:9410"
	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint("myService", addr)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// create global zipkin traced http client
	client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	// initiate a call to some_func
	addrServ := "127.0.0.1:8080"
	url := fmt.Sprintf("http://%s/list", addrServ)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("unable to create http request: %+v\n", err)
	}

	res, err := client.DoWithAppSpan(req, "client-list")
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("err read body %s", err)
	}

	// Output:
	log.Printf("result %s", body)
}
*/
