package httprpc

type RPCEncoder interface {
	EncodeRPC(v ...interface{}) (err error)
}

// RPCDecoder is a helper for decoding gob stream.
type RPCDecoder interface {
	// wraps gob.Decoder.Decode, always return ErrInfra if err != nil
	DecodeRPC(v ...interface{}) (err error)
}
