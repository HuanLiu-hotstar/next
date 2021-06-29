package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	// jaegerlog "github.com/uber/jaeger-client-go/log"
	// "github.com/uber/jaeger-lib/metrics"
)

func main() {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	localReportEndpoint := "localhost:6831"
	cfg := jaegercfg.Configuration{
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: localReportEndpoint,
		},
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, //sample type
			Param: 1,                       // sample prob
		},
	}

	// Initialize tracer with a logger and a metrics factory
	serviceName := "hello-jaeger"
	closer, err := cfg.InitGlobalTracer(
		serviceName,
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()
	doclient()
}

func doclient() {
	//tracer := opentracing.GlobalTracer()
	//span := tracer.StartSpan("world")
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "world")
	defer span.Finish()
	x := rand.Intn(100) + 50
	span.SetTag("cost-time", fmt.Sprintf("%d", x))
	time.Sleep(time.Duration(x) * time.Millisecond)
	otherclient(span.Context())
	other3client(ctx)
	log.Printf("finish doclient")
}
func otherclient(parent opentracing.SpanContext) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan(
		"GetFeed",
		opentracing.ChildOf(parent),
	)

	defer span.Finish()
	x := rand.Intn(100) + 100
	time.Sleep(time.Duration(x) * time.Millisecond)
	span.SetTag("other-cost", fmt.Sprintf("%d", x))
}

func other3client(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "client3")
	defer span.Finish()
	x := rand.Intn(100) + 100
	time.Sleep(time.Duration(x) * time.Millisecond)
	span.SetTag("client3-cost", fmt.Sprintf("%d", x))
}
