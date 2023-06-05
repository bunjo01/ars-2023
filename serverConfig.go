package main

import (
	tracer "ars-2023/tracer"
	"context"
	"fmt"
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
	span := tracer.StartSpanFromRequest("createConfigHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling config create at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	err := cs.CheckRequest(req, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	rt, err := DecodeFreeConfig(req.Body, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	config, err := cs.store.Config(rt, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	RenderJSON(w, config, ctx)
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
	span := tracer.StartSpanFromRequest("getAllConfigHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get all configs at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	allTasks, err := cs.store.GetAllConfigs(ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	RenderJSON(w, allTasks, ctx)
}

// swagger:route GET /config/{id}/all/ Configuration getAllConfigVersions
// Get all configuration versions
//
// responses:
//
//	404: ErrorResponse
//	200: []FreeConfig
func (cs *configServer) getConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigVersionsHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get config versions at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	task, err := cs.store.GetConfigVersions(id, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	RenderJSON(w, task, ctx)
}

// swagger:route GET /config/{id}/{version}/ Configuration getConfig
// Get specific configuration
//
// responses:
//
//	404: ErrorResponse
//	200: FreeConfig
func (cs *configServer) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get config at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := cs.store.GetConfig(id, version, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	RenderJSON(w, task, ctx)
}

// swagger:route DELETE /config/{id}/all/ Configuration deleteConfigVersions
// Delete all configuration versions
//
// responses:
//
//	404: ErrorResponse
//	201: []FreeConfig
func (cs *configServer) delConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("delConfigVersionsHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del config versions at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	task, err := cs.store.DeleteConfigVersions(id, ctx)
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

// swagger:route DELETE /config/{id}/{version}/ Configuration deleteConfig
// Delete specific configuration
//
// responses:
//
//	404: ErrorResponse
//	201: FreeConfig
func (cs *configServer) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("delConfigHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del config at: %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	msg, err := cs.store.DeleteConfig(id, version, ctx)
	if err != nil {
		throwError(w, err)
		return
	}
	RenderJSON(w, msg, ctx)
}

// Swagger routing handler

func (cs *configServer) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}
