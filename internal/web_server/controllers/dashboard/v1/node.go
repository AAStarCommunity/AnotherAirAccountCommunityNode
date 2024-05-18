package dashboard_v1

import (
	"another_node/internal/community"

	"github.com/gin-gonic/gin"
)

// Node get all members in network
// @Tags Dashboard
// @Description get node members
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /api/dashboard/v1/node [GET]
func Node(ctx *gin.Context) {
	members := community.ListNodes()

	ctx.JSON(200, members)
}
