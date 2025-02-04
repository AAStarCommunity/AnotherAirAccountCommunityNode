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
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var httpPlugins = make([]plugin.HttpPlugin, 0)

func preventSuspend(host string) {

	go func() {
		t := time.Tick(time.Second * 30)
		for {
			if resp, err := http.Get(host); err != nil {
				log.Default().Printf("error: " + err.Error())
			} else {
				if b, err := io.ReadAll(resp.Body); err == nil {
					log.Default().Printf("health: " + string(b))
				}
			}
			<-t
		}
	}()
}

func getFlags() (listen *uint16, name *string, joinAddrs *string, genesis *bool) {
	// 解析命令行参数
	listenTmp := flag.Int("listen", 0, "Listen port number")
	name = flag.String("name", "", "Node name")
	joinAddrs = flag.String("join", "", "Addresses of nodes to join (comma-separated)")
	explicitGenesisVar := flag.Bool("genesis", false, "Is this genesis node")
	flag.Parse()

	if listenTmp == nil {
		m := uint16(7946)
		listen = &m
	} else {
		m := uint16(*listenTmp)
		listen = &m
	}

	flag.CommandLine.Visit(func(f *flag.Flag) {
		if f.Name == "genesis" {
			genesis = explicitGenesisVar
		}
	})
	return
}

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @contact.name   AAStar Support
// @contact.url    https://aastar.xyz
// @BasePath /api
func main() {

	if healthz := os.Getenv("healthz"); len(healthz) > 0 {
		preventSuspend(healthz)
	}

	defer func() {
		storage.Close()
	}()

	// start community node
	if community, err := node.New(getFlags()); err != nil {
		panic(err)
	} else {
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
