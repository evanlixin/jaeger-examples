package main

import (
	"context"
	"log"
	"net/http"

	"github.com/evanlixin/jaeger-examples/pkg/tracing"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const DefaultComponentName = "go-restful"
const DefaultOperationName = "go-restful client"

func main() {
	tracer, closer := tracing.Init("restful-client")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	span, _ := opentracing.StartSpanFromContext(context.Background(), "getHello")
	defer span.Finish()

	// 2) Demonstrate nethttp client-side OpenTracing instrumentation works
	client := &http.Client{Transport: &nethttp.Transport{}}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8083/hello", nil)
	if err != nil {
		log.Fatal(err)
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, req.URL.Path)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	/*
	req, ht := nethttp.TraceRequest(tracer, req,
		nethttp.OperationName(DefaultOperationName), nethttp.ComponentName(DefaultComponentName))
	*/

	_, err = client.Do(req)
	if err != nil {
		ext.LogError(span, err)
		log.Fatal(err)
	}


	//ht.Finish()

}