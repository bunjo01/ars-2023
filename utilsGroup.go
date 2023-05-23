package main

import (
	"encoding/json"
	"io"
	"strings"
)

func decodeFreeGroup(r io.Reader) (*FreeGroup, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt FreeGroup
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (gr *FreeGroup) freeToDBGroup() *DBGroup {
	var rt DBGroup
	rt.Configs = make(map[string]*DBConfig)
	rt.Id = gr.Id + separator() + gr.Version
	for _, v := range gr.Configs {
		conf := v.groupConfigToDBConfig()
		rt.Configs[conf.Id] = conf
	}
	return &rt
}

func (gr *DBGroup) dBGroupToFree() *FreeGroup {
	var rt FreeGroup
	rt.Configs = make(map[string]*GroupConfig)
	comb := strings.Split(gr.Id, separator())
	rt.Id = comb[0]
	rt.Version = comb[1]
	for _, v := range gr.Configs {
		conf := v.dBToGroupConfig()
		rt.Configs[conf.Id] = conf
	}
	return &rt
}

func decodeGroupConfigs(r io.Reader) (*GroupConfigList, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt GroupConfigList
	if err := dec.Decode(&rt.Configs); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (con *GroupConfig) groupConfigToDBConfig() *DBConfig {
	var rt DBConfig
	lab := ""
	for k, v := range con.Labels {
		lab = lab + k + ":" + *v + labelSeparator()
	}
	lab = lab[:len(lab)-1]
	rt.Id = lab + separator() + con.Id
	rt.Entries = con.Entries
	return &rt
}

func (con *DBConfig) dBToGroupConfig() *GroupConfig {
	var rt GroupConfig
	comb := strings.Split(con.Id, separator())
	labels := mapConfigLabels(comb[0])
	rt.Labels = labels
	rt.Id = comb[1]
	rt.Entries = con.Entries
	return &rt
}

func (con *DBConfig) compareLabels(labelString string) bool {
	labels := mapConfigLabels(labelString)
	confLabels := con.dBToGroupConfig().Labels
	if len(labels) == len(confLabels) {
		for k, _ := range labels {
			if confLabels[k] != nil && *confLabels[k] == *labels[k] {
				continue
			} else {
				return false
			}
		}
	} else {
		return false
	}
	return true
}