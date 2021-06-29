package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"

	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinmw "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	tracer *openzipkin.Tracer
	port   = ":8089"
)

func main() {

	localEndpoint, err := openzipkin.NewEndpoint("zipkin-opentracing", "192.168.1.61:8082")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	nativeTracer, err := openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		panic(fmt.Sprintf("err:%s", err))
	}
	tracer = nativeTracer //zipkinot.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(zipkinot.Wrap(nativeTracer))
	mux := http.NewServeMux()
	mux.HandleFunc("/list", list)
	mux.HandleFunc("/client", client)

	// h := middler(mux)
	h := nethttp.Middleware(opentracing.GlobalTracer(), mux)
	log.Printf("Server listening! %s ...", port)
	log.Fatal(http.ListenAndServe(port, h))

}

func client(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	span, _ := opentracing.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	x := rand.Intn(10) + 3
	time.Sleep(time.Duration(x) * time.Millisecond)
	w.Write([]byte("hello client"))
}
func list(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	span, _ := opentracing.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	database(r)
	serviceb(r)
	res := strings.Repeat("o", rand.Intn(100)+1)
	time.Sleep(time.Duration(rand.Intn(10)+1) * time.Millisecond)
	w.Write([]byte("Hello, w" + res + "rld!"))
}

func database(r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "database")
	defer span.Finish()
	cache(r)
	x := rand.Intn(4) + 10
	time.Sleep(time.Duration(x) * time.Millisecond)
	span.SetTag("sleep-time", fmt.Sprintf("database-cost:%d", x))
}

func cache(r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "cache")
	defer span.Finish()
	x := rand.Intn(4) + 10
	time.Sleep(time.Duration(x) * time.Millisecond)
	// span.LogFields Annotate(time.Now(), fmt.Sprintf("cost:%d", x))
}

func serviceb(r *http.Request) {
	span, pc := opentracing.StartSpanFromContext(r.Context(), "serviceb")
	defer span.Finish()
	time.Sleep(time.Duration(rand.Intn(10)+10) * time.Millisecond)
	servicec(pc) // servicec is childof serviceb
	// span.Annotate(time.Now(), "endtime")
	//span.LogEvent("endtme")
	span.LogKV("endtme", time.Now().String())
}

//func servicec(r *http.Request) {
func servicec(c context.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "servicec")
	defer span.Finish()
	time.Sleep(time.Duration(rand.Intn(10)+10) * time.Millisecond)
	span.SetTag("servicec", "C") // set tags for search servicec
	doclient(ctx)
	AskGoogle(ctx)
}

func doclient(c context.Context) {

	// initiate a call to some_func
	addrServ := "127.0.0.1:8089"
	url := fmt.Sprintf("http://%s/client", addrServ)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("unable to create http request: %+v\n", err)
	}
	req = req.WithContext(c) // use parent context
	client, err := zipkinmw.NewClient(tracer, zipkinmw.ClientTrace(true), zipkinmw.ClientTags(map[string]string{"type:": "from-raw-http-client"}))
	if err != nil {
		log.Fatalf("err NewClient %s", err)
	}
	res, err := client.DoWithAppSpan(req, "other_svr_client")
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

func AskGoogle(ctx context.Context) error {
	client := &http.Client{Transport: &nethttp.Transport{}}
	addrServ := "127.0.0.1:8089"
	url := fmt.Sprintf("http://%s/client", addrServ)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx) // extend existing trace, if any

	req, ht := nethttp.TraceRequest(opentracing.GlobalTracer(), req, nethttp.OperationName("ask-google"))
	defer ht.Finish()

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	bye, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("client-resp:%s", bye)

	return nil
}
