package main

import (
	"github.com/gsusin/goexpert/cloud-run/internal"
	"net/http"
)

func main() {
	http.HandleFunc("/temp", GetTemperature)
	http.ListenAndServe(":8080", nil)

}

func GetTemperature(w http.ResponseWriter, r *http.Request) {
	code, body := internal.GetTemperature(w, r)
	w.WriteHeader(code)
	w.Write([]byte(body))
}
