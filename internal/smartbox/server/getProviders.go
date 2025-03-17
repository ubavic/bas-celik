package server

import (
	"encoding/json"
	"io"
)

type GetProvidersInput struct{}

type GetProvidersPayload struct {
	Providers []Provider `json:"providers"`
}

type Provider struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (s *SmartBoxServer) handleGetProviders(data []byte, w io.Writer) error {
	msg := Message[GetProvidersInput]{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	providers := make([]Provider, 0, len(s.modulePaths))
	for _, module := range s.modulePaths {
		providers = append(providers, Provider{
			Id:   int(module.Vendor),
			Name: module.Vendor.String(),
		})
	}

	rsp := Response[GetProvidersPayload]{
		Operation: operationGetProviders,
		Payload: GetProvidersPayload{
			Providers: providers,
		},
	}

	return json.NewEncoder(w).Encode(rsp)
}
