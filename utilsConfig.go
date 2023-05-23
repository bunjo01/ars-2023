package main

import (
	"encoding/json"
	"io"
	"strings"
)

func decodeFreeConfig(r io.Reader) (*FreeConfig, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt FreeConfig
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (con *FreeConfig) freeToDBConfig() *DBConfig {
	var rt DBConfig
	rt.Id = con.Id + separator() + con.Version
	rt.Entries = con.Entries
	return &rt
}
func (con *DBConfig) dBConfigToFree() *FreeConfig {
	var rt FreeConfig
	comb := strings.Split(con.Id, separator())
	rt.Id = comb[0]
	rt.Version = comb[1]
	rt.Entries = con.Entries
	return &rt
}
