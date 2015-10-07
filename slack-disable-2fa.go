package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var apiKey string = ""

type Users struct {
	Ok      bool   `json:"ok"`
	Error   string `json:error"`
	Members []struct {
		Fa      bool   `json:"has_2fa"`
		IsBot   bool   `json:"is_bot"`
		Deleted bool   `json:"deleted"`
		Name    string `json:"name"`
		Profile struct {
			Email string `json:"email"`
		} `json:"profile"`
	} `json:"members"`
}

type AppConfig struct {
	Slack_API_KEY string
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func annoyUser(user string) {
	var message = "You have been identied as a user that does not have 2 Factor Auth (2FA).\n" +
		"To you Slack you must have 2FA enabled as per the requirements and terms of use.  You will be reminded every 24 hours until you enabled 2FA.\n" +
		"Every so often we will disable accounts that do not have 2FA turned on.  To avoid this please turn on 2FA now.  Instructions for 2FA can be found at https://slack.zendesk.com/hc/en-us/articles/204509068-Enabling-two-factor-authentication."

	resp, err := http.PostForm("https://slack.com/api/chat.postMessage",
		url.Values{"token": {apiKey}, "channel": {"@" + user}, "text": {message}})
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}

func getUsers() {
	var data Users

	err := getJson("https://slack.com/api/users.list?token="+apiKey, &data)
	if err != nil {
		panic(err)
	}

	println(data.Error)
	if len(data.Error) > 0 {
		println("Invalid API Key for Slack.  Please correct.")
		os.Exit(-1)
	}

	for _, member := range data.Members {
		if member.Fa == false && member.Deleted == false && member.IsBot == false {
			annoyUser(member.Name)
			fmt.Printf("%s\n", member.Name)
		}
	}
}

func main() {

	apiKey = os.Getenv("SLACK_API_KEY")
	if apiKey == "" {
		println("You must set the enviroment variable SLACK_API_KEY")
		os.Exit(-1)
	}

	getUsers()

}
