package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cloudrun/internal/httpapi"
)

// helper: mock ViaCEP
func mockViaCEP(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// expected path: /ws/{cep}/json
		if r.URL.Path == "/ws/01001000/json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"localidade":"São Paulo","uf":"SP"}`))
			return
		}
		if r.URL.Path == "/ws/99999999/json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"erro": true}`))
			return
		}
		w.WriteHeader(404)
	}))
}

// helper: mock WeatherAPI
func mockWeather(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "São Paulo,SP,BR" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"current":{"temp_c":28.49}}`))
			return
		}
		w.WriteHeader(404)
	}))
}

func TestInvalidZipcode(t *testing.T) {
	os.Setenv("VIACEP_BASE", "http://invalid") // won't be called
	os.Setenv("WEATHERAPI_BASE", "http://invalid")
	os.Setenv("WEATHERAPI_KEY", "x")

	s := httpapi.NewServer()
	ts := httptest.NewServer(s.Router())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/weather?cep=ABC")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 422 {
		t.Fatalf("expected 422, got %d", res.StatusCode)
	}
}

func TestZipNotFound(t *testing.T) {
	viacep := mockViaCEP(t)
	defer viacep.Close()
	weather := mockWeather(t)
	defer weather.Close()

	os.Setenv("VIACEP_BASE", viacep.URL)
	os.Setenv("WEATHERAPI_BASE", weather.URL)
	os.Setenv("WEATHERAPI_KEY", "dummy")

	s := httpapi.NewServer()
	ts := httptest.NewServer(s.Router())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/weather?cep=99999999")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 404 {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
}

func TestSuccess(t *testing.T) {
	viacep := mockViaCEP(t)
	defer viacep.Close()
	weather := mockWeather(t)
	defer weather.Close()

	os.Setenv("VIACEP_BASE", viacep.URL)
	os.Setenv("WEATHERAPI_BASE", weather.URL)
	os.Setenv("WEATHERAPI_KEY", "dummy")

	s := httpapi.NewServer()
	ts := httptest.NewServer(s.Router())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/weather?cep=01001000")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var body struct {
		TempC float64 `json:"temp_C"`
		TempF float64 `json:"temp_F"`
		TempK float64 `json:"temp_K"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	// temp_c = 28.49 -> round1 = 28.5
	if body.TempC != 28.5 {
		t.Fatalf("TempC expected 28.5, got %v", body.TempC)
	}
	// F = 28.49*1.8+32 = 83.282 -> round1 = 83.3
	if body.TempF != 83.3 {
		t.Fatalf("TempF expected 83.3, got %v", body.TempF)
	}
	// K = 28.49 + 273 = 301.49 -> round1 = 301.5
	if body.TempK != 301.5 {
		t.Fatalf("TempK expected 301.5, got %v", body.TempK)
	}
}
