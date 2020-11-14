package http

import (
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/drashland/go-drash/errors"
	"github.com/valyala/fasthttp"
)

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - VARIABLE DECLARATIONS ////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

var resources = []Resource{}
var responseContentType = "application/json"

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - MEMBERS EXPORTED /////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type ServerOptions struct {
	Hostname string
	Port     int
}

type Server struct {
	Resources []func() Resource
	ResponseContentType string
}

// Handle all HTTP requests with this function
func (s Server) HandleRequest(ctx *fasthttp.RequestCtx) {

	// Make a "Drash" request -- basically a wrapper around fasthttp's request
	res := &Response{
		ContentType: responseContentType,
	}
	request := Request{
		Ctx: ctx,
		Response: res,
	}

	uri := string(request.Ctx.Path())
	method := string(request.Ctx.Method())

	// Find the resource that matches the request's URI the best
	resource, err := findResource(uri)
	if err != nil {
		request.SendError(err.Code, err.Message);
		return
	}

	// Make the request
	_, err = callHttpMethod(resource, method, request)
	if err != nil {
		request.SendError(err.Code, err.Message);
		return
	}

	// Finally, send the response. The response content type, status code, and
	// body should all be set before this method is called.
	request.Send()
}

// Run the server
func (s Server) Run(o ServerOptions) {
	addResources(s.Resources)

	if s.ResponseContentType != "" {
		responseContentType = s.ResponseContentType
	}

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
		resource.ParseUris()
		resources = append(resources, resource)
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
	args ...interface{},
) (response *Response, err *errors.HttpError) {
	f := reflect.ValueOf(resource.Methods[funcName])

	// Is the method defined?
	if !f.IsValid() || f.IsNil() {
		return nil, buildError(405, "Method Not Allowed")
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	var result []reflect.Value
	result = f.Call(in)

	if len(result) > 0 {
		data := result[0].Interface().(*Response)
		return data, nil
	}

	return nil, buildError(418, "I'm a teapot")
}

// Build an HTTP error response (e.g., a 404 Not Found error response)
func buildError(code int, message string) *errors.HttpError {
	e := new(errors.HttpError)
	e.Code = code
	e.Message = message
	return e
}

// Find the best matching resource based on the request's URI. If a resource
// cannot be found, then that is a 404 error -- most likey due to a resource
// not being defined to handle the URI in question.
func findResource(uri string) (Resource, *errors.HttpError) {
	for i := range resources {
		for k := range resources[i].UrisParsed {
			pathObj := resources[i].UrisParsed[k]
			re := regexp.MustCompile(pathObj.RegexPath)
			matches := re.FindAllString(uri, -1)
			if len(matches) > 0 {
				return resources[i], nil
			}
		}
	}

	return Resource{}, buildError(404, "Not Found")
}
