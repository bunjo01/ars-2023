package main

// swagger:parameters post createConfig
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
