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
        "/cms/{owner}/{repo}/{ref}/collections": {
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
                    "cms"
                ],
                "summary": "Get collections",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "git ref (branch, tag, sha)",
                        "name": "ref",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.treeItem"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.errorData"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "cms"
                ],
                "summary": "Create or Update collection",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "git ref (branch, tag, sha)",
                        "name": "ref",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "collection payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.collectionPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.errorData"
                        }
                    }
                }
            }
        },
        "/cms/{owner}/{repo}/{ref}/collections/{collection}": {
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
                    "cms"
                ],
                "summary": "Get entries",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "git ref (branch, tag, sha)",
                        "name": "ref",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "collection",
                        "name": "collection",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.treeItem"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.errorData"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "cms"
                ],
                "summary": "Create or Update entry",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "git ref (branch, tag, sha)",
                        "name": "ref",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "collection",
                        "name": "collection",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "entry payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.entryPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.errorData"
                        }
                    }
                }
            }
        },
        "/cms/{owner}/{repo}/{ref}/{collection}": {
            "delete": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "cms"
                ],
                "summary": "Delete document",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "git ref (branch, tag, sha)",
                        "name": "ref",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "collection",
                        "name": "collection",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "delete payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.deletePayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.errorData"
                        }
                    }
                }
            }
        },
        "/cms/{owner}/{repo}/{ref}/{collection}/{entry}": {
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
                    "cms"
                ],
                "summary": "Get entry",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "git ref (branch, tag, sha)",
                        "name": "ref",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "collection",
                        "name": "collection",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "entry",
                        "name": "entry",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.entryPayload"
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
        },
        "/repos": {
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
                    "repos"
                ],
                "summary": "Get repositories",
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
        },
        "/repos/{owner}/{repo}/blob/{ref}/{path}": {
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
                    "repos"
                ],
                "summary": "Get blob",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "git ref (branch, tag, sha)",
                        "name": "ref",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "contents path",
                        "name": "path",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.blobEntry"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.errorData"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "repos"
                ],
                "summary": "Post blob",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "git ref (branch, tag, sha)",
                        "name": "ref",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "contents path",
                        "name": "path",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "commit payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.commitPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.errorData"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "repos"
                ],
                "summary": "Delete blob",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "git ref (branch, tag, sha)",
                        "name": "ref",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "contents path",
                        "name": "path",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.errorData"
                        }
                    }
                }
            }
        },
        "/repos/{owner}/{repo}/branches": {
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
                    "repos"
                ],
                "summary": "Get branhces",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.branchList"
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
        },
        "/repos/{owner}/{repo}/tree/{ref}/{path}": {
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
                    "repos"
                ],
                "summary": "Get tree",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account owner of the repository (the name is not case sensitive)",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "the name of the repository (the name is not case sensitive)",
                        "name": "repo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "git ref (branch, tag, sha)",
                        "name": "ref",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "tree path",
                        "name": "path",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
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
        "api.blobEntry": {
            "type": "object",
            "properties": {
                "contents": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "api.branchItem": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "sha": {
                    "type": "string"
                }
            }
        },
        "api.branchList": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.branchItem"
                    }
                }
            }
        },
        "api.collectionPayload": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "api.commitPayload": {
            "type": "object",
            "properties": {
                "commitMessage": {
                    "type": "string"
                },
                "contents": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "api.deletePayload": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "api.entryPayload": {
            "type": "object",
            "properties": {
                "contents": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
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
                "defaultBranch": {
                    "type": "string"
                },
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
        },
        "api.treeItem": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                },
                "sha": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
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
	Title:            "Moonbase",
	Description:      "### Git-based headless CMS API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
