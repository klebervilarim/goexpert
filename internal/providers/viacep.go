package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrCEPNotFound = errors.New("cep not found")

type ViaCEPClient struct {
	http *http.Client
	base string
}

func NewViaCEPClient(h *http.Client, base string) *ViaCEPClient {
	return &ViaCEPClient{http: h, base: base}
}

type viaCEPResp struct {
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Erro       bool   `json:"erro"`
}

type Address struct {
	Localidade string
	UF         string
}

func (c *ViaCEPClient) Lookup(ctx context.Context, cep string) (*Address, error) {
	url := fmt.Sprintf("%s/ws/%s/json", c.base, cep)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrCEPNotFound
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("viacep status %d", res.StatusCode)
	}

	var vr viaCEPResp
	if err := json.NewDecoder(res.Body).Decode(&vr); err != nil {
		return nil, err
	}
	if vr.Erro || vr.Localidade == "" || vr.UF == "" {
		return nil, ErrCEPNotFound
	}
	return &Address{Localidade: vr.Localidade, UF: vr.UF}, nil
}
