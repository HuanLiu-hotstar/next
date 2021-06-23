# Monitoring,Tracing

## What's Tracing and Monitoring?

## Tracing function

## Tracing Concept
- tree
- span
- annotation:
- tag: user defined k-v pair for query 
## Architecture of Tracing System


## Implementation of Tracing System

### http & rpc implementation

- http: Write the Tracing info in http-header 


```go
// SpanContext holds the context of a Span.
type SpanContext struct {
	TraceID  TraceID `json:"traceId"`
	ID       ID      `json:"id"`
	ParentID *ID     `json:"parentId,omitempty"`
	Debug    bool    `json:"debug,omitempty"`
	Sampled  *bool   `json:"-"`
	Err      error   `json:"-"`
}

// Default B3 Header keys
const (
	TraceID      = "x-b3-traceid"
	SpanID       = "x-b3-spanid"
	ParentSpanID = "x-b3-parentspanid"
	Sampled      = "x-b3-sampled"
	Flags        = "x-b3-flags"
	Context      = "b3"
)

type tracer interface {
	StartSpan(name)
	InjectHttp(SpanContext)
	ExtractHttp() SpanContext
}



// delegate the http.Handler with tracing info
type handler struct {
  tracer *Tracer // tracing the request and response
  handle http.Handler // the true http-handler 
  //   ... 
}


// InterceptResponse dalegate the http-response with tracing  info
type InterceptResponse struct {
  w http.ResponseWriter 
  statusCode int 
}

// ServeHTTP implements http.Handler.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  //tracing parameter init from http.Request

  defer func(){
  //   reporte tracing 
  }()
  h.handle.ServeHTTP(h.wrap(w),h.wrap(r)) // true http-call but with the wrap ResponseWriter and http.Request.
}

func main() {
	h := handler{
		tracer : newTracer(),
		handler : func(w http.ResponseWriter, r *http.Request){
			w.Write("hello")
		}
	}
	http.HandleFunc("/uri",h)
}

```
- rpc : Write the tracing with the intercept handler 

- Open-Tracing standard
- zipkin: twitter
- jaeger : uber
- Skywalking:Huawei
- CAT: meituan

## standard http.Handler with zipkin 

- example for standard http.Handler delegate

```go
func main() {

	//tracer init
	LocalServerName := "My-Server"
	LocalAddr := "192.168.1.1"
	ZipkinServerAddr := "http://localhost:9411/api/v2/spans"
	localEndpoint, err := zipkin.NewEndpoint(LocalServerName, LocalAddr)
	
  	// create a reporter to report data the zipkin-server 
	reporter = zipkinHTTP.NewReporter(ServerAddr)
	// initialize our tracer
	tracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(localEndpoint))
	
	// for func( w http.ResponseWriter, r *http.Request) 
	// we can use r.Context() to start a span
	f := func (w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "database")
		defer span.End()
		log.Printf("Serving request: %s", r.URL.Path)
		dologic(ctx)
	}

	// tracer usage
	golocic := func (c context.Context /*other paramater*/) {
		ctx,span := tracer.SpartSpan(c,spanName)
		defer span.End()
	}
	http.HandleFunc("/uri",f)
	http.ListenAndServe(":8080",nil)
}
	// another tracer init way 
func anotherInit() {

	// or we can use global tracer
	localEndPoint := NewLocalEndPoint(localAddr)
  	// create a exporter for real send data to zipkin-server
	exporter := zzipkin.NewExporter(reporter, localEndpoint)
	//register an exporter to export data to zipkin-Server
	// we can register more than one exporter to report with diff server
	trace.RegisterExporter(exporter)
  	//configure the sample method and other configure 
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	// usage for span 
	trace.StartSpan()
}

	

```

- http.Client delegate detail are [this](https://github.com/HuanLiu-hotstar/next/blob/main/test/zipkin/test_zipkin_raw_http_client.go)

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinmw "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	tracer *openzipkin.Tracer
)

func main() {

	localEndpoint, err := openzipkin.NewEndpoint("local-client", "192.168.1.61:8082")
	// err handle omit
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	//exporter := zipkin.NewExporter(reporter, localEndpoint)
	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	// err handle omit
	// create global zipkin traced http client
	client, err := zipkinmw.NewClient(tracer, zipkinmw.ClientTrace(true), zipkinmw.ClientTags(map[string]string{"type:": "from-raw-http-client"}))
	// err handle omit

	// initiate a call to some_func
	url := "http://localhost:8080/list"
	req, err := http.NewRequest("GET", url, nil)
	
	res, err := client.DoWithAppSpan(req, "raw-http-client")
	
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	log.Printf("%s", body)
}



```

## gin with zipkin

- gin middleware