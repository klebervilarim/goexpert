package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"regexp"
	"time"

	"cloudrun/internal/providers"

	"cloudrun/internal/util"
)

type Server struct {
	viacep    *providers.ViaCEPClient
	weather   *providers.WeatherAPIClient
	router    *http.ServeMux
	cepRegexp *regexp.Regexp
}

func NewServer() *Server {
	viacepBase := getenv("VIACEP_BASE", "https://viacep.com.br")
	weatherBase := getenv("WEATHERAPI_BASE", "https://api.weatherapi.com")
	weatherKey := os.Getenv("WEATHERAPI_KEY")

	httpClient := &http.Client{Timeout: 6 * time.Second}

	return &Server{
		viacep:    providers.NewViaCEPClient(httpClient, viacepBase),
		weather:   providers.NewWeatherAPIClient(httpClient, weatherBase, weatherKey),
		router:    http.NewServeMux(),
		cepRegexp: regexp.MustCompile(`^\d{8}$`),
	}
}

func (s *Server) Router() *http.ServeMux {
	s.router.HandleFunc("GET /healthz", s.health)
	s.router.HandleFunc("GET /weather", s.handleWeather)
	return s.router
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type weatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func (s *Server) handleWeather(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if !s.cepRegexp.MatchString(cep) {
		httpError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	addr, err := s.viacep.Lookup(ctx, cep)
	if err != nil {
		if errors.Is(err, providers.ErrCEPNotFound) {
			httpError(w, http.StatusNotFound, "can not find zipcode")
			return
		}
		httpError(w, http.StatusBadGateway, "upstream error")
		return
	}

	// Ex.: "São Paulo,SP,BR"
	query := addr.Localidade + "," + addr.UF + ",BR"
	tempC, err := s.weather.CurrentTempC(ctx, query)
	if err != nil {
		httpError(w, http.StatusBadGateway, "upstream error")
		return
	}

	resp := weatherResponse{
		TempC: util.Round1(tempC),
		TempF: util.Round1(util.CtoF(tempC)),
		TempK: util.Round1(util.CtoK_Custom(tempC)), // +273 conforme requisito
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(w, msg, code)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
