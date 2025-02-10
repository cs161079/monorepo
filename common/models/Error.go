package models

import (
	"github.com/gin-gonic/gin"
)

const (
	INTERNALL_SERVER_ERROR = "INTERNAL SERVER ERROR"
	BAD_SYNTAX             = "REQUEST CONTAINS BAD SYNTAX OR CANNOT BE FULLFILLED"
)

type OasaError struct {
	Error string `json:"error" `
}

type CustomError struct {
	code        int32
	origMessage string
	userMessage string
}

func NewError(originalError, userMessage string, code int32) error {
	return &CustomError{
		code:        code,
		origMessage: originalError,
		userMessage: userMessage,
	}
}

func (errTracer *CustomError) Error() string {
	if errTracer.userMessage != "" {
		return errTracer.userMessage
	}

	return errTracer.origMessage
}

func HttpResponse(ctx *gin.Context, err error) {
	if wrappedErr, ok := err.(*CustomError); ok {
		// wrappedErr := err.(*CustomError)
		ctx.AbortWithStatusJSON(int(wrappedErr.code), map[string]string{"message": err.Error()})
	} else {
		ctx.AbortWithStatusJSON(500, map[string]string{"message": "Internal Server Error"})
	}
}
