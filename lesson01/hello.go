package main

import (
	"fmt"
	"os"

	"github.com/evanlixin/jaeger-examples/pkg/tracing"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}

	/*
	However, if we run this program, we will see no difference, and no traces in the tracing UI.
	That's because the function opentracing.GlobalTracer() returns a no-op tracer by default.
	 */
	// tracer := opentracing.GlobalTracer()

	tracer, closer := tracing.Init("hello-world")
	defer closer.Close()

	helloTo := os.Args[1]

	span := tracer.StartSpan("say-hello")
	span.SetTag("hello-to", helloTo)

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	println(helloStr)
	span.LogKV("event", "println")

	span.Finish()
}
