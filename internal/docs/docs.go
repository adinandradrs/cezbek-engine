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
            "name": "Kezbek Developer",
            "url": "https://kezbek.id",
            "email": "developer@kezbek.id"
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
        "/ping": {
            "get": {
                "description": "Ping the status of server, should be respond fastly.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Default"
                ],
                "summary": "Show the status of server.",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "403": {
                        "description": "Forbidden"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/v1/authorization/b2b": {
            "post": {
                "description": "API to authorize B2B officer account",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Authorization APIs"
                ],
                "summary": "API B2B Authorization",
                "parameters": [
                    {
                        "enum": [
                            "EBIZKEZBEK",
                            "B2BCLIENT"
                        ],
                        "type": "string",
                        "description": "Client Channel",
                        "name": "x-client-channel",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "android 10",
                        "description": "Client OS or Browser Agent",
                        "name": "x-client-os",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Device ID",
                        "name": "x-client-device",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "1.0.0",
                        "description": "Client Platform Version",
                        "name": "x-client-version",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Original Timestamp in UNIX format (EPOCH)",
                        "name": "x-client-timestamp",
                        "in": "header"
                    },
                    {
                        "description": "B2B Officer Authentication Payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.OfficerAuthenticationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/v1/authorization/client": {
            "post": {
                "description": "API to authorize client's signature and code",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Authorization APIs"
                ],
                "summary": "API Client Authorization",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Client signature using HMAC SHA256, signature formula is \u003cb\u003eHEX(HMAC(SHA256(UPPER(HTTP-METHOD):UPPER(CODE):UNIX-EPOCH:UPPER(API-KEY))))\u003c/b\u003e",
                        "name": "x-client-signature",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client API Key",
                        "name": "x-api-key",
                        "in": "header",
                        "required": true
                    },
                    {
                        "enum": [
                            "EBIZKEZBEK",
                            "B2BCLIENT"
                        ],
                        "type": "string",
                        "description": "Client Channel",
                        "name": "x-client-channel",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "android 10",
                        "description": "Client OS or Browser Agent",
                        "name": "x-client-os",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Device ID",
                        "name": "x-client-device",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "1.0.0",
                        "description": "Client Platform Version",
                        "name": "x-client-version",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Original Timestamp in UNIX format (EPOCH)",
                        "name": "x-client-timestamp",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Client Authentication Payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ClientAuthenticationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/v1/authorization/otp": {
            "post": {
                "description": "API to validate B2B officer account OTP",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Authorization APIs"
                ],
                "summary": "API B2B OTP Validation",
                "parameters": [
                    {
                        "enum": [
                            "EBIZKEZBEK",
                            "B2BCLIENT"
                        ],
                        "type": "string",
                        "description": "Client Channel",
                        "name": "x-client-channel",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "android 10",
                        "description": "Client OS or Browser Agent",
                        "name": "x-client-os",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Device ID",
                        "name": "x-client-device",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "1.0.0",
                        "description": "Client Platform Version",
                        "name": "x-client-version",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Original Timestamp in UNIX format (EPOCH)",
                        "name": "x-client-timestamp",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "Client Transaction ID",
                        "name": "x-client-trxid",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "B2B Officer Authentication Payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.OfficerValidationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/v1/cashbacks": {
            "post": {
                "description": "API to apply cashback on client's transaction",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client Cashback APIs"
                ],
                "summary": "API Apply Cashback",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer",
                        "description": "Your Token to Access",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "enum": [
                            "EBIZKEZBEK",
                            "B2BCLIENT"
                        ],
                        "type": "string",
                        "description": "Client Channel",
                        "name": "x-client-channel",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "android 10",
                        "description": "Client OS or Browser Agent",
                        "name": "x-client-os",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Device ID",
                        "name": "x-client-device",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "1.0.0",
                        "description": "Client Platform Version",
                        "name": "x-client-version",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Original Timestamp in UNIX format (EPOCH)",
                        "name": "x-client-timestamp",
                        "in": "header"
                    },
                    {
                        "description": "Transaction Payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.TransactionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "403": {
                        "description": "Forbidden"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    },
                    "503": {
                        "description": "Service Unavailable"
                    }
                }
            }
        },
        "/v1/cashbacks/{trxId}": {
            "post": {
                "description": "API to view detail of applied cashback based on the given Kezbek transaction reference",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client Cashback APIs"
                ],
                "summary": "API Detail Applied Cashback",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer",
                        "description": "Your Token to Access",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "enum": [
                            "EBIZKEZBEK",
                            "B2BCLIENT"
                        ],
                        "type": "string",
                        "description": "Client Channel",
                        "name": "x-client-channel",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "android 10",
                        "description": "Client OS or Browser Agent",
                        "name": "x-client-os",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Device ID",
                        "name": "x-client-device",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "1.0.0",
                        "description": "Client Platform Version",
                        "name": "x-client-version",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Original Timestamp in UNIX format (EPOCH)",
                        "name": "x-client-timestamp",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "Kezbek Transaction Reference",
                        "name": "trxId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "403": {
                        "description": "Forbidden"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    },
                    "503": {
                        "description": "Service Unavailable"
                    }
                }
            }
        },
        "/v1/partners": {
            "post": {
                "description": "API to register a new B2B Partner data as user and client",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Partner Management APIs"
                ],
                "summary": "API Add Partner",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer",
                        "description": "Your Token to Access",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "enum": [
                            "EBIZKEZBEK",
                            "B2BCLIENT"
                        ],
                        "type": "string",
                        "description": "Client Channel",
                        "name": "x-client-channel",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "android 10",
                        "description": "Client OS or Browser Agent",
                        "name": "x-client-os",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Device ID",
                        "name": "x-client-device",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "1.0.0",
                        "description": "Client Platform Version",
                        "name": "x-client-version",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Original Timestamp in UNIX format (EPOCH)",
                        "name": "x-client-timestamp",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "default": "PT. Lajada Piranti Commerce",
                        "description": "Partner Corporate",
                        "name": "partner",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "LAJADA",
                        "description": "Partner Code",
                        "name": "code",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "kezbek.support@lajada.net",
                        "description": "Partner Email",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "628123456789",
                        "description": "MSISDN",
                        "name": "msisdn",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "John Doe",
                        "description": "Partner Officer",
                        "name": "officer",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Bintaro Exchange Mall Blok A1",
                        "description": "Office Address",
                        "name": "address",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Logo",
                        "name": "logo",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "403": {
                        "description": "Forbidden"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    },
                    "503": {
                        "description": "Service Unavailable"
                    }
                }
            }
        }
    },
    "definitions": {
        "model.ClientAuthenticationRequest": {
            "type": "object",
            "required": [
                "code"
            ],
            "properties": {
                "code": {
                    "type": "string",
                    "example": "LAJADA"
                }
            }
        },
        "model.OfficerAuthenticationRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john.doe@lajada.id"
                }
            }
        },
        "model.OfficerValidationRequest": {
            "type": "object",
            "required": [
                "otp"
            ],
            "properties": {
                "otp": {
                    "type": "string",
                    "example": "123456"
                }
            }
        },
        "model.TransactionRequest": {
            "type": "object",
            "required": [
                "amount",
                "msisdn",
                "quantity"
            ],
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 750000
                },
                "email": {
                    "type": "string",
                    "example": "john.doe@gmailxyz.com"
                },
                "merchant_code": {
                    "type": "string",
                    "example": "LSAJA,GPAID,JOSVO"
                },
                "msisdn": {
                    "type": "string",
                    "example": "62812345678"
                },
                "quantity": {
                    "type": "integer",
                    "example": 2
                },
                "transaction_reference": {
                    "type": "string",
                    "example": "INV/001/002"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0-Beta",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Kezbek - Cashback Engine Sandbox",
	Description:      "This Cashback Engine Sandbox is only used for test and development purpose. To explore and serve all Kezbek operational APIs as a live data. It is not intended for production usage.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
