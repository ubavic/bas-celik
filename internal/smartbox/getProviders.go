package smartbox

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

	providers := make([]Provider, 0, len(s.loadedVendors))
	for _, vendor := range s.loadedVendors {
		providers = append(providers, Provider{
			Id:   int(vendor),
			Name: vendor.String(),
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
