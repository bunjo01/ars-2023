package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (ts *dbServerConfig) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	checkRequest(req, w)

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
		throwNotFoundError(w)
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
		throwNotFoundError(w)
	}
}
