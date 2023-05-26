package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

// swagger:route POST /config/ Configuration createConfig
// Create new configuration
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: FreeConfig
func (cs *configServer) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	err := cs.CheckRequest(req, w)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	rt, err := DecodeFreeConfig(req.Body)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	config, err := cs.store.Config(rt)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	RenderJSON(w, config)
}

// swagger:route GET /config/all/ Configuration getAllConfigs
// Get all configurations
//
// responses:
//
//	404: ErrorResponse
//	418: Teapot
//	200: []FreeConfig
func (cs *configServer) getAllConfigHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := cs.store.GetAllConfigs()
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	RenderJSON(w, allTasks)
}

// swagger:route GET /config/{id}/all/ Configuration getAllConfigVersions
// Get all configuration versions
//
// responses:
//
//	404: ErrorResponse
//	200: []FreeConfig
func (cs *configServer) getConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, err := cs.store.GetConfigVersions(id)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	RenderJSON(w, task)
}

// swagger:route GET /config/{id}/{version}/ Configuration getConfig
// Get specific configuration
//
// responses:
//
//	404: ErrorResponse
//	200: FreeConfig
func (cs *configServer) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := cs.store.GetConfig(id, version)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	RenderJSON(w, task)
}

// swagger:route DELETE /config/{id}/all/ Configuration deleteConfigVersions
// Delete all configuration versions
//
// responses:
//
//	404: ErrorResponse
//	201: []FreeConfig
func (cs *configServer) delConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, err := cs.store.DeleteConfigVersions(id)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	if len(task) > 0 {
		RenderJSON(w, task)
	} else {
		ThrowNotFoundError(w)
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
func (cs *configServer) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	msg, err := cs.store.DeleteConfig(id, version)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	RenderJSON(w, msg)
}

// Swagger routing handler

func (cs *configServer) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}
