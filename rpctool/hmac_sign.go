package rpctool

import (
	"crypto/hmac"
	"encoding/hex"
	"hash"
	"net/http"

	"github.com/raohwork/httprpc"
)

type ErrHMACSign struct {
	origin error
}

func (e ErrHMACSign) Unwrap() (ret error) {
	return e.origin
}

func (e ErrHMACSign) Error() (ret string) {
	const s = "rpctool: HMACSign: invalid signature"

	if e.origin == nil {
		return s
	}
	return s + ": " + e.origin.Error()
}

// HMACSign creates a middleware which addes/verifies signature of request/response body
func HMACSign(headerKey string, f func() hash.Hash, hmacKey []byte) (ret httprpc.Middleware) {
	return &hmacSign{
		f:   f,
		key: hmacKey,
		tag: headerKey,
	}
}

type hmacSign struct {
	f   func() hash.Hash
	key []byte
	tag string
}

func (m *hmacSign) e(err error) error {
	return ErrHMACSign{err}
}

func (m *hmacSign) hashfunc() (ret hash.Hash) {
	return hmac.New(m.f, m.key)
}

func (m *hmacSign) Send(h http.Header, body []byte) (ret []byte, err error) {
	fhash := m.hashfunc()
	fhash.Write(body)
	sum := fhash.Sum(nil)
	h.Set(m.tag, hex.EncodeToString(sum))

	return body, nil
}

func (m *hmacSign) Receive(h http.Header, body []byte) (ret []byte, err error) {
	fhash := m.hashfunc()
	fhash.Write(body)
	expect := fhash.Sum(nil)
	actual, err := hex.DecodeString(h.Get(m.tag))
	if err != nil {
		return
	}
	if !hmac.Equal(expect, actual) {
		err = m.e(nil)
		return
	}

	return body, nil
}
