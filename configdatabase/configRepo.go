package configdatabase

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/consul/api"
)

func (ps *ConfigStore) GetConfig(id string, version string) (*FreeConfig, error) {
	kv := ps.cli.KV()
	data, _, err := kv.Get(dbKeyGen("config", id, version), nil)
	if err != nil || data == nil {
		return nil, errors.New("not found")
	}
	config := &FreeConfig{}
	err = json.Unmarshal(data.Value, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (ps *ConfigStore) GetConfigVersions(id string) ([]*FreeConfig, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(dbKeyGen("config", id), nil)
	if err != nil {
		return nil, err
	}
	configs := []*FreeConfig{}
	for _, pair := range data {
		config := &FreeConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	if len(configs) == 0 {
		return nil, errors.New("not found")
	}
	return configs, nil
}

func (ps *ConfigStore) GetAllConfigs() ([]*FreeConfig, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List("config/", nil)
	if err != nil {
		return nil, err
	}
	configs := []*FreeConfig{}
	for _, pair := range data {
		config := &FreeConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	if len(configs) == 0 {
		return nil, errors.New("not found")
	}
	return configs, nil
}

func (ps *ConfigStore) DeleteConfig(id string, version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.Delete(dbKeyGen("config", id, version), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Deleted": id + "/" + version}, nil
}

func (ps *ConfigStore) DeleteConfigVersions(id string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(dbKeyGen("config", id), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) Config(config *FreeConfig) (*FreeConfig, error) {
	kv := ps.cli.KV()
	config.Id = CreateId(config.Id)
	key := dbKeyGen("config", config.Id, config.Version)
	if con, err := checkConflict("config", config.Id, config.Version, kv); con {
		return nil, err
	}
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return config, nil
}
