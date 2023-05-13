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
func (ts *dbServerConfig) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	checkRequest(req, w)

	rt, err := decodeFreeGroup(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := createId(rt.Id)
	rt.Id = id
	el := rt.freeToDBGroup()
	if ts.dataGroup[el.Id] != nil {
		throwForbiddenError(w)
		return
	}
	ts.dataGroup[el.Id] = el
	renderJSON(w, rt)
}

// swagger:route GET /group/all/ Group getAllGroups
// Get all groups
//
// responses:
//
//	404: ErrorResponse
//	418: Teapot
//	200: []FreeGroup
func (ts *dbServerConfig) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*FreeGroup{}
	for _, v := range ts.dataGroup {
		allTasks = append(allTasks, v.dBGroupToFree())
	}
	renderJSON(w, allTasks)
}

// swagger:route GET /group/{id}/all/ Group getAllGroupVersions
// Get all group versions
//
// responses:
//
//	404: ErrorResponse
//	200: []FreeGroup
func (ts *dbServerConfig) getGroupVersionsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*FreeGroup{}
	for _, v := range ts.dataGroup {
		gro := v.dBGroupToFree()
		if gro.Id == mux.Vars(req)["id"] {
			allTasks = append(allTasks, gro)
		}
	}
	renderJSON(w, allTasks)
}

// swagger:route GET /group/{id}/{version}/ Group getGroup
// Get specific group
//
// responses:
//
//	404: ErrorResponse
//	200: FreeGroup
func (ts *dbServerConfig) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"] + separator() + mux.Vars(req)["version"]
	task, ok := ts.dataGroup[id]
	if !ok {
		throwNotFoundError(w)
		return
	}
	renderJSON(w, task.dBGroupToFree())
}

// swagger:route DELETE /group/{id}/all/ Group deleteGroupVersions
// Delete all group versions
//
// responses:
//
//	404: ErrorResponse
//	201: []FreeGroup
func (ts *dbServerConfig) delGroupVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	allTasks := []*FreeGroup{}
	for _, v := range ts.dataGroup {
		gro := v.dBGroupToFree()
		if gro.Id == id {
			allTasks = append(allTasks, gro)
			delete(ts.dataGroup, v.Id)
		}
	}
	if len(allTasks) > 0 {
		renderJSON(w, allTasks)
	} else {
		throwNotFoundError(w)
	}
}

// swagger:route DELETE /group/{id}/{version}/ Group deleteGroup
// Delete specific group
//
// responses:
//
//	404: ErrorResponse
//	201: FreeGroup
func (ts *dbServerConfig) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"] + separator() + mux.Vars(req)["version"]
	if v, ok := ts.dataGroup[id]; ok {
		delete(ts.dataGroup, id)
		renderJSON(w, v.dBGroupToFree())
	} else {
		throwNotFoundError(w)
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
func (ts *dbServerConfig) appendGroupHandler(w http.ResponseWriter, req *http.Request) {
	checkRequest(req, w)
	id := mux.Vars(req)["id"]
	oldVersion := mux.Vars(req)["version"]
	newVersion := mux.Vars(req)["new"]
	oldData := id + separator() + oldVersion
	newData := id + separator() + newVersion
	rt, err := decodeGroupConfigs(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if v, ok := ts.dataGroup[oldData]; ok {
		var newGroup DBGroup
		newGroup.Configs = make(map[string]*DBConfig)
		newGroup.Id = newData
		for _, old := range v.Configs {
			x := old.dBToGroupConfig()
			x.Id = createId(x.Id)
			appen := x.groupConfigToDBConfig()
			newGroup.Configs[appen.Id] = appen
		}
		for _, val := range rt.Configs {
			dbg := val.groupConfigToDBConfig()
			if newGroup.Configs[dbg.Id] == nil {
				newGroup.Configs[dbg.Id] = dbg
			}
		}
		if ts.dataGroup[newGroup.Id] != nil {
			throwForbiddenError(w)
		} else {
			ts.dataGroup[newGroup.Id] = &newGroup
			renderJSON(w, newGroup.dBGroupToFree())
		}
	} else {
		throwNotFoundError(w)
	}
}

// swagger:route GET /group/{id}/{version}/{labels}/ Label getConfigsByLabel
// Get configs by label
//
// responses:
//
//	404: ErrorResponse
//	418: Teapot
//	200: []GroupConfig
func (ts *dbServerConfig) getConfigsByLabel(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"] + separator() + mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]
	allTasks := []*GroupConfig{}
	if v, ok := ts.dataGroup[id]; ok {
		for _, val := range v.Configs {
			if val.compareLabels(labels) {
				allTasks = append(allTasks, val.dBToGroupConfig())
			}
		}
	} else {
		throwNotFoundError(w)
	}
	if len(allTasks) == 0 {
		throwTeapot(w)
	} else {
		renderJSON(w, allTasks)
	}
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
func (ts *dbServerConfig) delConfigsByLabel(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	oldVersion := mux.Vars(req)["version"]
	newVersion := mux.Vars(req)["new"]
	oldData := id + separator() + oldVersion
	newData := id + separator() + newVersion
	labels := mux.Vars(req)["labels"]

	if v, ok := ts.dataGroup[oldData]; ok {
		var newGroup DBGroup
		newGroup.Configs = make(map[string]*DBConfig)
		newGroup.Id = newData
		for _, old := range v.Configs {
			if !old.compareLabels(labels) {
				x := old.dBToGroupConfig()
				x.Id = createId(x.Id)
				appen := x.groupConfigToDBConfig()
				newGroup.Configs[appen.Id] = appen
			}
		}
		if len(v.Configs) == len(newGroup.Configs) {
			throwTeapot(w)
		} else if ts.dataGroup[newGroup.Id] != nil {
			throwForbiddenError(w)
		} else {
			ts.dataGroup[newGroup.Id] = &newGroup
			renderJSON(w, newGroup.dBGroupToFree())
		}
	} else {
		throwNotFoundError(w)
	}
}
