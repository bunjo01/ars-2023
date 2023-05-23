package configdatabase

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	config       = "config/%s/%s"
	configLabels = "config/%s/%s/%s"
	allConfigs   = "config/all"
	group        = "group/%s/%s/config/%s"
	groupLabels  = "group/%s/%s/%s/config/%s"
	allGroups    = "group/all"
)

func generateConfigKey(version string, labels string) (string, string) {
	id := uuid.New().String()
	if labels != "" {
		return fmt.Sprintf(configLabels, id, version, labels), id

	} else {
		return fmt.Sprintf(config, id, version), id
	}
}

func constructConfigKey(id string, version string, labels string) string {

	if labels != "" {
		return fmt.Sprintf(configLabels, id, version, labels)

	} else {
		return fmt.Sprintf(config, id, version)
	}
}

func generateGroupKey(version string, labels string) (string, string) {
	idGroup := uuid.New().String()
	idConfig := uuid.New().String()
	if labels != "" {
		return fmt.Sprintf(groupLabels, idGroup, version, labels, idConfig), idGroup

	} else {
		return fmt.Sprintf(group, idGroup, version, idConfig), idGroup
	}
}

func constructGroupKey(idGroup string, version string, labels string, idConfig string) string {

	if labels != "" {
		return fmt.Sprintf(groupLabels, idGroup, version, labels, idConfig)

	} else {
		return fmt.Sprintf(group, idGroup, version, idConfig)
	}
}
