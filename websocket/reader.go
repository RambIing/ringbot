package websocket

import (
	"RingBot/serverManagement"
	"RingBot/twilio"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WebsocketMessage struct {
	Event     string       `json:"event"`
	StreamSid string       `json:"streamSid"`
	Media     Media        `json:"media"`
	Start     StartMessage `json:"start"`
}
type Media struct {
	Track     string `json:"track"`
	Chunk     string `json:"chunk"`
	Timestamp string `json:"timestamp"`
	Payload   string `json:"payload"`
}

type StartMessage struct {
	CustomParameters CustomParameters `json:"customParameters"`
}
type CustomParameters struct {
	ServerUUID string `json:"server"`
}

func Reader(conn *websocket.Conn, dgSession *discordgo.Session) {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading websocket: %s \n", err.Error())
			break
		}

		var mediaMessage WebsocketMessage
		err = json.Unmarshal(p, &mediaMessage)
		if err != nil {
			log.Printf("error unmarshaling json: %e \n", err)
			break
		}

		switch mediaMessage.Event {
		case "start":
			serverID := serverManagement.FindBasedOnUUID(mediaMessage.Start.CustomParameters.ServerUUID)
			serverMap := serverManagement.ServerMap[serverID]
			if serverID == "" {
				log.Fatalln("server not found in map!")
			}
			serverMap.StreamSID = mediaMessage.StreamSid
			if serverMap.Tm.IsZero() {
				serverMap.Tm = time.Now()
			}
			vc := dgSession.VoiceConnections[serverID]
			go func() {
				err := twilio.ProcessTwilioAndPlay(vc)
				if err != nil {
					log.Printf("error playing twilio audio: %e \n", err)
					return
				}
			}()
			go func() {
				twilio.DiscordToTwilio(vc, conn)
			}()
		case "media":
			if mediaMessage.Media.Track == "inbound" {
				serverID := serverManagement.FindBasedOnStreamSID(mediaMessage.StreamSid)
				if serverID == "" {
					log.Fatalln("server not found in map!")
				}
				err = twilio.ReadAndProcessTwilio(mediaMessage.Media.Payload, serverID)
				if err != nil {
					log.Printf("error processing audio: %e \n", err)
					break
				}
			}
		case "stop":
			serverID := serverManagement.FindBasedOnStreamSID(mediaMessage.StreamSid)
			serverMap := serverManagement.ServerMap[serverID]
			if serverMap.DigitsPlaying == true { // needs to be here for how playDigits works
				continue
			}
			serverMap.Established = false
			vc := dgSession.VoiceConnections[serverID]
			vc.Disconnect()
		}
	}
}
