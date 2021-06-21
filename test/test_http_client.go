package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "chat_infra"
	subsystem = "http_client"
)

var (
	reqLabels = []string{"status", "endpoint", "method"}

	ReqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "request_count",
			Help:      "Total number of http requests.",
		}, reqLabels,
	)

	ReqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "request_duration",
			Help:      "HTTP request latencies in seconds.",
		}, reqLabels,
	)
)

func Init() {
	prometheus.MustRegister(
		ReqDuration,
		ReqCount,
	)
}

type ClientOpt func(*Client)
type Func func(namespace, subname string, code int, duration time.Time)
type Client struct {
	// http.Client
	client  http.Client
	url     string
	timeout time.Duration
	Method  string
	monitor Func

	namespace string
	subname   string // is url
}

func WithTimeout(timeout time.Duration) ClientOpt {
	return func(c *Client) {
		c.timeout = timeout
	}
}
func WithMethod(method string) ClientOpt {
	return func(c *Client) {
		c.Method = method
	}
}
func WithUrl(url string) ClientOpt {
	return func(c *Client) {
		c.url = url
	}
}
func WithFunc(f Func) ClientOpt {
	return func(c *Client) {
		c.monitor = f
	}
}
func DefaultFunc(url, method string, code int, start time.Time) {
	// fmt.Printf("namespace:%s,subname:%s,code:%d,cost:%f\n", namespace, subname, code, float64(duration.Milliseconds())/1000)

	status := fmt.Sprintf("%d", code)
	endpoint := url
	// method := method
	lvs := []string{status, endpoint, method}
	ReqCount.WithLabelValues(lvs...).Inc()
	ReqDuration.WithLabelValues(lvs...).Observe(time.Since(start).Seconds())
}

// NewClient default timeout=5s default:Method=POST
func NewClient(opts ...ClientOpt) *Client {
	c := &Client{
		timeout: time.Second * 5,
		Method:  "POST",
		monitor: DefaultFunc,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
func getpath(url string) string {
	a := strings.Split(url, "/")
	z := []string{}
	if strings.HasPrefix(url, "http") && len(a) > 4 {
		z = a[4:]
	}
	return strings.Join(z, "/")
}
func (c *Client) Do(bye []byte) ([]byte, error) {
	start := time.Now()
	code := 0
	defer func() {
		c.monitor(c.url, c.Method, code, start)
	}()
	req, err := http.NewRequest(c.Method, c.url, bytes.NewBuffer((bye)))
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		code = -resp.StatusCode
		return nil, fmt.Errorf("http resp_code:%d", resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, err
}

func main() {
	opts := []ClientOpt{
		WithMethod("GET"),
		WithUrl("http://localhost:9001/prometheus"),
	}
	Init()
	c := NewClient(opts...)
	// fmt.Println(c)
	resp, err := c.Do(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(resp))

}
