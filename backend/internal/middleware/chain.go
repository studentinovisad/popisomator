package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(mws ...Middleware) http.Handler {
	var h http.Handler = http.NotFoundHandler() // never reached, last mw always overrides it
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// Handle turns a plain handler func into a terminal Middleware, so it can go last in Chain.
func Handle(h http.HandlerFunc) Middleware {
	return func(http.Handler) http.Handler {
		return h
	}
}
