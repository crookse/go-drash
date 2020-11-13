package http

import (
	"reflect"

	"github.com/drashland/godrash"
	"github.com/valyala/fasthttp"
)

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - VARIABLE DECLARATIONS ////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

var resources = []*Resource{}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - MEMBERS EXPORTED /////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Server struct {}

// Add resources to the server
func (s Server) AddResources(resourcesArr ... func() *Resource) {
	for i := range resourcesArr {
		resources = append(resources, resourcesArr[i]())
	}
}

// Handle all http requests with this function
func (s Server) HandleRequest(ctx *fasthttp.RequestCtx) {
	uri := string(ctx.Path())
	method := string(ctx.Method())

	resource, err := s.findResource(uri)

	if err != nil {
		ctx.SetBody([]byte(err.Message))
		return
	}
	
	_, err = callHttpMethod(resource, method, ctx)

	if err != nil {
		ctx.SetBody([]byte(err.Message))
		return
	}
}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - MEMBERS NOT EXPORTED /////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// This code was taken from the following article:
// medium.com/@vicky.kurniawan/go-call-a-function-from-string-name-30b41dcb9e12
func callHttpMethod(
	resource *Resource,
	funcName string,
	ctx interface{},
) (result interface{}, err *godrash.errors.HttpError) {
	f := reflect.ValueOf(resource.Methods[funcName])

	// Is the method defined?
	if !f.IsValid() {
		var err = new(godrash.errors.HttpError)
		err.Code = 405
		err.Message = "Method Not Allowed"
		return nil, err
	}

	args := []reflect.Value{reflect.ValueOf(ctx)}

	var res []reflect.Value
	res = f.Call(args)

	if len(res) > 0 {
		result = res[0].Interface()
	}

	return result, nil
}

// Find the resource in question given the URI
func (s Server) findResource(uri string) (*Resource, *godrash.errors.HttpError) {
	if uri == "/" {
		return resources[0], nil
	}

	return nil, s.handleError(404, "Not Found")
}

// Handle server errors -- making sure to send HTTP error responses. HTTP error
// responses should always have a code and a message.
func (s Server) handleError(code int, message string) (*godrash.errors.HttpError) {
	e := new(godrash.errors.HttpError)
	e.Code = code
	e.Message = message
	return e
}
