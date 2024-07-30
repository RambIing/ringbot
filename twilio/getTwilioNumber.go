package twilio

import (
	"fmt"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// GetNumber pulls all active phone numbers from the Twilio account and returns the first one, for use in StartCall.
func GetNumber() (string, error) {
	params := &api.ListIncomingPhoneNumberParams{}
	params.SetLimit(1)

	resp, err := GetClient().Api.ListIncomingPhoneNumber(params)
	if err != nil {
		return "", fmt.Errorf("error grabbing phone")
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("no phone numbers on account")
	}
	return *resp[0].PhoneNumber, nil
}
