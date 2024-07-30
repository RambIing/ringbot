package commands

import (
	"RingBot/serverManagement"
	"RingBot/twilio"
	"github.com/bwmarrin/discordgo"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func EndCall(s *discordgo.Session, i *discordgo.InteractionCreate) {
	serverMap := serverManagement.ServerMap[i.GuildID]
	params := &api.UpdateCallParams{}
	params.SetTwiml(`
				<Response>
					<Hangup />
				</Response>
			`)
	_, err := twilio.GetClient().Api.UpdateCall(serverMap.PhoneSID, params)
	if err != nil {
		// Errors when a call is ringing.
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Call is currently not active. (Might be ringing still)",
			},
		})
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Call has been ended.",
		},
	})
}
