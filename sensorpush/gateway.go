package sensorpush

import (
	"time"
)

type Gateway struct {
	LastAlert time.Time
	ID        string
	Name      string
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

func newGateway(r gatewayResponse) (*Gateway, error) {
	g := &Gateway{
		ID:   r.ID,
		Name: r.Name,
	}

	// Last Alert
	if r.LastAlert != "" {
		t, err := parseTime(r.LastAlert)
		if err != nil {
			return nil, err
		}
		g.LastAlert = t
	}

	return g, nil
}

type gatewaysRequest struct {
}

type gatewayResponse struct {
	LastAlert string `json:"last_alert"`
	ID        string `json:"id"`
	Name      string `json:"name"`
}

type gatewaysResponse map[string]gatewayResponse
