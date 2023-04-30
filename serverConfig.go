package main

import (
	"errors"
	"log"
	"mime"
	"net/http"
)

type dbServerConfig struct {
	data map[string]*Config
}

func (ts *dbServerConfig) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expected application/json Content-type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate random UUID on creation
	id := createId()
	rt.Id = id
	ts.data[id] = rt
	renderJSON(w, rt)

	// test for the sent data
	log.Println(rt.Id)

}
