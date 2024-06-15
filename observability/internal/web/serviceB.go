package web

import (
	"context"
	weather "github.com/gsusin/goexpert/observability/pkg"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
	//"time"
)

func ShowTemperature(ctx context.Context, w http.ResponseWriter, cep string) {
	tr := otel.GetTracerProvider().Tracer("component-serviceB")
	_, span := tr.Start(ctx, "serviceB")
	code, body := weather.GetTemperature(cep)
	log.Println("Executou GetTemperature()")
	span.End()
	w.WriteHeader(code)
	w.Write([]byte(body))
	//<-time.After(6 * time.Millisecond)
}
