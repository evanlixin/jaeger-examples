package main

import (
	"io"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/evanlixin/jaeger-examples/pkg/tracing"
	"github.com/opentracing/opentracing-go"
)

func main() {
	tracer, closer := tracing.Init("restful-server")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// install a global (=DefaultContainer) filter (processed before any webservice in the DefaultContainer)
	// to provide  OpenTracing instrument
	restful.Filter(NewOTFilter(tracer))

	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))

	restful.Add(ws)
	http.ListenAndServe(":8083", nil)
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}