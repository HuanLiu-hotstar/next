package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	// "strings"
	"time"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinmw "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	tracer     *openzipkin.Tracer
	serverName = "PC-Server"
	localAddr  = "192.168.1.63"
	port       = 8083
	umAddr     = "http://localhost:8084/um"
)

func main() {

	localEndpoint, err := openzipkin.NewEndpoint(serverName, fmt.Sprintf("%s:%d", localAddr, port))
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		panic(fmt.Sprintf("err:%s", err))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/pc", pc)
	// spanName := "root"
	middler := zipkinmw.NewServerMiddleware(
		tracer,
		zipkinmw.SpanName(serverName),
		zipkinmw.TagResponseSize(true),
	)

	h := middler(mux)
	log.Printf("Server listening! %d ...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), h))

}
func pc(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	callum(ctx)
	//span.Annotate(time.Now(), "some-work")
	x := rand.Intn(100) + 10
	//time.Sleep(time.Duration(x) * time.Millisecond)
	otherwork(ctx)

	w.Write([]byte(fmt.Sprintf(`{"pc":%d}`, x)))
	//call um
}
func otherwork(c context.Context) {
	span, _ := tracer.StartSpanFromContext(c, "other-work")
	defer span.Finish()
	x := rand.Intn(100) + 10
	time.Sleep(time.Duration(x) * time.Millisecond)

}
func callum(r context.Context) (string, error) {
	data := ""
	// create global zipkin traced http client
	co := zipkinmw.WithClient(&http.Client{Timeout: time.Second * 3})
	client, err := zipkinmw.NewClient(tracer, co, zipkinmw.ClientTrace(true), zipkinmw.ClientTags(map[string]string{"type:": "from-raw-http-client"}))
	if err != nil {
		log.Printf("unable to create client: %+v\n", err)
		return data, err
	}

	// initiate a call to some_func
	url := umAddr
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("unable to create http request: %+v\n", err)
		return data, err
	}

	req = req.WithContext(r)
	res, err := client.DoWithAppSpan(req, "um-client")
	if err != nil {
		log.Printf("unable to do http request: %+v\n", err)
		return data, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	log.Printf("%s", body)
	return string(body), err
}
