package twilio

import (
	"RingBot/serverManagement"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"github.com/zaf/g711"
	"layeh.com/gopus"
	"log"
	"time"
)

func DiscordToTwilio(vc *discordgo.VoiceConnection, conn *websocket.Conn) {
	decoder, err := gopus.NewDecoder(8000, 1)
	if err != nil {
		panic(err)
	}
	serverMap := serverManagement.ServerMap[vc.GuildID]
	c := vc.OpusRecv
	for p := range c {
		if serverMap.Muted == true {
			continue
		}
		if int(p.SSRC) == serverMap.SpeakerPhoneSSRC {
			p.PCM, err = decoder.Decode(p.Opus, 960, false)
			if err != nil {
				log.Printf("bad opus decode: %e \n", err)
			}
			err = ReadAndProcessDiscord(p.PCM, conn, serverMap)
			if err != nil {
				log.Printf("error processing discord audio: %e \n", err)
			}
			serverMap.Chunk++
		}
	}
}

type WebsocketMessage struct {
	Event     string `json:"event"`
	StreamSid string `json:"streamSid"`
	Media     Media  `json:"media"`
}
type Media struct {
	Track     string `json:"track"`
	Chunk     string `json:"chunk"`
	Timestamp string `json:"timestamp"`
	Payload   string `json:"payload"`
}

func ReadAndProcessDiscord(pcm []int16, conn *websocket.Conn, server *serverManagement.ServerData) error {
	lpcm := make([]byte, 0)
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, pcm)
	if err != nil {
		return err
	}
	for i := 0; i < len(buf.Bytes()); i++ {
		lpcm = append(lpcm, buf.Bytes()[i])
	}
	ulaw := g711.EncodeUlaw(lpcm)
	t := &WebsocketMessage{
		Event:     "media",
		StreamSid: server.StreamSID,
		Media: Media{
			Track:     "outbound",
			Chunk:     fmt.Sprintf("%d", server.Chunk),
			Timestamp: fmt.Sprintf("%d", time.Since(server.Tm).Milliseconds()),
			Payload:   base64.StdEncoding.EncodeToString(ulaw),
		},
	}
	marshal, err := json.Marshal(t)
	if err != nil {
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		return err
	}
	return nil
}
