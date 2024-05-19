package account_v1

import (
	"another_node/internal/community/storage"
	"another_node/internal/web_server/pkg"
	"another_node/internal/web_server/pkg/response"

	"github.com/gin-gonic/gin"
)

// Sync download accounts
// @Tags Account
// @Description download accounts
// @Accept json
// @Produce json
// @Success 200 {array} byte
// @Router /api/account/v1/sync [GET]
// @Param        count    query     int  true  "how many accounts to download"
// @Security JWT
func Sync(ctx *gin.Context) {
	if count, err := pkg.GetUIntParamFromQueryOrPath("count", ctx, false); err != nil {
		response.BadRequest(ctx, err)
		return
	} else {
		members := storage.GetAllMembers(uint32(count))

		b := storage.MarshalMembers(members)
		//community.SyncAccounts(count)
		response.Success(ctx, b)
	}
}
