package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"os"
	"time"
)

var apiKey string = ""

var slackPolicy = "To use Slack you must have 2FA enabled as per the requirements and terms of use.  You will be reminded every 24 hours until you enabled 2FA.\n\n" +
	"Every so often we will disable accounts that do not have 2FA turned on.  To avoid this please turn on 2FA now.  Instructions for 2FA can be found at https://slack.zendesk.com/hc/en-us/articles/204509068-Enabling-two-factor-authentication."

func sendMessage(api slack.Client, channel string, message string) error {
	params := slack.PostMessageParameters{}
	params.LinkNames = 1
	_, _, err := api.PostMessage(channel, message, params)
	return err
}

func annoyUser(api slack.Client, user string) error {
	message := "You have been identied as a user that does not have 2 Factor Auth (2FA).\n\n" + slackPolicy

	err := sendMessage(api, "@"+user, message)
	return err
}

func shameUsers(api slack.Client, userString string) error {
	message := "Hello everyone...\n\n" +
		"The following users have not enabled 2 Factor auth " + userString + "\n\n" +
		slackPolicy

	err := sendMessage(api, "#general", message)
	return err
}

func getUsers(api slack.Client) error {
	users, err := api.GetUsers()
	if err != nil {
		return err
	}

	var userList = ""

	for _, member := range users {
		if member.Has2FA == false && member.Deleted == false && member.IsBot == false {
			if err := annoyUser(api, member.Name); err != nil {
				return err
			}
			userList += "@" + member.Name + " "
			fmt.Printf("%s\n", member.Name)
		}
	}

	if err := shameUsers(api, userList); err != nil {
		return err
	}

	return nil
}

func main() {

	apiKey = os.Getenv("SLACK_API_KEY")
	if apiKey == "" {
		println("You must set the enviroment variable SLACK_API_KEY")
		os.Exit(-1)
	}

	api := slack.New(apiKey)

	for {
		if err := getUsers(*api); err != nil {
			fmt.Printf("Error getting users and annoying them %s\n", err)
		}
		fmt.Printf("Running again in 24 hours...")
		time.Sleep(time.Hour * 24)
	}

}
