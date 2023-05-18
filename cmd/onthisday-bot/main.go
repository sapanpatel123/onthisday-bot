package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sapanpatel123/onthisday-bot/internal/helper"
	"github.com/sapanpatel123/onthisday-bot/internal/search"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func truncate(s []string, to int) []string {
	return s[:to]
}

func run(path string, date string) (string, []string) {
	var msg string
	var files []string

	reqDate, _ := time.Parse("01/02", date)

	files, err := search.FindPhotos(path, reqDate)
	if err != nil {
		log.Fatal(err)
		return err.Error(), files
	}

	switch numFiles := len(files); {
	case numFiles == 0:
		msg = fmt.Sprintf("No photos found for %v", date)
		log.Println(msg)
	case numFiles > 10:
		msg = fmt.Sprintf("Found more than 10 photos for %v, only posting 5", date)
		files = truncate(files, 5)
		log.Println(msg)
	default:
		msg = fmt.Sprintf("Found %v photos for %v", numFiles, date)
		log.Print(msg)
	}

	return msg, files
}

func main() {
	var path string
	var jobChannel string
	var date string

	flag.StringVar(&path, "path", "~/", "path for the photos directory to use")
	flag.StringVar(&jobChannel, "job-channel", "test", "Slack channel to run as a daily job")

	flag.Parse()

	jobChannel = fmt.Sprintf("#%v", jobChannel)

	if path == "" {
		log.Fatal("Path argument not set")
	}

	isDir := helper.Exists(path)

	if !isDir {
		log.Fatalf("The path %v does not exist", path)
	}

	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	slackAppToken := os.Getenv("SLACK_APP_TOKEN")

	bot := slacker.NewClient(slackBotToken, slackAppToken)

	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	})

	bot.Command("run {date}", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			date = request.Param("date")
			if date == "" {
				date = time.Now().Format("01/02")
			}

			msg, files := run(path, date)

			apiClient := botCtx.ApiClient()
			event := botCtx.Event()
			if event.ChannelID != "" {
				_, ts, _ := apiClient.PostMessage(event.ChannelID, slack.MsgOptionText(msg, false))
				for _, file := range files {
					fOpen, err := os.Open(file)
					if err != nil {
						log.Fatalf("Error reading file: %v", err)
					}
					_, err = apiClient.UploadFile(slack.FileUploadParameters{
						Filename:        file,
						Reader:          fOpen,
						ThreadTimestamp: ts,
						Channels:        []string{event.ChannelID},
					})
					if err != nil {
						fmt.Printf("Error encountered when uploading file: %+v\n", err)
					}
				}
			}
		},
	})

	bot.Job("0 0 9 * * *", &slacker.JobDefinition{
		Description: "A cron job that runs every day at 9AM",
		Handler: func(jobCtx slacker.JobContext) {
			date = time.Now().Format("01/02")
			msg, files := run(path, date)

			apiClient := jobCtx.ApiClient()
			_, ts, _ := apiClient.PostMessage(jobChannel, slack.MsgOptionText(msg, false))
			for _, file := range files {
				fOpen, err := os.Open(file)
				if err != nil {
					log.Fatalf("Error reading file: %v", err)
				}
				_, err = apiClient.UploadFile(slack.FileUploadParameters{
					Filename:        file,
					Reader:          fOpen,
					ThreadTimestamp: ts,
					Channels:        []string{jobChannel},
				})
				if err != nil {
					fmt.Printf("Error encountered when uploading file: %+v\n", err)
				}
			}

		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatalf("Failed to start listening: %v", err)
	}
}
