package main

import (
	"errors"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

func (ts *dbServerConfig) createConfigHandler(w http.ResponseWriter, req *http.Request) {
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

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate random UUID on creation
	id := createId()
	// TESTING
	staticId := rt.Id
	if staticId == "1234" {
		id = staticId
	}

	rt.Id = id
	tsId := id + rt.Version
	ts.data[tsId] = rt
	renderJSON(w, rt)

	// test for the sent data
	log.Println(ts.data)

}

func (ts *dbServerConfig) getAllConfigHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*Config{}
	for _, v := range ts.data {
		allTasks = append(allTasks, v)
	}

	renderJSON(w, allTasks)
}

func (ts *dbServerConfig) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, ok := ts.data[id+version]
	if !ok {
		err := errors.New("id not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (ts *dbServerConfig) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	if v, ok := ts.data[id+version]; ok {
		delete(ts.data, id+version)
		renderJSON(w, v)
	} else {
		err := errors.New("id not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
