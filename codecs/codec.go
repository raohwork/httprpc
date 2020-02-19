package codecs

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"

	"github.com/raohwork/httprpc"
)

// GOB creates a codec that encode/decode in gob format
func GOB() (ret httprpc.Codec) {
	return Danger(
		func(w io.Writer) DangerEncoder {
			return gob.NewEncoder(w)
		},
		func(r io.Reader) DangerDecoder {
			return gob.NewDecoder(r)
		},
	)
}

// JSON creates a codec that encode/decode in json format
func JSON() (ret httprpc.Codec) {
	return Danger(
		func(w io.Writer) DangerEncoder {
			return json.NewEncoder(w)
		},
		func(r io.Reader) DangerDecoder {
			return json.NewDecoder(r)
		},
	)
}

// XML creates a codec that encode/decode in xml format
func XML() (ret httprpc.Codec) {
	return Danger(
		func(w io.Writer) DangerEncoder {
			return xml.NewEncoder(w)
		},
		func(r io.Reader) DangerDecoder {
			return xml.NewDecoder(r)
		},
	)
}
