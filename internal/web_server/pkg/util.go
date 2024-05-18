package pkg

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUIntParamFromQueryOrPath(paramKey string, ctx *gin.Context, allowNull bool) (uint, error) {

	keyStr := ctx.Param(paramKey)
	if len(keyStr) == 0 {
		keyStr = ctx.Query(paramKey)
	}
	if allowNull {
		if len(keyStr) == 0 {
			return 0, nil
		}
	}
	key, err := strconv.ParseUint(keyStr, 10, 64)
	if err != nil {
		return 0, err
	} else {
		return uint(key), nil
	}
}

func GetBoolParamFromQuery(paramKey string, ctx *gin.Context, defaultValue bool) bool {

	keyStr := ctx.Query(paramKey)
	if len(keyStr) == 0 {
		return defaultValue
	}

	return strings.EqualFold("true", keyStr)
}
