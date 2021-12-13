package main

import (
	"fmt"
	"net/http"
)

type Router struct {
	handlersMap map[API]http.HandlerFunc
}

type API struct {
	Method string
	Path   string
}

func NewRouter() *Router {
	return &Router{
		handlersMap: make(map[API]http.HandlerFunc),
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api := API{Method: r.Method, Path: r.URL.EscapedPath()}

	if handler, ok := router.handlersMap[api]; !ok {
		fmt.Printf("api for %s %s not found", api.Method, api.Path)
		http.NotFound(w, r)
		return
	} else {
		handler(w, r)
	}
}

func (router *Router) AddAPI(method string, path string, fn http.HandlerFunc) {
	api := API{Method: method, Path: path} // TODO: validate methods and clean path
	router.handlersMap[api] = fn
}
