package main

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
