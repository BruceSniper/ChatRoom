package main

import (
	"ChatRoom/global"
	"ChatRoom/server"
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"
)

var (
	addr   = ":2022"
	banner = `
    ____                _____
   |     |    |    /\     |
   |     |____|   /  \    | 
   |     |    |  /----\   |
   |____ |    | /      \  |
BruceCode ChatRoom，start on：%s
`
)

func init() {
	global.Init()
}

func main() {
	fmt.Printf(banner, addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
