package controllers

import (
	"siren/venus/venus-controller/controllers/errors"

	"github.com/gin-gonic/gin"
)

func ResponseError(c *gin.Context, e errors.Error) {
	c.AbortWithStatusJSON(e.HTTPStatus, errors.ErrorResponse{
		Errors: []*errors.Error{&e},
	})
}

func ResponseWithErrorCode(c *gin.Context, errorCode *errors.ErrorCode, message string) {
	ResponseError(c, errors.MakeErrorWithErrorCode(errorCode, message))
}

func ResponseDBError(c *gin.Context, message string) {
	e := errors.MakeDBError(message)
	ResponseError(c, e)
}

func ResponseTokenInvalid(c *gin.Context, errorCode *errors.ErrorCode, message string) {
	c.AbortWithStatusJSON(errorCode.HTTPStatus, errors.ErrorResponse{
		Errors: []*errors.Error{&errors.Error{ErrorCode: *errorCode, Detail: message}},
	})
}

func ResponseInvalidParams(c *gin.Context, message string) {
	e := errors.MakeInvalidaParamsError(message)
	ResponseError(c, e)
}

func ResponseNotFound(c *gin.Context, message string) {
	e := errors.MakeNotFoundError(message)
	ResponseError(c, e)
}

func ResponseUnexpected(c *gin.Context, message string) {
	e := errors.MakeUnexpectedError(message)
	ResponseError(c, e)
}

func ResponsePictureError(c *gin.Context, message string) {
	e := errors.MakePictureError(message)
	ResponseError(c, e)
}

func ResponseNotPermitted(c *gin.Context) {
	e := errors.MakeErrorWithErrorCode(&errors.ErrorNotPermitted, errors.ErrorNotPermitted.Title)
	ResponseError(c, e)
}
