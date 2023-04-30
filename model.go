package main

type Config struct {
	Id      string             `json:"id"`
	Entries map[string]*string `json:"entries"`
}

type Group struct {
	Id      string             `json:"id"`
	Configs map[string]*string `json:"configs"`
}

type AdoConfig struct {
	Appends map[string]*string `json:"configIds"`
}
