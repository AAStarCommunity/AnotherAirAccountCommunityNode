package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

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
