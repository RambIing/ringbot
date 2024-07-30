package twilio

import (
	"RingBot/serverManagement"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/zaf/g711"
	"layeh.com/gopus"
	"time"
)

func ProcessTwilioAndPlay(vc *discordgo.VoiceConnection) error {
	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Voip)
	if err != nil {
		return fmt.Errorf("new encoder error: %s", err.Error())
	}
	server := serverManagement.ServerMap[vc.GuildID]
	for {
		if vc == nil {
			break
		}
		time.Sleep(1 * time.Millisecond)
		for i, v := range server.Buf {
			if len(server.InBuf) <= 640 {
				server.Mu.Lock()
				server.InBuf = append(server.InBuf, v...)
				server.Mu.Unlock()
			} else {
				reader := bytes.NewReader(server.InBuf)
				intArr := make([]int16, reader.Len()/2)
				err = binary.Read(reader, binary.LittleEndian, &intArr)
				if err != nil {
					return fmt.Errorf("binary error: %s", err.Error())
				}

				opus, err := opusEncoder.Encode(intArr, 320, 1000)
				if err != nil {
					return fmt.Errorf("encoding Error: %s", err.Error())
				}
				go SendToDiscord(vc, opus)
				server.InBuf = make([]byte, 0)
				if i+1 > len(server.Buf) {
					continue
				}
				server.Mu.Lock()
				server.Buf = append(server.Buf[:0], server.Buf[i+1:]...)
				server.Mu.Unlock()
			}
		}
	}
	return err
}

func SendToDiscord(vc *discordgo.VoiceConnection, opus []byte) {
	vc.OpusSend <- opus
}

const (
	channels  int = 1    // 1 for mono, 2 for stereo
	frameRate int = 8000 // audio sampling rate
)

func ReadAndProcessTwilio(payloadString, guildID string) error {
	server := serverManagement.ServerMap[guildID]
	byteArray, err := base64.StdEncoding.DecodeString(payloadString)
	if err != nil {
		return fmt.Errorf("decode error: %s", err.Error())
	}
	byteArray = g711.DecodeUlaw(byteArray)
	server.Mu.Lock()
	server.Buf = append(server.Buf, byteArray)
	server.Mu.Unlock()
	return nil
}
