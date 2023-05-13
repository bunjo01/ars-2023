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

// swagger:parameters createGroup
type RequestGroupBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/FreeGroup"
	//  required: true
	Body FreeGroup `json:"body"`
}

// swagger:parameters appendGroup
type RequestAppendBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/GroupConfigList"
	//  required: true
	Configs map[string]*GroupConfig `json:"configs,string"`
}

// swagger:parameters getAllConfigVersions getConfig deleteConfigVersions deleteConfig getAllGroupVersions getGroup deleteGroup deleteGroupVersions appendGroup
type EntityId struct {
	// name: ID
	// description: Entity ID
	// required: true
	// in: path
	Id string `json:"id"`
}

// swagger:parameters getConfig deleteConfig getGroup deleteGroup appendGroup
type EntityVersion struct {
	// name: Version
	// description: Entity version
	// required: true
	// in: path
	Version string `json:"version"`
}

// swagger:parameters appendGroup
type EntityNewVersion struct {
	// name: Version
	// description: Entity version
	// required: true
	// in: path
	Version string `json:"new"`
}
