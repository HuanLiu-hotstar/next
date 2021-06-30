# demo of tracing

## usage of zipkin

## http server

- http server
- grpc client

```go
var tracer *openzipkin.Tracer
func newTracer()(*openzipkin.Tracer,func() ){
	// 1. create local Endpoint 
	localEndpoint, _ := openzipkin.NewEndpoint("Gateway", "192.168.1.61:8082")
	// 2. create reporter to report tracing data
	reportEndpoint := "http://localhost:9411/api/v2/spans"
	reporter := zipkinHTTP.NewReporter(reportEndpoint)
	
	// 3. create a tracer to tracing data 
	tracer, _ = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	return tracer,func(){ reporter.Close()}
}
func main() {
	
	tracer, close := newTracer()
	defer close()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/playback", playback)
	mux.HandleFunc("/client", client)
	// 4. create middleware to extract injection data to conetxt from http Header 
	middler := zipkinmw.NewServerMiddleware(tracer,zipkinmw.TagResponseSize(true))
	h := middler(mux)
	log.Fatal(http.ListenAndServe(port, h))
}
func playback(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish() // finish span and report to zipkin server 
	span.Tag("ID", req.ID) // set tag
	callGrpc(ctx,"127.0.0.1:50553","auth-client")
	w.Write([]byte("success"))
}

func callGrpc(ctx context.Context, address, name string) {
	//this is a sub-span from playback 
	span, nctx := tracer.StartSpanFromContext(ctx, name)
	defer span.Finish()

	// Set up a connection to the server.
	conn, _:= grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(tracer)))

	defer conn.Close()
	c := pb.NewAuthClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(nctx, time.Second)
	defer cancel()
	r, _:= c.Auth(ctx, &pb.AuthRequest{Data: name})
}

```


## grpc server 

- grpc server 

```go 

func main(){
	// init an tracer same as http server 
	tracer, f = newTracer()
	defer f()
	lis, err := net.Listen("tcp", port)
	//create middleware for grpc server 
	middleware := grpc.StatsHandler(zipkingrpc.NewServerHandler(tracer))

	s := grpc.NewServer(middleware)
	pb.RegisterRateLimitServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

```

##  client to call other http server

```go
var tracer *openzipkin.Tracer
func main() {
	call_http(c,addr,clientName,nil)
}
// init a tracer same as http sever 
func call_http(c context.Context, addr, clientName string, reqData *Req) error {

	// create a http request 
	req, _ := http.NewRequest("GET", addr , nil)

	// req with parent context to inject header 
	req = req.WithContext(c) // use parent context

	// create a client with zipkin lib
	client, _ := zipkinmw.NewClient(tracer, zipkinmw.ClientTrace(true) )

	// do the http request
	res, _ := client.DoWithAppSpan(req, clientName)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	return nil
}

```