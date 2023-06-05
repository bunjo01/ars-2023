package main

import (
	"bytes"
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

func TestCheckRequest(t *testing.T) {
	req := http.Request{Header: map[string][]string{}}
	req.Header.Add("Idempotency-key", "value")
	req.Header.Add("Content-Type", "application/json")
	serv := mockConfigServer()

	actual := serv.CheckRequest(&req, context.Background())
	if actual != nil {
		t.Errorf("Test failed. Expected: nil, but got: %s", actual.Error())
	}
}

func TestDecodeFreeConfig(t *testing.T) {
	jsonBody := []byte(`{
			"id": "1234",
    		"version": "v3.0",
			"entries": {"entry1": "1"}
			}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/", bodyReader)
	_, err := DecodeFreeConfig(req.Body, context.Background())
	if err != nil {
		t.Errorf("Test failed. Expected: FreeConfig, but got: %s", err.Error())
	}
}

func TestDecodeFreeGroup(t *testing.T) {
	jsonBody := []byte(`{
		"id": "1234",
		"version": "v1.0",
		"configs": {"2222": {
				"id": "2222",
				"labels": {"lol": "lol"},
				"entries": {"lol2": "lol2"}}}}`,
	)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/", bodyReader)
	_, err := DecodeFreeGroup(req.Body, context.Background())
	if err != nil {
		t.Errorf("Test failed. Expected: FreeGroup, but got: %s", err.Error())
	}
}
