package main

import (
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var configPath = flag.String("config", "", "Path to config file")

func main() {
	flag.Parse()
	config, err := GetConfig(*configPath)
	if err != nil {
		log.Panic(err)
		return
	}
	httpClient, err := ProxyHttpClient(config.ProxyAddr, config.ProxyUser, config.ProxyPassword)
	if err != nil {
		log.Panic(err)
		return
	}
	bot, err := tgbotapi.NewBotAPIWithClient(config.Token, httpClient)
	if err != nil {
		log.Panic(err)
		return
	}

	cloudJenkinsClient, err := GetNewJenkinsClinet(JenkinsParams{
		JenkinsUrl: config.CloudJenkinsUrl,
		Password:   config.CloudJenkinsPassword,
		Username:   config.CloudJenkinsUser,
	})
	if err != nil {
		log.Panic(err)
		return
	}
	chmodJenkinsClient, err := GetNewJenkinsClinet(JenkinsParams{
		JenkinsUrl: config.ChmodJenkinsUrl,
		Password:   config.ChmodJenkinsPassword,
		Username:   config.ChmodJenkinsUser,
	})
	if err != nil {
		log.Panic(err)
		return
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
		return
	}

	// Do not handle a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		guid := uuid.NewV4()
		repoDir := fmt.Sprintf(config.GitRepoCloneDir + guid.String())

		handler := NewMessageHandler(
			HandlerParams{
				bot:                bot,
				cloudJenkinsClient: cloudJenkinsClient,
				chModJenkinsClient: chmodJenkinsClient,
				git:                GetNewGitClient(config.GitRepoUrl, repoDir, config.GitUser, config.GitEmail, config.GitPassword),
				archiveEditor:      GetNewArchiveEditor(repoDir),
				users:              config.TgUsers,
			},
		)
		go handler.Handle(update)
	}
}
