package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// func ServeHTTP(rw, r)
// type Func struct {
// 	f func(w http.ResponseWriter, r *http.Request)
// }
type Func func(w http.ResponseWriter, r *http.Request)

func (f Func) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}
func main() {
	router := mux.NewRouter()

	// Serving static files
	// router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world http2"))
	}
	router.PathPrefix("/hello").Handler(Func(f))

	// Prometheus endpoint
	router.Path("/prometheus").Handler(promhttp.Handler())

	fmt.Println("Serving requests on port 9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}
