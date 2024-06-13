package plugin

import (
	"another_node/internal/community/node"

	"github.com/gin-gonic/gin"
)

type RouteName string
type Auth bool

type HttpPlugin interface {
	RegisterRoutes(*gin.Engine, *node.Community)
}
