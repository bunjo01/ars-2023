package main

import (
	"ars-2023/tracer"
	"context"
	"fmt"
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
	span := tracer.StartSpanFromRequest("createGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling group create at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	err := cs.CheckRequest(req, ctx)
	if err != nil {
		throwError(w, err)
		return
	}

	rt, err := DecodeFreeGroup(req.Body, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	group, err := cs.store.Group(rt, true, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	RenderJSON(w, group, ctx)
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
	span := tracer.StartSpanFromRequest("getAllGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get all groups at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	task, err := cs.store.GetAllGroups(ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	if len(task) > 0 {
		RenderJSON(w, task, ctx)
	} else {
		throwError(w, err)
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
	span := tracer.StartSpanFromRequest("getGroupVersionsHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get group versions at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	task, err := cs.store.GetGroupVersions(id, ctx)
	if err != nil || len(task) < 1 {
		throwError(w, err)
		return
	}
	RenderJSON(w, task, ctx)
}

// swagger:route GET /group/{id}/{version}/ Group getGroup
// Get specific group
//
// responses:
//
//	404: ErrorResponse
//	200: FreeGroup
func (cs *configServer) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get group at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := cs.store.GetGroup(id, version, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	RenderJSON(w, task, ctx)
}

// swagger:route DELETE /group/{id}/all/ Group deleteGroupVersions
// Delete all group versions
//
// responses:
//
//	404: ErrorResponse
//	201: []FreeGroup
func (cs *configServer) delGroupVersionsHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("delGroupVersionsHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del group versions at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	task, err := cs.store.DeleteGroupVersions(id, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	if len(task) > 0 {
		RenderJSON(w, task, ctx)
	} else {
		throwError(w, err)
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
	span := tracer.StartSpanFromRequest("delGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del group at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := cs.store.DeleteGroup(id, version, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	if len(task) > 0 {
		RenderJSON(w, task, ctx)
	} else {
		throwError(w, err)
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
	span := tracer.StartSpanFromRequest("appendGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling group append at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	er := cs.CheckRequest(req, ctx)
	if er != nil {
		http.Error(w, er.Message, er.Status)
		return
	}
	rt, err := DecodeFreeGroup(req.Body, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	group, err := cs.store.Group(rt, false, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	RenderJSON(w, group, ctx)
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
	span := tracer.StartSpanFromRequest("getLabelGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get group configs by label at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]
	task, err := cs.store.GetConfigsByLabels(id, version, labels, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	RenderJSON(w, task, ctx)
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
	span := tracer.StartSpanFromRequest("delLabelGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del group configs by label at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]
	newVersion := mux.Vars(req)["new"]

	task, err := cs.store.DeleteConfigsByLabels(id, version, labels, newVersion, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	RenderJSON(w, task, ctx)
}
