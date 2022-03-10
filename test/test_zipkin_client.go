package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//	"net/http/httptest"
	//"time"

	//"github.com/gorilla/mux"
	//"contrib.go.opencensus.io/exporter/prometheus"
	//"contrib.go.opencensus.io/exporter/zipkin"
	// openzipkin "github.com/openzipkin/zipkin-go"
	zzipkin "contrib.go.opencensus.io/exporter/zipkin"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	logreporter "github.com/openzipkin/zipkin-go/reporter/log"
	"go.opencensus.io/trace"
)

var (
	addrServ = "127.0.0.1:8080"
	port     = flag.Int("port", 8080, "server listen port")
)

func registerZipkin() reporter.Reporter {
	localEndpoint, err := zipkin.NewEndpoint("golangsvc", "localhost:8081")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
		return nil
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	exporter := zzipkin.NewExporter(reporter, localEndpoint)
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return reporter
}
func main() {
	flag.Parse()
	if port != nil {
		addrServ = fmt.Sprintf("127.0.0.1:%d", *port)
	}
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
	// addrServ := "127.0.0.1:8080"
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

/*
func someFunc(client *zipkinhttp.Client, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("some_function called with method: %s\n", r.Method)

		// retrieve span from context (created by server middleware)
		span := zipkin.SpanFromContext(r.Context())
		span.Tag("custom_key", "some value")

		// doing some expensive calculations....
		time.Sleep(25 * time.Millisecond)
		span.Annotate(time.Now(), "expensive_calc_done")

		newRequest, err := http.NewRequest("POST", url+"/other_function", nil)
		if err != nil {
			log.Printf("unable to create client: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}

		ctx := zipkin.NewContext(newRequest.Context(), span)

		newRequest = newRequest.WithContext(ctx)

		res, err := client.DoWithAppSpan(newRequest, "other_function")
		if err != nil {
			log.Printf("call to other_function returned error: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		res.Body.Close()
	}
}

func otherFunc(client *zipkinhttp.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("other_function called with method: %s\n", r.Method)
		time.Sleep(50 * time.Millisecond)
	}
}

*/
