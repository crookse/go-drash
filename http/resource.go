package http

import (
	"fmt"
	"regexp"
)

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - VARIABLE DECLARATIONS ////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

var _regexUriColon = regexp.MustCompile(":")
var _regexUriMatches = regexp.MustCompile("(:[^(/]+|{[^0-9][^}]*})")
var _regexUriReplacement = "([^/]+)"

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - STRUCTS //////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Resource struct {
	Methods    map[string]interface{}
	Uris       []string
	UrisParsed []ResourceUrisParsed
	response   Response

	// HTTP methods
	CONNECT func(r *Request) Response
	DELETE  func(r *Request) Response
	GET     func(r *Request) Response
	HEAD    func(r *Request) Response
	OPTIONS func(r *Request) Response
	PATCH   func(r *Request) Response
	POST    func(r *Request) Response
	PUT     func(r *Request) Response
	TRACE   func(r *Request) Response
}

type ResourceUrisParsed struct {
	RegexUri      string
	UriParamNames []string
}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - METHODS - EXPORTED ///////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// This method parses all URIs associated with this resource so that we can
// match request URIs to the resource's URIs.
func (r *Resource) ParseUris() {
	uris := make([]ResourceUrisParsed, len(r.Uris))

	for i, uri := range r.Uris {
		uris[i] = r.parseUri(uri)
	}

	r.UrisParsed = uris
}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - METHODS - NOT EXPORTED ///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// This methods returns the regex version of a URI. This regex version is used
// when matching a request URI to a resource.
func (r *Resource) getRegexVersionOfUri(uri string) string {
	result := _regexUriMatches.ReplaceAllString(uri, _regexUriReplacement)
	return fmt.Sprintf("^%s/?$", result)
}

// This method gets the names of all URI params in the resource's URIs. For
// example, if a resource defines /uri/:something in its URIs, then
// "something" will become a uri param name.
func (r *Resource) getUriParamNames(uri string) []string {
	matches := _regexUriMatches.FindAllString(uri, -1)

	for i := range matches {
		uri := matches[i]
		matches[i] = fmt.Sprintf(
			"%s", _regexUriColon.ReplaceAllString(uri, ""),
		)
	}

	return matches
}

// This method expands a URI into parsable parts for runtime purposes. During
// runtime, the RegexUri is used to match a request URI against and the
// UriParamNames are used to match URI param values to a URI param name.
func (r *Resource) parseUri(uri string) ResourceUrisParsed {
	return ResourceUrisParsed{
		RegexUri:     r.getRegexVersionOfUri(uri),
		UriParamNames: r.getUriParamNames(uri),
	}
}
