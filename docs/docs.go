// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Onkar Chendage",
            "email": "onkar.chendage@gmail.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/check": {
            "post": {
                "description": "Checks if a user has the specified action on a resource",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "privilege"
                ],
                "summary": "Check user privilege",
                "parameters": [
                    {
                        "description": "Privilege check request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CheckPrivilegeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.CheckPrivilegeResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.CheckPrivilegeRequest": {
            "type": "object",
            "properties": {
                "action": {
                    "description": "The action the user wants to perform",
                    "type": "string"
                },
                "resource": {
                    "description": "The resource the action is performed on",
                    "type": "string"
                },
                "user": {
                    "description": "The user whose privileges are being checked",
                    "type": "string"
                }
            }
        },
        "dto.CheckPrivilegeResponse": {
            "type": "object",
            "properties": {
                "allowed": {
                    "description": "Whether the user is allowed to perform the action",
                    "type": "boolean"
                }
            }
        },
        "dto.ErrorResponse": {
            "description": "Standard error response format used by the API",
            "type": "object",
            "properties": {
                "code": {
                    "description": "@Description\tHTTP status code for the error\n\t@Example\t\t400",
                    "type": "integer"
                },
                "error": {
                    "description": "@Description\tError type, like \"bad_request\", \"internal_server_error\", etc.\n\t@Example\t\t\"bad_request\"",
                    "type": "string"
                },
                "message": {
                    "description": "@Description\tDetailed error message explaining the issue\n\t@Example\t\t\"Failed to parse request\"",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{"http", "https"},
	Title:            "Aura API",
	Description:      "## About\n\n`aura` is an authorizer created by [Onkar Chendage](https://github.com/highonsemicolon)\n\n- Source Code: <https://github.com/highonsemicolon/aura> \n- API Docs: <http://localhost:8080/docs/index.html> ",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
