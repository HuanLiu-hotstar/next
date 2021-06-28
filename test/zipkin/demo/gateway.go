package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	// "strings"
	"time"

	"google.golang.org/grpc"

	openzipkin "github.com/openzipkin/zipkin-go"

	pb "github.com/HuanLiu-hotstar/proto/authority"
	pblt "github.com/HuanLiu-hotstar/proto/ratelimit"

	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	zipkinmw "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	tracer        *openzipkin.Tracer
	pcAddr        = "http://127.0.0.1:8083/pc"
	authAddr      = "http://127.0.0.1:8085/auth"
	auth2Addr     = "http://127.0.0.1:8085/auth2"
	rateAddr      = "http://127.0.0.1:8082/ratelimit"
	authGrpcAddr  = "127.0.0.1:50055"
	limitGrpcAddr = "127.0.0.1:50052"
)

func main() {

	localEndpoint, err := openzipkin.NewEndpoint("Gateway", "192.168.1.61:8082")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	// ratio := 0.001
	// seed := int64(1000)
	// sample, err := openzipkin.NewBoundarySampler(ratio, seed)
	// if err != nil {
	// 	log.Fatal("err:%s", err)
	// }
	//tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint), openzipkin.WithSampler(sample))
	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		panic(fmt.Sprintf("err:%s", err))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", list)
	mux.HandleFunc("/playback", playback)
	mux.HandleFunc("/client", client)
	// spanName := "root"
	middler := zipkinmw.NewServerMiddleware(
		tracer,
		zipkinmw.SpanName("gateway"),
		zipkinmw.TagResponseSize(true),
	)

	h := middler(mux)
	port := ":8080"
	log.Printf("Server listening! %s ...", port)
	log.Fatal(http.ListenAndServe(port, h))

}

type Req struct {
	ID string `json:"id"`
}
type Resp struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func getbody(r *http.Request) (*Req, error) {
	req := &Req{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return req, err
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, req); err != nil {
		return req, err
	}
	return req, nil

}
func Write(w http.ResponseWriter, code int32, msg string) {
	resp := Resp{
		Code: code,
		Msg:  msg,
	}
	bye, _ := json.Marshal(resp)
	w.Write(bye)
}
func playback(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	req, err := getbody(r)
	if err != nil {
		Write(w, 500, fmt.Sprintf("err:%s", err))
		return
	}
	span.Tag("ID", req.ID)
	log.Printf("req:%+v", req)
	callAuthGrpc(ctx, authGrpcAddr, "auth-client")
	callRateLimitGrpc(ctx, limitGrpcAddr, "ratelimit-client")
	// //call auth
	// if err := callauth(ctx, authAddr, "auth-client", req); err != nil {
	// 	Write(w, -101, fmt.Sprintf("err:%s", err))
	// 	return
	// }
	// if err := callauth(ctx, auth2Addr, "auth2-client", req); err != nil {
	// 	Write(w, -101, fmt.Sprintf("err:%s", err))
	// 	return
	// }
	// //call rate limit
	// if err := callratelimit(ctx, rateAddr, "ratelimit-client", req); err != nil {
	// 	Write(w, -102, fmt.Sprintf("err:%s", err))
	// 	return
	// }
	//call pc
	if err := callpc(ctx, pcAddr, "pc-client", req); err != nil {
		Write(w, -103, fmt.Sprintf("err:%s", err))
		return
	}

	Write(w, 0, "success")
}

func callauth(c context.Context, addr, clientName string, reqData *Req) error {
	var err error
	callpc(c, addr, clientName, reqData)
	return err
}
func callratelimit(c context.Context, addr, clientName string, reqData *Req) error {
	var err error
	callpc(c, addr, clientName, reqData)
	return err
}
func callpc(c context.Context, addr, clientName string, reqData *Req) error {
	// var err error
	// /return err
	log.Printf("addr:%s,name:%s", addr, clientName)
	// data := *reqData

	// initiate a call to some_func
	url := addr
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("unable to create http request: %+v\n", err)
		return err
	}
	req = req.WithContext(c) // use parent context
	client, err := zipkinmw.NewClient(tracer, zipkinmw.ClientTrace(true), zipkinmw.ClientTags(map[string]string{"type:": "from-raw-http-client"}))
	if err != nil {
		log.Printf("err NewClient %s", err)
		return err
	}
	//res, err := client.DoWithAppSpan(req, "other_svr_client")
	res, err := client.DoWithAppSpan(req, clientName)
	if err != nil {
		log.Printf("unable to do http request: %+v\n", err)
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("err read body %s", err)
		return err
	}

	log.Printf("pc result %s", body)
	return nil
}

func client(w http.ResponseWriter, r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	x := rand.Intn(10) + 3
	time.Sleep(time.Duration(x) * time.Millisecond)
	w.Write([]byte("hello client"))
}
func list(w http.ResponseWriter, r *http.Request) {
}

func callAuthGrpc(ctx context.Context, address, name string) {
	log.Printf("addr:%s name:%s", address, name)
	span, nctx := tracer.StartSpanFromContext(ctx, name)
	defer span.Finish()
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(tracer)))
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewAuthClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(nctx, time.Second)
	defer cancel()
	r, err := c.Auth(ctx, &pb.AuthRequest{Data: name})
	if err != nil {
		log.Printf("could not greet: %v", err)
		return
	}
	log.Printf("auth: %s", r.GetMessage())
}

func callRateLimitGrpc(ctx context.Context, address, name string) {
	log.Printf("addr:%s name:%s", address, name)
	span, nctx := tracer.StartSpanFromContext(ctx, name)
	defer span.Finish()
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(tracer)))
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pblt.NewRateLimitClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(nctx, time.Second)
	defer cancel()
	r, err := c.Limit(ctx, &pblt.RateLimitRequest{Data: name})
	if err != nil {
		log.Printf("could not greet: %v", err)
		return
	}
	log.Printf("auth: %s", r.GetMessage())
}
