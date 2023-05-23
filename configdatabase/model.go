package configdatabase

import "github.com/hashicorp/consul/api"

type FreeConfig struct {
	Id      string             `json:"id"`
	Version string             `json:"version"`
	Entries map[string]*string `json:"entries"`
}

type ConfigStore struct {
	cli *api.Client
}
