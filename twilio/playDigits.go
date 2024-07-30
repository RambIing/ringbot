package twilio

import (
	"RingBot/serverManagement"
	"fmt"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"strings"
)
func PlayDigits(digits, guildId string) error {
	digitsBuilder := ""
	for _, r := range digits {
		if !strings.ContainsAny(string(r), "01234567890*#w") { // w = 0.5s delay according to twilio
			return fmt.Errorf("input did not contain a correct key")
		}
		digitsBuilder += fmt.Sprintf("%sw", string(r)) // appending 0.5 delay to allow any menu to process input
	}
	serverMap := serverManagement.ServerMap[guildId]
	serverMap.DigitsPlaying = true
	params := &api.UpdateCallParams{}
	params.SetTwiml(fmt.Sprintf(`
<Response>
	<Play digits="%s"></Play>
	<Connect>
		<Stream url="wss://%s/mediastream">
		<Parameter name="server" value="%s" />
		</Stream>
	</Connect>
	<Dial>
		<Conference beep="false" 
		waitUrl=""
		startConferenceOnEnter="true"
		endConferenceOnExit="true">%s</Conference>
	</Dial>
</Response>
	`, digitsBuilder, WebsocketURL, serverMap.ConferenceUUID, serverMap.ConferenceUUID))
	_, err := GetClient().Api.UpdateCall(serverMap.PhoneSID, params)
	if err != nil {
		return err
	}
	return nil
}