// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/games": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Get all games",
                "tags": [
                    "games"
                ],
                "summary": "Get all games",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/presenter.Response"
                        }
                    }
                }
            }
        },
        "/games/{gameID}": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Get game given an ID",
                "tags": [
                    "games"
                ],
                "summary": "Get a single game",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Game identifier",
                        "name": "gameID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/presenter.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "presenter.GameBody": {
            "type": "object",
            "properties": {
                "kills": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                },
                "kills_by_means": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                },
                "players": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "total_kills": {
                    "type": "integer"
                }
            }
        },
        "presenter.Response": {
            "type": "object",
            "additionalProperties": {
                "$ref": "#/definitions/presenter.GameBody"
            }
        }
    },
    "securityDefinitions": {
        "JWT": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "0.0.0.0:8000",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Quake API",
	Description:      "This is a Quake API server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
