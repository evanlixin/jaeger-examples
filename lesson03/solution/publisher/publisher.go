package main

import (
	"log"
	"net/http"

	"github.com/evanlixin/jaeger-examples/pkg/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func main() {
	tracer, closer := tracing.Init("publisher")
	defer closer.Close()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
		}
		span := tracer.StartSpan("publish", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloStr := r.FormValue("helloStr")
		println(helloStr)
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
}
