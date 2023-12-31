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
            "name": "@mager",
            "url": "https://geotory.com",
            "email": "magerleagues@gmail.com"
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
        "/datasets/{username}": {
            "get": {
                "description": "Fetch datasets from a given user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dataset"
                ],
                "summary": "Get all datasets for a user",
                "operationId": "get-datasets",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.DatasetsResp"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    }
                }
            }
        },
        "/datasets/{username}/{slug}": {
            "get": {
                "description": "Fetch details about a dataset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dataset"
                ],
                "summary": "Get a dataset",
                "operationId": "get-dataset",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Slug",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.DatasetResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    }
                }
            },
            "put": {
                "description": "Syncing a dataset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dataset"
                ],
                "summary": "Sync a dataset",
                "operationId": "sync-dataset",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Slug",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.DatasetResp"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a dataset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dataset"
                ],
                "summary": "Delete a dataset",
                "operationId": "delete-dataset",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Slug",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.DeleteDatasetResp"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    }
                }
            }
        },
        "/datasets/{username}/{slug}/deleteFeatures": {
            "post": {
                "description": "Delete features for a given dataset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dataset"
                ],
                "summary": "Delete features",
                "operationId": "delete-features",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Slug",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Delete Features Req",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.DeleteFeaturesReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.DatasetResp"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    }
                }
            }
        },
        "/datasets/{username}/{slug}/zip": {
            "get": {
                "description": "Download a dataset as a zip file",
                "produces": [
                    "application/zip"
                ],
                "tags": [
                    "dataset"
                ],
                "summary": "Download a dataset as a zip file",
                "operationId": "download-dataset-zip",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Slug",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResp"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_paulmach_go_geojson.Feature": {
            "type": "object",
            "properties": {
                "bbox": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "crs": {
                    "description": "Coordinate Reference System Objects are not currently supported",
                    "type": "object",
                    "additionalProperties": true
                },
                "geometry": {
                    "$ref": "#/definitions/github_com_paulmach_go_geojson.Geometry"
                },
                "id": {},
                "properties": {
                    "type": "object",
                    "additionalProperties": true
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "github_com_paulmach_go_geojson.FeatureCollection": {
            "type": "object",
            "properties": {
                "bbox": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "crs": {
                    "description": "Coordinate Reference System Objects are not currently supported",
                    "type": "object",
                    "additionalProperties": true
                },
                "features": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_paulmach_go_geojson.Feature"
                    }
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "github_com_paulmach_go_geojson.Geometry": {
            "type": "object",
            "properties": {
                "bbox": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "crs": {
                    "description": "Coordinate Reference System Objects are not currently supported",
                    "type": "object",
                    "additionalProperties": true
                },
                "geometries": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_paulmach_go_geojson.Geometry"
                    }
                },
                "lineString": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "number"
                        }
                    }
                },
                "multiLineString": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "array",
                            "items": {
                                "type": "number"
                            }
                        }
                    }
                },
                "multiPoint": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "number"
                        }
                    }
                },
                "multiPolygon": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "type": "number"
                                }
                            }
                        }
                    }
                },
                "point": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "polygon": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "array",
                            "items": {
                                "type": "number"
                            }
                        }
                    }
                },
                "type": {
                    "$ref": "#/definitions/github_com_paulmach_go_geojson.GeometryType"
                }
            }
        },
        "github_com_paulmach_go_geojson.GeometryType": {
            "type": "string",
            "enum": [
                "Point",
                "MultiPoint",
                "LineString",
                "MultiLineString",
                "Polygon",
                "MultiPolygon",
                "GeometryCollection"
            ],
            "x-enum-varnames": [
                "GeometryPoint",
                "GeometryMultiPoint",
                "GeometryLineString",
                "GeometryMultiLineString",
                "GeometryPolygon",
                "GeometryMultiPolygon",
                "GeometryCollection"
            ]
        },
        "handler.DatasetResp": {
            "type": "object",
            "properties": {
                "bbox": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "centroid": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "geojson": {
                    "$ref": "#/definitions/github_com_paulmach_go_geojson.FeatureCollection"
                },
                "id": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "types": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.DatasetType"
                    }
                },
                "updatedAt": {
                    "type": "string"
                },
                "user": {
                    "type": "object",
                    "properties": {
                        "image": {
                            "type": "string"
                        },
                        "slug": {
                            "type": "string"
                        }
                    }
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "handler.DatasetType": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "handler.Datasets": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                }
            }
        },
        "handler.DatasetsResp": {
            "type": "object",
            "properties": {
                "datasets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.Datasets"
                    }
                }
            }
        },
        "handler.DeleteDatasetResp": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "handler.DeleteFeaturesReq": {
            "type": "object",
            "properties": {
                "dataset": {
                    "type": "string"
                }
            }
        },
        "handler.ErrorResp": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "api.geotory.com",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Bluedot",
	Description:      "Primary backend for Geotory",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
