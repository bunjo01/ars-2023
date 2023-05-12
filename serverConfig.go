package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

// swagger:route POST /config/ config createConfig
// Add new config
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseConfig
func (ts *dbServerConfig) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	checkRequest(req, w)

	rt, err := decodeFreeConfig(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := createId(rt.Id)
	rt.Id = id
	el := rt.freeToDBConfig()
	if ts.data[el.Id] != nil {
		throwForbiddenError(w)
		return
	}
	ts.data[el.Id] = el
	renderJSON(w, el)
}

func (ts *dbServerConfig) getAllConfigHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*FreeConfig{}
	for _, v := range ts.data {
		allTasks = append(allTasks, v.dBConfigToFree())
	}
	renderJSON(w, allTasks)
}

func (ts *dbServerConfig) getConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*FreeConfig{}
	for _, v := range ts.data {
		conf := v.dBConfigToFree()
		if conf.Id == mux.Vars(req)["id"] {
			allTasks = append(allTasks, conf)
		}
	}
	renderJSON(w, allTasks)
}

func (ts *dbServerConfig) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"] + separator() + mux.Vars(req)["version"]
	task, ok := ts.data[id]
	if !ok {
		throwNotFoundError(w)
		return
	}
	renderJSON(w, task.dBConfigToFree())
}

func (ts *dbServerConfig) delConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	allTasks := []*FreeConfig{}
	for _, v := range ts.data {
		conf := v.dBConfigToFree()
		if conf.Id == id {
			allTasks = append(allTasks, conf)
			delete(ts.data, v.Id)
		}
	}
	if len(allTasks) > 0 {
		renderJSON(w, allTasks)
	} else {
		throwNotFoundError(w)
	}
}

func (ts *dbServerConfig) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"] + separator() + mux.Vars(req)["version"]
	if v, ok := ts.data[id]; ok {
		delete(ts.data, id)
		renderJSON(w, v.dBConfigToFree())
	} else {
		throwNotFoundError(w)
	}
}
