package settingsManager

import (
	twilio2 "RingBot/twilio"
	"encoding/json"
	"github.com/twilio/twilio-go"
	"io"
	"log"
	"os"
)

type Settings struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Websocket string `json:"websocket"`
	Token     string `json:"token"`
}

var botToken string

func InitializeSettings() {
	jsonFile, err := os.Open("settings.json")
	if err != nil {
		log.Fatalf("error opening json: %e \n", err)

	}
	defer jsonFile.Close()

	readBytes, _ := io.ReadAll(jsonFile)

	var settings Settings
	json.Unmarshal(readBytes, &settings)

	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: settings.Password,
		Username: settings.Username,
	})
	twilio2.Client = twilioClient
	twilio2.WebsocketURL = settings.Websocket
	botToken = settings.Token
}

func GetBotToken() string {
	return botToken
}
