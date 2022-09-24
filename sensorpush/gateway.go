package sensorpush

import (
	"time"
)

type Gateway struct {
	LastAlert time.Time
	LastSeen  time.Time
	ID        string
	Message   string
	Name      string
	Paired    bool
	// TODO: Tags
	Version string
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
		ID:      r.ID,
		Message: r.Message,
		Name:    r.Name,
		Paired:  r.Paired,
		Version: r.Version,
	}

	// Last Alert
	if r.LastAlert != "" {
		t, err := parseTime(r.LastAlert)
		if err != nil {
			return nil, err
		}
		g.LastAlert = t
	}

	// Last Seen
	if r.LastSeen != "" {
		t, err := parseTime(r.LastSeen)
		if err != nil {
			return nil, err
		}
		g.LastSeen = t
	}

	return g, nil
}

type gatewaysRequest struct {
}

type gatewayResponse struct {
	LastAlert string `json:"last_alert"`
	LastSeen  string `json:"last_seen"`
	ID        string `json:"id"`
	Message   string `json:"message"`
	Name      string `json:"name"`
	Paired    bool   `json:"paired"`
	Version   string `json:"version"`
}

type gatewaysResponse map[string]gatewayResponse
