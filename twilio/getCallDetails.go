package twilio

import (
	"fmt"
	lookups "github.com/twilio/twilio-go/rest/lookups/v2"
)

func GetNumberDetails(call int64) interface{} {
	params := &lookups.FetchPhoneNumberParams{}
	params.SetFields("line_type_intelligence")
	resp, err := GetClient().LookupsV2.FetchPhoneNumber(fmt.Sprintf("+1%d", call), params)
	if err != nil {
		return err
	} else {
		if resp.LineTypeIntelligence != nil {
			return *resp.LineTypeIntelligence
		} else {
			return nil
		}
	}
}

func GetCallName(call int64) interface{} {
	params := &lookups.FetchPhoneNumberParams{}
	params.SetFields("caller_name")
	resp, err := GetClient().LookupsV2.FetchPhoneNumber(fmt.Sprintf("+1%d", call), params)
	if err != nil {
		return err
	} else {
		if resp.CallerName != nil {
			return *resp.CallerName
		} else {
			return nil
		}
	}
}
