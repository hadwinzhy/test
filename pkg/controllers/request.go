package controllers

import (
	"fmt"
	"siren/pkg/controllers/errors"
	"siren/pkg/logger"

	"github.com/gin-gonic/gin"
)

// CheckRequestForm will check binding of form data
func CheckRequestForm(c *gin.Context, form interface{}) error {
	bindErr := c.ShouldBind(form)

	if bindErr != nil {
		logger.L.Info("error happened checking form", bindErr.Error())
		errors.ResponseInvalidParams(c, bindErr.Error())
	}
	return bindErr
}

// CheckRequestBody will check binding of request body
func CheckRequestBody(c *gin.Context, param interface{}) error {
	bindErr := c.ShouldBindJSON(param)
	if bindErr != nil {
		logger.L.Info("error happened checking body", bindErr.Error())
		errors.ResponseInvalidParams(c, bindErr.Error())
	}
	return bindErr
}

// CheckRequestQuery will check binding of request query
func CheckRequestQuery(c *gin.Context, query interface{}) error {
	bindErr := c.ShouldBindQuery(query)
	fmt.Println(query)
	if bindErr != nil {
		logger.L.Info("error happened checking query", bindErr.Error())
		errors.ResponseInvalidParams(c, bindErr.Error())
	}

	return bindErr
}
