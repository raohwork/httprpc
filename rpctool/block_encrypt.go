package rpctool

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"net/http"

	"github.com/raohwork/httprpc"
)

type cbcEncrypt struct {
	b  cipher.Block
	iv []byte
}

func padding(data []byte, blockSize int) (ret []byte) {
	l := len(data)
	delta := blockSize - (l % blockSize)
	if delta == blockSize {
		return data
	}

	return append(data, bytes.Repeat([]byte{0}, delta)...)
}

func CBCEncrypt(b cipher.Block, iv []byte) (ret httprpc.Middleware) {
	if l := len(iv); l != b.BlockSize() {
		err := errors.New(
			"rpctool: CBCEncrypt: size of iv mismatch with block size",
		)
		panic(err)
	}

	return &cbcEncrypt{
		b:  b,
		iv: iv,
	}
}

func (m *cbcEncrypt) Send(h http.Header, body []byte) (ret []byte, err error) {
	l := len(body)
	size := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.PutUvarint(size, uint64(l))

	data := padding(body, m.b.BlockSize())
	dest := make([]byte, len(data))
	enc := cipher.NewCBCEncrypter(m.b, m.iv)
	enc.CryptBlocks(dest, data)

	return append(size, dest...), nil
}

func (m *cbcEncrypt) Receive(h http.Header, body []byte) (ret []byte, err error) {
	sizeBuf := body[:8]
	data := body[8:]
	size, _ := binary.Uvarint(sizeBuf)
	dest := make([]byte, len(data))
	dec := cipher.NewCBCDecrypter(m.b, m.iv)
	dec.CryptBlocks(dest, data)

	return dest[:size], nil
}
