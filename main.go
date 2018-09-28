package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
	"flag"
)

func main() {
    token := flag.String("token", "default", "Bot User OAuth Access Token which starts with \"xoxb-\".")
    flag.Parse()

	api := slack.New(*token)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			fmt.Println(ev.User, ev.BotID)
			// Response only user message.
			if ev.User != "" {
				if (ev.Text == "hoge") {
					api.PostMessage(ev.Channel, "Please input message.", slack.NewPostMessageParameters())
				} else {
					api.PostMessage(ev.Channel, "Your message : " + ev.Text, slack.NewPostMessageParameters())
				}
			}

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
			fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
