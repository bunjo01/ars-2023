package main

import (
	"context"
	"net/http"
	"testing"
)

func mockConfigServer() *configServer {
	x := configServer{
		store:  nil,
		Keys:   make(map[string]string),
		tracer: nil,
		closer: nil,
	}
	x.Keys["used"] = "used"
	return &x
}

func TestIsFirst(t *testing.T) {
	req := http.Request{Header: map[string][]string{}}
	req.Header.Add("Idempotency-key", "idempotent")
	serv := mockConfigServer()
	actual, _ := IsFirst(&req, serv, context.Background())
	if !actual {
		t.Error("Test failed. Expected: true, but got: false")
	}
}

func TestIsFirstUsed(t *testing.T) {
	req := http.Request{Header: map[string][]string{}}
	req.Header.Add("Idempotency-key", "used")
	serv := mockConfigServer()
	_, err := IsFirst(&req, serv, context.Background())
	if err == nil {
		t.Error("Test failed. Expected: conflict, but got: nil")
	}
}

func TestIsFirstEmpty(t *testing.T) {
	req := http.Request{Header: map[string][]string{}}
	req.Header.Add("Idempotency-key", "")
	serv := mockConfigServer()
	actual, err := IsFirst(&req, serv, context.Background())
	if !(!actual && err != nil) {
		t.Error("Test failed. Expected: teapot, but got: coffee")
	}
}
