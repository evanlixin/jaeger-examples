package main

import (
	"log"
	"net/http"

	"github.com/evanlixin/jaeger-examples/pkg/tracing"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

const DefaultComponentName = "go-restful"
const DefaultOperationName = "go-restful client"

func main() {
	tracer, closer := tracing.Init("restful-client")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// 2) Demonstrate nethttp client-side OpenTracing instrumentation works
	client := &http.Client{Transport: &nethttp.Transport{}}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8083/hello", nil)
	if err != nil {
		log.Fatal(err)
	}

	req, ht := nethttp.TraceRequest(tracer, req,
		nethttp.OperationName(DefaultOperationName), nethttp.ComponentName(DefaultComponentName))

	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	ht.Finish()

}