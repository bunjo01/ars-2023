package configdatabase

import (
	"ars-2023/tracer"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

func (ps *ConfigStore) GetGroup(id, version string, ctx context.Context) (*FreeGroup, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbGetGroup")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling get group in db\n")),
	)

	kv := ps.cli.KV()
	key := dbKeyGen("group", id, version)
	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, tracer.NewError(404, span)
	}
	group := &FreeGroup{id, version, make(map[string]*GroupConfig)}
	for _, pair := range data {
		config := &GroupConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, tracer.NewError(400, span)
		}
		group.Configs[config.Id] = config
	}
	if len(group.Configs) == 0 {
		return nil, tracer.NewError(404, span)
	}
	return group, nil
}

func (ps *ConfigStore) GetGroupVersions(id string, ctx context.Context) ([]*FreeGroup, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbGetGroup")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling get group in db\n")),
	)

	kv := ps.cli.KV()
	key := dbKeyGen("group", id)
	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, tracer.NewError(404, span)
	}
	groupMap := make(map[string]*FreeGroup)
	groupList := []*FreeGroup{}
	for _, ky := range data {
		config := &GroupConfig{}
		err := json.Unmarshal(ky.Value, config)
		if err != nil {
			return nil, tracer.NewError(400, span)
		}
		groupVersion := getKeyIndexInfo(group_version, ky.Key)
		if groupMap[groupVersion] == nil {
			newGroup := &FreeGroup{id, groupVersion, make(map[string]*GroupConfig)}
			newGroup.Configs[config.Id] = config
			groupMap[groupVersion] = newGroup
		} else {
			groupMap[groupVersion].Configs[config.Id] = config
		}
	}
	for _, g := range groupMap {
		groupList = append(groupList, g)
	}
	if len(groupList) == 0 {
		return nil, tracer.NewError(404, span)
	}
	return groupList, nil
}

func (ps *ConfigStore) GetAllGroups(ctx context.Context) ([]*FreeGroup, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbGetAllGroups")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling get all groups in db\n")),
	)

	kv := ps.cli.KV()
	key := "group/"
	data, _, err := kv.List(key, nil)

	if err != nil {
		return nil, tracer.NewError(404, span)
	}
	groupMap := make(map[string]*FreeGroup)
	groupList := []*FreeGroup{}
	for _, ky := range data {
		config := &GroupConfig{}
		err := json.Unmarshal(ky.Value, config)
		if err != nil {
			return nil, tracer.NewError(400, span)
		}

		groupVersion := getKeyIndexInfo(group_version, ky.Key)
		groupId := getKeyIndexInfo(group_id, ky.Key)
		if groupMap[groupVersion+groupId] == nil {
			newGroup := &FreeGroup{groupId, groupVersion, make(map[string]*GroupConfig)}
			newGroup.Configs[config.Id] = config
			groupMap[groupVersion+groupId] = newGroup
		} else {
			groupMap[groupVersion+groupId].Configs[config.Id] = config
		}
	}
	for _, g := range groupMap {
		groupList = append(groupList, g)
	}
	if len(groupList) == 0 {
		return nil, tracer.NewError(404, span)
	}
	return groupList, nil
}

func (ps *ConfigStore) DeleteGroup(id, version string, ctx context.Context) (map[string]string, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbDeleteGroup")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling delete group in db\n")),
	)

	kv := ps.cli.KV()
	_, err := kv.DeleteTree(dbKeyGen("group", id, version, ""), nil)
	if err != nil {
		return nil, tracer.NewError(404, span)
	}
	return map[string]string{"Deleted": id + "/" + version}, nil
}

func (ps *ConfigStore) DeleteGroupVersions(id string, ctx context.Context) (map[string]string, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbDeleteGroupVersions")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling delete group versions in db\n")),
	)

	kv := ps.cli.KV()
	_, err := kv.DeleteTree(dbKeyGen("group", id, ""), nil)
	if err != nil {
		return nil, tracer.NewError(404, span)
	}
	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) Group(group *FreeGroup, creating bool, ctx context.Context) (*FreeGroup, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbCreateGroup")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling create group in db\n")),
	)

	ctx = tracer.ContextWithSpan(ctx, span)

	kv := ps.cli.KV()

	if creating {
		group.Id = CreateId(group.Id)
	}
	if con, err := checkConflict("group", group.Id, group.Version, kv, ctx); con {
		return nil, err
	}
	for _, v := range group.Configs {
		cId := uuid.New().String()
		v.Id = cId
		labels := generateLabelString(v.Labels)
		data, err := json.Marshal(v)
		if err != nil {
			return nil, tracer.NewError(400, span)
		}
		p := &api.KVPair{Key: dbKeyGen("group", group.Id, group.Version, labels, cId), Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			return nil, tracer.NewError(400, span)
		}
	}
	return group, nil
}

func (ps *ConfigStore) DeleteConfigsByLabels(id, version, labels, newVersion string, ctx context.Context) (*FreeGroup, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbLabelGroupDelete")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling delete group configs by label in db\n")),
	)

	ctx = tracer.ContextWithSpan(ctx, span)

	kv := ps.cli.KV()
	labels = sortLabels(labels)
	if con, err := checkConflict("group", id, newVersion, kv, ctx); con {
		return nil, err
	}
	allKeys, _, err := kv.Keys(dbKeyGen("group", id, version, ""), "", nil)
	if err != nil {
		return nil, tracer.NewError(400, span)
	}
	matchKeys, _, err := kv.Keys(dbKeyGen("group", id, version, labels, ""), "", nil)
	if err != nil {
		return nil, tracer.NewError(400, span)
	}
	keyList := []string{}
	for _, key := range allKeys {
		flag := true
		for _, ky := range matchKeys {
			if ky == key {
				flag = false
			}
		}
		if flag {
			keyList = append(keyList, key)
		}
	}
	if len(keyList) == 0 {
		return nil, tracer.NewError(404, span)
	}
	ret := &FreeGroup{id, newVersion, make(map[string]*GroupConfig)}
	for _, key := range keyList {
		config := &GroupConfig{}
		conf, _, err := kv.Get(key, nil)
		if err != nil {
			return nil, tracer.NewError(400, span)
		}
		err = json.Unmarshal(conf.Value, config)
		if err != nil {
			return nil, tracer.NewError(400, span)
		}
		ret.Configs[config.Id] = config

	}
	group, erro := ps.Group(ret, true, ctx)
	if err != nil {
		return nil, erro
	}
	return group, nil
}

func (ps *ConfigStore) GetConfigsByLabels(id, version, labels string, ctx context.Context) (*FreeGroup, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbLabelGroupGet")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling get group configs by label in db\n")),
	)

	kv := ps.cli.KV()
	key := dbKeyGen("group", id, version, labels, "")
	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, tracer.NewError(404, span)
	}
	group := &FreeGroup{id, version, make(map[string]*GroupConfig)}
	for _, pair := range data {
		config := &GroupConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, tracer.NewError(400, span)
		}
		group.Configs[config.Id] = config
	}
	if len(group.Configs) == 0 {
		return nil, tracer.NewError(404, span)
	}
	return group, nil
}
