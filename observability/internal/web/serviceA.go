package web

import (
	"context"
	"net/http"
	"strings"
)

type CepProcessor func(context.Context, *http.ResponseWriter, string)

func ValidateCepAndForward(ctx context.Context, w *http.ResponseWriter, r *http.Request, f CepProcessor) {
	cep := r.FormValue("cep")
	cep = strings.Map(keepNumerals, cep)
	if len(cep) != 8 {
		(*w).WriteHeader(422)
		(*w).Write([]byte("invalid zipcode"))
		return
	}
	f(ctx, w, cep)
}

func keepNumerals(r rune) rune {
	if r >= '0' && r <= '9' {
		return r
	}
	return -1
}
