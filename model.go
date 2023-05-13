package main

type DBConfig struct {
	Id      string             `json:"id"`
	Entries map[string]*string `json:"entries"`
}

type DBGroup struct {
	Id      string               `json:"id"`
	Configs map[string]*DBConfig `json:"configs"`
}

type FreeGroup struct {
	Id      string                  `json:"id"`
	Version string                  `json:"version"`
	Configs map[string]*GroupConfig `json:"configs"`
}

type GroupConfigList struct {
	Configs map[string]*GroupConfig `json:"configs,string"`
}

// FreeConfig swagger: model FreeConfig
type FreeConfig struct {
	// Id of the config
	// in: string
	Id string `json:"id"`

	// Version of the config
	// in: string
	Version string `json:"version"`

	// Entries map of config options
	// in: map[string]*string
	Entries map[string]*string `json:"entries"`
}

type GroupConfig struct {
	Id      string             `json:"id"`
	Labels  map[string]*string `json:"labels"`
	Entries map[string]*string `json:"entries"`
}
