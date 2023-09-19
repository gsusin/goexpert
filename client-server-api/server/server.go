package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Rate struct {
	USDBRL struct {
		Price     string `json:"bid"`
		Timestamp string
	}
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "./exchange.db")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/cotacao", getPriceHandler)
	http.ListenAndServe(":8080", nil)
}

func getPriceHandler(w http.ResponseWriter, r *http.Request) {
	ctxPrice, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	rate, err := getPrice(ctxPrice)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctxPersist, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err = persistRate(ctxPersist, rate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rate.USDBRL.Price)
}

func getPrice(ctx context.Context) (*Rate, error) {
	type item struct {
		body []byte
		err  error
	}
	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	c := make(chan item, 1)
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
			return nil, err
		}
		body = i.body
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	var r Rate
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func persistRate(ctx context.Context, rate *Rate) error {
	c := make(chan error)
	go func() {
		stmt, err := db.Prepare("insert into rate(timestamp, price) values(?, ?)")
		if err != nil {
			c <- err
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(rate.USDBRL.Timestamp, rate.USDBRL.Price)
		if err != nil {
			c <- err
			return
		}
		c <- nil
	}()

	select {
	case err := <-c:
		if err != nil {
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
