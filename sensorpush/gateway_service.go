package sensorpush

import (
	"context"
	//"fmt"
	"net/http"
	"sort"
)

type GatewayService service

// List returns the gateways
func (s *GatewayService) List(ctx context.Context) (GatewaySlice, error) {
	var g0 []*Gateway

	sreq := gatewaysRequest{}

	req, err := s.c.NewRequest(ctx, http.MethodPost, "devices/gateways", sreq)
	if err != nil {
		return g0, err
	}

	gsresp := gatewaysResponse{}
	_, err = s.c.Do(req, &gsresp)
	if err != nil {
		return g0, err
	}

	gateways := make(GatewaySlice, 0, len(gsresp))
	for _, gresp := range gsresp {
		gateways = append(gateways, newGateway(gresp))
	}

	sort.Sort(gateways)

	return gateways, nil
}
