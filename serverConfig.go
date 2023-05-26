package main

import (
	cs "ars-2023/configdatabase"
	util "ars-2023/utilities"
	"github.com/gorilla/mux"
	"net/http"
)

type configServer struct {
	store *cs.ConfigStore
}

// swagger:route POST /config/ Configuration createConfig
// Create new configuration
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: FreeConfig
func (ts *configServer) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	util.CheckRequest(req, w)
	rt, err := util.DecodeFreeConfig(req.Body)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	config, err := ts.store.Config(rt)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	util.RenderJSON(w, config)
}

// swagger:route GET /config/all/ Configuration getAllConfigs
// Get all configurations
//
// responses:
//
//	404: ErrorResponse
//	418: Teapot
//	200: []FreeConfig
func (ts *configServer) getAllConfigHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := ts.store.GetAllConfigs()
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	util.RenderJSON(w, allTasks)
}

// swagger:route GET /config/{id}/all/ Configuration getAllConfigVersions
// Get all configuration versions
//
// responses:
//
//	404: ErrorResponse
//	200: []FreeConfig
func (ts *configServer) getConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, err := ts.store.GetConfigVersions(id)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	util.RenderJSON(w, task)
}

// swagger:route GET /config/{id}/{version}/ Configuration getConfig
// Get specific configuration
//
// responses:
//
//	404: ErrorResponse
//	200: FreeConfig
func (ts *configServer) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := ts.store.GetConfig(id, version)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	util.RenderJSON(w, task)
}

// swagger:route DELETE /config/{id}/all/ Configuration deleteConfigVersions
// Delete all configuration versions
//
// responses:
//
//	404: ErrorResponse
//	201: []FreeConfig
func (ts *configServer) delConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, err := ts.store.DeleteConfigVersions(id)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	if len(task) > 0 {
		util.RenderJSON(w, task)
	} else {
		util.ThrowNotFoundError(w)
		return
	}
}

// swagger:route DELETE /config/{id}/{version}/ Configuration deleteConfig
// Delete specific configuration
//
// responses:
//
//	404: ErrorResponse
//	201: FreeConfig
func (ts *configServer) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	msg, err := ts.store.DeleteConfig(id, version)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	util.RenderJSON(w, msg)
}

// Swagger routing handler

func (ts *configServer) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}
