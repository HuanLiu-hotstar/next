package zipkinClientHttp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	// "go.opencensus.io/trace"
)

func RunGinServer(rgin *gin.Engine, port int) {

	h := &ochttp.Handler{Handler: rgin}
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		log.Fatal("Failed to register ochttp.DefaultServerViews")
	}
	http.Handle("/", h)
	listenPort := fmt.Sprintf(":%d", port)
	log.Printf("listen:%d\n", port)
	if err := http.ListenAndServe(listenPort, nil); err != nil {
		log.Fatalf("err:%s", err)
	}
}
