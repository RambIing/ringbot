package commands

import (
	"RingBot/serverManagement"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func TransferSpeaker(s *discordgo.Session, i *discordgo.InteractionCreate) {
	server := serverManagement.ServerMap[i.GuildID]
	server.SpeakerPhoneID = i.Interaction.Member.User.ID
	server.SpeakerPhoneSSRC = server.IDtoSSRC[i.Interaction.Member.User.ID]
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Speaker has been transferred to %s", i.Interaction.Member.User.Mention()),
		},
	})
}
