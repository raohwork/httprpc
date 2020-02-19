package rpctool

import (
	"net/http"

	"github.com/raohwork/httprpc"
)

func init() {
	_ = httprpc.Middleware(Nop{})
}

type Nop struct{}

func (m Nop) Send(h http.Header, body []byte) (ret []byte, err error) {
	return body, nil
}

func (m Nop) Receive(h http.Header, body []byte) (ret []byte, err error) {
	return body, nil
}
