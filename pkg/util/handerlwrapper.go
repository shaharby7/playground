package util

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

// This rapper provides an easy interface to transform http handlers of
// net/http into valid httprouter handlers.
//
// This is useful if you want to combine packages like middlewares
// that are usually not compatible with the handlers of httprouter.
//
// Based on Nicolas Merouze's awesome guide about 3rd party routers
// in Go (www.nicolasmerouze.com/guide-routers-golang).
// This package is a stripped down version of Merouze's code example
// at Github (www.gist.github.com/nmerouze/5ed810218c661b40f5c4).

type router struct {
	*httprouter.Router
}

func (r *router) Get(path string, handler http.Handler) {
	r.GET(path, WrapHandler(handler))
}

func (r *router) Post(path string, handler http.Handler) {
	r.POST(path, WrapHandler(handler))
}

func (r *router) Put(path string, handler http.Handler) {
	r.PUT(path, WrapHandler(handler))
}

func (r *router) Patch(path string, handler http.Handler) {
	r.PATCH(path, WrapHandler(handler))
}

func (r *router) Delete(path string, handler http.Handler) {
	r.DELETE(path, WrapHandler(handler))
}

func (r *router) Head(path string, handler http.Handler) {
	r.HEAD(path, WrapHandler(handler))
}

func (r *router) Options(path string, handler http.Handler) {
	r.OPTIONS(path, WrapHandler(handler))
}

// NewRouter returns a router instance that acts
// like httprouter and enables access to the
// wrapper functions.
func NewRouter() *router {
	return &router{httprouter.New()}
}

// WrapHandler transforms ususal handlers (http.Handler) of the standard library
// package net/http into valid ones of httprouter by adding the params
// (httprouter.Params) parameter.
//
// Use this function to add a context to middlewares or other things that should be
// shared between both handler types.
func WrapHandler(h http.Handler) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		context.Set(req, "params", ps)
		h.ServeHTTP(rw, req)
	}
}
