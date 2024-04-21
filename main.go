package main

import (
	"another_node/conf"
	"another_node/internal/community/node"
	"another_node/internal/global"
	"another_node/internal/web_server/routers"
	"strconv"
	"strings"
)

func main() {
	if n, err := node.New(); err != nil {
		panic(err)
	} else {
		global.CommunityInstance = &global.Community{
			Node: n,
		}
	}

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
