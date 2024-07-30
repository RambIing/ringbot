package commands

import (
	"github.com/bwmarrin/discordgo"
)

func KeyPad(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "enter_digits",
			Title:    "Enter digits",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "digits",
							Label:       "Enter digits to play over the phone",
							Style:       discordgo.TextInputShort,
							Placeholder: "01234567890*# are all acceptable to input", // 'w' as an input also works, and is documented to pause for 0.5s on it.
							Required:    true,
							MaxLength:   300,
							MinLength:   1,
						},
					},
				},
			},
		},
	})
}
