package http

import (
	"fmt"
	"reflect"

	"../errors"

	"github.com/valyala/fasthttp"
)

type Server struct {}

var resources = []*Resource{}

// Add resources to the server
func (s Server) AddResource(resourcesArr ... *Resource) {
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
	
	CallHttpMethod(resource, method, ctx)
}

func (s Server) findResource(uri string) (*Resource, *errors.HttpError) {
	if uri == "/" {
		return resources[0], nil
	}

	return nil, s.handleError(404, "Not Found")
}

func (s Server) handleError(code int, message string) (*errors.HttpError) {
	e := new(errors.HttpError)
	e.Code = code
	e.Message = message
	return e
}


// This code was taken from the following article:
// medium.com/@vicky.kurniawan/go-call-a-function-from-string-name-30b41dcb9e12
func CallHttpMethod(
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
