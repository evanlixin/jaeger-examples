package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/evanlixin/jaeger-examples/pkg/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

func main() {
	tracer, closer := tracing.Init("formatter")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	http.HandleFunc("/format", func(w http.ResponseWriter, r *http.Request) {
		// map[Accept-Encoding:[gzip] Uber-Trace-Id:[218e63028becb349:7199bb351a797c0c:218e63028becb349:1] User-Agent:[Go-http-client/1.1]]
		fmt.Printf("%v\n", r.Header)
		spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := opentracing.GlobalTracer().StartSpan("format", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloTo := r.FormValue("helloTo")
		helloStr := fmt.Sprintf("Hello, %s!", helloTo)
		span.LogFields(
			otlog.String("event", "string-format"),
			otlog.String("value", helloStr),
		)
		w.Write([]byte(helloStr))
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
