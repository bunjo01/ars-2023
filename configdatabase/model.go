package configdatabase

import "github.com/hashicorp/consul/api"

type FreeConfig struct {
	Id      string             `json:"id"`
	Version string             `json:"version"`
	Entires map[string]*string `json:"entires"`
}

type ConfigStore struct {
	cli *api.Client
}
