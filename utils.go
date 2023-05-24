package main

import (
	cs "ars-2023/configdatabase"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

// JSON decoders

func decodeFreeConfig(r io.Reader) (*cs.FreeConfig, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt cs.FreeConfig
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}
func decodeFreeGroup(r io.Reader) (*cs.FreeGroup, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt cs.FreeGroup
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

// JSON render

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// TODO implement into repos instead of staticId block
// ID generator for static IDs

// Errors

func throwNotFoundError(w http.ResponseWriter) {
	err := errors.New("ID not found")
	http.Error(w, err.Error(), http.StatusNotFound)
}

func throwForbiddenError(w http.ResponseWriter) {
	err := errors.New("already exists")
	http.Error(w, err.Error(), http.StatusForbidden)
}

func throwTeapot(w http.ResponseWriter) {
	err := errors.New("The server refuses the attempt to brew coffee with a teapot.")
	http.Error(w, err.Error(), http.StatusTeapot)
}

// Http request validator

func checkRequest(req *http.Request, w http.ResponseWriter) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expected application/json Content-type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
}

// Swagger routing handler

func (ts *configServer) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}
