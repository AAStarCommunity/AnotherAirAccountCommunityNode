package main

import (
	"another_node/conf"
	"another_node/internal/community/node"
	"another_node/internal/community/storage"
	"another_node/internal/web_server/routers"
	"flag"
	"strconv"
	"strings"
)

func getFlags() (listen *uint16, name *string, joinAddrs *string, genesis *bool) {
	// 解析命令行参数
	listenTmp := flag.Uint("listen", 0, "Listen port number")
	listenVal := uint16(*listenTmp)
	listen = &listenVal
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

	defer func() {
		storage.Close()
	}()

	// start community node
	if _, err := node.New(getFlags()); err != nil {
		panic(err)
	}

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
