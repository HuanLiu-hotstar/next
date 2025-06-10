# Tracing 介绍

- Tracing 概念
- Tracing 标准 - opentracing
- Tracing 实现
- Tracing 使用方法
- Tracing 集成在Gateway中

## Tracing 概念

- 目标：追踪请求执行踪迹，特别是在分布式系统中
- 追踪请求，调用链
- 为什么需要tracing。一个请求案例，如果请求RequestX出现变慢，如何排查？
- ![image-20210624112353755](/Users/liuhuan/Library/Application Support/typora-user-images/image-20210624112353755.png)

## Tracing 标准 - OpenTracing

- Trace: 请求的踪迹汇总
- span: 请求的上下文，主要包括traceID, SpanID，发生时间，结束时间。可以认为是一次函数调用
- Annotation: 请求执行过程中的发生的事件，用于具体的事件记录，例如记录函数开始时间
- Tag：与请求相关的标记内容，相当于日志记录。多用于debug的具体case搜索，例如记录每个请求的业务ID。
- reporter：代表传输组件，一般是单独实现，集成在span内

### span 举例

- span relationship 
- ![image-20210624194003618](/Users/liuhuan/Library/Application Support/typora-user-images/image-20210624194003618.png)
- one span example 

![image-20210624193924703](/Users/liuhuan/Library/Application Support/typora-user-images/image-20210624193924703.png)

![image-20210624162821705](/Users/liuhuan/Library/Application Support/typora-user-images/image-20210624162821705.png)

## Tracing系统的基本架构

- Tracing结构图
- Tracing Client: 上报trace
- Tracing Server : 收集trace数据
- Storage：存储trace数据
- User API : 用户界面
- Tracing具体实现以lib包的形式提供给应用程序



![image-20210624115547871](/Users/liuhuan/Library/Application Support/typora-user-images/image-20210624115547871.png)



## tracing 实现

#### http实现

- http代理流程

![image-20210624150211529](/Users/liuhuan/Library/Application Support/typora-user-images/image-20210624150211529.png)

- 使用http对Header的代理示例

```go
┌─────────────┐ ┌───────────────────────┐  ┌─────────────┐  ┌──────────────────┐
│ User Code   │ │ Trace Instrumentation │  │ Http Client │  │ Zipkin Collector │
└─────────────┘ └───────────────────────┘  └─────────────┘  └──────────────────┘
       │                 │                         │                 │
           ┌─────────┐
       │ ──┤GET /foo ├─▶ │ ────┐                   │                 │
           └─────────┘         │ record tags
       │                 │ ◀───┘                   │                 │
                           ────┐
       │                 │     │ add trace headers │                 │
                           ◀───┘
       │                 │ ────┐                   │                 │
                               │ record timestamp
       │                 │ ◀───┘                   │                 │
                             ┌─────────────────┐
       │                 │ ──┤GET /foo         ├─▶ │                 │
                             │X-B3-TraceId: aa │     ────┐
       │                 │   │X-B3-SpanId: 6b  │   │     │           │
                             └─────────────────┘         │ invoke
       │                 │                         │     │ request   │
                                                         │
       │                 │                         │     │           │
                                 ┌────────┐          ◀───┘
       │                 │ ◀─────┤200 OK  ├─────── │                 │
                           ────┐ └────────┘
       │                 │     │ record duration   │                 │
            ┌────────┐     ◀───┘
       │ ◀──┤200 OK  ├── │                         │                 │
            └────────┘       ┌────────────────────────────────┐
       │                 │ ──┤ asynchronously report span     ├────▶ │
                             │                                │
                             │{                               │
                             │  "traceId": "aa",              │
                             │  "id": "6b",                   │
                             │  "name": "get",                │
                             │  "timestamp": 1483945573944000,│
                             │  "duration": 386000,           │
                             │  "annotations": [              │
                             │--snip--                        │
                             └────────────────────────────────┘
```

- 具体header注入的标记

```go
// Default B3 Header keys
const (
	TraceID      = "x-b3-traceid"
	SpanID       = "x-b3-spanid"
	ParentSpanID = "x-b3-parentspanid"
	Sampled      = "x-b3-sampled"
	Flags        = "x-b3-flags"
	Context      = "b3"
)

```

- 其他内容，性能监控，采样频率等

  

  

#### grpc实现

- Metadata with intercepter 

## Tracing 使用

### zipkin go的使用



```go
// create reporter to localhost zipkin-server 
reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
defer reporter.Close()
//create tracer 
tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))

```



### 客户端

- http.Client的代理

```go
  localEndpoint, err := openzipkin.NewEndpoint("local-client", "192.168.1.61:8082")

	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	//exporter := zipkin.NewExporter(reporter, localEndpoint)
	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	
	// create  zipkin traced http client
	client, err := zipkinmw.NewClient(tracer, zipkinmw.ClientTrace(true), zipkinmw.ClientTags(map[string]string{"type:": "from-raw-http-client"}))

	// initiate a call to list
	url := "http://localhost:8080/list"
	req, err := http.NewRequest("GET", url, nil)

	res, err := client.DoWithAppSpan(req, "raw-http-client")

	defer res.Body.Close() // close will send the span data
	body, err := ioutil.ReadAll(res.Body)
	log.Printf("%s", body)
```



### 服务端

- raw http 的使用

```go

func main() {

	localEndpoint, err := openzipkin.NewEndpoint("golangsvc", "192.168.1.61:8082")

	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	//exporter := zipkin.NewExporter(reporter, localEndpoint)
	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))

	mux := http.NewServeMux()
	mux.HandleFunc("/list", list)
	spanName := "root"
  // use the middleware around the true http.Handler
	middler := zipkinmw.NewServerMiddleware(
		tracer,
		zipkinmw.SpanName(spanName),
		zipkinmw.TagResponseSize(true),
	)

	h := middler(mux)
	port := ":8080"
	log.Printf("Server listening! %s ...", port)
	log.Fatal(http.ListenAndServe(port, h))

}

func list(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	span, _ := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	database(r)
	serviceb(r)
	res := strings.Repeat("o", rand.Intn(100)+1)
	time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)
	w.Write([]byte("Hello, w" + res + "rld!"))
}

func database(r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), "database")
	defer span.Finish()
	cache(r)
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
	span, pc := tracer.StartSpanFromContext(r.Context(), "serviceb")
	defer span.Finish()
	time.Sleep(time.Duration(rand.Intn(100)+100) * time.Millisecond)
	servicec(pc) // servicec is childof serviceb
	span.Annotate(time.Now(), "endtime")
}

//func servicec(r *http.Request) {
func servicec(c context.Context) {
	span, _ := tracer.StartSpanFromContext(c, "servicec")
	defer span.Finish()
	time.Sleep(time.Duration(rand.Intn(700)+100) * time.Millisecond)
	span.Tag("servicec", "C") // set tags for search servicec
}
```



- 服务端追踪展示效果

![image-20210624140808047](/Users/liuhuan/Library/Application Support/typora-user-images/image-20210624140808047.png)

- zipkin gin 的使用

  ```go
  
  //use gin.Engine as a http.Handler with middleware
  // init tracer omit 
  middler := zipkinmw.NewServerMiddleware(
  		tracer,
  		zipkinmw.SpanName(spanName),
  		zipkinmw.TagResponseSize(true),
  	)
  	rgin := gin.Default()
  	h := middler(rgin)
  	port := ":8080"
  	log.Fatal(http.ListenAndServe(port, h))
  
  ```

  



## Tracing在Gateway中集成

- 以插件的形式 提供给Gateway
- 在gateway里提供tracing服务的地址，Gateway以tracing server client形式存在

![image-20210624120222115](/Users/liuhuan/Library/Application Support/typora-user-images/image-20210624120222115.png)

# 参考

- Dapper paper from google
- Zipkin 
- Ambassador



## 附录



- zipkin server的搭建 java8以上

```shell
curl -sSL https://zipkin.io/quickstart.sh | bash -s
java -jar zipkin.jar
// open default url localhost:9411/zipkin
```

- http server端代码 [参见这里](https://github.com/HuanLiu-hotstar/next/blob/main/test/zipkin/test_zipkin_raw_http.go)
- http client 端代码 [参见这里](https://github.com/HuanLiu-hotstar/next/blob/main/test/zipkin/test_zipkin_raw_http_client.go)

