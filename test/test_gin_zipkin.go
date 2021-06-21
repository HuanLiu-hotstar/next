package main

import (
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	//"github.com/uber/jaeger-client-go/zipkin"
)

func main() {

	tracer, closer := jaeger.NewTracer(
		"serviceName",
		jaeger.NewConstSampler(true),
		jaeger.NewInMemoryReporter(),
	)
	defer closer.Close()

	fn := func(c *gin.Context) {
		span := opentracing.SpanFromContext(c.Request.Context())
		if span == nil {
			log.Fatal("Span is nil")
		}
	}

	r := gin.New()
	r.Use(ginhttp.Middleware(tracer))
	r.Run(":8081")
	group := r.Group("")
	group.GET("", fn)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal("Error non-nil %v", err)
	}
	r.ServeHTTP(w, req)
}
