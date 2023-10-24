package internal

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
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

type Item1 struct {
	Result Apicep
	Err    error
}

type Item2 struct {
	Result Viacep
	Err    error
}

func Query_apicep(ctx context.Context, c chan Item1, cep string) {
	var p Apicep
	url := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	req1, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c <- Item1{p, err}
		return
	}
	resp, err := http.DefaultClient.Do(req1)
	if err != nil {
		c <- Item1{p, err}
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c <- Item1{p, err}
		return
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		c <- Item1{p, err}
		return
	}
	//time.Sleep(time.Duration(500) * time.Millisecond)
	c <- Item1{p, nil}
}

func Query_viacep(ctx context.Context, c chan Item2, cep string) {
	var p Viacep
	url := "http://viacep.com.br/ws/" + cep + "/json"
	req1, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c <- Item2{p, err}
		return
	}
	resp, err := http.DefaultClient.Do(req1)
	if err != nil {
		c <- Item2{p, err}
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c <- Item2{p, err}
		return
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		c <- Item2{p, err}
		return
	}
	//time.Sleep(time.Duration(500) * time.Millisecond)
	c <- Item2{p, nil}
}
