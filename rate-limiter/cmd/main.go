package main

import (
	"net/http"

	"github.com/gsusin/goexpert/rate-limiter/limiter"
)

func main() {
	mux := http.NewServeMux()
	LS := LimitedServer{}
	LS.lim = limiter.NewLimitedHandler()
	mux.HandleFunc("/course", LS.LimitedGetCourseId)
	http.ListenAndServe(":8080", mux)
}

type LimitedServer struct {
	lim *limiter.LimitedHandler
}

func (ls *LimitedServer) LimitedGetCourseId(w http.ResponseWriter, r *http.Request) {
	ls.lim.LimitedFunc(w, r, GetCourseId)(w, r)
}

func GetCourseId(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Course queried"))
	println("API_KEY = ", r.Header.Get("API_KEY"))
}
