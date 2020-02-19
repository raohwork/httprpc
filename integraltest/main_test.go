package integraltest

import (
	"crypto/des"
	"crypto/sha256"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/raohwork/httprpc"
	"github.com/raohwork/httprpc/codecs"
	"github.com/raohwork/httprpc/rpctool"
)

type testsrv struct {
	srv *http.Server
}

func middlewares() (ret []httprpc.Middleware) {
	tripdes, err := des.NewTripleDESCipher([]byte("123456789012345678901234"))
	if err != nil {
		panic(err)
	}
	return []httprpc.Middleware{
		rpctool.CBCEncrypt(
			tripdes,
			[]byte("12345678"),
		),
		rpctool.SharedKeyAuth(
			"X-API-Key", "HackMe",
		),
		rpctool.HMACSign(
			"X-Signature", sha256.New, []byte("DeAdBeEf"),
		),
	}
}

func (s *testsrv) hello(resp httprpc.Response, req httprpc.Request) {
	var (
		name    string
		surname string
	)

	if err := req.Decode(&name); err != nil {
		log.Print("failed to decode name:", err)
		return
	}

	if err := req.Decode(&surname); err != nil {
		log.Print("failed to decode surname:", err)
		return
	}

	resp.Encode("hello, mr. " + surname)
	resp.Encode("hello, " + name)
}

func (s *testsrv) run() {
	m := &http.ServeMux{}
	mux := httprpc.NewMux(codecs.GOB(), m)
	s.srv = &http.Server{
		Addr:    ":29876",
		Handler: m,
	}

	for _, m := range middlewares() {
		mux = mux.With(m)
	}
	mux.Register("hello", httprpc.HandlerFunc(s.hello))
	s.srv.ListenAndServe()
}

func (s *testsrv) close() {
	s.srv.Close()
}

func TestOK(t *testing.T) {
	srv := &testsrv{}
	go srv.run()
	defer srv.close()

	time.Sleep(1 * time.Second)

	caller := httprpc.NewCaller(codecs.GOB(), nil)
	for _, m := range middlewares() {
		caller = caller.With(m)
	}

	for x := 0; x < 100; x++ {
		t.Log(x)
		func() {
			dec, err := caller.Call(
				"http://127.0.0.1:29876/hello",
				"qwe",
				"asd",
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var sur, name string
			if err := dec.DecodeRPC(&sur, &name); err != nil {
				t.Fatalf("unexpected decode error: %v", err)
			}
			if sur != "hello, mr. asd" {
				t.Errorf("unexpected surname: %s", sur)
			}

			if name != "hello, qwe" {
				t.Errorf("unexpected name: %s", name)
			}
		}()
	}
}

func BenchmarkOK(b *testing.B) {
	srv := &testsrv{}
	go srv.run()
	defer srv.close()

	time.Sleep(1 * time.Second)

	caller := httprpc.NewCaller(codecs.GOB(), nil)
	for _, m := range middlewares() {
		caller = caller.With(m)
	}

	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		dec, err := caller.Call(
			"http://127.0.0.1:29876/hello",
			"qwe", "asd",
		)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}

		var sur, name string
		if err := dec.DecodeRPC(&sur, &name); err != nil {
			b.Fatalf("unexpected decode error: %v", err)
		}
		if sur != "hello, mr. asd" {
			b.Fatalf("unexpected surname: %s", sur)
		}

		if name != "hello, qwe" {
			b.Fatalf("unexpected name: %s", name)
		}
	}
}
