package main

import (
	"net/http"

	"github.com/gsusin/goexpert/rate-limiter/limiter"
)

func main() {
	mux := http.NewServeMux()
	ls := LimitedServer{}
	as := limiter.NewRedisStorage()
	ls.lim = limiter.NewLimitedHandler(&as)
	mux.HandleFunc("/course", ls.LimitedGetCourseId)
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
}
