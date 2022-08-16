// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://github.com/DAN-AND-DNA?tab=stars",
        "contact": {
            "name": "Snk技术开发中心",
            "email": "danyang.chen@snkad.cn"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1.user-service.Create_company": {
            "post": {
                "description": "创建公司",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-service"
                ],
                "summary": "创建公司",
                "parameters": [
                    {
                        "description": "args",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user_service.Create_company_args"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "reply",
                        "schema": {
                            "$ref": "#/definitions/user_service.Create_company_reply"
                        }
                    }
                }
            }
        },
        "/v1.user-service.Create_system": {
            "post": {
                "description": "创建系统",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-service"
                ],
                "summary": "创建系统",
                "parameters": [
                    {
                        "description": "args",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user_service.Create_system_args"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "reply",
                        "schema": {
                            "$ref": "#/definitions/user_service.Create_system_reply"
                        }
                    }
                }
            }
        },
        "/v1.user-service.Create_user": {
            "post": {
                "description": "创建用户",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-service"
                ],
                "summary": "创建用户",
                "parameters": [
                    {
                        "description": "args",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user_service.Create_user_args"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "reply",
                        "schema": {
                            "$ref": "#/definitions/user_service.Create_user_reply"
                        }
                    }
                }
            }
        },
        "/v1.user-service.Gen_id": {
            "post": {
                "description": "创建id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-service"
                ],
                "summary": "创建id",
                "parameters": [
                    {
                        "description": "args",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user_service.Gen_id_args"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "reply",
                        "schema": {
                            "$ref": "#/definitions/user_service.Gen_id_reply"
                        }
                    }
                }
            }
        },
        "/v1.user-service.Get_auth_info": {
            "post": {
                "description": "获得授权信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-service"
                ],
                "summary": "获得授权信息",
                "parameters": [
                    {
                        "description": "args",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user_service.Get_auth_info_args"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "reply",
                        "schema": {
                            "$ref": "#/definitions/user_service.Get_auth_info_reply"
                        }
                    }
                }
            }
        },
        "/v1.user-service.Get_idle_toutiao_appid": {
            "post": {
                "description": "获得空闲头条appid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-service"
                ],
                "summary": "获得空闲头条appid",
                "parameters": [
                    {
                        "description": "args",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user_service.Get_idle_toutiao_appid_args"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "reply",
                        "schema": {
                            "$ref": "#/definitions/user_service.Get_idle_toutiao_appid_reply"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "errcode.Errcode": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "错误码",
                    "type": "integer"
                },
                "from": {
                    "description": "来源 {3000: user_service}",
                    "type": "integer"
                },
                "msg": {
                    "description": "错误消息",
                    "type": "string",
                    "example": "错误消息"
                }
            }
        },
        "user_service.Create_company_args": {
            "type": "object",
            "properties": {
                "company_describe": {
                    "type": "string",
                    "example": "公司描述"
                },
                "company_name": {
                    "type": "string",
                    "example": "公司名"
                },
                "sysid": {
                    "description": "系统id",
                    "type": "integer",
                    "example": 1
                },
                "tid": {
                    "type": "string",
                    "example": "事务id"
                }
            }
        },
        "user_service.Create_company_reply": {
            "type": "object",
            "properties": {
                "err": {
                    "$ref": "#/definitions/errcode.Errcode"
                }
            }
        },
        "user_service.Create_system_args": {
            "type": "object",
            "properties": {
                "system_describe": {
                    "type": "string",
                    "example": "系统描述"
                },
                "system_name": {
                    "type": "string",
                    "example": "系统名称"
                },
                "tid": {
                    "type": "string",
                    "example": "事务id"
                }
            }
        },
        "user_service.Create_system_reply": {
            "type": "object",
            "properties": {
                "err": {
                    "$ref": "#/definitions/errcode.Errcode"
                }
            }
        },
        "user_service.Create_user_args": {
            "type": "object",
            "properties": {
                "cid": {
                    "description": "公司id",
                    "type": "integer",
                    "example": 1
                },
                "password": {
                    "type": "string",
                    "example": "密码"
                },
                "tid": {
                    "type": "string",
                    "example": "事务id"
                },
                "username": {
                    "type": "string",
                    "example": "用户名"
                }
            }
        },
        "user_service.Create_user_reply": {
            "type": "object",
            "properties": {
                "err": {
                    "$ref": "#/definitions/errcode.Errcode"
                }
            }
        },
        "user_service.Gen_id_args": {
            "type": "object",
            "properties": {
                "tid": {
                    "type": "string",
                    "example": "事务id"
                },
                "type": {
                    "description": "1: 公司 2: 用户",
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "user_service.Gen_id_reply": {
            "type": "object",
            "properties": {
                "err": {
                    "$ref": "#/definitions/errcode.Errcode"
                },
                "id": {
                    "description": "产生的id",
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "user_service.Get_auth_info_args": {
            "type": "object"
        },
        "user_service.Get_auth_info_reply": {
            "type": "object",
            "properties": {
                "err": {
                    "$ref": "#/definitions/errcode.Errcode"
                }
            }
        },
        "user_service.Get_idle_toutiao_appid_args": {
            "type": "object"
        },
        "user_service.Get_idle_toutiao_appid_reply": {
            "type": "object",
            "properties": {
                "err": {
                    "$ref": "#/definitions/errcode.Errcode"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1.0",
	Host:        "",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "user-service",
	Description: "用户服务",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}