package main

import (
	"another_node/conf"
	"another_node/internal/community"
	"another_node/internal/community/node"
	"another_node/internal/community/storage/migrations"
	"another_node/internal/web_server/routers"
	"flag"
	"strconv"
	"strings"
)

func getFlags() (listen *int, name *string, joinAddrs *string, genesis *bool) {
	// 解析命令行参数
	listen = flag.Int("listen", 0, "Listen port number")
	name = flag.String("name", "", "Node name")
	joinAddrs = flag.String("join", "", "Addresses of nodes to join (comma-separated)")
	genesis = flag.Bool("genesis", false, "Is this genesis node")
	flag.Parse()

	return
}

// @contact.name   AAStar Support
// @contact.url    https://aastar.xyz
// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description Type 'Bearer \<TOKEN\>' to correctly set the AccessToken
// @BasePath /api
func main() {

	// start community node
	if n, err := node.New(getFlags()); err != nil {
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
