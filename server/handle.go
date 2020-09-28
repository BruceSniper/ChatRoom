package server

import (
	"ChatRoom/logic"
	"net/http"
)

func RegisterHandle() {
	// 广播消息处理
	go logic.Broadcaster.Start() //启动广播器Goroutine

	http.HandleFunc("/", homeHandleFunc) // "/"代表首页
	http.HandleFunc("/user_list", userListHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc) // "/ws"用来服务Websocket长连接
}
