package main

import (
	"another_node/conf"
	"another_node/internal/community"
	"another_node/internal/community/node"
	"another_node/internal/community/storage/migrations"
	"another_node/internal/web_server/routers"
	"strconv"
	"strings"
)

// @contact.name   AAStar Support
// @contact.url    https://aastar.xyz
// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description Type 'Bearer \<TOKEN\>' to correctly set the AccessToken
// @BasePath /api
func main() {

	// start community node
	if n, err := node.New(); err != nil {
		panic(err)
	} else {
		community.New(n)
	}

	migrations.AutoMigrate()

	// start web server
	routers.SetRouters().Run(func(port int) string {
		if port == 0 {
			panic("port is invalid")
		}

		portStr := strconv.Itoa(conf.GetWeb().Port)
		if !strings.HasPrefix(portStr, ":") {
			portStr = ":" + portStr
		}
		return portStr
	}(conf.GetWeb().Port))
}
