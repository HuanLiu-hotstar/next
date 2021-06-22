package zipkinClientHttp

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	// "go.opencensus.io/plugin/ochttp"
	// "go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

func TestClient(t *testing.T) {
	Init(WithLocalServerName("go-test"), WithServiceName("remote_server"), WithLocalAddr("192.168.0.2"))
	defer Destroy()

	// initiate a call to some_func
	addrServ := "127.0.0.1:8080"
	url := fmt.Sprintf("http://%s/list", addrServ)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("unable to create http request: %+v\n", err)
	}

	client, err := NewClient(WithClient(&http.Client{Timeout: time.Second * 5}))
	if err != nil {
		log.Fatalf("err NewClient %s", err)
	}
	res, err := client.DoWithAppSpan(req, "test-client-list")
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("err read body %s", err)
	}

	// Output:
	log.Printf("result %s", body)
}

func TestServer(t *testing.T) {
	Init(WithLocalServerName("go-test-svr"), WithServiceName("local_svr_name"), WithLocalAddr("192.168.0.2"))
	defer Destroy()
	rgin := gin.Default()
	// 第三步: 添加一个 middleWare, 为每一个请求添加span
	// r.Use(zipkinMiddle)
	rgin.GET("/",
		func(c *gin.Context) {
			time.Sleep(500 * time.Millisecond)
			c.JSON(200, gin.H{"code": 200, "msg": "OK-root"})
		})
	rgin.GET("/list",
		func(c *gin.Context) {
			// time.Sleep(500 * time.Millisecond)
			dohttp(c)
			c.JSON(200, gin.H{"code": 200, "msg": "OK-list"})
		})
	RunGinServer(rgin, 8080)
	// h := &ochttp.Handler{Handler: rgin}
	// if err := view.Register(ochttp.DefaultServerViews...); err != nil {
	// 	log.Fatal("Failed to register ochttp.DefaultServerViews")
	// }
	// http.Handle("/", h)
	// port := ":8080"
	// log.Printf("listen:%s", port)
	// if err := http.ListenAndServe(port, nil); err != nil {
	// 	log.Fatalf("err:%s", err)
	// }
}
func dohttp(c *gin.Context) {
	//span := r.zkTracer.StartSpan("dohttp")
	_, span := trace.StartSpan(c.Request.Context(), "dohttp")
	defer span.End()
	doOther(c)
	time.Sleep(400 * time.Millisecond)
}
func doOther(c context.Context) {
	_, span := trace.StartSpan(c, "do_other")
	defer span.End()
	log.Println("do_other")
	time.Sleep(100 * time.Millisecond)

}
