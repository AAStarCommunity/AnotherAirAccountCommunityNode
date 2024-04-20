package middlewares

import (
	"another_node/conf"
	"another_node/internal/web_server/pkg/response"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GenericRecovery is a general error (panic) interception middleware, which intercepts and uniformly records possible errors.
func GenericRecovery() gin.HandlerFunc {
	DefaultErrorWriter := &PanicExceptionRecord{}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		// a unified response is made for the panic and other exceptions that occur
		errStr := ""
		if conf.Environment.IsDevelopment() {
			errStr = fmt.Sprintf("%v", err)
		}
		response.GetResponse().SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
	})
}

// PanicExceptionRecord  panic exception record
type PanicExceptionRecord struct{}

func (p *PanicExceptionRecord) Write(b []byte) (n int, err error) {
	s1 := "An error occurred in the server's internal code："
	var build strings.Builder
	build.WriteString(s1)
	build.Write(b)
	errStr := build.String()
	return len(errStr), errors.New(errStr)
}
