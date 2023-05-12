package main

type DBConfig struct {
	Id      string             `json:"id"`
	Entries map[string]*string `json:"entries"`
}

type DBGroup struct {
	Id      string               `json:"id"`
	Version string               `json:"version"`
	Configs map[string]*DBConfig `json:"configs"`
}

type DTOConfig struct {
	NewConfigs []*string `json:"configIds"`
}

type FreeConfig struct {
	Id      string             `json:"id"`
	Vers    string             `json:"version"`
	Entries map[string]*string `json:"entries"`
}

type GroupConfig struct {
	Id      string             `json:"id"`
	Labels  string             `json:"labels"`
	Entries map[string]*string `json:"entries"`
}
