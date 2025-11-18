package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mutition/go_start/common/tracing"
)

type BaseResponse struct {}

type response struct {
	Errno int `json:"errno"`
	Message string `json:"message"`
	Data any `json:"data"`
	TraceID string `json:"trace_id"`
}

func (base *BaseResponse) Response(c *gin.Context, err error, data interface{}) {
	if err != nil {
		base.Error(c, err)
	}
	base.Success(c, data)
}

func (base *BaseResponse) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{
		Errno: 0,
		Message: "success",
		Data: data,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}

func (base *BaseResponse) Error(c *gin.Context, err error) *BaseResponse {
	c.JSON(http.StatusOK, response{
		Errno: 2,
		Message: err.Error(),
		Data: nil,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
	return base
}