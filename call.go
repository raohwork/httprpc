package httprpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Result represents response from remote endpoint
type Result struct {
	D RPCDecoder
	E error
}

// Caller defines a RPC client
type Caller interface {
	// execute remote procedure synchronously
	Call(uri string, params ...interface{}) (dec RPCDecoder, err error)
	// execute remote procedure asynchronously
	Do(uri string, params ...interface{}) (ret chan *Result)
	// add middleware to this caller
	With(m Middleware) (ret Caller)
}

func NewCaller(codec Codec, hc *http.Client) (ret Caller) {
	if hc == nil {
		hc = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	return &caller{
		hc:    hc,
		codec: codec,
	}
}

type caller struct {
	hc    *http.Client
	m     []Middleware
	codec Codec
}

func (c *caller) dupe() (ret *caller) {
	return &caller{
		hc:    c.hc,
		m:     c.m,
		codec: c.codec,
	}
}

func (c *caller) With(m Middleware) (ret Caller) {
	x := c.dupe()
	if m != nil {
		x.m = append(c.m, m)
	}

	return x
}

func (c *caller) Call(uri string, params ...interface{}) (dec RPCDecoder, err error) {
	resp := <-c.Do(uri, params...)
	if resp == nil {
		return
	}
	return resp.D, resp.E
}

func (c *caller) Do(uri string, params ...interface{}) (ret chan *Result) {
	ret = make(chan *Result)
	f := func(e error) {
		ret <- &Result{E: ErrInfra{Origin: e}}
		close(ret)
	}

	reqBody := &bytes.Buffer{}
	enc := c.codec.NewEncoder(reqBody)
	err := enc.EncodeRPC(params...)
	if err != nil {
		f(err)
		return
	}

	body := reqBody.Bytes()
	h := http.Header{}
	for _, m := range c.m {
		body, err = m.Send(h, body)
		if err != nil {
			f(err)
			return
		}
	}

	req, err := http.NewRequest("POST", uri, bytes.NewReader(body))
	if err != nil {
		f(err)
		return
	}
	for k, v := range h {
		if len(v) == 1 {
			req.Header.Set(k, v[0])
			continue
		}

		for _, val := range v {
			req.Header.Add(k, val)
		}
	}

	go c.do(req, ret)

	return
}

func (c *caller) do(req *http.Request, ch chan *Result) {
	result := &Result{}
	defer close(ch)
	f := func(e error) {
		result.E = ErrInfra{Origin: e}
		ch <- result
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		f(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		f(fmt.Errorf("server returns error code: %d", resp.StatusCode))
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	for x := len(c.m) - 1; x >= 0; x-- {
		respBody, err = c.m[x].Receive(resp.Header, respBody)
		if err != nil {
			f(err)
			return
		}
	}

	result.D = c.codec.NewDecoder(bytes.NewReader(respBody))
	ch <- result
	return
}
