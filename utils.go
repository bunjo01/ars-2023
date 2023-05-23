package main

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"mime"
	"net/http"
	"regexp"
	"strings"
)

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
func labelSeparator() string {
	return ";"
}

func throwNotFoundError(w http.ResponseWriter) {
	err := errors.New("ID not found")
	http.Error(w, err.Error(), http.StatusNotFound)
}

func throwForbiddenError(w http.ResponseWriter) {
	err := errors.New("already exists")
	http.Error(w, err.Error(), http.StatusForbidden)
}

func throwTeapot(w http.ResponseWriter) {
	err := errors.New("The server refuses the attempt to brew coffee with a teapot.")
	http.Error(w, err.Error(), http.StatusTeapot)
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

func mapConfigLabels(labelString string) map[string]*string {
	labels := make(map[string]*string)
	reg := regexp.MustCompile("(([^;: ]+:[^;: ]+)(;([^;: ]+:[^;: ]+))*)")
	if reg.MatchString(labelString) {
		for _, v := range strings.Split(labelString, labelSeparator()) {
			en := strings.Split(v, ":")
			labels[en[0]] = &en[1]
		}
		return labels
	} else {
		badMap := "bad_label"
		labels["403"] = &badMap
		return labels
	}
}

func (ts *dbServerConfig) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}
