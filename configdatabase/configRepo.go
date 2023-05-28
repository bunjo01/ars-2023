package configdatabase

import (
	"ars-2023/tracer"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
)

func (ps *ConfigStore) GetConfig(id string, version string, ctx context.Context) (*FreeConfig, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbGetConfig")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling get config in db\n")),
	)

	kv := ps.cli.KV()
	data, _, err := kv.Get(dbKeyGen("config", id, version), nil)
	if err != nil {
		return nil, tracer.NewError(400, span)
	}
	if data == nil {
		return nil, tracer.NewError(404, span)
	}
	config := &FreeConfig{}
	err = json.Unmarshal(data.Value, config)
	if err != nil {
		return nil, tracer.NewError(400, span)
	}
	return config, nil
}

func (ps *ConfigStore) GetConfigVersions(id string, ctx context.Context) ([]*FreeConfig, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbGetConfigVersions")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling get config versions in db\n")),
	)

	kv := ps.cli.KV()

	data, _, err := kv.List(dbKeyGen("config", id), nil)
	if err != nil {
		return nil, tracer.NewError(400, span)
	}
	configs := []*FreeConfig{}
	for _, pair := range data {
		config := &FreeConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, tracer.NewError(400, span)
		}
		configs = append(configs, config)
	}
	if len(configs) == 0 {
		return nil, tracer.NewError(404, span)
	}
	return configs, nil
}

func (ps *ConfigStore) GetAllConfigs(ctx context.Context) ([]*FreeConfig, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbGetAllConfigs")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling get all config in db\n")),
	)

	kv := ps.cli.KV()
	data, _, err := kv.List("config/", nil)
	if err != nil {
		return nil, tracer.NewError(400, span)
	}
	configs := []*FreeConfig{}
	for _, pair := range data {
		config := &FreeConfig{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, tracer.NewError(400, span)
		}
		configs = append(configs, config)
	}
	if len(configs) == 0 {
		return nil, tracer.NewError(404, span)
	}
	return configs, nil
}

func (ps *ConfigStore) DeleteConfig(id string, version string, ctx context.Context) (map[string]string, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbDeleteConfig")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling delete config in db\n")),
	)

	kv := ps.cli.KV()
	_, err := kv.Delete(dbKeyGen("config", id, version), nil)
	if err != nil {
		return nil, tracer.NewError(404, span)
	}
	return map[string]string{"Deleted": id + "/" + version}, nil
}

func (ps *ConfigStore) DeleteConfigVersions(id string, ctx context.Context) (map[string]string, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbDeleteConfigVersions")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling delete config versions in db\n")),
	)

	kv := ps.cli.KV()
	_, err := kv.DeleteTree(dbKeyGen("config", id), nil)
	if err != nil {
		return nil, tracer.NewError(404, span)
	}
	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) Config(config *FreeConfig, ctx context.Context) (*FreeConfig, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "dbCreateConfig")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("db", fmt.Sprintf("handling create config in db\n")),
	)

	ctx = tracer.ContextWithSpan(context.Background(), span)

	kv := ps.cli.KV()
	config.Id = CreateId(config.Id)
	key := dbKeyGen("config", config.Id, config.Version)
	if con, err := checkConflict("config", config.Id, config.Version, kv, ctx); con {
		return nil, err
	}
	data, err := json.Marshal(config)
	if err != nil {
		return nil, tracer.NewError(400, span)
	}
	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, tracer.NewError(400, span)
	}

	return config, nil
}
