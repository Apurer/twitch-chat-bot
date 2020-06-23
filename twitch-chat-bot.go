package main

import (
	"github.com/Apurer/eev/privatekey"
	"github.com/Apurer/eev"
	"encoding/json"
	"crypto/tls"
	"net/http"
	"time"
	"flag"
	"fmt"
	"log"
)

type User struct {
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	DisplayName     string `json:"display_name"`
	ID              string `json:"id"`
	Login           string `json:"login"`
	OfflineImageURL string `json:"offline_image_url"`
	ProfileImageURL string `json:"profile_image_url"`
	Type            string `json:"type"`
	ViewCount       int64  `json:"view_count"`
}

type Users struct {
	Data []User
}

type Followers struct {
	Data []struct {
		FollowedAt string `json:"followed_at"`
		FromID     string `json:"from_id"`
		FromName   string `json:"from_name"`
		ToID       string `json:"to_id"`
		ToName     string `json:"to_name"`
	} `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
	Total int64 `json:"total"`
}

func GetFollowers(channel chan<- Followers, userID string, oauthToken string, clientID string) {

	followers := new(Followers)

	url := fmt.Sprintf("https://api.twitch.tv/helix/users/follows?to_id=%s", userID)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauthToken))

	client := &http.Client{}

	for {
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(followers)
		if err != nil {
			panic(err)
		}
		channel <- *followers
		time.Sleep(3 * time.Second)
	}
}

func PrintFollowers(channel chan Followers){
	for {
		select {
        case followers := <-channel:
            fmt.Println(followers)
		}
	}
}

func GetUser(username string, oauthToken string, clientID string) User {
	users := new(Users)
	
	url := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", username)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauthToken))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(users)
	if err != nil {
		panic(err)
	}

	if e:=len(users.Data); e != 1 {
		panic(fmt.Errorf("number: %d of elements returned by request regarding user: %s does not equal 1", e, username))
	}

	return users.Data[0]
}

func IRCchat(user string, oauthToken string, channel string) {

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

	username := "apurertv"
	channel := "apurertv"
	clientID := "r5lf6jgtrj9f7jxay85o4q4vz4v8xa"
	
	c := make(chan Followers)
	user := GetUser(username, oauthToken, clientID)
	userID := user.ID
	go GetFollowers(c, userID, oauthToken, clientID)
	go PrintFollowers(c)
	IRCchat(username, oauthToken, channel)
}
