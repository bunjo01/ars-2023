package main

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"mime"
	"net/http"
)

func decodeBody(r io.Reader) (*Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil

}
func decodeGroupBody(r io.Reader) (*Group, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt Group
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeAppendBody(r io.Reader) (*DTOConfig, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt DTOConfig
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()
}

func throwNotFoundError(w http.ResponseWriter) {
	err := errors.New("id not found")
	http.Error(w, err.Error(), http.StatusNotFound)
}

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
