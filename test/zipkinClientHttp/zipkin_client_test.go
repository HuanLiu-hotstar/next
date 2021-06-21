package zipkinClientHttp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
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
