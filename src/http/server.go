package http

import (
	"fmt"
	"reflect"

	"../errors"

	"github.com/valyala/fasthttp"
)

var resources = []*Resource{}

///////////////////////////////////////////////////////////////////////////////
// EXPORTED ///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Server struct {}

// Add resources to the server
func (s Server) AddResources(resourcesArr ... *Resource) {
	for i := range resourcesArr {
		resources = append(resources, resourcesArr[i])
	}
}

// Handle all http requests with this function
func (s Server) HandleRequest(ctx *fasthttp.RequestCtx) {
	uri := string(ctx.Path())
	method := string(ctx.Method())

	resource, err := s.findResource(uri)

	if err != nil {
		fmt.Fprintf(ctx, err.Message)
		return
	}
	
	callHttpMethod(resource, method, ctx)
}

///////////////////////////////////////////////////////////////////////////////
// NON-EXPORTED ///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// This code was taken from the following article:
// medium.com/@vicky.kurniawan/go-call-a-function-from-string-name-30b41dcb9e12
func callHttpMethod(
	resource *Resource,
	funcName string,
	params ... interface{},
) (result interface{}, err error) {
	f := reflect.ValueOf(resource.Methods[funcName])

	if len(params) != f.Type().NumIn() {
		e := new(errors.HttpError)
		e.Code = 500
		e.Message = "Internal Server Error"
		return
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	var res []reflect.Value
	res = f.Call(in)

	if len(res) > 0 {
		result = res[0].Interface()
	}

	return
}

// Find the resource in question given the URI
func (s Server) findResource(uri string) (*Resource, *errors.HttpError) {
	if uri == "/" {
		return resources[0], nil
	}

	return nil, s.handleError(404, "Not Found")
}

// Handle server errors -- making sure to send HTTP error responses. HTTP error
// responses should always have a code and a message.
func (s Server) handleError(code int, message string) (*errors.HttpError) {
	e := new(errors.HttpError)
	e.Code = code
	e.Message = message
	return e
}
