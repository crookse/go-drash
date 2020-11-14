package http

import (
	"fmt"
	"regexp"
)

var regexUriMatches = regexp.MustCompile("(:[^(/]+|{[^0-9][^}]*})")
var regexUriReplacement= "([^/]+)"
var regexUriColon = regexp.MustCompile(":")

type ResourceUrisParsed struct {
	OgPath string
	RegexPath string
	PathParamNames []string
}

type methods map[string]interface{}

type Resource struct {
	// Http methods
	DELETE func(r Request) Response
	GET func(r Request) Response
	POST func(r Request) Response
	PUT func(r Request) Response

	Methods methods
	Uris []string
	UrisParsed []ResourceUrisParsed
	response Response
}

func (r *Resource) ParseUris() {
	var uris = []ResourceUrisParsed{}

	for i := range r.Uris {
		uri := r.Uris[i]
		uris = append(r.UrisParsed, parseUri(uri))
	}

	r.UrisParsed = uris
}

func parseUri(path string) ResourceUrisParsed {
	return ResourceUrisParsed{
		OgPath: path,
		RegexPath: getRegexPath(path),
		PathParamNames: getPathParamNames(path),
	}
}

func getRegexPath(path string) string {
	result := regexUriMatches.ReplaceAllString(path, regexUriReplacement)
	return fmt.Sprintf("^%s/?$", result)
}

func getPathParamNames(path string) []string {
	matches := regexUriMatches.FindAllString(path, -1)
	for i := range matches {
		uri := matches[i]
		matches[i] = fmt.Sprintf("%s", regexUriColon.ReplaceAllString(uri, ""))
	}
	return matches
}
