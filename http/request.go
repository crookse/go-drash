package http

import (
	"github.com/valyala/fasthttp"
)

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - VARIABLE DECLARATIONS ////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Request struct {
	Ctx      *fasthttp.RequestCtx
	Response Response
}

func (r Request) Error(code int, message string) Response {
	r.Response.StatusCode = code
	r.Response.Body = message
	return r.Response
}

func (r Request) Send() {
	r.Ctx.SetContentType(r.Response.ContentType)
	r.Ctx.SetStatusCode(r.Response.StatusCode)
	r.Ctx.SetBody([]byte(r.Response.Body))
}

func (r Request) SendError(code int, message string) {
	r.Ctx.SetContentType(r.Response.ContentType)
	r.Ctx.SetStatusCode(code)
	r.Ctx.SetBody([]byte(message))
}
