package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gsusin/goexpert/multithreading/internal"
)

func main() {
	const defaultCep = "88010-100"
	var cep string
	if len(os.Args) > 1 {
		cep = os.Args[1]
	} else {
		cep = defaultCep
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	v, err := internal.QueryFastest(ctx, cep)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch result := v.(type) {
	case internal.Apicep:
		if !result.Ok {
			fmt.Printf("API apicep.com returned error. Status: %d.\n", result.Status)
		} else {
			fmt.Println("Fastest API: apicep.com")
			fmt.Printf("CEP: %v\nAddress:\n%v\n%v\n%v %v\n", result.Code, result.Address, result.District, result.City, result.State)
		}
	case internal.Viacep:
		if result.Erro {
			fmt.Println("API viacep.com.br returned error")
		} else {
			fmt.Println("Fastest API: viacep.com.br")
			fmt.Printf("CEP: %v\nAddress:\n%v\n%v\n%v %v\n", result.Cep, result.Logradouro, result.Bairro, result.Localidade, result.Uf)
		}
	}
}
