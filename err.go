package httprpc

import "fmt"

// ErrInfra denotes an error cause by infrastructure, not the service itself
//
// It includes following errors:
//
//   - io
//   - networking (can't connect to remote, interrupted, ...)
//   - gob format error
//   - errors from http.NewRequest
type ErrInfra struct {
	Origin error
}

func (e ErrInfra) Unwrap() error {
	return e.Origin
}

func (e ErrInfra) Error() (ret string) {
	const s = "httprpc: error communicating with remote"
	if e.Origin == nil {
		return s
	}
	return fmt.Sprintf(s+": %+v", e.Origin)
}
