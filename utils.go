package main

import (
	cs "ars-2023/configdatabase"
	er "ars-2023/errors"
	"ars-2023/tracer"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"io"
	"mime"
	"net/http"
)

const (
	name = "config_service"
)

// Server access structure

type configServer struct {
	store *cs.ConfigStore
	Keys  map[string]string

	tracer opentracing.Tracer
	closer io.Closer
}

func NewConfigServer() (*configServer, error) {
	store, err := cs.New()
	if err != nil {
		return nil, err
	}

	tracer, closer := tracer.Init(name)
	opentracing.SetGlobalTracer(tracer)
	return &configServer{
		store:  store,
		Keys:   make(map[string]string),
		tracer: tracer,
		closer: closer,
	}, nil
}

func (s *configServer) GetTracer() opentracing.Tracer {
	return s.tracer
}

func (s *configServer) GetCloser() io.Closer {
	return s.closer
}

func (s *configServer) CloseTracer() error {
	return s.closer.Close()
}

// JSON decoders

func DecodeFreeConfig(r io.Reader) (*cs.FreeConfig, *er.ErrorResponse) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt cs.FreeConfig
	if err := dec.Decode(&rt); err != nil {
		return nil, er.NewError(400)
	}
	return &rt, nil
}
func DecodeFreeGroup(r io.Reader) (*cs.FreeGroup, *er.ErrorResponse) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt cs.FreeGroup
	if err := dec.Decode(&rt); err != nil {
		return nil, er.NewError(400)
	}
	return &rt, nil
}

// JSON render

func RenderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		throwError(w, er.NewError(400))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Errors

func throwError(w http.ResponseWriter, err *er.ErrorResponse) {
	http.Error(w, err.Message, err.Status)
}

// Http request validator

func (cv *configServer) CheckRequest(req *http.Request) *er.ErrorResponse {
	ok, err := IsFirst(req, cv)
	if err != nil {
		return err
	}
	if ok && err != nil {
		return err
	}
	contentType := req.Header.Get("Content-Type")
	mediaType, _, erro := mime.ParseMediaType(contentType)
	if erro != nil {
		return er.NewError(400)
	}
	if mediaType != "application/json" {
		return er.NewError(415)
	}
	return nil
}

// Idempotency check

func IsFirst(req *http.Request, cv *configServer) (bool, *er.ErrorResponse) {
	key := req.Header.Get("Idempotency-key")
	if key == "" {
		return false, er.NewError(418)
	}
	if cv.Keys[key] == "used" {
		return true, er.NewError(409)
	}
	cv.Keys[key] = "used"
	return true, nil
}
