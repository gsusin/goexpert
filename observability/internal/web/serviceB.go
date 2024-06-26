package web

import (
	"context"
	weather "github.com/gsusin/goexpert/observability/pkg"
	"net/http"
)

func ShowTemperature(ctx context.Context, w *http.ResponseWriter, cep string) {
	code, body := weather.GetTemperature(ctx, cep)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Cache-Control", "no-cache")
	(*w).WriteHeader(code)
	(*w).Write(body)
}
