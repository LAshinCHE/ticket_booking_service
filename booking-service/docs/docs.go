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
        "/booking/": {
            "post": {
                "description": "Создание бронирования пользователем",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "booking"
                ],
                "summary": "Создать бронь",
                "parameters": [
                    {
                        "description": "Данные для создания брони",
                        "name": "booking",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateBookingData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.CreateBookingResponse"
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
        },
        "/internal/booking/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Создание бронирования внутренним сервисом",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "internal"
                ],
                "summary": "Внутреннее создание бронирования",
                "parameters": [
                    {
                        "description": "Данные для создания брони",
                        "name": "booking",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.CreateBookingInternalRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.CreateBookingResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request",
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
        },
        "/internal/booking/delete": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Удаление бронирования внутренним сервисом",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "internal"
                ],
                "summary": "Внутреннее удаление бронирования",
                "parameters": [
                    {
                        "description": "Данные для удаления брони",
                        "name": "booking",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.DeleteBookingInternalRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "failed to create booking",
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
                    "type": "integer"
                },
                "status": {
                    "$ref": "#/definitions/models.BookingStatus"
                },
                "ticketID": {
                    "type": "integer"
                },
                "userID": {
                    "type": "integer"
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
        "models.CreateBookingData": {
            "type": "object",
            "properties": {
                "price": {
                    "type": "number"
                },
                "ticketID": {
                    "type": "integer"
                },
                "userID": {
                    "type": "integer"
                }
            }
        },
        "types.CreateBookingInternalRequest": {
            "type": "object",
            "properties": {
                "ticket_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "types.CreateBookingResponse": {
            "type": "object",
            "properties": {
                "bookingID": {
                    "type": "integer"
                }
            }
        },
        "types.DeleteBookingInternalRequest": {
            "type": "object",
            "properties": {
                "booking_id": {
                    "type": "integer"
                }
            }
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
	Host:             "localhost:8081",
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
