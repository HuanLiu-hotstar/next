# Monitoring,Tracing

## What's Tracing and Monitoring?

## Tracing function

## Tracing Concept

## Architecture of Tracing System


## Implementation of Tracing System

### http & rpc implementation

- http: Write the Tracing info in http-header 
  
```go
// delegate the http.Handler with tracing info
type handler struct {
  tracer *Tracer // tracing the request and response
  handle http.Handler // the true http-handler 
  //   ... 
}


// InterceptResponse dalegate the http-response with tracing  info
type InterceptResponse struct {
  w http.ResponseWriter 
  statusCode int 
}

// ServeHTTP implements http.Handler.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  //tracing parameter init from http.Request

  defer func(){
  //   reporte tracing 
  }()
  h.handle.ServeHTTP(h.wrap(w),h.wrap(r)) // true http-call but with the wrap ResponseWriter and http.Request.
}

func main() {
	h := handler{
		tracer : newTracer(),
		handler : func(w http.ResponseWriter, r *http.Request){
			w.Write("hello")
		}
	}
	http.HandleFunc("/uri",h)
}

```
- rpc : Write the tracing with the intercept handler 

- Open-Tracing standard
- zipkin: twitter
- jaeger : uber
- Skywalking:Huawei
- CAT: meituan