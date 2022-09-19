package sensorpush

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type StatusService service

type StatusEnum int

const (
	StatusUnknown StatusEnum = iota
	StatusOK
)

func (s StatusEnum) String() string {
	switch s {
	case StatusOK:
		return "ok"
	}
	return "unknown"
}

func newStatusEnum(s string) StatusEnum {
	switch s {
	case "ok":
		return StatusOK
	}
	return StatusUnknown
}

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

	fmt.Printf("sresp => %+v\n", sresp)

	depT, err := parseTime(sresp.Deployed)
	if err != nil {
		return s0, err
	}

	srvT, err := parseTime(sresp.Time)
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
