package main

type DBConfig struct {
	Id      string             `json:"id"`
	Entries map[string]*string `json:"entries"`
}

type DBGroup struct {
	Id      string               `json:"id"`
	Configs map[string]*DBConfig `json:"configs"`
}

// FreeGroup swagger: model FreeGroup
type FreeGroup struct {
	// Id of the group
	// in: string
	Id string `json:"id"`

	// Version of the group
	// in: string
	Version string `json:"version"`

	// Entries map of configs
	// in: map[string]*GroupConfig
	Configs map[string]*GroupConfig `json:"configs"`
}

// GroupConfigList swagger: model GroupConfigList
type GroupConfigList struct {
	// Entries map of configs
	// in: map[string]*GroupConfig
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

// GroupConfig swagger: model GroupConfig
type GroupConfig struct {
	// Id of the config
	// in: string
	Id string `json:"id"`

	// Labels of the config
	// in: map[string]*string
	Labels map[string]*string `json:"labels"`

	// Entries of the config
	// in: map[string]*string
	Entries map[string]*string `json:"entries"`
}
