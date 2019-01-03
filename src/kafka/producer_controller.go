package kafka

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeResponse(context *gin.Context, code int, value interface{}) {
	context.JSON(
		code, gin.H{
			"data": value,
		},
	)
}

func PostProducerDataHandler(context *gin.Context) {
	var params producerParams
	if err := context.ShouldBind(&params); err != nil {
		return
	}
	if params.Values == "" {
		MakeResponse(context, http.StatusBadRequest, "check params")
		return
	}
	ProducerServer.WithAccessLog(params.Topic, params.Key, params.Values)
	MakeResponse(context, http.StatusOK, params)
}
