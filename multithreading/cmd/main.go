package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type apicep struct {
	Code       string
	State      string
	City       string
	District   string
	Address    string `json:"address"`
	Status     int
	Ok         bool `json:"ok"`
	StatusText string
}

type viacep struct {
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

type item1 struct {
	result apicep
	err    error
}

type item2 struct {
	result viacep
	err    error
}

func query_apicep(ctx context.Context, c chan item1, cep string) {
	url := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	req1, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		//TODO
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req1)
	if err != nil {
		//TODO
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//TODO
		panic(err)
	}
	var p apicep
	err = json.Unmarshal(body, &p)
	if err != nil {
		//TODO
		panic(err)
	}
	//time.Sleep(time.Duration(500) * time.Millisecond)
	c <- item1{p, nil}
}

func query_viacep(ctx context.Context, c chan item2, cep string) {
	url := "http://viacep.com.br/ws/" + cep + "/json"
	req1, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req1)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var p viacep
	err = json.Unmarshal(body, &p)
	if err != nil {
		panic(err)
	}
	//time.Sleep(time.Duration(500) * time.Millisecond)
	c <- item2{p, nil}
}

func main() {
	const defaultCep = "88010-100"
	var cep string
	if len(os.Args) > 1 {
		cep = os.Args[1]
	} else {
		cep = defaultCep
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c1 := make(chan item1)
	c2 := make(chan item2)
	go query_apicep(ctx, c1, cep)
	go query_viacep(ctx, c2, cep)

	for {
		select {
		case i := <-c1:
			if i.err != nil {
				panic(i.err)
			}
			p := i.result
			if !p.Ok {
				fmt.Printf("API apicep.com returned error. Status: %d.\n", p.Status)
				if p.Status == 200 {
					return
				}
			} else {
				fmt.Println("Fastest API: apicep.com")
				fmt.Printf("CEP: %v\nAddress:\n%v\n%v\n%v %v\n", p.Code, p.Address, p.District, p.City, p.State)
				return
			}
		case i := <-c2:
			if i.err != nil {
				panic(i.err)
			}
			p := i.result
			if p.Erro {
				fmt.Println("API viacep.com.br returned error")
				return
			}
			fmt.Println("Fastest API: viacep.com.br")
			fmt.Printf("CEP: %v\nAddress:\n%v\n%v\n%v %v\n", p.Cep, p.Logradouro, p.Bairro, p.Localidade, p.Uf)
			return
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		}
	}

}
