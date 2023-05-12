package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

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

func (ts *dbServerConfig) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*FreeGroup{}
	for _, v := range ts.dataGroup {
		allTasks = append(allTasks, v.dBGroupToFree())
	}
	renderJSON(w, allTasks)
}

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

func (ts *dbServerConfig) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"] + separator() + mux.Vars(req)["version"]
	task, ok := ts.dataGroup[id]
	if !ok {
		throwNotFoundError(w)
		return
	}
	renderJSON(w, task.dBGroupToFree())
}

func (ts *dbServerConfig) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"] + separator() + mux.Vars(req)["version"]
	if v, ok := ts.dataGroup[id]; ok {
		delete(ts.dataGroup, id)
		renderJSON(w, v.dBGroupToFree())
	} else {
		throwNotFoundError(w)
	}
}

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
			renderJSON(w, newGroup)
		}
	} else {
		throwNotFoundError(w)
	}

}
