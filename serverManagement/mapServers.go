package serverManagement

import (
	"sync"
	"time"
)

type ServerData struct {
	Tm               time.Time
	ConferenceUUID   string
	IDtoSSRC         map[string]int
	PhoneSID         string
	StreamSID        string
	CallID           string
	SpeakerPhoneID   string
	SpeakerPhoneSSRC int
	DigitsPlaying    bool
	Muted            bool
	Established      bool
	Chunk            int
	Buf              [][]byte
	InBuf            []byte
	Mu               sync.Mutex
}

var ServerMap = make(map[string]*ServerData)

func AddServer(id string) ServerData {
	ServerMap[id] = &ServerData{IDtoSSRC: make(map[string]int), Buf: make([][]byte, 0), InBuf: make([]byte, 0)}
	return *ServerMap[id]
}

func FindBasedOnUUID(uuid string) string {
	for key, val := range ServerMap {
		if val.ConferenceUUID == uuid {
			return key
		}
	}
	return ""
}

func FindBasedOnStreamSID(streamSid string) string {
	for key, val := range ServerMap {
		if val.StreamSID == streamSid {
			return key
		}
	}
	return ""
}
