package zipkinClientHttp

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	//"os"

	//"contrib.go.opencensus.io/exporter/prometheus"
	//"contrib.go.opencensus.io/exporter/zipkin"
	// openzipkin "github.com/openzipkin/zipkin-go"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinreport "github.com/openzipkin/zipkin-go/reporter"
	// "go.opencensus.io/trace"
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

type resource struct {
	reporter zipkinreport.Reporter //zipkinHTTP.NewReporter(c.ServerAddr)
	tracer   *zipkin.Tracer
}

var res resource

//Destory release the resource
func Destroy() {
	fmt.Println("destroy")
	defer func() {
		if res.reporter != nil {
			res.reporter.Close()
		}
	}()
}

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

func WithClient(client *http.Client) zipkinhttp.ClientOption {
	return zipkinhttp.WithClient(client)
}
