package configdatabase

import (
	er "ars-2023/errors"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

func (ps *ConfigStore) GetGroup(id, version string) (*FreeGroup, *er.ErrorResponse) {
	kv := ps.cli.KV()
	key := dbKeyGen("group", id, version)
	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, er.NewError(404)
	}
	group := &FreeGroup{id, version, make(map[string]*GroupConfig)}
	for _, pair := range data {
		config := &GroupConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, er.NewError(400)
		}
		group.Configs[config.Id] = config
	}
	if len(group.Configs) == 0 {
		return nil, er.NewError(404)
	}
	return group, nil
}

func (ps *ConfigStore) GetGroupVersions(id string) ([]*FreeGroup, *er.ErrorResponse) {
	kv := ps.cli.KV()
	key := dbKeyGen("group", id)
	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, er.NewError(404)
	}
	groupMap := make(map[string]*FreeGroup)
	groupList := []*FreeGroup{}
	for _, ky := range data {
		config := &GroupConfig{}
		err := json.Unmarshal(ky.Value, config)
		if err != nil {
			return nil, er.NewError(400)
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
		return nil, er.NewError(404)
	}
	return groupList, nil
}

func (ps *ConfigStore) GetAllGroups() ([]*FreeGroup, *er.ErrorResponse) {
	kv := ps.cli.KV()
	key := "group/"
	data, _, err := kv.List(key, nil)

	if err != nil {
		return nil, er.NewError(404)
	}
	groupMap := make(map[string]*FreeGroup)
	groupList := []*FreeGroup{}
	for _, ky := range data {
		config := &GroupConfig{}
		err := json.Unmarshal(ky.Value, config)
		if err != nil {
			return nil, er.NewError(400)
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
		return nil, er.NewError(404)
	}
	return groupList, nil
}

func (ps *ConfigStore) DeleteGroup(id, version string) (map[string]string, *er.ErrorResponse) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(dbKeyGen("group", id, version, ""), nil)
	if err != nil {
		return nil, er.NewError(404)
	}
	return map[string]string{"Deleted": id + "/" + version}, nil
}

func (ps *ConfigStore) DeleteGroupVersions(id string) (map[string]string, *er.ErrorResponse) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(dbKeyGen("group", id, ""), nil)
	if err != nil {
		return nil, er.NewError(404)
	}
	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) Group(group *FreeGroup, creating bool) (*FreeGroup, *er.ErrorResponse) {
	kv := ps.cli.KV()

	if creating {
		group.Id = CreateId(group.Id)
	}
	if con, err := checkConflict("group", group.Id, group.Version, kv); con {
		return nil, err
	}
	for _, v := range group.Configs {
		cId := uuid.New().String()
		v.Id = cId
		labels := generateLabelString(v.Labels)
		data, err := json.Marshal(v)
		if err != nil {
			return nil, er.NewError(400)
		}
		p := &api.KVPair{Key: dbKeyGen("group", group.Id, group.Version, labels, cId), Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			return nil, er.NewError(400)
		}
	}
	return group, nil
}

func (ps *ConfigStore) DeleteConfigsByLabels(id, version, labels, newVersion string) (*FreeGroup, *er.ErrorResponse) {
	kv := ps.cli.KV()
	labels = sortLabels(labels)
	if con, err := checkConflict("group", id, newVersion, kv); con {
		return nil, err
	}
	allKeys, _, err := kv.Keys(dbKeyGen("group", id, version, ""), "", nil)
	if err != nil {
		return nil, er.NewError(400)
	}
	matchKeys, _, err := kv.Keys(dbKeyGen("group", id, version, labels, ""), "", nil)
	if err != nil {
		return nil, er.NewError(400)
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
		return nil, er.NewError(404)
	}
	ret := &FreeGroup{id, newVersion, make(map[string]*GroupConfig)}
	for _, key := range keyList {
		config := &GroupConfig{}
		conf, _, err := kv.Get(key, nil)
		if err != nil {
			return nil, er.NewError(400)
		}
		err = json.Unmarshal(conf.Value, config)
		if err != nil {
			return nil, er.NewError(400)
		}
		ret.Configs[config.Id] = config

	}
	group, erro := ps.Group(ret, true)
	if err != nil {
		return nil, erro
	}
	return group, nil
}

func (ps *ConfigStore) GetConfigsByLabels(id, version, labels string) (*FreeGroup, *er.ErrorResponse) {
	kv := ps.cli.KV()
	key := dbKeyGen("group", id, version, labels, "")
	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, er.NewError(404)
	}
	group := &FreeGroup{id, version, make(map[string]*GroupConfig)}
	for _, pair := range data {
		config := &GroupConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, er.NewError(400)
		}
		group.Configs[config.Id] = config
	}
	if len(group.Configs) == 0 {
		return nil, er.NewError(404)
	}
	return group, nil
}
