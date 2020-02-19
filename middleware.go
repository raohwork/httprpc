package httprpc

import (
	"net/http"
)

type Middleware interface {
	Send(h http.Header, body []byte) (ret []byte, err error)
	Receive(h http.Header, body []byte) (ret []byte, err error)
}
