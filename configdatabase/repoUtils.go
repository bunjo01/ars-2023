package configdatabase

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"sort"
	"strings"
)

const (
	// config/config_id/config_version
	config = "config/%s/%s"
	// group/group_id/group_version/labels/config_id
	group = "group/%s/%s/%s/%s"
)

type Position int

const (
	group_id      = 1
	group_version = 2
	group_labels  = 3
	config_id     = 4
)

func generateConfigKey(version string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(config, id, version), id
}

func constructConfigKey(id, version string) string {
	return fmt.Sprintf(config, id, version)
}

func generateGroupKey(version, labels string) (string, string, string) {
	idGroup := uuid.New().String()
	idConfig := uuid.New().String()
	return fmt.Sprintf(group, idGroup, version, labels, idConfig), idGroup, idConfig
}

func constructGroupKey(idGroup, version, labels, idConfig string) string {
	return fmt.Sprintf(group, idGroup, version, labels, idConfig)
}

// Generic key generator for any entity

func generateCustomKey(params []string, info string) string {
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

func CreateId(id string) string {
	ret := uuid.New().String()
	staticId := id
	if staticId == "1234" {
		ret = staticId
	}
	return ret
}

func getKeyIndexInfo(index Position, key string) string {
	separated := strings.Split(key, "/")
	return separated[index]
}

func TestKeys() {
	log.Println(generateConfigKey("version"))
	log.Println(constructConfigKey("id", "version"))
	log.Println(generateGroupKey("version", "labels"))
	log.Println(constructGroupKey("idGroup", "version", "labels", "idConfig"))
}
