package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

// swagger:route POST /group/ Group createGroup
// Create new group
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: FreeGroup
func (cs *configServer) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	err := cs.CheckRequest(req, w)
	if err != nil {
		ThrowBadRequest(w)
		return
	}

	rt, err := DecodeFreeGroup(req.Body)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	group, err := cs.store.Group(rt, true)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	RenderJSON(w, group)
}

// swagger:route GET /group/all/ Group getAllGroups
// Get all groups
//
// responses:
//
//	400: ErrorResponse
//	418: Teapot
//	200: []FreeGroup
func (cs *configServer) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	task, err := cs.store.GetAllGroups()
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	if len(task) > 0 {
		RenderJSON(w, task)
	} else {
		ThrowTeapot(w)
		return
	}
}

// swagger:route GET /group/{id}/all/ Group getAllGroupVersions
// Get all group versions
//
// responses:
//
//	404: ErrorResponse
//	200: []FreeGroup
func (cs *configServer) getGroupVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, err := cs.store.GetGroupVersions(id)
	if err != nil || len(task) < 1 {
		ThrowNotFoundError(w)
		return
	}
	RenderJSON(w, task)
}

// swagger:route GET /group/{id}/{version}/ Group getGroup
// Get specific group
//
// responses:
//
//	404: ErrorResponse
//	200: FreeGroup
func (cs *configServer) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := cs.store.GetGroup(id, version)
	if err != nil {
		ThrowNotFoundError(w)
		return
	}
	RenderJSON(w, task)
}

// swagger:route DELETE /group/{id}/all/ Group deleteGroupVersions
// Delete all group versions
//
// responses:
//
//	404: ErrorResponse
//	201: []FreeGroup
func (cs *configServer) delGroupVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, err := cs.store.DeleteGroupVersions(id)
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

// swagger:route DELETE /group/{id}/{version}/ Group deleteGroup
// Delete specific group
//
// responses:
//
//	404: ErrorResponse
//	201: FreeGroup
func (cs *configServer) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := cs.store.DeleteGroup(id, version)
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

// swagger:route POST /group/{id}/{version}/{new}/ Group appendGroup
// Create new group with appended config
//
// responses:
//
//	415: ErrorResponse
//	404: ErrorResponse
//	403: ErrorResponse
//	400: ErrorResponse
//	201: FreeGroup
func (cs *configServer) appendGroupHandler(w http.ResponseWriter, req *http.Request) {
	err := cs.CheckRequest(req, w)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	rt, err := DecodeFreeGroup(req.Body)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	group, err := cs.store.Group(rt, false)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	RenderJSON(w, group)
}

// swagger:route GET /group/{id}/{version}/{labels}/ Label getConfigsByLabel
// Get configs by label
//
// responses:
//
//	404: ErrorResponse
//	418: Teapot
//	200: []GroupConfig
func (cs *configServer) getConfigsByLabel(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]
	task, err := cs.store.GetConfigsByLabels(id, version, labels)
	if err != nil {
		ThrowNotFoundError(w)
		return
	}
	RenderJSON(w, task)
}

// swagger:route DELETE /group/{id}/{version}/{new}/{labels}/ Label delConfigsByLabel
// Delete configs by label
//
// responses:
//
//	404: ErrorResponse
//	403: ErrorResponse
//	418: Teapot
//	200: FreeGroup
func (cs *configServer) delConfigsByLabel(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]
	newVersion := mux.Vars(req)["new"]

	task, err := cs.store.DeleteConfigsByLabels(id, version, labels, newVersion)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	RenderJSON(w, task)
}
