package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gsusin/goexpert/observability/configs"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetTemperature(ctx context.Context, cep string) (int, []byte) {

	cep = strings.Map(keepNumerals, cep)
	if len(cep) != 8 {
		return 422, []byte("invalid Zipcode")
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tr := otel.GetTracerProvider().Tracer("component-serviceB")
	ctx, span := tr.Start(ctx, "zip-service")
	v, err := QueryFastest(ctx, cep)
	if err != nil {
		span.End()
		fmt.Println(err)
		return 503, []byte(fmt.Sprintf("%v", err))
	}
	var city string
	responseCode := 200
	switch result := v.(type) {
	case Apicep:
		if !result.Ok {
			responseCode = 404
		} else {
			city = result.City
		}
	case Viacep:
		if result.Erro {
			responseCode = 404
		} else {
			city = result.Localidade
		}
	}
	if responseCode == 404 {
		span.End()
		return 404, []byte("can not find zipcode")
	}
	span.End()

	_, err = configs.LoadConfig(".")
	key := viper.GetString("KEY")
	if key == "" {
		log.Fatal("Couldn't read key")
	}

	ctx, span = tr.Start(ctx, "temperature-service")
	defer span.End()
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", key, url.QueryEscape(city))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error in NewRequest(): ", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error calling Get: ", err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	var temperature struct {
		Current struct {
			City   string  `json:"city"`
			Temp_c float32 `json:"temp_C"`
			Temp_f float32 `json:"temp_F"`
			Temp_k float32 `json:"temp_K"`
		}
	}
	if err := json.Unmarshal(body, &temperature); err != nil {
		log.Fatal("Json unmarshaling failed ", err)
	}
	temperature.Current.City = city
	temperature.Current.Temp_k = temperature.Current.Temp_c + 273
	answer, err := json.Marshal(temperature)
	if err != nil {
		log.Fatal("Json marshaling failed: ", err)
	}
	return 200, answer
}

func keepNumerals(r rune) rune {
	if r >= '0' && r <= '9' {
		return r
	}
	return -1
}
