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
func (ts *configServer) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	checkRequest(req, w)

	rt, err := decodeFreeGroup(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	group, err := ts.store.Group(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}

// swagger:route GET /group/all/ Group getAllGroups
// Get all groups
//
// responses:
//
//	404: ErrorResponse
//	418: Teapot
//	200: []FreeGroup
func (ts *configServer) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	task, err := ts.store.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(task) > 0 {
		renderJSON(w, task)
	} else {
		throwNotFoundError(w)
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
	if err != nil {
		throwNotFoundError(w)
	}
	if len(task) > 0 {
		renderJSON(w, task)
	} else {
		throwNotFoundError(w)
	}
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
		throwNotFoundError(w)
	}
	renderJSON(w, task)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(task) > 0 {
		renderJSON(w, task)
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
func (ts *configServer) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := ts.store.DeleteGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(task) > 0 {
		renderJSON(w, task)
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
//func (ts *dbServerConfig) appendGroupHandler(w http.ResponseWriter, req *http.Request) {
//	checkRequest(req, w)
//	id := mux.Vars(req)["id"]
//	oldVersion := mux.Vars(req)["version"]
//	newVersion := mux.Vars(req)["new"]
//	oldData := id + separator() + oldVersion
//	newData := id + separator() + newVersion
//	rt, err := decodeGroupConfigs(req.Body)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	if v, ok := ts.dataGroup[oldData]; ok {
//		var newGroup DBGroup
//		newGroup.Configs = make(map[string]*DBConfig)
//		newGroup.Id = newData
//		for _, old := range v.Configs {
//			x := old.dBToGroupConfig()
//			x.Id = createId(x.Id)
//			appen := x.groupConfigToDBConfig()
//			newGroup.Configs[appen.Id] = appen
//		}
//		for _, val := range rt.Configs {
//			dbg := val.groupConfigToDBConfig()
//			if newGroup.Configs[dbg.Id] == nil {
//				newGroup.Configs[dbg.Id] = dbg
//			}
//		}
//		if ts.dataGroup[newGroup.Id] != nil {
//			throwForbiddenError(w)
//		} else {
//			ts.dataGroup[newGroup.Id] = &newGroup
//			renderJSON(w, newGroup.dBGroupToFree())
//		}
//	} else {
//		throwNotFoundError(w)
//	}
//}

// swagger:route GET /group/{id}/{version}/{labels}/ Label getConfigsByLabel
// Get configs by label
//
// responses:
//
//	404: ErrorResponse
//	418: Teapot
//	200: []GroupConfig
//func (ts *dbServerConfig) getConfigsByLabel(w http.ResponseWriter, req *http.Request) {
//	id := mux.Vars(req)["id"] + separator() + mux.Vars(req)["version"]
//	labels := mux.Vars(req)["labels"]
//	allTasks := []*GroupConfig{}
//	if v, ok := ts.dataGroup[id]; ok {
//		for _, val := range v.Configs {
//			if val.compareLabels(labels) {
//				allTasks = append(allTasks, val.dBToGroupConfig())
//			}
//		}
//	} else {
//		throwNotFoundError(w)
//	}
//	if len(allTasks) == 0 {
//		throwTeapot(w)
//	} else {
//		renderJSON(w, allTasks)
//	}
//}

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

	task, err := ts.store.DeleteConfigsByLabels(id, version, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(task) > 0 {
		renderJSON(w, task)
	} else {
		throwNotFoundError(w)
	}

}
