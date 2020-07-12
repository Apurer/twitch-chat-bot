package api

import (
	. "github.com/Apurer/twitch-chat-bot/structs/twitch"
	"encoding/json"
	"net/http"
	"time"
	"fmt"
)

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

		for k, v := range resp.Header {
			fmt.Printf("Header field %q, Value %q\n", k, v)
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