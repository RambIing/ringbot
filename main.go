package main

import (
	commands "RingBot/commands"
	"RingBot/serverManagement"
	"RingBot/settingsManager"
	"RingBot/twilio"
	websocket2 "RingBot/websocket"
	"net/http"
	"time"

	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"log"

	"os"
	"os/signal"
)

var (
	// Magic Numbers
	phoneOptionMinLength = 10 // minimum length a phone number can be without dashes or parenthesis
	phoneOptionMaxLength = 14 // max length assuming input is this: (123) 456-7890. Country code is not implemented,
	// although something you can add is accepting it.

)

var DgSession *discordgo.Session

var (
	/*
		Command handling, ugly code but at the fault of discordGo :c
	*/
	cmds = []*discordgo.ApplicationCommand{
		{
			Name:        "call",
			Description: "Calls a number. Accepts parenthesis and dashes in the phone number. Country code is not accepted.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "number",
					Description: "Phone number",
					MinLength:   &phoneOptionMinLength,
					MaxLength:   phoneOptionMaxLength,
					Required:    true,
				},
			},
		},
	}

	/*
		Handling each slash command.
	*/
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){

		"transfer_speaker": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.TransferSpeaker(s, i)
		},
		"mute_call": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.MuteCall(s, i)
		},
		"end_call": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.EndCall(s, i)
		},
		"key_pad": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.KeyPad(s, i)
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"call": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.CommandHandler(s, i)
		},
	}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	settingsManager.InitializeSettings()
	DgSession, _ = discordgo.New(fmt.Sprintf("Bot %s", settingsManager.GetBotToken()))

	// Intent for certain servers to view guild info
	DgSession.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	DgSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		s.UpdateGameStatus(0, "/call (number)")
		log.Printf("Logged in as: %v#%v \n", s.State.User.Username, s.State.User.Discriminator)
	})

	err := DgSession.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(cmds))
	for i, v := range cmds {
		cmd, err := DgSession.ApplicationCommandCreate(DgSession.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// Adds server to serverMap
	DgSession.AddHandler(func(s *discordgo.Session, guild *discordgo.GuildCreate) {
		serverManagement.AddServer(guild.ID)
	})

	DgSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		case discordgo.InteractionModalSubmit: //this is only connected to the sendDigits command
			data := i.ModalSubmitData()
			digits := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%s has entered %s on the keypad", i.Interaction.Member.User.Mention(), digits),
				},
			})
			err := twilio.PlayDigits(digits, i.Interaction.GuildID)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("There was an error entering %s: %e", digits, err),
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}
			time.Sleep(1 * time.Second)
			serverManagement.ServerMap[i.Interaction.GuildID].DigitsPlaying = false
		}
	})

	//
	// Websocket for mediastream
	//
	http.HandleFunc("/mediastream", func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
		)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("websocket upgrade failed: %e \n", err)
			return
		}
		websocket2.Reader(conn, DgSession)
	})
	http.ListenAndServe(":8000", nil)

	defer DgSession.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")

	// Disconnect from all voices otherwise it'll persist for a few minutes
	for _, voice := range DgSession.VoiceConnections {
		voice.Disconnect()
	}

	for _, cmd := range registeredCommands {
		err := DgSession.ApplicationCommandDelete(DgSession.State.User.ID, "", cmd.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
}
