// Package docs is a stub generated to allow compilation before running swag init.
// Run `swag init` from the backend/ directory to regenerate the full OpenAPI spec.
package docs

import "github.com/swaggo/swag"

// SwaggerInfo holds the API metadata used by the Swagger UI.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Aelidecx API",
	Description:      "Self-hosted personal knowledge base. Run `swag init` to regenerate.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "Aelidecx API",
    "description": "Self-hosted personal knowledge base. Run swag init to regenerate full spec.",
    "version": "1.0"
  },
  "host": "localhost:3000",
  "basePath": "/api",
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  },
  "paths": {}
}`
