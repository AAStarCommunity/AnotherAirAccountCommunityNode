package plugin_passkey_relay_party

import (
	"another_node/plugins/passkey_relay_party/conf"
	"another_node/plugins/passkey_relay_party/seedworks"
	"time"

	jwt2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type login struct {
	Account    string
	AccontType seedworks.AccountType
}

var jwtMiddleware *jwt2.GinJWTMiddleware

func ginJwtMiddleware() *jwt2.GinJWTMiddleware {
	return jwtMiddleware
}

func AuthHandler() gin.HandlerFunc {
	m, _ := jwt2.New(&jwt2.GinJWTMiddleware{
		Realm:       conf.GetJwtKey().Realm,
		Key:         []byte(conf.GetJwtKey().Security),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour / 2,
		IdentityKey: "jti",
		PayloadFunc: func(data interface{}) jwt2.MapClaims {
			if v, ok := data.(*login); ok {
				return jwt2.MapClaims{
					"jti": uuid.New().String(),
					"sub": v.Account,
				}
			}
			return jwt2.MapClaims{}
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"code":   code,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			if user, ok := c.Get("user"); !ok {
				email := c.Query("email")
				account := c.Query("account")
				if len(email) == 0 && len(account) == 0 {
					return nil, jwt2.ErrFailedAuthentication
				}

				return &login{
					Account: func() string {
						if len(email) != 0 {
							return email
						}
						return account
					}(),
					AccontType: func() seedworks.AccountType {
						if len(email) != 0 {
							return seedworks.Email
						}
						if c.Query("type") == "" {
							return seedworks.EOA
						}
						return seedworks.AccountType(c.Query("type"))
					}(),
				}, nil
			} else {
				if u, ok := user.(*seedworks.User); ok {
					return &login{
						Account: func() string {
							account, _ := u.GetAccount()
							return account
						}(),
						AccontType: func() seedworks.AccountType {
							_, accountType := u.GetAccount()
							return accountType
						}(),
					}, nil
				} else {
					return nil, jwt2.ErrFailedAuthentication
				}
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		HTTPStatusMessageFunc: func(e error, c *gin.Context) string {
			return "unauthorized"
		},
	})

	jwtMiddleware = m

	return m.MiddlewareFunc()
}
