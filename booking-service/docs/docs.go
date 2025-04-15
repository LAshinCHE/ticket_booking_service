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
        "/": {
            "get": {
                "description": "Проверка доступности сервиса",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "hello",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/booking/{booking_id}": {
            "get": {
                "description": "Возвращает информацию о бронировании",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "booking"
                ],
                "summary": "Получить бронь по ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Booking ID",
                        "name": "booking_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.GetBookingByIDHandlerResponse"
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Booking": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "status": {
                    "$ref": "#/definitions/models.BookingStatus"
                },
                "tikcets": {
                    "type": "string"
                },
                "userID": {
                    "type": "string"
                }
            }
        },
        "models.BookingStatus": {
            "type": "string",
            "enum": [
                "draft",
                "reserved",
                "paid",
                "canceled"
            ],
            "x-enum-comments": {
                "BookingStatusCanceled": "Отменено",
                "BookingStatusDraft": "Черновик брони",
                "BookingStatusPaid": "Оплачено",
                "BookingStatusReserved": "Забронировано, но не оплачено"
            },
            "x-enum-varnames": [
                "BookingStatusDraft",
                "BookingStatusReserved",
                "BookingStatusPaid",
                "BookingStatusCanceled"
            ]
        },
        "types.GetBookingByIDHandlerResponse": {
            "type": "object",
            "properties": {
                "booking": {
                    "$ref": "#/definitions/models.Booking"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Booking Service API",
	Description:      "API для сервиса бронирования",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
