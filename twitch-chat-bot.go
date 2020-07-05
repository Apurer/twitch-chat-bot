package main

import (
	WebSocket "github.com/Apurer/twitch-chat-bot/server/websocket"
	Server "github.com/Apurer/twitch-chat-bot/server"
	. "github.com/Apurer/twitch/structs"
	API "github.com/Apurer/twitch/api"
	IRC "github.com/Apurer/twitch/irc"
	"github.com/Apurer/eev/privatekey"
	"github.com/Apurer/eev"
	"flag"
)


func main() {

	key := flag.String("key", "", "path to private key which is to be used for dencryption of environment variable")
	passphrase := flag.String("passphrase", "", "passphrase by which private key is encrypted")
	flag.Parse()

	privkey, err := privatekey.Read(*key, *passphrase)
	if err != nil {
		panic(err)
	}

	oauthToken, err := eev.Get("OAUTH_TOKEN_TWITCH", privkey)
	if err != nil {
		panic(err)
	}

	username := "apurertv"
	channel := "apurertv"
	clientID := "r5lf6jgtrj9f7jxay85o4q4vz4v8xa"

	c := make(chan Followers)
	r := make(chan string)
	w := make(chan string)

	user := API.GetUser(username, oauthToken, clientID)
	userID := user.ID
	go API.GetFollowers(c, userID, oauthToken, clientID)
	var rw = WebSocket.IRC{Read: r, Write: w}
	go IRC.Chat(username, oauthToken, channel, c, r, w)
	Server.Host(rw)
}