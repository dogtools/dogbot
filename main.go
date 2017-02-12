package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nlopes/slack"

	"github.com/dogtools/dog"
)

var (
	apiKey      string
	botID       string
	dogfilePath string
	dogfile     dog.Dogfile
)

func main() {
	flag.StringVar(&apiKey, "key", "", "Slack Bot API Key")
	flag.StringVar(&botID, "id", "", "Bot ID")
	flag.StringVar(&dogfilePath, "dogfile", ".", "Dogfile path")
	flag.Parse()

	apiKey = os.Getenv("DOGBOT_API_KEY")
	botID = os.Getenv("DOGBOT_BOT_ID")

	if apiKey == "" || botID == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var err error
	dogfile, err = dog.ParseFromDisk(dogfilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	api := slack.New(apiKey)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		handleMessageEvent(rtm, msg)
	}
}

func handleMessageEvent(rtm *slack.RTM, msg slack.RTMEvent) {

	var helpMessage = `Hello human friend! Can I help you run some tasks?
• Type ` + "`@dogbot list`" + ` to list all tasks
• Type ` + "`@dogbot taskname`" + ` to run a task
`

	switch ev := msg.Data.(type) {
	case *slack.MessageEvent:
		if strings.HasPrefix(ev.Text, "<@"+botID+">") {

			parts := strings.Fields(ev.Text)
			if len(parts) == 2 {
				switch parts[1] {
				case "help":
					// print help
					rtm.SendMessage(rtm.NewOutgoingMessage(helpMessage, ev.Channel))

				case "list":
					// print list of tasks
					taskList := "I can run any of the following tasks, just ask me typing `@dobgot taskname`\n"
					for _, t := range dogfile.Tasks {
						taskList += fmt.Sprintf("• *%s*: _%s_\n", t.Name, t.Description)
					}
					rtm.SendMessage(rtm.NewOutgoingMessage(taskList, ev.Channel))

				default:
					taskName := parts[1]
					taskChain, err := dog.NewTaskChain(dogfile, taskName)
					if err != nil {
						rtm.SendMessage(rtm.NewOutgoingMessage(err.Error(), ev.Channel))
						break
					}

					startTime := time.Now()
					err = taskChain.Run(os.Stdout, os.Stderr)
					if err != nil {
						rtm.SendMessage(rtm.NewOutgoingMessage(err.Error(), ev.Channel))
						break
					}

					finishMsg := fmt.Sprintf("*%s* finished after %s", taskName, time.Since(startTime).String())
					rtm.SendMessage(rtm.NewOutgoingMessage(finishMsg, ev.Channel))
				}
			} else {
				noUnderstandMsg := "I don't understand what you are saying. Type `@dobgot help` or `@dogbot list`."
				rtm.SendMessage(rtm.NewOutgoingMessage(noUnderstandMsg, ev.Channel))
			}
		}

	case *slack.RTMError:
		fmt.Printf("Error: %s\n", ev.Error())

	case *slack.InvalidAuthEvent:
		fmt.Printf("Invalid credentials")
		return
	}
}
