package main

import (
	"another_node/conf"
	"another_node/internal/community/node"
	"another_node/internal/community/storage"
	"another_node/internal/plugin"
	"another_node/internal/web_server/routers"
	plugin_passkey_relay_party "another_node/plugins/passkey_relay_party"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var httpPlugins = make([]plugin.HttpPlugin, 0)

func getFlags() (listen *uint16, name *string, joinAddrs *string, genesis *bool) {
	// 解析命令行参数
	listenTmp := flag.Uint("listen", 0, "Listen port number")
	listenVal := uint16(*listenTmp)
	listen = &listenVal
	name = flag.String("name", "", "Node name")
	joinAddrs = flag.String("join", "", "Addresses of nodes to join (comma-separated)")
	genesisStr := flag.String("genesis", "", "Is this genesis node")
	if *genesisStr == "" {
		fmt.Println("genesis is empty")
		genesis = nil
	} else {
		genesis = flag.Bool("genesis", false, "Is this genesis node")
	}
	fmt.Println("[getFlags] : listen:", listen, "name:", name, "join:", joinAddrs, "genesis:", genesis)
	return
}

// @contact.name   AAStar Support
// @contact.url    https://aastar.xyz
// @BasePath /api
func main() {

	defer func() {
		storage.Close()
	}()

	// start community node
	if community, err := node.New(getFlags()); err != nil {
		panic(err)
	} else {
		storage.Init()
		routers := routers.SetRouters()

		// load http plugins
		for _, plugin := range RegisterHttpPlugin() {
			plugin.RegisterRoutes(routers, community)
		}
		// start web server
		routers.Run(func(port int) string {
			if port == 0 {
				panic("port is invalid")
			}
			fmt.Println("community node created Port is: ", port)
			portStr := strconv.Itoa(conf.GetWeb().Port)
			if !strings.HasPrefix(portStr, ":") {
				portStr = ":" + portStr
			}
			return portStr
		}(conf.GetWeb().Port))
	}
}

func RegisterHttpPlugin() []plugin.HttpPlugin {
	if len(httpPlugins) == 0 {
		httpPlugins = append(httpPlugins, plugin_passkey_relay_party.NewRelay())
	}

	return httpPlugins
}
