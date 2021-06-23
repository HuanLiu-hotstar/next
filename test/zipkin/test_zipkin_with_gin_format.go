package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinmw "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	wraphh "github.com/turtlemonvh/gin-wraphh"
)

var (
	tracer *openzipkin.Tracer
)

func main() {

	localEndpoint, err := openzipkin.NewEndpoint("gin_foramt_svr", "192.168.1.61:8082")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	//exporter := zipkin.NewExporter(reporter, localEndpoint)
	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		panic(fmt.Sprintf("err:%s", err))
	}
	// mux := http.NewServeMux()
	// mux.HandleFunc("/list", list)
	rgin := gin.Default()
	spanName := "root"
	middler := zipkinmw.NewServerMiddleware(
		tracer,
		zipkinmw.SpanName(spanName),
		zipkinmw.TagResponseSize(true),
	)
	// http.Handler http.Handler

	rgin.Use(wraphh.WrapHH(middler))
	rgin.GET("/",
		func(c *gin.Context) {
			time.Sleep(500 * time.Millisecond)
			c.JSON(200, gin.H{"code": 200, "msg": "OK1"})
		})
	// rgin.GET("/list", list)
	rgin.GET("/list", gin.WrapF(list1))
	port := ":8080"
	rgin.Run(port)
	h := middler(rgin)
	log.Printf("Server listening! %s ...", port)
	log.Fatal(http.ListenAndServe(port, h))

}

func list(c *gin.Context) {
	// span, _ := tracer.StartSpanFromContext(c.Request.Context(), c.FullPath())
	// defer span.Finish()
	time.Sleep(100 * time.Millisecond)
	//c.Writer.Write()
	c.String(200, "hello world")

}
func list1(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	span, _ := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	//time.Sleep(6 * time.Second)
	database(r)
	serviceb(r)
	res := strings.Repeat("o", rand.Intn(100)+1)
	time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)
	w.Write([]byte("Hello, w" + res + "rld!"))
}

func database(r *http.Request) {
	cache(r)
	span, _ := tracer.StartSpanFromContext(r.Context(), "database")
	defer span.Finish()
	x := rand.Intn(4) + 100
	time.Sleep(time.Duration(x) * time.Millisecond)
	span.Tag("sleep-time", fmt.Sprintf("database-cost:%d", x))
}

func cache(r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), "cache")
	defer span.Finish()
	x := rand.Intn(4) + 100
	time.Sleep(time.Duration(x) * time.Millisecond)
	span.Annotate(time.Now(), fmt.Sprintf("cost:%d", x))
}

func serviceb(r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), "serviceb")
	defer span.Finish()
	time.Sleep(time.Duration(rand.Intn(100)+100) * time.Millisecond)
	servicec(r)
}

func servicec(r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), "servicec")
	defer span.Finish()
	time.Sleep(time.Duration(rand.Intn(700)+100) * time.Millisecond)
}
