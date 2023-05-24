package configdatabase

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

//getAllGroupHandler
//getGroupVersionsHandler
//delGroupVersionsHandler
//delGroupHandler
//appendGroupHandler

func (ps *ConfigStore) GetGroup(id, version string) (*FreeGroup, error) {
	kv := ps.cli.KV()
	key := generateCustomKey([]string{id, version}, "group/")
	data, _, err := kv.List(key, nil)

	if err != nil {
		return nil, err
	}
	group := &FreeGroup{}
	group.Id = id
	group.Version = version
	group.Configs = make(map[string]*GroupConfig)
	for _, pair := range data {
		config := &GroupConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		group.Configs[config.Id] = config
	}
	return group, nil
}

func (ps *ConfigStore) GetGroupVersions(id string) ([]*FreeGroup, error) {
	kv := ps.cli.KV()
	key := generateCustomKey([]string{id}, "group/")
	data, _, err := kv.List(key, nil)

	if err != nil {
		return nil, err
	}
	groupMap := make(map[string]*FreeGroup)
	groupList := []*FreeGroup{}
	for _, ky := range data {
		config := &GroupConfig{}
		err := json.Unmarshal(ky.Value, config)
		if err != nil {
			return nil, err
		}

		groupVersion := getKeyIndexInfo(group_version, ky.Key)
		if groupMap[groupVersion] == nil {
			newGroup := &FreeGroup{}
			newGroup.Configs = make(map[string]*GroupConfig)
			newGroup.Version = groupVersion
			newGroup.Configs[config.Id] = config
			groupMap[groupVersion] = newGroup
		} else {
			groupMap[groupVersion].Configs[config.Id] = config
		}
	}
	for _, g := range groupMap {
		g.Id = id
		groupList = append(groupList, g)
	}
	return groupList, nil

}

func (ps *ConfigStore) GetAllGroups() ([]*FreeGroup, error) {
	kv := ps.cli.KV()
	key := generateCustomKey([]string{}, "group/")
	data, _, err := kv.List(key, nil)

	if err != nil {
		return nil, err
	}
	groupMap := make(map[string]*FreeGroup)
	groupList := []*FreeGroup{}
	for _, ky := range data {
		config := &GroupConfig{}
		err := json.Unmarshal(ky.Value, config)
		if err != nil {
			return nil, err
		}

		groupVersion := getKeyIndexInfo(group_version, ky.Key)
		groupId := getKeyIndexInfo(group_id, ky.Key)
		if groupMap[groupVersion+groupId] == nil {
			newGroup := &FreeGroup{}
			newGroup.Configs = make(map[string]*GroupConfig)
			newGroup.Version = groupVersion
			newGroup.Id = groupId
			newGroup.Configs[config.Id] = config
			groupMap[groupVersion+groupId] = newGroup
		} else {
			groupMap[groupVersion+groupId].Configs[config.Id] = config
		}
	}
	for _, g := range groupMap {
		groupList = append(groupList, g)
	}
	return groupList, nil
}

func (ps *ConfigStore) DeleteGroup(id, version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(generateCustomKey([]string{id, version}, "group/"), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) DeleteGroupVersions(id string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(generateCustomKey([]string{id}, "group/"), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) Group(group *FreeGroup) (*FreeGroup, error) {
	kv := ps.cli.KV()

	group.Id = CreateId(group.Id)
	for _, v := range group.Configs {
		cId := uuid.New().String()
		v.Id = cId
		labels := generateLabelString(v.Labels)
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		p := &api.KVPair{Key: generateCustomKey([]string{group.Id, group.Version, labels, cId}, "group"), Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			return nil, err
		}
	}
	return group, nil
}

func (ps *ConfigStore) DeleteConfigsByLabels(id, version, labels string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(generateCustomKey([]string{id, version, labels}, "group/"), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) GetConfigsByLabels(id, version, labels string) (*FreeGroup, error) {

}
