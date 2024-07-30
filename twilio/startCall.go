package twilio

import (
	"RingBot/serverManagement"
	"fmt"
	"github.com/google/uuid"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"log"
)

func StartCall(call int64, guildID string) error {
	params := &api.CreateCallParams{}
	// hacky way to create bi-directional audio through twilio!
	// uuid needed to prevent two servers from hearing each other's calls
	serverMap := serverManagement.ServerMap[guildID]
	serverMap.ConferenceUUID = uuid.NewString()

	params.SetTwiml(fmt.Sprintf(`
<Response>
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
	`, WebsocketURL, serverMap.ConferenceUUID, serverMap.ConferenceUUID))
	params.SetTo(fmt.Sprintf("+1%d", call))
	phone, err := GetNumber()
	if err != nil {
		log.Printf("error grabbing number: %e \n", err)
		return err
	}
	params.SetFrom(phone)

	resp, err := GetClient().Api.CreateCall(params)
	if err != nil {
		return err
	} else {
		if resp.Sid != nil {
			serverMap.PhoneSID = *resp.Sid
			return nil
		} else {
			return fmt.Errorf("no sid")
		}
	}
}
