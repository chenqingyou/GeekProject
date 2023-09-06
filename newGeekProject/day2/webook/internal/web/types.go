package web

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterRoutesCt(server *gin.Engine)
}
