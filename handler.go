package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

type HandlerParams struct {
	bot                *tgbotapi.BotAPI
	cloudJenkinsClient JenkinsClient
	chModJenkinsClient JenkinsClient
	git                GitClient
	archiveEditor      ArchiveEditor
	users              []string
}

type UpdateHandler interface {
	Handle(update tgbotapi.Update)
}

func NewMessageHandler(params HandlerParams) UpdateHandler {
	return &handler{
		archiveEditor:      params.archiveEditor,
		bot:                params.bot,
		cloudJenkinsClient: params.cloudJenkinsClient,
		chModJenkinsClient: params.chModJenkinsClient,
		git:                params.git,
		users:              params.users,
	}
}

type handler struct {
	archiveEditor      ArchiveEditor
	bot                *tgbotapi.BotAPI
	cloudJenkinsClient JenkinsClient
	chModJenkinsClient JenkinsClient
	git                GitClient
	users              []string
	send               func(m string) error
}

func (h *handler) Handle(update tgbotapi.Update) {
	h.send = func(m string) error {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, m)
		//msg.ReplyToMessageID = update.Message.MessageID
		_, err := h.bot.Send(msg)
		h.checkError(err)
		return nil
	}
	var auth bool
	for _, u := range h.users {
		if update.Message.From.UserName == u {
			auth = true
			break
		}
	}
	if !auth {
		err := h.send("Forbidden")
		h.checkError(err)
		return
	}
	t := update.Message.Text
	if t == "/start" {
		err := h.send("Please, choose the action:\n" +
			"/archive_builds - show last archive builds\n" +
			"/runvm_agent_builds - show last runvm builds\n" +
			"/runvm_controller_builds - show last controller builds\n" +
			"/vpn_manager_builds - show last vpn manager builds\n" +
			"/ganesha_builds - show last ganesha builds\n" +
			"/aakore_builds - show last aakore builds\n" +
			"/iaas_management_builds - show last iaas management tool builds\n" +
			"/runvm_image_builds - show last runvm image builds\n")
		if err != nil {
			h.checkError(err)
			return
		}
	} else if t == "/archive_builds" {
		h.printBuilds(h.cloudJenkinsClient, "runvm_archive_builder")
		return
	} else if t == "/runvm_agent_builds" {
		h.printBuilds(h.chModJenkinsClient, "mod-runvm-agent")
		return
	} else if t == "/runvm_controller_builds" {
		h.printBuilds(h.chModJenkinsClient, "mod-runvm-controller")
		return
	} else if t == "/vpn_manager_builds" {
		h.printBuilds(h.chModJenkinsClient, "mod-vpn-manager-image-builder")
		return
	} else if t == "/ganesha_builds" {
		h.printBuilds(h.chModJenkinsClient, "mod-ganesha-archive-plugin")
		return
	} else if t == "/aakore_builds" {
		h.printBuilds(h.chModJenkinsClient, "mod-aakore")
		return
	} else if t == "/iaas_management_builds" {
		h.printBuilds(h.chModJenkinsClient, "mod-iaas-management")
		return
	} else if t == "/runvm_image_builds" {
		h.printBuilds(h.chModJenkinsClient, "mod-runvm-image-builder")
		return
	} else if strings.Contains(update.Message.Text, "/change_vpn_build") {
		params := strings.TrimSpace(strings.Replace(update.Message.Text, "/change_vpn_build ", "", -1))
		paramsList := strings.Split(params, " ")
		build, err := strconv.Atoi(strings.TrimSpace(paramsList[1]))
		if err != nil {
			h.checkError(err)
			return
		}
		h.changeBuilds(
			strings.TrimSpace(paramsList[0]),
			"vpn",
			build,
		)
	}
}

func (h *handler) printBuilds(jk JenkinsClient, job string) {
	builds, err := jk.GetLastBuilds(job)
	if err != nil {
		h.checkError(err)
		return
	}
	for _, b := range builds {
		err = h.send(fmt.Sprintf("%s %d", b.branch, b.number))
		if err != nil {
			h.checkError(err)
			return
		}
	}
}

func (h *handler) changeBuilds(branch string, buildName string, build int) {
	err := h.git.Pull(branch)
	if err != nil {
		h.checkError(err)
		return
	}
	err = h.archiveEditor.ChangeBuild(buildName, build)
	if err != nil {
		h.checkError(err)
		return
	}
	err = h.git.CommitAndPush(fmt.Sprintf("New %s %d", buildName, build))
	if err != nil {
		h.checkError(err)
		return
	}
}

func (h *handler) checkError(err error) {
	if err != nil {
		errSend := h.send(err.Error())
		if errSend != nil {
			log.Print(errSend)
		}
		log.Print(err)
	}
}
