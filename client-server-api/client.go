package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	type item struct {
		body []byte
		err  error
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	c := make(chan item)
	go func() {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c <- item{nil, err}
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c <- item{nil, err}
			return
		}
		c <- item{body, nil}
	}()

	var body []byte
	select {
	case i := <-c:
		if i.err != nil {
			log.Println(i.err)
			panic(err)
		}
		body = i.body
	case <-ctx.Done():
		log.Println(ctx.Err())
		return
	}
	if len(body) == 0 {
		return
	}

	var p string
	err = json.Unmarshal(body, &p)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	f.WriteString("Dólar: " + p)
	f.Close()
}
