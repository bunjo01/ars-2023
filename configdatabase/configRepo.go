package configdatabase

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
)

func New() (*ConfigStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigStore{
		cli: client,
	}, nil
}

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

	data, _, err := kv.List(id, nil)
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

func (ps *ConfigStore) GetAll() ([]*FreeConfig, error) {
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
	_, err := kv.DeleteTree(constructConfigKey(id, version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) DeleteConfigVersions(id string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(id, nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) Config(config *FreeConfig) (*FreeConfig, error) {
	kv := ps.cli.KV()

	sid, rid := generateConfigKey(config.Version)
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
