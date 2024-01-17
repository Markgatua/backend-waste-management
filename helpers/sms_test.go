package helpers

import (
	"testing"
	"ttnmwastemanagementsystem/configs"
)

func TestSendSMS(t* testing.T) {
	configs.InitEnvConfigs("../")
	sms := SMS{}
	error := sms.SendSMS([]string{"0791507732"}, "Test SMS")
	if error != nil {
		t.Error(error.Error())
	}
}
