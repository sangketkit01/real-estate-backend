package apigin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) HomePage(ctx *gin.Context){
	ctx.JSON(http.StatusOK, gin.H{"message":"Hello from back-end"})
}