package limiter_test

import (
	"testing"
	"time"

	"github.com/gsusin/goexpert/rate-limiter/limiter"
	"io"
	"net/http"
)

type LimitedServer struct {
	lim *limiter.LimitedHandler
}

func (ls *LimitedServer) LimitedGetCourseId(w http.ResponseWriter, r *http.Request) {
	ls.lim.LimitedFunc(w, r, GetCourseId)(w, r)
}

func GetCourseId(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Course ok"))
}

func init() {
	mux := http.NewServeMux()
	ls := LimitedServer{}
	as := limiter.NewMemoryStorage()
	ls.lim = limiter.NewLimitedHandler(&as)
	mux.HandleFunc("/course", ls.LimitedGetCourseId)
	go http.ListenAndServe(":8080", mux)
}

func TestInsideLimit(t *testing.T) {
	time.Sleep(1000 * time.Millisecond)
	req, err := http.NewRequest("GET", "http://localhost:8080/course", nil)
	if err != nil {
		t.Fatal("Error creating Request")
	}
	req.Header.Add("API_KEY", "")
	client := &http.Client{}
	for i := 0; i < 10; i++ {
		time.Sleep(105 * time.Millisecond)
		var err error
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal("Error in GET request")
		}
		t.Log("Returned status code ", (*resp).StatusCode)
		body, err := io.ReadAll((*resp).Body)
		(*resp).Body.Close()
		t.Log(string(body))

		if (*resp).StatusCode != 200 {
			t.Errorf("Returned status code %d", (*resp).StatusCode)
		}
	}
}

func TestOutsideLimit(t *testing.T) {
	time.Sleep(1000 * time.Millisecond)
	req, err := http.NewRequest("GET", "http://localhost:8080/course", nil)
	if err != nil {
		t.Fatal("Error creating Request")
	}
	req.Header.Add("API_KEY", "")
	client := &http.Client{}
	for i := 0; i < 11; i++ {
		time.Sleep(95 * time.Millisecond)
		var err error
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal("Error in GET request")
		}
		t.Log("Returned status code ", (*resp).StatusCode)
		body, err := io.ReadAll((*resp).Body)
		(*resp).Body.Close()
		t.Log(string(body))

		if i == 10 {
			if (*resp).StatusCode == 429 {
				return
			}
			t.Errorf("Returned status code %d", (*resp).StatusCode)
		}
	}
}

func TestUnblock(t *testing.T) {
	time.Sleep(1000 * time.Millisecond)
	req, err := http.NewRequest("GET", "http://localhost:8080/course", nil)
	if err != nil {
		t.Fatal("Error creating Request")
	}
	req.Header.Add("API_KEY", "")
	client := &http.Client{}
	for i := 0; i < 24; i++ {
		time.Sleep(85 * time.Millisecond)
		var err error
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal("Error in GET request")
		}
		t.Log("Returned status code ", (*resp).StatusCode)
		body, err := io.ReadAll((*resp).Body)
		(*resp).Body.Close()
		t.Log(string(body))

		if i == 23 {
			if (*resp).StatusCode == 200 {
				return
			}
			t.Errorf("Returned status code %d", (*resp).StatusCode)
		}
	}
}

func TestToken(t *testing.T) {
	time.Sleep(1000 * time.Millisecond)
	req, err := http.NewRequest("GET", "http://localhost:8080/course", nil)
	if err != nil {
		t.Fatal("Error creating Request")
	}
	req.Header.Add("API_KEY", "HIGH")
	client := &http.Client{}
	for i := 0; i < 20; i++ {
		time.Sleep(55 * time.Millisecond)
		var err error
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal("Error in GET request")
		}
		t.Log("Returned status code ", (*resp).StatusCode)
		body, err := io.ReadAll((*resp).Body)
		(*resp).Body.Close()
		t.Log(string(body))

		if (*resp).StatusCode != 200 {
			t.Errorf("Returned status code %d", (*resp).StatusCode)
		}
	}
}
