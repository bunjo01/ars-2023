package configdatabase

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
)

type FreeConfig struct {
	Id      string             `json:"id"`
	Version string             `json:"version"`
	Entries map[string]*string `json:"entries"`
}

type GroupConfig struct {
	Id      string             `json:"id"`
	Labels  map[string]*string `json:"labels"`
	Entries map[string]*string `json:"entries"`
}

type FreeGroup struct {
	Id      string                  `json:"id"`
	Version string                  `json:"version"`
	Configs map[string]*GroupConfig `json:"configs"`
}

type ConfigStore struct {
	cli *api.Client
}

func New() (*ConfigStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigStore{
		cli: client,
	}, nil
}
