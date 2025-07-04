/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	//pb "helloworld/helloworld"

	openzipkin "github.com/openzipkin/zipkin-go"
	// zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"github.com/HuanLiu-hotstar/demo-tracing/zipkin/demo/tracelib"
	// "github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"

	// otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	span, _ := tracelib.StartSpanFromContext(ctx, "hello")
	defer span.Finish()
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

var (
	tracer     *openzipkin.Tracer
	serverName = "grpc-greeter-server"
	localAddr  = "192.168.1.61"
	zipkinAddr = "http://localhost:9411/api/v2/spans" //http://localhost:9411/api/v2/spans
	// port       = 18084
)

func newTracer() *openzipkin.Tracer {
	localEndpoint, err := openzipkin.NewEndpoint("grpc", "192.168.1.61:8082")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	// defer reporter.Close()

	// ratio := 0.001
	// seed := int64(1000)
	// sample, err := openzipkin.NewBoundarySampler(ratio, seed)
	// if err != nil {
	// 	log.Fatal("err:%s", err)
	// }
	// tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint), openzipkin.WithSampler(sample))
	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		panic(fmt.Sprintf("err:%s", err))
	}
	return tracer
}
func main() {
	// newTracer()
	opts := []tracelib.ConfigOpt{tracelib.WithLocalAddr(localAddr), tracelib.WithLocalName(serverName),
		tracelib.WithZipkinSerAddr(zipkinAddr),
	}
	f := tracelib.InitTracer(opts...)
	defer f()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// m := map[string]string{"key": "value"}
	// middleware := grpc.StatsHandler(zipkingrpc.NewServerHandler(tracer, zipkingrpc.ServerTags(m)))
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())))

	// s := grpc.NewServer(middleware)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
