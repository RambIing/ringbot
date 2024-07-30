package commands

import (
	"RingBot/serverManagement"
	"RingBot/twilio"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
)

func JoinVoice(s *discordgo.Session, guildID, channelID, userID string) (*discordgo.VoiceConnection, error) {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		return nil, err
	}
	serverMap := serverManagement.ServerMap[guildID]
	serverMap.SpeakerPhoneID = userID
	// ssrc = the unique ID for each user in voice call. requires the user to speak before we can determine
	// used later on for determining packets that go to twilio as we don't want to pickup everyone in the vc
	vc.AddHandler(func(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
		serverMap.IDtoSSRC[vs.UserID] = vs.SSRC
		if vs.UserID == serverMap.SpeakerPhoneID {
			serverMap.SpeakerPhoneSSRC = vs.SSRC
		}
	})
	s.VoiceConnections[guildID] = vc
	return vc, nil
}

func CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	// grabs number from the command
	numString := ""
	if opt, ok := optionMap["number"]; ok {
		numString = opt.StringValue()
	}

	// Remove unnecessary values that can be inputted
	replacer := strings.NewReplacer("(", ")", "-", " ")
	strippedNum := replacer.Replace(numString)

	number, _ := strconv.Atoi(strippedNum)

	// default values, makes it easier to read
	guildID := i.Interaction.GuildID
	userID := i.Interaction.Member.User.ID
	serverManagement.AddServer(guildID)
	serverMap := serverManagement.ServerMap[guildID]
	guild, err := s.State.Guild(guildID)
	if err != nil {
		log.Printf("error getting guild info: %e \n", err)
		return
	}

	serverMap.SpeakerPhoneID = userID
	serverMap.SpeakerPhoneSSRC = serverMap.IDtoSSRC[userID]

	// finds what voice channel the sender is in
	for _, vs := range guild.VoiceStates {
		if vs.UserID == i.Interaction.Member.User.ID {
			_, err := JoinVoice(s, guild.ID, vs.ChannelID, vs.UserID)
			if err != nil {
				log.Printf("error joining voice: %e \n", err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: "Unknown error starting call. Check console."},
				})
			}
			err = twilio.StartCall(int64(number), i.Interaction.GuildID)
			if err != nil {
				if err.Error() == "no phone numbers on account" {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{Content: "No phone numbers on Twilio account."},
					})
					return
				}
				log.Printf("error starting call: %e \n", err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: "Unknown error starting call. Check console."},
				})
				return
			}
		}
	}

	numberDetails := twilio.GetNumberDetails(int64(number)) // shows carrier
	details := numberDetails.(map[string]interface{})
	callerID := twilio.GetCallName(int64(number)) // caller id
	caller := callerID.(map[string]interface{})
	currentNumber, err := twilio.GetNumber()
	if err != nil {
		log.Printf("error grabbing number: %e \n", err)
		return
	}

	// ugly structs, basically just buttons and embeds.
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Get Speaker",
							Style:    discordgo.SuccessButton,
							Disabled: false,
							Emoji:    discordgo.ComponentEmoji{Name: "🎤"},
							CustomID: "transfer_speaker",
						},
						discordgo.Button{
							Label:    "Key Pad",
							Style:    discordgo.SecondaryButton,
							Disabled: false,
							Emoji:    discordgo.ComponentEmoji{Name: "🔢"},
							CustomID: "key_pad",
						},
						discordgo.Button{
							Label:    "Mute/Unmute Call",
							Style:    discordgo.PrimaryButton,
							Disabled: false,
							Emoji:    discordgo.ComponentEmoji{Name: "🔇"},
							CustomID: "mute_call",
						},
						discordgo.Button{
							Label:    "End Call",
							Style:    discordgo.DangerButton,
							Disabled: false,
							Emoji:    discordgo.ComponentEmoji{Name: "📵"},
							CustomID: "end_call",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					URL:         "",
					Type:        "",
					Title:       fmt.Sprintf("Calling from %s", currentNumber),
					Description: "",
					Timestamp:   "",
					Color:       3447003,
					Footer:      nil,
					Image:       nil,
					Thumbnail:   nil,
					Video:       nil,
					Provider:    nil,
					Author:      nil,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Phone Number",
							Value:  fmt.Sprintf("+1%d", number),
							Inline: true,
						},
						{
							Name:   "Description",
							Value:  fmt.Sprintf("%v", caller["caller_type"]),
							Inline: true,
						},
						{
							Name:   "Carrier",
							Value:  fmt.Sprintf("%v", details["carrier_name"]),
							Inline: true,
						},
						{
							Name:   "Caller ID",
							Value:  fmt.Sprintf("%v", caller["caller_name"]),
							Inline: true,
						},
					},
				},
			},
		},
	})
}
