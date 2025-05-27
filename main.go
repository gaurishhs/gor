// Package gor provides a simple HTTP router with middleware support.
package gor

import (
	"net/http"
	"slices"
)

// Middleware defines a function type that takes an http.Handler and returns an http.Handler.
// This allows for the creation of middleware functions that can wrap handlers to add functionality such as logging, authentication, etc.
// The middleware can be applied to all routes in the router or to specific routes as needed.
type Middleware func(handler http.Handler) http.Handler
type router struct {
	*http.ServeMux
	mx []Middleware
}

// NewRouter creates a new router with optional middlewares.
func NewRouter(mx ...Middleware) *router {
	return &router{
		ServeMux: &http.ServeMux{},
		mx:       mx,
	}
}

// Use adds middlewares to the router which will be applied to all routers in that router.
func (r *router) Use(mx ...Middleware) {
	r.mx = append(r.mx, mx...)
}

// Group creates a new router with the same middlewares and allows for further grouping.
func (r *router) Group(fn func(r *router)) {
	fn(&router{ServeMux: r.ServeMux, mx: slices.Clone(r.mx)})
}

// wrap wraps a handler function with the router's middlewares.
func (r *router) wrap(fn http.HandlerFunc, mx []Middleware) http.Handler {
	mx = append(slices.Clone(r.mx), mx...)
	h := http.Handler(fn)
	slices.Reverse(mx)
	for _, m := range mx {
		h = m(h)
	}
	return h
}

// handle registers a handler function for the specified method, path and middlewares.
func (r *router) handle(method, path string, fn http.HandlerFunc, mx []Middleware) {
	r.Handle(method+" "+path, r.wrap(fn, mx))
}

// Get registers a handler function for the GET method at the specified path with optional middlewares.
func (r *router) Get(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodGet, path, fn, mx)
}

// Post registers a handler function for the POST method at the specified path with optional middlewares.
func (r *router) Post(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodPost, path, fn, mx)
}

// Put registers a handler function for the PUT method at the specified path with optional middlewares.
func (r *router) Put(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodPut, path, fn, mx)
}

// Delete registers a handler function for the DELETE method at the specified path with optional middlewares.
func (r *router) Delete(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodDelete, path, fn, mx)
}

// Patch registers a handler function for the PATCH method at the specified path with optional middlewares.
func (r *router) Patch(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodPatch, path, fn, mx)
}

// Head registers a handler function for the HEAD method at the specified path with optional middlewares.
func (r *router) Head(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodHead, path, fn, mx)
}

// Options registers a handler function for the OPTIONS method at the specified path with optional middlewares.
func (r *router) Options(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodOptions, path, fn, mx)
}

// Custom registers a handler function for a custom method at the specified path with optional middlewares.
func (r *router) Custom(method, path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(method, path, fn, mx)
}
