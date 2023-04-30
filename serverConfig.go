package main

import (
	"errors"
	"github.com/gorilla/mux"
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

func (ts *dbServerConfig) getAllHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*Config{}
	for _, v := range ts.data {
		allTasks = append(allTasks, v)
	}

	renderJSON(w, allTasks)
}

func (ts *dbServerConfig) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.data[id]
	if !ok {
		err := errors.New("id not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (ts *dbServerConfig) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if v, ok := ts.data[id]; ok {
		delete(ts.data, id)
		renderJSON(w, v)
	} else {
		err := errors.New("id not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
