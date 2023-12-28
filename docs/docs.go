// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/syllabus": {
            "get": {
                "description": "シラバス全データを取得します。重すぎてswaggerで表示できないので注意。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tags"
                ],
                "summary": "シラバス全データ取得",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SyllabusViewModel"
                        }
                    }
                }
            }
        },
        "/syllabus/random": {
            "get": {
                "description": "シラバスデータ1つをランダムに取得します。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tags"
                ],
                "summary": "シラバスデータをランダム取得",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SyllabusViewModel"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.SyllabusViewModel": {
            "type": "object",
            "properties": {
                "credits": {
                    "type": "integer"
                },
                "day": {
                    "type": "string"
                },
                "faculty": {
                    "type": "string"
                },
                "lectureId": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "period": {
                    "type": "string"
                },
                "season": {
                    "type": "string"
                },
                "teacher": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
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
