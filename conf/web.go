package conf

import "sync"

type Web struct {
	Port int
}

var web *Web
var onceWeb sync.Once

// GetJwtKey 获取JWT私钥
func GetWeb() *Web {
	onceWeb.Do(func() {
		if web == nil {
			j := getConf().Web
			web = &Web{
				Port: j.Port,
			}
		}
	})

	return web
}
