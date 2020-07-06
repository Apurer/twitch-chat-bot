package websocket

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	. "github.com/Apurer/twitch-chat-bot/server/websocket/hub"
)

var upgrader = websocket.Upgrader{    
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type IRC struct {
	Read *chan string
	Write *chan string
}

func RWWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

func (channel *IRC) RW(w http.ResponseWriter, r *http.Request){
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	go func() { 
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			go func() { *channel.Write <- string(message) }()
			log.Printf("recv: %s", message)
		}
	}()

	for {
		select {
		case read := <-*channel.Read:
			err = c.WriteMessage(websocket.BinaryMessage, []byte(read))
			log.Printf("write: %s", read)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func Echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}