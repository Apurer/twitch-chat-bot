package server 

import (
	//WebSocket "github.com/Apurer/twitch-chat-bot/server/websocket"
	. "github.com/Apurer/twitch-chat-bot/server/websocket"
	. "github.com/Apurer/twitch-chat-bot/server/websocket/hub"
	"log"
	"flag"
	"net/http"
)

var (
	listen = flag.String("listen", ":443", "listen address")
	dir    = flag.String("dir", "./wwwroot", "directory to serve")
)

func Host(irc IRC) {
	hub := NewHub()
	go hub.Run()
	fs := http.FileServer(http.Dir(*dir))
	http.Handle("/", fs)
	//http.HandleFunc("/irc", irc.RW)
	http.HandleFunc("/irc", func(w http.ResponseWriter, r *http.Request) {
		RWWs(hub, w, r)
	})
	http.HandleFunc("/echo", Echo)
	err := http.ListenAndServeTLS(*listen, "server.rsa.crt", "server.rsa.key", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
