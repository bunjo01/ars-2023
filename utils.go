package main

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"mime"
	"net/http"
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

func (con *FreeConfig) freeToDBConfig() *Config {
	var rt Config
	rt.Id = con.Id + separator() + con.Vers
	rt.Entries = con.Entries
	return &rt
}
func (con *Config) dBConfigToFree() *FreeConfig {
	var rt FreeConfig
	comb := strings.Split(con.Id, separator())
	rt.Id = comb[0]
	rt.Vers = comb[1]
	rt.Entries = con.Entries
	return &rt
}

func decodeGroupBody(r io.Reader) (*Group, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt Group
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

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId(id string) string {
	ret := uuid.New().String()
	staticId := id
	if staticId == "1234" {
		ret = staticId
	}
	return ret
}
func separator() string {
	return "|"
}

func throwNotFoundError(w http.ResponseWriter) {
	err := errors.New("id not found")
	http.Error(w, err.Error(), http.StatusNotFound)
}

func checkRequest(req *http.Request, w http.ResponseWriter) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expected application/json Content-type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
}
