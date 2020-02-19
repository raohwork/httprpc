package httprpc

import (
	"net/http"
)

// Mux wraps http.ServeMux so you can register your handlers with ease
//
// Since it only works as a wrapper, all restriction of http.ServeMux applies.
//
// Order of Middleware
//
// Middleware is stored in slice internally, so you can register same Middleware
// multiple times.
type Mux interface {
	// Adds a middleware
	With(m Middleware) (ret Mux)
	// Register a handler to specified endpoint
	Register(name string, h Handler)
}

// NewMux creates a Mux instance
func NewMux(codec Codec, httpmux *http.ServeMux) (ret Mux) {
	if httpmux == nil {
		httpmux = http.DefaultServeMux
	}
	return &mux{
		Codec: codec,
		mux:   httpmux,
	}
}

type mux struct {
	Middleware []Middleware
	Codec      Codec
	mux        *http.ServeMux
}

func (m *mux) With(w Middleware) (ret Mux) {
	m.Middleware = append(m.Middleware, w)
	return m
}

func (m *mux) Register(name string, h Handler) {
	m.mux.Handle("/"+name, &httpHandler{
		h: h,
		m: m.Middleware,
		c: m.Codec,
	})
}
