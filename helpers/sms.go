package helpers

import (
	"ttnmwastemanagementsystem/configs"
	"ttnmwastemanagementsystem/logger"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	_"gopkg.in/guregu/null.v3"
)

const TAG string = "helper/SMS"

type SMS struct{}

type AfricasTalkingSendMessageResponse struct {
	SMS `json:"SMSMessageData"`
}

func (sms SMS) SendSMS(recipients []string, message string) error {
	if len(recipients) == 0 {
		logger.Log(TAG, "Empty recipients", logger.LOG_LEVEL_ERROR)
		return errors.New("Empty recipients")
	}

	if len(strings.TrimSpace(message)) == 0 {
		logger.Log(TAG, "Empty message", logger.LOG_LEVEL_ERROR)
		return errors.New("Empty message")
	}
	return SendSMSAfricaIsTalking(recipients, message)
}

func SendSMSAfricaIsTalking(recipients []string, message string) error {
	
	var apiKey string = configs.EnvConfigs.AfricasTalkingAPIKey
	var username string = configs.EnvConfigs.AfricasTalkingUsername
	var endPoint string = "https://api.africastalking.com/version1/messaging"

	//if !apiKey.Valid{
	//	return errors.New("Invalid AfricasTalking ApiKey")
	//}
	//if !username.Valid{
//		return errors.New("Invalid AfricasTalking Username")

//	}
	values := url.Values{}
	values.Set("username", username)
	values.Set("to", strings.Join(recipients, ","))
	values.Set("message", message)

	reader := strings.NewReader(values.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", endPoint, reader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("apikey", apiKey)
	req.Header.Add("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		bodyString := string(body)
		logger.Log(TAG, bodyString, logger.LOG_LEVEL_INFO)
	}

	return nil
}

