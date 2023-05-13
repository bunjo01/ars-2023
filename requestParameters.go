package main

// swagger:parameters createConfig
type RequestConfigBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/FreeConfig"
	//  required: true
	Body FreeConfig `json:"body"`
}

// swagger:parameters getAllConfigVersions getConfig deleteConfigVersions deleteConfig
type ConfigId struct {
	// name: ID
	// description: Config ID
	// required: true
	// in: path
	Id string `json:"id"`
}

// swagger:parameters getConfig deleteConfig
type ConfigVersion struct {
	// name: Version
	// description: Config description
	// required: true
	//in: path
	Version string `json:"version"`
}
