package irc

import (
	. "github.com/Apurer/twitch-chat-bot/structs/twitch"
	"crypto/tls"
	"fmt"
	"log"
)

func AppreciateFollowers(channel string, con *tls.Conn, c chan Followers){
	var compare_followers []string
	
	for {
		select {
        case followers := <-c:
			if len(compare_followers) > 0 {
				for _, follower := range followers.Data {
					found := false
					for _, compare_follower := range compare_followers {
						if follower.FromID == compare_follower {
							found = true
							break
						}
					}
					if found == false {
						fmt.Fprintf(con, fmt.Sprintf("PRIVMSG #%s :Dziekuje za follow %s\r\n", channel, follower.FromName))
						compare_followers = append(compare_followers, follower.FromID)
					}
				}
			} else {
				for _, follower := range followers.Data {
					compare_followers = append(compare_followers, follower.FromID)
				}
			}
		}
	}
}

func Chat(user string, oauthToken string, channel string, f chan Followers, r chan []byte, w chan []byte) {

	conf := &tls.Config{}
	con, err := tls.Dial("tcp", "irc.chat.twitch.tv:6697", conf)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintf(con, fmt.Sprintf("PASS oauth:%s\r\n", oauthToken))

	fmt.Fprintf(con, "CAP REQ :twitch.tv/tags\r\n")
	fmt.Fprintf(con, fmt.Sprintf("NICK %s\r\n", user)) 
	fmt.Fprintf(con, fmt.Sprintf("USER %s\r\n", user))
	fmt.Fprintf(con, fmt.Sprintf("JOIN #%s\r\n", channel))

	go AppreciateFollowers(channel, con, f)

	var msg = make([]byte, 1024)
	var b int
	b, _ = con.Read(msg)
	PING := "PING :tmi.twitch.tv\r\n"
	PONG := "PONG :tmi.twitch.tv\r\n"

	go func() {
		for {
			select {
			case write := <-w:
				fmt.Fprintf(con, fmt.Sprintf("PRIVMSG #%s :%s\r\n", channel, write))
			}
		}
    }()

	for b != 0 {
		message := msg[:b]
		fmt.Printf("%s\n", message)
		
		if PING == string(message) {
			fmt.Printf("%s\n", PONG)
			fmt.Fprintf(con, PONG)
		} else {
			r <- message
		}
		
		b = 0
		b, _ = con.Read(msg)
	}
}