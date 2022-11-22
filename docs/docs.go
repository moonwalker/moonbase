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
        "contact": {},
        "license": {
            "name": "MIT",
            "url": "https://github.com/moonwalker/moonbase/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/list": {
            "get": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "repositories"
                ],
                "summary": "List repositories",
                "parameters": [
                    {
                        "type": "string",
                        "description": "page of results to retrieve (default: ` + "`" + `1` + "`" + `)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "number of results to include per page (default: ` + "`" + `30` + "`" + `)",
                        "name": "per_page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "how to sort the repository list, can be one of ` + "`" + `created` + "`" + `, ` + "`" + `updated` + "`" + `, ` + "`" + `pushed` + "`" + `, ` + "`" + `full_name` + "`" + ` (default: ` + "`" + `full_name` + "`" + `)",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "direction in which to sort repositories, can be one of ` + "`" + `asc` + "`" + ` or ` + "`" + `desc` + "`" + ` (default when using ` + "`" + `full_name` + "`" + `: ` + "`" + `asc` + "`" + `; otherwise: ` + "`" + `desc` + "`" + `)",
                        "name": "direction",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.repositoryList"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.errorData"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.errorData": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "api.repositoryItem": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                }
            }
        },
        "api.repositoryList": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.repositoryItem"
                    }
                },
                "lastPage": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "bearerToken": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Moonbase API",
	Description:      "### Git-based headless CMS API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
