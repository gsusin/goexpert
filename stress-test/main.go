/*
Copyright © 2024 Giancarlo Susin <giancarlosusin@gmail.com>

*/
package main

import (
	"github.com/gsusin/goexpert/stress-test/cmd"
	"net/http"
)

func init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/course", GetCourseId)
	go http.ListenAndServe(":8080", mux)
}

func GetCourseId(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Course queried"))
}

func main() {
	cmd.Execute()
}
