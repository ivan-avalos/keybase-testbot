package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

const backs = "```"

func makeAdvertisement() kbchat.Advertisement {
	helpExtended := fmt.Sprintf(`Get some help!
	Example:%s
			!help%s`, backs, backs)

	fuckExtended := fmt.Sprintf(`Tell someone to fuck! If you forget to mention someone, it will tell YOU to fuck!
	Example:%s
			!fuck
			!fuck @avalos%s`, backs, backs)

	cmds := []chat1.UserBotCommandInput{
		{
			Name:        "help",
			Description: "Get some help!",
			ExtendedDescription: &chat1.UserBotExtendedDescription{
				Title:       `*!help*`,
				DesktopBody: helpExtended,
				MobileBody:  helpExtended,
			},
		},
		{
			Name:        "fuck",
			Description: "Tell someone to fuck!",
			ExtendedDescription: &chat1.UserBotExtendedDescription{
				Title:       `*!fuck* [@username]`,
				DesktopBody: fuckExtended,
				MobileBody:  fuckExtended,
			},
		},
	}
	return kbchat.Advertisement{
		Alias: "Ávalos Testbot",
		Advertisements: []chat1.AdvertiseCommandAPIParam{
			{
				Typ:      "public",
				Commands: cmds,
			},
		},
	}
}

func main() {
	var kbLoc string
	var kbc *kbchat.API
	var err error

	flag.StringVar(&kbLoc, "keybase", "keybase", "/usr/bin/")
	flag.Parse()

	if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc}); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	if _, err = kbc.AdvertiseCommands(makeAdvertisement()); err != nil {
		fail("Error advertising commands: %s", err.Error())
	}

	sub, err := kbc.ListenForNewTextMessages()
	if err != nil {
		fail("Error listening: %s", err.Error())
	}

	for {
		msg, err := sub.Read()
		if err != nil {
			fail("Failed to read message: %s", err.Error())
		}
		if msg.Message.Content.TypeName != "text" {
			log.Println(msg.Message.Content.TypeName)
			continue
		}

		body := strings.TrimSpace(msg.Message.Content.Text.Body)
		answer := ""

		fmt.Println("Received: ", body)

		if m, _ := regexp.MatchString("^!help$", body); m {
			answer = "@" + msg.Message.Sender.Username + " No."
		} else if m, _ := regexp.MatchString("^!fuck$", body); m {
			answer = "Fuck you, " + "@" + msg.Message.Sender.Username
		} else if m, _ := regexp.MatchString("^!fuck @[a-zA-Z0-9_]+$", body); m {
			answer = "Fuck you, " + strings.Split(body, " ")[1]
		} else {
			continue
		}
		if _, err := kbc.SendMessageByConvID(msg.Message.ConvID, answer); err != nil {
			fail("Error echo'ing message: %s", err.Error())
		}
	}
}
