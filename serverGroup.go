package main

import (
	"errors"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

func (ts *dbServerConfig) createGroupHandler(w http.ResponseWriter, req *http.Request) {
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

	rt, err := decodeGroupBody(req.Body)
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
	ts.dataGroup[tsId] = rt
	renderJSON(w, rt)

	log.Println(rt.Id)
}

func (ts *dbServerConfig) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*Group{}
	for _, v := range ts.dataGroup {
		allTasks = append(allTasks, v)
	}
	renderJSON(w, allTasks)
}

func (ts *dbServerConfig) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, ok := ts.dataGroup[id+version]
	if !ok {
		err := errors.New("id not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (ts *dbServerConfig) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	if v, ok := ts.dataGroup[id+version]; ok {
		delete(ts.dataGroup, id+version)
		renderJSON(w, v)
	} else {
		err := errors.New("id not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *dbServerConfig) appendGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
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
	rt, err := decodeAppendBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if v, ok := ts.dataGroup[id+version]; ok {
		for _, el := range rt.Appends {
			if ts.data[*el] != nil {
				v.Configs[*el] = el
			}
		}
		renderJSON(w, v)
	} else {
		err := errors.New("id not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
