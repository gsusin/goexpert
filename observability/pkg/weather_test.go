package weather

import (
	"net/http"
	"testing"
)

func init() {
	http.HandleFunc("/temp", GetTempHandler)
	go http.ListenAndServe(":8080", nil)

}

func GetTempHandler(w http.ResponseWriter, r *http.Request) {
	code, body := GetTemperature(w, r)
	w.WriteHeader(code)
	w.Write([]byte(body))
}

func TestValidCep(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/temp?cep=01324-040")
	if err != nil {
		t.Fatal("Error in GET request")
	}
	if resp.StatusCode != 200 {

		t.Errorf("Returned status code %d", resp.StatusCode)
	}
}

func TestInvalidCep(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/temp?cep=01324-04")
	if err != nil {
		t.Fatal("Error in GET request")
	}
	if resp.StatusCode != 422 {

		t.Errorf("Returned status code %d", resp.StatusCode)
	}
}

func TestNotFoundCep(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/temp?cep=89222541")
	if err != nil {
		t.Fatal("Error in GET request")
	}
	if resp.StatusCode != 404 {

		t.Errorf("Returned status code %d", resp.StatusCode)
	}
}
