package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	openzipkin "github.com/openzipkin/zipkin-go"
	// zipkinmw "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	// wraphh "github.com/turtlemonvh/gin-wraphh"
	"github.com/openzipkin/zipkin-go/propagation/b3"
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
	// spanName := "root"
	// middler := zipkinmw.NewServerMiddleware(
	// 	tracer,
	// 	zipkinmw.SpanName(spanName),
	// 	zipkinmw.TagResponseSize(true),
	// )
	// http.Handler http.Handler

	// rgin.Use(wraphh.WrapHH(middler))
	// rgin.Use(func(c *gin.Context) {
	// 	span, ctx := tracer.StartSpanFromContext(c.Request.Context(), c.FullPath()+"-http-")
	// 	defer span.Finish()
	// 	c.Request = c.Request.WithContext(ctx)

	// 	c.Next()
	// })
	rgin.Use(middlewareGin(tracer))

	rgin.GET("/",
		func(c *gin.Context) {
			time.Sleep(500 * time.Millisecond)
			c.JSON(200, gin.H{"code": 200, "msg": "OK1"})
		})
	rgin.GET("/list", func(c *gin.Context) {
		glist(c)
		// c.Next()
	})
	// rgin.GET("/list", gin.WrapF(list1))
	port := ":8080"
	rgin.Run(port)
	// h := middler(rgin)
	log.Printf("Server listening! %s ...", port)
	// log.Fatal(http.ListenAndServe(port, h))

}
func glist2(r *gin.Context) {
	span, _ := tracer.StartSpanFromContext(r.Request.Context(), "glist2")
	//	span, ctx := tracer.StartSpanFromContext(r, "glist2")
	defer span.Finish()
	// r.Request = r.Request.WithContext(ctx)
	glist3(r.Request.Context())
	x := rand.Intn(50) + 10
	time.Sleep(time.Duration(x) * time.Millisecond)
}
func glist3(c context.Context) {
	span, _ := tracer.StartSpanFromContext(c, "glist3")
	defer span.Finish()
	x := rand.Intn(80) + 10
	time.Sleep(time.Duration(x) * time.Millisecond)
}
func glist(r *gin.Context) {
	log.Printf("Serving request: %s", r.FullPath())
	glist2(r)
	res := strings.Repeat("o", rand.Intn(100)+1)
	time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)
	r.String(200, "Hello, w"+res+"rld!")
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

func middlewareGin(tracer *openzipkin.Tracer) gin.HandlerFunc {

	return func(gc *gin.Context) {
		parent := tracer.Extract(b3.ExtractHTTP(gc.Request))
		var span openzipkin.Span
		// no parent span
		if parent.Err == nil {
			var ctx context.Context
			span, ctx = tracer.StartSpanFromContext(gc.Request.Context(), gc.FullPath(), openzipkin.Kind(model.Server))
			defer span.Finish()
			gc.Request = gc.Request.WithContext(ctx)
		} else {
			span = tracer.StartSpan(gc.FullPath(), openzipkin.Parent(parent), openzipkin.Kind(model.Server))
			defer span.Finish()
		}
		openzipkin.TagHTTPMethod.Set(span, gc.Request.Method)
		openzipkin.TagHTTPPath.Set(span, gc.Request.URL.Path)
		openzipkin.TagHTTPRequestSize.Set(span, fmt.Sprintf("%d", gc.Request.ContentLength))
		//
		gc.Next()
		//
		if statusCode := gc.Writer.Status(); statusCode < 200 || statusCode > 299 {
			openzipkin.TagHTTPStatusCode.Set(span, strconv.Itoa(gc.Writer.Status()))
			openzipkin.TagHTTPResponseSize.Set(span, strconv.Itoa(gc.Writer.Size()))
		}
	}
}
