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
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/exibition": {
            "get": {
                "description": "Get nft info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exibition"
                ],
                "summary": "Get specific NFT",
                "responses": {}
            }
        },
        "/file-upload": {
            "post": {
                "description": "upload file and thumbnail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "upload file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "file you want to upload",
                        "name": "upload_file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "thumbnail file",
                        "name": "thumbnail_file",
                        "in": "formData"
                    }
                ],
                "responses": {}
            }
        },
        "/getNFTInfoWithId": {
            "get": {
                "description": "Get nft info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NFT"
                ],
                "summary": "Get specific NFT",
                "parameters": [
                    {
                        "type": "string",
                        "description": "nft_id",
                        "name": "nft_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/like": {
            "post": {
                "description": "update like",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Like"
                ],
                "summary": "update like",
                "parameters": [
                    {
                        "description": "id of user",
                        "name": "UserId",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "id of work",
                        "name": "WorkId",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "type": "file",
                        "description": "test",
                        "name": "test",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/work/specific": {
            "get": {
                "description": "Get works",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Work"
                ],
                "summary": "get specific work",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/work/top10": {
            "get": {
                "description": "get top 10 works",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Work"
                ],
                "summary": "get top 10 works",
                "responses": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "34.212.84.161",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "NFTime Sample Swagger API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
