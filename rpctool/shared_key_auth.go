package rpctool

import (
	"errors"
	"net/http"

	"github.com/raohwork/httprpc"
)

type sharedKeyAuth struct {
	tag string
	key string
}

func (m sharedKeyAuth) Send(h http.Header, body []byte) (ret []byte, err error) {
	h.Set(m.tag, m.key)
	return body, nil
}

func (m sharedKeyAuth) Receive(h http.Header, body []byte) (ret []byte, err error) {
	v := h.Get(m.tag)
	if m.key != v {
		err = errors.New("rpctool: auth error: invalid shared secret")
		return
	}
	return body, nil
}

func SharedKeyAuth(headerKey, sharedKey string) (ret httprpc.Middleware) {
	return sharedKeyAuth{
		tag: headerKey,
		key: sharedKey,
	}
}
