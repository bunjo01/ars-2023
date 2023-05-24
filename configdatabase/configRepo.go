package configdatabase

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/consul/api"
)

func (ps *ConfigStore) GetConfig(id string, version string) ([]*FreeConfig, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(constructConfigKey(id, version), nil)
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
	return configs, nil
}

func (ps *ConfigStore) GetConfigVersions(id string) ([]*FreeConfig, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List("config/"+id+"/", nil)
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

	return configs, nil
}

func (ps *ConfigStore) DeleteConfig(id string, version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.Delete(constructConfigKey(id, version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) DeleteConfigVersions(id string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree("config/"+id+"/", nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) Config(config *FreeConfig) (*FreeConfig, error) {
	kv := ps.cli.KV()

	sid, rid := generateConfigKey(config.Version)

	staticId := config.Id
	if staticId == "1234" {
		static := constructConfigKey(staticId, config.Version)
		rid = staticId
		sid = static
	}

	val, _, err := kv.Get(sid, nil)

	if val != nil {
		return nil, errors.New("DESI POŠO GARIŠON")
	}

	config.Id = rid

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return config, nil
}
