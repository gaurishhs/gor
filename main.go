// Package gor provides a simple HTTP router with middleware support.
package gor

import (
	"net/http"
	"slices"
)

// Middleware is the type for HTTP middleware functions.
type Middleware func(next http.Handler) http.Handler

// Router is a simple HTTP router that supports middleware.
type Router struct {
	*http.ServeMux
	mx []Middleware
}

// NewRouter creates a new router with optional middlewares.
func NewRouter(mx ...Middleware) *Router {
	return &Router{
		ServeMux: &http.ServeMux{},
		mx:       mx,
	}
}

// Use adds middlewares to the router which will be applied to all routers in that router.
func (r *Router) Use(mx ...Middleware) {
	r.mx = append(r.mx, mx...)
}

// Group creates a new router with the same middlewares and allows for further grouping.
func (r *Router) Group(fn func(r *Router)) {
	fn(&Router{ServeMux: r.ServeMux, mx: slices.Clone(r.mx)})
}

// wrap wraps a handler function with the router's middlewares.
func (r *Router) wrap(fn http.HandlerFunc, mx []Middleware) http.Handler {
	mx = append(slices.Clone(r.mx), mx...)
	h := http.Handler(fn)
	slices.Reverse(mx)
	for _, m := range mx {
		h = m(h)
	}
	return h
}

// handle registers a handler function for the specified method, path and middlewares.
func (r *Router) handle(method, path string, fn http.HandlerFunc, mx []Middleware) {
	r.Handle(method+" "+path, r.wrap(fn, mx))
}

// Get registers a handler function for the GET method at the specified path with optional middlewares.
func (r *Router) Get(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodGet, path, fn, mx)
}

// Post registers a handler function for the POST method at the specified path with optional middlewares.
func (r *Router) Post(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodPost, path, fn, mx)
}

// Put registers a handler function for the PUT method at the specified path with optional middlewares.
func (r *Router) Put(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodPut, path, fn, mx)
}

// Delete registers a handler function for the DELETE method at the specified path with optional middlewares.
func (r *Router) Delete(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodDelete, path, fn, mx)
}

// Patch registers a handler function for the PATCH method at the specified path with optional middlewares.
func (r *Router) Patch(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodPatch, path, fn, mx)
}

// Head registers a handler function for the HEAD method at the specified path with optional middlewares.
func (r *Router) Head(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodHead, path, fn, mx)
}

// Options registers a handler function for the OPTIONS method at the specified path with optional middlewares.
func (r *Router) Options(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodOptions, path, fn, mx)
}

// Custom registers a handler function for a custom method at the specified path with optional middlewares.
func (r *Router) Custom(method, path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(method, path, fn, mx)
}
