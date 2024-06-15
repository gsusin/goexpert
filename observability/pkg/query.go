package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Apicep struct {
	Code       string
	State      string
	City       string
	District   string
	Address    string
	Status     int
	Ok         bool
	StatusText string
}

type Viacep struct {
	Cep         string
	Logradouro  string
	Complemento string
	Bairro      string
	Localidade  string
	Uf          string
	Ibge        string
	Gia         string
	Ddd         string
	Siafi       string
	Erro        bool
}

const (
	apicepId = iota
	viacepId
)

type queryDelay struct {
	duration time.Duration
	api      uint8
}

type delayContextKey string

func query[A Apicep | Viacep, CA chan A](ctx context.Context, c CA, ce chan error, url string) {
	req1, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ce <- err
		return
	}
	resp, err := http.DefaultClient.Do(req1)
	if err != nil {
		ce <- err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusBadRequest {
		ce <- fmt.Errorf("StatusBadRequest")
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ce <- err
		return
	}
	var addr A
	err = json.Unmarshal(body, &addr)
	if err != nil {
		ce <- err
		return
	}
	if d := ctx.Value(delayContextKey("delay")); d != nil {
		time.Sleep(d.(time.Duration))
	}
	c <- addr
}

func QueryFastest(ctx context.Context, cep string) (any, error) {
	c1 := make(chan Apicep)
	c2 := make(chan Viacep)
	ce := make(chan error)

	ctx1, ctx2 := ctx, ctx
	if v := ctx.Value(delayContextKey("delayPlan")); v != nil {
		switch dP := v.(queryDelay); dP.api {
		case apicepId:
			ctx1 = context.WithValue(ctx, delayContextKey("delay"), dP.duration)
		case viacepId:
			ctx2 = context.WithValue(ctx, delayContextKey("delay"), dP.duration)
		}
	}

	go query(ctx1, c1, ce, "https://cdn.apicep.com/file/apicep/"+cep+".json")
	go query(ctx2, c2, ce, "http://viacep.com.br/ws/"+cep+"/json")

	for {
		select {
		case v := <-ce:
			return nil, v
		case v := <-c1:
			if !v.Ok {
				if v.Status == 200 {
					return v, nil
				}
			}
			return v, nil
		case v := <-c2:
			return v, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
