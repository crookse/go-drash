package http

import (
	"fmt"
	"log"
	"reflect"

	"github.com/drashland/go-drash/services"
	"github.com/valyala/fasthttp"
)

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - VARIABLE DECLARATIONS ////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

var _resources = []Resource{}

var _responseContentType = "application/json"

var _services = map[string]interface{}{
	"ResourceIndexService": services.IndexService{
		Cache:       map[string][]services.IndexServiceSearchResult{},
		LookupTable: map[int]interface{}{},
		Index:       map[string][]int{},
	},
}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - STRUCTS //////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Server struct {
	Resources           []func() Resource
	ResponseContentType string
}

type ServerOptions struct {
	Hostname string
	Port     int
}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - METHODS - EXPORTED ///////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// This method handles all incoming requests.
func (s Server) HandleIncomingRequest(ctx *fasthttp.RequestCtx) {

	// Make a "Drash" request -- basically a wrapper around fasthttp's request
	request := s.buildRequest(ctx)

	// Find the resource that matches the request's URI the best
	resource := s.findResource(request.Uri)
	if resource == nil {
		request.SendError(404, "Not Found")
		return
	}

	// If the HTTP method does not exist on the resource, then that method is
	// not allowed
	httpMethod := resource.Methods[request.Method]
	if reflect.ValueOf(httpMethod).IsNil() {
		request.SendError(405, "Method Not Allowed")
		return
	}

	// Make the request
	_ = httpMethod.(func(r *Request) Response)(request)

	// Finally, send the response. The response content type, status code, and
	// body should all be set before this method is called.
	request.Send()
}

// This method runs the server at the specified hostname and port given in the
// ServerOptions.
func (s Server) Run(o ServerOptions) {
	s.buildServer()

	if s.ResponseContentType != "" {
		_responseContentType = s.ResponseContentType
	}

	address := fmt.Sprintf("%s:%d", o.Hostname, o.Port)
	err := fasthttp.ListenAndServe(address, s.HandleIncomingRequest)

	if err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - METHODS - NOT EXPORTED ///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// This method builds a request. It is essentially a "Drash" request wrapped
// around fasthttp's request.
func (s Server) buildRequest(ctx *fasthttp.RequestCtx) *Request {
	return &Request{
		Ctx: ctx,
		Method: string(ctx.Method()),
		Response: Response{
			ContentType: _responseContentType,
			StatusCode:  200,
		},
		Uri: string(ctx.Path()),
	}
}

// This methods builds and returns a map of HTTP methods that the resource has
// or does not have.
func (s Server) buildResourceHttpMethodsMap(
	resource Resource,
) map[string]interface{} {
	return map[string]interface{}{
		"CONNECT": resource.CONNECT,
		"DELETE":  resource.DELETE,
		"GET":     resource.GET,
		"HEAD":    resource.HEAD,
		"OPTIONS": resource.OPTIONS,
		"PATCH":   resource.PATCH,
		"POST":    resource.POST,
		"PUT":     resource.PUT,
		"TRACE":   resource.TRACE,
	}
}

// This builds the server. Run anything that needs to process during compile
// time in this function.
func (s Server) buildServer() {
	s.buildServerResourcesTable()
}

// This builds the resources table. During runtime, the resources table is used
// to match a request's URI to a resource. If a resource is found, the resource
// takes responsibility of handling the request. If a resource is not found,
// then a 404 error is thrown.
func (s Server) buildServerResourcesTable() {
	for i1 := range s.Resources {
		resource := s.Resources[i1]()

		resource.Methods = s.buildResourceHttpMethodsMap(resource)

		resource.ParseUris()

		for i2 := range resource.UrisParsed {
			s.indexResource(resource, i2)
		}
	}
}

// Find the best matching resource based on the request's URI. If a resource
// cannot be found, then that is a 404 error -- most likey due to a resource
// not being defined to handle the URI in question.
func (s *Server) findResource(uri string) *Resource {

	var results = _services["ResourceIndexService"].(services.IndexService).Search(uri)

	if len(results) > 0 {
		return results[0].Item.(*Resource)
	}

	return nil
}

// This method indexes the given resource.
func (s Server) indexResource(resource Resource, index int) {
	_services["ResourceIndexService"].(services.IndexService).AddItem(
		[]string{resource.UrisParsed[index].RegexUri},
		&resource,
	)
}
