package httprpc

import "io"

type Codec interface {
	NewEncoder(w io.Writer) RPCEncoder
	NewDecoder(r io.Reader) RPCDecoder
}
