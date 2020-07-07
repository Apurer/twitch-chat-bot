package server 

import (
	. "github.com/Apurer/twitch-chat-bot/server/websocket/hub"
	. "github.com/Apurer/twitch-chat-bot/server/websocket"
	"net/http"
	"flag"
	"log"
)

var (
	listen = flag.String("listen", ":443", "listen address")
	dir    = flag.String("dir", "./wwwroot", "directory to serve")
)

func Host(irc IRC) {
	hub := NewHub()
	go hub.Run(*irc.Write)
	go hub.Chat(*irc.Read)
	fs := http.FileServer(http.Dir(*dir))
	http.Handle("/", fs)
	http.HandleFunc("/irc", func(w http.ResponseWriter, r *http.Request) {
		irc.RW(hub, w, r)
	})
	http.HandleFunc("/echo", Echo)
	err := http.ListenAndServeTLS(*listen, "server.rsa.crt", "server.rsa.key", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
