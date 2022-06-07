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
        "termsOfService": "TODO",
        "contact": {
            "name": "TODO",
            "url": "TODO",
            "email": "TODO"
        },
        "license": {
            "name": "GPL-3.0 license",
            "url": "https://www.gnu.org/licenses/gpl-3.0.en.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/geopoint": {
            "post": {
                "description": "create the geopoint in the database and receive the sound and picture file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Geopoint"
                ],
                "summary": "create a geopoint",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "example": [
                            0,
                            1,
                            2,
                            3,
                            45,
                            3,
                            2,
                            1
                        ],
                        "name": "amplitudes",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "example": "2022-05-26T11:17:35.079344Z",
                        "name": "date",
                        "in": "formData"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "number"
                        },
                        "example": [
                            38.652608,
                            -120.357448
                        ],
                        "name": "location",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "example": "Forêt à l'aube",
                        "name": "title",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "example": 1,
                        "name": "user_id",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "geopoint sound",
                        "name": "sound",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "geopoint picture",
                        "name": "picture",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/geopoint.GeoPoint"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/geopoint/{id}": {
            "get": {
                "description": "retrieve the geopoint in the database using its name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Geopoint"
                ],
                "summary": "get a geopoint",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "geopoint id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/geopoint.GeoPoint"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/geopoint/{id}/picture": {
            "get": {
                "description": "create the geopoint in the database and receive the sound and picture file",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Geopoint"
                ],
                "summary": "get the url of the picture",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "geopoint id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "used to check if api is alive",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "pings the api",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/user": {
            "post": {
                "description": "create a user in the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "create user",
                "parameters": [
                    {
                        "description": "Add user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.AddUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/user.AddUser"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/user/{username}": {
            "get": {
                "description": "retrieve the user in the database using its name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "get a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user name",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "geopoint.GeoPoint": {
            "type": "object",
            "properties": {
                "amplitudes": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    },
                    "example": [
                        0,
                        1,
                        2,
                        3,
                        45,
                        3,
                        2,
                        1
                    ]
                },
                "created_on": {
                    "type": "string",
                    "example": "2022-05-26T11:17:35.079344Z"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "location": {
                    "$ref": "#/definitions/geopoint.Location"
                },
                "picture": {
                    "type": "string",
                    "example": "https://example.com/picture-1.jpg"
                },
                "sound": {
                    "type": "string",
                    "example": "https://example.com/sound-2.mp3"
                },
                "title": {
                    "type": "string",
                    "example": "Forêt à l'aube"
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "geopoint.Location": {
            "type": "object",
            "properties": {
                "todo": {
                    "type": "string"
                }
            }
        },
        "httputil.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        },
        "user.AddUser": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "bob"
                }
            }
        },
        "user.User": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "created_on": {
                    "type": "string",
                    "example": "2022-05-26T11:17:35.079344Z"
                },
                "last_login": {
                    "type": "string",
                    "example": "2022-05-26T11:17:35.079344Z"
                },
                "name": {
                    "type": "string",
                    "example": "bob"
                },
                "token": {
                    "type": "string",
                    "example": "auinrsetanruistnstnaustie"
                },
                "user_id": {
                    "type": "integer",
                    "example": 123
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Swagger for biophonie-api",
	Description:      "API of biophonie (https://secret-garden-77375.herokuapp.com/).",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
