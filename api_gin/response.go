package apigin

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func messageResponse(message string) gin.H {
	return gin.H{"message": message}
}

func errorResponse(err error) gin.H{
	return gin.H{"error" : err.Error()}
}

func newErrorResponse(message string) gin.H{
	return gin.H{"error" : errors.New(message)}
}