package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/samber/lo"
)

type handler func(http.ResponseWriter, *http.Request, httprouter.Params) error

type httpError struct {
	msg  string
	code int
}

func (e httpError) Error() string {
	return fmt.Sprintf("http %d: %s", e.code, e.msg)
}

func handleErrors(h handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := h(w, r, p)
		var (
			he   = httpError{}
			code = 500
			msg  = ""
		)

		if errors.As(err, &he) {
			code = he.code
		}

		if err != nil {
			log.Println("request failed:", err.Error())
			msg = err.Error()
			_, _ = w.Write([]byte(msg))
			w.WriteHeader(code)
		}
	}
}

// struct2Map is a hacky way of turning structs into maps by abusing the json encoder
func struct2Map(input any) map[string]any {
	out := make(map[string]any)
	pr, pw := io.Pipe()
	go func() { lo.Must0(json.NewEncoder(pw).Encode(input)) }()
	lo.Must0(json.NewDecoder(pr).Decode(&out))
	return out
}
