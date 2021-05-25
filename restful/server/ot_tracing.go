package main

import (
	"log"

	"github.com/emicklei/go-restful"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const DefaultComponentName = "go-restful"
const DefaultOperationName = "go-restful client"

var (
	DefaultOperationNameFunc = func(r *restful.Request) string {
		// extract the route that the request maps to and use it as the operation name.
		return r.SelectedRoutePath()
	}
)

type filterOptions struct {
	operationNameFunc func(r *restful.Request) string
	componentName     string
}

// FilterOption controls the behavior of the Filter.
type FilterOption func(*filterOptions)

// OperationNameFunc returns a FilterOption that uses given function f
// to generate operation name for each server-side span.
func OperationNameFunc(f func(r *restful.Request) string) FilterOption {
	return func(options *filterOptions) {
		options.operationNameFunc = f
	}
}

// ComponentName returns a FilterOption that sets the component name
// name for the server-side span.
func ComponentName(componentName string) FilterOption {
	return func(options *filterOptions) {
		options.componentName = componentName
	}
}

// NewOTFilter returns a go-restful filter which add OpenTracing instrument
func NewOTFilter(tracer opentracing.Tracer, options ...FilterOption) restful.FilterFunction {
	opts := filterOptions{
		operationNameFunc: DefaultOperationNameFunc,
		componentName:     DefaultComponentName,
	}
	for _, opt := range options {
		opt(&opts)
	}

	// return a go-restful filter which add OpenTracing instrument
	// NOTE: "filter" in go-restful is similar with "middleware" mechanism in modern web framework
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		ctx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Request.Header))
		if err != nil {
			log.Println(err)
		}
		// record operation name
		sp := tracer.StartSpan(opts.operationNameFunc(req), ext.RPCServerOption(ctx))
		defer sp.Finish()

		// record HTTP method
		ext.HTTPMethod.Set(sp, req.Request.Method)
		// record HTTP url
		ext.HTTPUrl.Set(sp, req.Request.URL.String())
		// record component name
		ext.Component.Set(sp, opts.componentName)
		// record HTTP status code
		ext.HTTPStatusCode.Set(sp, uint16(resp.StatusCode()))

		chain.ProcessFilter(req, resp)
	}
}
