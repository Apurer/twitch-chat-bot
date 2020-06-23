package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/Apurer/eev"
	"github.com/Apurer/eev/privatekey"
	"log"
	"time"
	"net/http"
	"io/ioutil"
)

func GetFollowers(oauthToken string) {

	url := "https://api.twitch.tv/helix/users/follows?to_id=488014220"
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Client-ID", "r5lf6jgtrj9f7jxay85o4q4vz4v8xa")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauthToken))

	client := &http.Client{}

	for {
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		time.Sleep(3 * time.Second)
	}
}

func Connect(user string, oauthToken string, channel string) {

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

	go GetFollowers(oauthToken)
	Connect("ApurerTV", oauthToken, "ApurerTV")
}
