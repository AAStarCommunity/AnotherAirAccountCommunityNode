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
            "name": "AAStar Support",
            "url": "https://aastar.xyz"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/account/v1/bind": {
            "post": {
                "description": "bind a account to community node",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account"
                ],
                "parameters": [
                    {
                        "description": "Account Binding",
                        "name": "bind",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.Bind"
                        }
                    },
                    {
                        "type": "string",
                        "description": "apiKey",
                        "name": "apiKey",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    }
                }
            }
        },
        "/api/account/v1/transfer": {
            "post": {
                "description": "transfer a TX",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account"
                ],
                "parameters": [
                    {
                        "description": "Transfer TX",
                        "name": "tx",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.Transfer"
                        }
                    },
                    {
                        "type": "string",
                        "description": "apiKey",
                        "name": "apiKey",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    }
                }
            }
        },
        "/api/dashboard/v1/node": {
            "get": {
                "description": "get node members",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dashboard"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/healthz": {
            "get": {
                "description": "Get Healthz",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Healthz"
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "request.Bind": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string"
                },
                "publicKey": {
                    "type": "string"
                }
            }
        },
        "request.Transfer": {
            "type": "object"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
