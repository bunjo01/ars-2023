package configdatabase

import (
	er "ars-2023/errors"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"sort"
	"strings"
)

// config/config_id/config_version
// group/group_id/group_version/labels/config_id

// DB Key position enum

type Position int

const (
	group_id      = 1
	group_version = 2
)

// Generic db key generator for any entity

func dbKeyGen(info string, params ...string) string {
	var sb strings.Builder
	sb.WriteString(info)
	for _, v := range params {
		sb.WriteString("/" + v)
	}
	return sb.String()
}

// Label map to string converter

func generateLabelString(labels map[string]*string) string {
	var sb strings.Builder
	labelList := make([]string, 0, len(labels))
	for k := range labels {
		labelList = append(labelList, k)
	}
	sort.Strings(labelList)
	for _, k := range labelList {
		sb.WriteString(k + ":" + *labels[k] + ";")
	}
	l := sb.String()
	return l[:len(l)-1]
}

func sortLabels(labels string) string {
	if strings.Contains(labels, ";") {
		labs := strings.Split(labels, ";")
		sort.Strings(labs)
		labels = strings.Join(labs, ";")
	}
	return labels

}

// Static ID check for testing examples

func CreateId(id string) string {
	ret := uuid.New().String()
	staticId := id
	if staticId == "1234" {
		ret = staticId
	}
	return ret
}

// Key information getter

func getKeyIndexInfo(index Position, key string) string {
	separated := strings.Split(key, "/")
	return separated[index]
}

//

func checkConflict(info, id, version string, kv *api.KV) (bool, *er.ErrorResponse) {
	key := dbKeyGen(info, id, version)
	val, _, err := kv.List(key, nil)
	if err != nil {
		return true, er.NewError(404)
	}
	if (len(val) > 0 && info == "group") || (val != nil && info == "config") {
		return true, er.NewError(409)
	}
	return false, nil
}
