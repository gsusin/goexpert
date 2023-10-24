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

	c1 := make(chan internal.Item1)
	c2 := make(chan internal.Item2)
	go internal.Query_apicep(ctx, c1, cep)
	go internal.Query_viacep(ctx, c2, cep)

	for {
		select {
		case i := <-c1:
			if i.Err != nil {
				panic(i.Err)
			}
			p := i.Result
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
			if i.Err != nil {
				panic(i.Err)
			}
			p := i.Result
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
