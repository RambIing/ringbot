package commands

import (
	"RingBot/serverManagement"
	"github.com/bwmarrin/discordgo"
)

func MuteCall(s *discordgo.Session, i *discordgo.InteractionCreate) {
	server := serverManagement.ServerMap[i.GuildID]
	if server.Muted == false {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:        "Call is now muted 🔇",
			},
		})
		server.Muted = true
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:        "Call is now un-muted 🔊",
			},
		})
		server.Muted = false
	}
}