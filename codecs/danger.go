package codecs

import (
	"io"

	"github.com/raohwork/httprpc"
)

// DnagerEncoder is only used for Danger()
type DangerEncoder interface {
	Encode(v interface{}) error
}

// DnagerDecoder is only used for Danger()
type DangerDecoder interface {
	Decode(v interface{}) error
}

type dangerE struct {
	enc DangerEncoder
}

func (e *dangerE) EncodeRPC(val ...interface{}) (err error) {
	for _, v := range val {
		if err = e.enc.Encode(v); err != nil {
			break
		}
	}

	return
}

type dangerD struct {
	dec DangerDecoder
}

func (d *dangerD) DecodeRPC(val ...interface{}) (err error) {
	for _, v := range val {
		if err = d.dec.Decode(v); err != nil {
			break
		}
	}

	return
}

type danger struct {
	enc func(w io.Writer) DangerEncoder
	dec func(r io.Reader) DangerDecoder
}

func (d *danger) NewEncoder(w io.Writer) (ret httprpc.RPCEncoder) {
	return &dangerE{
		enc: d.enc(w),
	}
}

func (d *danger) NewDecoder(r io.Reader) (ret httprpc.RPCDecoder) {
	return &dangerD{
		dec: d.dec(r),
	}
}

// Danger creates httprpc.Codec using builtin encoder/decoder
//
// I named it "Danger" since there's no way to validate the relation between
// encoder and decoder. Although nonsense, it is possible to pass gob.Encoder
// along with json.Decoder. Use it with care.
func Danger(
	enc func(w io.Writer) DangerEncoder,
	dec func(r io.Reader) DangerDecoder,
) (ret httprpc.Codec) {
	return &danger{enc, dec}
}
