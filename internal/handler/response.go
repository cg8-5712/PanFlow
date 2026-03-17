package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the unified API response structure
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Success returns a 200 response with data
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessMsg returns a 200 response with a custom message and no data
func SuccessMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: msg,
		Data:    nil,
	})
}

// Fail returns an error response with the given HTTP status, error code, and message
func Fail(c *gin.Context, httpStatus, code int, msg string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

// FailBadRequest returns a 400 response
func FailBadRequest(c *gin.Context, code int, msg string) {
	Fail(c, http.StatusBadRequest, code, msg)
}

// FailUnauthorized returns a 401 response
func FailUnauthorized(c *gin.Context, code int, msg string) {
	Fail(c, http.StatusUnauthorized, code, msg)
}

// FailForbidden returns a 403 response
func FailForbidden(c *gin.Context, code int, msg string) {
	Fail(c, http.StatusForbidden, code, msg)
}

// FailNotFound returns a 404 response
func FailNotFound(c *gin.Context, code int, msg string) {
	Fail(c, http.StatusNotFound, code, msg)
}

// FailInternal returns a 500 response
func FailInternal(c *gin.Context, msg string) {
	Fail(c, http.StatusInternalServerError, 50000, msg)
}
