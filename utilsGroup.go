package main

import (
	"encoding/json"
	"io"
)

func decodeGroupBody(r io.Reader) (*DBGroup, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt DBGroup
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeAppendBody(r io.Reader) (*DTOConfig, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt DTOConfig
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}
