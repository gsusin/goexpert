package main

import (
	"context"
	"github.com/gsusin/goexpert/observability/internal/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var logger = log.New(os.Stderr, "zipkin-collector", log.Ldate|log.Ltime|log.Llongfile)

func main() {
	http.HandleFunc("/temp", serviceA)
	http.ListenAndServe(":8080", nil)
}

func serviceA(w http.ResponseWriter, r *http.Request) {
	url := "http://zipkin-collector:9411/api/v2/spans"

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initTracer(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tr := otel.GetTracerProvider().Tracer("component-main")
	ctx, span := tr.Start(ctx, "main", trace.WithSpanKind(trace.SpanKindServer))
	web.ValidateCepAndForward(ctx, &w, r, web.ShowTemperature)
	span.End()
}

func initTracer(url string) (func(context.Context) error, error) {
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("zipkin-collector"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
