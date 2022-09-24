package sensorpush

import (
	"context"
	"net/http"
	"time"
)

type StatusService service

type Status struct {
	Deployed time.Time
	Message  string
	MS       int
	Stack    string
	Status   StatusEnum
	Time     time.Time
	Version  string
}

type statusResponse struct {
	Deployed string `json:"deployed"`
	Message  string `json:"message"`
	MS       int    `json:"ms"`
	Stack    string
	Status   string
	Time     string
	Version  string
}

func (s *StatusService) Get(ctx context.Context) (*Status, error) {
	s0 := &Status{}

	req, err := s.c.NewBaseRequest(ctx, http.MethodPost, "", nil)
	if err != nil {
		return s0, err
	}

	sresp := statusResponse{}
	_, err = s.c.Do(req, &sresp)
	if err != nil {
		return s0, err
	}

	depT, err := parseTimeStatus(sresp.Deployed)
	if err != nil {
		return s0, err
	}

	srvT, err := parseTimeStatus(sresp.Time)
	if err != nil {
		return s0, err
	}

	st := &Status{
		Message:  sresp.Message,
		Deployed: depT,
		MS:       sresp.MS,
		Stack:    sresp.Stack,
		Status:   newStatusEnum(sresp.Status),
		Time:     srvT,
		Version:  sresp.Version,
	}
	return st, nil
}
