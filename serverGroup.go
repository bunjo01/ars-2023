package main

import (
	util "ars-2023/utilities"
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
func (ts *configServer) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	util.CheckRequest(req, w)

	rt, err := util.DecodeFreeGroup(req.Body)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	group, err := ts.store.Group(rt, true)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	util.RenderJSON(w, group)
}

// swagger:route GET /group/all/ Group getAllGroups
// Get all groups
//
// responses:
//
//	400: ErrorResponse
//	418: Teapot
//	200: []FreeGroup
func (ts *configServer) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	task, err := ts.store.GetAllGroups()
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	if len(task) > 0 {
		util.RenderJSON(w, task)
	} else {
		util.ThrowTeapot(w)
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
func (ts *configServer) getGroupVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, err := ts.store.GetGroupVersions(id)
	if err != nil || len(task) < 1 {
		util.ThrowNotFoundError(w)
		return
	}
	util.RenderJSON(w, task)
}

// swagger:route GET /group/{id}/{version}/ Group getGroup
// Get specific group
//
// responses:
//
//	404: ErrorResponse
//	200: FreeGroup
func (ts *configServer) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := ts.store.GetGroup(id, version)
	if err != nil {
		util.ThrowNotFoundError(w)
		return
	}
	util.RenderJSON(w, task)
}

// swagger:route DELETE /group/{id}/all/ Group deleteGroupVersions
// Delete all group versions
//
// responses:
//
//	404: ErrorResponse
//	201: []FreeGroup
func (ts *configServer) delGroupVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, err := ts.store.DeleteGroupVersions(id)
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

// swagger:route DELETE /group/{id}/{version}/ Group deleteGroup
// Delete specific group
//
// responses:
//
//	404: ErrorResponse
//	201: FreeGroup
func (ts *configServer) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := ts.store.DeleteGroup(id, version)
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
func (ts *configServer) appendGroupHandler(w http.ResponseWriter, req *http.Request) {
	util.CheckRequest(req, w)

	rt, err := util.DecodeFreeGroup(req.Body)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	group, err := ts.store.Group(rt, false)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	util.RenderJSON(w, group)
}

// swagger:route GET /group/{id}/{version}/{labels}/ Label getConfigsByLabel
// Get configs by label
//
// responses:
//
//	404: ErrorResponse
//	418: Teapot
//	200: []GroupConfig
func (ts *configServer) getConfigsByLabel(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]
	task, err := ts.store.GetConfigsByLabels(id, version, labels)
	if err != nil {
		util.ThrowNotFoundError(w)
		return
	}
	util.RenderJSON(w, task)
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
func (ts *configServer) delConfigsByLabel(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]
	newVersion := mux.Vars(req)["new"]

	task, err := ts.store.DeleteConfigsByLabels(id, version, labels, newVersion)
	if err != nil {
		util.ThrowBadRequest(w)
		return
	}
	util.RenderJSON(w, task)
}
