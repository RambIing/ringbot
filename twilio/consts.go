package twilio

import "github.com/twilio/twilio-go"

// Client password & username are in settings.json
var Client *twilio.RestClient

// WebsocketURL update this in settings.json
var WebsocketURL = ""

// GetClient returns initialized Twilio client
func GetClient() *twilio.RestClient {
	return Client
}
