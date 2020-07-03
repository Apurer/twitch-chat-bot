package irc

import (
	. "github.com/Apurer/twitch-chat-bot/structs"
	"crypto/tls"
	"fmt"
	"log"
)

func AppreciateFollowers(channel string, con *tls.Conn, c chan Followers){
	var compare_followers []Follower
	for {
		select {
        case followers := <-c:
			if len(compare_followers) > 0 {
				var new_followers []string
				
				for _, follower := range followers.Data {
					found := false
					for _, compare_follower := range compare_followers {
						if follower.FromID == compare_follower.FromID {
							found = true
							break
						}
					}
					if found == false {
						new_followers = append(new_followers, follower.FromName)
					}
				}

				for _, username := range new_followers {
					fmt.Fprintf(con, fmt.Sprintf("PRIVMSG #%s :Dziekuje za follow %s\r\n", channel, username)) 
				}
			}
			compare_followers = make([]Follower, len(followers.Data))
			copy(compare_followers, followers.Data)
		}
	}
}

func Chat(user string, oauthToken string, channel string, c chan Followers) {

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

	go AppreciateFollowers(channel, con, c)

	var msg = make([]byte, 1024)
	var b int
	b, _ = con.Read(msg)
	PING := "PING :tmi.twitch.tv\r\n"
	PONG := "PONG :tmi.twitch.tv\r\n"
	for b != 0 {
		comment := string(msg[:b])
		fmt.Printf("%s\n", comment)

		if comment == PING {
			fmt.Printf("%s\n", PONG)
			fmt.Fprintf(con, PONG)
		}

		b = 0
		b, _ = con.Read(msg)
	}
}