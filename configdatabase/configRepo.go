package configdatabase

import (
	er "ars-2023/errors"
	"encoding/json"
	"github.com/hashicorp/consul/api"
)

func (ps *ConfigStore) GetConfig(id string, version string) (*FreeConfig, *er.ErrorResponse) {
	kv := ps.cli.KV()
	data, _, err := kv.Get(dbKeyGen("config", id, version), nil)
	if err != nil {
		return nil, er.NewError(400)
	}
	if data == nil {
		return nil, er.NewError(404)
	}
	config := &FreeConfig{}
	err = json.Unmarshal(data.Value, config)
	if err != nil {
		return nil, er.NewError(400)
	}
	return config, nil
}

func (ps *ConfigStore) GetConfigVersions(id string) ([]*FreeConfig, *er.ErrorResponse) {
	kv := ps.cli.KV()

	data, _, err := kv.List(dbKeyGen("config", id), nil)
	if err != nil {
		return nil, er.NewError(400)
	}
	configs := []*FreeConfig{}
	for _, pair := range data {
		config := &FreeConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, er.NewError(400)
		}
		configs = append(configs, config)
	}
	if len(configs) == 0 {
		return nil, er.NewError(404)
	}
	return configs, nil
}

func (ps *ConfigStore) GetAllConfigs() ([]*FreeConfig, *er.ErrorResponse) {
	kv := ps.cli.KV()
	data, _, err := kv.List("config/", nil)
	if err != nil {
		return nil, er.NewError(400)
	}
	configs := []*FreeConfig{}
	for _, pair := range data {
		config := &FreeConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, er.NewError(400)
		}
		configs = append(configs, config)
	}
	if len(configs) == 0 {
		return nil, er.NewError(404)
	}
	return configs, nil
}

func (ps *ConfigStore) DeleteConfig(id string, version string) (map[string]string, *er.ErrorResponse) {
	kv := ps.cli.KV()
	_, err := kv.Delete(dbKeyGen("config", id, version), nil)
	if err != nil {
		return nil, er.NewError(404)
	}
	return map[string]string{"Deleted": id + "/" + version}, nil
}

func (ps *ConfigStore) DeleteConfigVersions(id string) (map[string]string, *er.ErrorResponse) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(dbKeyGen("config", id), nil)
	if err != nil {
		return nil, er.NewError(404)
	}
	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) Config(config *FreeConfig) (*FreeConfig, *er.ErrorResponse) {
	kv := ps.cli.KV()
	config.Id = CreateId(config.Id)
	key := dbKeyGen("config", config.Id, config.Version)
	if con, err := checkConflict("config", config.Id, config.Version, kv); con {
		return nil, err
	}
	data, err := json.Marshal(config)
	if err != nil {
		return nil, er.NewError(400)
	}
	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, er.NewError(400)
	}

	return config, nil
}
