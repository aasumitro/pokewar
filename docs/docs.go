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
        "contact": {
            "name": "@aasumitro",
            "url": "https://aasumitro.id/",
            "email": "hello@aasumitro.id"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/aasumitro/pokewar/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/battles": {
            "get": {
                "description": "Get Battle List.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Battles"
                ],
                "summary": "Battle List",
                "parameters": [
                    {
                        "type": "string",
                        "description": "data limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "data offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "PAGINATION RESPOND",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.SuccessRespondWithPagination"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/domain.Battle"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "NOT FOUND",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorRespond"
                        }
                    },
                    "500": {
                        "description": "INTERNAL SERVER ERROR RESPOND",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorRespond"
                        }
                    }
                }
            }
        },
        "/api/v1/monsters": {
            "get": {
                "description": "Get Monster List.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Monsters"
                ],
                "summary": "Monster List",
                "parameters": [
                    {
                        "type": "string",
                        "description": "data limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "data offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "PAGINATION RESPOND",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.SuccessRespondWithPagination"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/domain.Monster"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "NOT FOUND",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorRespond"
                        }
                    },
                    "500": {
                        "description": "INTERNAL SERVER ERROR RESPOND",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorRespond"
                        }
                    }
                }
            }
        },
        "/api/v1/ranks": {
            "get": {
                "description": "Get Rank List.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Ranks"
                ],
                "summary": "Rank List",
                "parameters": [
                    {
                        "type": "string",
                        "description": "data limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "data offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "PAGINATION RESPOND",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.SuccessRespondWithPagination"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/domain.Rank"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "NOT FOUND",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorRespond"
                        }
                    },
                    "500": {
                        "description": "INTERNAL SERVER ERROR RESPOND",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorRespond"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Battle": {
            "type": "object",
            "properties": {
                "ended_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "logs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Log"
                    }
                },
                "players": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Player"
                    }
                },
                "started_at": {
                    "type": "integer"
                }
            }
        },
        "domain.Log": {
            "type": "object",
            "properties": {
                "battle_id": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "domain.Monster": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "base_exp": {
                    "type": "integer"
                },
                "height": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "origin_id": {
                    "type": "integer"
                },
                "skills": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Skill"
                    }
                },
                "stats": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Stat"
                    }
                },
                "types": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "weight": {
                    "type": "integer"
                }
            }
        },
        "domain.Player": {
            "type": "object",
            "properties": {
                "annulled_at": {
                    "type": "integer"
                },
                "avatar": {
                    "type": "string"
                },
                "battle_id": {
                    "type": "integer"
                },
                "eliminated_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "monster_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "point": {
                    "type": "integer"
                },
                "rank": {
                    "type": "integer"
                }
            }
        },
        "domain.Rank": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lose_battles": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "origin_id": {
                    "type": "integer"
                },
                "points": {
                    "type": "integer"
                },
                "total_battles": {
                    "type": "integer"
                },
                "types": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "win_battles": {
                    "type": "integer"
                }
            }
        },
        "domain.Skill": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "pp": {
                    "description": "Power Points",
                    "type": "integer"
                }
            }
        },
        "domain.Stat": {
            "type": "object",
            "properties": {
                "base_stat": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "utils.ErrorRespond": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "utils.Paging": {
            "type": "object",
            "properties": {
                "path": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "utils.SuccessRespond": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "status": {
                    "type": "string"
                }
            }
        },
        "utils.SuccessRespondWithPagination": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "current_page": {
                    "type": "integer"
                },
                "data": {},
                "next": {
                    "$ref": "#/definitions/utils.Paging"
                },
                "previous": {
                    "$ref": "#/definitions/utils.Paging"
                },
                "status": {
                    "type": "string"
                },
                "total_page": {
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
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
