package http

import (
	"github.com/valyala/fasthttp"
)

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - STRUCTS //////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Request struct {
	Ctx      *fasthttp.RequestCtx
	Method   string
	Response Response
	Uri      string
}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - STRUCTS //////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// This methods returns a response with an error status code and a custom body.
// The body should be related to the code. For example, if the code was 404,
// then the message should relate to an something not being found.
func (r Request) Error(code int, message string) Response {
	r.Response.StatusCode = code
	r.Response.Body = message
	return r.Response
}

// This method sends a response.
func (r Request) Send() {
	r.Ctx.SetContentType(r.Response.ContentType)
	r.Ctx.SetStatusCode(r.Response.StatusCode)
	r.Ctx.SetBody([]byte(r.Response.Body))
}

// This method sends a response with an error code and message related to that
// error code.
func (r Request) SendError(code int, message string) {
	r.Ctx.SetContentType(r.Response.ContentType)
	r.Ctx.SetStatusCode(code)
	r.Ctx.SetBody([]byte(message))
}
