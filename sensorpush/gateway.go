package sensorpush

import ()

type Gateway struct {
	ID   string
	Name string
}

type GatewaySlice []*Gateway

func (s GatewaySlice) Len() int {
	return len(s)
}

func (s GatewaySlice) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s GatewaySlice) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

func newGateway(r gatewayResponse) *Gateway {
	return &Gateway{
		ID:   r.ID,
		Name: r.Name,
	}
}

type gatewaysRequest struct {
}

type gatewayResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type gatewaysResponse map[string]gatewayResponse
