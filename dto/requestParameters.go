package dto

import (
	cbd "ars-2023/configdatabase"
)

// swagger:parameters createConfig
type RequestConfigBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/FreeConfig"
	//  required: true
	Body cbd.FreeConfig `json:"body"`
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
	Body cbd.FreeGroup `json:"body"`
}

// swagger:parameters getAllConfigVersions getConfig deleteConfigVersions deleteConfig delConfigsByLabel
// swagger:parameters getAllGroupVersions getGroup deleteGroup deleteGroupVersions appendGroup getConfigsByLabel
type EntityId struct {
	// name: ID
	// description: Entity ID
	// required: true
	// in: path
	Id string `json:"id"`
}

// swagger:parameters getConfig deleteConfig getGroup deleteGroup appendGroup getConfigsByLabel delConfigsByLabel
type EntityVersion struct {
	// name: Version
	// description: Entity version
	// required: true
	// in: path
	Version string `json:"version"`
}

// swagger:parameters appendGroup delConfigsByLabel
type EntityNewVersion struct {
	// name: Version
	// description: Entity version
	// required: true
	// in: path
	Version string `json:"new"`
}

// swagger:parameters getConfigsByLabel delConfigsByLabel
type EntityLabels struct {
	// name: Labels
	// description: Entity labels
	// required: true
	// in: path
	Labels string `json:"labels"`
}
