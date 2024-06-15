package web

import (
	"context"
	"go.opentelemetry.io/otel"
	"net/http"
	"strings"
)

type CepProcessor func(context.Context, http.ResponseWriter, string)

func ValidateCepAndForward(ctx context.Context, w http.ResponseWriter, r *http.Request, f CepProcessor) {

	tr := otel.GetTracerProvider().Tracer("component-serviceA")
	_, span := tr.Start(ctx, "serviceA")
	cep := r.Header.Get("cep")
	cep = strings.Map(keepNumerals, cep)
	if len(cep) != 8 {
		w.WriteHeader(422)
		w.Write([]byte("invalid zipcode"))
		return
	}
	f(ctx, w, cep)
	span.End()
}

func keepNumerals(r rune) rune {
	if r >= '0' && r <= '9' {
		return r
	}
	return -1
}
