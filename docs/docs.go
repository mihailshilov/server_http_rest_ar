// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "soberkoder@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/servicerequests/": {
            "post": {
                "description": "Создать заявку",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "requests"
                ],
                "summary": "Создать заявку",
                "operationId": "create-request",
                "parameters": [
                    {
                        "description": "list info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.DataRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Response"
                        }
                    }
                }
            }
        },
        "/authentication/": {
            "post": {
                "description": "Auth Login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Auth Login",
                "operationId": "auth-login",
                "responses": {}
            }
        }
    },
    "definitions": {
        "model.DataRequest": {
            "type": "object",
            "required": [
                "date_time_req",
                "id_dep",
                "id_org",
                "id_request",
                "responsible"
            ],
            "properties": {
                "date_time_inf": {
                    "type": "string"
                },
                "date_time_req": {
                    "type": "string"
                },
                "id_dep": {
                    "type": "string"
                },
                "id_org": {
                    "type": "string"
                },
                "id_request": {
                    "type": "string"
                },
                "responsible": {
                    "type": "string"
                }
            }
        },
        "model.Response": {
            "type": "object",
            "properties": {
                "response": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "/",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Repair App API",
	Description:      "API Server for service stations Application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
