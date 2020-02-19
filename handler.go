package httprpc

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type httpHandler struct {
	h Handler
	m []Middleware
	c Codec
}

func (m *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	for x := len(m.m) - 1; x >= 0; x-- {
		buf, err := m.m[x].Receive(r.Header, reqBody)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		reqBody = buf
	}

	req := request{
		q:   r,
		dec: m.c.NewDecoder(bytes.NewReader(reqBody)),
	}

	respBody := &bytes.Buffer{}
	resp := response{
		h:   w.Header(),
		enc: m.c.NewEncoder(respBody),
	}

	m.h.ServeRPC(resp, req)

	body := respBody.Bytes()
	for _, m := range m.m {
		body, err = m.Send(resp.h, body)
		if err != nil {
			w.WriteHeader(400)
			return
		}
	}

	w.Write(body)
}

// Response defines how your handler write response to client
type Response interface {
	// Use this to set http header, supports http trailer
	Header() (ret http.Header)
	// Use this to write data to http response body
	Encode(v ...interface{}) (err error)
}

type response struct {
	h   http.Header
	enc RPCEncoder
}

func (r response) Header() (ret http.Header) {
	return r.h
}

func (r response) Encode(v ...interface{}) (err error) {
	return r.enc.EncodeRPC(v...)
}

// Request defines how your handler read client requests
type Request interface {
	// Get raw http request with this
	Request() (ret *http.Request)
	// Read http request body with this
	Decode(v ...interface{}) (err error)
}

type request struct {
	q   *http.Request
	dec RPCDecoder
}

func (r request) Request() (ret *http.Request) {
	return r.q
}

func (r request) Decode(v ...interface{}) (err error) {
	return r.dec.DecodeRPC(v...)
}

// Handler defines interface of httprpc handler
type Handler interface {
	ServeRPC(resp Response, req Request)
}

// HandlerFunc is a simplified implementation of Handler
type HandlerFunc func(resp Response, req Request)

func (h HandlerFunc) ServeRPC(resp Response, req Request) {
	h(resp, req)
}
