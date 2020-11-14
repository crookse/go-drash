package http

import (
	"fmt"
	"log"
	"reflect"

	"github.com/drashland/go-drash/errors"
	"github.com/valyala/fasthttp"
)

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - VARIABLE DECLARATIONS ////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

var resources = []Resource{}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - MEMBERS EXPORTED /////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Server struct {
	Resources []func() Resource
}

type HttpOptions struct {
	Hostname string
	Port     int
}

// Handle all HTTP requests with this function
func (s Server) HandleRequest(ctx *fasthttp.RequestCtx) {

	// Make a "Drash" request -- basically a wrapper around fasthttp's request
	request := Request{Ctx: ctx}

	uri := string(request.Ctx.Path())
	method := string(request.Ctx.Method())

	// Find the resource that matches the request's URI the best
	resource, err := findResource(uri)
	if err != nil {
		request.Ctx.SetBody([]byte(err.Message))
		return
	}

	// Make the request
	response, err := callHttpMethod(resource, method, request)
	if err != nil {
		request.Ctx.SetBody([]byte(err.Message))
		return
	}

	// Send the response
	request.Ctx.SetBody([]byte(response.Body))
}

// Run the server
func (s *Server) Run(o HttpOptions) {
	addResources(s.Resources)

	address := fmt.Sprintf("%s:%d", o.Hostname, o.Port)
	err := fasthttp.ListenAndServe(address, s.HandleRequest)

	if err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - MEMBERS NOT EXPORTED /////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// Add resources to the server
func addResources(resourcesArr []func() Resource) {

	for i := range resourcesArr {
		resource := resourcesArr[i]()
		resource.Methods = map[string]interface{}{
			"GET":    resource.GET,
			"POST":   resource.POST,
			"PUT":    resource.PUT,
			"DELETE": resource.DELETE,
		}
		resources = append(resources, resource)
	}

	for i := range resources {
		resource := resources[i]
		resource.ParseUris()
	}
}

// This code was taken from the following article:
// medium.com/@vicky.kurniawan/go-call-a-function-from-string-name-30b41dcb9e12
//
// This code is used to allow "indexing" of a resource's HTTP methods. Without
// this code, we would not be be able to make calls like the following:
//
//     resource[request.Method()].
//
// This is similar to how deno-drash makes its HTTP calls.
func callHttpMethod(
	resource Resource,
	funcName string,
	params ...interface{},
) (response Response, err *errors.HttpError) {
	f := reflect.ValueOf(resource.Methods[funcName])

	// Is the method defined?
	if !f.IsValid() || f.IsNil() {
		return Response{}, errorResponse(405, "Method Not Allowed")
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	var result []reflect.Value
	result = f.Call(in)

	if len(result) > 0 {
		data := result[0].Interface().(Response)
		return data, nil
	}

	return Response{}, errorResponse(418, "I'm a teapot")
}

// Handle server errors -- making sure to send HTTP error responses. HTTP error
// responses should always have a code and a message.
func errorResponse(code int, message string) *errors.HttpError {
	e := new(errors.HttpError)
	e.Code = code
	e.Message = message
	return e
}

// Find the resource in question given the URI
func findResource(uri string) (Resource, *errors.HttpError) {
	if uri == "/" {
		return resources[0], nil
	}

	return Resource{}, errorResponse(404, "Not Found")
}
