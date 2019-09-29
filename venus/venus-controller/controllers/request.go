package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// CheckRequestForm will check binding of form data
func CheckRequestForm(c *gin.Context, form interface{}) error {
	bindErr := c.ShouldBind(form)

	if bindErr != nil {
		if c.Request != nil {
			fmt.Println("invalid form", c.Request.URL, bindErr.Error())
		}
		ResponseInvalidParams(c, bindErr.Error())
	}
	return bindErr
}

// CheckRequestBody will check binding of request body
func CheckRequestBody(c *gin.Context, param interface{}) error {
	bindErr := c.ShouldBindJSON(param)
	if bindErr != nil {
		if c.Request != nil {
			fmt.Println("invalid body", c.Request.URL, bindErr.Error())
		}
		ResponseInvalidParams(c, bindErr.Error())
	}
	return bindErr
}

// CheckRequestQuery will check binding of request query
func CheckRequestQuery(c *gin.Context, query interface{}) error {
	bindErr := c.ShouldBindQuery(query)
	fmt.Println(query)
	if bindErr != nil {
		if c.Request != nil {
			fmt.Println("invalid body", c.Request.URL, bindErr.Error())
		}
		ResponseInvalidParams(c, bindErr.Error())
	}

	return bindErr
}
