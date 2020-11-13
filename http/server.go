package http

import (
	"fmt"
	"reflect"
	"log"

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

type Server struct {}

// Add resources to the server
func (s Server) AddResources(resourcesArr ... func() Resource) {
	for i := range resourcesArr {
		resource := resourcesArr[i]()
		resources = append(resources, resource)
	}
	for i := range resources {
		resource := resources[i]
		resource.ParseUris()
		fmt.Println(resource)
	}
}

// Handle all http requests with this function
func (s Server) HandleRequest(ctx *fasthttp.RequestCtx) {

	request := Request{Ctx: ctx}

	uri := string(request.Ctx.Path())
	method := string(request.Ctx.Method())

	resource, err := s.findResource(uri)

	if err != nil {
		ctx.SetBody([]byte(err.Message))
		return
	}
	
	response := new(Response)
	resourceResponse, err := callHttpMethod(resource, method, request, *response)

	if err != nil {
		ctx.SetBody([]byte(err.Message))
		return
	}

	ctx.SetBody([]byte(resourceResponse.Body))
}

// Run the server
func (s Server) Run(addr string) {
	err := fasthttp.ListenAndServe(addr, s.HandleRequest)

	if err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - MEMBERS NOT EXPORTED /////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// This code was taken from the following article:
// medium.com/@vicky.kurniawan/go-call-a-function-from-string-name-30b41dcb9e12
func callHttpMethod(
	resource Resource,
	funcName string,
	params ... interface{},
) (response Response, err *errors.HttpError) {
	f := reflect.ValueOf(resource.Methods[funcName])

	// Is the method defined?
	if !f.IsValid() {
		var err = new(errors.HttpError)
		err.Code = 405
		err.Message = "Method Not Allowed"
		var r = Response{}
		return r, err
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

	err = new(errors.HttpError)
	err.Code = 418
	err.Message = "I'm a teapot"

	r := Response{}
	return r, err
}

// Find the resource in question given the URI
func (s Server) findResource(uri string) (Resource, *errors.HttpError) {
	if uri == "/" {
		return resources[0], nil
	}

	return Resource{}, s.handleError(404, "Not Found")
}

// Handle server errors -- making sure to send HTTP error responses. HTTP error
// responses should always have a code and a message.
func (s Server) handleError(code int, message string) (*errors.HttpError) {
	e := new(errors.HttpError)
	e.Code = code
	e.Message = message
	return e
}
