package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

const (
	saveDirectory = "savedTimezones"
	timezoneFile  = saveDirectory + "/timezones.json"
)

// Stores the user's timezones when loaded into memory
var timezones = make(map[string]string)

func main() {
	botToken := strings.Trim(os.Getenv("botToken"), "\n")
	if botToken == "" {
		log.Fatalf("Bot token is empty, you must specify one")
	}

	// Try and load the previously saved timezone data
	dataPresent := true
	loadedTimezones, err := os.ReadFile(timezoneFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("Pre-created timezone data was not found, not loading it in")
			dataPresent = false
		} else {
			log.Fatalf("Error reading timezone data from file %v: %v", timezoneFile, err)
		}
	}

	if dataPresent {
		err = json.Unmarshal(loadedTimezones, &timezones)
		if err != nil {
			log.Errorf("Could not read timezone data, skipping: %v", err)
		}
		log.Info("Loaded in saved timezone data")
	}

	// Make sure that the data gets saved when the bot quits
	defer func() {
		timezoneData, err := json.Marshal(&timezones)
		if err != nil {
			log.Errorf("Unable to convert timezones to JSON: %v", err)
			return
		}

		err = os.WriteFile(timezoneFile, timezoneData, 0755)
		if err != nil {
			log.Errorf("Unable to write timezones: %v", err)
		}
	}()

	log.Infof("Starting bot with token %v", botToken)
	bot, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}

	// Called when a message is sent on the server
	bot.AddHandler(messageCreate)

	bot.Identify.Intents = discordgo.IntentsGuildMessages

	if err := bot.Open(); err != nil {
		log.Fatalf("Error connecting to discord: %v", err)
	}
	log.Info("Started listening for requests")

	// Block until a message to stop is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Prevent the bot from messaging itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content[0] == '!' {
		log.Infof("Received command: %v", m.Content[1:])
		switch m.Content[1:strings.IndexRune(m.Content, ' ')] {
		// Sets a user's time zone
		case "settime":
			// Extract the time zone from the command
			timezone := strings.Trim(m.Content[len("settime")+1:], " ")
			zone, err := time.LoadLocation(timezone)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error setting timezone to [%v]: %v", timezone, err))
				return
			}
			// Person mentioned another person, set the time for them
			for _, otherPerson := range m.Mentions {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Set timezone for %v to %v", otherPerson.Mention(), timezone))
				timezones[otherPerson.ID] = zone.String()
			}

			// If nobody else is mentioned, set the own person's timezone
			if len(m.Mentions) <= 0 {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Set timezone to %v", zone))
				timezones[m.Author.ID] = zone.String()
			}
		// Sets their time
		case "timefor":
			authorLoc, err := time.LoadLocation(timezones[m.Author.ID])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error parsing timezone: "+err.Error())
				return
			}
			// Checks if the message author has their own timezone set, since the time will be individual for them
			if _, setAuthorTime := timezones[m.Author.ID]; !setAuthorTime {
				s.ChannelMessageSend(m.ChannelID, "It looks like you don't have your own timezone set. Please set one with !settime for the result to be displayed properly")
				return
			}

			// Loop through all the mentioned people and display their timezone
			for _, mentionedUser := range m.Mentions {
				if _, targetHasTimezone := timezones[mentionedUser.ID]; !targetHasTimezone {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %v has not set their timezone, set it with !settime", mentionedUser.Mention()))
					continue
				}
				// Parse the location of the user
				loc, err := time.LoadLocation(timezones[mentionedUser.ID])
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Could not load timezone of user %v: %v", mentionedUser.Mention(), err))
					continue
				}
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("From your perspective, the time for user %v is %v", mentionedUser.Mention(), time.Now().In(authorLoc).In(loc).Format(time.Kitchen)))
			}
		}
	}
}
