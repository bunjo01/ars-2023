package main

// swagger: response ResponseConfig
type ResponseFreeConfig struct {
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

// swagger:response ErrorResponse
type ErrorResponse struct {
	// Error status code
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:response NoContentResponse
type NoContentResponse struct{}
