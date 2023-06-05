package main

import (
	cs "ars-2023/configdatabase"
	"ars-2023/tracer"
	"context"
	"encoding/json"
	"fmt"
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

	trace, closer := tracer.Init(name)
	opentracing.SetGlobalTracer(trace)
	return &configServer{
		store:  store,
		Keys:   make(map[string]string),
		tracer: trace,
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

func DecodeFreeConfig(r io.Reader, ctx context.Context) (*cs.FreeConfig, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "decodeConfig")
	defer span.Finish()
	span.LogFields(
		tracer.LogString("requestUtility", fmt.Sprintf("decoding configuration")),
	)

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt cs.FreeConfig
	if err := dec.Decode(&rt); err != nil {
		return nil, tracer.NewError(400, span)
	}
	return &rt, nil
}
func DecodeFreeGroup(r io.Reader, ctx context.Context) (*cs.FreeGroup, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "decodeGroup")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("requestUtility", fmt.Sprintf("decoding group")),
	)

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt cs.FreeGroup
	if err := dec.Decode(&rt); err != nil {
		return nil, tracer.NewError(400, span)
	}
	return &rt, nil
}

// JSON render

func RenderJSON(w http.ResponseWriter, v interface{}, ctx context.Context) {
	span := tracer.StartSpanFromContext(ctx, "renderJSON")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("requestUtility", fmt.Sprintf("rendering JSON")),
	)

	js, err := json.Marshal(v)
	if err != nil {
		throwError(w, tracer.NewError(400, span))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Errors

func throwError(w http.ResponseWriter, err *tracer.ErrorResponse) {
	http.Error(w, err.Message, err.Status)
}

// Http request validator

func (cv *configServer) CheckRequest(req *http.Request, ctx context.Context) *tracer.ErrorResponse {
	span := tracer.StartSpanFromContext(ctx, "requestCheck")
	defer span.Finish()
	span.LogFields(
		tracer.LogString("requestUtility", fmt.Sprintf("checking request content type")),
	)
	ctx = tracer.ContextWithSpan(ctx, span)
	ok, err := IsFirst(req, cv, ctx)
	if err != nil {
		return err
	}
	if ok && err != nil {
		return err
	}
	contentType := req.Header.Get("Content-Type")
	mediaType, _, erro := mime.ParseMediaType(contentType)
	if erro != nil {
		return tracer.NewError(400, span)
	}
	if mediaType != "application/json" {
		return tracer.NewError(415, span)
	}
	return nil
}

// Idempotency check

func IsFirst(req *http.Request, cv *configServer, ctx context.Context) (bool, *tracer.ErrorResponse) {
	span := tracer.StartSpanFromContext(ctx, "idempotencyCheck")
	defer span.Finish()

	span.LogFields(
		tracer.LogString("requestUtility", fmt.Sprintf("checking request idempotency")),
	)

	key := req.Header.Get("Idempotency-key")
	if key == "" {
		return false, tracer.NewError(418, span)
	}
	if cv.Keys[key] == "used" {
		return true, tracer.NewError(409, span)
	}
	cv.Keys[key] = "used"
	return true, nil
}
