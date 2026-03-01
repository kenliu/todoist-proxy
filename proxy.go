package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

var todoistTarget = &url.URL{
	Scheme: "https",
	Host:   "api.todoist.com",
}

// newReverseProxy returns an httputil.ReverseProxy that forwards requests to
// api.todoist.com verbatim, preserving the path, query, and all headers
// (including the client's Authorization header).
func newReverseProxy() *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(todoistTarget)
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		// httputil.NewSingleHostReverseProxy sets req.URL.Host but leaves
		// req.Host as the original incoming value. Todoist requires the Host
		// header to match the target, so we override it explicitly.
		req.Host = todoistTarget.Host
	}
	return proxy
}
