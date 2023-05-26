package utilities

import (
	cs "ars-2023/configdatabase"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

// JSON decoders

func DecodeFreeConfig(r io.Reader) (*cs.FreeConfig, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt cs.FreeConfig
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}
func DecodeFreeGroup(r io.Reader) (*cs.FreeGroup, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt cs.FreeGroup
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

// JSON render

func RenderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Error throws

func ThrowBadRequest(w http.ResponseWriter) {
	err := errors.New("bad request")
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func ThrowNotFoundError(w http.ResponseWriter) {
	err := errors.New("not found")
	http.Error(w, err.Error(), http.StatusNotFound)
}

func ThrowTeapot(w http.ResponseWriter) {
	err := errors.New("The server refuses the attempt to brew coffee with a teapot.")
	http.Error(w, err.Error(), http.StatusTeapot)
}

// Http request validator

func CheckRequest(req *http.Request, w http.ResponseWriter) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		ThrowBadRequest(w)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expected application/json Content-type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
}
